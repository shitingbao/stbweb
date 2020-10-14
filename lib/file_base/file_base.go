package base

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"strings"
)

//GetAllDirFile 便利所有文件内文件，反馈所有文件路径,isAbsolute代表是否反馈完整路径，或者只反馈文件名称
//isAbsolute为true反馈当前开始的完整路径，[file/aa/aa.txt file/aa/bb/bb.txt]，为false只反馈文件名，[aa.txt bb.txt]
func GetAllDirFile(url string, isAbsolute bool) ([]string, error) {
	fList := []string{}
	ft, err := ioutil.ReadDir(url)
	if err != nil {
		return fList, err
	}
	for _, v := range ft {
		if v.IsDir() {
			ft, err := GetAllDirFile(path.Join(url, v.Name()), isAbsolute)
			if err != nil {
				return fList, err
			}
			fList = append(fList, ft...)
			continue
		}
		if isAbsolute {
			fList = append(fList, path.Join(url, v.Name()))
		} else {
			fList = append(fList, v.Name())
		}
	}
	return fList, nil
}

//文件后缀操作
func fileNameOpera() {
	fullFilename := "/Users/itfanr/Documents/test.txt"
	log.Println(path.Dir(fullFilename)) //获取当前目录，"/Users/itfanr/Documents"
	var filenameWithSuffix string
	filenameWithSuffix = path.Base(fullFilename) //获取文件名带后缀(test.txt)
	fmt.Println("filenameWithSuffix =", filenameWithSuffix)

	var fileSuffix string
	fileSuffix = path.Ext(fullFilename) //获取文件后缀(.txt)
	fmt.Println("fileSuffix =", fileSuffix)

	var filenameOnly string
	filenameOnly = strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名(test)
	fmt.Println("filenameOnly =", filenameOnly)
}
