package nest

import (
	"errors"
	"fmt"
	"log"
	"math"
	"sort"
	"sync"
	"time"

	. "github.com/mojinfu/point"
	"github.com/rlds/rlog"
)

type PairKeyStruct struct {
	A         int
	B         int
	Arotation int
	Brotation int
	inside    bool
}

func (this *PairKeyStruct) ToString() string {
	return fmt.Sprintf(`{"A":%d,"B":%d,"inside":%v,"Arotation":%d,"Brotation":%d}`, this.A, this.B, this.inside, this.Arotation, this.Brotation)
}

type pairStruct struct {
	A     *PolygonStruct
	B     *PolygonStruct
	key   *PairKeyStruct
	value [][]*Point
}

var TOL = math.Pow(10, -9)

func _almostEqual3(a, b, tolerance float64) bool {
	return math.Abs(a-b) < tolerance
}
func _almostEqual2(a, b float64) bool {

	return math.Abs(a-b) < TOL
}

func _onSegment(A, B *Point, p *Point) bool {
	// vertical line
	if _almostEqual2(A.X, B.X) && _almostEqual2(p.X, A.X) {
		if !_almostEqual2(p.Y, B.Y) && !_almostEqual2(p.Y, A.Y) && p.Y < math.Max(B.Y, A.Y) && p.Y > math.Min(B.Y, A.Y) {
			return true
		} else {
			return false
		}
	}

	// horizontal line
	if _almostEqual2(A.Y, B.Y) && _almostEqual2(p.Y, A.Y) {
		if !_almostEqual2(p.X, B.X) && !_almostEqual2(p.X, A.X) && p.X < math.Max(B.X, A.X) && p.X > math.Min(B.X, A.X) {
			return true
		} else {
			return false
		}
	}

	//range check
	if (p.X < A.X && p.X < B.X) || (p.X > A.X && p.X > B.X) || (p.Y < A.Y && p.Y < B.Y) || (p.Y > A.Y && p.Y > B.Y) {
		return false
	}

	// exclude end points
	if (_almostEqual2(p.X, A.X) && _almostEqual2(p.Y, A.Y)) || (_almostEqual2(p.X, B.X) && _almostEqual2(p.Y, B.Y)) {
		return false
	}

	var cross = (p.Y-A.Y)*(B.X-A.X) - (p.X-A.X)*(B.Y-A.Y)

	if math.Abs(cross) > TOL {
		return false
	}

	var dot = (p.X-A.X)*(B.X-A.X) + (p.Y-A.Y)*(B.Y-A.Y)

	if dot < 0 || _almostEqual2(dot, 0) {
		return false
	}

	var len2 = (B.X-A.X)*(B.X-A.X) + (B.Y-A.Y)*(B.Y-A.Y)

	if dot > len2 || _almostEqual2(dot, len2) {
		return false
	}

	return true
}

func pointInPolygon(myPoint *Point, myPolygon Polygon) int64 {
	if len(myPolygon) < 3 {
		return 0
	}

	var inside int64 = -1
	// var offsetx = Polygon.offsetx || 0
	// var offsety = Polygon.offsety || 0
	var offsetx float64 = 0
	var offsety float64 = 0
	i := 0
	j := len(myPolygon) - 1
	for i < len(myPolygon) {
		var xi = myPolygon[i].X + offsetx
		var yi = myPolygon[i].Y + offsety
		var xj = myPolygon[j].X + offsetx
		var yj = myPolygon[j].Y + offsety

		if _almostEqual2(xi, myPoint.X) && _almostEqual2(yi, myPoint.Y) {
			return 0 // no result
		}

		if (_onSegment(&Point{X: xi, Y: yi}, &Point{X: xj, Y: yj}, myPoint)) {
			return 0 // exactly on the segment
		}

		if _almostEqual2(xi, xj) && _almostEqual2(yi, yj) { // ignore very small lines
			j = i
			i++
			continue
		}

		var intersect = ((yi > myPoint.Y) != (yj > myPoint.Y)) && (myPoint.X < (xj-xi)*(myPoint.Y-yi)/(yj-yi)+xi)
		if intersect {
			inside = inside * -1
		}
		j = i
		i++
	}

	return inside
}
func CaluMinXMinY(myPairPolygon Polygon) (float64, float64) {
	if len(myPairPolygon) == 0 {
		return 0, 0
	}
	x := myPairPolygon[0].X
	y := myPairPolygon[0].Y
	for index := range myPairPolygon {
		if x > myPairPolygon[index].X {
			x = myPairPolygon[index].X
		}
		if y > myPairPolygon[index].Y {
			y = myPairPolygon[index].Y
		}
	}
	return x, y
}
func CaluMaxXMaxY(myPairPolygon Polygon) (float64, float64) {
	if len(myPairPolygon) == 0 {
		return 0, 0
	}
	x := myPairPolygon[0].X
	y := myPairPolygon[0].Y
	for index := range myPairPolygon {
		if x < myPairPolygon[index].X {
			x = myPairPolygon[index].X
		}
		if y < myPairPolygon[index].Y {
			y = myPairPolygon[index].Y
		}
	}
	return x, y
}
func TranslatePolygon(myPairPolygon Polygon, x, y float64) Polygon {
	newPolygon := Polygon{}
	for index := range myPairPolygon {
		newPolygon = append(newPolygon, &Point{
			X: myPairPolygon[index].X + x,
			Y: myPairPolygon[index].Y + y,
		})
	}
	return newPolygon
}

func (this *PolygonStruct) TranslateToEndPolygon(x, y float64) {
	for index := range this.RootPoly.EndPolygon {
		this.RootPoly.EndPolygon[index].X = this.RootPoly.EndPolygon[index].X + x
		this.RootPoly.EndPolygon[index].Y = this.RootPoly.EndPolygon[index].Y + y
	}
}

