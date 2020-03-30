package nest

import (
	"fmt"
	"log"
	"math"
)

func (this *SVG) logIfDebug(v ...interface{}) {
	if this.config.IfDebug {
		log.Println(v...)
	}
}

type Polygon = []*Point
type polygonWithOffset struct {
	Polygon []*Point
	offsetx float64
	offsety float64
}
type PolyNode struct {
	OriginPolygon []*Point //原始多边形
	polygonAfterRotaion       []*Point //简化后 膨胀后的多边形
	polygonBeforeRotation       []*Point //简化后 膨胀后的多边形
	EndPolygon    []*Point //未简化 未膨胀  旋转 平移后的多边形
	children []*PolyNode
	parent   *PolyNode
}

type PolygonStruct1 struct {
	OriginPolygon []*Point //原始多边形
	polygon       []*Point //简化后 膨胀后的多边形
	EndPolygon    []*Point //未简化 未膨胀  旋转 平移后的多边形
	children []*PolygonStruct
	parent   *PolygonStruct
	id       int //全局唯一id
	typeID   int //同一类型的  NFP 无需重复生成//但要注意这个字段可能没被初始化的  没被初始化 则用id代替
//	source   int
	width    float64
	height   float64
	rotation int

	Name        string
	isWart     bool
	AngleList   []int64
}
type BinPolygonStruct struct {
	myPolygon Polygon
	height    float64
	id        int
	width     float64
}

