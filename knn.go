package knn

import (
	"github.com/intdxdt/mbr"
	"github.com/intdxdt/rtree"
	"github.com/TopoSimplify/ctx"
	"github.com/TopoSimplify/box"
	"github.com/TopoSimplify/node"
	"github.com/intdxdt/geom"
)

const EpsilonDist = 1.0e-5

//find knn
func Find(database *rtree.RTree, g geom.Geometry, dist float64,
	score func(*mbr.MBR, *rtree.KObj) float64,
	predicate ... func(*rtree.KObj) (bool, bool)) []*rtree.Obj {

	var fn func(*rtree.KObj) (bool, bool)
	if len(predicate) > 0 {
		fn = predicate[0]
	} else {
		fn = PredicateFn(dist)
	}
	return database.Knn(g.Bounds(), -1, score, fn)
}

//score function
func ScoreFn(query geom.Geometry) func(_ *mbr.MBR, item *rtree.KObj) float64 {
	return func(_ *mbr.MBR, item *rtree.KObj) float64 {
		var ok bool
		var mb *mbr.MBR
		var other geom.Geometry
		//item is box from rtree
		if mb, ok = item.GetItem().Object.(*mbr.MBR); ok {
			other = box.MBRToPolygon(*mb)
		} else { //item is either ctxgeom or node.Node
			if item.GetItem().Object == nil {
				other = box.MBRToPolygon(*item.MBR)
			} else if o, ok := item.GetItem().Object.(*node.Node); ok {
				other = o.Geometry
			} else if o, ok := item.GetItem().Object.(*ctx.ContextGeometry); ok {
				other = o.Geom
			} else {
				panic("unimplemented !")
			}
		}
		return query.Distance(other)
	}
}

//predicate function
func PredicateFn(dist float64) func(*rtree.KObj) (bool, bool) {
	return func(candidate *rtree.KObj) (bool, bool) {
		if candidate.Dist <= dist {
			return true, false
		}
		return false, true
	}
}
