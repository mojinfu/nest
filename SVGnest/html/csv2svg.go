package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/mojinfu/UselessHelper"
)

func GetHtmlPol(a string) string {
	return fmt.Sprintf(`<Polygon fill="none" stroke="#010101" stroke-miterlimit="10" points=%s/>`, a)
}
func GetHtmlPolPoints(a [][2]float64) string {
	all := `"`
	for index := range a {
		all = all + fmt.Sprintf("%v,%v", a[index][0], a[index][1])
		if index != len(a)-1 {
			all = all + " "
		} else {
			all = all + `"`
		}
	}
	return all
}

func GetArray() []*withcode {
	myCsv := LoadCsvCfg("./L000"+PC+"_lingjian.csv", 100)
	//log.Println(myCsv.Records[0].Record["外轮廓"])
	// log.Println(len(myCsv.Records))
	allPol := []*withcode{}
	for index := range myCsv.Records {
		mywithcode := &withcode{}
		mywithcode.lingjianCode = myCsv.Records[index].Record["零件号"]
		mywithcode.mianliaoCode = myCsv.Records[index].Record["面料号"]
		mywithcode.xialiaoCode = myCsv.Records[index].Record["下料批次号"]
		PolStr := myCsv.Records[index].Record["外轮廓"]
		Pol := [][2]float64{}
		err := json.Unmarshal([]byte(PolStr), &Pol)
		if err != nil {
			//log.Println(err)
			panic("")
		} else {
			//log.Println(Pol)

			mywithcode.polo = Pol
		}
		allPol = append(allPol, mywithcode)
	}
	return allPol
}
func svg2Csv() string {
	result := "下料批次,零件号,面料号,零件外轮廓线坐标\n"
	filename := "./img/out" + PC + ".html"
	file, err := os.Open(filename)
	if err != nil {
		beego.Error(err)
		return ""
	}
	defer file.Close()

	records, err := ioutil.ReadAll(file)
	if err != nil {

		return ""
	}
	recordList := strings.Split(string(records), "</g>")
	allArr := GetArray()
	for index := range recordList {
		myPol := parsePointsFromHtml(UselessHelper.GetHtmlValue(recordList[index], "points"))
		myRolate := parseRotateFromHtml(UselessHelper.GetHtmlValue(recordList[index], "transform"))
		myTranslate := parseTranslateFromHtml(UselessHelper.GetHtmlValue(recordList[index], "transform"))

		myWithCode := GetWithCodeInfo(myPol, allArr)
		if myTranslate[0] == 67996.43416300001 {
			log.Println("found ", getOutCSVLine(myRolate, myTranslate, myWithCode))
		}
		result = result + getOutCSVLine(myRolate, myTranslate, myWithCode)
	}
	return result
}
func Rolate180(pol [][2]float64) [][2]float64 {
	after := [][2]float64{}
	for index := range pol {
		after = append(after, [2]float64{
			pol[index][0] * -1,
			pol[index][1] * -1,
		})

	}
	return after
}
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.4f", value), 64)
	return value
}
func Translate(pol [][2]float64, myTranslate [2]float64) [][2]float64 {
	var xOffset float64 = 0
	var yOffset float64 = 0
	//初赛暂时不启用
	//材料边角
	after := [][2]float64{}
	for index := range pol {
		after = append(after, [2]float64{
			Decimal(pol[index][0] + myTranslate[0] + xOffset),
			Decimal(pol[index][1] + myTranslate[1] + yOffset),
		})
	}
	return after
}
func getOutCSVLine(myRolate int64, myTranslate [2]float64, myWithCode *withcode) string {
	myPolo := [][2]float64{}
	if myRolate == 0 {
		myPolo = myWithCode.polo
	} else if myRolate == 180 {
		myPolo = Rolate180(myWithCode.polo)
	} else {
		panic("")
	}
	myPolo = Translate(myPolo, myTranslate)
	Area = Area + GetArea(myPolo)
	// if MaxX > 20000 {
	// 	log.Println(myWithCode)
	// 	panic("err")
	// }
	return fmt.Sprintf("%s,%s,%s,\"%s\"\n", "L000"+PC, myWithCode.lingjianCode, myWithCode.mianliaoCode, getSvcPointListFmt(myPolo))
}

