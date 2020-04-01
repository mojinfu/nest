package nest

import (
	"math"

	. "github.com/mojinfu/point"
)

func polygonArea(myPolygon Polygon) float64 {
	var area float64 = 0
	var i, j int
	i = 0
	j = len(myPolygon) - 1
	for i < len(myPolygon) {
		area += (myPolygon[j].X + myPolygon[i].X) * (myPolygon[j].Y - myPolygon[i].Y)
		j = i
		i++

	}
	return 0.5 * area
}

type BoundStruct struct {
	x      float64
	y      float64
	width  float64
	height float64
}

func noFitPolygonRectangle(A, B Polygon) [][]*Point {
	var minAx = A[0].X
	var minAy = A[0].Y
	var maxAx = A[0].X
	var maxAy = A[0].Y
	emptyList := [][]*Point{}
	for i := 1; i < len(A); i++ {
		if A[i].X < minAx {
			minAx = A[i].X
		}
		if A[i].Y < minAy {
			minAy = A[i].Y
		}
		if A[i].X > maxAx {
			maxAx = A[i].X
		}
		if A[i].Y > maxAy {
			maxAy = A[i].Y
		}
	}

	var minBx = B[0].X
	var minBy = B[0].Y
	var maxBx = B[0].X
	var maxBy = B[0].Y
	for i := 1; i < len(B); i++ {
		if B[i].X < minBx {
			minBx = B[i].X
		}
		if B[i].Y < minBy {
			minBy = B[i].Y
		}
		if B[i].X > maxBx {
			maxBx = B[i].X
		}
		if B[i].Y > maxBy {
			maxBy = B[i].Y
		}
	}

	if maxBx-minBx > maxAx-minAx {
		return emptyList
	}
	if maxBy-minBy > maxAy-minAy {
		return emptyList
	}

	return [][]*Point{[]*Point{
		&Point{X: minAx - minBx + B[0].X, Y: minAy - minBy + B[0].Y},
		&Point{X: maxAx - maxBx + B[0].X, Y: minAy - minBy + B[0].Y},
		&Point{X: maxAx - maxBx + B[0].X, Y: maxAy - maxBy + B[0].Y},
		&Point{X: minAx - minBx + B[0].X, Y: maxAy - maxBy + B[0].Y},
	}}
}
func isRectangle(poly Polygon, tolerance float64) bool {
	var bb = getPolygonBounds(poly)
	if tolerance == 0 {
		tolerance = TOL
	}

	for i := 0; i < len(poly); i++ {
		if !_almostEqual2(poly[i].X, bb.x) && !_almostEqual2(poly[i].X, bb.x+bb.width) {
			return false
		}
		if !_almostEqual2(poly[i].Y, bb.y) && !_almostEqual2(poly[i].Y, bb.y+bb.height) {
			return false
		}
	}

	return true
}
func (this *SVG) intersect(oldA, oldB *polygonWithOffset) bool {

	var Aoffsetx = oldA.offsetx
	var Aoffsety = oldA.offsety

	var Boffsetx = oldB.offsetx
	var Boffsety = oldB.offsety

	A := oldA.Polygon[:]
	B := oldB.Polygon[:]

	for i := 0; i < len(A)-1; i++ {
		for j := 0; j < len(B)-1; j++ {
			var a1 = &Point{X: A[i].X + Aoffsetx, Y: A[i].Y + Aoffsety}
			var a2 = &Point{X: A[i+1].X + Aoffsetx, Y: A[i+1].Y + Aoffsety}
			var b1 = &Point{X: B[j].X + Boffsetx, Y: B[j].Y + Boffsety}
			var b2 = &Point{X: B[j+1].X + Boffsetx, Y: B[j+1].Y + Boffsety}
			var prevbindex int
			if j == 0 {
				prevbindex = len(B) - 1
			} else {
				prevbindex = j - 1
			}

			var prevaindex int
			if i == 0 {
				prevaindex = 0
			} else {
				prevaindex = i - 1
			}

			var nextbindex int
			if j+1 == len(B)-1 {
				nextbindex = 0
			} else {
				nextbindex = j + 2
			}

			var nextaindex int
			if i+1 == len(A)-1 {
				nextaindex = 0
			} else {
				nextaindex = i + 2
			}

			// go even further back if we happen to hit on a loop end Point
			if B[prevbindex] == B[j] || (_almostEqual2(B[prevbindex].X, B[j].X) && _almostEqual2(B[prevbindex].Y, B[j].Y)) {

				if prevbindex == 0 {
					prevbindex = len(B) - 1
				} else {
					prevbindex = prevbindex - 1
				}
			}

			if A[prevaindex] == A[i] || (_almostEqual2(A[prevaindex].X, A[i].X) && _almostEqual2(A[prevaindex].Y, A[i].Y)) {

				if prevaindex == 0 {
					prevaindex = len(A) - 1
				} else {
					prevaindex = prevaindex - 1
				}
			}

			// go even further forward if we happen to hit on a loop end Point
			if B[nextbindex] == B[j+1] || (_almostEqual2(B[nextbindex].X, B[j+1].X) && _almostEqual2(B[nextbindex].Y, B[j+1].Y)) {
				if nextbindex == len(B)-1 {
					nextbindex = 0
				} else {
					nextbindex = nextbindex + 1
				}
			}
			if A[nextaindex] == A[i+1] || (_almostEqual2(A[nextaindex].X, A[i+1].X) && _almostEqual2(A[nextaindex].Y, A[i+1].Y)) {
				if nextaindex == len(A)-1 {
					nextaindex = 0
				} else {
					nextaindex = nextaindex + 1
				}
			}

			var a0 = &Point{X: A[prevaindex].X + Aoffsetx, Y: A[prevaindex].Y + Aoffsety}
			var b0 = &Point{X: B[prevbindex].X + Boffsetx, Y: B[prevbindex].Y + Boffsety}

			var a3 = &Point{X: A[nextaindex].X + Aoffsetx, Y: A[nextaindex].Y + Aoffsety}
			var b3 = &Point{X: B[nextbindex].X + Boffsetx, Y: B[nextbindex].Y + Boffsety}

			if _onSegment(a1, a2, b1) || (_almostEqual2(a1.X, b1.X) && _almostEqual2(a1.Y, b1.Y)) {
				// if a Point is on a segment, it could intersect or it could not. Check via the neighboring points
				var b0in = pointInPolygon(b0, A)
				var b2in = pointInPolygon(b2, A)
				if (b0in > 0 && b2in < 0) || (b0in < 0 && b2in > 0) {
					return true
				} else {
					continue
				}
			}

			if _onSegment(a1, a2, b2) || (_almostEqual2(a2.X, b2.X) && _almostEqual2(a2.Y, b2.Y)) {
				// if a Point is on a segment, it could intersect or it could not. Check via the neighboring points
				var b1in = pointInPolygon(b1, A)
				var b3in = pointInPolygon(b3, A)

				if (b1in < 0 && b3in > 0) || (b1in > 0 && b3in < 0) {
					return true
				} else {
					continue
				}
			}

			if _onSegment(b1, b2, a1) || (_almostEqual2(a1.X, b2.X) && _almostEqual2(a1.Y, b2.Y)) {
				// if a Point is on a segment, it could intersect or it could not. Check via the neighboring points
				var a0in = pointInPolygon(a0, B)
				var a2in = pointInPolygon(a2, B)

				if (a0in > 0 && a2in < 0) || (a0in < 0 && a2in > 0) {
					return true
				} else {
					continue
				}
			}

			if _onSegment(b1, b2, a2) || (_almostEqual2(a2.X, b1.X) && _almostEqual2(a2.Y, b1.Y)) {
				// if a Point is on a segment, it could intersect or it could not. Check via the neighboring points
				var a1in = pointInPolygon(a1, B)
				var a3in = pointInPolygon(a3, B)

				if (a1in > 0 && a3in < 0) || (a1in < 0 && a3in > 0) {
					return true
				} else {
					continue
				}
			}

			var p = _lineIntersect(b1, b2, a1, a2, false)

			if p != nil {
				return true
			}
		}
	}

	return false
}
func _lineIntersect(A, B, E, F *Point, infinite bool) *Point {
	var a1, a2, b1, b2, c1, c2, x, y float64

	a1 = B.Y - A.Y
	b1 = A.X - B.X
	c1 = B.X*A.Y - A.X*B.Y
	a2 = F.Y - E.Y
	b2 = E.X - F.X
	c2 = F.X*E.Y - E.X*F.Y

	var denom = a1*b2 - a2*b1

	x = (b1*c2 - b2*c1) / denom
	y = (a2*c1 - a1*c2) / denom

	if math.IsInf(x, 0) || math.IsInf(y, 0) {
		return nil
	}

	// lines are colinear
	/*var crossABE = (E.y - A.y) * (B.x - A.x) - (E.x - A.x) * (B.y - A.y);
	var crossABF = (F.y - A.y) * (B.x - A.x) - (F.x - A.x) * (B.y - A.y);
	if(_almostEqual2(crossABE,0) && _almostEqual2(crossABF,0)){
		return nil;
	}*/

	if !infinite {
		// coincident points do not count as intersecting
		if math.Abs(A.X-B.X) > TOL {
			if A.X < B.X {
				if x < A.X || x > B.X {
					return nil
				}
			} else {
				if x > A.X || x < B.X {
					return nil
				}
			}

		}
		if math.Abs(A.Y-B.Y) > TOL {
			if A.Y < B.Y {
				if y < A.Y || y > B.Y {
					return nil
				}
			} else {
				if y > A.Y || y < B.Y {
					return nil
				}
			}

		}

		if math.Abs(E.X-F.X) > TOL {
			if E.X < F.X {
				if x < E.X || x > F.X {
					return nil
				}
			} else {
				if x > E.X || x < F.X {
					return nil
				}
			}

		}
		if math.Abs(E.Y-F.Y) > TOL {
			if E.Y < F.Y {
				if y < E.Y || y > F.Y {
					return nil
				}
			} else {
				if y > E.Y || y < F.Y {
					return nil
				}
			}
		}
	}

	return &Point{X: x, Y: y}
}

