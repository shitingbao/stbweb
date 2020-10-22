package imagetowordapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	accessTokenURL = "https://aip.baidubce.com/oauth/2.0/token?grant_type=client_credentials&client_id=WhGsmv5uTul6WUVdqmQjAbv3&client_secret=owaOpOjMUVt3zXIweepNQPIpgEDxSeTt"
	// access_token   = "24.89e545a55e7425d87864341b99429dd8.2592000.1581213789.282335-17903904"
)

//accessTokenType 获取access_token的类型，只有expires_in和access_token有用，其他可以忽略
type accessTokenType struct {
	RefreshToken  string `json:"refresh_token"`
	ExpiresIn     string `json:"expires_in"` //使用时间，秒
	SessionKey    string `json:"session_key"`
	AccessToken   string `json:"access_token"` //正式用来使用的校验值
	Scope         string `json:"scope"`
	SessionSecret string `json:"session_secret"`
}

//AcceptResultWord 接收识别的文字信息
type AcceptResultWord struct {
	LogID          int64   `json:"log_id"`
	WordsResultNum int64   `json:"words_result_num"`
	WordsResult    []words `json:"words_result"`
}
type words struct {
	Words string
}

//getAccessToken 获取accessToken,内部的三个参数,一个月的有效期
//grant_type： 必须参数，固定为client_credentials
///client_id： 必须参数，应用的API Key
//client_secret： 必须参数，应用的Secret Key
//https://ai.baidu.com/docs#/Auth/top
func getAccessToken() (accessTokenType, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", accessTokenURL, nil)
	if err != nil {
		return accessTokenType{}, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return accessTokenType{}, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return accessTokenType{}, err
	}

	accessToken := accessTokenType{}
	json.Unmarshal(b, &accessToken)
	return accessToken, nil
}

//imageObject post中json中的格式
type imageObject struct {
	Image string `json:"image"`
}

//getImageWord 发送表单数据,返回word，需要一个query值和表单中放入image64数据，大小不能超过4M,具体返回和参数去参考api文档
//https://ai.baidu.com/docs#/OCR-API-GeneralBasic/top
//注意：这里的数据不能使用上面那种imageObject形式在body中放json，只能用表单数据提交（下面这种），亲测无效
//每次一张图片
func getImageWord(address, accessToken, accessTokenDate string, imageBase64 []string) (AcceptResultWord, error) {
	// checkTokenEffect(accessTokenDate)
	client := &http.Client{}
	// res, err := client.PostForm(core.WebConfig.BaidubceAddress+"?access_token="+core.WebConfig.AccessToken, url.Values{
	res, err := client.PostForm(address+"?access_token="+accessToken, url.Values{
		"image": imageBase64,
	})
	if err != nil {
		return AcceptResultWord{}, err
	}
	defer res.Body.Close()
	resData := AcceptResultWord{}
	resdatabyte, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return AcceptResultWord{}, err
	}
	if err := json.Unmarshal(resdatabyte, &resData); err != nil {
		return AcceptResultWord{}, err
	}
	return resData, nil
}

//GetImageWord 输入images base64格式数组，反馈解析内容和err，反馈的解析内容中包含对应word数组
//address, accessToken 为core.WebConfig.BaidubceAddress  core.WebConfig.AccessToken  core.WebConfig.AccessTokenDate
func GetImageWord(address, accessToken, accessTokenDate string, imageBase64 []string) (AcceptResultWord, error) {
	return getImageWord(address, accessToken, accessTokenDate, imageBase64)
}

//judge30Date 日期是否在30天之内，是返回true
func judge30Date(date string) bool {
	historyTime, err := time.Parse("2006-01-02 15:04:05", date)
	if err != nil {
		logrus.WithFields(logrus.Fields{"AccessTokenDate": err.Error()}).Panic("AccessTokenDate")
		panic(err)
	}
	if time.Now().AddDate(0, 0, -30).After(historyTime) {
		return false
	}
	return true
}

//CheckTokenEffect 检查百度接口的token是否过期，如过期，反馈新token
func CheckTokenEffect(accessTokenDate string) (string, error) {
	if judge30Date(accessTokenDate) {
		return "", nil
	}
	logrus.WithFields(logrus.Fields{"baidu-word-token": "get new wordAPI token"}).Info("wordAPItoken")
	at, err := getAccessToken()
	if err != nil {
		logrus.WithFields(logrus.Fields{"baidu-word-token": err}).Error("get baidu api err") //出错就直接异常
		return "", err
	}
	return at.AccessToken, nil
}
