package luarocks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/luarocks/parser"
	"github.com/murphysecurity/murphysec/utils"
	"github.com/repeale/fp-go"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ParsingErrors struct {
	Total  int
	Errors []ParsingErrorItem
}

func (p ParsingErrors) Error() string {
	var r = []string{"total " + strconv.Itoa(p.Total) + " errors:"}
	var prefix = func(s string) string { return "  " + s }
	r = append(r, fp.Pipe2(fp.Map(ParsingErrorItem.String), fp.Map(prefix))(p.Errors)...)
	return strings.Join(r, "\n")
}

type ParsingErrorItem struct {
	Line    int
	Column  int
	Message string
}

func (p ParsingErrorItem) String() string {
	return strconv.Itoa(p.Line) + ":" + strconv.Itoa(p.Column) + " " + p.Message
}

func analyze(ctx context.Context, tree parser.IStart_Context) ([]model.DependencyItem, []model.DependencyItem) {
	var chunk = tree.Chunk()
	if chunk == nil {
		return nil, nil
	}
	var block = chunk.Block()
	if block == nil {
		return nil, nil
	}
	var ret = block.Retstat()
	if ret == nil {
		return nil, nil
	}
	var epList = ret.Explist()
	if epList == nil {
		return nil, nil
	}
	var expContent = epList.Exp(0)
	if expContent == nil {
		return nil, nil
	}
	var tableConstruct = expContent.Tableconstructor()
	if tableConstruct == nil {
		return nil, nil
	}
	var fieldList = tableConstruct.Fieldlist()
	if fieldList == nil {
		return nil, nil
	}
	var deps []model.DependencyItem
	var testDeps []model.DependencyItem
	for _, it := range fieldList.AllField() {
		var name = it.NAME()
		if name == nil {
			continue
		}
		var nameText = name.GetText()
		var exp = it.Exp(0)
		if exp == nil {
			continue
		}
		var tc = exp.Tableconstructor()
		if tc == nil {
			continue
		}
		switch nameText {
		case "dependencies":
			deps = analyzeTc(ctx, tc)
		case "test_dependencies":
			testDeps = analyzeTc(ctx, tc)
		}
	}
	for i := range testDeps {
		testDeps[i].IsOnline.SetOnline(false)
	}
	return deps, testDeps
}

func analyzeTc(ctx context.Context, tc parser.ITableconstructorContext) []model.DependencyItem {
	var deps []model.DependencyItem
	var fl = tc.Fieldlist()
	if fl == nil {
		return nil
	}
	for _, fieldContext := range fl.AllField() {
		var name, value = analyzeField(ctx, fieldContext)
		if name == "" || value == "" {
			continue
		}
		deps = append(deps, model.DependencyItem{
			Component: model.Component{
				CompName:    name,
				CompVersion: value,
				EcoRepo:     _EcoRepo,
			},
			IsDirectDependency: true,
			IsOnline:           model.IsOnlineTrue(),
		})
	}
	return deps
}

var _EcoRepo = model.EcoRepo{
	Ecosystem:  "luarocks",
	Repository: "",
}

func analyzeField(ctx context.Context, field parser.IFieldContext) (string, string) {
	var name = analyzeFieldName(field)
	var value = analyzeFieldValue(field)
	return name, value
}

func analyzeFieldValue(field parser.IFieldContext) (r string) {
	var expIdx int
	if field.NAME() == nil {
		expIdx = 1
	}
	var vExp = field.Exp(expIdx)
	if vExp == nil {
		return
	}
	var strc = vExp.String_()
	if strc == nil {
		return
	}
	return deEscapeLuaText(strc.GetText())
}

func analyzeFieldName(field parser.IFieldContext) string {
	var nameNode = field.NAME()
	var name string
	if nameNode != nil {
		name = nameNode.GetText()
	} else {
		// maybe string array
		var exp = field.Exp(0)
		if exp != nil {
			var strc = exp.String_()
			if strc != nil {
				name = deEscapeLuaText(strc.GetText())
			}
		}
	}
	return name
}

func deEscapeLuaText(text string) (r string) {
	_ = json.Unmarshal([]byte(text), &r)
	return
}

func parse(input string) (parser.IStart_Context, error) {
	var listener recListener
	var lexer = parser.NewLuaLexer(antlr.NewInputStream(input))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(&listener)
	var cts = antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	var p = parser.NewLuaParser(cts)
	p.RemoveErrorListeners()
	p.AddErrorListener(&listener)
	var r = p.Start_()
	if len(listener.errors) > 0 {
		return nil, &ParsingErrors{
			Total:  len(listener.errors),
			Errors: listener.errors,
		}
	}
	return r, nil
}

type recListener struct {
	count int
	antlr.DefaultErrorListener
	errors []ParsingErrorItem
}

func (r *recListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line int, col int, msg string, _ antlr.RecognitionException) {
	r.count++
	if len(r.errors) > 10 {
		return
	}
	r.errors = append(r.errors, ParsingErrorItem{
		Line:    line,
		Column:  col,
		Message: msg,
	})
}

const LockfileName = "luarocks.lock"

type Inspector struct {
}

func (Inspector) String() string {
	return "LuaRocks"
}

func (Inspector) CheckDir(dir string) bool {
	return utils.IsFile(filepath.Join(dir, LockfileName))
}

func (Inspector) InspectProject(ctx context.Context) error {
	var it = model.UseInspectionTask(ctx)
	var lfPath = filepath.Join(it.Dir(), LockfileName)
	data, e := os.ReadFile(lfPath)
	if e != nil {
		return nil
	}
	tree, e := parse(string(data))
	if e != nil {
		return fmt.Errorf("parse: %w", e)
	}
	var deps, testDeps = analyze(ctx, tree)
	if len(deps) != 0 || len(testDeps) == 0 {
		return nil
	}
	var module = model.Module{
		ModulePath:     lfPath,
		PackageManager: "LuaRocks",
		ScanStrategy:   model.ScanStrategyNormal,
	}
	module.Dependencies = append(module.Dependencies, deps...)
	module.Dependencies = append(module.Dependencies, testDeps...)
	it.AddModule(module)
	return nil
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	return false
}

var _ model.Inspector = (*Inspector)(nil)