func (this *PolyNode) RotateToEndPolygon(degrees int) {
	if len(this.EndPolygon) != 0 {
		//	panic("rotateToEndPolygon"
		this.EndPolygon = []*Point{}
	} else {

	}
	angle := float64(degrees) * math.Pi / 180
	// for i := 0; i < len(this.polygon); i++ {
	// 	var x = this.polygon[i].X
	// 	var y = this.polygon[i].Y
	// 	var x1 = x*math.Cos(angle) - y*math.Sin(angle)
	// 	var y1 = x*math.Sin(angle) + y*math.Cos(angle)
	// 	this.RootPoly.EndPolygon = append(this.RootPoly.EndPolygon, &Point{X: x1, Y: y1})
	// }
	//画 切换
	for i := 0; i < len(this.CleanedPolygon); i++ {
		var x = this.CleanedPolygon[i].X
		var y = this.CleanedPolygon[i].Y
		var x1 = x*math.Cos(angle) - y*math.Sin(angle)
		var y1 = x*math.Sin(angle) + y*math.Cos(angle)
		this.EndPolygon = append(this.EndPolygon, &Point{X: x1, Y: y1})
	}

	if len(this.children) > 0 {
		//	panic("不支持child")//saya
		// rotated.children = []*PolygonStruct{}
		for j := 0; j < len(this.children); j++ {
			this.children[j].RotateToEndPolygon(degrees)
		}
	}
}
func (this *PolyNode) rotatePolygon(degrees int) {
	angle := float64(degrees) * math.Pi / 180
	for i := 0; i < len(this.polygonBeforeRotation); i++ {
		var x = this.polygonBeforeRotation[i].X
		var y = this.polygonBeforeRotation[i].Y
		var x1 = x*math.Cos(angle) - y*math.Sin(angle)
		var y1 = x*math.Sin(angle) + y*math.Cos(angle)
		this.polygonAfterRotaion[i].X = x1
		this.polygonAfterRotaion[i].Y = y1
	}
	for index := range this.children {
		this.children[index].rotatePolygon(degrees)
	}
}
func (myPairPolygon *PolygonStruct) rotatePolygon(degrees int) {
	if myPairPolygon.rotation == degrees {
		return
	}
	myPairPolygon.rotation = degrees
	myPairPolygon.RootPoly.rotatePolygon(degrees)
}

type SVG struct {
	typeIDIndex      int
	idIndex          int
	parts            []*PolygonStruct
	bins             []*PolygonStruct
	warts            [][]*PolygonStruct
	debugnfp         [][]*Point
	debugnfpList     [][]*Point
	debugnfpList2    [][]*Point
	debugnfpHitPoint int
	working          bool
	nfpCache         map[string][][]*Point
	cleanedRootPoly  map[int][]*Point

	Best           *placementsStruct
	BestRecord     []*placementsStruct
	config         *ConfigStruct
	progress       int64 //进度条
	isNeedStopLoop bool
	// LoopMaxNum        int
	// RunTimeOut        int
	lastDrawedFitness float64
	GA                *GAStruct
	//tree              []*PolygonStruct
	//bin               *BinPolygonStruct
	maxX []int64
	maxY []int64
	//BestCSV           string
	PaperList      map[int][]*Paper
	DrawedPaperNum int
	runID          string
	runStartAt     time.Time
	runEndAt       *time.Time
	LengthWeight   float64
	WidthWeight    float64
}

func (this *SVG) SetRunID(runID string) {
	this.runID = runID
}
func (this *SVG) RunID() string {
	return this.runID
}
func NewSVG(config *ConfigStruct) *SVG {
	this := &SVG{
		debugnfpHitPoint: -1,
	}
	if config == nil {
		this = &SVG{
			config: PublicConfig,
			//LoopMaxNum: loopMaxNum,
		}
	} else {
		this = &SVG{
			config: config,
			//LoopMaxNum: loopMaxNum,
		}
	}
	if config.CanNotPutLoopMaxNum <= 0 {
		config.CanNotPutLoopMaxNum = 100
	}
	if config.WidthWeight <= 0 {
		config.WidthWeight = 1
	}
	if config.LengthWeight <= 0 {
		this.LengthWeight = 0
		this.WidthWeight = 1
	} else {
		this.LengthWeight = config.LengthWeight / (config.WidthWeight + config.LengthWeight)
		this.WidthWeight = config.WidthWeight / (config.WidthWeight + config.LengthWeight)
	}
	// log.Println("长权重：", this.LengthWeight)
	// log.Println("宽权重：", this.WidthWeight)
	this.DrawedPaperNum = 0
	this.PaperList = make(map[int][]*Paper, 10)
	if this.config.RunTimeOut <= 0 {
		this.config.RunTimeOut = 20
	}
	if this.config.LoopMaxNum <= 0 {
		this.config.LoopMaxNum = 20
	}
	return this
}

