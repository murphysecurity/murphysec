package shared

type Component struct {
	Name     string
	Version  string
	Dev      bool
	Children []*Component
}
