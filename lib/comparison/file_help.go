//Package comparison 该包适用于txt以及csv文件，excel文件使用excel中的比对
package comparison

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"stbweb/lib/excel"
	"strings"

	"github.com/Sirupsen/logrus"
)

const (
	space = " "
	comma = ","
	other = ""
)

//逐行读取的三种基础方法
//csv内容实质也是文本，所以txt和csv的解析格式和过程都一样，而excel不同，需要另外解析
func readLineFile(fileName string) {
	if file, err := os.Open(fileName); err != nil {
		panic(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log.Println("NewScanner:", scanner.Text())
		}
	}
}

//如果有空行，这个方法会多一行，因为最后一行也可能有回车转义符
func readFileLine(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		log.Println("n:", line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func readLine(fileName string) {
	fi, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		log.Println("line:", string(a))
	}
}

//////////////////////////////////////////////////////////

//LineMode 列对象
type LineMode []string

//csv,txt获取行组
//返回的key是行号
//csv按文本形式解析时，会以最长的行为基准，短的行的列不足也会有空字符，用制表符（逗号）隔开
func getLineGroup(fileName, sep string) map[int]LineMode {
	file, err := os.Open(fileName)
	if err != nil {
		logrus.WithFields(logrus.Fields{"file error": err.Error()}).Error("parsing file have err")
		return nil
	}
	result := make(map[int]LineMode)
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		strList := strings.Split(scanner.Text(), sep)
		result[i] = strList
		i++
	}
	return result
}

//excel获取行组
func getExcelLineGroup(fileName string) map[int]LineMode {
	resultList, err := excel.LoadCsvCfg(fileName)
	if err != nil {
		return nil
	}
	result := make(map[int]LineMode)
	for i, v := range resultList {
		result[i+1] = v
	}
	return result
}
