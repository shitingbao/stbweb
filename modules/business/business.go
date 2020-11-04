package business

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"stbweb/core"
	"time"

	"github.com/pborman/uuid"
)

type business struct{}

type order struct {
	OrderID  string
	FoodName string
	Price    float64
	Number   float64
	Customer string
	OutTime  string
	PayTime  sql.NullString
}

func init() {
	core.RegisterFun("business", new(business), false)
}

func (bn *business) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("del", nil, delProtect) {
		return
	}
}

func delProtect(param interface{}, p *core.ElementHandleArgs) error {
	codeid := p.Req.URL.Query()["codeid"]
	//默认会有一个空字符串，数组索引0不会为nil log.Println("codeid:", codeid[0])

	if len(codeid[0]) == 0 {
		core.SendJSON(p.Res, http.StatusOK, map[string]string{"msg": "codeid can not null"})
		return nil
	}
	stmt, err := core.Ddb.Prepare(`delete from order_info where id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(codeid[0]); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

func (bn *business) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("add", new(order), addProduct) ||
		p.APIInterceptionPost("update", new(order), updateProduct) ||
		p.APIInterceptionPost("result", new(tab), resProtect) {
		return
	}
}

func updateProduct(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*order)
	log.Println(pa.Customer)
	stmt, err := core.Ddb.Prepare(`UPDATE order_info SET food_name=? WHERE id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.OrderID, pa.FoodName); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

func addProduct(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*order)
	stmt, err := core.Ddb.Prepare(`INSERT INTO order_info(id,food_name,price,num,customer,out_time) VALUES(?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(uuid.NewUUID().String(), pa.FoodName, pa.Price, pa.Number, pa.Customer, time.Now()); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

type tab struct {
	Page         int    //页码
	Max          int    //每页约束
	OrderBy      string //排序标识
	ConditionCol string //条件列
	Condition    string //条件
}

func resProtect(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*tab)
	rows, err := resProtectRows(pa)
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err})
		return err
	}
	var results []order
	for rows.Next() {
		var res order
		rows.Scan(&res.OrderID, &res.FoodName, &res.Price, &res.Number, &res.Customer, &res.OutTime, &res.PayTime)
		results = append(results, res)
	}
	if err := rows.Err(); err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err})
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "data": results})
	return nil
}

func resProtectRows(pa *tab) (*sql.Rows, error) {
	var rows *sql.Rows
	var err error
	sql := `SELECT id,food_name,price,num,customer,out_time,pay_time FROM order_info ORDER BY price asc LIMIT ?,?`
	switch pa.ConditionCol {
	case "all":
		rows, err = core.Ddb.Query(sql, (pa.Page-1)*pa.Max, pa.Max)
		if err != nil {
			return nil, err
		}
	case "":
		sql = `SELECT id,food_name,price,num,customer,out_time,pay_time FROM order_info 
					WHERE 
				LOCATE(?,id)>0 or 
				LOCATE(?,food_name)>0 or
				LOCATE(?,price)>0 or
				LOCATE(?,num)>0 or
				LOCATE(?,customer)>0 ORDER BY price asc LIMIT ?,?`
		rows, err = core.Ddb.Query(sql, pa.Condition, pa.Condition, pa.Condition, pa.Condition, pa.Condition, (pa.Page-1)*pa.Max, pa.Max)
		if err != nil {
			return nil, err
		}
	default:
		sql = fmt.Sprintf(`SELECT id,food_name,price,num,customer,out_time,pay_time FROM order_info 
					WHERE 
				LOCATE(?,%s)>0 ORDER BY price asc LIMIT ?,?`, pa.ConditionCol)
		rows, err = core.Ddb.Query(sql, pa.Condition, (pa.Page-1)*pa.Max, pa.Max)
		if err != nil {
			return nil, err
		}
	}
	return rows, err
}
