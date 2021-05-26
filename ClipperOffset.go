package nest

import (
	"math"

	. "github.com/mojinfu/point"
)

func (this *ClipperOffsetStruct) FixOrientations() {
	if 0 <= this.m_lowest.X && !Orientation(this.m_polyNodes.Childs()[this.m_lowest.X].m_polygon) {
		for a := 0; a < len(this.m_polyNodes.m_Childs); a++ {
			var b = this.m_polyNodes.Childs()[a]
			if b.m_endtype == etClosedPolygon || b.m_endtype == etClosedLine && Orientation(b.m_polygon) {
				b.m_polygon = IntPointReverse(b.m_polygon)
			}
		}
	} else {
		for a := 0; a < len(this.m_polyNodes.m_Childs); a++ {
			var b = this.m_polyNodes.Childs()[a]
			if !(b.m_endtype != etClosedLine || Orientation(b.m_polygon)) {
				b.m_polygon = IntPointReverse(b.m_polygon)
			}
		}

	}
}
func (this *ClipperOffsetStruct) GetUnitNormal(a, b *IntPoint) *Point {
	c := float64(b.X - a.X)
	e := float64(b.Y - a.Y)
	if 0 == c && 0 == e {
		return DoublePoint(0, 0)
	}

	var f = 1 / math.Sqrt(c*c+e*e)
	return DoublePoint(e*f, -1*(c*f))
}
func Round(a float64) int64 {
	if 0 > a {
		return int64(math.Ceil(a - 0.5))
	} else {
		return int64(math.Round(a))
	}

}
func near_zero(a float64) bool {
	return a > -1*ClipperBase_tolerance && a < ClipperBase_tolerance
}

func DoublePoint(X, Y float64) *Point {
	return &Point{
		X: X,
		Y: Y,
	}
}

func (this *ClipperOffsetStruct) DoMiter(a, b int, c float64) {
	c = this.m_delta / c
	this.m_destPoly = append(this.m_destPoly, &IntPoint{
		X: Round(float64(this.m_srcPoly[a].X) + (this.m_normals[b].X+this.m_normals[a].X)*c),
		Y: Round(float64(this.m_srcPoly[a].Y) + (this.m_normals[b].Y+this.m_normals[a].Y)*c),
	})
}
func (this *ClipperOffsetStruct) DoSquare(a, b int) {
	var c = math.Tan(math.Atan2(this.m_sinA, this.m_normals[b].X*this.m_normals[a].X+this.m_normals[b].Y*this.m_normals[a].Y) / 4)
	this.m_destPoly = append(this.m_destPoly, &IntPoint{
		X: Round(float64(this.m_srcPoly[a].X) + this.m_delta*(this.m_normals[b].X-this.m_normals[b].Y*c)),
		Y: Round(float64(this.m_srcPoly[a].Y) + this.m_delta*(this.m_normals[b].Y+this.m_normals[b].X*c)),
	})
	this.m_destPoly = append(this.m_destPoly, &IntPoint{
		X: Round(float64(this.m_srcPoly[a].X) + this.m_delta*(this.m_normals[a].X+this.m_normals[a].Y*c)),
		Y: Round(float64(this.m_srcPoly[a].Y) + this.m_delta*(this.m_normals[a].Y-this.m_normals[a].X*c)),
	})
}
func (this *ClipperOffsetStruct) DoRound(a, b int) {
	cFloat := math.Atan2(this.m_sinA, this.m_normals[b].X*this.m_normals[a].X+this.m_normals[b].Y*this.m_normals[a].Y) //saya
	c := int(Round(this.m_StepsPerRad * math.Abs(cFloat)))
	e := this.m_normals[b].X
	f := this.m_normals[b].Y
	var g float64
	var h int = 0
	for ; h < c; h++ {
		this.m_destPoly = append(this.m_destPoly, &IntPoint{X: Round(float64(this.m_srcPoly[a].X) + e*this.m_delta), Y: Round(float64(this.m_srcPoly[a].Y) + f*this.m_delta)})
		g = e
		e = e*this.m_cos - this.m_sin*f
		f = g*this.m_sin + f*this.m_cos
	}

	this.m_destPoly = append(this.m_destPoly, &IntPoint{X: Round(float64(this.m_srcPoly[a].X) + this.m_normals[a].X*this.m_delta), Y: Round(float64(this.m_srcPoly[a].Y) + this.m_normals[a].Y*this.m_delta)})
}