// func (this *SVG) SetSVGConfig(config *ConfigStruct) {
// 	if config == nil {
// 		this.config = PublicConfig
// 	}
// 	this.config = config
// }
func (this *SVG) GetPolByID(id int) *PolygonStruct {
	for index := range this.parts {
		if this.parts[index].id == id {
			return this.parts[index]
		}
	}
	for index := range this.bins {
		if this.bins[index].id == id {
			return this.bins[index]
		}
	}
	for binIndex := range this.warts {
		for index := range this.warts[binIndex] {
			if this.warts[binIndex][index].id == id {
				return this.warts[binIndex][index]
			}
		}
	}
	return nil
}
func (this *SVG) Draw(i int) {
	if this.Best == nil {
		return
	}
	//生成 csv
	name := fmt.Sprintf("%dline", i)
	if this.lastDrawedFitness == this.Best.fitness {
		return
	}
	////	this.BestCSV = "下料批次,零件号,面料号,零件外轮廓线坐标\n"
	//log.Printf()
	rlog.V(1).Info("%f-->%f\n", this.lastDrawedFitness, this.Best.fitness)
	myPaper := this.NewPaper()
	for index := range this.bins {
		myPaper[index].AddBinPolygon(this.bins[index].RootPoly.OriginPolygon, false)
	}

	for i := range this.Best.Placements {
		for j := range this.Best.Placements[i] {
			this.logIfDebug("id:", this.Best.Placements[i][j].id)
			this.logIfDebug("x:", this.Best.Placements[i][j].x)
			this.logIfDebug("y:", this.Best.Placements[i][j].y)
			this.logIfDebug("rotation:", this.Best.Placements[i][j].rotation)
			//this.tree[this.Best.Placements[i][j].id]
			myPol := this.GetPolByID(this.Best.Placements[i][j].id)
			if myPol.isWart && (this.Best.Placements[i][j].rotation != 0 || this.Best.Placements[i][j].x != 0) {
				panic("")
			}
			myPol.RootPoly.RotateToEndPolygon(this.Best.Placements[i][j].rotation)
			myPol.TranslateToEndPolygon(this.Best.Placements[i][j].x, this.Best.Placements[i][j].y)
			this.logIfDebug("translate x:", this.Best.Placements[i][j].x)
			this.logIfDebug("translate y:", this.Best.Placements[i][j].y)
			myPaper[i].isEmpty = false
			if myPol.isWart {
				myPaper[i].AddPolygon(myPol.RootPoly.OriginPolygon, myPol.isWart)
			} else {
				myPaper[i].AddPolygon(myPol.RootPoly.EndPolygon, myPol.isWart)
			}
			if !myPol.isWart {
				//this.BestCSV = this.BestCSV + fmt.Sprintf("%s,%s,%s,\"%s\"\n", myPol.XiaLiaoPiCi, myPol.Name, myPol.MianLiaoHao, getSvcPointListFmt(myPol.RootPoly.EndPolygon))
			}
		}
	}
	this.PaperList[i] = myPaper
	this.DrawedPaperNum++
	if this.config.IfDraw {
		for index := range myPaper {
			myPaper[index].Draw(name, fmt.Sprintf("_%d", index))
		}
	}
	this.lastDrawedFitness = this.Best.fitness
}
func getSvcPointListFmt(pointList []*Point) string {
	all := "["
	for index := range pointList {
		all = all + fmt.Sprintf("[%v, %v]", pointList[index].X, pointList[index].Y)
		if index == len(pointList)-1 {
			all = all + "]"
		} else {
			all = all + ", "
		}
	}
	return all
}
func (this *GAStruct) randomWeightedIndividual(exclude *populationStruct) *populationStruct {
	pop := this.population[:]
	if exclude != nil {
		temp := []*populationStruct{}
		for index := range pop {
			if pop[index] == exclude {
			} else {
				temp = append(temp, pop[index])
			}

		}
		pop = temp
	}
	myRand := this.randomSeed.Float64()
	var lower float64 = 0
	var weight float64 = float64(1) / float64(len(pop))
	var upper float64 = weight

	for i := 0; i < len(pop); i++ {
		// if the random number falls between lower and upper bounds, select this individual
		if myRand > lower && myRand < upper {
			return pop[i]
		}
		lower = upper
		upper += 2 * weight * (float64(len(pop)-i) / float64(len(pop)))
	}
	return pop[0]
}
func (this *GAStruct) mate(male, female *populationStruct) []*populationStruct {

	myRand := this.randomSeed.Intn(10000)

	var cutpoint = int(math.Round(math.Min(math.Max(float64(myRand)/10000, 0.1), 0.9) * float64(len(male.placement)-1)))

	this.logIfDebug("cutpoint:", cutpoint)
	this.logIfDebug("len(male.placement):", len(male.placement))
	this.logIfDebug("len(male.rotation):", len(male.rotation))
	this.logIfDebug("len(female.placement):", len(female.placement))
	this.logIfDebug("len(female.rotation):", len(female.rotation))
	//log.Println("cutpoint:", cutpoint)
	var gene1 = male.placement[0:cutpoint]
	var rot1 = male.rotation[0:cutpoint]

	var gene2 = female.placement[0:cutpoint]
	var rot2 = female.rotation[0:cutpoint]
	contains := func(gene []*PolygonStruct, id int) bool {
		for i := 0; i < len(gene); i++ {
			if gene[i].id == id {
				return true
			}
		}
		return false
	}
	for i := 0; i < len(female.placement); i++ {
		if !contains(gene1, female.placement[i].id) {
			gene1 = append(gene1, female.placement[i])
			rot1 = append(rot1, female.rotation[i])
		}
	}

	for i := 0; i < len(male.placement); i++ {
		if !contains(gene2, male.placement[i].id) {
			gene2 = append(gene2, male.placement[i])
			rot2 = append(rot2, male.rotation[i])
		}
	}
	return []*populationStruct{&populationStruct{placement: gene1, rotation: rot1}, &populationStruct{placement: gene2, rotation: rot2}}
}
func (this *GAStruct) generation() {

	sort.Sort(populationSlice(this.population))
	// // Individuals with higher fitness are more likely to be selected for mating
	// this.population.sort(function (a, b) {
	// 	return a.fitness - b.fitness;
	// });

	// fittest individual is preserved in the new generation (elitism)
	var newpopulation []*populationStruct = []*populationStruct{this.population[0]}

	for len(newpopulation) < len(this.population) {
		var male = this.randomWeightedIndividual(nil)
		var female = this.randomWeightedIndividual(male)

		// each mating produces two children
		var children = this.mate(male, female)

		// slightly mutate children
		newpopulation = append(newpopulation, this.mutate(children[0]))

		if len(newpopulation) < len(this.population) {
			newpopulation = append(newpopulation, this.mutate(children[1]))
		}
	}
	this.population = newpopulation
}

func (this *SVG) getPartsWithInfo(paths []*PolygonStruct) []*PolygonStruct {

	//var polygons = []*PolygonStruct{}
	var numChildren = len(paths)
	for i := 0; i < numChildren; i++ {
		// myStructPolygon := paths[i]
		// for index := range myStructPolygon.OriginPolygon {
		// 	myStructPolygon.polygon = append(myStructPolygon.polygon, myStructPolygon.OriginPolygon[index])
		// }
		// //var poly = SvgParser.polygonify(paths[i]);
		// //poly = this.cleanPolygon(poly);//saya
		// // todo: warn user if poly could not be processed and is excluded from the nest
		if len(paths[i].RootPoly.polygonBeforeRotation) > 2 && math.Abs(PolygonArea(paths[i].RootPoly.polygonBeforeRotation)) > this.config.CurveTolerance*this.config.CurveTolerance {
			//	paths[i].source = i
			//polygons = append(polygons, paths[i])
		} else {
			log.Print("PolygonArea(paths[i].RootPoly.polygonBeforeRotation)):", PolygonArea(paths[i].RootPoly.polygonBeforeRotation))
		}
		// paths[i].source = i
	}
	// turn the list into a tree
	return paths
}

// func (this *SVG) getParts(paths []Polygon) []*PolygonStruct {

// 	var polygons = []*PolygonStruct{}

// 	var numChildren = len(paths)
// 	for i := 0; i < numChildren; i++ {
// 		myStructPolygon := &PolygonStruct{
// 			RootPoly: &PolyNode{
// 				polygon: paths[i],
// 			},
// 		}
// 		//var poly = SvgParser.polygonify(paths[i]);
// 		//poly = this.cleanPolygon(poly);//saya

