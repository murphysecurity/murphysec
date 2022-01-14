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
	"time"
)

func PkgPrefixScan(dir string, pkgPrefixes []string) map[string][]string {
	prefixSet := map[string]struct{}{}
	for _, it := range pkgPrefixes {
		prefixSet[it] = struct{}{}
	}
	filepathCh := make(chan string, 50)
	// 文件读取
	go func() {
		// 去重，防环
		visited := map[string]struct{}{}
		// 广度优先遍历目录树
		q := list.New()
		q.PushBack(dir)

		for q.Len() > 0 {
			curr := q.Front().Value.(string)
			q.Remove(q.Front())
			if _, ok := visited[curr]; ok {
				continue
			}
			visited[curr] = struct{}{}
			flist, e := ioutil.ReadDir(curr)
			if e != nil {
				logger.Err.Println("ReadDir failed,", curr, ".", e.Error())
				continue
			}
			for _, it := range flist {
				if it.IsDir() {
					q.PushBack(filepath.Join(curr, it.Name()))
					continue
				}
				if strings.HasSuffix(it.Name(), ".java") {
					filepathCh <- filepath.Join(curr, it.Name())
				}
			}
		}
		close(filepathCh)
	}()
	type pkgImport struct {
		pkg      string
		filePath string
	}
	// parsing
	parserWg := sync.WaitGroup{}
	commCh := make(chan pkgImport, 100)
	for i := 0; i < 2; i++ {
		parserWg.Add(1)
		go func() {
			for {
				filePath := <-filepathCh
				if filePath == "" {
					break
				}
				imports := parseJavaFileImport(filePath)
				for prefix := range prefixSet { // todo: optimize
					for _, it := range imports {
						if strings.HasPrefix(it, prefix) {
							commCh <- pkgImport{pkg: prefix, filePath: filePath}
						}
					}
				}
			}
			parserWg.Done()
		}()
	}
	// 统计
	summary := make(chan map[string]map[string]struct{}, 1)
	go func() {
		m := map[string]map[string]struct{}{}
		for {
			p := <-commCh
			if p.pkg == "" {
				break
			}
			if m[p.pkg] == nil {
				m[p.pkg] = map[string]struct{}{}
			}
			m[p.pkg][p.filePath] = struct{}{}
		}
		summary <- m
	}()
	parserWg.Wait()
	close(commCh)
	logger.Info.Println("Finish file parse.")
	// result
	rs := map[string][]string{}
	for pkg, files := range <-summary {
		for it := range files {
			rs[pkg] = append(rs[pkg], it)
		}
	}
	return rs
}

type JavaImportListener struct {
	javaparser.BaseJavaParserListener
	Imports []string
}

func (s *JavaImportListener) EnterImportDeclaration(ctx *javaparser.ImportDeclarationContext) {
	s.Imports = append(s.Imports, ctx.QualifiedName().GetText())
}

func parseJavaFileImport(path string) []string {
	data, e := ioutil.ReadFile(path)
	if e != nil {
		logger.Err.Println("Read java file error:", e.Error())
		return nil
	}
	st := time.Now()
	logger.Debug.Println("Parse java file:", path)
	listener := &JavaImportListener{}
	parser := javaparser.NewJavaParser(antlr.NewCommonTokenStream(javaparser.NewJavaLexer(antlr.NewInputStream(string(data))), antlr.TokenDefaultChannel))
	walker := new(antlr.ParseTreeWalker)
	walker.Walk(listener, parser.CompilationUnit())
	logger.Debug.Println("Parse done,", path, ", time:", time.Now().Sub(st).Microseconds())
	return listener.Imports
}
