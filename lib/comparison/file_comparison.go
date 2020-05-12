package comparison

import (
	"errors"
	"os"
	"path"
	"reflect"
)

//ParisonFileObject 比较文件对象
type ParisonFileObject struct {
	FileName string
	Sep      string
	IsTitle  bool
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
func lineModeComparison(objData, objSepData map[int]LineMode) ParisonResult {
	sameData := []FileSameLineList{}
	for i, v := range objData {
		for idx, val := range objSepData {
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
		delete(objData, v.LeftRow)
		delete(objSepData, v.RightRow)
	}
	return ParisonResult{
		SameDataLists: sameData,
		LeftAims:      objData,
		RightAims:     objSepData,
	}
}

//FileComparison 是否有标题，分隔符时什么，csv也可以定义分隔符
//主要三种不同文件类型的相互比较，txt，csv和excel，对是否有标题有影响
//文件带有标题，则以标题key为基准，不带有标题或者忽略标题比较则仅仅对内容进行比较
func FileComparison(fn, ofn ParisonFileObject) error {
	// getLineGroup("./assets/gg.csv", ",")
	if tp, err := checkFile(fn.FileName); err != nil {

	} else {
		switch tp {
		case ".csv", ".txt":
			getLineGroup(fn.FileName, fn.Sep)
		default:
			getExcelLineGroup(fn.FileName)
		}
	}
	if tp, err := checkFile(ofn.FileName); err != nil {

	} else {
		switch tp {
		case ".csv", ".txt":
			getLineGroup(fn.FileName, fn.Sep)
		default:
			getExcelLineGroup(fn.FileName)
		}
	}

	//判断文件类型
	return nil
}

func checkFile(fileName string) (string, error) {
	filenameWithSuffix := path.Base(fileName) //获取文件名带后缀
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀

	switch fileSuffix {
	case ".csv", ".txt", ".xlsx", ".xls":
	default:
		return "", errors.New("文件类型不匹配")
	}
	_, err := os.Stat(fileName)
	if err != nil {
		return "", err
	}
	return fileSuffix, nil
}
