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
