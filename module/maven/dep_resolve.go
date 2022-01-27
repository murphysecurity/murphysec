package maven

import (
	"context"
	"fmt"
	"murphysec-cli-simple/logger"
	"sync"
)

type depTree struct {
	coordinate Coordinate
	children   []*depTree
}

func _resolve(ctx context.Context, resolver *Resolver, pom *PomFile, cache *DepTreeCacheMap, visited *CoordinateNode, depth int) *Dependency {
	currentPomCoor := pom.Coordinate()
	if depth < 0 || visited.Has(currentPomCoor) {
		return nil
	}
	visited = visited.Append(currentPomCoor)
	depth--
	if v := cache.Get(currentPomCoor); v != nil {
		return v
	}
	node := &Dependency{
		Coordinate: currentPomCoor,
		Children:   []Dependency{},
	}
	wg := sync.WaitGroup{}
	for _, it := range pom.Dependencies() {
		if !it.HasVersion() {
			logger.Debug.Println(fmt.Sprintf("Can't resolve version of %v, skip", it))
			continue
		}
		wg.Add(1)
		it := it
		go func() {
			defer wg.Done()
			p, e := resolver.ResolvePomFile(ctx, it)
			if e != nil {
				logger.Warn.Println("Resolve dependency failed.", it.String(), e)
				return
			}
			n := _resolve(ctx, resolver, p, cache, visited, depth)
			if n != nil {
				node.Children = append(node.Children, *n)
			}
		}()
	}
	wg.Wait()
	cache.Put(currentPomCoor, node)
	return node
}
