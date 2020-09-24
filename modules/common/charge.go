package common

import (
	"errors"
	"net/http"
	"stbweb/core"

	"github.com/sirupsen/logrus"
)

type charge struct {
	User     string
	Scroe    string
	PaidTime string
}

func init() {
	core.RegisterFun("charge", new(charge), true)
}

func (c *charge) Get() {}
func (c *charge) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("charge", new(charge), recharge) {
		return
	}

}

//大量操作使用协程池处理，使用全局workpool处理
func recharge(param interface{}, p *core.ElementHandleArgs) error {
	pm := param.(*charge)
	if p.Usr == "" {
		return errors.New("no login")
	}
	if core.WorkPool == nil {
		return errors.New("workpool nil")
	}

	if err := core.WorkPool.Submit(func() {
		stmt, err := core.Ddb.Prepare("INSERT INTO recharge(user,scroe,paid_time) VALUES(?,?,?)")
		if err != nil {
			logrus.WithFields(logrus.Fields{"INSERT": err.Error(), "user": p.Usr}).Error("recharge")
			return
		}
		if _, err := stmt.Exec(p.Usr, pm.Scroe, pm.PaidTime); err != nil {
			logrus.WithFields(logrus.Fields{"Exec:": err.Error(), "user": p.Usr}).Error("recharge")
			return
		}
	}); err != nil {
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}