func getSvcPointListFmt(pointList [][2]float64) string {
	all := "["
	for index := range pointList {
		all = all + fmt.Sprintf("[%v, %v]", pointList[index][0], pointList[index][1])
		if index == len(pointList)-1 {
			all = all + "]"
		} else {
			all = all + ", "
		}
	}
	return all
}
func parseFloatToStr(v string) float64 {
	v1, err := strconv.ParseFloat(v, 64)
	if err != nil {
		panic("float")
	}
	return v1
}
func parseRotateFromHtml(pointsStr string) int64 {
	pointsStr = strings.Trim(pointsStr, `"`)
	pointsAList := strings.Split(pointsStr, "(")
	pointsBList := strings.Split(pointsAList[2], ")")

	i, err := strconv.ParseInt(pointsBList[0], 10, 64)
	if err != nil {
		panic("it")
	}
	return i

}
func parseTranslateFromHtml(pointsStr string) [2]float64 {
	pointsStr = strings.Trim(pointsStr, `"`)
	pointsAList := strings.Split(pointsStr, "(")
	pointsBList := strings.Split(pointsAList[1], ")")
	pointsCList := strings.Split(pointsBList[0], " ")
	return [2]float64{
		parseFloatToStr(pointsCList[0]),
		parseFloatToStr(pointsCList[1]),
	}
}
func parsePointsFromHtml(pointsStr string) [][2]float64 {
	all := [][2]float64{}
	//	pointsStr = pointsStr[1 : len(pointsStr)-2]
	pointsStr = strings.Trim(pointsStr, `"`)
	pointsList := strings.Split(pointsStr, " ")
	for index := range pointsList {
		Point := strings.Split(pointsList[index], ",")
		all = append(all,
			[2]float64{
				parseFloatToStr(Point[0]), parseFloatToStr(Point[1]),
			},
		)

	}
	return all
}

type withcode struct {
	polo         [][2]float64
	lingjianCode string
	mianliaoCode string
	xialiaoCode  string
}

func isEqual(a, b [][2]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i][0] == b[i][0] && a[i][1] == b[i][1] {
			continue
		} else {
			return false
		}
	}
	return true
}
func GetWithCodeInfo(pointInfo [][2]float64, arr []*withcode) *withcode {
	for i := range arr {
		if isEqual(arr[i].polo, pointInfo) {
			return arr[i]
		}
	}
	panic("未找到")
	return nil
}

func CanSubmit(fileName string) {
	myCsv := LoadCsvCfg(fileName, 100)
	//log.Println(myCsv.Records[0].Record["外轮廓"])
	// log.Println(len(myCsv.Records))
	allPol := [][][2]float64{}
	for index := range myCsv.Records {
		PolStr := myCsv.Records[index].Record["零件外轮廓线坐标"]
		Pol := [][2]float64{}
		err := json.Unmarshal([]byte(PolStr), &Pol)
		if err != nil {
			//log.Println(err)
			panic("")
		} else {
			//log.Println(Pol)
		}
		allPol = append(allPol, Pol)
	}

	for i := 0; i < len(allPol)-1; i++ {
		for j := 0; j < len(allPol)-1; j++ {
			if i == j {
				continue
			} else {
				GetDistencePolo(allPol[i], allPol[j])
			}
		}
	}

	for i := 0; i < len(allPol)-1; i++ {
		for j := range allPol[i] {
			if allPol[i][j][0] > 20000 || allPol[i][j][0] < 0 {
				log.Println(allPol[i][j][0])
				panic("out range")
			}
			if allPol[i][j][1] > 1600 || allPol[i][j][1] < 0 {
				log.Println(allPol[i][j][1])
				panic("out range")
			}
		}
	}

	if MinDistencePoint < 25 {
		panic("too close")
	}

}

