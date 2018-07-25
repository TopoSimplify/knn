package knn

import (
	"github.com/intdxdt/geom"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/hdb"
)

//find context neighbours
func FindNeighbours(database *hdb.Hdb, query geom.Geometry, dist float64) []*node.Node {
	return Find(database, query, dist, ScoreFn(query))
}

//find context hulls
func FindNodeNeighbours(database *hdb.Hdb, hull *node.Node, dist float64) []*node.Node {
	return Find(database, hull.Geometry, dist, ScoreFn(hull.Geometry), NodePredicateFn(hull, dist))
}

//hull predicate within index range i, j.
func NodePredicateFn(query *node.Node, dist float64) func(*hdb.KObj) (bool, bool) {
	//@formatter:off
	return func(candidate *hdb.KObj) (bool, bool) {
		var candhull = candidate.GetNode()

		// same hull
		if candhull.Range.Equals(query.Range) {
			return false, false
		}

		// if intersects or distance from context neighbours is within dist
		if query.Geometry.Intersects(candhull.Geometry) || (candidate.Dist <= dist) {
			return true, false
		}
		return false, true
	}
}
