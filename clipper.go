package nest

import (
	"math"
)

func ScaleUpPaths(a []Polygon, times int64) []IntPolygon {

	if times <= 0 {
		times = 1
	}
	b := []IntPolygon{}
	for c := len(a) - 1; c >= 0; c-- {
		oneb := IntPolygon{}
		for d := len(a[c]) - 1; d >= 0; d-- {
			oneb = append(oneb, &IntPoint{int64(math.Round(a[c][d].X * float64(times))), int64(math.Round(a[c][d].Y * float64(times)))})
		}
		b = append(b, oneb)
	}
	return b
}

func ScaleUpPath(a Polygon, times int64) IntPolygon {
	if times <= 0 {
		times = 1
	}
	b := IntPolygon{}
	for index := range a {
		b = append(b, &IntPoint{
			X: int64(math.Round(a[index].X * float64(times))),
			Y: int64(math.Round(a[index].Y * float64(times))),
		})
	}
	return b
}

func ScaleDownPath(a IntPolygon, times int64) Path {
	if times <= 0 {
		times = 1
	}
	b := Path{}
	for index := range a {
		b = append(b, &Point{
			X: float64(a[index].X) / float64(times),
			Y: float64(a[index].Y) / float64(times),
		})
	}
	return b
}

func toClipperCoordinates(a Polygon) Polygon {
	b := Polygon{}
	for index := range a {
		b = append(b, &Point{
			X: a[index].X,
			Y: a[index].Y,
		})
	}
	return b
}

func minkowskiDifference(A, B Polygon) [][]*Point {
	Acpolygon := toClipperCoordinates(A)
	Ac := ScaleUpPath(Acpolygon, clipperScaleTimes)
	Bcpolygon := toClipperCoordinates(B)
	Bc := ScaleUpPath(Bcpolygon, clipperScaleTimes)
	for i := 0; i < len(Bc); i++ {
		Bc[i].X *= -1
		Bc[i].Y *= -1
	}
	solution := MinkowskiSum3(Ac, Bc, true)
	var largestArea float64 = 0
	clipperNfp := []*Point{}
	for i := 0; i < len(solution); i++ {
		n := toNestCoordinates(solution[i], clipperScaleTimes)
		sarea := polygonArea(n)
		if math.Abs(largestArea) < math.Abs(sarea) { //saya
			clipperNfp = n
			largestArea = sarea
		}
	}
	for i := 0; i < len(clipperNfp); i++ {
		clipperNfp[i].X += B[0].X
		clipperNfp[i].Y += B[0].Y
	}
	return [][]*Point{clipperNfp}
}

func toNestCoordinates(Polygon []*IntPoint, scale int64) []*Point {
	clone := []*Point{}
	for i := 0; i < len(Polygon); i++ {
		clone = append(clone, &Point{
			X: float64(Polygon[i].X) / float64(scale),
			Y: float64(Polygon[i].Y) / float64(scale),
		})
	}
	return clone
}
func PointsAreClose(a, b *Point, c float64) bool {
	el := a.X - b.X
	al := a.Y - b.Y
	return el*el+al*al <= c
}
func ExcludeOp(a *OutFloatPtStruct) *OutFloatPtStruct {
	var b = a.Prev
	b.Next = a.Next
	a.Next.Prev = b
	b.Idx = 0
	return b
}
func SlopesNearCollinear(a, b, c *Point, e float64) bool {
	return DistanceFromLineSqrd(b, a, c) < e
}
func DistanceFromLineSqrd(a, b, c *Point) float64 {
	var el = b.Y - c.Y
	cl := c.X - b.X
	bl := el*b.X + cl*b.Y
	bl = el*a.X + cl*a.Y - bl
	return bl * bl / (el*el + cl*cl)
}

