package core

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//SendMap send type
type SendMap map[string]interface{}

//SendJSON 将数据传递到json转码，并传到前端
func SendJSON(w http.ResponseWriter, statuscode int, data interface{}) {

	bt, err := json.Marshal(data)
	if err != nil {
		return
	}
	w.WriteHeader(statuscode)
	fmt.Fprintf(w, string(bt))
}
