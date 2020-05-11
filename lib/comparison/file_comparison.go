package comparison

import (
	"log"
	"reflect"
)

//FileComparison 是否有标题，分隔符时什么，csv也可以定义分隔符
//主要三种不同文件类型的相互比较，txt，csv和excel，对是否有标题有影响
//文件带有标题，则以标题key为基准，不带有标题或者忽略标题比较则仅仅对内容进行比较
func FileComparison() {}

func GetTitleLineGroup(fileName, sep string) {
	res := getTitleLineGroup(fileName, sep)
	for i, v := range res {
		log.Println("i:", i)
		log.Println("v:", v)
	}
}
func GetLineGroup(fileName, sep string) {
	res := getLineGroup(fileName, sep)
	for i, v := range res {
		log.Println("i:", i)
		log.Println("v:", v)
	}
}
func ExcelTitleLineGroup(fileName string) {
	res := excelTitleLineGroup(fileName)
	for i, v := range res {
		log.Println("i:", i)
		log.Println("v:", v)
	}
}
func ExcelLineGroup(fileName string) {
	res := excelLineGroup(fileName)
	for i, v := range res {
		log.Println("i:", i)
		log.Println("v:", v)
	}
}

//FileSameLineList 反馈两个文件之间相同数据集和对应文件内的行号
type FileSameLineList struct {
	LeftRow  int
	RightRow int
	Data     LineMode
}

//ParisonResult 行号和文件数据集,相同数据集，以及左右两个目标文件不同的数据集
type ParisonResult struct {
	SameDataLists []FileSameLineList
	LeftAims      map[int]LineMode
	RightAims     map[int]LineMode
}

//比较lineMode不同
//记录相同记录的内容和对应文件行号
//根据相同数据去除对应源数据内容，筛选剩余数据
func lineModeComparison(obj, objSep map[int]LineMode) ParisonResult {
	sameData := []FileSameLineList{}
	for i, v := range obj {
		for idx, val := range objSep {
			if reflect.DeepEqual(v, val) {
				sData := FileSameLineList{
					LeftRow:  i,
					RightRow: idx,
					Data:     v,
				}
				sameData = append(sameData, sData)
				break
			}
		}
	}
	for _, v := range sameData {
		delete(obj, v.LeftRow)
		delete(objSep, v.RightRow)
	}
	return ParisonResult{
		SameDataLists: sameData,
		LeftAims:      obj,
		RightAims:     objSep,
	}
}

//比较lineModeBool不同
func lineModeBoolComparison(obj, objSep map[int]LineModeBool) {

}

//比较lineMode和lineModeBool不同
func otherComparison(obj map[int]LineMode, objSep map[int]LineModeBool) {}

func searchLineMode() {

}
