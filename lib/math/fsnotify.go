package math

import (
	"log"

	"github.com/fsnotify/fsnotify"
)

func fsnFunotify() {
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Close()
	err = watch.Add("./tmp")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			select {
			case ev := <-watch.Events:
				{
					//这是位的开关操作，Op的每一位（个位，十位，百位。。。都代表了一个开关，&（与）上自己，表示是否这一位是否是开启）
					//注意的是，删除文件，触发的也是重命名，因为不是直接删除，涉及系统删除文件，详细查资料
					if ev.Op&fsnotify.Create == fsnotify.Create {
						log.Println("创建文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Write == fsnotify.Write {
						log.Println("写入文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Remove == fsnotify.Remove {
						log.Println("删除文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Rename == fsnotify.Rename {
						log.Println("重命名文件 : ", ev.Name)
					}
					if ev.Op&fsnotify.Chmod == fsnotify.Chmod {
						log.Println("修改权限 : ", ev.Name)
					}
				}
			case err := <-watch.Errors:
				{
					log.Println("error : ", err)
					return
				}
			}
		}
	}()

	//循环
	select {}
}
