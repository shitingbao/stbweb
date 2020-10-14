package images

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"
	"time"
)

//base64ToPngimage base64转为png图片
func base64ToPngimage(imagebase64 string) error {
	debytes, err := base64.StdEncoding.DecodeString(imagebase64)
	if err != nil {
		return err
	}
	bt := bytes.NewReader(debytes)

	image, err := png.Decode(bt)
	if err != nil {
		log.Println("png 编辑出错")
		return err
	}
	f, err := os.OpenFile("./file/"+getUniqueFileName()+".png", os.O_WRONLY|os.O_CREATE, 0777) //等待拆分
	if err != nil {
		return err
	}
	// f.Write(debytes)
	defer f.Close()
	png.Encode(f, image) //Options是编码参数，它的取值范围是1-100，值越高质量越好
	return nil
}

//base64ToJpgimage base64转为jpeg图片
func base64ToJpgimage(imagebase64 string) error {
	debytes, err := base64.StdEncoding.DecodeString(imagebase64) //这里需要注意，data:image/jpeg;base64的前缀需要先去掉
	if err != nil {
		return err
	}
	bt := bytes.NewReader(debytes)

	image, err := jpeg.Decode(bt)
	if err != nil {
		log.Println("jpeg 编辑出错")
		return err
	}
	f, err := os.OpenFile("./file/"+getUniqueFileName()+".jpeg", os.O_WRONLY|os.O_CREATE, 0777) //等待拆分
	if err != nil {
		return err
	}
	// f.Write(debytes)//不需要其他精度的话，直接这一句完成图片加载即可
	// defer f.Close()
	jpeg.Encode(f, image, &jpeg.Options{Quality: 100}) //Options是编码参数，它的取值范围是1-100，值越高质量越好
	return nil
}

//ImageToBase64 传入图片路径，转化为base64返回
func ImageToBase64(url string) (string, error) {
	imgFile, err := os.Open(url) // a QR code image
	if err != nil {
		return "", err
	}
	defer imgFile.Close()
	fInfo, _ := imgFile.Stat() //返回文件结构
	size := fInfo.Size()       //获取文件大小
	buf := make([]byte, size)  //根据大小分配一个byte数组
	fReader := bufio.NewReader(imgFile)
	fReader.Read(buf)
	//这里不加前缀，是为了将base64转回来方便，生成的时候没有，解析回图片的时候也不需要
	//"data:image/jpeg;base64,默认是不产生头标志的，如果需要自己添加
	imgBase64Str := base64.StdEncoding.EncodeToString(buf)

	return imgBase64Str, nil
}

//ImageHandtoBase64 传入一个*file对象，返回base64字符串
func ImageHandtoBase64(f *os.File) (string, error) {
	finfo, err := f.Stat()
	if err != nil {
		return "", err
	}
	buf := make([]byte, finfo.Size())
	fRead := bufio.NewReader(f)
	fRead.Read(buf)
	imageBase64 := base64.StdEncoding.EncodeToString(buf)
	return imageBase64, nil
}

//dataTofile 将数据写入记录文档中
func dataTofile(data string) error {
	f, err := os.Create("./assets/" + getUniqueFileName() + ".txt")
	if err != nil {
		return err
	}
	defer f.Close()
	f.Write([]byte(data))
	return nil
}

//getUniqueFileName 返回一个根据时间的唯一数字型字符串
func getUniqueFileName() string {
	name := time.Now().Format("2006-01-02 15:04:05")
	name = strings.Replace(name, "-", "", -1)
	name = strings.Replace(name, " ", "", -1)
	name = strings.Replace(name, ":", "", -1)
	return name
}
