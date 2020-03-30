package nest

import (
	"math"
)

func MinWidthAndAtLeft(path *PolygonStruct, finalNfpFloat [][]*Point, placed []*PolygonStruct, placements []*PositionStruct) (position *PositionStruct) {
	var nf Path
	var shiftvector *PositionStruct
	var minarea *float64 = nil
	var minx *float64 = nil
	var area *float64 = nil

	for j := 0; j < len(finalNfpFloat); j++ {
		nf = finalNfpFloat[j]
		if math.Abs(polygonArea(nf)) < 2 {
			continue
		}

		for k := 0; k < len(nf); k++ {
			//以前放置的
			var allpoints []*Point
			for m := 0; m < len(placed); m++ {
				for n := 0; n < len(placed[m].RootPoly.polygonAfterRotaion); n++ {
					if !placed[m].isWart {
						allpoints = append(allpoints, &Point{X: placed[m].RootPoly.polygonAfterRotaion[n].X + placements[m].x, Y: placed[m].RootPoly.polygonAfterRotaion[n].Y + placements[m].y})
					}
				}
			}

			shiftvector = &PositionStruct{
				x:             nf[k].X - path.RootPoly.polygonAfterRotaion[0].X,
				y:             nf[k].Y - path.RootPoly.polygonAfterRotaion[0].Y,
				id:            path.id,
				rotation:      path.rotation,
				finalNfpFloat: finalNfpFloat,
			}
			//当前放置的
			for m := 0; m < len(path.RootPoly.polygonAfterRotaion); m++ {
				allpoints = append(allpoints, &Point{X: path.RootPoly.polygonAfterRotaion[m].X + shiftvector.x, Y: path.RootPoly.polygonAfterRotaion[m].Y + shiftvector.y})
			}
			rectbounds := getPolygonBounds(allpoints)
			// weigh width more, to help compress in direction of gravity
			shiftvector.width = rectbounds.width
			shiftvector.height = rectbounds.height
			temparea := rectbounds.width*2 + rectbounds.height
			area = &temparea
			if minarea == nil || *area < *minarea || (_almostEqual2(*minarea, *area) && (*minx < 0 || shiftvector.x < *minx)) {
				minarea = area
				position = shiftvector
				minx = &shiftvector.x
			} else {
				//this.logIfDebug("")
			}
		}
	}
	return position
}

func MinWidthAndAtMinNfpLeft(path *PolygonStruct, finalNfpFloat [][]*Point, placed []*PolygonStruct, placements []*PositionStruct) (position *PositionStruct) {
	var nf Path
	//var shiftvector *PositionStruct
	var allpoints []*Point
	for m := 0; m < len(placed); m++ {
		for n := 0; n < len(placed[m].RootPoly.polygonAfterRotaion); n++ {
			if !placed[m].isWart {
				allpoints = append(allpoints, &Point{X: placed[m].RootPoly.polygonAfterRotaion[n].X + placements[m].x, Y: placed[m].RootPoly.polygonAfterRotaion[n].Y + placements[m].y})
			}
		}
	}
	var minWidth *float64 = nil
	var minWidthNfpPoints []*Point
	withoutPathRectbounds := getPolygonBounds(allpoints)
	positionList := []*PositionStruct{}
	for j := 0; j < len(finalNfpFloat); j++ {
		nf = finalNfpFloat[j]
		for k := 0; k < len(nf); k++ {
			//以前放置的
			shiftvector := &PositionStruct{
				x:             nf[k].X - path.RootPoly.polygonAfterRotaion[0].X,
				y:             nf[k].Y - path.RootPoly.polygonAfterRotaion[0].Y,
				id:            path.id,
				rotation:      path.rotation,
				finalNfpFloat: finalNfpFloat,
				nfpArea:       math.Abs(polygonArea(finalNfpFloat[j])),
			}
			//当前放置的
			var pathPoints []*Point
			for m := 0; m < len(path.RootPoly.polygonAfterRotaion); m++ {
				pathPoints = append(pathPoints, &Point{X: path.RootPoly.polygonAfterRotaion[m].X + shiftvector.x, Y: path.RootPoly.polygonAfterRotaion[m].Y + shiftvector.y})
			}
			pathRectbounds := getPolygonBounds(pathPoints)
			if minWidth != nil && RecBoundSum(pathRectbounds, withoutPathRectbounds).width == *minWidth {
				minWidthNfpPoints = append(minWidthNfpPoints, nf[k])
				shiftvector.width = RecBoundSum(pathRectbounds, withoutPathRectbounds).width
				shiftvector.height = RecBoundSum(pathRectbounds, withoutPathRectbounds).height
				positionList = append(positionList, shiftvector)
				continue
			}
			if minWidth == nil || RecBoundSum(pathRectbounds, withoutPathRectbounds).width < *minWidth {
				minWidth = &RecBoundSum(pathRectbounds, withoutPathRectbounds).width
				minWidthNfpPoints = []*Point{nf[k]}
				shiftvector.width = RecBoundSum(pathRectbounds, withoutPathRectbounds).width
				shiftvector.height = RecBoundSum(pathRectbounds, withoutPathRectbounds).height
				positionList = []*PositionStruct{shiftvector}
			} else {
				continue
			}
		}

	}
	var minArea float64 = -1
	var minAreaPointIndex int
	for index := range positionList {
		if minArea < 0 || positionList[index].nfpArea < minArea {
			minAreaPointIndex = index
		}
	}
	return positionList[minAreaPointIndex]
}