// func PrintPolygonList(a [][]*IntPoint) {
// 	for index := range a {
// 		logIfDebug("Point:", index)
// 		for index2 := range a[index] {
// 			logIfDebug(a[index][index2].Print())
// 		}
// 	}
// }
func MinkowskiSum3(Ac, Bc IntPolygon, unkonwBool bool) [][]*IntPoint {
	return Minkowski(Ac, Bc, true, unkonwBool)
}
func Minkowski(a, b IntPolygon, c_bool bool, e_bool bool) [][]*IntPoint {
	var f int
	var aPolygonLen, bPolygonLen int
	aPolygonLen = len(a)
	bPolygonLen = len(b)
	if e_bool {
		f = 1
	} else {
		f = 0
	}

	ePolygonList := [][]*IntPoint{}

	if c_bool {
		for c := 0; c < bPolygonLen; c++ {
			l := []*IntPoint{}
			k := 0
			n := len(a)
			m := a[k]
			for k < n {
				l = append(l, IntPoint2(b[c].X+m.X, b[c].Y+m.Y))
				k++
				if k < len(a) {
					m = a[k]
				}
			}
			ePolygonList = append(ePolygonList, l)
		}
	} else {
		for c := 0; c < bPolygonLen; c++ {
			l := []*IntPoint{}
			k := 0
			n := len(a)
			m := a[k]
			for k < n {
				l = append(l, IntPoint2(b[c].X-m.X, b[c].Y-m.Y))
				k++
				if k < len(a) {
					m = a[k]
				}
			}
			ePolygonList = append(ePolygonList, l)
		}
	}
	//PrintPolygonList(ePolygonList)
	aPolygonList := [][]*IntPoint{}
	for c := 0; c < bPolygonLen-1+f; c++ {

		for k := 0; k < aPolygonLen; k++ {

			b = []*IntPoint{}
			b = append(b, ePolygonList[c%bPolygonLen][k%aPolygonLen])
			b = append(b, ePolygonList[(c+1)%bPolygonLen][k%aPolygonLen])
			b = append(b, ePolygonList[(c+1)%bPolygonLen][(k+1)%aPolygonLen])
			b = append(b, ePolygonList[c%bPolygonLen][(k+1)%aPolygonLen])
			if !Orientation(b) {
				// logIfDebug("前：", b[0])
				// logIfDebug("前：", b[1])
				// logIfDebug("前：", b[2])
				// logIfDebug("前：", b[3])
				b = IntPointReverse(b)
				// logIfDebug("后：", b[0])
				// logIfDebug("后：", b[1])
				// logIfDebug("后：", b[2])
				// logIfDebug("后：", b[3])
			}
			aPolygonList = append(aPolygonList, b)
		}
	}
	//logIfDebug("e----")
	//PrintPolygonList(ePolygonList)
	//logIfDebug("a----")
	//PrintPolygonList(aPolygonList)

	fClipper := Clipper(0)
	fClipper.AddPaths(aPolygonList, ptSubject, true)

	executePolygonList, _ := fClipper.Execute(ctUnion, pftNonZero, pftNonZero)
	return executePolygonList
}

func IntPoint2(x int64, y int64) *IntPoint {

	return &IntPoint{
		X: x,
		Y: y,
	}
}
func Orientation(a IntPolygon) bool {
	return 0 <= Area(a)
}
func Area(a IntPolygon) int64 {
	pointNum := len(a)
	if 3 > pointNum {
		return 0
	}
	//相邻两点依次操作
	var c int64 = 0
	var e = 0
	var d = pointNum - 1
	for ; e < pointNum; e++ {
		// logIfDebug("d:", d)
		// logIfDebug("e:", e)
		// logIfDebug("a[d].X + a[e].X:", a[d].X+a[e].X)
		// logIfDebug("(a[d].Y - a[e].Y):", (a[d].Y - a[e].Y))
		// logIfDebug("(a[d].X + a[e].X) * (a[d].Y - a[e].Y):", (a[d].X+a[e].X)*(a[d].Y-a[e].Y))
		c += (a[d].X + a[e].X) * (a[d].Y - a[e].Y)
		d = e
	}
	return c / 2 * -1
}
func op_Equality_IntPoint(a, b *IntPoint) bool {
	return a.X == b.X && a.Y == b.Y
}
func op_InEquality_IntPoint(a, b *IntPoint) bool {
	return a.X != b.X || a.Y != b.Y
}

type TEdgeStruct struct {
	Bot       *IntPoint
	Curr      *IntPoint
	Top       *IntPoint
	Delta     *IntPoint
	Dx        float64
	PolyTyp   PolyType
	Side      EdgeSide
	OutIdx    int
	WindCnt2  int64
	WindCnt   int64
	WindDelta int64

	PrevInSEL *TEdgeStruct
	NextInSEL *TEdgeStruct
	PrevInAEL *TEdgeStruct
	NextInAEL *TEdgeStruct
	NextInLML *TEdgeStruct

	Prev *TEdgeStruct
	Next *TEdgeStruct
}

func NewTEdge() *TEdgeStruct {
	return &TEdgeStruct{
		Bot:   &IntPoint{},
		Curr:  &IntPoint{},
		Top:   &IntPoint{},
		Delta: &IntPoint{},
		Dx:    0,
	}

}