func (this *ClipperOffsetStruct) OffsetPoint(a, b int, c int64) int {

	this.m_sinA = this.m_normals[b].X*this.m_normals[a].Y - this.m_normals[a].X*this.m_normals[b].Y
	if 5E-5 > this.m_sinA && -5E-5 < this.m_sinA {
		return b
	}

	if 1 < this.m_sinA {
		this.m_sinA = 1
	} else {
		if -1 > this.m_sinA {
			this.m_sinA = -1
		}
	}
	//此处的round 是否造成最终的精度问题
	if 0 > this.m_sinA*this.m_delta {
		this.m_destPoly = append(this.m_destPoly, &IntPoint{X: Round(float64(this.m_srcPoly[a].X) + this.m_normals[b].X*this.m_delta), Y: Round(float64(this.m_srcPoly[a].Y) + this.m_normals[b].Y*this.m_delta)})

		this.m_destPoly = append(this.m_destPoly, &IntPoint{X: Round(float64(this.m_srcPoly[a].X)), Y: Round(float64(this.m_srcPoly[a].Y))})

		this.m_destPoly = append(this.m_destPoly, &IntPoint{X: Round(float64(this.m_srcPoly[a].X) + this.m_normals[a].X*this.m_delta), Y: Round(float64(this.m_srcPoly[a].Y) + this.m_normals[a].Y*this.m_delta)})
	} else {
		switch c {
		case jtMiter:
			{
				cFloat := 1 + (this.m_normals[a].X*this.m_normals[b].X + this.m_normals[a].Y*this.m_normals[b].Y)
				if cFloat >= this.m_miterLim {
					this.DoMiter(a, b, cFloat)
				} else {
					this.DoSquare(a, b)
				}
				break
			}

		case jtSquare:
			{
				this.DoSquare(a, b)
				break
			}

		case jtRound:
			{
				this.DoRound(a, b)
			}

		}
	}

	return a
}
func (this *ClipperOffsetStruct) DoOffset(a float64) {
	this.m_destPolys = []IntPolygon{}
	this.m_delta = a
	if near_zero(a) {
		// for b := 0; b < this.m_polyNodes.ChildCount(); b++ {
		// 	var c = this.m_polyNodes.Childs()[b];
		// 	c.m_endtype == etClosedPolygon && this.m_destPolys.push(c.m_polygon)
		// }
	} else {
		if 2 < this.MiterLimit {
			this.m_miterLim = 2 / (this.MiterLimit * this.MiterLimit)
		} else {
			this.m_miterLim = 0.5
		}
		var b float64
		if 0 >= this.ArcTolerance {
			b = ClipperOffset_def_arc_tolerance
		} else {
			if this.ArcTolerance > math.Abs(a)*ClipperOffset_def_arc_tolerance {
				b = math.Abs(a) * ClipperOffset_def_arc_tolerance
			} else {
				b = this.ArcTolerance
			}
		}
		e := 3.14159265358979 / math.Acos(1-b/math.Abs(a))

		this.m_sin = math.Sin(ClipperOffset_two_pi / e)
		this.m_cos = math.Cos(ClipperOffset_two_pi / e)
		this.m_StepsPerRad = e / ClipperOffset_two_pi
		if 0 > a {
			this.m_sin = -1 * this.m_sin
		}
		for b := 0; b < len(this.m_polyNodes.Childs()); b++ {
			c := this.m_polyNodes.Childs()[b]
			this.m_srcPoly = c.m_polygon //saya
			var f = len(this.m_srcPoly)
			if !(0 == f || 0 >= a && (3 > f || c.m_endtype != etClosedPolygon)) {
				this.m_destPoly = IntPolygon{}
				if 1 == f {
					if c.m_jointype == jtRound {
						var c float64 = 1
						var f float64 = 0
						var g float64 = 1
						for ; g <= e; g++ {
							this.m_destPoly = append(this.m_destPoly,
								&IntPoint{
									X: Round(float64(this.m_srcPoly[0].X) + c*a),
									Y: Round(float64(this.m_srcPoly[0].Y) + f*a),
								})
							var h = c
							c = c*this.m_cos - this.m_sin*f
							f = h*this.m_sin + f*this.m_cos
						}
					} else {
						var c float64 = -1
						var f float64 = -1

						g := 0
						for ; 4 > g; g++ {
							this.m_destPoly = append(this.m_destPoly, &IntPoint{
								X: Round(float64(this.m_srcPoly[0].X) + c*a),
								Y: Round(float64(this.m_srcPoly[0].Y) + f*a)})

							if 0 > c {
								c = 1
							} else {
								if 0 > f {
									f = 1
								} else {
									c = -1
								}
							}

						}
					}
				} else {
					g := 0
					this.m_normals = []*Point{}
					for ; g < f-1; g++ {
						this.m_normals = append(this.m_normals, this.GetUnitNormal(this.m_srcPoly[g], this.m_srcPoly[g+1]))
					}

					if c.m_endtype == etClosedLine || c.m_endtype == etClosedPolygon {
						this.m_normals = append(this.m_normals, this.GetUnitNormal(this.m_srcPoly[f-1], this.m_srcPoly[0]))
					} else {
						this.m_normals = append(this.m_normals, DoublePoint(this.m_normals[f-2].X, this.m_normals[f-2].Y))
					}
					if c.m_endtype == etClosedPolygon {
						h := f - 1
						var g int = 0
						for ; g < f; g++ {
							h = this.OffsetPoint(g, h, c.m_jointype)
						}

					} else if c.m_endtype == etClosedLine {
						h := f - 1
						g = 0
						for ; g < f; g++ {
							h = this.OffsetPoint(g, h, c.m_jointype)
						}

						this.m_destPolys = append(this.m_destPolys, this.m_destPoly)
						this.m_destPoly = IntPolygon{}
						hPoint := this.m_normals[f-1]
						g = f - 1
						for ; 0 < g; g-- {
							this.m_normals[g] = DoublePoint(-this.m_normals[g-1].X, -this.m_normals[g-1].Y)
						}

						this.m_normals[0] = DoublePoint(-1*hPoint.X, -hPoint.Y)
						h = 0
						g = f - 1
						for ; 0 <= g; g-- {
							h = this.OffsetPoint(g, h, c.m_jointype)
						}

					} else {
						h := 0
						g = 1
						for ; g < f-1; g++ {
							h = this.OffsetPoint(g, h, c.m_jointype)
						}
						if c.m_endtype == etOpenButt {
							g = f - 1
							hPoint := &IntPoint{
								X: Round(float64(this.m_srcPoly[g].X) + this.m_normals[g].X*a),
								Y: Round(float64(this.m_srcPoly[g].Y) + this.m_normals[g].Y*a),
							}

							this.m_destPoly = append(this.m_destPoly, hPoint)
							hPoint = &IntPoint{

								X: Round(float64(this.m_srcPoly[g].X) - this.m_normals[g].X*a), Y: Round(float64(this.m_srcPoly[g].Y) - this.m_normals[g].Y*a),
							}

							this.m_destPoly = append(this.m_destPoly, hPoint)
						} else {
							g = f - 1
							h = f - 2
							this.m_sinA = 0
							this.m_normals[g] = DoublePoint(-this.m_normals[g].X, -this.m_normals[g].Y)
							if c.m_endtype == etOpenSquare {
								this.DoSquare(g, h)
							} else {
								this.DoRound(g, h)
							}
						}
						g = f - 1
						for ; 0 < g; g-- {
							this.m_normals[g] = DoublePoint(-this.m_normals[g-1].X, -this.m_normals[g-1].Y)
						}

						this.m_normals[0] = DoublePoint(-this.m_normals[1].X, -this.m_normals[1].Y)
						h = f - 1
						g = h - 1
						for ; 0 < g; g-- {
							h = this.OffsetPoint(g, h, c.m_jointype)
						}

						if c.m_endtype == etOpenButt {
							hPoint := &IntPoint{
								X: Round(float64(this.m_srcPoly[0].X) - this.m_normals[0].X*a),
								Y: Round(float64(this.m_srcPoly[0].Y) - this.m_normals[0].Y*a)}

							this.m_destPoly = append(this.m_destPoly, hPoint)
							hPoint = &IntPoint{
								X: Round(float64(this.m_srcPoly[0].X) + this.m_normals[0].X*a),
								Y: Round(float64(this.m_srcPoly[0].Y) + this.m_normals[0].Y*a)}

							this.m_destPoly = append(this.m_destPoly, hPoint)
						} else {
							this.m_sinA = 0
							if c.m_endtype == etOpenSquare {
								this.DoSquare(0, 1)
							} else {
								this.DoRound(0, 1)
							}
						}
					}
				}
				this.m_destPolys = append(this.m_destPolys, this.m_destPoly)
			}
		}
	}
	//saya
}
func (this *ClipperOffsetStruct) ExecutePath(c float64) []IntPolygon {
	//d.Clear(b);
	b := []IntPolygon{}
	this.FixOrientations()
	this.DoOffset(c)

	a := Clipper(0)
	a.AddPaths(this.m_destPolys, ptSubject, true)
	if 0 < c {
		b, _ = a.Execute(ctUnion, pftPositive, pftPositive)
	} else {
		c := a.GetBounds(this.m_destPolys)
		e := IntPolygon{}
		e = append(e, &IntPoint{X: c.left - 10, Y: c.bottom + 10})
		e = append(e, &IntPoint{X: c.right + 10, Y: c.bottom + 10})
		e = append(e, &IntPoint{X: c.right + 10, Y: c.top - 10})
		e = append(e, &IntPoint{X: c.left - 10, Y: c.top - 10})
		a.AddPath(e, ptSubject, true)
		a.ReverseSolution = true
		b, _ = a.Execute(ctUnion, pftNegative, pftNegative)
		if 0 < len(b) {
			b = b[1:]
		}
	}
	return b
}
func (this *ClipperOffsetStruct) AddPath(a IntPolygon, b, c int64) {
	e := len(a) - 1
	if !(0 > e) {
		f := NewPolyNode()
		f.m_jointype = b
		f.m_endtype = c
		if c == etClosedLine || c == etClosedPolygon {
			for 0 < e && op_Equality_IntPoint(a[0], a[e]) {
				e--
			}

		}

		f.m_polygon = append(f.m_polygon, a[0])
		var g int64 = 0
		b = 0
		for h := 1; h <= e; h++ {
			if op_InEquality_IntPoint(f.m_polygon[g], a[h]) {
				g++
				f.m_polygon = append(f.m_polygon, a[h])
				if a[h].Y > f.m_polygon[b].Y || a[h].Y == f.m_polygon[b].Y && a[h].X < f.m_polygon[b].X {
					b = g
				}

			}

		}
		if !(c == etClosedPolygon && 2 > g || c != etClosedPolygon && 0 > g) {
			this.m_polyNodes.AddChild(f)
			if c == etClosedPolygon {
				if 0 > this.m_lowest.X {
					this.m_lowest = &IntPoint{
						X: 0,
						Y: b,
					}
				} else {
					aPoint := this.m_polyNodes.m_Childs[this.m_lowest.X].m_polygon[this.m_lowest.Y]
					if f.m_polygon[b].Y > aPoint.Y || f.m_polygon[b].Y == aPoint.Y && f.m_polygon[b].X < aPoint.X {
						this.m_lowest = &IntPoint{
							X: int64(len(this.m_polyNodes.m_Childs)) - 1,
							Y: b,
						}
					}

				}

			}
		}

	}
}
