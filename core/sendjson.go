package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	// ContentType 类型
	ContentType = "Content-Type"
	//ContentJSON 反馈json
	ContentJSON = "application/json"
	//defaultCharset 默认反馈编码格式
	defaultCharset = "UTF-8"
)

//SendMap send type
type SendMap map[string]interface{}

//SendJSON 将数据传递到json转码，并传到前端,这里需要重新设置header，不然不被允许反馈数据，就是前端接收不到write中的数据，会引起AllowCORS的问题
func SendJSON(w http.ResponseWriter, statuscode int, data interface{}) {

	bt, err := json.Marshal(data)
	if err != nil {
		return
	}

	w.Header().Set(ContentType, ContentJSON+";"+defaultCharset)
	if WebConfig.AllowCORS {
		allowOrigin := WebConfig.AllowOrigin
		if len(allowOrigin) == 0 {
			allowOrigin = "*" //待定，跨域允许的指定地址
		}
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin) //设置允许跨域的请求地址
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", fmt.Sprintf(
			"%s,Origin, X-Requested-With, Content-Type, Accept, Connection, User-Agent, Cookie",
			WebAPIHanderName)) //这里可以增加对应handle
	}

	w.WriteHeader(statuscode)
	w.Write(bt)
}
