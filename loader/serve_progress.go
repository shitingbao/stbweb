package loader

import (
	"log"
	"net/http"
	"stbweb/core"
)

type appRest struct{}

func (ap *appRest) Get(arge *core.ElementHandleArgs) {
	log.Println("this is apprest")
	arge.Res.Write([]byte("this is apprest"))
	return
}

func httpProcess(w http.ResponseWriter, r *http.Request) {

	control := r.Header.Get("dbweb-api")
	log.Println("control:", control)
	if r.Method == "GET" && control == "rest" {
		arge := core.NewElementHandleArgs(w, r)
		var ats interface{}
		ats = new(appRest) //这里一定得使用new新建对象

		f, ok := ats.(core.Element)
		log.Println("ok-tt:", ok)
		if !ok {
			log.Println("no ok")
		}
		f.Get(arge)
	} else {
		log.Println("nothing")
	}

}
