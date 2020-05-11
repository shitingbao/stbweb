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

type LineMode map[string]string
type LineModeBool map[string]bool

//csv,txt第一行为标题时,获取行组
//返回的key是行号，下同
//csv按文本形式解析时，会以最长的行为基准，短的行的列不足也会有空字符，用制表符（逗号）隔开
func getTitleLineGroup(fileName, sep string) map[int]LineMode {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	result := make(map[int]LineMode)
	var title []string
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		strList := strings.Split(scanner.Text(), sep)
		if i == 1 {
			title = strList
		} else {
			result[i] = sTmap(title, strList)
		}
		i++
	}
	return result
}

//将两个字符串抓转化成map
func sTmap(title, cot []string) LineMode {
	res := make(LineMode)
	switch {
	case len(title) == len(cot):
		for i, v := range title {
			res[v] = cot[i]
		}
	case len(title) > len(cot):
		for i, v := range title {
			if i <= len(cot)-1 {
				res[v] = cot[i]
				continue
			}
			res[v] = ""
		}
	case len(title) < len(cot): //如果标题比内容短，那就反过来设置键值对，所有都采用这种方法就不会错
		for i, v := range cot {
			if i <= len(title)-1 {
				res[v] = title[i]
				continue
			}
			res[v] = ""
		}
	}
	return res
}

//csv,txt第一行不为标题时,获取行组
func getLineGroup(fileName, sep string) map[int]LineModeBool {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	result := make(map[int]LineModeBool)
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		strList := strings.Split(scanner.Text(), sep)
		result[i] = sTBoolMap(strList)
		i++
	}
	return result
}

func sTBoolMap(cot []string) LineModeBool {
	res := make(LineModeBool)
	for _, v := range cot {
		res[v] = true
	}
	return res
}

//excel第一行为标题时，获取行组
func excelTitleLineGroup(fileName string) map[int]LineMode {
	resultList, err := excel.LoadCsvCfg(fileName)
	if err != nil {
		return nil
	}
	var title []string
	result := make(map[int]LineMode)
	for i, v := range resultList {
		if i == 0 {
			title = v
			continue
		}
		result[i+1] = sTmap(title, v)
	}
	return result
}

//excel第一行不为标题时，获取行组
func excelLineGroup(fileName string) map[int]LineModeBool {
	resultList, err := excel.LoadCsvCfg(fileName)
	if err != nil {
		return nil
	}
	result := make(map[int]LineModeBool)
	for i, v := range resultList {
		result[i+1] = sTBoolMap(v)
	}
	return result
}
