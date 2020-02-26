package formopera

import (
	"fmt"
	"linux_test_golang/lib/images"
	"net/http"
	"net/url"
)

//sendForm 发送表单数据,form内的数据，后台用key := r.PostFormValue("key")接收
func sendForm(address string) {
	client := &http.Client{}
	res, err := client.PostForm(address, url.Values{
		"key": []string{"values"},
	})
	if err != nil {
		return
	}
	defer res.Body.Close()
}

//getFormQuery 解析form中的值
func getFormQuery(r *http.Request) {
	//这里是接收query的值，需要使用ParseForm解析
	if err := r.ParseForm(); err != nil {
		return
	}
	for k, v := range r.Form {
		fmt.Printf("%s:%s\n", k, v)
	}
}

//getFormBodyVal 接收表单中body中的form值
func getFormBodyVal(r *http.Request) {
	//这里是接收body中form表单内的元素值，ParseMultipartForm需要先调用，用来分配给接收到文件的大小，不然会异常
	r.ParseMultipartForm(20 << 20)
	for k, v := range r.MultipartForm.Value { //获取表单字段
		fmt.Printf("%s:%s\n", k, v)
	}
}

//GetFromOnceImage 解析出表单内的图片
//单张图片内容解析
func GetFromOnceImage(r *http.Request) (string, error) {
	r.ParseMultipartForm(20 << 20)
	//也需要调用ParseMultipartForm
	file, _, err := r.FormFile("file")
	if err != nil {
		return "", err
	}

	defer file.Close()

	return images.ByteToImage(file)
}

//getFormEvenFile 便利获取所有文件内容
func getFormEvenFile(r *http.Request) {
	r.ParseMultipartForm(20 << 20)
	//获取表单中的文件
	//多个同时接受
	for _, v := range r.MultipartForm.File {
		for _, f := range v {
			fil, err := f.Open()
			if err != nil {
				return
			}
			defer fil.Close()
		}
	}
}
