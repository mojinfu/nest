package nest

import (
	"fmt"
	"log"
)

func op_Equality(a, b int64) bool {
	return a == b
}

func Int128Mul(a, b int64) int64 {
	return a * b
}

func (this *ClipperStruct) Check(AllMinimaList string) {
	var AllMinimaListGo = ""
	for b := this.m_MinimaList; b != nil; b = b.Next {
		AllMinimaListGo = AllMinimaListGo + fmt.Sprintf("%d", b.Y)
	}
	if AllMinimaListGo != AllMinimaList {
		log.Println("AllMinimaListGo!=AllMinimaList")
		panic("AllMinimaListGo!=AllMinimaList")
	}
}

func (this *ClipperStruct) SortedEdges() string {
	var AllMinimaListGo = ""
	for b := this.m_SortedEdges; b != nil; b = b.NextInSEL {
		AllMinimaListGo = AllMinimaListGo + fmt.Sprintf("%d", b.Curr.X)
	}
	return AllMinimaListGo
}
