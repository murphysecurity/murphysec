package luarocks

import (
	"context"
	"fmt"
	"github.com/antlr4-go/antlr/v4"
	"github.com/murphysecurity/murphysec/model"
	"github.com/murphysecurity/murphysec/module/luarocks/parser"
)

func parse(input string) parser.IStart_Context {
	var errlistener = &recListener{}
	var lexer = parser.NewLuaLexer(antlr.NewInputStream(input))
	lexer.RemoveErrorListeners()
	lexer.AddErrorListener(errlistener)
	var cts = antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	var p = parser.NewLuaParser(cts)
	p.RemoveErrorListeners()
	p.AddErrorListener(errlistener)
	return p.Start_()
}

type recListener struct {
	antlr.DefaultErrorListener
	errors []string
}

func (r *recListener) SyntaxError(_ antlr.Recognizer, _ interface{}, line int, col int, msg string, _ antlr.RecognitionException) {
	r.errors = append(r.errors, fmt.Sprintf("%d:%d %s", line, col, msg))
}

type Inspector struct {
}

func (Inspector) String() string {
	//TODO implement me
	panic("implement me")
}

func (Inspector) CheckDir(dir string) bool {
	//TODO implement me
	panic("implement me")
}

func (Inspector) InspectProject(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (Inspector) SupportFeature(feature model.InspectorFeature) bool {
	//TODO implement me
	panic("implement me")
}

var _ model.Inspector = (*Inspector)(nil)
