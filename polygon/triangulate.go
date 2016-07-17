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

type triangulation struct {
	nodes []node
	start int
}

// see:
// https://www.geometrictools.com/Documentation/TriangulationByEarClipping.pdf
//
func Triangulate(vertices []image.Point) []IndexedTriangle {

	n := len(vertices)
	N := make([]node, n)
	tri := &triangulation{
		nodes: N,
		start: 0,
	}

	// 1. build the circular list of vertices, identifying the
	// reflex vertices in the process (since those are the ones
	// that we have to check for containment.  Note that for
	// convex polygons, there are no reflex vertices)

	for i, p := range vertices {
		pred := (i + n - 1) % n
		succ := (i + 1) % n
		N[i].prev = pred
		N[i].next = succ
		N[i].index = i
		N[i].vertex = p
		N[i].reflex = isReflex(vertices[pred], p, vertices[succ])
		if N[i].reflex {
			log.Printf("v[%d] is reflex", i)
		}
	}
	
	// 2. find ears and carve them off
	for i := 3; i < n; i++ {
		e := tri.findEar()
		ear := N[e]
		log.Printf("*[%d] ear: {%d} %v  {%d}-{%d}-{%d}",
			i, e, ear.vertex,
			ear.prev, e, ear.next)

		// pop it out of the list
		N[ear.next].prev = ear.prev
		N[ear.prev].next = ear.next
		if e == tri.start {
			tri.start = ear.next
		}

		// check to see if adjacent nodes are no longer reflex
		if N[ear.next].reflex {
			i := ear.next
			N[i].reflex = isReflex(
				N[N[i].prev].vertex,
				N[i].vertex,
				N[N[i].next].vertex)
			if !N[i].reflex {
				log.Printf("  de-reflexed {%d}", i)
			}
		}
		if N[ear.prev].reflex {
			i := ear.prev
			N[i].reflex = isReflex(
				N[N[i].prev].vertex,
				N[i].vertex,
				N[N[i].next].vertex)
			if !N[i].reflex {
				log.Printf("  de-reflexed {%d}", i)
			}
		}
	}

	e := tri.start
	ear := N[e]
	log.Printf("*[LAST] {%d}-{%d}-{%d}",
		ear.prev, e, ear.next)
		
	panic("done")
	return nil
}


func (tri *triangulation) findEar() int {

	N := tri.nodes
	for i := tri.start;; i = N[i].next {
		// is node[i] an ear?
		v := N[i]
		if v.reflex {
			continue
		}
		vprev := N[v.prev]
		vnext := N[v.next]
		log.Printf("is {%d} %v an ear?", i, v.vertex)

		contained := false
		for j := tri.start;; j = N[j].next {
			if j != v.prev && j != v.next && N[j].reflex {
				contained = isContained(
					N[j].vertex,
					vprev.vertex,
					v.vertex,
					vnext.vertex)
				if contained {
					log.Printf("nope, contains {%d} %v",
						j, N[j].vertex)
					break
				}
			}
			if tri.nodes[j].next == tri.start {
				break
			}
		}
		if !contained {
			return i
		}
	}

	/*
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
*/
}

func isReflex(a, b, c image.Point) bool {
	u := b.Sub(a)
	v := c.Sub(b)
	// rotate u 90 degrees counterclockwise, so the dot product
	// is positive if the angle bends around to the right
	u.X, u.Y = u.Y, -u.X
	dot := u.X*v.X + u.Y*v.Y
	return dot < 0
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