// 		// todo: warn user if poly could not be processed and is excluded from the nest
// 		if len(myStructPolygon.RootPoly.polygon) > 2 && math.Abs(PolygonArea(myStructPolygon.RootPoly.polygon)) > this.config.CurveTolerance*this.config.CurveTolerance {
// 			myStructPolygon.source = i
// 			polygons = append(polygons, myStructPolygon)
// 		} else {
// 			log.Print("PolygonArea(myStructPolygon.RootPoly.polygon)):", PolygonArea(myStructPolygon.RootPoly.polygon))
// 		}
// 	}

// 	// turn the list into a tree
// 	toTree(polygons, 0)

// 	return polygons
// }

// func toTree(list []*PolygonStruct, idstart int) int {
// 	var parents = []*PolygonStruct{}

// 	// assign a unique id to each leaf
// 	var id = idstart

// 	for i := 0; i < len(list); i++ {
// 		var p = list[i]

// 		var ischild = false
// 		for j := 0; j < len(list); j++ {
// 			if j == i {
// 				continue
// 			}
// 			// if pointInPolygon(p.RootPoly.polygon[0], list[j].RootPoly.polygon) == 1 {
// 			// 	if len(list[j].RootPoly.children) == 0 {
// 			// 		list[j].RootPoly.children = []*PolyNode{}
// 			// 	}
// 			// 	list[j].RootPoly.children = append(list[j].children, p)
// 			// 	p.parent = list[j]
// 			// 	ischild = true
// 			// 	break
// 			// }
// 		}

// 		if !ischild {
// 			parents = append(parents, p)
// 		}
// 	}

// 	for i := 0; i < len(list); i++ {
// 		isExist := false
// 		for indexParent := range parents {
// 			if parents[indexParent] == list[i] {
// 				isExist = true
// 			}
// 		}
// 		if !isExist {
// 			newList := []*PolygonStruct{}
// 			for index := range list {
// 				if index == i {

// 				} else {
// 					newList = append(newList, list[i])
// 				}
// 			}
// 			list = newList
// 			i--
// 		}

// 	}

// 	for i := 0; i < len(parents); i++ {
// 		parents[i].id = id
// 		id++
// 	}

// 	for i := 0; i < len(parents); i++ {
// 		//if parents[i].children {
// 		if true {
// 			id = toTree(parents[i].children, id)
// 		}
// 	}

// 	return id
// }
func (this *SVG) polygonOffset(myPolygon Polygon, offset float64) [][]*Point {

	if offset == 0 || _almostEqual2(offset, 0) {
		return []Polygon{myPolygon}
	}

	var p = ScaleUpPath(myPolygon, this.config.ClipperScale)

	var miterLimit float64 = 2 //saya 原来默认值2
	var co = NewClipperOffset(miterLimit, this.config.CurveTolerance*float64(this.config.ClipperScale))
	co.AddPath(p, jtRound, etClosedPolygon)

	newpaths := co.ExecutePath(offset * float64(this.config.ClipperScale))

	var result = [][]*Point{}
	for i := 0; i < len(newpaths); i++ {
		result = append(result, ScaleDownPath(newpaths[i], this.config.ClipperScale))
	}

	return result
}
func (this *PolyNode) offsetTree(offset float64, offsetFunction func(myPolygon Polygon, offset float64) [][]*Point) {
	offsetpaths := offsetFunction(this.polygonBeforeRotation, offset)
	if len(offsetpaths) == 1 {
		// replace array items in place
		this.polygonBeforeRotation = offsetpaths[0]
		//Array.prototype.splice.apply(t[i], [0, len(t[i]),offsetpaths[0]]);
	} else {
		//?
		this.polygonBeforeRotation = offsetpaths[0]
		log.Println("offset err")
	}
	for index := range this.children {
		this.children[index].offsetTree(-offset, offsetFunction)
	}
}

func (this *PolyNode) copyBeforeToAfter() {
	newAfter := []*Point{}
	for index := range this.polygonBeforeRotation {
		newAfter = append(newAfter, &Point{
			X: this.polygonBeforeRotation[index].X,
			Y: this.polygonBeforeRotation[index].Y,
		})
	}
	this.polygonAfterRotaion = newAfter
	for index := range this.children {
		this.children[index].copyBeforeToAfter()
	}
}

func (this *SVG) offsetTree(t []*PolygonStruct, offset float64, offsetFunction func(myPolygon Polygon, offset float64) [][]*Point) {
	for i := 0; i < len(t); i++ {
		t[i].RootPoly.offsetTree(offset, offsetFunction)
	}
}

