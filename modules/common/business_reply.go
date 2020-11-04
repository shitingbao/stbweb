package common

import (
	"net/http"
	"stbweb/core"
	"time"
)

//评论，回复消息
//（参考抖音回复）
//它里面的评论结构，一个视频里面，首先是一级评论（根据点赞，时间，如果是本用户的评论，显示在最前面，排序）
//二级评论，在一级下面展开，二级评论里面的回复就是指定用户，不再向下划分层次了，排序同理
//他的层次其实就两层,下面都是i指定，就像@那种，谁对谁说,这样的话，在增加，删除都查一次就行,前端只要点击展开更多，根据这个条件，向下再取几条即可,然后加进去展示

//完整的回复数据结构
type replyResult struct {
	ID         int64     `json:"id"`
	LikeNumber int64     `json:"like_number"`
	CreateTime time.Time `json:"create_time"`
	Name       string    `json:"name"`
	AimName    string    `json:"aim_name"`
	reply
}

//部分数据结构，用于接受数据
type reply struct {
	CommodityID int64  `json:"commodity_id"`
	ParentID    int64  `json:"parent_id"`
	Common      string `json:"common"`
	UserID      int64  `json:"user_id"`
	AimsUserID  int64  `json:"aims_user_id"`
}

//搜索条件结构
type replyWhere struct {
	Limit       int64 `json:"limit"`
	CommodityID int64 `json:"commodity_id"`
	ParentID    int64 `json:"parent_id"`
}

func init() {
	core.RegisterFun("reply", new(replyResult), false)
}

func (c *replyResult) Get(p *core.ElementHandleArgs) {
	if p.APIInterceptionGet("del", nil, delReply) {
		return
	}
}

func delReply(pa interface{}, p *core.ElementHandleArgs) error {
	codeid := p.Req.URL.Query()["id"]
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

func (c *replyResult) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("info", new(replyWhere), getReply) ||
		p.APIInterceptionPost("add", new(reply), addReply) {
		return
	}
}

//获取一部分评论信息
func getReply(pa interface{}, p *core.ElementHandleArgs) error {
	param := pa.(*replyWhere)
	var results []replyResult
	sql := `SELECT id,commodity_id,parent_id,like_number,common,create_time,user_id,
				IFNULL(aims_user_id,0),
				IFNULL((SELECT NAME FROM user WHERE id=user_id),'') as NAME,
				IFNULL((SELECT NAME FROM user WHERE id=aims_user_id),'')	as aim_name  
			FROM reply 
			WHERE commodity_id=? AND parent_id=? 
			order BY like_number DESC,create_time asc 
			LIMIT 0,?`
	//因为是评论，默认从头拿，所以是limit 0开始
	rows, err := core.Ddb.Query(sql, param.CommodityID, param.ParentID, param.Limit)
	if err != nil {
		return err
	}
	for rows.Next() {
		var result replyResult
		if err := rows.Scan(&result.ID, &result.CommodityID, &result.ParentID, &result.LikeNumber, &result.Common, &result.CreateTime, &result.UserID, &result.AimsUserID, &result.Name, &result.AimName); err != nil {
			return err
		}
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
	stmt, err := core.Ddb.Prepare(`INSERT INTO reply(commodity_id,parent_id,common,create_time,user_id,aims_user_id) VALUES(?,?,?,?,?,?)`)
	if err != nil {
		return err
	}
	if _, err := stmt.Exec(pa.CommodityID, pa.ParentID, pa.Common, time.Now(), pa.UserID, pa.AimsUserID); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}
