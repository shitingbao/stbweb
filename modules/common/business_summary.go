package common

import (
	"log"
	"net/http"
	"stbweb/core"
	"time"
)

//统计,反馈顾客和对应统计的信息
type summaryInfo struct {
	Customer string
	Summary  string
}

func init() {
	core.RegisterFun("summary", new(summaryInfo), false)
}

func (s *summaryInfo) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("info", nil, summaryProtectProMon) {
		return
	}
}

//上月汇总
func summaryProtectProMon(param interface{}, p *core.ElementHandleArgs) error {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	sql := "SELECT customer,SUM(price*num) as summary FROM order_info WHERE out_time BETWEEN ? AND ? GROUP BY customer"
	results, err := getSummaryInfo(sql, thisMonth.AddDate(0, -1, 0), thisMonth.AddDate(0, 0, -1))
	if err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, results)
	return nil
}

//本月汇总
func summaryProtectThisMon(param interface{}, p *core.ElementHandleArgs) error {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	sql := "SELECT customer,SUM(price*num) as summary FROM order_info WHERE out_time BETWEEN ? AND ? GROUP BY customer"
	results, err := getSummaryInfo(sql, thisMonth, time.Now())
	if err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, results)
	return nil
}

func getSummaryInfo(sql string, start, stop time.Time) ([]summaryInfo, error) {
	rows, err := core.Ddb.Query(sql, start, stop)
	if err != nil {
		return nil, err
	}
	var results []summaryInfo
	for rows.Next() {
		var res summaryInfo
		rows.Scan(&res.Customer, &res.Summary)
		results = append(results, res)
	}
	return results, nil
}

//条件定义，入库时间，支付时间和是否支付条件分别选择
//如果开始时间或者结束时间只写了一个，那就是对应的时间大于结束时间或者开始时间
//example,OutTime只有一个StartTime,那就是入库时间小于该时间，只有一个StopTime，就是入库时间大于该时间，都存在则取中间值，其他同理
//IsPay内容为0，1，2，代表无条件，已支付或未支付
type customizeParam struct {
	OutTime outTime
	PayTime outTime
	IsPay   int
}

type outTime struct {
	StartTime time.Time
	StopTime  time.Time
}

//Post 自定义
func (s *summaryInfo) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("customize", new(customizeParam), customize) {
		return
	}
}

func customize(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*customizeParam)
	log.Println("pa:", pa)
	log.Println("StartTime:", pa.OutTime.StartTime.IsZero())
	log.Println("StopTime:", pa.OutTime.StopTime.IsZero())
	return nil
}

func getWhere(cus customizeParam) string {
	switch {
	case cus.OutTime.StartTime.IsZero() || cus.OutTime.StopTime.IsZero(): //入库时间参数有一个
	case cus.OutTime.StartTime.IsZero() && cus.OutTime.StopTime.IsZero(): //入库时间参数都有
	case cus.PayTime.StartTime.IsZero() || cus.PayTime.StopTime.IsZero(): //支付时间参数有一个
	case cus.PayTime.StartTime.IsZero() && cus.PayTime.StopTime.IsZero(): //支付时间参数都有
	default:
	}
	return ""
}