func (this *SVG) GetRunStartAt() time.Time {
	return this.runStartAt
}
func (this *SVG) GetRunEndAt() *time.Time {
	return this.runEndAt
}
func (this *SVG) lazyCleanPolygon(p *PolygonStruct) {
	cleanedRootPoly, isOK := this.cleanedRootPoly[p.TypeID()]
	if isOK {
		p.RootPoly.polygonBeforeRotation = CopyPointList(cleanedRootPoly)
		p.RootPoly.CleanedPolygon = CopyPointList(cleanedRootPoly)
	} else {
		p.CleanPolygon(-1)
		this.cleanedRootPoly[p.TypeID()] = CopyPointList(p.RootPoly.CleanedPolygon)
	}
}
func (this *SVG) Start() error {
	defer func() {
		// if r := recover(); r != nil {
		// 	fmt.Printf("SVG Start 捕获到的错误：%s\n", r)
		// }
		now := time.Now()
		this.runEndAt = &now
	}() //saya 错误处理
	if len(this.bins) == 0 {
		return WithoutBinErr
	}
	if len(this.parts) == 0 {
		return WithoutPartsErr
	}
	this.nfpCache = make(map[string][][]*Point, 10000)
	this.cleanedRootPoly = make(map[int][]*Point, 10000)
	this.runStartAt = time.Now()
	/*----*/
	for index := range this.parts {
		this.lazyCleanPolygon(this.parts[index])
	}
	for index := range this.bins {
		this.lazyCleanPolygon(this.bins[index])
	}
	//log.Println("OriginPolygon len:", len(this.bins[0].RootPoly.OriginPolygon))
	//log.Println("len:", len(this.bins[0].RootPoly.polygonBeforeRotation))
	if len(this.bins[0].RootPoly.polygonBeforeRotation) == 0 {
		return errors.New("定义的bin有些问题")
	}
	for i := range this.warts {
		for j := range this.warts[i] {
			this.lazyCleanPolygon(this.warts[i][j])
		}
	}
	tree := this.parts
	if false {
		for index := range tree {
			paper := NewPaper(1000, 1000, this.config.PaperSavePath)
			paper.AddPolygon(tree[index].RootPoly.polygonBeforeRotation, true)
			paper.AddPolygon(tree[index].RootPoly.OriginPolygon, false)
			paper.Draw("polygonBeforeRotationA:", fmt.Sprintf("%d", tree[index].TypeID()))
		}
	}
	//return nil
	/*----*/

	this.offsetTree(tree, 0.5*this.config.PartPartSpacing, this.polygonOffset) //先0空隙 saya

	if false {
		for index := range tree {
			paper := NewPaper(1000, 1000, this.config.PaperSavePath)
			paper.AddPolygon(tree[index].RootPoly.polygonBeforeRotation, true)
			paper.AddPolygon(tree[index].RootPoly.OriginPolygon, false)
			paper.Draw("polygonBeforeRotationB:", fmt.Sprintf("%d", tree[index].TypeID()))
		}
	}
	// remove duplicate endpoints, ensure counterclockwise winding direction
	for i := 0; i < len(tree); i++ {
		var start = tree[i].RootPoly.polygonBeforeRotation[0]
		var end = tree[i].RootPoly.polygonBeforeRotation[len(tree[i].RootPoly.polygonBeforeRotation)-1]
		if start == end || (_almostEqual2(start.X, end.X) && _almostEqual2(start.Y, end.Y)) {
			tree[i].RootPoly.polygonBeforeRotation = tree[i].RootPoly.polygonBeforeRotation[1:]
		}
		if PolygonArea(tree[i].RootPoly.polygonBeforeRotation) > 0 {
			tree[i].RootPoly.polygonBeforeRotation = PolygonReverse(tree[i].RootPoly.polygonBeforeRotation)
		}

	}
	for index := range tree {
		if len(tree[index].RootPoly.polygonBeforeRotation) == 0 {
			return InvalidPoly
		} else {
			tree[index].RootPoly.copyBeforeToAfter()
		}
	}
	/*----*/
	for i := range this.warts {
		this.offsetTree(this.warts[i], 0.5*this.config.PartPartSpacing, this.polygonOffset) //先0空隙 saya
		for j := range this.warts[i] {
			this.warts[i][j].RootPoly.copyBeforeToAfter()
		}
	}
	/*----*/
	for index := range this.bins {
		if this.config.PartPartSpacing > 0 || this.config.BinPartSpacing > 0 {
			var offsetBin = this.polygonOffset(this.bins[index].RootPoly.polygonBeforeRotation, 0.5*this.config.PartPartSpacing-0.5*this.config.BinPartSpacing) //轮廓 saya
			if len(offsetBin) == 1 {
				this.bins[index].RootPoly.polygonBeforeRotation = offsetBin[0]
			} else {
				panic("")
			}
		}

		var xbinmax = this.bins[index].RootPoly.polygonBeforeRotation[0].X
		var xbinmin = this.bins[index].RootPoly.polygonBeforeRotation[0].X
		var ybinmax = this.bins[index].RootPoly.polygonBeforeRotation[0].Y
		var ybinmin = this.bins[index].RootPoly.polygonBeforeRotation[0].Y

		for i := 1; i < len(this.bins[index].RootPoly.polygonBeforeRotation); i++ {
			if this.bins[index].RootPoly.polygonBeforeRotation[i].X > xbinmax {
				xbinmax = this.bins[index].RootPoly.polygonBeforeRotation[i].X
			} else if this.bins[index].RootPoly.polygonBeforeRotation[i].X < xbinmin {
				xbinmin = this.bins[index].RootPoly.polygonBeforeRotation[i].X
			}
			if this.bins[index].RootPoly.polygonBeforeRotation[i].Y > ybinmax {
				ybinmax = this.bins[index].RootPoly.polygonBeforeRotation[i].Y
			} else if this.bins[index].RootPoly.polygonBeforeRotation[i].Y < ybinmin {
				ybinmin = this.bins[index].RootPoly.polygonBeforeRotation[i].Y
			}
		}
		// //让bin 回到原点  （可能点不在原点）
		// for i := 0; i < len(this.bins[index].RootPoly.polygonBeforeRotation); i++ {
		// 	this.bins[index].RootPoly.polygonBeforeRotation[i].X -= xbinmin
		// 	this.bins[index].RootPoly.polygonBeforeRotation[i].Y -= ybinmin
		// }

		this.bins[index].width = xbinmax - xbinmin
		this.bins[index].height = ybinmax - ybinmin
		this.maxX = append(this.maxX, int64(xbinmax)+int64(0.5*this.config.BinPartSpacing-0.5*this.config.PartPartSpacing))
		this.maxY = append(this.maxY, int64(ybinmax)+int64(0.5*this.config.BinPartSpacing-0.5*this.config.PartPartSpacing))
		//确立坐标轴
		// all paths need to have the same winding direction
		if PolygonArea(this.bins[index].RootPoly.polygonBeforeRotation) > 0 {
			this.bins[index].RootPoly.polygonBeforeRotation = PolygonReverse(this.bins[index].RootPoly.polygonBeforeRotation)
		}
		this.bins[index].RootPoly.copyBeforeToAfter()
	}

	this.working = false
	if this.config == nil {
		this.config = PublicConfig
	}
	//this.tree = tree
	var err error = CanNotPutErr
	canNotPutNum := 0
	for i := 0; i < this.config.LoopMaxNum && time.Now().Before(this.runStartAt.Add(time.Second*time.Duration(this.config.RunTimeOut))) && canNotPutNum < this.config.CanNotPutLoopMaxNum; i++ {
		if this.isNeedStopLoop {
			break
		}
		if !this.working {
			_, errWork := this.launchWorkers()
			if errWork == CanNotPutErr && err == CanNotPutErr {
				canNotPutNum++
			}
			if errWork != nil {
				time.Sleep(time.Millisecond * 500)
			} else {
				this.Draw(i)
				err = nil
			}
		} else {
			time.Sleep(time.Second * 1)
		}
	}
	return err
}
func (this *SVG) GetOrCreateNFP(pairKey *PairKeyStruct, ANotInTree *PolygonStruct, BNotInTree *PolygonStruct) [][]*Point {
	var key = pairKey.ToString()
	value, isOK := this.nfpCache[key]
	if isOK {
		return value
	}
	A := this.GetPolByID(ANotInTree.id)
	B := this.GetPolByID(BNotInTree.id)
	pair := &pairStruct{
		key: pairKey,
		A:   A,
		B:   B,
	}
	log.Println("!!!需要一个未生成的NFP:", pairKey.ToString())
	this.MapNfpPairsList(pair)
	this.nfpCache[pair.key.ToString()] = pair.value
	// if A.isWart && !B.isWart {
	// 	for i := range newPair.value {
	// 		fmt.Println("--")
	// 		for j := range newPair.value[i] {
	// 			fmt.Printf("(%f,%f)", newPair.value[i][j].X, newPair.value[i][j].Y)
	// 		}
	// 	}
	// }
	return pair.value
}
func GenPairKey(A *PolygonStruct, B *PolygonStruct, inside bool, Arotation int, Brotation int) *PairKeyStruct {
	AKey := 0
	BKey := 0
	if false {
		AKey = A.id
	} else {
		AKey = A.typeID
	}
	if false {
		BKey = B.id
	} else {
		BKey = B.typeID
	}
	return &PairKeyStruct{A: AKey, B: BKey, inside: inside, Arotation: Arotation, Brotation: Brotation}
}
func (this *SVG) launchWorkers() (*placementsStruct, error) {
	if this.GA == nil {
		// initiate new GA
		adam := this.parts[:]
		sort.Sort(PolygonStructSlice(adam)) //saya  此处必须要排序
		this.GA = this.NewGeneticAlgorithm(adam, this.bins, this.config)
	}
	var individual *populationStruct = nil
	// // evaluate all members of the population
	for i := 0; i < len(this.GA.population); i++ {
		if this.GA.population[i].fitness == defaultfitness {
			individual = this.GA.population[i]
			break
		}
	}
	//说明之前已经遍历过一次GA 重新生成
	if individual == nil {
		// all individuals have been evaluated, start next generation
		this.GA.generation()
		individual = this.GA.population[1]
	}
	var placelist = individual.placement
	var rotations = individual.rotation

	nfpPairs := this.GenPairs(placelist, rotations)
	//以上  两两组合 生产 nfp
	this.multiThreadingBuildNfp(nfpPairs)
	worker := this.NewPlacementWorker(this.bins, placelist[:], this.warts, rotations, this.config, this.nfpCache)
	this.setNFPToWorker(nfpPairs, worker)
	myplacements, err := worker.newPlacePaths()
	if err != nil {
		return nil, err
	}
	// for i := range myplacements.Placements {
	// 	for j := range myplacements.Placements[i] {
	// 		log.Println("x:", myplacements.Placements[i][j])
	// 	}
	// }
	individual.fitness = myplacements.fitness
	this.WhenWorkerPlacePathsSuccess(myplacements)
	return myplacements, nil
}
func (this *SVG) GenPairs(placelist []*PolygonStruct, rotations []int) []*pairStruct {

	dupMap := map[string]*pairStruct{}
	nfpPairs := []*pairStruct{}
	//两两组合点结构体 包含了npf
	var key *PairKeyStruct
	//两两组合点结构体 不包含npf
	//newCache := map[string][][]*Point{}
	for i := 0; i < len(placelist); i++ {
		for binIndex := 0; binIndex < len(this.warts); binIndex++ {
			for j := 0; j < len(this.warts[binIndex]); j++ {
				key = GenPairKey(
					this.warts[binIndex][j],
					placelist[i],
					false, 0, rotations[i])
				if _, isOK := this.nfpCache[key.ToString()]; !isOK {
					pair := &pairStruct{A: this.warts[binIndex][j], B: placelist[i], key: key}
					dupMap[key.ToString()] = pair
					//nfpPairs = append(nfpPairs)
				}
			}
		}
	}
	for i := 0; i < len(placelist); i++ {
		for j := 0; j < len(this.bins); j++ {
			var part = placelist[i]
			key = GenPairKey(
				this.bins[j],
				part, true, 0, rotations[i])
			if _, isOK := this.nfpCache[key.ToString()]; !isOK {
				//nfpPairs = append(nfpPairs)
				pair := &pairStruct{A: this.bins[j], B: part, key: key}
				dupMap[key.ToString()] = pair
			}
			for j := 0; j < i; j++ {
				var placed = placelist[j]
				key = GenPairKey(placed, part, false, rotations[j], rotations[i])
				if _, isOK := this.nfpCache[key.ToString()]; !isOK {
					pair := &pairStruct{A: placed, B: part, key: key}
					dupMap[key.ToString()] = pair
					//nfpPairs = append(nfpPairs)
				}
			}
		}
	}
	for _, v := range dupMap {
		nfpPairs = append(nfpPairs, v)
	}
	//log.Println("nfpPairs:", len(nfpPairs))
	return nfpPairs
}
func (this *SVG) WhenWorkerPlacePathsSuccess(placement *placementsStruct) {
	if placement == nil {
		return
	}
	// var bestresult = placements[0]
	// for i := 1; i < len(placements); i++ {
	// 	if placements[i].fitness < bestresult.fitness {
	// 		bestresult = placements[i]
	// 	}
	// }
	if this.Best == nil || placement.fitness < this.Best.fitness {
		this.Best = placement
		this.BestRecord = append(this.BestRecord, placement)
		// var placedArea float64
		// var totalArea float64
		// //var numParts = len(placelist)
		// var numPlacedParts = 0
		// for i := 0; i < len(this.Best.Placements); i++ {
		// 	totalArea += math.Abs(PolygonArea(binPolygon.RootPoly.polygonBeforeRotation))
		// 	for j := 0; j < len(this.Best.Placements[i]); j++ {
		// 		placedArea += math.Abs(PolygonArea(tree[this.Best.Placements[i][j].id].RootPoly.polygonBeforeRotation))
		// 		numPlacedParts++
		// 	}
		// }
		//	displayCallback(self.applyPlacement(this.Best.Placements), placedArea/totalArea, numPlacedParts, numParts)//图形化todo
	} else {
		//	displayCallback()//图形化todo
	}
	this.working = false
}
func (this *SVG) setNFPToWorker(nfpPairs []*pairStruct, worker *PlacementWorkerStruct) {
	if nfpPairs != nil {
		for i := 0; i < len(nfpPairs); i++ {
			var Nfp = nfpPairs[i]
			if Nfp != nil {
				// a null nfp means the nfp could not be generated, either because the parts simply don't fit or an error in the nfp algo
				var key = Nfp.key.ToString()
				this.nfpCache[key] = Nfp.value
				if this.config.IfDebug {
					this.logIfDebug("放入 key:", key)
					for index := range Nfp.value {
						log.Println("第", index, "个NFP")
						for j := range Nfp.value[index] {
							log.Printf("(%f,%f)", Nfp.value[index][j].X, Nfp.value[index][j].Y)
						}
						log.Printf("\n")
						//this.logIfDebug(Nfp.value[index])
					}
				}
			}
		}
	}
	worker.nfpCache = this.nfpCache
}
func (this *SVG) multiThreadingBuildNfp(pairList []*pairStruct) {
	if this.config.MutilThread <= 0 {
		this.config.MutilThread = 1
	}
	var wg sync.WaitGroup
	for i := 0; i < this.config.MutilThread; i++ {
		wg.Add(1)
		go this.threadingBuildNfp(i, pairList, &wg)
	}
	wg.Wait()
	for index := range pairList {
		this.nfpCache[pairList[index].key.ToString()] = pairList[index].value
	}
	return
}
func (this *SVG) threadingBuildNfp(threadID int, pairList []*pairStruct, wg *sync.WaitGroup) {
	for i := 0; this.config.MutilThread*i+threadID < len(pairList); i++ {
		this.MapNfpPairsList(pairList[this.config.MutilThread*i+threadID])
	}
	wg.Done()
}
func (this *SVG) MapNfpPairsList(pair *pairStruct) {
	if pair == nil {
		return //saya 可能有问题  pair.length == 0
	}
	var searchEdges = this.config.ExploreConcave
	//var UseHoles = config.UseHoles
	//fmt.Println("A rotation :", pair.key.A, pair.key.Arotation)
	pair.A.rotatePolygon(pair.key.Arotation)
	//fmt.Println("B rotation :", pair.key.B, pair.key.Brotation)
	pair.B.rotatePolygon(pair.key.Brotation)
	nfp := [][]*Point{}
	if pair.key.inside {
		if isRectangle(pair.A.RootPoly.polygonAfterRotaion, 0.001) {
			nfp = noFitPolygonRectangle(pair.A.RootPoly.polygonAfterRotaion, pair.B.RootPoly.polygonAfterRotaion)
		} else {
			//当一个多边形可以放置在另一个多边形内 且不是长方形当情况//saya todo
			this.logIfDebug("inside and not rectangle")
			nfp = this.noFitPolygon(pair.A.RootPoly.polygonAfterRotaion, pair.B.RootPoly.polygonAfterRotaion, true, searchEdges)
		}
		// ensure all interior NFPs have the same winding direction
		if len(nfp) > 0 {
			for i := 0; i < len(nfp); i++ {
				if PolygonArea(nfp[i]) > 0 {
					nfp[i] = PolygonReverse(nfp[i])
				}
			}
		}
		if len(nfp) == 0 {
			// log.Printf("A:%+v\n", pair.A.RootPoly.String())
			// log.Printf("B:%+v\n", pair.B.RootPoly.String())
			// log.Println("nfp gen err:" + pair.key.ToString())
		}
	} else {
		if searchEdges {
			panic("searchEdges")
			nfp = this.noFitPolygon(pair.A.RootPoly.polygonAfterRotaion, pair.B.RootPoly.polygonAfterRotaion, false, searchEdges) //saya todo
		} else {
			nfp = minkowskiDifference(pair.A.RootPoly.polygonAfterRotaion, pair.B.RootPoly.polygonAfterRotaion)
		}
		if len(nfp) == 0 {
			panic("") //saya
			log.Println("NFP Error: ", pair.key)
			this.logIfDebug("NFP Error: ", pair.key)
			return
		}

		for i := 0; i < len(nfp); i++ {
			if !searchEdges || i == 0 { // if searchedges is active, only the first NFP is guaranteed to pass sanity check
				if math.Abs(PolygonArea(nfp[i])) < math.Abs(PolygonArea(pair.A.RootPoly.polygonBeforeRotation)) {
					newNfp := [][]*Point{}
					for index := range nfp {
						if index == i {

						} else {
							newNfp = append(newNfp, nfp[index])
						}
					}
					nfp = newNfp
					//splice
					return
				}
			}
		}
		if len(nfp) == 0 {
			//不应该return
			return
		}

		// for outer NFPs, the first is guaranteed to be the largest. Any subsequent NFPs that lie inside the first are holes
		for i := 0; i < len(nfp); i++ {
			if PolygonArea(nfp[i]) > 0 {
				nfp[i] = PolygonReverse(nfp[i])
			}

			if i > 0 {
				if pointInPolygon(nfp[i][0], nfp[0]) == 1 {
					if PolygonArea(nfp[i]) < 0 {
						nfp[i] = PolygonReverse(nfp[i])
					}
				}
			}
		}
		//UseHoles saya todo
		// // generate nfps for children (holes of parts) if any exist
		// if UseHoles && A.childNodes && A.childNodes.length > 0 {
		// 	var Bbounds = getPolygonBounds(B)

		// 	for i := 0; i < len(A.childNodes); i++ {
		// 		var Abounds = getPolygonBounds(A.childNodes[i])

		// 		// no need to find nfp if B's bounding box is too big
		// 		if Abounds.width > Bbounds.width && Abounds.height > Bbounds.height {

		// 			var cnfp = noFitPolygon(A.childNodes[i], B, true, searchEdges)
		// 			// ensure all interior NFPs have the same winding direction
		// 			if cnfp && cnfp.length > 0 {
		// 				for j := 0; j < len(cnfp); j++ {
		// 					if PolygonArea(cnfp[j]) < 0 {
		// 						cnfp[j] = PolygonReverse(cnfp[j])
		// 					}
		// 					nfp.push(cnfp[j])
		// 				}
		// 			}

		// 		}
		// 	}
		// }
	}

	// fmt.Println("key str:", pair.key.ToString())
	// if pair.A.isWart && !pair.B.isWart && !pair.key.inside {
	// 	for i := range nfp {
	// 		fmt.Println("----")
	// 		for j := range nfp[i] {
	// 			fmt.Printf("(%f,%f)", nfp[i][j].X, nfp[i][j].Y)
	// 		}
	// 		fmt.Printf("\n")
	// 	}
	// 	fmt.Println("A:", A.rotation)
	// 	for i := range pair.A.polygon {
	// 		fmt.Printf("(%f,%f)", pair.A.polygon[i].X, pair.A.polygon[i].Y)
	// 	}
	// 	fmt.Printf("\n")
	// 	fmt.Println("B:", B.rotation)
	// 	for i := range pair.B.polygon {
	// 		fmt.Printf("(%f,%f)", pair.B.polygon[i].X, pair.B.polygon[i].Y)
	// 	}
	// 	fmt.Printf("\n")
	// }
	pair.value = nfp
	return
}
func cnFang(x, y, x1, y1, x2, y2 float64) float64 {
	var cross float64 = (x2-x1)*(x-x1) + (y2-y1)*(y-y1) // |AB| * |AC|*cos(x)
	if cross <= 0 {
		return (x-x1)*(x-x1) + (y-y1)*(y-y1)
	} //积小于等于0，说明 角BAC 是直角或钝角

	var d2 float64 = (x2-x1)*(x2-x1) + (y2-y1)*(y2-y1) // |AB|
	if cross >= d2 {
		return (x-x2)*(x-x2) + (y-y2)*(y-y2)
	} //角ABC是直角或钝角

	//锐角三角形
	var r float64 = cross / d2
	var px float64 = x1 + (x2-x1)*r // C在 AB上的垂足点（px，py）
	var py float64 = y1 + (y2-y1)*r
	return (x-px)*(x-px) + (y-py)*(y-py) //两点间距离公式
}

