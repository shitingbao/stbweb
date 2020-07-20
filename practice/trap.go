package main

import "log"

//求水滴，分层计算，首先计算出最高高度，用于循环次数
//计算首位和末尾的有效位置（就是有墙的地方）
//从原始列表截取这一段数组，循环，小于该层数的，就是凹点，就是有效值，计入总和
func trap(height []int) int {
	sum := 0
	dataMax := getTrapMax(height)
	for i := 1; i <= dataMax; i++ {
		max := getTrapDataMax(height, i)
		min := getTrapDataMin(height, i)
		for _, val := range height[min : max+1] {
			if val < i {
				sum++
			}
		}
	}
	return sum
}
func getTrapMax(list []int) int {
	max := 0
	for _, v := range list {
		if v > max {
			max = v
		}
	}
	return max
}
func getTrapDataMin(list []int, lg int) int {
	logo := 0
	for i := 0; i < len(list); i++ {
		if list[i] >= lg {
			logo = i
			break
		}
	}
	return logo
}
func getTrapDataMax(list []int, lg int) int {
	logo := 0
	for i := len(list) - 1; i > 0; i-- {
		if list[i] >= lg {
			logo = i
			break
		}
	}
	return logo
}
func trapLoad() {
	list := []int{0, 1, 0, 2, 1, 0, 1, 3, 2, 1, 2, 1}
	da := trap(list)
	log.Println(da)
}
