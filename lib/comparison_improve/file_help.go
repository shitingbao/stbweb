package comparisonimprove

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pborman/uuid"
)

const (
	defaultComma = "," //默认文件内容分隔符
)

var (
	//LineModel 文件数据保存总字典
	LineModel = make(map[string]FileLineRow)
)

//FileLineRow 反馈一列文件的行对象，包含文件名称，行号和未拆分的内容,sep制表符
type FileLineRow struct {
	FileName   string
	LineNumber int
	RowVal     string
	Sep        string
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

//FileComparise csv,txt比较，暂时不支持excel，因为这里是按行读取
//参数是文件名，对应文件制表符，连续两对
//返回的key是行内容，先不做按列切割
//csv按文本形式解析时，会以最长的行为基准，短的行的列不足也会有空字符，用制表符（逗号）隔开
//这里生成一个唯一uuid，作为两个解析内容完成协程的结束标识
//sep制表符这里，如果替换后数据不同，那原来肯定不同，所以不同制表符的问价替换比较没问题，数据包含该制表符也没问题,所以把输入的制表符都换成默认的比较即可
//注意输出的时候将原来的数据输出，因为数据中也可能包含制表符
func FileComparise(fileName, sep, otherFileName, osep string) {
	inData := make(chan FileLineRow, 1)
	outData := make(chan sameFileRow)
	stopLogo := uuid.New()
	go dataMatch(inData, outData)
	go getFileData(fileName, sep, stopLogo, inData)
	go getFileData(otherFileName, osep, stopLogo, inData)
	for {
		select {
		case mes := <-outData:
			log.Println(mes)
			if mes.FRow.LineNumber == -1 && mes.OtherFRow.LineNumber == -1 {
				goto out
			}
		}
	}
out:
	log.Println(LineModel)
}

//数据匹配，in通道中接受文件解析的行数据，out中反馈匹配成功的相同数据，两方不同的数据则遗留在总map中
//对应数据分别会包含对应的文件名和行号
func dataMatch(inData <-chan FileLineRow, outData chan<- sameFileRow) {
	for {
		line := <-inData
		if line.RowVal == "" { //为空就掠过
			continue
		}
		lineKey := delStrBlank(line.RowVal, line.Sep) //这个生成的key只是用于比较，数据还是原样输出
		lineData := LineModel[lineKey]
		if lineData.LineNumber == 0 { //从总数据表中获取内容,行号不会为0，0的就是无该值,将改行存入总map中
			LineModel[lineKey] = line //无该数据则放入
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
		delete(LineModel, lineKey) //行不为0说明有值，那就有相同的行，则删除，免得占用内存资源
		outData <- sm
		if line.LineNumber == -1 && lineData.LineNumber == -1 {
			break
		}
	}
}

//读取文件内容，完成后反馈一个包含-1的数据包
func getFileData(fileName, sep, stopLogo string, inData chan<- FileLineRow) {
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
			Sep:        sep,
		}
		rowNumber++
	}
	inData <- FileLineRow{
		FileName:   "",
		LineNumber: -1,
		RowVal:     stopLogo,
	}
}

//切除尾部空白,并用默认制表符代替，方便下面比较
func delStrBlank(col, sep string) string {
	str := strings.Split(col, sep)
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] != "" {
			str = str[0 : i+1]
			break
		}
	}
	pstr := ""
	for i, v := range str {
		if i == 0 {
			pstr += v
			continue
		}
		pstr += sep + v
	}
	return strings.Replace(pstr, sep, defaultComma, -1)
}