// func (placement *placementsStruct) GetDistencePolo(a, b []*Point) {
// 	for i := range a {
// 		for j := 0; j < len(b)-2; j++ {
// 			dis := cnFang(a[i].X, a[i].Y, b[j].X, b[j].Y, b[j+1].X, b[j+1].Y)
// 			if placement.MinDistencePoint > dis {
// 				placement.MinDistencePoint = dis
// 			}
// 		}

// 	}

// }
// func (placement *placementsStruct) CanSubmit() {
// 	placement.MinDistencePoint = 100000000
// 	for i := 0; i < len(placement.placements)-1; i++ {
// 		for j := 0; j < len(placement.placements)-1; j++ {
// 			if i == j {
// 				continue
// 			} else {
// 				placement.GetDistencePolo(placement.placements[i], placement.placements[j])
// 			}
// 		}
// 	}

// 	for i := 0; i < len(placement.placements)-1; i++ {
// 		for j := range placement.placements[i] {
// 			if placement.placements[i][j][0] > 20000 || placement.placements[i][j][0] < 0 {
// 				log.Println(placement.placements[i][j][0])
// 				panic("out range")
// 			}
// 			if placement.placements[i][j][1] > 1600 || placement.placements[i][j][1] < 0 {
// 				log.Println(placement.placements[i][j][1])
// 				panic("out range")
// 			}
// 		}
// 	}