func (this *SVG) polygonProjectionDistance(oldA, oldB *polygonWithOffset, direction *vectorsStruct) *float64 {

	var Boffsetx = oldB.offsetx
	var Boffsety = oldB.offsety

	var Aoffsetx = oldA.offsetx
	var Aoffsety = oldA.offsety

	A := oldA.Polygon[:]
	B := oldB.Polygon[:]

	// close the loop for polygons
	if A[0] != A[len(A)-1] {
		A = append(A, A[0])
	}

	if B[0] != B[len(B)-1] {
		B = append(B, B[0])
	}

	var edgeA = A
	var edgeB = B

	var distance *float64 = nil
	var p, s1, s2 *Point
	var d *float64
	for i := 0; i < len(edgeB); i++ {
		// the shortest/most negative projection of B onto A
		var minprojection *float64 = nil
		//var minp *Point = nil
		for j := 0; j < len(edgeA)-1; j++ {
			p = &Point{X: edgeB[i].X + Boffsetx, Y: edgeB[i].Y + Boffsety}
			s1 = &Point{X: edgeA[j].X + Aoffsetx, Y: edgeA[j].Y + Aoffsety}
			s2 = &Point{X: edgeA[j+1].X + Aoffsetx, Y: edgeA[j+1].Y + Aoffsety}

			if math.Abs((s2.Y-s1.Y)*direction.x-(s2.X-s1.X)*direction.y) < TOL {
				continue
			}

			// project Point, ignore edge boundaries
			d = this.pointDistance(p, s1, s2, direction, false)

			if d != nil && (minprojection == nil || *d < *minprojection) {
				minprojection = d
				//minp = p
			}
		}
		if minprojection != nil && (distance == nil || *minprojection > *distance) {
			distance = minprojection
		}
	}

	return distance

}
func (this *SVG) searchStartPoint(oldA, oldB *polygonWithOffset, inside bool, NFP []Polygon) *Point {

	A := &polygonWithOffset{
		Polygon: oldA.Polygon[:],
	}
	B := &polygonWithOffset{
		Polygon: oldB.Polygon[:],
	}
	if A.Polygon[0] != A.Polygon[len(A.Polygon)-1] {
		A.Polygon = append(A.Polygon, A.Polygon[0])
	}

	if B.Polygon[0] != B.Polygon[len(B.Polygon)-1] {
		B.Polygon = append(B.Polygon, B.Polygon[0])
	}

	for i := 0; i < len(A.Polygon)-1; i++ {
		if !A.Polygon[i].Marked {
			A.Polygon[i].Marked = true
			for j := 0; j < len(B.Polygon); j++ {
				B.offsetx = A.Polygon[i].X - B.Polygon[j].X
				B.offsety = A.Polygon[i].Y - B.Polygon[j].Y

				var Binside int64 = 0
				for k := 0; k < len(B.Polygon); k++ {
					inpoly := pointInPolygon(&Point{X: B.Polygon[k].X + B.offsetx, Y: B.Polygon[k].Y + B.offsety}, A.Polygon)
					if inpoly != 0 {
						Binside = inpoly
						break
					}
				}

				if Binside == 0 { // A and B are the same
					return nil
				}
				// returns true if Point already exists in the given nfp
				inNfp := func(p *Point, nfp []Polygon) bool {
					if len(nfp) == 0 {
						return false
					}

					for i := 0; i < len(nfp); i++ {
						for j := 0; j < len(nfp[i]); j++ {
							if _almostEqual2(p.X, nfp[i][j].X) && _almostEqual2(p.Y, nfp[i][j].Y) {
								return true
							}
						}
					}

					return false
				}

				var startPoint = &Point{X: B.offsetx, Y: B.offsety}
				if ((Binside > 0 && inside) || (Binside < 0 && !inside)) && !this.intersect(A, B) && !inNfp(startPoint, NFP) {
					return startPoint
				}

				// slide B along vector
				var vx = A.Polygon[i+1].X - A.Polygon[i].X
				var vy = A.Polygon[i+1].Y - A.Polygon[i].Y

				var d1 = this.polygonProjectionDistance(A, B, &vectorsStruct{x: vx, y: vy})
				var d2 = this.polygonProjectionDistance(B, A, &vectorsStruct{x: -vx, y: -vy})

				var d *float64 = nil

				// todo: clean this up
				if d1 == nil && d2 == nil {
					// nothin
				} else if d1 == nil {
					d = d2
				} else if d2 == nil {
					d = d1
				} else {
					temp := math.Min(*d1, *d2)
					d = &temp
				}

				// only slide until no longer negative
				// todo: clean this up
				if d != nil && !_almostEqual2(*d, 0) && *d > 0 {

				} else {
					continue
				}

				var vd2 = vx*vx + vy*vy

				if *d**d < vd2 && !_almostEqual2(*d**d, vd2) {
					var vd = math.Sqrt(vx*vx + vy*vy)
					vx *= *d / vd
					vy *= *d / vd
				}

				B.offsetx += vx
				B.offsety += vy

				for k := 0; k < len(B.Polygon); k++ {
					var inpoly = pointInPolygon(&Point{X: B.Polygon[k].X + B.offsetx, Y: B.Polygon[k].Y + B.offsety}, A.Polygon)
					if inpoly != 0 {
						Binside = inpoly
						break
					}
				}
				startPoint = &Point{X: B.offsetx, Y: B.offsety}
				if ((Binside > 0 && inside) || (Binside < 0 && !inside)) && !this.intersect(A, B) && !inNfp(startPoint, NFP) {
					return startPoint
				}
			}
		}
	}

	return nil

}
func _normalizeVector(v *vectorsStruct) *vectorsStruct {
	if _almostEqual2(v.x*v.x+v.y*v.y, 1) {
		return v // given vector was already a unit vector
	}
	var len = math.Sqrt(v.x*v.x + v.y*v.y)
	var inverse = 1 / len

	return &vectorsStruct{
		x: v.x * inverse,
		y: v.y * inverse,
	}
}

