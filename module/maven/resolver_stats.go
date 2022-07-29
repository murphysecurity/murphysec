package maven

type resolverStats struct {
	totalReq int
	cacheHit int
}

func newResolverStats() *resolverStats {
	return &resolverStats{}
}