// 	if MinDistencePoint < 25 {
// 		panic("too close")
// 	}
// }
type Bin struct {
	binIndex int
	svg      *SVG
}

func (this *Bin) AddWart(part *PolygonStruct) *PolygonStruct {
	log.Println("AddWart")
	this.svg.typeIDIndex++
	part.setTypeID(this.svg.typeIDIndex)
	newPart := NewPoly(part.RootPoly.OriginPolygon)
	newPart.SetAngelList([]int64{0})
	newPart.setTypeID(this.svg.typeIDIndex)
	newPart.SetName(part.Name)
	newPart.isWart = true
	this.svg.addWart(newPart, this.binIndex)
	return newPart
}
func (this *SVG) addWart(newPart *PolygonStruct, binIndex int) {
	this.idIndex++
	newPart.id = this.idIndex
	this.warts[binIndex] = append(this.warts[binIndex], newPart)
}
func (this *SVG) AddBinPoly(part *PolygonStruct) *Bin {
	this.typeIDIndex++
	part.setTypeID(this.typeIDIndex)
	newPart := NewPoly(part.RootPoly.OriginPolygon)
	newPart.SetAngelList([]int64{0})
	newPart.setTypeID(this.typeIDIndex)
	newPart.SetName(part.Name)
	newPart.isWart = false
	return this.addBinPoly(newPart)
}

