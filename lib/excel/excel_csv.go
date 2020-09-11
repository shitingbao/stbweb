package excel

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/text/transform"
)

//新建csv文件,输入文件名，间隔字符，数据源
//该方法生产的是utf8编码文件，需要生成gbk文件使用createGBKCsvFile方法
//注意这里要保留去除\n的操作，防止单行数据中分裂为多行（比如excel中解析出来一个单元格中多行，字符串中就会有回车符）
func createCsvFile(fileName, set string, data [][]string) error {
	if set == "" {
		set = ","
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	for _, v := range data {
		for i, val := range v {
			f.WriteString(strings.Replace(val, "\n", " ", -1))
			if i == len(v)-1 {
				break
			}
			f.WriteString(set)
		}
		f.WriteString("\r\n")
	}
	defer f.Close()
	return nil
}

//同上创建csv文件，不过这里使用了gbk编码，单独拿出来写
func createGBKCsvFile(fileName, set string, data [][]string) (err error) {
	if set == "" {
		set = ","
	}
	f, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer f.Close()
	wf := transform.NewWriter(f, enc.NewEncoder())
	defer wf.Close()
	for _, v := range data {
		for i, val := range v {
			data := strings.Replace(val, "\n", " ", -1)
			if _, err = fmt.Fprint(wf, data); err != nil {
				return err
			}
			if i == len(v)-1 {
				break
			}
			if _, err = fmt.Fprint(wf, set); err != nil {
				return err
			}
		}
		if _, err = fmt.Fprint(wf, "\r\n"); err != nil {
			return err
		}
	}
	return nil
}

//LoadCsvCfg 解析csv,utf8和gbk两种都可以，通过参数控制
//输入完整文件路径
//format为编码格式，需要解码gbk模式时，输入内容‘GBK’，utf8默认不用输入
func LoadCsvCfg(filename string, format ...string) ([][]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if filename == "" || (len(format) > 0 && format[0] != "GBK") {
		return nil, errors.New("param error")
	}

	var reader *csv.Reader
	if len(format) > 0 && format[0] == "GBK" {
		r := transform.NewReader(f, enc.NewDecoder())
		reader = csv.NewReader(r)
	} else {
		reader = csv.NewReader(f)
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
