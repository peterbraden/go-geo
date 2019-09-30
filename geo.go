// Package geo contains utilities for working with geospatial, such as bounding
// boxes etc.
package geo

import (
	"fmt"
	"github.com/twpayne/go-polyline"
	"math"
)

type BBox struct {
	W, S, E, N float64
}

// Used for Slippy Tile Calculations
// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames
type Tile struct {
	X, Y, Z int
}

func BBoxFromPolyline(s string) (p BBox, empty bool) {
	coords, _, _ := polyline.DecodeCoords([]byte(s))
	if len(coords) > 0 {
		p = BBox{
			W: coords[0][1],
			E: coords[0][1],
			N: coords[0][0],
			S: coords[0][0],
		}
		for _, n := range coords {
			p = p.UnionPoint(n[0], n[1])
		}
		return
	} else {
		empty = true
		return
	}
}

func BBoxFromPolylines(s []string) (p BBox, empty bool) {
	if len(s) > 0 {
		p, _ = BBoxFromPolyline(s[0]) // TODO deal empty
		for _, n := range s {
			if x, empty := BBoxFromPolyline(n); !empty {
				p = p.Union(x)
			}
		}
	} else {
		empty = true
	}
	return
}

func TileFromBBox(b BBox) Tile {
	z := 1.
	t := Deg2Tile(b.N, b.W, z)
	for {
		// Max out at zoom 15
		if t.Contains(b.S, b.E) && z < 15. {
			z++
			t = Deg2Tile(b.N, b.W, z)
		} else {
			return t
		}
	}
}

func Deg2Tile(lat, long, z float64) Tile {
	// From https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames#Lon..2Flat._to_tile_numbers
	e2z := math.Exp2(z)
	x := int(
		math.Floor(
			(long + 180.0) / 360.0 *
				(e2z)))

	y := int(
		math.Floor(
			(1.0 -
				math.Log(
					math.Tan(lat*math.Pi/180.0)+
						1.0/math.Cos(
							lat*math.Pi/180.0))/math.Pi) / 2.0 * (e2z)))

	return Tile{
		x, y, int(z),
	}
}

func Tile2Deg(t Tile) (lat, long float64) {
	n := math.Pi - 2.0*math.Pi*float64(t.Y)/math.Exp2(float64(t.Z))
	r2d := 180 / math.Pi
	lat = r2d * math.Atan(0.5*(math.Exp(n)-math.Exp(-n)))
	long = float64(t.X)/math.Exp2(float64(t.Z))*360.0 - 180.0
	return
}

func TileToBBox(tile Tile) BBox {
	n, w := Tile2Deg(tile)

	next_t := tile
	next_t.X++
	next_t.Y++
	s, e := Tile2Deg(next_t)
	return BBox{w, s, e, n}
}

func (self BBox) BBoxToSQL() string {
	return fmt.Sprintf("BOX(%f %f,%f %f)", self.N, self.E, self.S, self.W)
}

func (self BBox) Union(b BBox) BBox {
	return BBox{
		W: math.Min(self.W, b.W),
		E: math.Max(self.E, b.E),
		N: math.Max(self.N, b.N),
		S: math.Min(self.S, b.S),
	}
}

func (self BBox) UnionPoint(lat, lng float64) BBox {
	return BBox{
		W: math.Min(self.W, lng),
		E: math.Max(self.E, lng),
		N: math.Max(self.N, lat),
		S: math.Min(self.S, lat),
	}
}

func (self BBox) Intersects(b BBox) bool {
	if self.E < b.W {
		return false // Self is to the west of b
	}
	if self.W > b.E {
		return false // Self is to the east of b
	}
	if self.N < b.S {
		return false // Self is to the south of b
	}
	if self.S > b.N {
		return false // Self is to the north of b
	}
	return true
}

func (self BBox) Contains(lat, long float64) bool {
	return lat >= self.S && lat <= self.N && long >= self.W && long <= self.E
}

func (self BBox) ContainsBBox(b BBox) bool {
	return self.Contains(b.N, b.W) && self.Contains(b.S, b.E)
}

func (self Tile) Contains(lat, long float64) bool {
	return TileToBBox(self).Contains(lat, long)
}
