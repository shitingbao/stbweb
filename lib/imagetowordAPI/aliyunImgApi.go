package imagetowordapi

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"stbweb/lib/images"
)

var appcode = "13c126d6576b4209a8c781db06c0c1c0"                    //appcode,接口身份唯一标识
var host = "https://ocrapi-document.taobao.com/ocrservice/document" //提取接口地址

type photoBodyimage struct {
	Img string
}

//imageObject post中json中的格式
// type imageObject struct {
// 	Image string
// }

//postPhoto 访问aliyun接口，提取图片中的文字
func postPhoto(url string) (string, error) {
	image64, err := images.ImageToBase64(url)
	if err != nil {
		return "", err
	}
	bodyData := photoBodyimage{}
	bodyData.Img = image64
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return "", err
	}
	sendBody := bytes.NewReader(jsonBody)
	client := &http.Client{}
	req, err := http.NewRequest("POST", host, sendBody)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "APPCODE "+appcode)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
