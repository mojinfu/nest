package nest

import (
	"fmt"
	"log"
	"math"
	"sort"

	. "github.com/mojinfu/point"
)

func (this *ClipperStruct) logIfDebug(v ...interface{}) {
	if this.IfDebug {
		log.Println(v...)
	}
}
func (this *ClipperStruct) Execute(b ClipType, c, a PolyFillType) ([][]*IntPoint, bool) {
	if this.m_ExecuteLocked {
		return nil, false
	}
	if this.m_HasOpenPaths {
		this.logIfDebug("Error: PolyTree struct is need for open path clipping.")
	}
	this.m_ExecuteLocked = true
	this.m_SubjFillType = c
	this.m_ClipFillType = a
	this.m_ClipType = b
	this.m_UsingPolyTree = false
	e := [][]*IntPoint{}

	//okok
	//try {
	var f = this.ExecuteInternal()
	if f {
		e = this.BuildResult()
	}
	//} finally {
	this.DisposeAllPolyPts()
	this.m_ExecuteLocked = false
	//}
	return e, f
}
func (this *ClipperStruct) DisposeAllPolyPts() {
	a := 0
	b := len(this.m_PolyOuts)
	for ; a < b; a++ {
		this.DisposeOutRec(a)
	}
	this.m_PolyOuts = []*OutRecStruct{}
}
func (this *ClipperStruct) DisposeOutRec(a int) {
	b := this.m_PolyOuts[a]
	if nil != b.Pts {
		this.DisposeOutPts(b.Pts)
	}
	this.m_PolyOuts[a] = nil
}

func (this *ClipperStruct) DisposeOutPts(a *OutPtStruct) {
	if nil != a {
		for a.Prev.Next = nil; nil != a; {
			a = a.Next
		}
	}
}

func (this *ClipperStruct) PointCount(a *OutPtStruct) int {
	if nil == a {
		return 0
	}
	b := 0
	c := a
	b++
	c = c.Next
	for c != a {
		b++
		c = c.Next
	}
	return b
}
func (this *ClipperStruct) BuildResult() [][]*IntPoint {
	a := [][]*IntPoint{}
	b := 0
	c := len(this.m_PolyOuts)
	for ; b < c; b++ {
		eOutReg := this.m_PolyOuts[b]
		if nil != eOutReg.Pts {
			eOutPt := eOutReg.Pts.Prev
			f := this.PointCount(eOutPt)
			if !(2 > f) {
				g := []*IntPoint{}
				for h := 0; h < f; h++ {
					g = append(g, eOutPt.Pt)
					//g[h] = eOutPt.Pt
					eOutPt = eOutPt.Prev
				}
				a = append(a, g)

			}
		}
	}
	return a
}

func (this *ClipperStruct) BaseReset() {
	this.m_CurrentLM = this.m_MinimaList
	if nil != this.m_CurrentLM {
		var a = this.m_MinimaList
		for nil != a {
			var b = a.LeftBound
			if nil != b {
				b.Curr.X = b.Bot.X
				b.Curr.Y = b.Bot.Y
				b.Side = esLeft
				b.OutIdx = ClipperBase_Unassigned
			}
			b = a.RightBound
			if nil != b {
				b.Curr.X = b.Bot.X
				b.Curr.Y = b.Bot.Y
				b.Side = esRight
				b.OutIdx = ClipperBase_Unassigned
			}
			a = a.Next
		}
	}

}
func (this *ClipperStruct) Reset() {
	this.BaseReset()
	this.m_SortedEdges = nil
	this.m_ActiveEdges = nil
	this.m_Scanbeam = nil
	a := this.m_MinimaList
	for nil != a {
		this.InsertScanbeam(a.Y)
		a = a.Next
	}
}

func (this *ClipperStruct) GetBounds(a []IntPolygon) *IntRectStruct {
	b := 0
	c := len(a)
	for b < c && 0 == len(a[b]) {
		b++
	}

	if b == c {
		return NewIntRect_4(0, 0, 0, 0)
	}

	e := &IntRectStruct{}
	e.left = a[b][0].X
	e.right = e.left
	e.top = a[b][0].Y
	for e.bottom = e.top; b < c; b++ {
		f := 0
		g := len(a[b])
		for ; f < g; f++ {
			if a[b][f].X < e.left {
				e.left = a[b][f].X
			} else {
				if a[b][f].X > e.right {
					e.right = a[b][f].X
				}
			}
			if a[b][f].Y < e.top {
				e.top = a[b][f].Y
			} else {
				if a[b][f].Y > e.bottom {
					e.bottom = a[b][f].Y
				}
			}
		}

	}

	return e
}

func (this *ClipperStruct) InsertScanbeam(a int64) {
	if nil == this.m_Scanbeam {
		this.m_Scanbeam = &ScanbeamSturct{}
		this.m_Scanbeam.Next = nil
		this.m_Scanbeam.Y = a
	} else if a > this.m_Scanbeam.Y {
		b := &ScanbeamSturct{}
		b.Y = a
		b.Next = this.m_Scanbeam
		this.m_Scanbeam = b
	} else {
		c := this.m_Scanbeam
		for nil != c.Next && a <= c.Next.Y {
			c = c.Next
		}
		if a != c.Y {
			b := &ScanbeamSturct{}
			b.Y = a
			b.Next = c.Next
			c.Next = b
		}
	}
	return
}

type ScanbeamSturct struct {
	Next *ScanbeamSturct
	Y    int64
}

func (this *ClipperStruct) PopScanbeam() int64 {
	a := this.m_Scanbeam.Y
	this.m_Scanbeam = this.m_Scanbeam.Next
	return a
}
func (this *ClipperStruct) PopLocalMinima() {
	if nil != this.m_CurrentLM {
		this.m_CurrentLM = this.m_CurrentLM.Next
	}
}
func (this *ClipperStruct) E2InsertsBeforeE1(a *TEdgeStruct, b *TEdgeStruct) bool {
	if b.Curr.X == a.Curr.X {
		if b.Top.Y > a.Top.Y {
			return b.Top.X < this.TopX(a, b.Top.Y)
		} else {
			return a.Top.X > this.TopX(b, a.Top.Y)
		}
	} else {
		return b.Curr.X < a.Curr.X
	}
}
func (this *ClipperStruct) TopX(a *TEdgeStruct, b int64) int64 {
	if b == a.Top.Y {
		return a.Top.X
	} else {
		return a.Bot.X + int64(math.Round(a.Dx*float64(b-a.Bot.Y)))
	}
}
func (this *ClipperStruct) InsertEdgeIntoAEL(a *TEdgeStruct, b *TEdgeStruct) {
	if nil == this.m_ActiveEdges {
		a.PrevInAEL = nil
		a.NextInAEL = nil
		this.m_ActiveEdges = a
	} else if nil == b && this.E2InsertsBeforeE1(this.m_ActiveEdges, a) {
		a.PrevInAEL = nil
		a.NextInAEL = this.m_ActiveEdges
		this.m_ActiveEdges.PrevInAEL = a
		this.m_ActiveEdges = a

	} else {
		if nil == b {
			b = this.m_ActiveEdges
		}
		for nil != b.NextInAEL && !this.E2InsertsBeforeE1(b.NextInAEL, a) {
			b = b.NextInAEL
		}

		a.NextInAEL = b.NextInAEL
		if nil != b.NextInAEL {
			b.NextInAEL.PrevInAEL = a
		}
		a.PrevInAEL = b
		b.NextInAEL = a
	}
}
func (this *ClipperStruct) SetWindingCount(a *TEdgeStruct) {
	b := a.PrevInAEL
	for nil != b && (b.PolyTyp != a.PolyTyp || 0 == b.WindDelta) {
		b = b.PrevInAEL
	}

	if nil == b {
		if 0 == a.WindDelta {
			a.WindCnt = 1
		} else {
			a.WindCnt = a.WindDelta
		}

		a.WindCnt2 = 0
		b = this.m_ActiveEdges
	} else {
		if 0 == a.WindDelta && this.m_ClipType != ctUnion {
			a.WindCnt = 1
		} else if this.IsEvenOddFillType(a) {
			if 0 == a.WindDelta {
				c := true
				e := b.PrevInAEL
				for nil != e {
					if e.PolyTyp == b.PolyTyp && 0 != e.WindDelta {
						c = !c
					}
					e = e.PrevInAEL
				}
				if c {
					a.WindCnt = 0
				} else {
					a.WindCnt = 1
				}

			} else {
				a.WindCnt = a.WindDelta
			}

		} else {

			if b.WindCnt*b.WindDelta < 0 {
				if Abs(b.WindCnt) > 1 {
					if b.WindDelta*a.WindDelta < 0 {
						a.WindCnt = b.WindCnt
					} else {
						a.WindCnt = b.WindCnt + a.WindDelta
					}
				} else {
					a.WindCnt = b.WindCnt + b.WindDelta + a.WindDelta
				}
			} else {
				if Abs(b.WindCnt) > 1 && b.WindDelta*a.WindDelta < 0 {
					a.WindCnt = b.WindCnt
				} else if b.WindCnt+a.WindDelta == 0 {
					a.WindCnt = b.WindCnt
				} else {
					a.WindCnt = b.WindCnt + a.WindDelta
				}
			}

		}
		a.WindCnt2 = b.WindCnt2
		b = b.NextInAEL
	}
	if this.IsEvenOddAltFillType(a) {
		for b != a {
			if 0 != b.WindDelta {
				if 0 == a.WindCnt2 {
					a.WindCnt2 = 1
				} else {
					a.WindCnt2 = 0
				}

				b = b.NextInAEL
			}
		}
	} else {
		for b != a {
			a.WindCnt2 += b.WindDelta
			b = b.NextInAEL
		}
	}
}
func (this *ClipperStruct) IsEvenOddAltFillType(a *TEdgeStruct) bool {
	if a.PolyTyp == ptSubject {
		return this.m_ClipFillType == pftEvenOdd
	} else {
		return this.m_SubjFillType == pftEvenOdd
	}
}
func (this *ClipperStruct) IsContributing(a *TEdgeStruct) bool {
	var b, c PolyFillType
	if a.PolyTyp == ptSubject {
		b = this.m_SubjFillType
		c = this.m_ClipFillType
	} else {
		b = this.m_ClipFillType
		c = this.m_SubjFillType
	}
	switch b {
	case pftEvenOdd:
		if 0 == a.WindDelta && 1 != a.WindCnt {
			return false
		}

		break
	case pftNonZero:
		if 1 != Abs(a.WindCnt) {
			return false
		}

		break
	case pftPositive:
		if 1 != a.WindCnt {
			return false
		}

		break
	default:
		if -1 != a.WindCnt {
			return false
		}
	}
	switch this.m_ClipType {
	case ctIntersection:
		switch c {
		case pftEvenOdd:
		case pftNonZero:
			return 0 != a.WindCnt2
		case pftPositive:
			return 0 < a.WindCnt2
		default:
			return 0 > a.WindCnt2
		}
	case ctUnion:
		switch c {
		case pftEvenOdd:
		case pftNonZero:
			return 0 == a.WindCnt2
		case pftPositive:
			return 0 >= a.WindCnt2
		default:
			return 0 <= a.WindCnt2
		}
	case ctDifference:
		if a.PolyTyp == ptSubject {
			switch c {
			case pftEvenOdd:
				{

				}
			case pftNonZero:
				return 0 == a.WindCnt2
			case pftPositive:
				return 0 >= a.WindCnt2
			default:
				return 0 <= a.WindCnt2
			}
		} else {
			switch c {
			case pftEvenOdd:
			case pftNonZero:
				return 0 != a.WindCnt2
			case pftPositive:
				return 0 < a.WindCnt2
			default:
				return 0 > a.WindCnt2
			}
		}
	case ctXor:
		if 0 == a.WindDelta {
			switch c {
			case pftEvenOdd:
			case pftNonZero:
				return 0 == a.WindCnt2
			case pftPositive:
				return 0 >= a.WindCnt2
			default:
				return 0 <= a.WindCnt2
			}
		}
	}
	return true
}

func (this *ClipperStruct) AddLocalMinPoly(a, b *TEdgeStruct, c *IntPoint) *OutPtStruct {
	var e *OutPtStruct = &OutPtStruct{}
	var f *TEdgeStruct = &TEdgeStruct{}
	if this.IsHorizontal(b) || a.Dx > b.Dx {
		e = this.AddOutPt(a, c)
		b.OutIdx = a.OutIdx
		a.Side = esLeft
		b.Side = esRight
		f = a
		if f.PrevInAEL == b {
			a = b.PrevInAEL
		} else {
			a = f.PrevInAEL
		}
	} else {
		e = this.AddOutPt(b, c)
		a.OutIdx = b.OutIdx
		a.Side = esRight
		b.Side = esLeft
		f = b
		if f.PrevInAEL == a {
			a = a.PrevInAEL
		} else {
			a = f.PrevInAEL
		}
	}

	if nil != a && 0 <= a.OutIdx && this.TopX(a, c.Y) == this.TopX(f, c.Y) && this.SlopesEqual_3(f, a, this.m_UseFullRange) && 0 != f.WindDelta && 0 != a.WindDelta {
		cOutPt := this.AddOutPt(a, c)
		this.AddJoin(e, cOutPt, f.Top)
	}
	return e
}

