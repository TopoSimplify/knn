package knn

import (
	"time"
	"testing"
	"github.com/intdxdt/geom"
	"github.com/franela/goblin"
	"github.com/TopoSimplify/node"
	"github.com/TopoSimplify/pln"
	"github.com/TopoSimplify/rng"
	"github.com/TopoSimplify/dp"
	"github.com/TopoSimplify/hdb"
	"github.com/intdxdt/mbr"
)

type iG struct{ g geom.Geometry }

func (o *iG) Geometry() geom.Geometry {
	return o.g
}

func linearCoords(wkt string) []geom.Point {
	return geom.NewLineStringFromWKT(wkt).Coordinates()
}

func createNodes(indxs [][]int, coords []geom.Point) []*node.Node {
	var poly = pln.CreatePolyline(coords)
	var hulls = make([]*node.Node, 0, len(indxs))
	for i := range indxs {
		var r = rng.Range(indxs[i][0], indxs[i][1])
		hulls = append(hulls, node.CreateNode(poly.SubCoordinates(r), r, dp.NodeGeometry))
	}
	return hulls
}

func TestDB(t *testing.T) {
	var g    = goblin.Goblin(t)
	var wkts = []string{
		"POINT ( 190 310 )", "POINT ( 220 400 )", "POINT ( 260 200 )", "POINT ( 260 340 )",
		"POINT ( 260 290 )", "POINT ( 310 280 )", "POINT ( 350 250 )", "POINT ( 350 330 )",
		"POINT ( 380 370 )", "POINT ( 400 240 )", "POINT ( 410 310 )",
		"POLYGON (( 160 340, 160 380, 180 380, 180 340, 160 340 ))",
		"POLYGON (( 180 240, 180 280, 210 280, 210 240, 180 240 ))",
		"POLYGON (( 280 370, 280 400, 300 400, 300 370, 280 370 ))",
		"POLYGON (( 340 210, 340 230, 360 230, 360 210, 340 210 ))",
		"POLYGON (( 410 340, 410 430, 420 430, 420 340, 410 340 ))",
	}
	g.Describe("rtree knn", func() {
		var scoreFn = func(q *mbr.MBR, item *hdb.KObj) float64 {
			return q.Distance(item.MBR)
		}

		g.It("should test k nearest neighbour", func() {
			var objects = make([]*node.Node, 0)
			for i := range wkts {
				var g = geom.ReadGeometry(wkts[i])
				objects = append(objects, &node.Node{Id: i, MBR: g.Bounds(), Geom: g})
			}
			var tree = hdb.NewHdb()
			tree.Load(objects)
			var q = geom.ReadGeometry("POLYGON (( 370 300, 370 330, 400 330, 400 300, 370 300 ))")
			var results = Find(tree, q, 15, scoreFn)

			g.Assert(len(results) == 2)
			results = Find(tree, q, 20, scoreFn)
			g.Assert(len(results) == 3)
		})

		g.It("should test k nearest node neighbour", func() {
			g.Timeout(1 * time.Hour)

			var coords = linearCoords("LINESTRING ( 780 600, 740 620, 720 660, 720 700, 760 740, 820 760, 860 740, 880 720, 900 700, 880 660, 840 680, 820 700, 800 720, 760 700, 780 660, 820 640, 840 620, 860 580, 880 620, 820 660 )")
			var hulls = createNodes([][]int{{0, 3}, {3, 8}, {8, 13}, {13, 17}, {17, len(coords) - 1}}, coords)
			var tree = hdb.NewHdb(2)
			for i := range hulls {
				tree.insert(hulls[i])
			}

			var q = hulls[0]
			var vs = FindNeighbours(tree, q.Geom, 0)
			g.Assert(len(vs)).Equal(2)
			vs = FindNodeNeighbours(tree, q, 0)
			g.Assert(len(vs)).Equal(1)
		})
	})
}
