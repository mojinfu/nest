package nest

import (
	"math"
	"strings"

	"github.com/fogleman/gg"
	. "github.com/mojinfu/point"
)

type Paper struct {
	dc              *gg.Context
	LineWidth       float64
	PaperSavingPath string
	polyNum         int
	isEmpty         bool
}

func (this *Paper) IsEmpty() bool {
	return this.isEmpty
}
func (this *SVG) NewPaper() []*Paper {
	myPaperList := []*Paper{}
	for index := range this.maxX {
		myPaper := &Paper{
			LineWidth:       math.Min(math.Min(float64(this.maxX[index])/100+1, 5), float64(this.maxY[index])/100+1),
			PaperSavingPath: this.config.PaperSavePath,
			isEmpty:         true,
		}
		dc := gg.NewContext(int(this.maxX[index]), int(this.maxY[index])) //上下文，含长和宽
		dc.SetRGBA(1, 1, 1, 1)                                            //设置当前色
		dc.Clear()
		myPaper.dc = dc
		myPaperList = append(myPaperList, myPaper)
	}
	return myPaperList
}
func NewPaper(x, y int, path string) *Paper {
	myPaper := &Paper{
		LineWidth:       math.Min(math.Min(float64(x)/100+1, 5), float64(y)/100+1),
		PaperSavingPath: path,
		isEmpty:         false,
	}
	dc := gg.NewContext(int(x), int(y)) //上下文，含长和宽
	dc.SetRGBA(1, 1, 1, 0)              //设置当前色
	dc.Clear()
	myPaper.dc = dc
	return myPaper
}
func (this *Paper) AddBinPolygon(myPolygon []*Point, isWart bool) {
	if len(myPolygon) <= 2 {
		return
	}
	if isWart {
		this.dc.SetRGBA(1, 0, 0, 1)
	} else {
		//this.dc.SetRGBA(0+rand.Float64()*0.3, rand.Float64()*0.3, rand.Float64()*0.5, 1) //saya 应该是1
		this.dc.SetRGBA(0, 0, 0, 1)
	}

	this.dc.SetLineWidth(this.LineWidth / 3)
	// for index := range myPolygon {
	// 	this.dc.DrawPoint(myPolygon[index].X, myPolygon[index].Y, this.LineWidth/2)
	// 	this.dc.Stroke()
	// }
	// this.dc.SetRGBA(0.5, 0.5, 0, 1)
	// this.dc.SetLineWidth(this.LineWidth / 3)
	for index := range myPolygon {
		if index < len(myPolygon)-1 {
			//this.dc.DrawPoint(myPolygon[index].X, myPolygon[index].Y, this.LineWidth/2)
			this.dc.DrawLine(myPolygon[index].X, myPolygon[index].Y, myPolygon[index+1].X, myPolygon[index+1].Y) //画线
		} else {
			this.dc.DrawLine(myPolygon[index].X, myPolygon[index].Y, myPolygon[0].X, myPolygon[0].Y) //画线
		}
		this.dc.Stroke()
	}
}
func (this *Paper) AddPolygon(myPolygon []*Point, isWart bool) {
	if len(myPolygon) <= 2 {
		return
	}
	this.polyNum++
	if isWart {
		this.dc.SetRGBA(1, 0, 0, 1)
	} else {
		//this.dc.SetRGBA(0+rand.Float64()*0.3, rand.Float64()*0.3, rand.Float64()*0.5, 1) //saya 应该是1
		this.dc.SetRGBA(0, 0, 0, 1)
	}

	this.dc.SetLineWidth(this.LineWidth / 3)
	// for index := range myPolygon {
	// 	this.dc.DrawPoint(myPolygon[index].X, myPolygon[index].Y, this.LineWidth/2)
	// 	this.dc.Stroke()
	// }
	// this.dc.SetRGBA(0.5, 0.5, 0, 1)
	// this.dc.SetLineWidth(this.LineWidth / 3)
	for index := range myPolygon {
		if index < len(myPolygon)-1 {
			//this.dc.DrawPoint(myPolygon[index].X, myPolygon[index].Y, this.LineWidth/2)
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
	//	this.dc.SetRGBA(0+rand.Float64()*0.3, rand.Float64()*0.3, rand.Float64()*0.5, 1)
	this.dc.SetRGBA(0, 0, 0, 1)
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

func (this *Paper) Draw(name string, binIndex string) {
	if this.polyNum <= 0 {
		//fmt.Println("this.polyNum <= 0 跳过绘制")
		return
	}
	if this.PaperSavingPath == "" {
		this.PaperSavingPath = "./output/" + name + binIndex + ".png"
		this.dc.SavePNG(this.PaperSavingPath)
	} else {
		if strings.HasSuffix(this.PaperSavingPath, "/") {
			this.dc.SavePNG(this.PaperSavingPath + "best" + binIndex + ".png")
			this.PaperSavingPath = this.PaperSavingPath + name + binIndex + ".png"
			this.dc.SavePNG(this.PaperSavingPath)
		} else {
			this.dc.SavePNG(this.PaperSavingPath + "/" + "best" + binIndex + ".png")
			this.PaperSavingPath = this.PaperSavingPath + "/" + name + binIndex + ".png"
			this.dc.SavePNG(this.PaperSavingPath)
		}
	}
}

func (this *Paper) DrawPartPreview(name string, binIndex string) {
	if this.polyNum <= 0 {
		//fmt.Println("this.polyNum <= 0 跳过绘制")
		return
	}
	if this.PaperSavingPath == "" {
		this.dc.SavePNG("./output/" + name + binIndex + ".png")
	} else {
		if strings.HasSuffix(this.PaperSavingPath, "/") {
			this.dc.SavePNG(this.PaperSavingPath + name + binIndex + ".png")
		} else {
			this.dc.SavePNG(this.PaperSavingPath + "/" + name + binIndex + ".png")
		}
	}
}
