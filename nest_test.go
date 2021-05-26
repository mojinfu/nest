package nest

import (
	"fmt"
	"testing"

	. "github.com/mojinfu/point"
)

func Test_Intersect(t *testing.T) {
	part := []*Point{

		&Point{362.09, -1.372, false},
		&Point{362.5, 0, false},
		&Point{362.5, 25, false},
		&Point{361.372, 27.09, false},
		&Point{360, 27.5, false},
		&Point{335.143, 27.5, false},
		&Point{304.349, 31.087, false},
		&Point{273.326, 33.892, false},
		&Point{242.241, 35.896, false},
		&Point{211.115, 37.099, false},
		&Point{179.968, 37.5, false},
		&Point{148.82, 37.097, false},
		&Point{117.695, 35.893, false},
		&Point{86.61, 33.887, false},
		&Point{55.587, 31.08, false},
		&Point{24.857, 27.5, false},
		&Point{0, 27.5, false},
		&Point{-2.09, 26.372, false},
		&Point{-2.5, 25, false},
		&Point{-2.5, 0, false},
		&Point{-1.372, -2.09, false},
		&Point{0, -2.5, false},
		&Point{25.289, -2.483, false},
		&Point{56.118, 1.109, false},
		&Point{87.026, 3.904, false},
		&Point{117.997, 5.902, false},
		&Point{148.969, 7.099, false},
		&Point{180, 7.5, false},
		&Point{211.03, 7.099, false},
		&Point{242.003, 5.902, false},
		&Point{272.974, 3.904, false},
		&Point{303.882, 1.108, false},
		&Point{335, -2.5, false},
		&Point{360, -2.5, false},
	}
	mySVG := NewSVG(PublicConfig)
	a := &polygonWithOffset{
		Polygon: []*Point{
			&Point{2000, 2000, false},
			&Point{0, 2000, false},
			&Point{0, 0, false},
			&Point{2000, 0, false},
		}}
	b := &polygonWithOffset{
		offsetx: 1999.59 - part[0].X,
		offsety: 1971.128 - part[0].Y,
		Polygon: part}

	c := &polygonWithOffset{
		offsetx: 0,
		offsety: 0,
		Polygon: []*Point{
			&Point{0, 500, false},
			&Point{0, 0, false},
			&Point{500, 0, false},
			&Point{500, -500, false},
			&Point{2500, -500, false},
			&Point{2500, 500, false},
			&Point{2000, 500, false},
			&Point{2000, 0, false},
			&Point{1000, 0, false},
			&Point{1000, 500, false},
		}}

	fmt.Println(mySVG.intersect(a, b))
	fmt.Println("ac:", mySVG.intersect(a, c))
	DrawWithOffset("相交问题图", b.offsetx, b.offsety, (a.Polygon), b.Polygon)
	DrawWithOffset("ac相交问题图", 700, 700, TranslatePolygon(a.Polygon, 700, 700), c.Polygon)

}