func (this *SVG) AddBinPolyWithNum(part *PolygonStruct, num int) int {
	this.typeIDIndex++
	part.setTypeID(this.typeIDIndex)
	for i := 0; i < num; i++ {
		newPart := NewPoly(part.RootPoly.OriginPolygon)
		newPart.SetAngelList([]int64{0})
		newPart.setTypeID(this.typeIDIndex)
		newPart.SetName(part.Name)
		newPart.isWart = false
		this.addBinPoly(newPart)
	}
	return this.typeIDIndex
}

func (this *SVG) addBinPoly(newPart *PolygonStruct) *Bin {
	this.idIndex++
	newPart.id = this.idIndex
	this.bins = append(this.bins, newPart)
	this.warts = append(this.warts, []*PolygonStruct{})
	return &Bin{
		binIndex: len(this.bins) - 1,
		svg:      this,
	}
}
func (this *SVG) addPart(newPart *PolygonStruct) {
	this.idIndex++
	newPart.id = this.idIndex
	this.parts = append(this.parts, newPart)
}
func (this *SVG) AddPart(part *PolygonStruct, num int) int {
	this.typeIDIndex++
	part.setTypeID(this.typeIDIndex)
	for i := 0; i < num; i++ {
		newPart := NewPoly(part.RootPoly.OriginPolygon)
		newPart.SetAngelList(part.AngleList)
		newPart.setTypeID(this.typeIDIndex)
		newPart.SetName(part.Name)
		newPart.isWart = false
		this.addPart(newPart)
	}
	return this.typeIDIndex
}

// func (this *SVG) GetPolyById(id int) *PolygonStruct {
// 	for index := range this.parts {
// 		if this.parts[index].id == id {
// 			return this.parts[index]
// 		}
// 	}
// 	return nil
// }