func (this *ClipperStruct) IsHorizontal(a *TEdgeStruct) bool {
	return 0 == a.Delta.Y
}
func (this *ClipperStruct) AddJoin(a, b *OutPtStruct, c *IntPoint) {
	e := NewJoin()
	e.OutPt1 = a
	e.OutPt2 = b
	e.OffPt.X = c.X
	e.OffPt.Y = c.Y
	this.m_Joins = append(this.m_Joins, e)
}

func (this *ClipperStruct) CreateOutRec() *OutRecStruct {
	a := &OutRecStruct{}
	a.Idx = -1
	a.IsHole = false
	a.IsOpen = false
	a.FirstLeft = nil
	a.Pts = nil
	a.BottomPt = nil
	a.PolyNode = nil
	this.m_PolyOuts = append(this.m_PolyOuts, a)
	a.Idx = len(this.m_PolyOuts) - 1
	return a
}
func (this *ClipperStruct) AddOutPt(a *TEdgeStruct, b *IntPoint) *OutPtStruct {
	f := &OutPtStruct{
		Pt: &IntPoint{},
	}
	c := (a.Side == esLeft)
	if 0 > a.OutIdx {
		var e = this.CreateOutRec()
		e.IsOpen = (0 == a.WindDelta)

		e.Pts = f
		f.Idx = e.Idx
		f.Pt.X = b.X
		f.Pt.Y = b.Y
		f.Next = f
		f.Prev = f
		if !e.IsOpen {
			this.SetHoleState(a, e)
		}
		a.OutIdx = e.Idx
	} else {
		e := this.m_PolyOuts[a.OutIdx]
		g := e.Pts
		if c && op_Equality_IntPoint(b, g.Pt) {
			return g
		}
		if !c && op_Equality_IntPoint(b, g.Prev.Pt) {
			return g.Prev
		}

		f.Idx = e.Idx
		f.Pt.X = b.X
		f.Pt.Y = b.Y
		f.Next = g
		f.Prev = g.Prev
		f.Prev.Next = f
		g.Prev = f
		if c {
			e.Pts = f
		}
	}
	return f
}
func (this *ClipperStruct) SetHoleState(a *TEdgeStruct, b *OutRecStruct) {
	c := false
	e := a.PrevInAEL
	for nil != e {
		if 0 <= e.OutIdx && 0 != e.WindDelta {
			c = !c
			if nil == b.FirstLeft {
				b.FirstLeft = this.m_PolyOuts[e.OutIdx]
			}
		}
		e = e.PrevInAEL
	}
	if c {
		b.IsHole = true
	}
}
func (this *ClipperStruct) AddEdgeToSEL(a *TEdgeStruct) {
	if nil == this.m_SortedEdges {
		this.m_SortedEdges = a
		a.PrevInSEL = nil
		a.NextInSEL = nil
	} else {
		a.NextInSEL = this.m_SortedEdges
		a.PrevInSEL = nil

		//this.m_SortedEdges = this.m_SortedEdges.PrevInSEL = a

		this.m_SortedEdges.PrevInSEL = a
		this.m_SortedEdges = a
	}
}

func (this *ClipperStruct) HorzSegmentsOverlap(a, b, c, e *IntPoint) bool {
	if (a.X > c.X) == (a.X < e.X) {
		return true
	} else {
		if (b.X > c.X) == (b.X < e.X) {
			return true
		} else {
			if (c.X > a.X) == (c.X < b.X) {
				return true
			} else {
				if (e.X > a.X) == (e.X < b.X) {
					return true
				} else {
					if a.X == c.X && b.X == e.X {
						return true
					} else {
						if a.X == e.X && b.X == c.X {
							return true
						} else {
							return false
						}
					}
				}
			}
		}

	}

}
func (this *ClipperStruct) InsertLocalMinimaIntoAEL(a int64) {
	for nil != this.m_CurrentLM && this.m_CurrentLM.Y == a {
		b := this.m_CurrentLM.LeftBound
		c := this.m_CurrentLM.RightBound
		this.PopLocalMinima()
		var e *OutPtStruct = nil
		if nil == b {
			this.InsertEdgeIntoAEL(c, nil)
			this.SetWindingCount(c)
			if this.IsContributing(c) {
				e = this.AddOutPt(c, c.Bot)
			}
		} else {
			if nil == c {
				this.InsertEdgeIntoAEL(b, nil)
				this.SetWindingCount(b)
				if this.IsContributing(b) {
					e = this.AddOutPt(b, b.Bot)
				}
			} else {
				this.InsertEdgeIntoAEL(b, nil)
				this.InsertEdgeIntoAEL(c, b)
				this.SetWindingCount(b)
				c.WindCnt = b.WindCnt
				c.WindCnt2 = b.WindCnt2 //okok
				if this.IsContributing(b) {
					e = this.AddLocalMinPoly(b, c, b.Bot)
				}
			}
			this.InsertScanbeam(b.Top.Y)
		}
		if nil != c { //okok
			if this.IsHorizontal(c) {
				this.AddEdgeToSEL(c)
			} else {
				this.InsertScanbeam(c.Top.Y)
			}
		}
		if nil != b && nil != c {
			if nil != e && this.IsHorizontal(c) && 0 < len(this.m_GhostJoins) && 0 != c.WindDelta {
				f := 0
				g := len(this.m_GhostJoins)
				for ; f < g; f++ {
					h := this.m_GhostJoins[f]

					if this.HorzSegmentsOverlap(h.OutPt1.Pt, h.OffPt, c.Bot, c.Top) {
						this.AddJoin(h.OutPt1, e, h.OffPt)
					}
				}
			}
			if 0 <= b.OutIdx && nil != b.PrevInAEL && b.PrevInAEL.Curr.X == b.Bot.X && 0 <= b.PrevInAEL.OutIdx {
				if this.SlopesEqual_3(b.PrevInAEL, b, this.m_UseFullRange) && 0 != b.WindDelta && 0 != b.PrevInAEL.WindDelta {
					f := this.AddOutPt(b.PrevInAEL, b.Bot)
					this.AddJoin(e, f, b.Top)
				}
			}

			if b.NextInAEL != c {
				if 0 <= c.OutIdx && 0 <= c.PrevInAEL.OutIdx && this.SlopesEqual_3(c.PrevInAEL, c, this.m_UseFullRange) && 0 != c.WindDelta && 0 != c.PrevInAEL.WindDelta {
					f := this.AddOutPt(c.PrevInAEL, c.Bot)
					this.AddJoin(e, f, c.Top)
				}
				eTEdge := b.NextInAEL
				if nil != eTEdge {
					for eTEdge != c {
						this.IntersectEdges(c, eTEdge, b.Curr, false)
						eTEdge = eTEdge.NextInAEL
					}
				}
			}
		}
	}
}
func (this *ClipperStruct) Param1RightOfParam2(outRec1, outRec2 *OutRecStruct) bool {

	outRec1 = outRec1.FirstLeft
	if outRec1 == outRec2 {
		return true
	}
	for outRec1 != nil {
		outRec1 = outRec1.FirstLeft
		if outRec1 == outRec2 {
			return true
		}
	}
	return false
}
func (this *ClipperStruct) GetDx(pt1, pt2 *IntPoint) float64 {
	if pt1.Y == pt2.Y {
		return ClipperBase_horizontal
	} else {
		return float64(pt2.X-pt1.X) / float64(pt2.Y-pt1.Y)
	}
}
func (this *ClipperStruct) FirstIsBottomPt(btmPt1, btmPt2 *OutPtStruct) bool {
	p := btmPt1.Prev
	for op_Equality_IntPoint(p.Pt, btmPt1.Pt) && (p != btmPt1) {
		p = p.Prev
	}

	var dx1p = math.Abs(this.GetDx(btmPt1.Pt, p.Pt))
	p = btmPt1.Next
	for op_Equality_IntPoint(p.Pt, btmPt1.Pt) && (p != btmPt1) {
		p = p.Next
	}

	dx1n := math.Abs(this.GetDx(btmPt1.Pt, p.Pt))
	p = btmPt2.Prev
	for op_Equality_IntPoint(p.Pt, btmPt2.Pt) && (p != btmPt2) {
		p = p.Prev
	}

	dx2p := math.Abs(this.GetDx(btmPt2.Pt, p.Pt))
	p = btmPt2.Next
	for op_Equality_IntPoint(p.Pt, btmPt2.Pt) && (p != btmPt2) {
		p = p.Next
	}

	dx2n := math.Abs(this.GetDx(btmPt2.Pt, p.Pt))
	return (dx1p >= dx2p && dx1p >= dx2n) || (dx1n >= dx2p && dx1n >= dx2n)
}
func (this *ClipperStruct) GetBottomPt(a *OutPtStruct) *OutPtStruct {
	var b *OutPtStruct = nil
	c := a.Next
	for c != a {
		if c.Pt.Y > a.Pt.Y {
			a = c
			b = nil
		} else {
			if c.Pt.Y == a.Pt.Y && c.Pt.X <= a.Pt.X {
				if c.Pt.X < a.Pt.X {
					b = nil
					a = c
				} else {
					if c.Next != a && c.Prev != a {
						b = c
					}
				}
			}
		}
		c = c.Next
	}

	if nil != b {
		for b != c {
			if !this.FirstIsBottomPt(c, b) {
				a = b
			}
			b = b.Next
			for op_InEquality_IntPoint(b.Pt, a.Pt) {
				b = b.Next
			}

		}

	}

	return a

}
func (this *ClipperStruct) GetLowermostRec(outRec1, outRec2 *OutRecStruct) *OutRecStruct {

	if nil == outRec1.BottomPt {
		outRec1.BottomPt = this.GetBottomPt(outRec1.Pts)
	}
	if nil == outRec2.BottomPt {
		outRec2.BottomPt = this.GetBottomPt(outRec2.Pts)
	}

	var bPt1 = outRec1.BottomPt
	var bPt2 = outRec2.BottomPt
	if bPt1.Pt == nil {
		bPt1.Pt = &IntPoint{}
	}
	if bPt2.Pt == nil {
		bPt2.Pt = &IntPoint{}
	}
	if bPt1.Pt.Y > bPt2.Pt.Y {
		return outRec1
	} else if bPt1.Pt.Y < bPt2.Pt.Y {
		return outRec2
	} else if bPt1.Pt.X < bPt2.Pt.X {
		return outRec1
	} else if bPt1.Pt.X > bPt2.Pt.X {
		return outRec2
	} else if bPt1.Next == bPt1 {
		return outRec2
	} else if bPt2.Next == bPt2 {
		return outRec1
	} else if this.FirstIsBottomPt(bPt1, bPt2) {
		return outRec1
	} else {
		return outRec2
	}
}
func (this *ClipperStruct) ReversePolyPtLinks(pp *OutPtStruct) {
	if pp == nil {
		return
	}
	var pp1 *OutPtStruct
	var pp2 *OutPtStruct
	pp1 = pp

	pp2 = pp1.Next
	pp1.Next = pp1.Prev
	pp1.Prev = pp2
	pp1 = pp2

	for pp1 != pp {
		pp2 = pp1.Next
		pp1.Next = pp1.Prev
		pp1.Prev = pp2
		pp1 = pp2
	}
}
func (this *ClipperStruct) AppendPolygon(e1 *TEdgeStruct, e2 *TEdgeStruct) {
	var outRec1 = this.m_PolyOuts[e1.OutIdx]
	var outRec2 = this.m_PolyOuts[e2.OutIdx]
	var holeStateRec *OutRecStruct
	if this.Param1RightOfParam2(outRec1, outRec2) {
		holeStateRec = outRec2
	} else if this.Param1RightOfParam2(outRec2, outRec1) {
		holeStateRec = outRec1
	} else {
		holeStateRec = this.GetLowermostRec(outRec1, outRec2)
	}
	var p1_lft = outRec1.Pts
	var p1_rt = p1_lft.Prev
	var p2_lft = outRec2.Pts
	var p2_rt = p2_lft.Prev
	var side EdgeSide
	if e1.Side == esLeft {
		if e2.Side == esLeft {
			this.ReversePolyPtLinks(p2_lft)
			p2_lft.Next = p1_lft
			p1_lft.Prev = p2_lft
			p1_rt.Next = p2_rt
			p2_rt.Prev = p1_rt
			outRec1.Pts = p2_rt
		} else {
			p2_rt.Next = p1_lft
			p1_lft.Prev = p2_rt
			p2_lft.Prev = p1_rt
			p1_rt.Next = p2_lft
			outRec1.Pts = p2_lft
		}
		side = esLeft
	} else {
		if e2.Side == esRight {
			this.ReversePolyPtLinks(p2_lft)
			p1_rt.Next = p2_rt
			p2_rt.Prev = p1_rt
			p2_lft.Next = p1_lft
			p1_lft.Prev = p2_lft
		} else {
			p1_rt.Next = p2_lft
			p2_lft.Prev = p1_rt
			p1_lft.Prev = p2_rt
			p2_rt.Next = p1_lft
		}
		side = esRight
	}
	if holeStateRec == outRec2 {
		outRec1.BottomPt = outRec2.BottomPt
		if outRec1.BottomPt == nil {
			outRec1.BottomPt = &OutPtStruct{}
		} //saya
		outRec1.BottomPt.Idx = outRec1.Idx
		if outRec2.FirstLeft != outRec1 {
			outRec1.FirstLeft = outRec2.FirstLeft
		}
		outRec1.IsHole = outRec2.IsHole
	}
	outRec2.Pts = nil
	outRec2.BottomPt = nil
	outRec2.FirstLeft = outRec1

	var OKIdx = e1.OutIdx
	var ObsoleteIdx = e2.OutIdx
	e1.OutIdx = -1
	e2.OutIdx = -1
	var e = this.m_ActiveEdges
	for e != nil {
		if e.OutIdx == ObsoleteIdx {
			e.OutIdx = OKIdx
			e.Side = side
			break
		}
		e = e.NextInAEL
	}
	outRec2.Idx = outRec1.Idx
}
func (this *ClipperStruct) AddLocalMaxPoly(e1 *TEdgeStruct, e2 *TEdgeStruct, pt *IntPoint) {
	this.AddOutPt(e1, pt)
	if e1.OutIdx == e2.OutIdx {
		e1.OutIdx = -1
		e2.OutIdx = -1
	} else if e1.OutIdx < e2.OutIdx {
		this.AppendPolygon(e1, e2)
	} else {
		this.AppendPolygon(e2, e1)
	}
}
func (this *ClipperStruct) DeleteFromAEL(e *TEdgeStruct) {
	var AelPrev = e.PrevInAEL
	var AelNext = e.NextInAEL
	if AelPrev == nil && AelNext == nil && (e != this.m_ActiveEdges) {
		return
	}
	if AelPrev != nil {
		AelPrev.NextInAEL = AelNext
	} else {
		this.m_ActiveEdges = AelNext
	}
	if AelNext != nil {
		AelNext.PrevInAEL = AelPrev
	}
	e.NextInAEL = nil
	e.PrevInAEL = nil
}

