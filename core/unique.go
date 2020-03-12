package core

import (
	"strings"
	"time"
)

//GetUniqueFileName 返回一个根据时间的唯一数字型字符串
func GetUniqueFileName() string {
	name := time.Now().Format("2006-01-02 15:04:05")
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, " ", "", -1)
	name = strings.Replace(name, ":", "", -1)
	return name
}

func Test() {
	redirectStderr()
}