func CleanIntPolygon(this []*IntPoint, b float64) []*IntPoint {

	if b < 0 {
		b = 1.415 //？ saya
	}

	c := len(this)
	if 0 == c {
		return []*IntPoint{}
	}

	pointChain := []*OutFloatPtStruct{}

	for f := 0; f < c; f++ {
		pointChain = append(pointChain, &OutFloatPtStruct{})
	}

	for f := 0; f < c; f++ {
		pointChain[f].Pt = &Point{X: float64(this[f].X), Y: float64(this[f].Y)}
		pointChain[f].Next = pointChain[(f+1)%c]
		pointChain[f].Next.Prev = pointChain[f]
		pointChain[f].Idx = 0
	}

	fFloat := b * b
	eOutPt := pointChain[0]
	for 0 == eOutPt.Idx && eOutPt.Next != eOutPt.Prev {
		if PointsAreClose(eOutPt.Pt, eOutPt.Prev.Pt, fFloat) {
			eOutPt = ExcludeOp(eOutPt)
			c = c - 1
		} else {
			if PointsAreClose(eOutPt.Prev.Pt, eOutPt.Next.Pt, fFloat) {
				ExcludeOp(eOutPt.Next)
				eOutPt = ExcludeOp(eOutPt)
				c -= 2
			} else {
				if SlopesNearCollinear(eOutPt.Prev.Pt, eOutPt.Pt, eOutPt.Next.Pt, fFloat) {
					// if this.isWart {
					// 	if c > 4 {
					// 		eOutPt = ExcludeOp(eOutPt)
					// 		c = c - 1
					// 	} else {
					// 		break
					// 	}
					// } else {
					eOutPt = ExcludeOp(eOutPt)
					c = c - 1
					//}
				} else {
					eOutPt.Idx = 1
					eOutPt = eOutPt.Next
				}
			}
		}
	}
	if 3 > c {
		c = 0
	}
	g := []*IntPoint{}
	for f := 0; f < c; f++ {
		g = append(g, &IntPoint{X: int64(eOutPt.Pt.X), Y: int64(eOutPt.Pt.Y)})
		eOutPt = eOutPt.Next
	}
	// //多边形订正
	// fmt.Println("多边形订正")
	// for index := range g {
	// 	fmt.Println(g[index].Print())
	// }
	IsIn := func(intPoint *IntPoint) (int, int) {
		start := -1
		end := -1
		for index := range g {
			if g[index].X == intPoint.X && g[index].Y == intPoint.Y {
				if start < 0 {
					start = index
					continue
				} else if end < 0 {
					end = index
					continue
				}
			}
		}
		return start, end
	}
	for i := range g {
		start, end := IsIn(g[i])
		if start >= 0 && end >= 0 {
			//log.Println("clean works~")
			if end-start > len(g)/2 {
				return g[start:end]
			} else {
				cleanPol := []*IntPoint{}
				for j := range g {
					if j <= start || j > end {
						cleanPol = append(cleanPol, g[j])
					}
				}
				return cleanPol
			}
		}
	}
	return g
}
func (this *PolygonStruct) CleanPolygon(b float64) {
	if len(this.RootPoly.polygonBeforeRotation) != 0 {
		panic(" polygonBeforeRotation")
	}
	if b < 0 {
		b = 1.415 //？ saya
	}

	c := len(this.RootPoly.OriginPolygon)
	if 0 == c {
		return
	}

	pointChain := []*OutFloatPtStruct{}

	for f := 0; f < c; f++ {
		pointChain = append(pointChain, &OutFloatPtStruct{})
	}

	for f := 0; f < c; f++ {
		pointChain[f].Pt = this.RootPoly.OriginPolygon[f]
		pointChain[f].Next = pointChain[(f+1)%c]
		pointChain[f].Next.Prev = pointChain[f]
		pointChain[f].Idx = 0
	}

	fFloat := b * b
	eOutPt := pointChain[0]
	for 0 == eOutPt.Idx && eOutPt.Next != eOutPt.Prev {
		if PointsAreClose(eOutPt.Pt, eOutPt.Prev.Pt, fFloat) {
			eOutPt = ExcludeOp(eOutPt)
			c = c - 1
		} else {
			if PointsAreClose(eOutPt.Prev.Pt, eOutPt.Next.Pt, fFloat) {
				ExcludeOp(eOutPt.Next)
				eOutPt = ExcludeOp(eOutPt)
				c -= 2
			} else {
				if SlopesNearCollinear(eOutPt.Prev.Pt, eOutPt.Pt, eOutPt.Next.Pt, fFloat) {
					if this.isWart {
						if c > 4 {
							eOutPt = ExcludeOp(eOutPt)
							c = c - 1
						} else {
							break
						}
					} else {
						eOutPt = ExcludeOp(eOutPt)
						c = c - 1
					}
				} else {
					eOutPt.Idx = 1
					eOutPt = eOutPt.Next
				}
			}
		}
	}
	if 3 > c {
		c = 0
	}
	g := []*Point{}
	g2 := []*Point{}
	for f := 0; f < c; f++ {
		g = append(g, &Point{X: eOutPt.Pt.X, Y: eOutPt.Pt.Y})
		g2 = append(g2, &Point{X: eOutPt.Pt.X, Y: eOutPt.Pt.Y})
		eOutPt = eOutPt.Next
	}
	this.RootPoly.polygonBeforeRotation = g
	this.RootPoly.polygonAfterRotaion = g2
}

type PolygonStructSlice []*PolygonStruct

func (this PolygonStructSlice) Len() int {
	return len(this)
}
func (this PolygonStructSlice) Less(i, j int) bool {
	if this[i].isWart && !this[j].isWart {
		return true
	}
	if !this[i].isWart && this[j].isWart {
		return false
	}
	return math.Abs(polygonArea(this[i].RootPoly.polygonBeforeRotation)) > math.Abs(polygonArea(this[j].RootPoly.polygonBeforeRotation))
}