func (this *ClipperStruct) IntersectEdges(a, b *TEdgeStruct, c *IntPoint, e bool) {

	f := (!e && nil == a.NextInLML && a.Top.X == c.X && a.Top.Y == c.Y)
	e = (!e && nil == b.NextInLML && b.Top.X == c.X && b.Top.Y == c.Y)
	g := (0 <= a.OutIdx)
	h := (0 <= b.OutIdx)
	if 0 == a.WindDelta || 0 == b.WindDelta {

		if 0 == a.WindDelta && 0 == b.WindDelta {
			if (f || e) && g && h {
				this.AddLocalMaxPoly(a, b, c)
			}

		} else {
			if a.PolyTyp == b.PolyTyp && a.WindDelta != b.WindDelta && this.m_ClipType == ctUnion {

				if 0 == a.WindDelta {
					if h {
						this.AddOutPt(a, c)
						if g {
							a.OutIdx = -1
						}
					}
				} else {
					if g {
						this.AddOutPt(b, c)
						if h {
							b.OutIdx = -1
						}
					}
				}
			} else {
				if a.PolyTyp != b.PolyTyp {

					if 0 != a.WindDelta || 1 != Abs(b.WindCnt) || (this.m_ClipType == ctUnion && 0 != b.WindCnt2) {

						if !(0 != b.WindDelta || 1 != Abs(a.WindCnt) || (this.m_ClipType == ctUnion && 0 != a.WindCnt2)) {

							this.AddOutPt(b, c)
							if h {
								b.OutIdx = -1
							}

						}
					} else {

						this.AddOutPt(a, c)
						if g {
							a.OutIdx = -1
						}

					}

				}

			}

		}

		if f {
			if 0 > a.OutIdx {
				this.DeleteFromAEL(a)
			} else {
				this.logIfDebug("Error intersecting polylines")
			}

		}

		if e {
			if 0 > b.OutIdx {
				this.DeleteFromAEL(b)
			} else {
				this.logIfDebug("Error intersecting polylines")
			}

		}

	} else {
		if a.PolyTyp == b.PolyTyp {
			if this.IsEvenOddFillType(a) {
				var l = a.WindCnt
				a.WindCnt = b.WindCnt
				b.WindCnt = l
			} else if 0 == a.WindCnt+b.WindDelta {
				a.WindCnt = -a.WindCnt
			} else {
				a.WindCnt = a.WindCnt + b.WindDelta
			}
			if 0 == b.WindCnt-a.WindDelta {
				b.WindCnt = -b.WindCnt
			} else {
				b.WindCnt = b.WindCnt - a.WindDelta
			}

		} else {
			if this.IsEvenOddFillType(b) {
				if 0 == a.WindCnt2 {
					a.WindCnt2 = 1
				} else {
					a.WindCnt2 = 0
				}
			} else {
				a.WindCnt2 += b.WindDelta
			}
			if this.IsEvenOddFillType(a) {
				if 0 == b.WindCnt2 {
					b.WindCnt2 = 1
				} else {
					b.WindCnt2 = 0
				}
			} else {
				b.WindCnt2 -= a.WindDelta
			}
		}
		var k, n, m, l PolyFillType
		if a.PolyTyp == ptSubject {
			k = this.m_SubjFillType
			m = this.m_ClipFillType
		} else {
			k = this.m_ClipFillType
			m = this.m_SubjFillType
		}

		if b.PolyTyp == ptSubject {
			n = this.m_SubjFillType
			l = this.m_ClipFillType
		} else {
			n = this.m_ClipFillType
			l = this.m_SubjFillType
		}
		switch k {
		case pftPositive:
			k = PolyFillType(a.WindCnt)
			break
		case pftNegative:
			k = PolyFillType(-a.WindCnt)
			break
		default:
			k = PolyFillType(Abs(a.WindCnt))
		}
		switch n {
		case pftPositive:
			n = PolyFillType(b.WindCnt)
			break
		case pftNegative:
			n = PolyFillType(-b.WindCnt)
			break
		default:
			n = PolyFillType(Abs(b.WindCnt))
		}
		if g && h {
			if f || e || 0 != k && 1 != k || 0 != n && 1 != n || a.PolyTyp != b.PolyTyp && this.m_ClipType != ctXor {
				this.AddLocalMaxPoly(a, b, c)

			} else {
				this.AddOutPt(a, c)
				this.AddOutPt(b, c)
				this.SwapSides(a, b)
				this.SwapPolyIndexes(a, b)
			}
		} else if g {
			if 0 == n || 1 == n {
				this.AddOutPt(a, c)
				this.SwapSides(a, b)
				this.SwapPolyIndexes(a, b)
			}
		} else if h {
			if 0 == k || 1 == k {
				this.AddOutPt(b, c)
				this.SwapSides(a, b)
				this.SwapPolyIndexes(a, b)
			}

		} else if !(0 != k && 1 != k || 0 != n && 1 != n || f || e) {
			var g, h PolyFillType
			switch m {
			case pftPositive:
				g = a.WindCnt2
				break
			case pftNegative:
				g = -a.WindCnt2
				break
			default:
				g = Abs(a.WindCnt2)
			}
			switch l {
			case pftPositive:
				h = b.WindCnt2
				break
			case pftNegative:
				h = -b.WindCnt2
				break
			default:
				h = Abs(b.WindCnt2)
			}
			if a.PolyTyp != b.PolyTyp {
				this.AddLocalMinPoly(a, b, c)
			} else if 1 == k && 1 == n {
				switch this.m_ClipType {
				case ctIntersection:
					if 0 < g && 0 < h {
						this.AddLocalMinPoly(a, b, c)
					}
					break
				case ctUnion:
					if 0 >= g && 0 >= h {
						this.AddLocalMinPoly(a, b, c)
					}
					break
				case ctDifference:
					if a.PolyTyp == ptClip && 0 < g && 0 < h || a.PolyTyp == ptSubject && 0 >= g && 0 >= h {
						this.AddLocalMinPoly(a, b, c)
					}
					break
				case ctXor:
					this.AddLocalMinPoly(a, b, c)
				}
			} else {
				this.SwapSides(a, b)
			}

		}
		if f != e && (f && 0 <= a.OutIdx || e && 0 <= b.OutIdx) {
			this.SwapSides(a, b)
			this.SwapPolyIndexes(a, b)

		}
		if f {
			this.DeleteFromAEL(a)
		}
		if e {
			this.DeleteFromAEL(b)
		}
	}
}
func (this *ClipperStruct) SwapSides(edge1, edge2 *TEdgeStruct) {
	side := edge1.Side
	edge1.Side = edge2.Side
	edge2.Side = side
}
func (this *ClipperStruct) SwapPolyIndexes(edge1, edge2 *TEdgeStruct) {
	outIdx := edge1.OutIdx
	edge1.OutIdx = edge2.OutIdx
	edge2.OutIdx = outIdx
}
func (this *ClipperStruct) IsEvenOddFillType(edge *TEdgeStruct) bool {
	if edge.PolyTyp == ptSubject {
		return this.m_SubjFillType == pftEvenOdd
	} else {
		return this.m_ClipFillType == pftEvenOdd
	}
}
func (this *ClipperStruct) GetHorzDirection(a *TEdgeStruct, b *CProcHStruct) {
	if a.Bot.X < a.Top.X {
		b.Left = a.Bot.X
		b.Right = a.Top.X
		b.Dir = dLeftToRight
	} else {
		b.Left = a.Top.X
		b.Right = a.Bot.X
		b.Dir = dRightToLeft
	}
}

type CProcHStruct struct {
	Dir   Direction
	Left  int64
	Right int64
}

func (this *ClipperStruct) GetMaximaPair(a *TEdgeStruct) *TEdgeStruct {
	var b *TEdgeStruct = nil
	if op_Equality_IntPoint(a.Next.Top, a.Top) && nil == a.Next.NextInLML {
		b = a.Next
	} else {
		if op_Equality_IntPoint(a.Prev.Top, a.Top) && nil == a.Prev.NextInLML {
			b = a.Prev
		}
	}
	if nil == b || -2 != b.OutIdx && (b.NextInAEL != b.PrevInAEL || this.IsHorizontal(b)) {
		return b
	} else {
		return nil
	}
}
func (this *ClipperStruct) GetNextInAEL(a *TEdgeStruct, b Direction) *TEdgeStruct {
	if b == dLeftToRight {
		return a.NextInAEL
	} else {
		return a.PrevInAEL
	}
}

func (this *ClipperStruct) AddGhostJoin(a *OutPtStruct, b *IntPoint) {
	c := &JoinStruct{
		OffPt: &IntPoint{},
	}
	c.OutPt1 = a
	c.OffPt.X = b.X
	c.OffPt.Y = b.Y
	this.m_GhostJoins = append(this.m_GhostJoins, c)

}
func (this *ClipperStruct) PrepareHorzJoins(a *TEdgeStruct, b bool) {
	c := this.m_PolyOuts[a.OutIdx].Pts
	if a.Side != esLeft {
		c = c.Prev
	}
	if b {
		if op_Equality_IntPoint(c.Pt, a.Top) {
			this.AddGhostJoin(c, a.Bot)
		} else {
			this.AddGhostJoin(c, a.Top)
		}
	}

}
func (this *ClipperStruct) SwapPositionsInAEL(a, b *TEdgeStruct) {
	if a.NextInAEL != a.PrevInAEL && b.NextInAEL != b.PrevInAEL {
		if a.NextInAEL == b {
			c := b.NextInAEL
			if nil != c {
				c.PrevInAEL = a
			}
			e := a.PrevInAEL
			if nil != e {
				e.NextInAEL = b
			}
			b.PrevInAEL = e
			b.NextInAEL = a
			a.PrevInAEL = b
			a.NextInAEL = c
		} else if b.NextInAEL == a {
			c := a.NextInAEL
			if nil != c {
				c.PrevInAEL = b
			}
			e := b.PrevInAEL
			if nil != e {
				e.NextInAEL = a
			}
			a.PrevInAEL = e
			a.NextInAEL = b
			b.PrevInAEL = a
			b.NextInAEL = c
		} else {
			c := a.NextInAEL
			e := a.PrevInAEL
			a.NextInAEL = b.NextInAEL
			if nil != a.NextInAEL {
				a.NextInAEL.PrevInAEL = a
			}
			a.PrevInAEL = b.PrevInAEL
			if nil != a.PrevInAEL {
				a.PrevInAEL.NextInAEL = a
			}
			b.NextInAEL = c
			if nil != b.NextInAEL {
				b.NextInAEL.PrevInAEL = b
			}
			b.PrevInAEL = e
			if nil != b.PrevInAEL {
				b.PrevInAEL.NextInAEL = b
			}
		}

		if nil == a.PrevInAEL {
			this.m_ActiveEdges = a
		} else {
			if nil == b.PrevInAEL {
				this.m_ActiveEdges = b
			}
		}
	}
}

