package knn

import (
    "time"
    "testing"
    "simplex/node"
    "simplex/pln"
    "simplex/rng"
    "simplex/dp"
    "simplex/box"
    "github.com/intdxdt/mbr"
    "github.com/intdxdt/geom"
    "github.com/intdxdt/rtree"
    "github.com/franela/goblin"
)

type iG struct{ g geom.Geometry }

func (o *iG) Geometry() geom.Geometry {
    return o.g
}

func linearCoords(wkt string) []*geom.Point {
    return geom.NewLineStringFromWKT(wkt).Coordinates()
}

func createNodes(indxs [][]int, coords []*geom.Point) []*node.Node {
    poly := pln.New(coords)
    hulls := make([]*node.Node, 0)
    for _, o := range indxs {
        hulls = append(hulls, node.NewFromPolyline(poly, rng.NewRange(o[0], o[1]), dp.NodeGeometry))
    }
    return hulls
}

func TestDB(t *testing.T) {
    g := goblin.Goblin(t)
    wkts := []string{
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
        score_fn := func(q, item rtree.BoxObj) float64 {
            g := q.(geom.Geometry)
            var other geom.Geometry
            if o, ok := item.(*mbr.MBR); ok {
                other = box.MBRToPolygon(o)
            } else {
                other = item.(geom.Geometry)
            }
            return g.Distance(other)
        }
        g.It("should test k nearest neighbour", func() {
            objs := make([]rtree.BoxObj, 0)
            for _, wkt := range wkts {
                objs = append(objs, geom.NewGeometry(wkt))
            }
            tree := rtree.NewRTree(8)
            tree.Load(objs)
            q := geom.NewGeometry("POLYGON (( 370 300, 370 330, 400 330, 400 300, 370 300 ))")

            results := Find(tree, q, 15, score_fn)

            g.Assert(len(results) == 2)
            results = Find(tree, q, 20, score_fn)
            g.Assert(len(results) == 3)
        })

        g.It("should test k nearest node neighbour", func() {
            g.Timeout(1 * time.Hour)

            var coords = linearCoords("LINESTRING ( 780 600, 740 620, 720 660, 720 700, 760 740, 820 760, 860 740, 880 720, 900 700, 880 660, 840 680, 820 700, 800 720, 760 700, 780 660, 820 640, 840 620, 860 580, 880 620, 820 660 )")
            var hulls = createNodes([][]int{{0, 3}, {3, 8}, {8, 13}, {13, 17}, {17, len(coords) - 1}}, coords)
            tree := rtree.NewRTree(2)
            for _, h := range hulls {
                tree.Insert(h)
            }
            var q = hulls[0]
            var vs = FindNeighbours(tree, q, 0)
            g.Assert(len(vs)).Equal(2)
            vs = FindNodeNeighbours(tree, q, 0)
            g.Assert(len(vs)).Equal(1)
        })
    })
}
