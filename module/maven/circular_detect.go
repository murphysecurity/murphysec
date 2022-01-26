package maven

type CoordinateNode struct {
	Coordinate
	Prev *CoordinateNode
}

func (c *CoordinateNode) Append(coordinate Coordinate) *CoordinateNode {
	return &CoordinateNode{
		Coordinate: coordinate,
		Prev:       c,
	}
}

func (c *CoordinateNode) Has(coordinate Coordinate) bool {
	if c == nil {
		return false
	}
	p := c.Prev
	for p != nil {
		if p.Coordinate == coordinate {
			return true
		}
		p = p.Prev
	}
	return false
}
