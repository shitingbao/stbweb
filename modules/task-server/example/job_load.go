package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"stbweb/lib/task"
	"time"
)

func main() {
	// defer c.Stop()
	taskHandle := task.NewTaskHandle("site", "@every 10s", test)
	taskHandle2 := task.NewTaskHandle("site1", "@every 5s", test1)
	t := task.NewTaskElement("iflow", "job", "iflow", 1.10, taskHandle, taskHandle2)
	defer t.Close()
	if err := task.Watch(t); err != nil {
		log.Println("Watch:", err)
		return
	}

	lend := make(chan bool)
	sign := make(chan os.Signal, 1)
	signal.Notify(sign, os.Interrupt)
	go func() {
		for range sign {
			lend <- true
			break
		}
	}()
	<-lend
}

func test1() error {
	log.Println("iflow:", time.Now())
	return nil
}
func test() error {
	return errors.New("this is a error")
}
