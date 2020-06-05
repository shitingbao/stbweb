package main

import "log"

type vmap struct {
	Key int
	Val int
}

func topKFrequent(nums []int, k int) []int {
	res := []vmap{}
	resval := make(map[int]int)
	for _, v := range nums {
		resval[v]++
	}
	for k, v := range resval {
		tp := vmap{
			Key: k,
			Val: v,
		}
		res = intomin(res, tp)
	}
	result := []int{}
	for _, v := range res {
		log.Println(v)
		result = append(result, v.Key)
	}
	return result[:k]
}

func intomin(nums []vmap, k vmap) []vmap {
	res := []vmap{}
	isMax := true
	for i, v := range nums {
		if k.Val > v.Val {
			isMax = false
			tres := []vmap{k}
			tres = append(tres, nums[i:]...)
			res = append(nums[:i], tres...)
			break
		}
	}
	if isMax {
		res = append(nums, k)
	}
	return res
}
