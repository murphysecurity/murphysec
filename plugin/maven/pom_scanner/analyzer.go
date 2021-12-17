package pom_scanner

func Analyze(dir string) (*ResolvedDependency, error) {
	pom, err := scanPom(dir, map[string]bool{})
	if err != nil {
		return nil, err
	}
	dependency, err := resolve(pom)
	if err != nil {
		return nil, err
	}
	return dependency, nil
}
