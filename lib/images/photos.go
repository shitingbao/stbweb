package images

import (
	"io"
	"mime/multipart"
	"os"
	"path"
	"stbweb/core"

	"github.com/pborman/uuid"
)

var (
	//文件保存的位置,使用默认路径
	fileDir = core.DefaultFilePath
)

//ByteToImage 转为图片暂存
//返回生成图片的路径和error
func ByteToImage(file multipart.File) (string, error) {

	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return "", err
	}
	fileAdree := path.Join(fileDir, uuid.NewUUID().String()+".jpeg")
	f, err := os.OpenFile(fileAdree, os.O_WRONLY|os.O_CREATE, 0777) //等待拆分
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}
	return fileAdree, nil
}
