package timeformatg

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

//gmt时间转化为time
func gmtTotime() {
	t, err := time.Parse(http.TimeFormat, "Wed, 22 Apr 2020 06:27:59 GMT")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(t.Format("2006-01-02 15:04:05"))
}

//生成gmt时间
func getGMTtime() {
	fmt.Println(time.Now().UTC().Format(http.TimeFormat))
}

//获取服务器文件更新时间信息
//header中的Last-Modified属性就是该文件的更新时间信息，格式是给gmt
func getServeFileInfo(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	lastModified := resp.Header.Get("Last-Modified")
	t, err := time.Parse(http.TimeFormat, lastModified)
	if err != nil {
		return "", err
	}
	return t.Add(time.Hour * 8).Format("2006-01-02 15:04:05"), nil
}
