package comparison

import (
	"bufio"
	"os"
	"strings"
)

type lineMode map[string]string
type lineModeBool map[string]bool

//FileComparison 是否有标题，分隔符时什么，csv也可以定义分隔符
func FileComparison() {}

//第一行为标题时,获取行组
func getTitleLineGroup(fileName, sep string) map[int]lineMode {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	result := make(map[int]lineMode)
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
func sTmap(title, cot []string) lineMode {
	res := make(lineMode)
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

//第一行不为标题时,获取行组
func getLineGroup(fileName, sep string) map[int]lineModeBool {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	result := make(map[int]lineModeBool)
	scanner := bufio.NewScanner(file)
	i := 1
	for scanner.Scan() {
		strList := strings.Split(scanner.Text(), sep)
		result[i] = sTBoolMap(strList)
		i++
	}
	return result
}

func sTBoolMap(cot []string) lineModeBool {
	res := make(lineModeBool)
	for _, v := range cot {
		res[v] = true
	}
	return res
}
