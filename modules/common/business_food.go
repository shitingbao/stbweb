package common

import (
	"net/http"
	"stbweb/core"
)

type food struct {
	ID   string
	Name string
}

func init() {
	core.RegisterFun("food", new(food), false)
}

func (c food) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("info", nil, getFood) ||
		p.APIInterceptionGet("del", nil, delFood) {
		return
	}
}

func getFood(pa interface{}, p *core.ElementHandleArgs) error {
	var results []food
	sql := `SELECT id,name FROM food`
	rows, err := core.Ddb.Query(sql)
	if err != nil {
		return err
	}
	for rows.Next() {
		var result food
		rows.Scan(&result.ID, &result.Name)
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err})
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "data": results})
	return nil
}

func delFood(pa interface{}, p *core.ElementHandleArgs) error {
	codeid := p.Req.URL.Query()["id"]
	//默认会有一个空字符串，数组索引0不会为nil log.Println("codeid:", codeid[0])

	if len(codeid[0]) == 0 {
		core.SendJSON(p.Res, http.StatusOK, map[string]string{"msg": "id can not null"})
		return nil
	}
	stmt, err := core.Ddb.Prepare(`delete from food where id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(codeid[0]); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

func (c food) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("add", new(food), addFood) ||
		p.APIInterceptionPost("update", new(food), updateFood) {
		return
	}
}

func updateFood(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*food)
	stmt, err := core.Ddb.Prepare(`UPDATE food SET name=? WHERE id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.Name, pa.ID); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}

func addFood(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*food)
	stmt, err := core.Ddb.Prepare(`INSERT INTO food(name) VALUES(?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.Name); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, map[string]bool{"success": true})
	return nil
}