// func (this *ClipperStruct) AddOutPt(a  *TEdgeStruct,b *IntPoint ) *OutPtStruct {
// 	 c := (a.Side == esLeft)
//         if (0 > a.OutIdx) {
//             e := this.CreateOutRec()
//             e.IsOpen = (0 == a.WindDelta)
//              f := &OutPtStruct{}
//             e.Pts = f
//             f.Idx = e.Idx
//             f.Pt.X = b.X
//             f.Pt.Y = b.Y
//             f.Next = f
//             f.Prev = f
//            if  !e.IsOpen {
// 			this.SetHoleState(a, e)
// 		   }
//             a.OutIdx = e.Idx
//         } else {
//              e := this.m_PolyOuts[a.OutIdx]
//                 g := e.Pts
//             if (c && op_Equality_IntPoint(b, g.Pt)){
// 				return g
// 			}

//             if (!c && op_Equality_IntPoint(b, g.Prev.Pt)){
// 				return g.Prev
// 			}

//             f :=&OutPtStruct{}
//             f.Idx = e.Idx
//             f.Pt.X = b.X
//             f.Pt.Y = b.Y
//             f.Next = g
//             f.Prev = g.Prev
//             f.Prev.Next = f
//             g.Prev = f
//             if c {
// 				e.Pts = f
// 			}
//         }
//         return f
// }
func (this *ClipperStruct) UpdateEdgeIntoAEL(a *TEdgeStruct) *TEdgeStruct {
	if nil == a.NextInLML {
		this.logIfDebug("UpdateEdgeIntoAEL: invalid call")
	}
	b := a.PrevInAEL
	c := a.NextInAEL
	a.NextInLML.OutIdx = a.OutIdx
	if nil != b {
		b.NextInAEL = a.NextInLML
	} else {
		this.m_ActiveEdges = a.NextInLML
	}
	if nil != c {
		c.PrevInAEL = a.NextInLML
	}
	a.NextInLML.Side = a.Side
	a.NextInLML.WindDelta = a.WindDelta
	a.NextInLML.WindCnt = a.WindCnt
	a.NextInLML.WindCnt2 = a.WindCnt2
	a = a.NextInLML
	a.Curr.X = a.Bot.X
	a.Curr.Y = a.Bot.Y
	a.PrevInAEL = b
	a.NextInAEL = c
	if !this.IsHorizontal(a) {
		this.InsertScanbeam(a.Top.Y)
	}
	return a
}
func (this *ClipperStruct) ProcessHorizontal(a *TEdgeStruct, b bool) {
	c := &CProcHStruct{}
	this.GetHorzDirection(a, c)
	e := c.Dir
	f := c.Left
	g := c.Right
	h := a
	var l *TEdgeStruct = nil
	for nil != h.NextInLML && this.IsHorizontal(h.NextInLML) {
		h = h.NextInLML
	}
	if nil == h.NextInLML {
		l = this.GetMaximaPair(h)
	}
	for {
		k := (a == h)
		n := this.GetNextInAEL(a, e)
		for nil != n && !(n.Curr.X == a.Top.X && nil != a.NextInLML && n.Dx < a.NextInLML.Dx) {
			cEdge := this.GetNextInAEL(n, e)
			if (e == dLeftToRight && n.Curr.X <= g) || (e == dRightToLeft && n.Curr.X >= f) {
				if 0 <= a.OutIdx && 0 != a.WindDelta {
					this.PrepareHorzJoins(a, b)
				}
				if n == l && k {
					if e == dLeftToRight {
						this.IntersectEdges(a, n, n.Top, false)
					} else {
						this.IntersectEdges(n, a, n.Top, false)
					}
					if 0 <= l.OutIdx {
						this.logIfDebug("ProcessHorizontal error")
					}
					return
				}
				if e == dLeftToRight {
					m := &IntPoint{X: n.Curr.X, Y: a.Curr.Y}
					this.IntersectEdges(a, n, m, true)
				} else {
					m := &IntPoint{X: n.Curr.X, Y: a.Curr.Y}
					this.IntersectEdges(n, a, m, true)
				}

				this.SwapPositionsInAEL(a, n)
			} else if e == dLeftToRight && n.Curr.X >= g || e == dRightToLeft && n.Curr.X <= f {
				break
			}

			n = cEdge
		}
		if 0 <= a.OutIdx && 0 != a.WindDelta {
			this.PrepareHorzJoins(a, b)
		}
		if nil != a.NextInLML && this.IsHorizontal(a.NextInLML) {
			a = this.UpdateEdgeIntoAEL(a)
			if 0 <= a.OutIdx {
				this.AddOutPt(a, a.Bot)
			}
			c = &CProcHStruct{
				Dir:   e,
				Left:  f,
				Right: g,
			}
			this.GetHorzDirection(a, c)
			e = c.Dir
			f = c.Left
			g = c.Right
		} else {
			break
		}

	}
	if nil != a.NextInLML {
		if 0 <= a.OutIdx {
			eOutPt := this.AddOutPt(a, a.Top)
			a = this.UpdateEdgeIntoAEL(a)
			if 0 != a.WindDelta {
				fEdge := a.PrevInAEL
				cEdge := a.NextInAEL
				if nil != fEdge && fEdge.Curr.X == a.Bot.X && fEdge.Curr.Y == a.Bot.Y && 0 != fEdge.WindDelta && 0 <= fEdge.OutIdx && fEdge.Curr.Y > fEdge.Top.Y && this.SlopesEqual_3(a, fEdge, this.m_UseFullRange) {
					cOutPt := this.AddOutPt(fEdge, a.Bot)
					this.AddJoin(eOutPt, cOutPt, a.Top)

				} else {
					if nil != cEdge && cEdge.Curr.X == a.Bot.X && cEdge.Curr.Y == a.Bot.Y && 0 != cEdge.WindDelta && 0 <= cEdge.OutIdx && cEdge.Curr.Y > cEdge.Top.Y && this.SlopesEqual_3(a, cEdge, this.m_UseFullRange) {
						cOutPt := this.AddOutPt(cEdge, a.Bot)
						this.AddJoin(eOutPt, cOutPt, a.Top)

					}

				}

			}
		} else {
			this.UpdateEdgeIntoAEL(a)
		}
	} else {
		if nil != l {
			if 0 <= l.OutIdx {
				if e == dLeftToRight {
					this.IntersectEdges(a, l, a.Top, false)
				} else {
					this.IntersectEdges(l, a, a.Top, false)
					if 0 <= l.OutIdx {
						this.logIfDebug("ProcessHorizontal error")
					}
				}
			} else {
				this.DeleteFromAEL(a)
				this.DeleteFromAEL(l)
			}
		} else {
			if 0 <= a.OutIdx {
				this.AddOutPt(a, a.Top)
			}
			this.DeleteFromAEL(a)
		}
	}
}

func (this *ClipperStruct) DeleteFromSEL(a *TEdgeStruct) {
	b := a.PrevInSEL
	c := a.NextInSEL
	if nil != b || nil != c || a == this.m_SortedEdges {
		if nil != b {
			b.NextInSEL = c
		} else {
			this.m_SortedEdges = c
		}
		if nil != c {
			c.PrevInSEL = b
		}
		a.NextInSEL = nil
		a.PrevInSEL = nil
	}
}
func (this *ClipperStruct) Area(a *OutRecStruct) float64 {
	b := a.Pts
	if nil == b {
		return 0
	}
	var c float64
	c += float64(b.Prev.Pt.X+b.Pt.X) * float64(b.Prev.Pt.Y-b.Pt.Y)
	b = b.Next
	for b != a.Pts {
		c += float64(b.Prev.Pt.X+b.Pt.X) * float64(b.Prev.Pt.Y-b.Pt.Y)
		b = b.Next
	}
	return 0.5 * c
}
func (this *ClipperStruct) ProcessHorizontals(a bool) {
	horzEdge := this.m_SortedEdges
	for horzEdge != nil {
		this.DeleteFromSEL(horzEdge)
		this.ProcessHorizontal(horzEdge, a)
		horzEdge = this.m_SortedEdges
	}
}
func (this *ClipperStruct) ProcessIntersections(a, b int64) bool {
	if nil == this.m_ActiveEdges {
		return true
	}

	defer func() {
		if err := recover(); err != nil {
			this.m_SortedEdges = nil
			//	this.m_IntersectList
			// this.m_IntersectList.length = 0,
			// this.logIfDebug("ProcessIntersections error")
		}
	}()
	this.BuildIntersectList(a, b)
	if 0 == len(this.m_IntersectList) {
		return true
	}

	if 1 == len(this.m_IntersectList) || this.FixupIntersectionOrder() {
		this.ProcessIntersectList()
	} else {
		return false
	}
	this.m_SortedEdges = nil
	return true
}
func (this *ClipperStruct) CopyAELToSEL() { //saya?
	e := this.m_ActiveEdges
	this.m_SortedEdges = e
	if this.m_ActiveEdges == nil {
		return
	}
	this.m_SortedEdges.PrevInSEL = nil
	e = e.NextInAEL
	for e != nil {
		e.PrevInSEL = e.PrevInAEL
		e.PrevInSEL.NextInSEL = e
		e.NextInSEL = nil
		e = e.NextInAEL
	}
}
func (this *ClipperStruct) EdgesAdjacent(a *IntersectNodeStruct) bool {
	return a.Edge1.NextInSEL == a.Edge2 || a.Edge1.PrevInSEL == a.Edge2
}

func (this *ClipperStruct) SwapPositionsInSEL(a, b *TEdgeStruct) {
	if nil != a.NextInSEL || nil != a.PrevInSEL {
		if nil != b.NextInSEL || nil != b.PrevInSEL {
			if a.NextInSEL == b {
				var c = b.NextInSEL
				if nil != c {
					c.PrevInSEL = a
				}
				var e = a.PrevInSEL
				if nil != e {
					e.NextInSEL = b
				}
				b.PrevInSEL = e
				b.NextInSEL = a
				a.PrevInSEL = b
				a.NextInSEL = c
			} else if b.NextInSEL == a {
				c := a.NextInSEL
				if nil != c {
					c.PrevInSEL = b
				}

				e := b.PrevInSEL
				if nil != e {
					e.NextInSEL = a
				}
				a.PrevInSEL = e
				a.NextInSEL = b
				b.PrevInSEL = a
				b.NextInSEL = c
			} else {
				c := a.NextInSEL
				e := a.PrevInSEL
				a.NextInSEL = b.NextInSEL
				if nil != a.NextInSEL {
					a.NextInSEL.PrevInSEL = a
				}
				a.PrevInSEL = b.PrevInSEL
				if nil != a.PrevInSEL {
					a.PrevInSEL.NextInSEL = a
				}
				b.NextInSEL = c
				if nil != b.NextInSEL {
					b.NextInSEL.PrevInSEL = b
				}
				b.PrevInSEL = e
				if nil != b.PrevInSEL {
					b.PrevInSEL.NextInSEL = b
				}
			}
			if nil == a.PrevInSEL {
				this.m_SortedEdges = a
			} else {
				if nil == b.PrevInSEL {
					this.m_SortedEdges = b
				}
			}
		}
	}

}

