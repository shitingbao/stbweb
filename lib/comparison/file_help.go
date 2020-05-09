//Package comparison 该包适用于txt以及csv文件，excel文件使用excel中的比对
package comparison

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	space = " "
	comma = ","
	other = ""
)

//逐行读取的三种方法
func readLineFile(fileName string) {
	if file, err := os.Open(fileName); err != nil {
		panic(err)
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			log.Println("NewScanner:", scanner.Text())
		}
	}
}

//如果有空行，这个方法会多一行，因为最后一行也可能有回车转义符
func readFileLine(fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		log.Println("n:", line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

func readLine(fileName string) {
	fi, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		log.Println("line:", string(a))
	}
}
