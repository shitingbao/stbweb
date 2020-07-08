package main

func divingBoard(shorter int, longer int, k int) []int {
	if k == 0 {
		return []int{}
	}
	if shorter == longer {
		return []int{shorter * k}
	}
	res := make([]int, k+1)
	for i := 0; i <= k; i++ {
		long := i*longer + (k-i)*shorter
		res[i] = long
	}
	return res
}
