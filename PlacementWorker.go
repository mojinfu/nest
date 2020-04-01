package nest

import (
	"log"
	"math"

	. "github.com/mojinfu/nest/point"
)

func (this *PlacementWorkerStruct) logIfDebug(v ...interface{}) {
	if this.IfDebug {
		log.Println(v...)
	}
}

type PlacementWorkerStruct struct {
	binPolygon *BinPolygonStruct
	paths      []*PolygonStruct
	ids        []int
	Rotations  []int
	config     *ConfigStruct
	nfpCache   map[string][]Polygon
	IfDebug    bool
	nest       *SVG
}

func (this *SVG) NewPlacementWorker(binPolygon *BinPolygonStruct, paths []*PolygonStruct, ids []int, Rotations []int, config *ConfigStruct, nfpCache map[string][]Polygon) *PlacementWorkerStruct {

	return &PlacementWorkerStruct{
		binPolygon: binPolygon,
		paths:      paths,
		ids:        ids,
		Rotations:  Rotations,
		config:     config,
		nfpCache:   nfpCache,
		IfDebug:    this.config.IfDebug,
		nest:       this,
	}
}

type PositionStruct struct {
	x             float64
	y             float64
	id            int
	rotation      int
	finalNfpFloat [][]*Point //最终可放置位置
	width         float64    //此次选择放置位置 所有placed的宽
	height        float64    //此次选择放置位置 所有placed的高
	nfpArea       float64
}

func (this *PositionStruct) GetTranslateX() float64 {
	return this.x
}
func (this *PositionStruct) GetTranslateY() float64 {
	return this.y
}
func (this *PositionStruct) GetRotation() int {
	return this.rotation
}
func (this *PositionStruct) GetPoloID() int {
	return this.id
}

var placePathsNum int64 = 0

