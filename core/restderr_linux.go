package core

import "log"

// redirectStderr to the file passed in
//panic重定向至指定文件中
//该方法只适用于linux下编译，因为windows下没有Dup2这个，会编译报错
//注意这个文件操作如果是局部变量，在他被回收时，就无效了
func redirectStderr() {
	// err := syscall.Dup2(int(f.Fd()), int(os.Stderr.Fd()))
	// if err != nil {
	// 	log.Fatalf("Failed to redirect stderr to file: %v", err)
	// }
	log.Println("this is linux")
}