func (this *ClipperStruct) FixupIntersectionOrder() bool {
	sort.Sort(IntersectNodeList(this.m_IntersectList))
	this.CopyAELToSEL()
	a := len(this.m_IntersectList)
	b := 0
	for ; b < a; b++ {
		if !this.EdgesAdjacent(this.m_IntersectList[b]) {
			c := b + 1
			for c < a && !this.EdgesAdjacent(this.m_IntersectList[c]) {
				c++
			}
			if c == a {
				return false
			}
			e := this.m_IntersectList[b]
			this.m_IntersectList[b] = this.m_IntersectList[c]
			this.m_IntersectList[c] = e
		}
		this.SwapPositionsInSEL(this.m_IntersectList[b].Edge1, this.m_IntersectList[b].Edge2)
	}
	return true
}
func (this *ClipperStruct) BuildIntersectList(a, b int64) {
	if nil != this.m_ActiveEdges { //todo
		c := this.m_ActiveEdges
		this.m_SortedEdges = c
		// check := this.SortedEdges()
		// this.logIfDebug(check)
		for nil != c {
			c.PrevInSEL = c.PrevInAEL
			c.NextInSEL = c.NextInAEL

			c.Curr.X = this.TopX(c, b)
			//this.logIfDebug("c.Curr.X:", c.Curr.X)
			c = c.NextInAEL
		}
		// check = this.SortedEdges()
		// this.logIfDebug(check)
		for e := true; e && nil != this.m_SortedEdges; {
			e = false
			for c = this.m_SortedEdges; nil != c.NextInSEL; {
				f := c.NextInSEL

				g := &IntPoint{}
				if c.Curr.X > f.Curr.X {
					if !this.IntersectPoint(c, f, g) && c.Curr.X > f.Curr.X+1 {
						this.logIfDebug("Intersection error")
					}
					if g.Y > a {
						g.Y = a
						if math.Abs(c.Dx) > math.Abs(f.Dx) {
							g.X = this.TopX(f, a)

						} else {
							g.X = this.TopX(c, a)
						}
					}

					eIntersectNode := &IntersectNodeStruct{
						Pt: &IntPoint{},
					}
					eIntersectNode.Edge1 = c
					eIntersectNode.Edge2 = f
					eIntersectNode.Pt.X = g.X
					eIntersectNode.Pt.Y = g.Y
					this.m_IntersectList = append(this.m_IntersectList, eIntersectNode)
					this.SwapPositionsInSEL(c, f)
					e = true
				} else {
					c = f
				}
			}
			if nil != c.PrevInSEL {
				c.PrevInSEL.NextInSEL = nil
			} else {
				break
			}
		}
		this.m_SortedEdges = nil
	}
}
func (this *ClipperStruct) IntersectPoint(a, b *TEdgeStruct, c *IntPoint) bool {
	c.X = 0
	c.Y = 0
	var e, f float64
	if this.SlopesEqual_3(a, b, this.m_UseFullRange) || a.Dx == b.Dx {
		if b.Bot.Y > a.Bot.Y {
			c.X = b.Bot.X
			c.Y = b.Bot.Y
		} else {
			c.X = a.Bot.X
			c.Y = a.Bot.Y
		}
		return false
	}

	if 0 == a.Delta.X {
		c.X = a.Bot.X
		if this.IsHorizontal(b) {
			c.Y = b.Bot.Y
		} else {
			f = float64(b.Bot.Y) - float64(b.Bot.X)/b.Dx
			c.Y = int64(math.Round(float64(c.X)/b.Dx + f))
		}
	} else if 0 == b.Delta.X {
		c.X = b.Bot.X
		if this.IsHorizontal(a) {
			c.Y = a.Bot.Y
		} else {
			e = float64(a.Bot.Y) - float64(a.Bot.X)/a.Dx
			c.Y = int64(math.Round(float64(c.X)/a.Dx + e))
		}
	} else {
		e = float64(a.Bot.X) - float64(a.Bot.Y)*a.Dx
		f = float64(b.Bot.X) - float64(b.Bot.Y)*b.Dx
		g := (f - e) / (a.Dx - b.Dx)
		c.Y = int64(math.Round(g))
		if math.Abs(a.Dx) < math.Abs(b.Dx) {
			c.X = int64(math.Round(a.Dx*g + e))
		} else {
			c.X = int64(math.Round(b.Dx*g + f))
		}
	}
	if c.Y < a.Top.Y || c.Y < b.Top.Y {
		if a.Top.Y > b.Top.Y {
			c.Y = a.Top.Y
			c.X = this.TopX(b, a.Top.Y)
			return c.X < a.Top.X
		}

		c.Y = b.Top.Y
		if math.Abs(a.Dx) < math.Abs(b.Dx) {
			c.X = this.TopX(a, c.Y)
		} else {
			c.X = this.TopX(b, c.Y)
		}
	}
	return true
}
func (this *ClipperStruct) ProcessIntersectList() {
	a := 0
	b := len(this.m_IntersectList)
	for ; a < b; a++ {
		c := this.m_IntersectList[a]
		this.IntersectEdges(c.Edge1, c.Edge2, c.Pt, true)
		this.SwapPositionsInAEL(c.Edge1, c.Edge2)
	}
	this.m_IntersectList = []*IntersectNodeStruct{}
}

func (this *ClipperStruct) IsMaxima(a *TEdgeStruct, b int64) bool {
	return nil != a && a.Top.Y == b && nil == a.NextInLML
}
func (this *ClipperStruct) DoMaxima(a *TEdgeStruct) {
	var b = this.GetMaximaPair(a)
	if nil == b {
		if 0 <= a.OutIdx {
			this.AddOutPt(a, a.Top)
		}
		this.DeleteFromAEL(a)
	} else {
		c := a.NextInAEL
		for nil != c && c != b {
			this.IntersectEdges(a, c, a.Top, true)
			this.SwapPositionsInAEL(a, c)
			c = a.NextInAEL
		}

		if -1 == a.OutIdx && -1 == b.OutIdx {
			this.DeleteFromAEL(a)
			this.DeleteFromAEL(b)
		} else {
			if 0 <= a.OutIdx && 0 <= b.OutIdx {
				this.IntersectEdges(a, b, a.Top, false)
			} else {
				if 0 == a.WindDelta {
					if 0 <= a.OutIdx {
						this.AddOutPt(a, a.Top)
						a.OutIdx = -1
					}
					this.DeleteFromAEL(a)
					if 0 <= b.OutIdx {
						this.AddOutPt(b, a.Top)
						b.OutIdx = -1
					}
					this.DeleteFromAEL(b)
				} else {
					log.Println("DoMaxima error")
					panic("DoMaxima error")
				}
			}
		}

	}
}

func (this *ClipperStruct) IsIntermediate(a *TEdgeStruct, b int64) bool {
	return a.Top.Y == b && nil != a.NextInLML
}
func (this *ClipperStruct) ProcessEdgesAtTopOfScanbeam(a int64) {
	b := this.m_ActiveEdges
	for nil != b {
		c := this.IsMaxima(b, a)
		if c {
			cEdge := this.GetMaximaPair(b)
			c = (nil == cEdge || !this.IsHorizontal(cEdge))
		}
		if c {
			e := b.PrevInAEL
			this.DoMaxima(b)

			if nil == e {
				b = this.m_ActiveEdges
			} else {
				b = e.NextInAEL
			}
		} else {
			if this.IsIntermediate(b, a) && this.IsHorizontal(b.NextInLML) {
				b = this.UpdateEdgeIntoAEL(b)
				if 0 <= b.OutIdx {
					this.AddOutPt(b, b.Bot)
				}
				this.AddEdgeToSEL(b)
			} else {
				b.Curr.X = this.TopX(b, a)
				b.Curr.Y = a
			}
			if this.StrictlySimple {
				e := b.PrevInAEL
				cOutPt := &OutPtStruct{
					Pt: &IntPoint{},
				}
				if 0 <= b.OutIdx && 0 != b.WindDelta && nil != e && 0 <= e.OutIdx && e.Curr.X == b.Curr.X && 0 != e.WindDelta {
					cOutPt = this.AddOutPt(e, b.Curr)
					eOutPt := this.AddOutPt(b, b.Curr)
					this.AddJoin(cOutPt, eOutPt, b.Curr)
				}
			}
			b = b.NextInAEL
		}
	}
	this.ProcessHorizontals(true)
	bEdge := this.m_ActiveEdges
	for nil != bEdge {
		if this.IsIntermediate(bEdge, a) {
			var c *OutPtStruct = nil
			if 0 <= bEdge.OutIdx {
				c = this.AddOutPt(bEdge, bEdge.Top)
			}
			bEdge = this.UpdateEdgeIntoAEL(bEdge)
			e := bEdge.PrevInAEL
			f := bEdge.NextInAEL
			if nil != e && e.Curr.X == bEdge.Bot.X && e.Curr.Y == bEdge.Bot.Y && nil != c && 0 <= e.OutIdx && e.Curr.Y > e.Top.Y && this.SlopesEqual_3(bEdge, e, this.m_UseFullRange) && 0 != bEdge.WindDelta && 0 != e.WindDelta {
				eOutPt := this.AddOutPt(e, bEdge.Bot)
				this.AddJoin(c, eOutPt, bEdge.Top)
			} else {
				if nil != f && f.Curr.X == bEdge.Bot.X && f.Curr.Y == bEdge.Bot.Y && nil != c && 0 <= f.OutIdx && f.Curr.Y > f.Top.Y && this.SlopesEqual_3(bEdge, f, this.m_UseFullRange) && 0 != bEdge.WindDelta && 0 != f.WindDelta {
					eOutPt := this.AddOutPt(f, bEdge.Bot)
					this.AddJoin(c, eOutPt, bEdge.Top)
				}
			}
		}
		bEdge = bEdge.NextInAEL
	}
}
func recoverName() {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
	}
}
func (this *ClipperStruct) ExecuteInternal() bool {
	defer func() {
		if r := recover(); r != nil {

		}
		this.m_Joins = []*JoinStruct{}
		this.m_GhostJoins = []*JoinStruct{}
	}()
	// defer func() {
	// 	this.m_Joins = []*JoinStruct{}
	// 	this.m_GhostJoins = []*JoinStruct{}
	// }()//saya
	this.Reset()
	if nil == this.m_CurrentLM {
		return false
	}

	a := this.PopScanbeam()
	isBreak := false
	this.InsertLocalMinimaIntoAEL(a)
	this.m_GhostJoins = []*JoinStruct{}
	//d.Clear(this.m_GhostJoins);
	this.ProcessHorizontals(false)
	if nil == this.m_Scanbeam {
		isBreak = true
	}

	b := this.PopScanbeam()
	if !this.ProcessIntersections(a, b) {
		return false
	}
	this.ProcessEdgesAtTopOfScanbeam(b)
	a = b
	for !isBreak && (nil != this.m_Scanbeam || nil != this.m_CurrentLM) {
		this.InsertLocalMinimaIntoAEL(a)
		this.m_GhostJoins = []*JoinStruct{}
		// d.Clear(this.m_GhostJoins);
		this.ProcessHorizontals(false)
		if nil == this.m_Scanbeam {
			break
		}

		b := this.PopScanbeam()
		if !this.ProcessIntersections(a, b) {
			return false
		}
		this.ProcessEdgesAtTopOfScanbeam(b)
		a = b
	}
	aIndex := 0
	c := len(this.m_PolyOuts)
	for ; aIndex < c; aIndex++ {
		e := this.m_PolyOuts[aIndex]
		if nil == e.Pts {
			if !e.IsOpen {
				if e.IsHole && !this.ReverseSolution && 0 < this.Area(e) {
					this.ReversePolyPtLinks(e.Pts)
				}
				if !e.IsHole && this.ReverseSolution && 0 < this.Area(e) {
					this.ReversePolyPtLinks(e.Pts)
				}
				if e.IsHole && this.ReverseSolution && !(0 < this.Area(e)) {
					this.ReversePolyPtLinks(e.Pts)
				}
				if !e.IsHole && !this.ReverseSolution && !(0 < this.Area(e)) {
					this.ReversePolyPtLinks(e.Pts)
				}

			}
		}

	}
	cIndex := len(this.m_PolyOuts)

	//cJoinIndex := len(this.m_Joins)
	//	this.logIfDebug(cJoinIndex)
	this.JoinCommonEdges()
	aIndex = 0
	cIndex = len(this.m_PolyOuts)

	for ; aIndex < cIndex; aIndex++ {
		e := this.m_PolyOuts[aIndex]
		if !(nil == e.Pts || e.IsOpen) {
			this.FixupOutPolygon(e)
		}
	}

	if this.StrictlySimple {
		this.DoSimplePolygons()
	}
	return true

}

func (this *ClipperStruct) GetOutRec(a int) *OutRecStruct {
	aOutRec := this.m_PolyOuts[a]
	for aOutRec != this.m_PolyOuts[aOutRec.Idx] {
		aOutRec = this.m_PolyOuts[aOutRec.Idx]
	}
	return aOutRec
}

func (this *ClipperStruct) DupOutPt(a *OutPtStruct, b bool) *OutPtStruct {
	c := &OutPtStruct{
		Pt: &IntPoint{},
	}
	c.Pt.X = a.Pt.X
	c.Pt.Y = a.Pt.Y
	c.Idx = a.Idx
	if b {
		c.Next = a.Next
		c.Prev = a
		a.Next.Prev = c
		a.Next = c
	} else {
		c.Prev = a.Prev
		c.Next = a
		a.Prev.Next = c
		a.Prev = c
	}
	return c
}

