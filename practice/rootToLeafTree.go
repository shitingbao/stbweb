package main

import "log"

//TreeNode node
type rootToLeafTreeNode struct {
	Val   int
	Left  *rootToLeafTreeNode
	Right *rootToLeafTreeNode
}

func sumNumbersRootToLeaf(root *rootToLeafTreeNode) int {
	if root == nil {
		return 0
	}
	str := getAllLeafRootToLeaf(0, root)
	sum := 0
	for _, v := range str {
		sum += v
	}
	return sum
}

func getAllLeafRootToLeaf(sum int, root *rootToLeafTreeNode) []int {
	str := []int{}
	if root.Left == nil && root.Right == nil {
		str = append(str, sum*10+root.Val)
	}
	if root.Left != nil {
		str = append(str, getAllLeafRootToLeaf(sum*10+root.Val, root.Left)...)
	}
	if root.Right != nil {
		str = append(str, getAllLeafRootToLeaf(sum*10+root.Val, root.Right)...)
	}
	return str
}

func rootToLeafSumload() {
	t5 := &rootToLeafTreeNode{
		Val:   1,
		Left:  nil,
		Right: nil,
	}

	t4 := &rootToLeafTreeNode{
		Val:   5,
		Left:  nil,
		Right: nil,
	}

	t3 := &rootToLeafTreeNode{
		Val:   0,
		Left:  nil,
		Right: nil,
	}

	t2 := &rootToLeafTreeNode{
		Val:   9,
		Left:  t4,
		Right: t5,
	}

	t1 := &rootToLeafTreeNode{
		Val:   4,
		Left:  t2,
		Right: t3,
	}
	sum := sumNumbersRootToLeaf(t1)
	log.Println(sum)
}
