package knn

import (
	"simplex/igeom"
	"simplex/node"
	"github.com/intdxdt/rtree"
)

//find context neighbours
func FindNeighbours(database *rtree.RTree, query igeom.IGeom, dist float64) []rtree.BoxObj {
	return Find(database, query.Geometry(), dist, ScoreFn(query))
}

//find context hulls
func FindNodeNeighbours(hulldb *rtree.RTree, hull *node.Node, dist float64) []rtree.BoxObj {
	return Find(hulldb, hull.Geometry(), dist, ScoreFn(hull), NodePredicateFn(hull, dist))
}

//hull predicate within index range i, j.
func NodePredicateFn(query *node.Node, dist float64) func(*rtree.KObj) (bool, bool) {
	//@formatter:off
	return func(candidate *rtree.KObj) (bool, bool) {
		var candhull = candidate.GetItem().(*node.Node)
		var qgeom    = query.Geom
		var cgeom    = candhull.Geom

		// same hull
		if candhull.Range.Equals(query.Range) {
			return false, false
		}

		// if intersects or distance from context neighbours is within dist
		if qgeom.Intersects(cgeom) || (candidate.Score() <= dist) {
			return true, false
		}
		return false, true
	}
}