type CStruct struct {
	Left  int64
	Right int64
}

func (this *ClipperStruct) GetOverlap(a, b, c, e int64, d *CStruct) bool {
	if a < b {
		if c < e {
			d.Left = Max(a, c)
			d.Right = Min(b, e)
		} else {
			d.Left = Max(a, e)
			d.Right = Min(b, c)
		}
	} else {
		if c < e {
			d.Left = Max(b, c)
			d.Right = Min(a, e)
		} else {
			d.Left = Max(b, e)
			d.Right = Min(a, c)
		}
	}
	return d.Left < d.Right
}
func (this *ClipperStruct) JoinPoints(a *JoinStruct, b *OutRecStruct, c *OutRecStruct) bool {
	e := a.OutPt1
	f := &OutPtStruct{
		Pt: &IntPoint{},
	}
	g := a.OutPt2
	h := &OutPtStruct{
		Pt: &IntPoint{},
	}
	hBool := a.OutPt1.Pt.Y == a.OffPt.Y
	if hBool && op_Equality_IntPoint(a.OffPt, a.OutPt1.Pt) && op_Equality_IntPoint(a.OffPt, a.OutPt2.Pt) {
		for f = a.OutPt1.Next; f != e && op_Equality_IntPoint(f.Pt, a.OffPt); {
			f = f.Next
		}

		fBool := f.Pt.Y > a.OffPt.Y
		for h = a.OutPt2.Next; h != g && op_Equality_IntPoint(h.Pt, a.OffPt); {
			h = h.Next
		}

		if fBool == (h.Pt.Y > a.OffPt.Y) {
			return false
		}

		if fBool {
			f = this.DupOutPt(e, false)
			h = this.DupOutPt(g, true)
			e.Prev = g
			g.Next = e
			f.Next = h
			h.Prev = f
		} else {
			f = this.DupOutPt(e, true)
			h = this.DupOutPt(g, false)
			e.Next = g
			g.Prev = e
			f.Prev = h
			h.Next = f
		}
		a.OutPt1 = e
		a.OutPt2 = f
		return true
	}
	if hBool {
		f = e
		for e.Prev.Pt.Y == e.Pt.Y && e.Prev != f && e.Prev != g {
			e = e.Prev
		}

		for f.Next.Pt.Y == f.Pt.Y && f.Next != e && f.Next != g {
			f = f.Next
		}

		if f.Next == e || f.Next == g {
			return false
		}
		for h = g; g.Prev.Pt.Y == g.Pt.Y && g.Prev != h && g.Prev != f; {
			g = g.Prev
		}
		for h.Next.Pt.Y == h.Pt.Y && h.Next != g && h.Next != e {
			h = h.Next
		}
		if h.Next == g || h.Next == e {
			return false
		}
		c := &CStruct{}

		if !this.GetOverlap(e.Pt.X, f.Pt.X, g.Pt.X, h.Pt.X, c) {
			return false
		}

		bInt := c.Left
		lInt := c.Right
		cPoint := &IntPoint{}
		bBool := false
		if e.Pt.X >= bInt && e.Pt.X <= lInt {
			cPoint.X = e.Pt.X
			cPoint.Y = e.Pt.Y
			bBool = e.Pt.X > f.Pt.X
		} else {
			if g.Pt.X >= bInt && g.Pt.X <= lInt {
				cPoint.X = g.Pt.X
				cPoint.Y = g.Pt.Y
				bBool = g.Pt.X > h.Pt.X
			} else {
				if f.Pt.X >= bInt && f.Pt.X <= lInt {
					cPoint.X = f.Pt.X
					cPoint.Y = f.Pt.Y
					bBool = f.Pt.X > e.Pt.X
				} else {
					cPoint.X = h.Pt.X
					cPoint.Y = h.Pt.Y
					bBool = h.Pt.X > g.Pt.X
				}
			}
		}
		a.OutPt1 = e
		a.OutPt2 = g
		return this.JoinHorz(e, f, g, h, cPoint, bBool)
	}
	for f = e.Next; op_Equality_IntPoint(f.Pt, e.Pt) && f != e; {
		f = f.Next
	}
	lBool := f.Pt.Y > e.Pt.Y || !this.SlopesEqual_4(e.Pt, f.Pt, a.OffPt, this.m_UseFullRange)
	if lBool {
		for f = e.Prev; op_Equality_IntPoint(f.Pt, e.Pt) && f != e; {
			f = f.Prev
		}

		if f.Pt.Y > e.Pt.Y || !this.SlopesEqual_4(e.Pt, f.Pt, a.OffPt, this.m_UseFullRange) {
			return false
		}

	}
	for h = g.Next; op_Equality_IntPoint(h.Pt, g.Pt) && h != g; {
		h = h.Next
	}

	k := (h.Pt.Y > g.Pt.Y || !this.SlopesEqual_4(g.Pt, h.Pt, a.OffPt, this.m_UseFullRange))
	if k {
		for h = g.Prev; op_Equality_IntPoint(h.Pt, g.Pt) && h != g; {
			h = h.Prev
		}

		if h.Pt.Y > g.Pt.Y || !this.SlopesEqual_4(g.Pt, h.Pt, a.OffPt, this.m_UseFullRange) {
			return false
		}

	}
	if f == e || h == g || f == h || b == c && lBool == k {
		return false
	}

	if lBool {
		f = this.DupOutPt(e, false)
		h = this.DupOutPt(g, true)
		e.Prev = g
		g.Next = e
		f.Next = h
		h.Prev = f
	} else {
		f = this.DupOutPt(e, true)
		h = this.DupOutPt(g, false)
		e.Next = g
		g.Prev = e
		f.Prev = h
		h.Next = f
	}
	a.OutPt1 = e
	a.OutPt2 = f
	return true
}
func (this *ClipperStruct) JoinCommonEdges() {
	a := 0
	b := len(this.m_Joins)

	for ; a < b; a++ {
		// if a == 99 {
		// 	this.logIfDebug("a:")
		// }
		// if this.m_Joins[99].OutPt2.Idx == 128 {
		// 	this.logIfDebug("a:")
		// }
		this.logIfDebug("JoinCommonEdges:", a)
		c := this.m_Joins[a]
		e := this.GetOutRec(c.OutPt1.Idx)
		f := this.GetOutRec(c.OutPt2.Idx)
		if a == 12 {
			this.logIfDebug("a:")
		}
		if nil != e.Pts && nil != f.Pts {
			g := &OutRecStruct{}
			if e == f {
				g = e
			} else {
				if this.Param1RightOfParam2(e, f) {
					g = f
				} else {

					if this.Param1RightOfParam2(f, e) {
						g = e
					} else {
						g = this.GetLowermostRec(e, f)
					}
				}

			}

			if this.JoinPoints(c, e, f) {
				if e == f {
					e.Pts = c.OutPt1
					e.BottomPt = nil
					f = this.CreateOutRec()
					f.Pts = c.OutPt2
					this.UpdateOutPtIdxs(f)
					if this.m_UsingPolyTree {
						g := 0
						h := len(this.m_PolyOuts)
						for ; g < h-1; g++ {
							l := this.m_PolyOuts[g]
							if nil != l.Pts && this.ParseFirstLeft(l.FirstLeft) == e && l.IsHole != e.IsHole && this.Poly2ContainsPoly1(l.Pts, c.OutPt2) {
								l.FirstLeft = f
							}
						}
					}
					if this.Poly2ContainsPoly1(f.Pts, e.Pts) {
						f.IsHole = !e.IsHole
						f.FirstLeft = e
						if this.m_UsingPolyTree {
							this.FixupFirstLefts2(f, e)
						}
						isHole := 0
						if e.IsHole {
							isHole = 1
						} else {
							isHole = 0
						}
						reverseSolution := 0
						if this.ReverseSolution {
							reverseSolution = 1
						} else {
							reverseSolution = 0
						}

						if (isHole^reverseSolution == 1) == (0 < this.Area(f)) {
							this.ReversePolyPtLinks(f.Pts)
						} else {
							if this.Poly2ContainsPoly1(e.Pts, f.Pts) {
								f.IsHole = e.IsHole
								e.IsHole = !f.IsHole
								f.FirstLeft = e.FirstLeft
								e.FirstLeft = f
								if this.m_UsingPolyTree {
									this.FixupFirstLefts2(e, f)
								}
								isHole := 0
								if e.IsHole {
									isHole = 1
								} else {
									isHole = 0
								}
								reverseSolution := 0
								if this.ReverseSolution {
									reverseSolution = 1
								} else {
									reverseSolution = 0
								}

								if (isHole^reverseSolution == 1) == (0 < this.Area(e)) {
									this.ReversePolyPtLinks(e.Pts)
								}
							} else {
								f.IsHole = e.IsHole
								f.FirstLeft = e.FirstLeft
								if this.m_UsingPolyTree {
									this.FixupFirstLefts1(e, f)
								}
							}

						}
					}
				} else {
					f.Pts = nil
					f.BottomPt = nil
					f.Idx = e.Idx
					e.IsHole = g.IsHole
					if g == f {
						e.FirstLeft = f.FirstLeft
					}
					f.FirstLeft = e
					if this.m_UsingPolyTree {
						this.FixupFirstLefts2(f, e)
					}
				}

			}
		}
	}
}

func (this *ClipperStruct) FixupFirstLefts1(a, b *OutRecStruct) {
	for c := 0; c < len(this.m_PolyOuts); c++ {
		d := this.m_PolyOuts[c]
		if nil != d.Pts && d.FirstLeft == a && this.Poly2ContainsPoly1(d.Pts, b.Pts) {
			d.FirstLeft = b
		}
	}
}
func (this *ClipperStruct) FixupFirstLefts2(a, b *OutRecStruct) {
	c := 0
	e := this.m_PolyOuts
	d := len(e)
	g := e[c]
	for c < d {
		if g.FirstLeft == a {
			g.FirstLeft = b
		}
		c++
		g = e[c]
	}
}
func (this *ClipperStruct) ParseFirstLeft(a *OutRecStruct) *OutRecStruct {
	for nil != a && nil == a.Pts {
		a = a.FirstLeft
	}
	return a
}
func (this *ClipperStruct) Poly2ContainsPoly1(a, b *OutPtStruct) bool {
	c := a

	e := this.PointInPolygon(c.Pt, b)
	if 0 <= e {
		return 0 != e
	}
	c = c.Next
	for c != a {
		e = this.PointInPolygon(c.Pt, b)
		if 0 <= e {
			return 0 != e
		}
		c = c.Next
	}
	return true

}
func (this *ClipperStruct) PointInPolygon(a *IntPoint, b *OutPtStruct) int64 {
	var c int64 = 0
	e := b
	for {
		d := b.Pt.X
		g := b.Pt.Y
		h := b.Next.Pt.X
		l := b.Next.Pt.Y
		if l == a.Y && (h == a.X || g == a.Y && (h > a.X) == (d < a.X)) {
			return -1
		}

		if (g < a.Y) != (l < a.Y) {
			if d >= a.X {
				if h > a.X {
					c = 1 - c
				} else {
					d = (d-a.X)*(l-a.Y) - (h-a.X)*(g-a.Y)
					if 0 == d {
						return -1
					}

					if (0 < d) == (l > g) {
						c = 1 - c
					}
				}
			} else if h > a.X {
				d = (d-a.X)*(l-a.Y) - (h-a.X)*(g-a.Y)
				if 0 == d {
					return -1
				}

				if (0 < d) == (l > g) {
					c = 1 - c
				}
			}
		}
		b = b.Next
		if e == b {
			break
		}

	}
	return c
}
func (this *ClipperStruct) UpdateOutPtIdxs(a *OutRecStruct) {
	this.logIfDebug("UpdateOutPtIdxs A:", a.Idx)
	b := a.Pts
	b.Idx = a.Idx
	b = b.Prev
	for b != a.Pts {
		b.Idx = a.Idx
		b = b.Prev
	}
}
func (this *ClipperStruct) JoinHorz(a *OutPtStruct, b *OutPtStruct, c *OutPtStruct, e *OutPtStruct, f *IntPoint, g bool) bool {

	var hDirection Direction
	if a.Pt.X > b.Pt.X {
		hDirection = dRightToLeft
	} else {
		hDirection = dLeftToRight
	}
	var eDirection Direction
	if c.Pt.X > e.Pt.X {
		eDirection = dRightToLeft
	} else {
		eDirection = dLeftToRight
	}
	if hDirection == eDirection {
		return false
	}

	if hDirection == dLeftToRight {
		for a.Next.Pt.X <= f.X && a.Next.Pt.X >= a.Pt.X && a.Next.Pt.Y == f.Y {
			a = a.Next
		}

		if g && a.Pt.X != f.X {
			a = a.Next
		}
		b = this.DupOutPt(a, !g)
		if op_InEquality_IntPoint(b.Pt, f) {
			a = b
			a.Pt.X = f.X
			a.Pt.Y = f.Y
			b = this.DupOutPt(a, !g)
		}
	} else {
		for a.Next.Pt.X >= f.X && a.Next.Pt.X <= a.Pt.X && a.Next.Pt.Y == f.Y {
			a = a.Next
		}

		if !g && a.Pt.X != f.X {
			a = a.Next
		}
		b = this.DupOutPt(a, g)
		if op_InEquality_IntPoint(b.Pt, f) {
			a = b
			a.Pt.X = f.X
			a.Pt.Y = f.Y
			b = this.DupOutPt(a, g)
		}
	}
	if eDirection == dLeftToRight {
		for c.Next.Pt.X <= f.X && c.Next.Pt.X >= c.Pt.X && c.Next.Pt.Y == f.Y {
			c = c.Next
		}

		if g && c.Pt.X != f.X {
			c = c.Next
		}
		e = this.DupOutPt(c, !g)
		if op_InEquality_IntPoint(e.Pt, f) {
			c = e
			c.Pt.X = f.X
			c.Pt.Y = f.Y
			e = this.DupOutPt(c, !g)
		}
	} else {
		for c.Next.Pt.X >= f.X && c.Next.Pt.X <= c.Pt.X && c.Next.Pt.Y == f.Y {
			c = c.Next
		}

		if !g && c.Pt.X != f.X {
			c = c.Next
		}
		e = this.DupOutPt(c, g)
		if op_InEquality_IntPoint(e.Pt, f) {
			c = e
			c.Pt.X = f.X
			c.Pt.Y = f.Y
			e = this.DupOutPt(c, g)
		}
	}
	if (hDirection == dLeftToRight) == g {
		a.Prev = c

		c.Next = a
		b.Next = e
		e.Prev = b
	} else {
		a.Next = c
		c.Prev = a
		b.Prev = e
		e.Next = b
	}
	return true
}

