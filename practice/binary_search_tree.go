package main

//二叉排序树的构建
type searchTree struct {
	Value int
	Left  *searchTree
	Right *searchTree
}

// type searchTree struct {
// 	Value int
// 	Left  *searchTree
// 	Right *searchTree
// }
func createSearchTree(start, stop int) []*searchTree {
	allTree := []*searchTree{}
	if start == stop {
		node := &searchTree{start, nil, nil}
		allTree = append(allTree, node)
		return allTree
	}
	if start > stop {
		return []*searchTree{nil}
	}
	for i := start; i <= stop; i++ {
		// log.Println("start:", start, "--stop:", stop)
		leftTree := createSearchTree(start, i-1)
		rightTree := createSearchTree(i+1, stop)

		for _, left := range leftTree {
			for _, right := range rightTree {
				node := &searchTree{
					Value: i,
					Left:  left,
					Right: right,
				}
				// node.Left = left
				// node.Right = right
				allTree = append(allTree, node)
			}
		}
	}
	// log.Println("start:", start, "-res:", allTree)
	return allTree
}

func generateTrees(n int) []*searchTree {
	if n == 0 {
		return []*searchTree{}
	}
	return createSearchTree(1, n)
}
