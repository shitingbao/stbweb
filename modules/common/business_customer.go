package common

import (
	"net/http"
	"stbweb/core"
)

//顾客信息操作
type customer struct {
	ID     string
	Name   string
	Phone  string
	Adress string
}

func init() {
	core.RegisterFun("customer", new(customer), false)
}

func (c customer) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("info", nil, getCustomer) ||
		p.APIInterceptionGet("del", nil, delCustomer) {
		return
	}
}
func getCustomer(pa interface{}, p *core.ElementHandleArgs) error {
	var results []customer
	sql := `SELECT id,name,phone,adress FROM customer`
	rows, err := core.Ddb.Query(sql)
	if err != nil {
		return err
	}
	for rows.Next() {
		var result customer
		rows.Scan(&result.ID, &result.Name, &result.Phone, &result.Adress)
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err})
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "data": results})
	return nil
}

func delCustomer(pa interface{}, p *core.ElementHandleArgs) error {
	codeid := p.Req.URL.Query()["id"]
	//默认会有一个空字符串，数组索引0不会为nil log.Println("codeid:", codeid[0])

	if len(codeid[0]) == 0 {
		core.SendJSON(p.Res, http.StatusOK, map[string]string{"msg": "id can not null"})
		return nil
	}
	stmt, err := core.Ddb.Prepare(`delete from customer where id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(codeid[0]); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

func (c customer) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("add", new(customer), addCustomer) ||
		p.APIInterceptionPost("update", new(customer), updateCustomer) {
		return
	}
}

func updateCustomer(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*customer)
	stmt, err := core.Ddb.Prepare(`UPDATE customer SET name=?,phone=?,adress=? WHERE id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.Name, pa.Phone, pa.Adress, pa.ID); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

func addCustomer(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*customer)
	stmt, err := core.Ddb.Prepare(`INSERT INTO customer(name,phone,adress) VALUES(?,?,?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.Name, pa.Phone, pa.Adress); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}
