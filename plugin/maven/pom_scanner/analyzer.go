package pom_scanner

func Analyze(dir string) (map[string]interface{}, error) {
	pom, err := scanPom(dir, map[string]bool{})
	if err != nil {
		return nil, err
	}
	dependency, err := resolve(pom)
	if err != nil {
		return nil, err
	}
	return _convMap(*dependency), nil
}

func _convMap(rd ResolvedDependency) map[string]interface{} {
	dm := map[string]interface{}{}
	for _, it := range rd.Dependencies {
		dm[it.Name] = _convMap(it)
	}
	return map[string]interface{}{
		"name":         rd.Name,
		"version":      rd.Version,
		"dependencies": dm,
	}
}
