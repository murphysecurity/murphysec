package maven

import "sync"

type pomCache struct {
	cachedPom   map[Coordinate]*UnresolvedPom
	cachedError map[Coordinate]error
	mutex       sync.Mutex
}

func newPomCache() *pomCache {
	return &pomCache{
		cachedPom:   map[Coordinate]*UnresolvedPom{},
		cachedError: map[Coordinate]error{},
		mutex:       sync.Mutex{},
	}
}

func (r *pomCache) add(pom *UnresolvedPom) {
	r.write(pom.Coordinate(), pom, nil)
}

func (r *pomCache) fetch(coordinate Coordinate) (*UnresolvedPom, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	p, ok := r.cachedPom[coordinate]
	if ok {
		return p, nil
	}
	e, ok := r.cachedError[coordinate]
	if ok {
		return nil, e
	}
	return nil, nil
}

func (r *pomCache) write(coordinate Coordinate, pom *UnresolvedPom, err error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	if pom != nil {
		r.cachedPom[coordinate] = pom
		return
	}
	r.cachedError[coordinate] = err
}