func (this PolygonStructSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

const ClipperBase_horizontal float64 = -9007199254740992
const ClipperBase_Skip int = -2
const ClipperBase_Unassigned int = -1

const ClipperBase_tolerance float64 = 1E-20
const ClipperBase_loRange int64 = 47453132
const ClipperBase_hiRange int64 = 0xfffffffffffff

type IntPolygon = []*IntPoint

func Abs(windCnt int64) int64 {
	if windCnt >= 0 {
		return windCnt
	} else {
		return -1 * windCnt
	}
}
func Max(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}
func Min(a, b int64) int64 {
	if a > b {
		return b
	} else {
		return a
	}
}

type OutPtStruct struct {
	Idx  int
	Pt   *IntPoint
	Prev *OutPtStruct
	Next *OutPtStruct
}
type OutFloatPtStruct struct {
	Idx  int
	Pt   *Point
	Prev *OutFloatPtStruct
	Next *OutFloatPtStruct
}

// func NewOutPt()*OutPtStruct{
// 	return &OutPtStruct{
// 		OutPt2: &TEdgeStruct{},
// 		OutPt1: &TEdgeStruct{},
// 		OffPt:&IntPoint{
// 		},
// 	}
// }
type OutRecStruct struct {
	Idx    int
	IsOpen bool
	IsHole bool

	PolyNode  *IntPoint
	BottomPt  *OutPtStruct
	Pts       *OutPtStruct
	FirstLeft *OutRecStruct
}
type JoinStruct struct {
	OutPt2 *OutPtStruct
	OutPt1 *OutPtStruct
	OffPt  *IntPoint
}

func NewJoin() *JoinStruct {
	return &JoinStruct{
		OutPt2: &OutPtStruct{},
		OutPt1: &OutPtStruct{},
		OffPt:  &IntPoint{},
	}
}

type Point struct {
	X      float64
	Y      float64
	marked bool
}

type IntPoint struct {
	X int64
	Y int64
}

func (this *IntPoint) Print() string {
	return fmt.Sprintf("{X:%d,Y:%d}", this.X, this.Y)
}
func IntPointReverse(a IntPolygon) IntPolygon {
	b := IntPolygon{}
	for index := range a {
		Point := &IntPoint{
			X: a[index].X,
			Y: a[index].Y,
		}
		b = append(b, Point)
	}
	for i := 0; i <= (len(b)/2)-1; i++ {
		bPoiont := &IntPoint{}
		bPoiont.X = b[len(b)-1-i].X
		bPoiont.Y = b[len(b)-1-i].Y
		b[len(b)-1-i].X = b[i].X
		b[len(b)-1-i].Y = b[i].Y
		b[i].X = bPoiont.X
		b[i].Y = bPoiont.Y
	}
	return b
}
func PolygonReverse(a Polygon) Polygon {
	b := Polygon{}
	for index := range a {
		Point := &Point{
			X: a[index].X,
			Y: a[index].Y,
		}
		b = append(b, Point)
	}
	for i := 0; i <= (len(b)/2)-1; i++ {
		bPoiont := &Point{}
		bPoiont.X = b[len(b)-1-i].X
		bPoiont.Y = b[len(b)-1-i].Y
		b[len(b)-1-i].X = b[i].X
		b[len(b)-1-i].Y = b[i].Y
		b[i].X = bPoiont.X
		b[i].Y = bPoiont.Y
	}
	return b
}

var ClipperOffset_def_arc_tolerance float64 = 0.25
var ClipperOffset_two_pi float64 = 6.28318530717959

type PolyType int8

const ptSubject PolyType = 0
const ptClip PolyType = 1

type ClipType int8

const ctIntersection ClipType = 0
const ctUnion ClipType = 1
const ctDifference ClipType = 2
const ctXor ClipType = 3

type PolyFillType = int64

const pftEvenOdd PolyFillType = 0
const pftNonZero PolyFillType = 1
const pftPositive PolyFillType = 2
const pftNegative PolyFillType = 3

type JoinType = int64

const jtSquare JoinType = 0
const jtRound JoinType = 1
const jtMiter JoinType = 2

type EndType = int64

const etOpenSquare EndType = 0
const etOpenRound EndType = 1
const etOpenButt EndType = 2
const etClosedLine EndType = 3
const etClosedPolygon EndType = 4

type IntersectNodeStruct struct {
	Edge1 *TEdgeStruct
	Edge2 *TEdgeStruct
	Pt    *IntPoint
}
type EdgeSide int8

const esLeft EdgeSide = 0
const esRight EdgeSide = 1

type Direction int8

const dRightToLeft Direction = 0
const dLeftToRight Direction = 1

type ConfigStruct struct {
	PaperSavePath  string
	ClipperScale   int64
	CurveTolerance float64
	ExploreConcave bool
	MutationRate   int
	PopulationSize int
	Rotations      int
	//Spacing        float64
	UseHoles        bool
	IfDebug         bool
	IfDraw          bool
	PartPartSpacing float64
	BinPartSpacing  float64

	LoopMaxNum   int
	RunTimeOut   int
	LengthWeight float64
	WidthWeight  float64
	MutilThread  int
}

const clipperScaleTimes int64 = 10000

var PublicConfig *ConfigStruct = &ConfigStruct{
	ClipperScale:   clipperScaleTimes, // 扩大倍数
	CurveTolerance: 0.3,
	ExploreConcave: false, //searchEdges
	MutationRate:   20,
	PopulationSize: 2,
	Rotations:      360,
	//	Spacing:        0,
	UseHoles:        false,
	IfDebug:         false,
	PartPartSpacing: 50,
	BinPartSpacing:  0,
}

type Path = []*Point
type PolyNodeStruct struct {
	m_Parent *PolyNodeStruct

	m_polygon  IntPolygon
	m_endtype  int64
	m_jointype int64
	m_Index    int

	m_Childs []*PolyNodeStruct

	IsOpen bool
}

func NewPolyNode() *PolyNodeStruct {
	return &PolyNodeStruct{
		m_Parent:   &PolyNodeStruct{},
		m_polygon:  IntPolygon{},
		m_endtype:  0,
		m_jointype: 0,
		m_Index:    0,
		m_Childs:   []*PolyNodeStruct{},
		IsOpen:     false,
	}
}
func (this *PolyNodeStruct) Childs() []*PolyNodeStruct {
	return this.m_Childs
}
func (this *PolyNodeStruct) AddChild(a *PolyNodeStruct) {
	b := len(this.m_Childs)
	this.m_Childs = append(this.m_Childs, a)
	a.m_Parent = this
	a.m_Index = b
}

type ClipperOffsetStruct struct {
	// "undefined" == typeof a && (a = 2);
	//     "undefined" == typeof b && (b = d.ClipperOffset.def_arc_tolerance);
	m_destPolys   []IntPolygon
	m_srcPoly     IntPolygon
	m_destPoly    IntPolygon
	m_normals     Polygon
	m_StepsPerRad float64

	m_miterLim float64
	m_cos      float64
	m_sin      float64
	m_sinA     float64
	m_delta    float64

	m_lowest *IntPoint

	m_polyNodes  *PolyNodeStruct
	MiterLimit   float64
	ArcTolerance float64
	// m_lowest.X int64

}

func NewClipperOffset(a float64, b float64) *ClipperOffsetStruct {
	this := &ClipperOffsetStruct{}
	this.m_destPolys = []IntPolygon{}
	this.m_srcPoly = IntPolygon{}
	this.m_destPoly = IntPolygon{}
	//this.m_normals = [];
	this.m_StepsPerRad = 0
	this.m_miterLim = 0
	this.m_cos = 0
	this.m_sin = 0
	this.m_sinA = 0
	this.m_delta = 0
	this.m_lowest = &IntPoint{}
	this.m_polyNodes = NewPolyNode()
	this.MiterLimit = a
	this.ArcTolerance = b
	this.m_lowest.X = -1
	return this
}

type IntRectStruct struct {
	left   int64
	top    int64
	right  int64
	bottom int64
}

func NewIntRect_4(a, b, c, d int64) *IntRectStruct {
	this := &IntRectStruct{}

	this.left = a
	this.top = b
	this.right = c
	this.bottom = d
	return this

}

type OneTouchStruct struct {
	Type int
	A    int
	B    int
}
type vectorsStruct struct {
	x     float64
	y     float64
	start *Point
	end   *Point
}
