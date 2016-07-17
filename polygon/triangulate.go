// Package polygon implements utility functions for dealing with
// (simple) polygons, most notably triangulation
package polygon

import (
	"image"
	"log"
)

type IndexedTriangle [3]int

type node struct {
	prev, next int // doubly-linked-list pointers
	index      int
	vertex     image.Point // save some indirection
	reflex bool	 // is it a reflex vertex?
}

// see:
// https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf
//
func Triangulate(vertices []image.Point) []IndexedTriangle {
	n := len(vertices)

	// 1. find all the reflex vertices, since those are the ones
	// that we have to check for containment.  Note that for convex
	// polygons, there are no reflex vertices
	eartips := make([]int, 0, n)
	reflex := []int{}
	for i, p := range vertices {
		pred := (i + n - 1) % n
		succ := (i + 1) % n
		if isReflex(vertices[pred], p, vertices[succ]) {
			log.Printf("v[%d] is reflex", i)
			reflex = append(reflex, i)
		} else {
			log.Printf("v[%d] is a candidate", i)
			eartips = append(eartips, i)
		}
	}

	// 2. discard candidates that are not ear tips by checking to see
	// if any reflex vertex is contained in it
	for i := 0; i < len(eartips); {
		contained := false
		e := eartips[i]
		pred := (e + n - 1) % n
		succ := (e + 1) % n

		for _, r := range reflex {
			if r == pred || r == succ {
				// reflex vertices that are part of this
				// triangle don't count
				continue
			}
			contained = isContained(
				vertices[r],
				vertices[pred],
				vertices[e],
				vertices[succ])
			if contained {
				break
			}
		}
		if contained {
			// discard
			ne := len(eartips)
			eartips[i] = eartips[ne-1]
			eartips = eartips[:ne-1]
		} else {
			log.Printf("v[%d] is definitely an ear tip", e)
			// continue
			i++
		}
	}

	return nil
}

func isReflex(a, b, c image.Point) bool {
	u := b.Sub(a)
	v := c.Sub(b)
	// rotate u 90 degrees counterclockwise, so the dot product
	// is positive if the angle bends around to the right
	u.X, u.Y = u.Y, -u.X
	dot := u.X*v.X + u.Y*v.Y
	return dot > 0
}

func neg(a, b, c image.Point) bool {
	return ((a.X-c.X)*(b.Y-c.Y) - (b.X-c.X)*(a.Y-c.Y)) < 0
}

// check to see if the point p is contained in the
// triangle a, b, c
func isContained(p, a, b, c image.Point) bool {
	b1 := neg(p, a, b)
	b2 := neg(p, b, c)
	b3 := neg(p, c, a)
	return (b1 == b2) && (b2 == b3)
}
