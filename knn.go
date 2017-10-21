package knn

import (
	"github.com/intdxdt/rtree"
	"simplex/igeom"
	"github.com/intdxdt/mbr"
	"simplex/ctx"
	"simplex/box"
	"simplex/node"
)

//find knn
func Find(db *rtree.RTree, g rtree.BoxObj, dist float64,
	score func(rtree.BoxObj, rtree.BoxObj) float64,
	predicate ... func(*rtree.KObj) (bool, bool)) []rtree.BoxObj {

	var pred func(*rtree.KObj) (bool, bool)
	if len(predicate) > 0 {
		pred = predicate[0]
	} else {
		pred = PredicateFn(dist)
	}

	return db.KNN(g, -1, score, pred)
}

//score function
func ScoreFn(query igeom.IGeom) func(_, item rtree.BoxObj) float64 {
	return func(_, item rtree.BoxObj) float64 {
		var ok bool
		var mb *mbr.MBR
		var other igeom.IGeom
		//item is box from rtree
		if mb, ok = item.(*mbr.MBR); ok {
			other = box.MBRToPolygon(mb)
		} else { //item is either ctxgeom or node.Node
			if other, ok = item.(*ctx.CtxGeom); !ok {
				other = item.(*node.Node)
			}
		}
		return query.Geometry().Distance(other.Geometry())
	}
}

//predicate function
func PredicateFn(dist float64) func(*rtree.KObj) (bool, bool) {
	return func(candidate *rtree.KObj) (bool, bool) {
		if candidate.Score() <= dist {
			return true, false
		}
		return false, true
	}
}