package common

import (
	"fmt"
	"net/http"
	"stbweb/core"
	"time"
)

//回复信息
//每次只查询一部分，不管一级回复还是二级回复，前端点击展开，再次查询下一步的分页信息
type reply struct {
	ID          int64 `json:"id"`
	CommodityID int64 `json:"commodity_id"`
	ParentID    int64 `json:"parent_id"`
	LikeNumber  int64 `json:"like_number"`
	Common      int64 `json:"common"`
	CreateTime  int64 `json:"create_time"`
	UserID      int64 `json:"user_id"`
}

func init() {
	core.RegisterFun("reply", new(reply), false)
}

func (c *reply) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("del", nil, delReply) {
		return
	}
}

func delReply(pa interface{}, p *core.ElementHandleArgs) error {
	codeid := p.Req.URL.Query()["id"]
	//默认会有一个空字符串，数组索引0不会为nil log.Println("codeid:", codeid[0])
	if len(codeid) == 0 {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"msg": "id can not null"})
		return nil
	}
	stmt, err := core.Ddb.Prepare(`delete from reply where id=?`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(codeid[0]); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}

func (c *reply) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("info", new(replyWhere), getReply) ||
		p.APIInterceptionPost("add", new(reply), addReply) {
		return
	}
}

//搜索条件待定
type replyWhere struct {
	Limit    int64 `json:"limit"`
	FlagID   int64 `json:"flag_id"`   //commodity_id,或者是parent_id,根据is_parent，true为commodity_id
	IsParent bool  `json:"is_parent"` //是否是一级回复
}

func getReply(pa interface{}, p *core.ElementHandleArgs) error {
	param := pa.(*replyWhere)
	var results []reply

	sql := `SELECT id,commodity_id,parent_id,like_number,common,create_time,user_id FROM reply where %s=? order by create_time desc limit ?,?`
	flagWhere := "commodity_id"
	if !param.IsParent {
		flagWhere = "parent_id"
	}
	//因为是评论，默认从头拿，所以是1开始
	rows, err := core.Ddb.Query(fmt.Sprintf(sql, flagWhere), 1, param.Limit)
	if err != nil {
		return err
	}
	for rows.Next() {
		var result reply
		rows.Scan(&result.ID, &result.CommodityID, &result.ParentID, &result.LikeNumber, &result.Common, &result.CreateTime, &result.UserID)
		results = append(results, result)
	}
	if err := rows.Err(); err != nil {
		core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": false, "msg": err})
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true, "data": results})
	return nil
}

func addReply(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*reply)
	stmt, err := core.Ddb.Prepare(`INSERT INTO reply(commodity_id,parent_id,common,create_time,user_id) VALUES(?,?,?,?,?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.CommodityID, pa.ParentID, pa.Common, time.Now(), pa.UserID); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}
