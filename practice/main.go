package main

import "log"

func main() {
	sli := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
	log.Println(getMaxArea(sli))
}

//ali 求水滴数量
func getWaterCount() {
	sli := []int{1, 0, 2, 0, 0, 1, 3, 0, 1}
	log.Println(getResCount(sli))
}
