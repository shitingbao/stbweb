package main

var (
	num = 0
)

//pass
func uniquePathsWithObstacles(obstacleGrid [][]int) int {
	// m := len(obstacleGrid) - 1
	n := len(obstacleGrid[0])
	line := []int{}
	for i := 0; i < n; i++ {
		line = append(line, 0)
	}
	for i, v := range obstacleGrid {
		for j, val := range v {
			if val == 1 {
				line[j] = 0
				continue
			}
			if j == 0 && i == 0 {
				line[j] = 1
				continue
			}
			if j == 0 {
				continue
			}
			line[j] += line[j-1]

		}
	}
	return line[len(line)-1]
}
func foreachList(obstacleGrid [][]int) {
	// m := len(obstacleGrid) - 1
	n := len(obstacleGrid[0]) - 1
	line := []int{}
	for i := 0; i < n; i++ {
		line = append(line, 0)
	}
	for _, v := range obstacleGrid {
		for j, val := range v {
			if val == 1 {
				line[j] = 0
				continue
			}
			if j == 0 {
				line[j] = 1
			} else {
				line[j] += line[j-1]
			}
		}
	}
}

func reChanLists(m, n int, obstacleGrid [][]int) int {
	if obstacleGrid[m][n] == 1 {
		return 0
	}
	if m == 0 && n == 0 {
		return 1
	}

	xNum := 0
	ynum := 0
	if m == 0 {
		xNum = 0
		// ynum = reChanLists(m, n-1, obstacleGrid)
	} else {
		xNum = reChanLists(m-1, n, obstacleGrid)
		// ynum = reChanLists(m, n-1, obstacleGrid)
	}
	if n == 0 {
		// xNum = reChanLists(m-1, n, obstacleGrid)
		ynum = 0
	} else {
		// xNum = reChanLists(m-1, n, obstacleGrid)
		ynum = reChanLists(m, n-1, obstacleGrid)
	}
	return xNum + ynum
}

func reLineList(m, n, objM, objN int, list *[]int, obstacleGrid [][]int) []int {
	// log.Println(m, n, objM, objN, list)
	if obstacleGrid[objM][objN] == 1 {
		return *list
	}
	if m == objM && n == objN {
		*list = append(*list, 0)
	}
	if objM < m {
		reLineList(m, n, objM+1, objN, list, obstacleGrid)
	}
	if objN < n {
		reLineList(m, n, objM, objN+1, list, obstacleGrid)
	}
	return *list
}

func reCalculation(m, n, objM, objN int, obstacleGrid [][]int) {
	// log.Println(m, n, objM, objN)
	if obstacleGrid[objM][objN] == 1 {
		return
	}
	if m == objM && n == objN {
		num++
	}
	if objM < m {
		reCalculation(m, n, objM+1, objN, obstacleGrid)
	}
	if objN < n {
		reCalculation(m, n, objM, objN+1, obstacleGrid)
	}
}
