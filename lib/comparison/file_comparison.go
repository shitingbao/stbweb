package comparison

import (
	"errors"
	"os"
	"path"
	"reflect"
)

//ParisonFileObject 比较文件对象
//传入文件名，文件分隔符（excel文件不用，并且默认是","）和第一行是否是标题的标志(true为是标题)
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
//标题行不参与比对
//如果第一行是标题，就直接delete第一行，不参与比较
func FileComparison(fn, ofn ParisonFileObject) (ParisonResult, error) {
	fnData, err := getFileDataLists(fn)
	if err != nil {
		return ParisonResult{}, err
	}
	if fn.IsTitle && len(fnData) > 1 {
		delete(fnData, 1)
	}
	ofnData, err := getFileDataLists(ofn)
	if err != nil {
		return ParisonResult{}, err
	}
	if fn.IsTitle && len(ofnData) > 1 {
		delete(ofnData, 1)
	}
	return lineModeComparison(fnData, ofnData), nil
}

//根据文件类型解析文件，反馈文件内容对象map
func getFileDataLists(param ParisonFileObject) (map[int]LineMode, error) {
	tp, err := checkFile(param.FileName)
	if err != nil {
		return nil, err
	}
	switch tp {
	case ".csv", ".txt":
		return getLineGroup(param.FileName, param.Sep), nil
	default:
		return getExcelLineGroup(param.FileName), nil
	}
}

//检查文件存在和类型
func checkFile(fileName string) (string, error) {
	if fileName == "" {
		return "", errors.New("file name can not nil")
	}
	filenameWithSuffix := path.Base(fileName) //获取文件名带后缀
	var fileSuffix string
	fileSuffix = path.Ext(filenameWithSuffix) //获取文件后缀

	switch fileSuffix {
	case ".csv", ".txt", ".xlsx":
	default:
		return "", errors.New("file type not match")
	}
	_, err := os.Stat(fileName)
	if err != nil {
		return "", err
	}
	return fileSuffix, nil
}