func (this *ClipperStruct) FixupOutPolygon(a *OutRecStruct) {

	var b *OutPtStruct = nil
	a.BottomPt = nil
	c := a.Pts
	for {
		if c.Prev == c || c.Prev == c.Next {
			this.DisposeOutPts(c)
			a.Pts = nil
			return
		}
		if op_Equality_IntPoint(c.Pt, c.Next.Pt) || op_Equality_IntPoint(c.Pt, c.Prev.Pt) || this.SlopesEqual_4(c.Prev.Pt, c.Pt, c.Next.Pt, this.m_UseFullRange) && (!this.PreserveCollinear || !this.Pt2IsBetweenPt1AndPt3(c.Prev.Pt, c.Pt, c.Next.Pt)) {
			b = nil
			c.Prev.Next = c.Next
			c.Next.Prev = c.Prev
			c = c.Prev
		} else if c == b {
			break
		} else {
			if nil == b {
				b = c
			}
			c = c.Next
		}
	}
	a.Pts = c

}

func (this *ClipperStruct) DoSimplePolygons() {
	a := 0
	for a < len(this.m_PolyOuts) {
		b := this.m_PolyOuts[a]
		c := b.Pts
		a++
		if nil != c {

			for e := c.Next; e != b.Pts; {
				if op_Equality_IntPoint(c.Pt, e.Pt) && e.Next != c && e.Prev != c {
					f := c.Prev
					g := e.Prev
					c.Prev = g
					g.Next = c
					e.Prev = f
					f.Next = e
					b.Pts = c
					fOutRec := this.CreateOutRec()
					fOutRec.Pts = e
					this.UpdateOutPtIdxs(fOutRec)
					if this.Poly2ContainsPoly1(fOutRec.Pts, b.Pts) {
						fOutRec.IsHole = !b.IsHole
						fOutRec.FirstLeft = b
					} else {
						if this.Poly2ContainsPoly1(b.Pts, fOutRec.Pts) {
							fOutRec.IsHole = b.IsHole
							b.IsHole = !fOutRec.IsHole
							fOutRec.FirstLeft = b.FirstLeft
							b.FirstLeft = fOutRec
						} else {
							fOutRec.IsHole = b.IsHole
							fOutRec.FirstLeft = b.FirstLeft
						}
					}
					e = c
				}
				e = e.Next
			}
			c = c.Next

			for c != b.Pts {
				for e := c.Next; e != b.Pts; {
					if op_Equality_IntPoint(c.Pt, e.Pt) && e.Next != c && e.Prev != c {
						f := c.Prev
						g := e.Prev
						c.Prev = g
						g.Next = c
						e.Prev = f
						f.Next = e
						b.Pts = c
						fOutRec := this.CreateOutRec()
						fOutRec.Pts = e
						this.UpdateOutPtIdxs(fOutRec)
						if this.Poly2ContainsPoly1(fOutRec.Pts, b.Pts) {
							fOutRec.IsHole = !b.IsHole
							fOutRec.FirstLeft = b
						} else {
							if this.Poly2ContainsPoly1(b.Pts, fOutRec.Pts) {
								fOutRec.IsHole = b.IsHole
								b.IsHole = !fOutRec.IsHole
								fOutRec.FirstLeft = b.FirstLeft
								b.FirstLeft = fOutRec
							} else {
								fOutRec.IsHole = b.IsHole
								fOutRec.FirstLeft = b.FirstLeft
							}
						}
						e = c
					}
					e = e.Next
				}
				c = c.Next
			}
		}
	}
}

type IntersectNodeList []*IntersectNodeStruct

func (s IntersectNodeList) Len() int           { return len(s) }
func (s IntersectNodeList) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s IntersectNodeList) Less(i, j int) bool { return s[i].Pt.Y < s[j].Pt.Y }

type ClipperStruct struct {
	m_ClipType ClipType

	m_ExecuteLocked bool
	m_SubjFillType  PolyFillType
	m_ClipFillType  PolyFillType
	m_GhostJoins    []*JoinStruct
	StrictlySimple  bool
	// 	ClipperBase_call(this);
	m_SortedEdges   *TEdgeStruct
	m_ActiveEdges   *TEdgeStruct
	m_Scanbeam      *ScanbeamSturct
	m_CurrentLM     *LocalMinimaStruct
	m_IntersectList []*IntersectNodeStruct
	// 	m_IntersectNodeComparer
	m_UsingPolyTree bool
	m_PolyOuts      []*OutRecStruct
	m_Joins         []*JoinStruct

	m_MinimaList *LocalMinimaStruct

	m_edges [][]*TEdgeStruct

	PreserveCollinear bool
	m_HasOpenPaths    bool
	m_UseFullRange    bool

	ReverseSolution bool
	IfDebug         bool
}

func Clipper(a int) *ClipperStruct {
	this := &ClipperStruct{
		m_ClipType:     ctIntersection,
		m_SubjFillType: pftEvenOdd,
		m_ClipFillType: pftEvenOdd,
	}
	this.ReverseSolution = 0 != (1 & a)
	this.StrictlySimple = 0 != (2 & a)
	this.PreserveCollinear = 0 != (4 & a)
	return this
	// this.m_PolyOuts = nil;
	// this.m_IntersectNodeComparer = this.m_IntersectList = this.m_SortedEdges = this.m_ActiveEdges = this.m_Scanbeam = nil;
	// this.m_ExecuteLocked = false;
	// this.m_SubjFillType = this.m_ClipFillType = pftEvenOdd;
	// this.m_GhostJoins = this.m_Joins = nil;
	// this.StrictlySimple = this.ReverseSolution = this.m_UsingPolyTree = false;
	// ClipperBase_call(this);
	// this.m_SortedEdges = this.m_ActiveEdges = this.m_Scanbeam = nil;
	// this.m_IntersectList = [];
	// this.m_IntersectNodeComparer = d.MyIntersectNodeSort.Compare;
	// this.m_UsingPolyTree = this.m_ExecuteLocked = false;
	// this.m_PolyOuts = [];
	// this.m_Joins = [];
	// this.m_GhostJoins = [];
	// this.ReverseSolution = 0 != (1 & a);
	// this.StrictlySimple = 0 != (2 & a);
	// this.PreserveCollinear = 0 != (4 & a)

}

func (this *ClipperStruct) AddPaths(a []IntPolygon, b PolyType, c bool) bool {
	e := false
	d := 0
	for g := len(a); d < g; d++ {
		if this.AddPath(a[d], b, c) {
			e = true
		}
	}
	return e
}

const hiRange = 4503599627370495
const loRange = 47453132

func (this *ClipperStruct) RangeTest(a *IntPoint, b *lStruct) {
	if b.Value {
		if a.X > hiRange || a.Y > hiRange || -a.X > hiRange || -a.Y > hiRange {
			this.logIfDebug("Coordinate outside allowed range in RangeTest().")
		}
	} else if a.X > loRange || a.Y > loRange || -a.X > loRange || -a.Y > ClipperBase_loRange {
		b.Value = true
		this.RangeTest(a, b)
	}
}
func (this *ClipperStruct) InitEdge(a, b, c *TEdgeStruct, e *IntPoint) {
	a.Next = b
	a.Prev = c
	a.Curr.X = e.X
	a.Curr.Y = e.Y
	a.OutIdx = -1
}
func (this *ClipperStruct) InitEdge2(a *TEdgeStruct, b PolyType) {
	if a.Curr.Y >= a.Next.Curr.Y {
		a.Bot.X = a.Curr.X
		a.Bot.Y = a.Curr.Y
		a.Top.X = a.Next.Curr.X
		a.Top.Y = a.Next.Curr.Y
	} else {
		a.Top.X = a.Curr.X
		a.Top.Y = a.Curr.Y
		a.Bot.X = a.Next.Curr.X
		a.Bot.Y = a.Next.Curr.Y
	}
	this.SetDx(a)
	a.PolyTyp = b
}
func (this *ClipperStruct) SetDx(a *TEdgeStruct) {
	a.Delta.X = a.Top.X - a.Bot.X
	a.Delta.Y = a.Top.Y - a.Bot.Y
	if a.Delta.Y == 0 {
		a.Dx = ClipperBase_horizontal
	} else {
		a.Dx = float64(a.Delta.X) / float64(a.Delta.Y)
	}
}
func (this *ClipperStruct) Pt2IsBetweenPt1AndPt3(a, b, c *IntPoint) bool {
	tempResult := op_Equality_IntPoint(a, c) || op_Equality_IntPoint(a, b)
	if op_Equality_IntPoint(c, b) {
		return tempResult || false
	} else {
		if a.X != c.X {
			return tempResult || (b.X > a.X) == (b.X < c.X)
		} else {
			return tempResult || (b.Y > a.Y) == (b.Y < c.Y)
		}

	}
}
func (this *ClipperStruct) op_Equality(a, b bool) bool {
	return true
}
func (this *ClipperStruct) SlopesEqual_3(b, c *TEdgeStruct, a bool) bool {
	if a {
		return op_Equality(Int128Mul(b.Delta.Y, c.Delta.X), Int128Mul(b.Delta.X, c.Delta.Y))
	} else {
		return 0 == (b.Delta.Y*c.Delta.X - b.Delta.X*c.Delta.Y)
	}
}
func (this *ClipperStruct) SlopesEqual_4(b, c, e *IntPoint, a bool) bool {
	if a {
		return op_Equality(Int128Mul(b.Y-c.Y, c.X-e.X), Int128Mul(b.X-c.X, c.Y-e.Y))
	} else {
		return 0 == ((b.Y-c.Y)*(c.X-e.X) - (b.X-c.X)*(c.Y-e.Y))
	}
}

type lStruct struct {
	Value bool
}

