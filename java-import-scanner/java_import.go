package java_import_scanner

import (
	"container/list"
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"io/ioutil"
	"murphysec-cli-simple/java-import-scanner/javaparser"
	"murphysec-cli-simple/logger"
	"path/filepath"
	"strings"
	"sync"
)

type JavaImportListener struct {
	javaparser.BaseJavaParserListener
	Imports []string
}

func (s *JavaImportListener) EnterImportDeclaration(ctx *javaparser.ImportDeclarationContext) {
	s.Imports = append(s.Imports, ctx.QualifiedName().GetText())
}

type JavaFileImportItem struct {
	FilePath string
	Imports  []string
}

// JavaImportScan recursive scan dir, write all java import info to rsCh.
func JavaImportScan(dir string, rsCh chan JavaFileImportItem) {
	fileListCh := make(chan string, 100)
	go func() {
		dirJavaFileScan(dir, fileListCh)
		close(fileListCh)
	}()
	wg := sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			filePath := <-fileListCh
			for filePath != "" {
				imports := parseJavaFileImport(filePath)
				if len(imports) == 0 {
					continue
				}
				rsCh <- JavaFileImportItem{
					FilePath: filePath,
					Imports:  imports,
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func parseJavaFileImport(path string) []string {
	data, e := ioutil.ReadFile(path)
	if e != nil {
		logger.Err.Println("Read java file error:", e.Error())
		return nil
	}
	listener := &JavaImportListener{}
	parser := javaparser.NewJavaParser(antlr.NewCommonTokenStream(javaparser.NewJavaLexer(antlr.NewInputStream(string(data))), antlr.TokenDefaultChannel))
	parser.RemoveErrorListeners()
	walker := new(antlr.ParseTreeWalker)
	walker.Walk(listener, parser.CompilationUnit())
	return listener.Imports
}

// recursive scan all *.java files, write path to resultCh.
func dirJavaFileScan(dir string, resultCh chan string) {
	visited := map[string]struct{}{}
	q := list.New()
	q.PushBack(dir)
	for q.Len() > 0 {
		curr := q.Front().Value.(string)
		q.Remove(q.Front())
		if _, ok := visited[curr]; ok {
			continue
		}
		visited[curr] = struct{}{}
		d, e := ioutil.ReadDir(curr)
		if e != nil {
			logger.Err.Println("Read dir failed.", curr, e.Error())
			continue
		}
		for _, it := range d {
			if dirScanBlackList[it.Name()] {
				continue
			}
			if it.IsDir() {
				q.PushBack(filepath.Join(curr, it.Name()))
				continue
			}
			if strings.HasSuffix(it.Name(), ".java") {
				resultCh <- filepath.Join(curr, it.Name())
			}
		}
	}
}

var dirScanBlackList = map[string]bool{
	"node_modules": true,
	".git":         true,
	".gradle":      true,
	".mvn":         true,
	".m2":          true,
	"gradle":       true,
	".idea":        true,
}
