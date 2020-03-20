package core

import (
	"net/http"

	"github.com/Sirupsen/logrus"
)

var (
	//WebAPIHanderName hander中存放api名称
	WebAPIHanderName = "dbweb-api"
)

//APIInterceptionGet 拦截api请求，统一处理
func (e *ElementHandleArgs) APIInterceptionGet(methodName string, param interface{},
	cb func(pa interface{}, content *ElementHandleArgs) error) bool {
	//名称对应判断
	if methodName != e.apiName() {
		return false
	}
	if err := cb(param, e); err != nil {
		LOG.WithFields(logrus.Fields{
			"elename": e.Element.Name,
			"method":  methodName,
		}).Error(err)
		SendJSON(e.Res, http.StatusBadRequest, SendMap{"msg": err})
	}
	return true
}

func (e *ElementHandleArgs) apiName() string {
	return e.Req.Header.Get(WebAPIHanderName)
}
