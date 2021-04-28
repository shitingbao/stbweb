package spider

import (
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	img = "https://www.pximg.com/wp-content/uploads/2020/04/cdd72379d64dd08-scaled.jpg"
	url = "https://www.pximg.com/meinv/39211.html"
)

func laod() {
	n := getFileName(img)

	log.Println(strings.Replace(n, "-", "_", -1))
}
func createImage(imageURL string) error {
	resp, err := http.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return err
	}
	bt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f, err := os.Create("file/aa.jpg")
	if err != nil {
		return err
	}

	if _, err := f.Write(bt); err != nil {
		return err
	}
	return nil
}

func getFileName(u string) string {
	filenameWithSuffix := path.Base(u)                        //获取文件名带后缀(test.txt)
	fileSuffix := path.Ext(u)                                 //获取文件后缀(.txt)
	return strings.TrimSuffix(filenameWithSuffix, fileSuffix) //获取文件名(test)
}
