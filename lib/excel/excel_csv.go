package excel

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//LoadCsvCfg 解析csv
func LoadCsvCfg(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	if reader == nil {
		return nil, err
	}
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

//Lod 列对象
type Lod struct {
	Class1 string
	Class2 string
	Name   string
	Data   string
	Unit   string
}

//ComparisonFile 比对两个csv文件不同,输入两个文件路径
func ComparisonFile(aimFilePath, contrastPath string) (*map[int]Lod, *map[int]Lod, error) {
	recodeE := []int{}
	recodeQ := []int{}
	recordE, err := LoadCsvCfg(aimFilePath)
	if err != nil {
		return nil, nil, err
	}
	reMapE := getRes(recordE)
	recordQ, err := LoadCsvCfg(contrastPath)
	if err != nil {
		return nil, nil, err
	}
	reMapQ := getRes(recordQ)
	//需要取交集的反集，所以需要两边获取两次，由于相同记录的对应位置不同，不能用固定行数比对，记录对应文件中的行数，分别比对
	for k, v := range *reMapE {
		for _, vel := range *reMapQ {
			if v == vel {
				recodeE = append(recodeE, k)
				break
			}
		}
	}
	for k, v := range *reMapQ {
		for _, vel := range *reMapE {
			if v == vel {
				recodeQ = append(recodeQ, k)
				break
			}
		}
	}
	for _, v := range recodeE {
		delete(*reMapE, v)
	}
	for _, v := range recodeQ {
		delete(*reMapQ, v)
	}
	return reMapE, reMapQ, nil
}

func getRes(recordE [][]string) *map[int]Lod {
	reMapE := make(map[int]Lod) //int记录行数
	for i, v := range recordE {
		da := Lod{}
		for idx, vel := range v {
			switch idx {
			case 0:
				da.Class1 = vel
			case 1:
				da.Class2 = vel
			case 2:
				da.Name = vel
			case 3:
				da.Data = vel
			case 4:
				da.Unit = vel
			}
		}
		reMapE[i] = da
	}
	return &reMapE
}

//新建csv文件
func createFile(fileName string, data [][]string) error {
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	for _, v := range data {
		for i, val := range v {
			f.WriteString(val)
			if i == len(v)-1 {
				break
			}
			f.WriteString(",")
		}
		f.WriteString("\r\n")
	}
	defer f.Close()
	return nil
}

var enc = simplifiedchinese.GBK

//GBK形式读csv文件
func exampleReadGBK(filename string) {
	// Read UTF-8 from a GBK encoded file.
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	r := transform.NewReader(f, enc.NewDecoder())

	// Read converted UTF-8 from `r` as needed.
	// As an example we'll read line-by-line showing what was read:
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		fmt.Printf("Read line: %s\n", sc.Bytes())
	}
	if err = sc.Err(); err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}

//GBK形式写csv文件
func exampleWriteGBK(filename string) {
	// Write UTF-8 to a GBK encoded file.
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	w := transform.NewWriter(f, enc.NewEncoder())

	// Write UTF-8 to `w` as desired.
	// As an example we'll write some text from the Wikipedia
	// GBK page that includes Chinese.
	_, err = fmt.Fprintln(w,
		`In 1995, China National Information Technology Standardization
Technical Committee set down the Chinese Internal Code Specification
(Chinese: 汉字内码扩展规范（GBK）; pinyin: Hànzì Nèimǎ
Kuòzhǎn Guīfàn (GBK)), Version 1.0, known as GBK 1.0, which is a
slight extension of Codepage 936. The newly added 95 characters were not
found in GB 13000.1-1993, and were provisionally assigned Unicode PUA
code points.`)
	if err != nil {
		log.Fatal(err)
	}

	if err = f.Close(); err != nil {
		log.Fatal(err)
	}
}

//使用utf8的包，判断是否是utf8,还有很多其他包的使用，具体参考官网 https://golang.org/pkg/unicode/utf8/#pkg-examples
func isUTF8() bool {
	valid := 'a'
	return utf8.ValidRune(valid)
}
