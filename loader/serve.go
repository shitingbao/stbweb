package loader

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"stbweb/core"
	"stbweb/lib/formopera"
	"stbweb/lib/images"
	imagetowordapi "stbweb/lib/imagetowordAPI"
)

//AutoLoader 启动项
func AutoLoader() {
	serve()

	lend := make(chan bool)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			log.Println("received ctrl+c,wait back job finished...")
			core.TaskWaitGroup.Wait()
			log.Println("all back job finished,now shutdown http server...")
			// Shutdown()
			log.Println("success shutdown")
			lend <- true
			break
		}
	}()
	<-lend
}

func serve() {
	go func() {
		log.Println(http.ListenAndServe(":8088", nil))
	}()
	chatHub, ctrlHub := initChatWebsocket()
	core.Initinal(chatHub, ctrlHub)
	// http.HandleFunc("/", httpProcess) //设置访问的路由
	http.Handle("/", http.HandlerFunc(httpProcess))
}

func loadering(w http.ResponseWriter, r *http.Request) {
	if core.WebConfig.AllowCORS {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Action, Module")
	}

	imgurl, err := formopera.GetFromOnceImage(r)
	if err != nil {
		core.SendJSON(w, core.SendMap{"err": err.Error()})
		return
	}
	log.Println("imgurl:", imgurl)
	base64, err := images.ImageToBase64(imgurl)
	if err != nil {
		core.SendJSON(w, core.SendMap{"err": err.Error()})
		return
	}
	imagesBase64 := []string{}
	imagesBase64 = append(imagesBase64, base64)
	res, err := imagetowordapi.GetImageWord(imagesBase64)
	if err != nil {
		core.SendJSON(w, core.SendMap{"err": err.Error()})
		return
	}
	log.Println("res:", res)
	core.SendJSON(w, core.SendMap{"data": res})
}
