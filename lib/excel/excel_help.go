package excel

import (
	"time"

	"github.com/Sirupsen/logrus"
)

var (
	defaultExcelDate     = "1900-01-01"
	defaultExcelDateTime = "1900-01-01 00:00:00"
	defaultDate          = "2006-01-02"
	defaultDateTime      = "2006-01-02 15:04:05"
)

//numToDate 间隔的天数，转化为日期
func numToDate(distace int) time.Time {
	dDate, err := time.Parse(defaultDate, defaultExcelDate)
	if err != nil {
		logrus.WithFields(logrus.Fields{"date": err.Error()}).Error("excel")
	}

	return dDate.AddDate(0, 0, distace)
}

//numToDateTime 间隔的天数，转化为日期
func numToDateTime(distace float64) {
	// Excel中的时间24小时对应数字1，相应的1/24/60/60=0.0000115740740740741对应1秒
	//反式转化可能会有误差的问题，因为他这里已经使用小数了
	//完整的时间格式推荐使用文本格式来做
}
