package shared

type Visited struct {
	name    string
	version string
	parent  *Visited
}

func (v *Visited) Visited(name, version string) bool {
	var curr = v
	for curr != nil {
		if curr.name == name && curr.version == version {
			return true
		}
		curr = curr.parent
	}
	return false
}

func (v *Visited) AddOrNil(name, version string) *Visited {
	if name == "" {
		panic("name is empty")
	}
	if v != nil && v.Visited(name, version) {
		return nil
	}
	return &Visited{
		name:    name,
		version: version,
		parent:  v,
	}
}