func csv2svg() string {
	myCsv := LoadCsvCfg("./L0003_lingjian.csv", 100)
	//log.Println(myCsv.Records[0].Record["外轮廓"])

	// log.Println(len(myCsv.Records))
	allPol := [][][2]float64{}
	for index := range myCsv.Records {
		PolStr := myCsv.Records[index].Record["外轮廓"]
		Pol := [][2]float64{}
		err := json.Unmarshal([]byte(PolStr), &Pol)
		if err != nil {
			//log.Println(err)
			return ""
		} else {
			//log.Println(Pol)
		}
		allPol = append(allPol, Pol)
	}

	//log.Println(GetHtmlPolPoints(allPol[0]))
	all := ""
	for index := range allPol {
		all = all + GetHtmlPol(GetHtmlPolPoints(allPol[index]))
		all = all + "\n"
	}

	return all
}

type CsvTable struct {
	FileName string
	Records  []CsvRecord
}

type CsvRecord struct {
	Record map[string]string
}

func (c *CsvRecord) GetInt(field string) int {
	var r int
	var err error
	if r, err = strconv.Atoi(c.Record[field]); err != nil {
		beego.Error(err)
		panic(err)
	}
	return r
}

func (c *CsvRecord) GetString(field string) string {
	data, ok := c.Record[field]
	if ok {
		return data
	} else {
		beego.Warning("Get fileld failed! fileld:", field)
		return ""
	}
}

func LoadCsvCfg(filename string, row int) *CsvTable {
	file, err := os.Open(filename)
	if err != nil {
		beego.Error(err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		beego.Error("NewReader return nil, file:", file)
		return nil
	}
	records, err := reader.ReadAll()
	if err != nil {
		beego.Error(err)
		return nil
	}
	// if len(records) < row {
	// 	//beego.Warning(filename, " is empty")
	// 	return nil
	// }
	colNum := len(records[0])
	recordNum := len(records)
	var allRecords []CsvRecord
	for i := 1; i < recordNum; i++ {
		record := &CsvRecord{make(map[string]string)}
		for k := 0; k < colNum; k++ {
			record.Record[records[0][k]] = records[i][k]
		}
		allRecords = append(allRecords, *record)
	}
	var result = &CsvTable{
		filename,
		allRecords,
	}
	return result
}

func GetArea(a [][2]float64) float64 {
	for index := range a {
		if a[index][0] > MaxX {
			MaxX = a[index][0]
		}
	}
	pointNum := len(a)
	if 3 > pointNum {
		return 0
	}
	//相邻两点依次操作
	var c float64 = 0
	var e = 0
	var d = pointNum - 1
	for ; e < pointNum; e++ {
		// log.Println("d:", d)
		// log.Println("e:", e)
		// log.Println("a[d].X + a[e].X:", a[d].X+a[e].X)
		// log.Println("(a[d].Y - a[e].Y):", (a[d].Y - a[e].Y))
		c += (a[d][0] + a[e][0]) * (a[d][1] - a[e][1])
		d = e
	}
	if c < 0 {
		return c / 2 * -1
	} else {
		return c / 2
	}

}
func GetDistencePolo(a, b [][2]float64) {
	for i := range a {
		for j := 0; j < len(b)-2; j++ {
			dis := cnFang(a[i][0], a[i][1], b[j][0], b[j][1], b[j+1][0], b[j+1][1])
			if MinDistencePoint > dis {
				MinDistencePoint = dis
			}
		}

	}

}
func GetDistence(a, b [2]float64) float64 {

	//if (a[0]-b[0])*(a[0]-b[0])+(a[1]-b[1])*(a[1]-b[1]) < 20 {
	return (a[0]-b[0])*(a[0]-b[0]) + (a[1]-b[1])*(a[1]-b[1])
	//	panic("err:")
	//}

}