func (this *SVG) pointDistance(p, s1, s2 *Point, normal *vectorsStruct, infinite bool) *float64 {

	normal = _normalizeVector(normal)

	var dir = &vectorsStruct{
		x: normal.y,
		y: -normal.x,
	}

	var pdot = p.X*dir.x + p.Y*dir.y
	var s1dot = s1.X*dir.x + s1.Y*dir.y
	var s2dot = s2.X*dir.x + s2.Y*dir.y

	var pdotnorm = p.X*normal.x + p.Y*normal.y
	var s1dotnorm = s1.X*normal.x + s1.Y*normal.y
	var s2dotnorm = s2.X*normal.x + s2.Y*normal.y

	if !infinite {
		if ((pdot < s1dot || _almostEqual2(pdot, s1dot)) && (pdot < s2dot || _almostEqual2(pdot, s2dot))) || ((pdot > s1dot || _almostEqual2(pdot, s1dot)) && (pdot > s2dot || _almostEqual2(pdot, s2dot))) {
			return nil // dot doesn't collide with segment, or lies directly on the vertex
		}
		if (_almostEqual2(pdot, s1dot) && _almostEqual2(pdot, s2dot)) && (pdotnorm > s1dotnorm && pdotnorm > s2dotnorm) {
			rst := math.Min(pdotnorm-s1dotnorm, pdotnorm-s2dotnorm)

			return &rst
		}
		if (_almostEqual2(pdot, s1dot) && _almostEqual2(pdot, s2dot)) && (pdotnorm < s1dotnorm && pdotnorm < s2dotnorm) {
			rst := -math.Min(s1dotnorm-pdotnorm, s2dotnorm-pdotnorm)
			return &rst
		}
	}
	rst := -(pdotnorm - s1dotnorm + (s1dotnorm-s2dotnorm)*(s1dot-pdot)/(s1dot-s2dot))
	return &rst
}
func (this *SVG) segmentDistance(A, B, E, F *Point, direction *vectorsStruct) *float64 {

	var normal = &Point{
		X: direction.y,
		Y: -direction.x,
	}

	var reverse = &vectorsStruct{
		x: -direction.x,
		y: -direction.y,
	}

	var dotA = A.X*normal.X + A.Y*normal.Y
	var dotB = B.X*normal.X + B.Y*normal.Y
	var dotE = E.X*normal.X + E.Y*normal.Y
	var dotF = F.X*normal.X + F.Y*normal.Y

	var crossA = A.X*direction.x + A.Y*direction.y
	var crossB = B.X*direction.x + B.Y*direction.y
	var crossE = E.X*direction.x + E.Y*direction.y
	var crossF = F.X*direction.x + F.Y*direction.y

	// var crossABmin = math.Min(crossA, crossB)
	// var crossABmax = math.Max(crossA, crossB)

	// var crossEFmax = math.Max(crossE, crossF)
	// var crossEFmin = math.Min(crossE, crossF)//saya ?

	var ABmin = math.Min(dotA, dotB)
	var ABmax = math.Max(dotA, dotB)

	var EFmax = math.Max(dotE, dotF)
	var EFmin = math.Min(dotE, dotF)

	// segments that will merely touch at one Point
	if _almostEqual3(ABmax, EFmin, TOL) || _almostEqual3(ABmin, EFmax, TOL) {
		return nil
	}
	// segments miss eachother completely
	if ABmax < EFmin || ABmin > EFmax {
		return nil
	}
	var overlap float64
	if (ABmax > EFmax && ABmin < EFmin) || (EFmax > ABmax && EFmin < ABmin) {
		overlap = 1
	} else {
		var minMax = math.Min(ABmax, EFmax)
		var maxMin = math.Max(ABmin, EFmin)

		var maxMax = math.Max(ABmax, EFmax)
		var minMin = math.Min(ABmin, EFmin)

		overlap = (minMax - maxMin) / (maxMax - minMin)
	}

	var crossABE = (E.Y-A.Y)*(B.X-A.X) - (E.X-A.X)*(B.Y-A.Y)
	var crossABF = (F.Y-A.Y)*(B.X-A.X) - (F.X-A.X)*(B.Y-A.Y)

	// lines are colinear
	if _almostEqual2(crossABE, 0) && _almostEqual2(crossABF, 0) {

		var ABnorm = &Point{X: B.Y - A.Y, Y: A.X - B.X}
		var EFnorm = &Point{X: F.Y - E.Y, Y: E.X - F.X}

		var ABnormlength = math.Sqrt(ABnorm.X*ABnorm.X + ABnorm.Y*ABnorm.Y)
		ABnorm.X /= ABnormlength
		ABnorm.Y /= ABnormlength

		var EFnormlength = math.Sqrt(EFnorm.X*EFnorm.X + EFnorm.Y*EFnorm.Y)
		EFnorm.X /= EFnormlength
		EFnorm.Y /= EFnormlength

		// segment normals must Point in opposite directions
		if math.Abs(ABnorm.Y*EFnorm.X-ABnorm.X*EFnorm.Y) < TOL && ABnorm.Y*EFnorm.Y+ABnorm.X*EFnorm.X < 0 {
			// normal of AB segment must Point in same direction as given direction vector
			var normdot = ABnorm.Y*direction.y + ABnorm.X*direction.x
			// the segments merely slide along eachother
			if _almostEqual3(normdot, 0, TOL) {
				return nil
			}
			if normdot < 0 {
				var rst0 float64 = 0
				return &rst0
			}
		}
		return nil
	}

	var distances = []float64{}

	// coincident points
	if _almostEqual2(dotA, dotE) {
		distances = append(distances, crossA-crossE)
	} else if _almostEqual2(dotA, dotF) {
		distances = append(distances, crossA-crossF)
	} else if dotA > EFmin && dotA < EFmax {
		var d = this.pointDistance(A, E, F, reverse, false)
		if d != nil && _almostEqual2(*d, 0) { //  A currently touches EF, but AB is moving away from EF
			var dB = this.pointDistance(B, E, F, reverse, true)
			if *dB < 0 || _almostEqual2(*dB*overlap, 0) {
				d = nil
			}
		}
		if d != nil {
			distances = append(distances, *d)
		}
	}

	if _almostEqual2(dotB, dotE) {
		distances = append(distances, crossB-crossE)
	} else if _almostEqual2(dotB, dotF) {
		distances = append(distances, crossB-crossF)
	} else if dotB > EFmin && dotB < EFmax {
		var d = this.pointDistance(B, E, F, reverse, false)

		if d != nil && _almostEqual2(*d, 0) { // crossA>crossB A currently touches EF, but AB is moving away from EF
			var dA = this.pointDistance(A, E, F, reverse, true)
			if *dA < 0 || _almostEqual2(*dA*overlap, 0) {
				d = nil
			}
		}
		if d != nil {
			distances = append(distances, *d)
		}
	}

	if dotE > ABmin && dotE < ABmax {
		var d = this.pointDistance(E, A, B, direction, false)
		if d != nil && _almostEqual2(*d, 0) { // crossF<crossE A currently touches EF, but AB is moving away from EF
			var dF = this.pointDistance(F, A, B, direction, true)
			if *dF < 0 || _almostEqual2(*dF*overlap, 0) {
				d = nil
			}
		}
		if d != nil {
			distances = append(distances, *d)
		}
	}

	if dotF > ABmin && dotF < ABmax {
		var d = this.pointDistance(F, A, B, direction, false)
		if d != nil && _almostEqual2(*d, 0) { // && crossE<crossF A currently touches EF, but AB is moving away from EF
			var dE = this.pointDistance(E, A, B, direction, true)
			if *dE < 0 || _almostEqual2(*dE*overlap, 0) {
				d = nil
			}
		}
		if d != nil {
			distances = append(distances, *d)
		}
	}

	if len(distances) == 0 {
		return nil
	}
	if len(distances) <= 0 {
		return nil
	}
	var minDis float64 = distances[0]
	for index := range distances {
		minDis = math.Min(minDis, distances[index])
	}
	return &minDis
}

