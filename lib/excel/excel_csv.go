package excel

import (
	"encoding/csv"
	"os"
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
