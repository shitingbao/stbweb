package common

import (
	"net/http"
	"stbweb/core"
	"time"
)

//商品信息操作
type commodity struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	Price       float64   `json:"price"`
	CreateTime  time.Time `json:"create_time"`
}

func init() {
	core.RegisterFun("commodity", new(commodity), false)
}

func (c *commodity) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("del", nil, delCommodity) {
		return
	}
}

func delCommodity(pa interface{}, p *core.ElementHandleArgs) error {
	codeid := p.Req.URL.Query()["id"]
	//默认会有一个空字符串，数组索引0不会为nil log.Println("codeid:", codeid[0])
	if len(codeid) == 0 {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": "id can not null"})
		return nil
	}
	stmt, err := core.Ddb.Prepare(`delete from commodity where id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(codeid[0]); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}

func (c *commodity) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("info", new(commodityWhere), getCommodity) ||
		p.APIInterceptionPost("add", new(commodity), addCommodity) ||
		p.APIInterceptionPost("update", new(commodity), updateCommodity) {
		return
	}
}

//搜索条件待定
type commodityWhere struct {
	Page  int64 `json:"page"`
	Limit int64 `json:"limit"`
}

func getCommodity(pa interface{}, p *core.ElementHandleArgs) error {
	param := pa.(*commodityWhere)
	startLimit := (param.Page - 1) * param.Limit
	var results []commodity
	sql := `SELECT id,name,category,description,image,price,create_time FROM commodity limit ?,?`
	rows, err := core.Ddb.Query(sql, startLimit, param.Limit)
	if err != nil {
		return err
	}
	for rows.Next() {
		var result commodity
		rows.Scan(&result.ID, &result.Name, &result.Category, &result.Description, &result.Image, &result.Price, &result.CreateTime)
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err})
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "data": results})
	return nil
}

func updateCommodity(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*commodity)
	stmt, err := core.Ddb.Prepare(`UPDATE commodity SET name=?,category=?,description=?,image=?,price=? WHERE id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.Name, pa.Category, pa.Description, pa.Image, pa.Price, pa.ID); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}

func addCommodity(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*commodity)
	stmt, err := core.Ddb.Prepare(`INSERT INTO commodity(name,category,description,price,create_time) VALUES(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.Name, pa.Category, pa.Description, pa.Price, time.Now()); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}
