package main

import "log"

//动态规划基本,反馈多少种可能性，输入为一个矩形路线
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

//动态规划基本,反馈最短路线，输入为一个矩形路线
func minPathSum(grid [][]int) int {
	nums := []int{}
	for j := 0; j < len(grid[0]); j++ {
		nums = append(nums, 0)
	}

	for i, v := range grid {
		for idx, val := range v {
			if idx == 0 {
				nums[0] += val
				continue
			}
			if i == 0 {
				nums[idx] = val + nums[idx-1]
				continue
			}
			if nums[idx]+val > nums[idx-1]+val {
				nums[idx] = nums[idx-1] + val
			} else {
				nums[idx] = nums[idx] + val
			}
		}
		log.Println(nums)
	}
	return nums[len(nums)-1]
}
