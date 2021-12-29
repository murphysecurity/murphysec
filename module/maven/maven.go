package maven

type Dependency struct {
	Coordination
	Children []Dependency
}