// func PrintPolygonStructList(list []*PolygonStruct) {
// 	for index := range list {
// 		this.logIfDebug("id:", list[index].id)
// 	}
// }
func (this *PlacementWorkerStruct) placePaths(paths []*PolygonStruct) *placementsStruct {
	placePathsNum++
	// if placePaths == 1 {
	// 	PrintPolygonStructList(paths)
	// }
	if this.binPolygon == nil {
		return nil
	}
	// rotate paths by given rotation
	var rotated []*PolygonStruct
	var notXiaCiRotated []*PolygonStruct
	//把所有多边形 按照所需旋转角度，旋转
	for i := 0; i < len(paths); i++ {
		paths[i].RootPoly.rotatePolygon(paths[i].rotation)
		if paths[i].isWart {
			rotated = append(rotated, paths[i])
		} else {
			notXiaCiRotated = append(notXiaCiRotated, paths[i])
		}
	}
	for index := range notXiaCiRotated {
		rotated = append(rotated, notXiaCiRotated[index])
	}
	paths = rotated
	//旋转完毕
	allplacements := [][]*PositionStruct{}
	//所有多边形 放置后的位置记载
	var fitness float64 = 10000
	binarea := math.Abs(polygonArea(this.binPolygon.myPolygon))
	var position *PositionStruct
	//paths : 未放置的碎片
	for len(paths) > 0 {
		this.logIfDebug("现在paths长度:", len(paths))
		var placed []*PolygonStruct
		var placements = []*PositionStruct{}
		//paths 有关
		fitness += 1 // add 1 for each new bin opened (lower fitness is better)
		for i := 0; i < len(paths); i++ {
			path := paths[i]
			if path.isWart {
				//如果是瑕疵，就不移动
				placements = append(placements, &PositionStruct{
					x:        0,
					y:        0,
					id:       path.id,
					rotation: 0,
				})
				placed = append(placed, path)
				continue
			}

			// inner NFP
			//该part对bin对npf
			myKey := GenPairKey(&PolygonStruct{id: -1}, path, true, 0, path.rotation)
			key := myKey.ToString()
			binNfp, isOK := this.nfpCache[key]
			// part unplaceable, skip
			if !isOK || len(binNfp) == 0 {
				//不可放置的碎片
				log.Println("出现未知NPF")
				continue
			}
			// ensure all necessary NFPs exist
			isError := false
			for j := 0; j < len(placed); j++ {
				mykey := GenPairKey(placed[j], path, false, placed[j].rotation, path.rotation)
				this.nest.GetOrCreateNFP(mykey, placed[j], path)
			}
			// part unplaceable, skip
			if isError {
				continue
			}
			if len(placed) == 0 {
				//log.Println("len(placed) == 0 ")
				// first placement, put it on the left
				for j := 0; j < len(binNfp); j++ {
					//某个npf
					for k := 0; k < len(binNfp[j]); k++ {
						//npf中的一个点
						if position == nil || binNfp[j][k].X-path.RootPoly.polygonAfterRotaion[0].X < position.x {
							//x:  npf中的一个点 和 要放置的碎片的 第一个点的 x 偏差值
							if path.isWart {
								panic("瑕疵应该先被处理")
							} else {
								position = &PositionStruct{
									x:        binNfp[j][k].X - path.RootPoly.polygonAfterRotaion[0].X,
									y:        binNfp[j][k].Y - path.RootPoly.polygonAfterRotaion[0].Y,
									id:       path.id,
									rotation: path.rotation,
								}
							}
						}
					}
					//找到要放置碎片移到npf多边形上的位移方式 ： 目前是x位移最小
				}
				placements = append(placements, position)
				placed = append(placed, path)
				continue
			}
			//?
			var clipperBinNfp []Polygon
			for j := 0; j < len(binNfp); j++ {
				clipperBinNfp = append(clipperBinNfp, toClipperCoordinates(binNfp[j]))
			}
			//扩大
			clipperBinNfpInt := ScaleUpPaths(clipperBinNfp, this.config.ClipperScale)
			clipper := Clipper(0)
			var combinedNfp []IntPolygon
			//clipperBinNfpInt 可放置的多边形区域
			//循环检测此前放置的多边形

			//找到多个NPF多边形的公共可行区域
			for j := 0; j < len(placed); j++ {
				//已放置多边形 和 该零件的 可放置 可行域 的关系
				keyPair := GenPairKey(placed[j], path, false, placed[j].rotation, path.rotation)
				if path.typeID == 9999 && path.rotation != 0 {
					log.Println("不该有的旋转")
					panic("")
				}
				key := keyPair.ToString()
				nfp, isOK := this.nfpCache[key]
				if !isOK {
					log.Println("出现未知NPF")
					continue
				}
				this.logIfDebug("取出 key:", key)
				for k := 0; k < len(nfp); k++ {
					var clone = toClipperCoordinates(nfp[k])
					for m := 0; m < len(clone); m++ {
						clone[m].X += placements[j].x
						clone[m].Y += placements[j].y
					}

					intclone := ScaleUpPath(clone, this.config.ClipperScale)
					intclone = CleanIntPolygon(intclone, 0.0001*float64(this.config.ClipperScale))
					// Execute执行的所有剪纸  必须是规整的多边形
					// Execute执行的后的所有剪纸  必须是clean为规整的多边形
					var area = Abs(Area(intclone))
					if len(intclone) > 2 && area*10 > this.config.ClipperScale*this.config.ClipperScale {
						if path.typeID == 9999 {
							this.nest.debugnfpList = append(this.nest.debugnfpList, toNestCoordinates(intclone, this.config.ClipperScale))
						}
						clipper.AddPath(intclone, ptSubject, true)
					} else {
						log.Println("什么情况啊")
						if placed[j].isWart {
							panic("跳过了一个瑕疵！")
						}
					}
				}
			}
			//计算可行域之间的关系
			combinedNfp, isExecuteOK := clipper.Execute(ctUnion, pftNonZero, pftNonZero)
			if !isExecuteOK {
				log.Println("!!!!!计算新零件和已放置零件之间位置关系时 出现不可行的情况 ")
				//如果不可行 则跳过// 此处可能出现bug//也可能可以改成多bin操作
				continue
			}

			isClockwise := Area(combinedNfp[0]) > 0
			for index2 := range combinedNfp {
				//	this.nest.debugnfpList2 = append(this.nest.debugnfpList2, toNestCoordinates(combinedNfp[index2], this.config.ClipperScale))
				if index2 == 0 {
					continue
				}
				if isClockwise == (Area(combinedNfp[index2]) > 0) {
					combinedNfp[index2] = IntPointReverse(combinedNfp[index2])
				}
			}

			// difference with bin Polygon
			var finalNfpBeforeClean = []IntPolygon{}
			clipper = Clipper(0)
			clipper.AddPaths(combinedNfp, ptClip, true)
			clipper.AddPaths(clipperBinNfpInt, ptSubject, true)
			//
			//Bin 和  可放置可行域 的关系
			finalNfpBeforeClean, isExecuteOK = clipper.Execute(ctDifference, pftNonZero, pftNonZero)
			if !isExecuteOK {
				log.Println("!!!!!计算新零件和已放置零件、bin之间位置关系时 出现不可行的情况 ")
				continue
			}
			if len(finalNfpBeforeClean) == 0 {
				log.Println("!!!!!出现 len(finalNfpBeforeClean) ==0")
				//出现这个表示无处安放
				continue
			}
			///此处必须清理！！
			var finalNfp = []IntPolygon{}
			for index := range finalNfpBeforeClean {
				finalNfp = append(finalNfp, CleanIntPolygon(finalNfpBeforeClean[index], 0.0001*float64(this.config.ClipperScale)))
			}
			//对所有对nfp多边形  进行合理性筛选 不过筛选可以更谨慎
			for j := 0; j < len(finalNfp); j++ {
				if len(finalNfp[j]) < 3 {
					if len(finalNfp[j]) < 3 {
						log.Println("！！！对所有对nfp多边形  进行合理性筛选时 出现了：len(finalNfp[j]) < 3  跳过")
					}
					tempList := []IntPolygon{}
					for finalNfpIndex := range finalNfp {
						if finalNfpIndex == j {

						} else {
							tempList = append(tempList, finalNfp[finalNfpIndex])
						}
					}
					finalNfp = tempList
					j--
				}

			}
			if len(finalNfp) == 0 {
				log.Println("!!!!!出现 len(finalNfp) ==0")
				//出现这个表示无处安放
				continue
			}
			var finalNfpScaleDown = [][]*Point{}
			for j := 0; j < len(finalNfp); j++ {
				// back to normal scale
				finalNfpScaleDown = append(finalNfpScaleDown, toNestCoordinates(finalNfp[j], this.config.ClipperScale))
			}
			finalNfpFloat := finalNfpScaleDown
			if path.typeID == 9999 {
				this.nest.debugnfp = finalNfpFloat
			}
			// choose placement that results in the smallest bounding box
			// could use convex hull instead, but it can create oddly shaped nests (triangles or long slivers) which are not optimal for real-world use
			// todo: generalize gravity direction

			position = MinWidthAndAtLeft(path, finalNfpFloat, placed, placements)
			if position != nil {
				placed = append(placed, path)
				placements = append(placements, position)
			}
		}
		//minwidth 有关
		//总之就是拼接出 所有已放置碎片的bound 宽最小的的矩形
		if position.width > 0 {
			fitness += (position.width*this.nest.WidthWeight + this.nest.LengthWeight*position.height) / binarea
		}

		for ii := 0; ii < len(placed); ii++ {
			newPaths := []*PolygonStruct{}
			for pathsIndex := range paths {
				if paths[pathsIndex] == placed[ii] {

				} else {
					newPaths = append(newPaths, paths[pathsIndex])
				}
			}
			paths = newPaths
		}

		if len(placements) > 0 {
			allplacements = append(allplacements, placements)
		} else {
			break // something went wrong
		}
	}
	//剩余碎片 有关
	fitness += 2 * float64(len(paths))
	//结果不是最佳结果 ？  为什么
	return &placementsStruct{Placements: allplacements, fitness: fitness, paths: paths, area: binarea}
}
