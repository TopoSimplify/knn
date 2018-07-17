package knn

import (
	"github.com/TopoSimplify/node"
	"github.com/intdxdt/rtree"
	"github.com/intdxdt/geom"
)

//find context neighbours
func FindNeighbours(database *rtree.RTree, query geom.Geometry, dist float64) []*rtree.Obj {
	return Find(database, query, dist, ScoreFn(query))
}

//find context hulls
func FindNodeNeighbours(database *rtree.RTree, hull *node.Node, dist float64) []*rtree.Obj {
	return Find(database, hull.Geometry, dist, ScoreFn(hull.Geometry), NodePredicateFn(hull, dist))
}

//hull predicate within index range i, j.
func NodePredicateFn(query *node.Node, dist float64) func(*rtree.KObj) (bool, bool) {
	//@formatter:off
	return func(candidate *rtree.KObj) (bool, bool) {
		var candhull = candidate.GetItem().Object.(*node.Node)

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
