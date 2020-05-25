package comparisonimprove

import (
	"bufio"
	"log"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

// const (
// 	defaultComma = "," //默认文件内容分隔符
// )

var (
	//LineModel 文件数据保存总字典
	LineModel = make(map[string]FileLineRow)
)

//FileLineRow 反馈一列文件的行对象，包含文件名称，行号和未拆分的内容
type FileLineRow struct {
	FileName   string
	LineNumber int
	RowVal     string
}

//反馈对应相同的数据行，包含对应文件行号
type sameFileRow struct {
	FRow      FileLine
	OtherFRow FileLine
	RowVal    string
}

//FileLine 包含文件名和对应行号
type FileLine struct {
	FileName   string
	LineNumber int
}

//csv,txt获取行组
//返回的key是行号
//csv按文本形式解析时，会以最长的行为基准，短的行的列不足也会有空字符，用制表符（逗号）隔开
func getLineGroup(fileName, sep string) {
	st := time.Now()
	inData := make(chan FileLineRow)
	outData := make(chan sameFileRow)
	go dataMatch(inData, outData)
	go getFileData("", inData)
	go getFileData("", inData)
	for {
		select {
		case mes := <-outData:
			log.Println(mes)
		default:
			nt := time.Since(st)
			if nt > time.Second*2 {
				log.Println(nt)
				goto out
			}
		}
	}
out:
	log.Println(LineModel)
}

// //切除尾部空白
// func deleteStrBlank(str []string) LineMode {
// 	for i := len(str) - 1; i >= 0; i-- {
// 		if str[i] != "" {
// 			str = str[0 : i+1]
// 			return str
// 		}
// 	}
// 	return make(LineMode, 0)
// }

func dataMatch(inData <-chan FileLineRow, outData chan<- sameFileRow) {
	for {
		line := <-inData
		if line.RowVal == "" { //为空就掠过
			continue
		}
		lineData := LineModel[line.RowVal]
		if lineData.LineNumber == 0 { //从总数据表中获取内容,行号不会为0，0的就是无该值,将改行存入总map中
			LineModel[line.RowVal] = FileLineRow{
				FileName:   line.FileName,
				LineNumber: lineData.LineNumber,
				RowVal:     line.RowVal,
			} //无该数据则放入
			continue
		}
		f := FileLine{
			FileName:   line.FileName,
			LineNumber: line.LineNumber,
		}
		of := FileLine{
			FileName:   lineData.FileName,
			LineNumber: lineData.LineNumber,
		}
		sm := sameFileRow{
			FRow:      f,
			OtherFRow: of,
			RowVal:    line.RowVal,
		}
		delete(LineModel, line.RowVal) //行不为0说明有值，那就有相同的行，则删除，免得占用内存资源
		outData <- sm
	}
}

func getFileData(fileName string, inData chan<- FileLineRow) {
	file, err := os.Open(fileName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"file error": err.Error()}).Error("parsing file have err")
		return
	}
	scanner := bufio.NewScanner(file)
	rowNumber := 1
	for scanner.Scan() {
		inData <- FileLineRow{
			FileName:   fileName,
			LineNumber: rowNumber,
			RowVal:     scanner.Text(),
		}
		rowNumber++
	}
}