func (this *SVG) polygonSlideDistance(oldA, oldB *polygonWithOffset, direction *vectorsStruct, ignoreNegative bool) *float64 {

	var A1, A2, B1, B2 *Point

	Aoffsetx := oldA.offsetx
	Aoffsety := oldA.offsety

	Boffsetx := oldB.offsetx
	Boffsety := oldB.offsety

	A := oldA.Polygon[:]
	B := oldB.Polygon[:]

	// close the loop for polygons
	if A[0] != A[len(A)-1] {
		A = append(A, A[0])
	}

	if B[0] != B[len(B)-1] {
		B = append(B, B[0])
	}

	var edgeA = A
	var edgeB = B

	var distance *float64
	//var p, s1, s2 float64
	var d *float64
	var dir = _normalizeVector(direction)

	// var normal = &Point{
	// 	X: dir.y,
	// 	Y: -dir.x,
	// }

	// var reverse = &Point{
	// 	X: -dir.x,
	// 	Y: -dir.y,
	// }

	for i := 0; i < len(edgeB)-1; i++ {
		//var mind = nil
		for j := 0; j < len(edgeA)-1; j++ {
			A1 = &Point{X: edgeA[j].X + Aoffsetx, Y: edgeA[j].Y + Aoffsety}
			A2 = &Point{X: edgeA[j+1].X + Aoffsetx, Y: edgeA[j+1].Y + Aoffsety}
			B1 = &Point{X: edgeB[i].X + Boffsetx, Y: edgeB[i].Y + Boffsety}
			B2 = &Point{X: edgeB[i+1].X + Boffsetx, Y: edgeB[i+1].Y + Boffsety}

			if (_almostEqual2(A1.X, A2.X) && _almostEqual2(A1.Y, A2.Y)) || (_almostEqual2(B1.X, B2.X) && _almostEqual2(B1.Y, B2.Y)) {
				continue // ignore extremely small lines
			}

			d = this.segmentDistance(A1, A2, B1, B2, dir)

			if d != nil && (distance == nil || *d < *distance) {
				if !ignoreNegative || *d > 0 || _almostEqual2(*d, 0) {
					distance = d
				}
			}
		}
	}
	return distance

}

