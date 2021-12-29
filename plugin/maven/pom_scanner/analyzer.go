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
	return _mapAnalyzeMap(*dependency), nil
}

func _mapAnalyzeMap(rd ResolvedDependency) map[string]interface{} {
	dm := map[string]interface{}{}
	for _, it := range rd.Dependencies {
		r := _mapAnalyzeMap(it)
		dm[it.Name] = r
	}
	m := map[string]interface{}{
		"name":         rd.Name,
		"version":      rd.Version,
		"dependencies": dm,
	}
	return m
}
