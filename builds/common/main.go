package main

import (
	"stbweb/loader"
	_ "stbweb/modules/business" //这引用不能省，因为默认没被引用的包内，init函数不会被执行，不执行会导致对应操作元素没注册引起controller为空或无法匹配
	_ "stbweb/modules/chatroom" //这引用不能省
	_ "stbweb/modules/common"   //这引用不能省

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	loader.AutoLoader()
}
