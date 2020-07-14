package main

//动态规划基本
func minimumTotal(triangle [][]int) int {
	for i, v := range triangle {
		if i == 0 {
			continue
		}
		for idx, val := range v {
			if idx == len(v)-1 {
				triangle[i][idx] = val + triangle[i-1][idx-1]
				continue
			}
			if idx == 0 {
				triangle[i][idx] = val + triangle[i-1][idx]
				continue
			}
			if val+triangle[i-1][idx] > val+triangle[i-1][idx-1] {
				triangle[i][idx] = val + triangle[i-1][idx-1]
			} else {
				triangle[i][idx] = val + triangle[i-1][idx]
			}
		}
	}
	min := triangle[len(triangle)-1][0]
	for _, v := range triangle[len(triangle)-1] {
		if v < min {
			min = v
		}
	}
	return min
}
