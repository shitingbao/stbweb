package common

import (
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
	if p.APIInterceptionGet("info", nil, summaryProtect) {
		return
	}
}

//上月汇总
func summaryProtect(param interface{}, p *core.ElementHandleArgs) error {
	year, month, _ := time.Now().Date()
	thisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	sql := "SELECT customer,SUM(price*num) as summary FROM order_info WHERE out_time BETWEEN ? AND ? GROUP BY customer"
	rows, err := core.Ddb.Query(sql, thisMonth.AddDate(0, -1, 0), thisMonth.AddDate(0, 0, -1))
	if err != nil {
		return err
	}
	var results []summaryInfo
	for rows.Next() {
		var res summaryInfo
		rows.Scan(&res.Customer, &res.Summary)
		results = append(results, res)
	}
	core.SendJSON(p.Res, http.StatusOK, results)
	return nil
}
