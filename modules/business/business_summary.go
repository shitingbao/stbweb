package business

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

//条件定义，入库时间，支付时间和是否支付条件分别选择
//如果开始时间或者结束时间只写了一个，那就是对应的时间小于结束时间或者大于开始时间
//example,OutTime只有一个StartTime,那就是入库时间大于该时间，只有一个StopTime，就是入库时间小于该时间，都存在则取中间值，其他同理
//IsPay内容为0，1，2，0代表无条件，可以对paytime参数进行赋值，反馈不同条件查询，1和2代表已支付或未支付，pattime参数将不起作用
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
	strWhere, paramList := getWhere(pa)
	if strWhere != "" {
		strWhere = " where " + strWhere
	}
	sql := "SELECT customer,SUM(price*num) as summary FROM order_info " + strWhere + " GROUP BY customer"
	log.Println("sql:", sql)
	res, err := getSummaryInfo(sql, paramList...)
	if err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, res)
	return nil
}

//出库时间，支付时间的判断
func getWhere(cus *customizeParam) (string, []interface{}) {
	outWhere, outList := setTimeToWhere("out_time", cus.OutTime.StartTime, cus.OutTime.StopTime)
	isPayWhere := ""
	switch cus.IsPay {
	case 0:
		//可以设置条件区间，或者没有条件
	case 1:
		//已支付
		isPayWhere = "pay_time is not null"
	case 2:
		//未支付
		isPayWhere = "pay_time is null"
	}
	if cus.IsPay != 0 {
		return outWhere + " and " + isPayWhere, outList
	}
	payWhere, payList := setTimeToWhere("pay_time", cus.PayTime.StartTime, cus.PayTime.StopTime)
	if outWhere != "" && payWhere != "" {
		return outWhere + " and " + payWhere, append(outList, payList...)
	}
	return outWhere + payWhere, append(outList, payList...)
}

//将时间参数转化为条件
func setTimeToWhere(column string, start, stop time.Time) (string, []interface{}) {
	switch {
	case !(start.IsZero() || stop.IsZero()): //都存在
		return column + " between ? and ?", []interface{}{start, stop}
	case start.IsZero() && stop.IsZero():
		return "", []interface{}{}
	default:
		if start.IsZero() {
			return column + " <= ?", []interface{}{stop}
		}
		return column + " >= ?", []interface{}{start}
	}
}

//getSummaryInfo sql执行，返回数据组和err
func getSummaryInfo(sql string, paramList ...interface{}) ([]summaryInfo, error) {
	rows, err := core.Ddb.Query(sql, paramList...)
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