func PrintaEdge(aEdge *TEdgeStruct) {
	// this.logIfDebug("---------")
	// if aEdge != nil {
	// 	this.logIfDebug("aEdge.Bot:", aEdge.Bot)
	// }
	// if aEdge != nil && aEdge.Prev != nil {
	// 	this.logIfDebug("aEdge.Prev.Bot:", aEdge.Prev.Bot)
	// }
}
func (this *ClipperStruct) AddPath(aPolygon IntPolygon, b PolyType, c bool) bool {
	if !c {
		if b == ptClip {
			this.logIfDebug("AddPath: Open paths must be subject.")
		}
	}

	var eIdnex = len(aPolygon) - 1
	if c {
		for 0 < eIdnex && op_Equality_IntPoint(aPolygon[eIdnex], aPolygon[0]) {
			eIdnex--
		}
	}

	for 0 < eIdnex && op_Equality_IntPoint(aPolygon[eIdnex], aPolygon[eIdnex-1]) {
		eIdnex--
	}
	if c && 2 > eIdnex || !c && 1 > eIdnex {
		return false
	}
	f := []*TEdgeStruct{}
	for g := 0; g <= eIdnex; g++ {
		f = append(f, NewTEdge())
	}
	var hReal bool
	var h *bool = &hReal
	*h = true
	f[1].Curr.X = aPolygon[1].X
	f[1].Curr.Y = aPolygon[1].Y
	l := &lStruct{
		Value: this.m_UseFullRange,
	}

	this.RangeTest(aPolygon[0], l)

	this.m_UseFullRange = l.Value
	l.Value = this.m_UseFullRange

	this.RangeTest(aPolygon[eIdnex], l)
	this.m_UseFullRange = l.Value
	this.InitEdge(f[0], f[1], f[eIdnex], aPolygon[0])
	this.InitEdge(f[eIdnex], f[0], f[eIdnex-1], aPolygon[eIdnex])
	for g := eIdnex - 1; 1 <= g; g-- {
		this.RangeTest(aPolygon[g], l)
		this.m_UseFullRange = l.Value
		this.InitEdge(f[g], f[g+1], f[g-1], aPolygon[g])
	}
	var gEdge *TEdgeStruct = f[0]
	var aEdge *TEdgeStruct = f[0]
	PrintaEdge(aEdge)
	var eEdge *TEdgeStruct = f[0]
	for {
		if op_Equality_IntPoint(aEdge.Curr, aEdge.Next.Curr) {
			if aEdge == aEdge.Next {
				break
			}
			if aEdge == eEdge {
				eEdge = aEdge.Next
			}
			temp := this.RemoveEdge(aEdge)
			aEdge = temp
			PrintaEdge(aEdge)
			gEdge = temp
		} else {
			if aEdge.Prev == aEdge.Next {
				break
			} else if c && this.SlopesEqual_4(aEdge.Prev.Curr, aEdge.Curr, aEdge.Next.Curr, this.m_UseFullRange) && (!this.PreserveCollinear || !this.Pt2IsBetweenPt1AndPt3(aEdge.Prev.Curr, aEdge.Curr, aEdge.Next.Curr)) {
				if aEdge == eEdge {
					eEdge = aEdge.Next
				}
				aEdge = this.RemoveEdge(aEdge)
				PrintaEdge(aEdge)
				temp := aEdge.Prev
				aEdge = temp
				PrintaEdge(aEdge)
				gEdge = temp
				continue
			}
			aEdge = aEdge.Next
			PrintaEdge(aEdge)
			this.logIfDebug("触发 next")
			if aEdge == gEdge {
				break
			}
		}
	}
	if !c && aEdge == aEdge.Next || c && aEdge.Prev == aEdge.Next {
		return false
	}
	if !c {
		this.m_HasOpenPaths = true
		eEdge.Prev.OutIdx = ClipperBase_Skip
	}
	aEdge = eEdge
	PrintaEdge(aEdge)
	this.InitEdge2(aEdge, b)
	PrintaEdge(aEdge)
	aEdge = aEdge.Next
	PrintaEdge(aEdge)
	this.logIfDebug("触发InitEdge2 next")
	if *h {
		if aEdge.Curr.Y != eEdge.Curr.Y {
			*h = false
		}
	}
	for aEdge != eEdge {
		this.InitEdge2(aEdge, b)

		PrintaEdge(aEdge)
		this.logIfDebug("触发InitEdge2 next")
		aEdge = aEdge.Next
		PrintaEdge(aEdge)
		if *h {
			if aEdge.Curr.Y != eEdge.Curr.Y {
				*h = false
			}
		}
	}
	if *h {
		if c {
			return false
		}
		aEdge.Prev.OutIdx = ClipperBase_Skip
		if aEdge.Prev.Bot.X < aEdge.Prev.Top.X {
			this.ReverseHorizontal(aEdge.Prev)
		}
		bLocalMinima := LocalMinima()
		bLocalMinima.Next = nil
		bLocalMinima.Y = aEdge.Bot.Y
		bLocalMinima.LeftBound = nil
		bLocalMinima.RightBound = aEdge
		bLocalMinima.RightBound.Side = esRight
		for bLocalMinima.RightBound.WindDelta = 0; aEdge.Next.OutIdx != ClipperBase_Skip; {
			aEdge.NextInLML = aEdge.Next
			if aEdge.Bot.X != aEdge.Prev.Top.X {
				this.ReverseHorizontal(aEdge)
				PrintaEdge(aEdge)
			}
			aEdge = aEdge.Next
			PrintaEdge(aEdge)
			this.logIfDebug("触发ClipperBase_Skip next")
		}
		this.InsertLocalMinima(bLocalMinima)
		this.m_edges = append(this.m_edges, f)
		return true
	}
	this.m_edges = append(this.m_edges, f)
	var hEdge *TEdgeStruct = nil
	for {

		aEdge = this.FindNextLocMin(aEdge)
		if aEdge == hEdge {
			break
		} else {
			if nil == hEdge {
				hEdge = aEdge
			}
		}
		bLocalMinima := LocalMinima()
		bLocalMinima.Next = nil
		bLocalMinima.Y = aEdge.Bot.Y
		var fBool bool
		if aEdge.Dx < aEdge.Prev.Dx {
			bLocalMinima.LeftBound = aEdge.Prev
			bLocalMinima.RightBound = aEdge
			fBool = false
		} else {
			bLocalMinima.LeftBound = aEdge
			bLocalMinima.RightBound = aEdge.Prev
			fBool = true
		}
		bLocalMinima.LeftBound.Side = esLeft
		bLocalMinima.RightBound.Side = esRight
		if c {
			if bLocalMinima.LeftBound.Next == bLocalMinima.RightBound {
				bLocalMinima.LeftBound.WindDelta = -1
			} else {
				bLocalMinima.LeftBound.WindDelta = 1
			}
		} else {
			bLocalMinima.LeftBound.WindDelta = 0
		}

		bLocalMinima.RightBound.WindDelta = -bLocalMinima.LeftBound.WindDelta
		aEdge = this.ProcessBound(bLocalMinima.LeftBound, fBool)
		eEdge = this.ProcessBound(bLocalMinima.RightBound, !fBool)
		if bLocalMinima.LeftBound.OutIdx == ClipperBase_Skip {
			bLocalMinima.LeftBound = nil
		} else {
			if bLocalMinima.RightBound.OutIdx == ClipperBase_Skip {
				bLocalMinima.RightBound = nil
			}
		}

		this.InsertLocalMinima(bLocalMinima)
		if !fBool {
			aEdge = eEdge
		}
	}
	return true
}

type LocalMinimaStruct struct {
	Y          int64
	Next       *LocalMinimaStruct
	RightBound *TEdgeStruct
	LeftBound  *TEdgeStruct
}

func LocalMinima() *LocalMinimaStruct {
	return &LocalMinimaStruct{
		Y:          0,
		Next:       nil,
		RightBound: nil,
		LeftBound:  nil,
	}
}
func (this *ClipperStruct) InsertLocalMinima(a *LocalMinimaStruct) {
	//this.logIfDebug("a:", a.Y)
	if nil == this.m_MinimaList {
		this.m_MinimaList = a
	} else if a.Y >= this.m_MinimaList.Y {
		a.Next = this.m_MinimaList
		this.m_MinimaList = a
	} else {
		b := this.m_MinimaList
		for nil != b.Next && a.Y < b.Next.Y {
			b = b.Next
		}
		a.Next = b.Next
		b.Next = a
	}
}

func (this *ClipperStruct) RemoveEdge(aEdge *TEdgeStruct) *TEdgeStruct {
	aEdge.Prev.Next = aEdge.Next
	aEdge.Next.Prev = aEdge.Prev
	bEdge := aEdge.Next
	aEdge.Prev = nil
	return bEdge
}
func (this *ClipperStruct) ReverseHorizontal(a *TEdgeStruct) {
	var b int64 = a.Top.X
	a.Top.X = a.Bot.X
	a.Bot.X = b
}
func (this *ClipperStruct) ProcessBound(a *TEdgeStruct, b bool) *TEdgeStruct {
	this.logIfDebug("go in ProcessBound")
	this.logIfDebug(a.Curr)
	// if a.Curr.X == 4506 {
	// 	this.logIfDebug("should stop")
	// }
	c := a
	e := a
	var f int64
	var fEdge *TEdgeStruct
	if a.Dx == ClipperBase_horizontal {
		if b {
			f = a.Prev.Bot.X

		} else {
			f = a.Next.Bot.X
		}
		if a.Bot.X != f {
			this.ReverseHorizontal(a)
		}

	}
	if e.OutIdx != ClipperBase_Skip {
		if b {
			for e.Top.Y == e.Next.Bot.Y && e.Next.OutIdx != ClipperBase_Skip {
				e = e.Next
			}
			if e.Dx == ClipperBase_horizontal && e.Next.OutIdx != ClipperBase_Skip {
				for fEdge = e; fEdge.Prev.Dx == ClipperBase_horizontal; {
					fEdge = fEdge.Prev
				}
				if fEdge.Prev.Top.X == e.Next.Top.X {
					if !b {
						e = fEdge.Prev
					}
				} else {
					if fEdge.Prev.Top.X > e.Next.Top.X {
						e = fEdge.Prev
					}
				}
			}
			for a != e {
				a.NextInLML = a.Next
				if a.Dx == ClipperBase_horizontal {
					if a != c {
						if a.Bot.X != a.Prev.Top.X {
							this.ReverseHorizontal(a)
						}
					}
				}
				a = a.Next
			}

			if a.Dx == ClipperBase_horizontal {
				if a != c {
					if a.Bot.X != a.Prev.Top.X {
						this.ReverseHorizontal(a)
					}
				}
			}
			e = e.Next
		} else {
			for e.Top.Y == e.Prev.Bot.Y && e.Prev.OutIdx != ClipperBase_Skip {
				e = e.Prev
			}

			if e.Dx == ClipperBase_horizontal && e.Prev.OutIdx != ClipperBase_Skip {
				for fEdge = e; fEdge.Next.Dx == ClipperBase_horizontal; {
					fEdge = fEdge.Next
				}
				if fEdge.Next.Top.X == e.Prev.Top.X {
					if !b {
						e = fEdge.Next
					}
				} else {
					if fEdge.Next.Top.X > e.Prev.Top.X {
						e = fEdge.Next
					}
				}
			}
			for a != e {
				a.NextInLML = a.Prev
				if a.Dx == ClipperBase_horizontal {
					if a != c {
						if a.Bot.X != a.Next.Top.X {
							this.ReverseHorizontal(a)
						}
					}
				}
				a = a.Prev
			}
			if a.Dx == ClipperBase_horizontal {
				if a != c {
					if a.Bot.X != a.Next.Top.X {
						this.ReverseHorizontal(a)
					}
				}
			}
			e = e.Prev
		}
	}
	//this.logIfDebug("e.OutIdx == ClipperBase_Skip")
	if e.OutIdx == ClipperBase_Skip {
		a = e
		if b {
			for a.Top.Y == a.Next.Bot.Y {
				a = a.Next
			}

			for a != e && a.Dx == ClipperBase_horizontal {
				a = a.Prev
			}

		} else {
			for a.Top.Y == a.Prev.Bot.Y {
				a = a.Prev
			}

			for a != e && a.Dx == ClipperBase_horizontal {
				a = a.Next
			}

		}
		if a == e {
			if b {
				e = a.Next
			} else {
				e = a.Prev
			}

		} else {
			if b {
				a = e.Next
			} else {
				a = e.Prev
			}
			c := LocalMinima()
			c.Next = nil
			c.Y = a.Bot.Y
			c.LeftBound = nil
			c.RightBound = a
			c.RightBound.WindDelta = 0
			e = this.ProcessBound(c.RightBound, b)
			this.InsertLocalMinima(c)
		}

	}
	return e
}
func (this *ClipperStruct) FindNextLocMin(aEdge *TEdgeStruct) *TEdgeStruct {
	var bEdge *TEdgeStruct
	for {
		for op_InEquality_IntPoint(aEdge.Bot, aEdge.Prev.Bot) || op_Equality_IntPoint(aEdge.Curr, aEdge.Top) {
			aEdge = aEdge.Next
		}

		if aEdge.Dx != ClipperBase_horizontal && aEdge.Prev.Dx != ClipperBase_horizontal {
			break
		}

		for aEdge.Prev.Dx == ClipperBase_horizontal {
			aEdge = aEdge.Prev
		}

		for bEdge = aEdge; aEdge.Dx == ClipperBase_horizontal; {
			aEdge = aEdge.Next
		}

		if aEdge.Top.Y != aEdge.Prev.Bot.Y {
			if bEdge.Prev.Bot.X < aEdge.Bot.X {
				aEdge = bEdge
			}
			break
		}
	}
	return aEdge
}
