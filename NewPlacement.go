package nest

import (
	"log"
	"math"

	. "github.com/mojinfu/point"
)

func (this *placementsStruct) CaculateFitness() {
	var fitness float64 = 10000
	for index := range this.Placements {
		if len(this.Placements[index]) > 0 {
			fitness++
			var allpoints []*Point
			for m := 0; m < len(this.placedPaths[index]); m++ {
				for n := 0; n < len(this.placedPaths[index][m].RootPoly.polygonAfterRotaion); n++ {
					if !this.placedPaths[index][m].isWart {
						allpoints = append(allpoints, &Point{X: this.placedPaths[index][m].RootPoly.polygonAfterRotaion[n].X + this.Placements[index][m].x, Y: this.placedPaths[index][m].RootPoly.polygonAfterRotaion[n].Y + this.Placements[index][m].y})
					}
				}
			}
			rectbounds := getPolygonBounds(allpoints)
			fitness += rectbounds.width / this.binArea[index]
		}
	}
	this.fitness = fitness
}
func (this *PlacementWorkerStruct) newPlacePaths() (*placementsStruct, error) {
	if len(this.binPolygons) == 0 {
		panic("")
	}
	if len(this.paths) == 0 {
		panic("")
	}
	for i := 0; i < len(this.paths); i++ {
		this.paths[i].rotatePolygon(this.Rotations[i])
	}

	var placed [][]*PolygonStruct
	binareaList := []float64{}
	var placements = [][]*PositionStruct{}
	for index := range this.binPolygons {
		binareaList = append(binareaList, math.Abs(polygonArea(this.binPolygons[index].RootPoly.polygonBeforeRotation)))
		placed = append(placed, []*PolygonStruct{})
		placements = append(placements, []*PositionStruct{})
	}
	for binIndex := range this.warts {
		for j := range this.warts[binIndex] {
			wartPosition := &PositionStruct{
				x:        0,
				y:        0,
				id:       this.warts[binIndex][j].id,
				rotation: 0,
			}
			placements[binIndex] = append(placements[binIndex], wartPosition)
			placed[binIndex] = append(placed[binIndex], this.warts[binIndex][j])
		}
	}

	for i := 0; i < len(this.paths); i++ {
		path := this.paths[i]
		var position *PositionStruct = nil
		for binIndex := range this.binPolygons {
			position = nil
			myKey := GenPairKey(this.binPolygons[binIndex], path, true, 0, path.rotation)
			key := myKey.ToString()
			binNfp, isOK := this.nfpCache[key]
			if !isOK {
				//不可放置的碎片
				log.Println("出现未知NPF:", myKey.ToString())
				panic("")
			}
			if len(binNfp) == 0 {
				continue
			}
			for j := 0; j < len(placed[binIndex]); j++ {
				mykey := GenPairKey(placed[binIndex][j], path, false, placed[binIndex][j].rotation, path.rotation)
				this.nest.GetOrCreateNFP(mykey, placed[binIndex][j], path)
			}
			if len(placed[binIndex]) == 0 {
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
				placements[binIndex] = append(placements[binIndex], position)
				placed[binIndex] = append(placed[binIndex], path)
				break
			} else {
				var clipperBinNfp []Polygon
				for j := 0; j < len(binNfp); j++ {
					clipperBinNfp = append(clipperBinNfp, toClipperCoordinates(binNfp[j]))
				}
				//扩大
				clipperBinNfpInt := ScaleUpPaths(clipperBinNfp, this.config.ClipperScale)
				clipper := Clipper(0)
				var combinedNfp []IntPolygon

				//找到多个NPF多边形的公共可行区域
				for j := 0; j < len(placed[binIndex]); j++ {
					//已放置多边形 和 该零件的 可放置 可行域 的关系
					keyPair := GenPairKey(placed[binIndex][j], path, false, placed[binIndex][j].rotation, path.rotation)
					key := keyPair.ToString()
					nfp, isOK := this.nfpCache[key]
					if !isOK {
						panic("出现未知NPF:" + keyPair.ToString())
					}

					for k := 0; k < len(nfp); k++ {
						var clone = toClipperCoordinates(nfp[k])
						for m := 0; m < len(clone); m++ {
							clone[m].X += placements[binIndex][j].x
							clone[m].Y += placements[binIndex][j].y
						}
						intclone := ScaleUpPath(clone, this.config.ClipperScale)
						intclone = CleanIntPolygon(intclone, 0.0001*float64(this.config.ClipperScale))
						// Execute执行的所有剪纸  必须是规整的多边形
						// Execute执行的后的所有剪纸  必须是clean为规整的多边形
						var area = Abs(Area(intclone))
						if len(intclone) > 2 && area*10 > this.config.ClipperScale*this.config.ClipperScale {
							// if path.typeID == 9999 {
							// 	this.nest.debugnfpList = append(this.nest.debugnfpList, toNestCoordinates(intclone, this.config.ClipperScale))
							// }
							clipper.AddPath(intclone, ptSubject, true)
						} else {
							log.Println("什么情况啊")
							if placed[binIndex][j].isWart {
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
				//Bin 和  可放置可行域 的关系
				finalNfpBeforeClean, isExecuteOK = clipper.Execute(ctDifference, pftNonZero, pftNonZero)
				if !isExecuteOK {
					log.Println("!!!!!计算新零件和已放置零件、bin之间位置关系时 出现不可行的情况 ")
					continue
				}
				if len(finalNfpBeforeClean) == 0 {
					log.Println("!!!!!出现 len(finalNfpBeforeClean) ==0")
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
				position = MinWidthAndAtLeft(path, finalNfpFloat, placed[binIndex], placements[binIndex])
				if position != nil {
					placed[binIndex] = append(placed[binIndex], path)
					placements[binIndex] = append(placements[binIndex], position)
					break
				}
			}
			//结束bin的循环
		}
		if position == nil {
			return nil, CanNotPutErr
		}
	}
	myPlace := &placementsStruct{Placements: placements, placedPaths: placed, binArea: binareaList}
	myPlace.CaculateFitness()
	return myPlace, nil
}
