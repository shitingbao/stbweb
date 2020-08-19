package oss

import (
	"fmt"
	"os"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

//管理地址https://oss.console.aliyun.com/bucket/oss-cn-beijing/shitingbao/overview
func ossUp() {
	// 创建OSSClient实例。基本信息
	client, err := oss.New("oss-cn-beijing.aliyuncs.com", "LTAI4G799EsiusjQBuXMGjBR", "eb3i3VkTqhjTa5rcDlTXRGUrcOwkit")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 获取存储空间。Bucket名称
	bucket, err := client.Bucket("shitingbao")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}

	// 上传本地文件。上传到服务器使用的文件名和本地文件地址，比如这里本地上传aa图片，服务器上名称就是bb
	err = bucket.PutObjectFromFile("bb", "./aa.jpg")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}
