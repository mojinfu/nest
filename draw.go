package nest

import (
	"math/rand"
	"strings"

	"github.com/fogleman/gg"
	. "github.com/mojinfu/nest/point"
)

type Paper struct {
	dc              *gg.Context
	LineWidth       float64
	paperSavingPath string
}

func (this *SVG) NewPaper() *Paper {
	myPaper := &Paper{
		LineWidth:       5,
		paperSavingPath: this.config.PaperSavePath,
	}
	dc := gg.NewContext(int(this.maxX), int(this.maxY)) //上下文，含长和宽
	dc.SetRGB(1, 1, 1)                                  //设置当前色
	dc.Clear()
	myPaper.dc = dc
	return myPaper
}
func (this *Paper) AddPolygon(myPolygon []*Point, isWart bool) {
	if len(myPolygon) <= 2 {
		return
	}
	if isWart {
		this.dc.SetRGBA(1, 0, 0, 1)
	} else {
		this.dc.SetRGBA(0+rand.Float64()*0.3, rand.Float64()*0.3, rand.Float64()*0.5, 0.5) //saya 应该是1
	}

	this.dc.SetLineWidth(this.LineWidth)

	for index := range myPolygon {
		if index < len(myPolygon)-1 {
			this.dc.DrawLine(myPolygon[index].X, myPolygon[index].Y, myPolygon[index+1].X, myPolygon[index+1].Y) //画线
		} else {
			this.dc.DrawLine(myPolygon[index].X, myPolygon[index].Y, myPolygon[0].X, myPolygon[0].Y) //画线
		}
		this.dc.Stroke()
	}
}

func (this *Paper) AddDebugNFP(myPolygon []*Point) {
	if len(myPolygon) == 1 {
		this.dc.DrawPoint(myPolygon[0].X, myPolygon[0].Y, 2)
		return
	}
	if len(myPolygon) == 2 {
		this.dc.DrawLine(myPolygon[0].X, myPolygon[0].Y, myPolygon[1].X, myPolygon[1].Y)
		return
	}
	this.dc.SetRGBA(0+rand.Float64()*0.3, rand.Float64()*0.3, rand.Float64()*0.5, 1)
	this.dc.SetLineWidth(this.LineWidth)
	for index := range myPolygon {
		if index < len(myPolygon)-1 {
			this.dc.DrawLine(myPolygon[index].X, myPolygon[index].Y, myPolygon[index+1].X, myPolygon[index+1].Y) //画线
		} else {
			this.dc.DrawLine(myPolygon[index].X, myPolygon[index].Y, myPolygon[0].X, myPolygon[0].Y) //画线
		}
		this.dc.Stroke()
	}
}

func (this *Paper) Draw(name string) {
	if this.paperSavingPath == "" {
		this.dc.SavePNG("./output/" + name + ".png")
	} else {
		if strings.HasSuffix(this.paperSavingPath, "/") {
			this.dc.SavePNG(this.paperSavingPath + name + ".png")
			this.dc.SavePNG(this.paperSavingPath + "best" + ".png")
		} else {
			this.dc.SavePNG(this.paperSavingPath + "/" + name + ".png")
			this.dc.SavePNG(this.paperSavingPath + "/" + "best" + ".png")
		}

	}
}
