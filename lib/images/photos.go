package images

import (
	"io"
	"linux_test_golang/core"
	"mime/multipart"
	"os"
)

//ByteToImage 转为图片暂存
//返回生成图片的路径和error
func ByteToImage(file multipart.File) (string, error) {

	fileAdree := "./file/" + core.GetUniqueFileName() + ".jpeg"
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