func (this *SVG) noFitPolygon(A, B Polygon, inside bool, searchEdges bool) [][]*Point {

	if len(A) < 3 || len(B) < 3 {
		return nil
	}
	AWithOffset := &polygonWithOffset{Polygon: A}
	BWithOffset := &polygonWithOffset{Polygon: B}

	var minA = A[0].Y
	var minAindex = 0

	var maxB = B[0].Y
	var maxBindex = 0

	for i := 1; i < len(A); i++ {
		A[i].Marked = false
		if A[i].Y < minA {
			minA = A[i].Y
			minAindex = i
		}
	}

	for i := 1; i < len(B); i++ {
		B[i].Marked = false
		if B[i].Y > maxB {
			maxB = B[i].Y
			maxBindex = i
		}
	}
	var startpoint *Point = nil
	if !inside {
		// shift B such that the bottom-most Point of B is at the top-most Point of A. This guarantees an initial placement with no intersections
		startpoint = &Point{
			X: A[minAindex].X - B[maxBindex].X,
			Y: A[minAindex].Y - B[maxBindex].Y,
		}
	} else {
		// no reliable heuristic for inside
		startpoint = this.searchStartPoint(AWithOffset, BWithOffset, true, nil)
	}

	var NFPlist = []Polygon{}

	for startpoint != nil {

		BWithOffset.offsetx = startpoint.X
		BWithOffset.offsety = startpoint.Y

		// maintain a list of touching points/edges
		touching := []*OneTouchStruct{}

		var prevvector *vectorsStruct = nil // keep track of previous vector
		NFP := Polygon{
			&Point{
				X: B[0].X + BWithOffset.offsetx,
				Y: B[0].Y + BWithOffset.offsety,
			},
		}

		var referencex = B[0].X + BWithOffset.offsetx
		var referencey = B[0].Y + BWithOffset.offsety
		var startx = referencex
		var starty = referencey
		var counter = 0

		for counter < 10*(len(A)+len(B)) { // sanity check, prevent infinite loop
			touching = []*OneTouchStruct{}
			// find touching vertices/edges
			for i := 0; i < len(A); i++ {
				var nexti int
				if i == len(A)-1 {
					nexti = 0
				} else {
					nexti = i + 1
				}

				for j := 0; j < len(B); j++ {
					var nextj int
					if j == len(B)-1 {
						nextj = 0
					} else {
						nextj = j + 1
					}

					if _almostEqual2(A[i].X, B[j].X+BWithOffset.offsetx) && _almostEqual2(A[i].Y, B[j].Y+BWithOffset.offsety) {
						touching = append(touching, &OneTouchStruct{Type: 0, A: i, B: j})
					} else if (_onSegment(A[i], A[nexti], &Point{X: B[j].X + BWithOffset.offsetx, Y: B[j].Y + BWithOffset.offsety})) {
						touching = append(touching, &OneTouchStruct{Type: 1, A: nexti, B: j})
					} else if (_onSegment(&Point{X: B[j].X + BWithOffset.offsetx, Y: B[j].Y + BWithOffset.offsety}, &Point{X: B[nextj].X + BWithOffset.offsetx, Y: B[nextj].Y + BWithOffset.offsety}, A[i])) {
						touching = append(touching, &OneTouchStruct{Type: 2, A: i, B: nextj})
					}
				}
			}

			// generate translation vectors from touching vertices/edges
			var vectors = []*vectorsStruct{}
			for i := 0; i < len(touching); i++ {
				var vertexA = A[touching[i].A]
				vertexA.Marked = true

				// adjacent A vertices
				var prevAindex = touching[i].A - 1
				var nextAindex = touching[i].A + 1

				if prevAindex < 0 {
					prevAindex = len(A) - 1
				} else {
					//prevAindex = prevAindex
				}
				if nextAindex >= len(A) {
					nextAindex = 0
				} else {
					//nextAindex = nextAindex
				}

				var prevA = A[prevAindex]
				var nextA = A[nextAindex]

				// adjacent B vertices
				var vertexB = B[touching[i].B]

				var prevBindex = touching[i].B - 1
				var nextBindex = touching[i].B + 1

				if prevBindex < 0 {
					prevBindex = len(B) - 1
				} else {
					//prevBindex = prevBindex
				}

				if nextBindex >= len(B) {
					nextBindex = 0
				} else {
					//nextBindex = nextBindex
				}

				var prevB = B[prevBindex]
				var nextB = B[nextBindex]

				if touching[i].Type == 0 {

					vA1 := &vectorsStruct{
						x:     prevA.X - vertexA.X,
						y:     prevA.Y - vertexA.Y,
						start: vertexA,
						end:   prevA,
					}

					vA2 := &vectorsStruct{
						x:     nextA.X - vertexA.X,
						y:     nextA.Y - vertexA.Y,
						start: vertexA,
						end:   nextA,
					}

					// B vectors need to be inverted
					vB1 := &vectorsStruct{
						x:     vertexB.X - prevB.X,
						y:     vertexB.Y - prevB.Y,
						start: prevB,
						end:   vertexB,
					}
					vB2 := &vectorsStruct{
						x:     vertexB.X - nextB.X,
						y:     vertexB.Y - nextB.Y,
						start: nextB,
						end:   vertexB,
					}
					vectors = append(vectors, vA1)
					vectors = append(vectors, vA2)
					vectors = append(vectors, vB1)
					vectors = append(vectors, vB2)
				} else if touching[i].Type == 1 {
					vectors = append(vectors, &vectorsStruct{
						x:     vertexA.X - (vertexB.X + BWithOffset.offsetx),
						y:     vertexA.Y - (vertexB.Y + BWithOffset.offsety),
						start: prevA,
						end:   vertexA,
					})
					vectors = append(vectors, &vectorsStruct{
						x:     prevA.X - (vertexB.X + BWithOffset.offsetx),
						y:     prevA.Y - (vertexB.Y + BWithOffset.offsety),
						start: vertexA,
						end:   prevA,
					})
				} else if touching[i].Type == 2 {
					vectors = append(vectors, &vectorsStruct{
						x:     vertexA.X - (vertexB.X + BWithOffset.offsetx),
						y:     vertexA.Y - (vertexB.Y + BWithOffset.offsety),
						start: prevB,
						end:   vertexB,
					})

					vectors = append(vectors, &vectorsStruct{
						x:     vertexA.X - (prevB.X + BWithOffset.offsetx),
						y:     vertexA.Y - (prevB.Y + BWithOffset.offsety),
						start: vertexB,
						end:   prevB,
					})
				}
			}

			// todo: there should be a faster way to reject vectors that will cause immediate intersection. For now just check them all

			var translate *vectorsStruct = nil
			var maxd float64 = 0

			for i := 0; i < len(vectors); i++ {
				if vectors[i].x == 0 && vectors[i].y == 0 {
					continue
				}

				// if this vector points us back to where we came from, ignore it.
				// ie cross product = 0, dot product < 0
				if prevvector != nil && vectors[i].y*prevvector.y+vectors[i].x*prevvector.x < 0 {
					// compare magnitude with unit vectors
					var vectorlength = math.Sqrt(vectors[i].x*vectors[i].x + vectors[i].y*vectors[i].y)
					unitv := &vectorsStruct{x: vectors[i].x / vectorlength, y: vectors[i].y / vectorlength}

					var prevlength = math.Sqrt(prevvector.x*prevvector.x + prevvector.y*prevvector.y)
					prevunit := &vectorsStruct{x: prevvector.x / prevlength, y: prevvector.y / prevlength}

					// we need to scale down to unit vectors to normalize len(vector). Could also just do a tan here
					if math.Abs(unitv.y*prevunit.x-unitv.x*prevunit.y) < 0.0001 {
						continue
					}
				}

				d := this.polygonSlideDistance(AWithOffset, BWithOffset, vectors[i], true)

				var vecd2 = vectors[i].x*vectors[i].x + vectors[i].y*vectors[i].y

				if d == nil || *d**d > vecd2 {
					var vecd = math.Sqrt(vectors[i].x*vectors[i].x + vectors[i].y*vectors[i].y)
					d = &vecd
				}
				if d != nil && *d > maxd {
					maxd = *d
					translate = vectors[i]
				}
			}

			if translate == nil || _almostEqual2(maxd, 0) {
				// didn't close the loop, something went wrong here
				NFP = nil
				break
			}
			this.logIfDebug("translate:", translate.x, translate.y, translate.start, translate.end)
			translate.start.Marked = true
			translate.end.Marked = true

			prevvector = translate

			// trim
			var vlength2 = translate.x*translate.x + translate.y*translate.y
			if maxd*maxd < vlength2 && !_almostEqual2(maxd*maxd, vlength2) {
				var scale = math.Sqrt((maxd * maxd) / vlength2)
				translate.x *= scale
				translate.y *= scale
			}

			referencex += translate.x
			referencey += translate.y

			if _almostEqual2(referencex, startx) && _almostEqual2(referencey, starty) {
				// we've made a full loop
				break
			}

			// if A and B start on a touching horizontal line, the end Point may not be the start Point
			var looped = false
			if len(NFP) > 0 {
				for i := 0; i < len(NFP)-1; i++ {
					if _almostEqual2(referencex, NFP[i].X) && _almostEqual2(referencey, NFP[i].Y) {
						looped = true
					}
				}
			}

			if looped {
				// we've made a full loop
				break
			}

			NFP = append(NFP, &Point{
				X: referencex,
				Y: referencey,
			})

			BWithOffset.offsetx += translate.x
			BWithOffset.offsety += translate.y

			counter++
		}

		if len(NFP) > 0 {
			NFPlist = append(NFPlist, NFP)
		}

		if !searchEdges {
			// only get outer NFP or first inner NFP
			break
		}

		startpoint = this.searchStartPoint(AWithOffset, BWithOffset, inside, NFPlist)

	}

	return NFPlist
}
func RecBoundSum(A, B *BoundStruct) *BoundStruct {
	C := &BoundStruct{}
	pointArr := []*Point{
		&Point{
			X: A.x,
			Y: A.y,
		},
		&Point{
			X: A.x + A.width,
			Y: A.y + A.height,
		},
		&Point{
			X: B.x,
			Y: B.y,
		},
		&Point{
			X: B.x + B.width,
			Y: B.y + B.height,
		},
	}
	var left *float64
	var right *float64
	var top *float64
	var button *float64

	for index := range pointArr {
		if left == nil || *left > pointArr[index].X {
			left = &pointArr[index].X
		}
		if right == nil || *right < pointArr[index].X {
			right = &pointArr[index].X
		}
		if top == nil || *top < pointArr[index].Y {
			top = &pointArr[index].Y
		}
		if button == nil || *button > pointArr[index].X {
			button = &pointArr[index].Y
		}
	}
	C.x = *left
	C.y = *button
	C.width = *right - *left
	C.height = *top - *button
	return C
}
func getPolygonBounds(myPolygon Polygon) *BoundStruct {
	if myPolygon == nil {
		return &BoundStruct{
			x:      0,
			y:      0,
			width:  0,
			height: 0,
		}
	}
	if len(myPolygon) < 3 {
		return nil
	}

	var xmin = myPolygon[0].X
	var xmax = myPolygon[0].X
	var ymin = myPolygon[0].Y
	var ymax = myPolygon[0].Y

	for i := 1; i < len(myPolygon); i++ {
		if myPolygon[i].X > xmax {
			xmax = myPolygon[i].X
		} else if myPolygon[i].X < xmin {
			xmin = myPolygon[i].X
		}

		if myPolygon[i].Y > ymax {
			ymax = myPolygon[i].Y
		} else if myPolygon[i].Y < ymin {
			ymin = myPolygon[i].Y
		}
	}

	return &BoundStruct{
		x:      xmin,
		y:      ymin,
		width:  xmax - xmin,
		height: ymax - ymin,
	}
}
func rotatePolygonA(myPolygon *PolygonStruct, angle int) ([]*Point, *BoundStruct) {
	var rotated = []*Point{}
	floatAngle := float64(angle) * math.Pi / 180
	for i := 0; i < len(myPolygon.RootPoly.polygonBeforeRotation); i++ {
		var x = myPolygon.RootPoly.polygonBeforeRotation[i].X
		var y = myPolygon.RootPoly.polygonBeforeRotation[i].Y
		var x1 = x*math.Cos(floatAngle) - y*math.Sin(floatAngle)
		var y1 = x*math.Sin(floatAngle) + y*math.Cos(floatAngle)

		rotated = append(rotated, &Point{X: x1, Y: y1})
	}
	// reset bounding box
	var bounds = getPolygonBounds(rotated)
	// rotated.X = bounds.X
	// rotated.y = bounds.y
	// rotated.width = bounds.width
	// rotated.height = bounds.height

	return rotated, bounds
}
