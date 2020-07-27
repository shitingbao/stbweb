package main

import "log"

//深度优先，加记忆存储（动态规划的思想，找出每一个点的极值时，记录，后面要用到该点时，不需要重新计算）
func longestIncreasingPath(matrix [][]int) int {
	height := len(matrix)
	if height == 0 {
		return 0
	}
	width := len(matrix[0])
	logo := [][]int{}
	for i := 0; i < height; i++ {
		ld := []int{}
		for j := 0; j < width; j++ {
			ld = append(ld, -1)
		}
		logo = append(logo, ld)
	}
	max := 0
	for idx, val := range matrix {
		for i := range val {
			lmax := search(matrix, logo, 1, idx, i)
			// logo[idx][i] = lmax
			if lmax > max {
				max = lmax
			}
		}
	}
	return max
}

func search(matrix, logo [][]int, max, x, y int) int {
	height := len(matrix)
	width := len(matrix[0])
	dtop, dbotton, dleft, dright := 0, 0, 0, 0
	t := 0
	if x > 0 && matrix[x-1][y] > matrix[x][y] { //上
		if logo[x-1][y] != -1 {
			dtop = logo[x-1][y]
		} else {
			dtop = search(matrix, logo, max, x-1, y)
		}
		if dtop > t {
			t = dtop
		}
	}
	if x < height-1 && matrix[x+1][y] > matrix[x][y] { //下
		if logo[x+1][y] != -1 {
			dbotton = logo[x+1][y]
		} else {
			dbotton = search(matrix, logo, max, x+1, y)
		}
		if dbotton > t {
			t = dbotton
		}
	}
	if y > 0 && matrix[x][y-1] > matrix[x][y] { //左
		if logo[x][y-1] != -1 {
			dleft = logo[x][y-1]
		} else {
			dleft = search(matrix, logo, max, x, y-1)
		}
		if dleft > t {
			t = dleft
		}
	}
	if y < width-1 && matrix[x][y+1] > matrix[x][y] { //右
		if logo[x][y+1] != -1 {
			dright = logo[x][y+1]
		} else {
			dright = search(matrix, logo, max, x, y+1)
		}
		if dright > t {
			t = dright
		}
	}
	logo[x][y] = max + t
	return max + t
}

//给定一个矩形空间，求递增最长路径
func deepLoad() {
	// list := [][]int{{9, 9, 4}, {6, 6, 8}, {2, 1, 1}}
	// list := [][]int{{3, 4, 5}, {3, 2, 6}, {2, 2, 1}}
	// list := [][]int{{1, 2}}
	list := [][]int{{7, 7, 5}, {2, 4, 6}, {8, 2, 0}}
	max := longestIncreasingPath(list)
	log.Println(max)
}
