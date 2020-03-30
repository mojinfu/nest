package nest

import "math"

type PolygonStruct struct {
	RootPoly  *PolyNode
	id        int //全局唯一id
	typeID    int //同一类型的  NFP 无需重复生成//但要注意这个字段可能没被初始化的  没被初始化 则用id代替
	width     float64
	height    float64
	rotation  int //当前旋转的角度
	Name      string
	isWart    bool
	AngleList []int64
}

func NewArc() {

}
func NewCircle(x, y float64, r float64) []*Point {
	poly := []*Point{}
	var num = int(math.Ceil((2 * math.Pi) / math.Acos(1-(2/r))))
	if num < 3 {
		num = 3
	}
	for i := 0; i < num; i++ {
		theta := float64(i) * ((2 * math.Pi) / float64(num))
		point := &Point{
			X: r*math.Cos(theta) + x,
			Y: r*math.Sin(theta) + y,
		}
		poly = append(poly, point)
	}
	return poly
}
func (this *SVG) SetBinPoly(part *PolygonStruct) {
	this.typeIDIndex++
	part.setTypeID(this.typeIDIndex)

	newPart := NewPoly(part.RootPoly.OriginPolygon)
	newPart.SetAngelList(part.AngleList)
	newPart.setTypeID(this.typeIDIndex)
	newPart.SetName(part.Name)
	newPart.isWart = part.isWart
	this.bins = append(this.bins, newPart)
	if len(this.bins) != 1 {
		panic("")
	}
}
func (this *SVG) AddPoly(part *PolygonStruct, num int) {
	this.typeIDIndex++
	part.setTypeID(this.typeIDIndex)
	for i := 0; i < num; i++ {
		newPart := NewPoly(part.RootPoly.OriginPolygon)
		newPart.SetAngelList(part.AngleList)
		newPart.setTypeID(this.typeIDIndex)
		newPart.SetName(part.Name)
		newPart.isWart = part.isWart
		this.parts = append(this.parts, newPart)
	}
}
func (this *SVG) addPoly(part *PolygonStruct) {
	this.parts = append(this.parts, part)
}

func NewWartPoly(points []*Point) *PolygonStruct {
	return &PolygonStruct{
		RootPoly: &PolyNode{
			OriginPolygon: points,
		},
		isWart: true,
		AngleList: []int64{
			0,
		},
	}
}
func NewPolyWithName(points []*Point, name string) *PolygonStruct {
	a := NewPoly(points)
	a.SetName(name)
	return a
}
func NewPoly(points []*Point) *PolygonStruct {
	poly := &PolygonStruct{
		RootPoly: newNodePoly(points),
	}
	return poly
}
func newNodePoly(points []*Point) *PolyNode {
	poly := &PolyNode{}
	for index := range points {
		poly.OriginPolygon = append(poly.OriginPolygon, &Point{
			X: points[index].X,
			Y: points[index].Y,
		})
	}
	return poly
}
func (this *PolygonStruct) setTypeID(id int) {
	if this.typeID <= 0 {
		this.typeID = id
	}
}
func (this *PolygonStruct) getTypeID() int {
	return this.typeID
}
func (this *PolygonStruct) IsWart() bool {
	return this.isWart
}
func (this *PolygonStruct) SetName(name string) {
	this.Name = name
}
func (this *PolygonStruct) GetName() string {
	return this.Name
}
func (this *PolygonStruct) SetAngelList(angleList []int64) {
	if this.isWart {
		return
	}
	this.AngleList = angleList
}
func (this *PolygonStruct) AddChildPoly(points []*Point) string {
	//检查child是不是都在root poly里面
	this.RootPoly.children = append(this.RootPoly.children, newNodePoly(points))
	return this.Name
}