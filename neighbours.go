package knn

import (
	"github.com/intdxdt/geom"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/hdb"
)

//find context neighbours by a certain distance
func ContextNeighbours(database *hdb.Hdb, query geom.Geometry, dist float64) []*node.Node {
	return find(database, query, dist, ScoreFn(query))
}

//find context hulls
func NodeNeighbours(database *hdb.Hdb, hull *node.Node, dist float64) []*node.Node {
	//var ns = find(database, hull.Geom, dist, ScoreFn(hull.Geom), NodePredicateFn(hull, dist))
	var inters = database.Search(*hull.BBox())
	var ns = make([]*node.Node, 0, len(inters))
	for _, nd := range inters {
		if nd.Id != hull.Id && nd.Geom.Distance(hull.Geom) <= dist {
			ns = append(ns, nd)
		}
	}
	return ns
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
		if query.Geom.Intersects(candhull.Geom) || (candidate.Distance <= dist) {
			return true, false
		}
		return false, true
	}
}
