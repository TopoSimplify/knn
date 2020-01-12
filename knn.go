package knn

import (
	"github.com/TopoSimplify/box"
	"github.com/TopoSimplify/hdb"
	"github.com/TopoSimplify/node"
	"github.com/intdxdt/geom"
	"github.com/intdxdt/mbr"
)

const EpsilonDist = 1.0e-5

//find knn
func find(database *hdb.Hdb, g geom.Geometry, dist float64,
	score func(*mbr.MBR, *hdb.KObj) float64,
	predicate ...func(*hdb.KObj) (bool, bool)) []*node.Node {

	var fn func(*hdb.KObj) (bool, bool)

	if len(predicate) > 0 {
		fn = predicate[0]
	} else {
		fn = PredicateFn(dist)
	}

	return database.Knn(g.Bounds(), -1, score, fn)
}

//score function
func ScoreFn(query geom.Geometry) func(_ *mbr.MBR, item *hdb.KObj) float64 {
	return func(_ *mbr.MBR, item *hdb.KObj) float64 {
		var other geom.Geometry
		//item is box from rtree
		var nd = item.GetNode()
		if nd == nil {
			other = box.MBRToPolygon(*item.MBR)
		} else {
			other = nd.Geom
		}
		return query.Distance(other)
	}
}

//predicate function
func PredicateFn(dist float64) func(*hdb.KObj) (bool, bool) {
	return func(candidate *hdb.KObj) (bool, bool) {
		if candidate.Distance <= dist {
			return true, false
		}
		return false, true
	}
}