func Test_Fit(t *testing.T) {
	mySVG := NewSVG(PublicConfig)
	bin := []*Point{
		&Point{2000, 2000, false},
		&Point{0, 2000, false},
		&Point{0, 0, false},
		&Point{2000, 0, false},
	}

	part := []*Point{
		&Point{362.09, -1.372, false},
		&Point{362.5, 0, false},
		&Point{362.5, 25, false},
		&Point{361.372, 27.09, false},
		&Point{360, 27.5, false},
		&Point{335.143, 27.5, false},
		&Point{304.349, 31.087, false},
		&Point{273.326, 33.892, false},
		&Point{242.241, 35.896, false},
		&Point{211.115, 37.099, false},
		&Point{179.968, 37.5, false},
		&Point{148.82, 37.097, false},
		&Point{117.695, 35.893, false},
		&Point{86.61, 33.887, false},
		&Point{55.587, 31.08, false},
		&Point{24.857, 27.5, false},
		&Point{0, 27.5, false},
		&Point{-2.09, 26.372, false},
		&Point{-2.5, 25, false},
		&Point{-2.5, 0, false},
		&Point{-1.372, -2.09, false},
		&Point{0, -2.5, false},
		&Point{25.289, -2.483, false},
		&Point{56.118, 1.109, false},
		&Point{87.026, 3.904, false},
		&Point{117.997, 5.902, false},
		&Point{148.969, 7.099, false},
		&Point{180, 7.5, false},
		&Point{211.03, 7.099, false},
		&Point{242.003, 5.902, false},
		&Point{272.974, 3.904, false},
		&Point{303.882, 1.108, false},
		&Point{335, -2.5, false},
		&Point{360, -2.5, false},
	}
	// part := []*Point{

	// 	&Point{179.968, 37.5, false},
	// 	&Point{148.82, 37.097, false},
	// 	&Point{117.695, 35.893, false},
	// 	&Point{86.61, 33.887, false},
	// 	&Point{55.587, 31.08, false},
	// 	&Point{24.857, 27.5, false},
	// 	&Point{0, 27.5, false},
	// 	&Point{-2.09, 26.372, false},
	// 	&Point{-2.5, 25, false},
	// 	&Point{-2.5, 0, false},
	// 	&Point{-1.372, -2.09, false},
	// 	&Point{0, -2.5, false},
	// 	&Point{25.289, -2.483, false},
	// 	&Point{56.118, 1.109, false},
	// 	&Point{87.026, 3.904, false},
	// 	&Point{117.997, 5.902, false},
	// 	&Point{148.969, 7.099, false},
	// 	&Point{180, 7.5, false},
	// 	&Point{211.03, 7.099, false},
	// 	&Point{242.003, 5.902, false},
	// 	&Point{272.974, 3.904, false},
	// 	&Point{303.882, 1.108, false},
	// 	&Point{335, -2.5, false},
	// 	&Point{360, -2.5, false},
	// 	&Point{362.09, -1.372, false},
	// 	&Point{362.5, 0, false},
	// 	&Point{362.5, 25, false},
	// 	&Point{361.372, 27.09, false},
	// 	&Point{360, 27.5, false},
	// 	&Point{335.143, 27.5, false},
	// 	&Point{304.349, 31.087, false},
	// 	&Point{273.326, 33.892, false},
	// 	&Point{242.241, 35.896, false},
	// 	&Point{211.115, 37.099, false},
	// }

	// for i, p := range []*Point{
	// 	// &Point{X: 1800, Y: 1600, Marked: false},
	// 	// &Point{X: 200, Y: 1600, Marked: false},
	// 	// &Point{X: 200, Y: 0, Marked: false},
	// 	// &Point{X: 1800, Y: 0, Marked: false}

	// 	&Point{X: 1999.59, Y: 1971.128, Marked: false},
	// 	&Point{X: 364.5899999999999, Y: 1971.128, Marked: false},
	// } {
	// 	DrawWithOffset("NPF 拐点"+fmt.Sprintf("%d", i), p.X-part[0].X, p.Y-part[0].Y, bin, part)
	// }
	//DrawWithOffset("NPF 拐点A", 1999.59-part[0].X, 1971.128-part[0].Y, bin, part)
	//DrawWithOffset("NPF 拐点B", 364.5899999999999-part[0].X, 1971.128-part[0].Y, bin, part)

	poly := mySVG.noFitPolygon(
		bin, part,
		true,
		false,
	)
	if len(poly) == 0 {
		fmt.Println("nfp 失败")
		DrawAtO(0, bin)
		DrawAtO(1, part)
	} else {
		fmt.Println("nfp 成功")
		DrawAtO(10, bin)
		DrawAtO(11, part)
	}
	maxx, maxy := CaluMaxXMaxY(bin)
	paper := NewPaper(int(maxx)+100, int(maxy)+100, "./debugger/output/")
	paper.AddPolygon(bin, false)
	for i := range poly {
		DrawAtO(i, poly[i])
		for j := range poly[i] {
			fmt.Printf("%+v\n", poly[i][j])
			partNew := TranslatePolygon(part, poly[i][j].X-part[0].X, poly[i][j].Y-part[0].Y)
			paper.AddPolygon(partNew, true)
		}
	}
	paper.DrawPartPreview("NFP最终形式", "0")
}
func DrawAtO(name int, a []*Point) {
	minx, miny := CaluMinXMinY(a)
	a = TranslatePolygon(a, -minx, -miny)
	maxx, maxy := CaluMaxXMaxY(a)
	paper := NewPaper(int(maxx)+100, int(maxy)+100, "./debugger/output/")
	paper.AddPolygon(a, true)
	paper.DrawPartPreview(" dxf解析形状", fmt.Sprintf("%d", name))
}
func DrawWithOffset(name string, xOffset, yOffset float64, bin []*Point, aList ...[]*Point) {
	maxx, maxy := CaluMaxXMaxY(bin)
	paper := NewPaper(int(maxx)+2000, int(maxy)+2000, "./debugger/output/")
	paper.AddPolygon(bin, false)
	for index := range aList {
		a := TranslatePolygon(aList[index], xOffset, yOffset)
		paper.AddPolygon(a, true)
	}
	paper.DrawPartPreview(name, "0")
}
