package common

import (
	"net/http"
	"stbweb/core"
)

type register struct{}

func init() {
	core.RegisterFun("register", new(register), false)
}

type registerParam struct {
	Name     string
	Email    string
	Phone    string
	Password string
	Avatar   string
}

func (r *register) Post(p *core.ElementHandleArgs) {
	if p.APIInterceptionPost("register", new(registerParam), userRegister) {
		return
	}
}

func userRegister(param interface{}, p *core.ElementHandleArgs) error {
	pa := param.(*registerParam)
	if pa.Name == "" || pa.Password == "" {
		core.SendJSON(p.Res, http.StatusOK, "必填内容不能为空")
		return nil
	}
	if core.IsExistUser(pa.Name) {
		core.SendJSON(p.Res, http.StatusOK, "用户已存在")
		return nil
	}

	//两次加密一次解密，双向加单向
	salt := core.BuildIserSalt(pa.Name)
	bPwd := core.BuildPas(pa.Password, salt)
	stmt, err := core.Ddb.Prepare("INSERT INTO user(name,password,avatar,email,phone,salt)VALUES(?,?,?,?,?,?)")
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, err.Error())
		return err
	}
	_, err = stmt.Exec(pa.Name, bPwd, pa.Avatar, pa.Email, pa.Phone, salt)
	if err != nil {
		core.SendJSON(p.Res, http.StatusOK, err.Error())
		return err
	}
	core.SendJSON(p.Res, http.StatusOK, core.SendMap{"success": true})
	return nil
}
