package common

import (
	"net/http"
	"stbweb/core"
)

type login struct{}

type loginParam struct {
	Name string
	Pwd  string
}

func init() {
	core.RegisterFun("login", new(login))
}
func (l *login) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("login", new(loginParam), loginAPI) {
		return
	}

}

func (l *login) Get(p *core.ElementHandleArgs) {

}

type user struct {
	Name string
	Pwd  string
}

func loginAPI(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*loginParam)
	core.LOG.Info("param:", pa.Name, pa.Pwd)
	u := user{}
	// if err := core.Ddb.QueryRow("select * from user where id='1'").Scan(&u).Error(); err != "" {
	// 	return errors.New(err)
	// }
	// core.LOG.Info("user:", pa.Name, pa.Pwd)
	core.Ddb.QueryRow("select name,password from user where id='1'").Scan(&u.Name, &u.Pwd)
	core.SendJSON(p.Res, http.StatusOK, pa.Name+":"+pa.Pwd)
	return nil
}
