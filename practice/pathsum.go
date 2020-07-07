package main

import "log"

//TreeNode tree
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

//保存有左右节点的节点，每次先读左
type keepNode struct {
	Node    *TreeNode //保存有兄弟的节点
	NodeSum int       //根到该点的sum
	IsRead  int       //是否左右都已经读取,0star,1读左，2读右
}

//pass Wide-order first traversal
func wideOrder(root *TreeNode, sum int) bool {
	if root == nil {
		return false
	}
	stack := []keepNode{}
	logoNode := keepNode{
		Node:    root,
		NodeSum: root.Val,
	}
	stack = append(stack, logoNode)
	log.Println(len(stack))
	for len(stack) == 0 {
		log.Println(stack[0].NodeSum)
		log.Println(stack[0].Node.Val)
		if stack[0].Node.Left == nil && stack[0].Node.Right == nil {
			if stack[0].NodeSum == sum {
				return true
			}
			continue
		}
		if stack[0].Node.Left != nil {
			stack = append(stack, keepNode{
				Node:    stack[0].Node.Left,
				NodeSum: stack[0].NodeSum + stack[0].Node.Left.Val,
			})
		}
		if stack[0].Node.Right != nil {
			stack = append(stack, keepNode{
				Node:    stack[0].Node.Right,
				NodeSum: stack[0].NodeSum + stack[0].Node.Right.Val,
			})
		}
		stack = stack[1:]
	}
	return false
}

func hasPathSum(root *TreeNode, sum int) bool {
	if root == nil {
		return false
	}
	stack := []keepNode{}
	logoNode := keepNode{
		Node:    root,
		NodeSum: root.Val,
		IsRead:  2,
	}
	stack = append(stack, logoNode)
	for {
		// log.Println("sum:", logoNode.NodeSum)
		// log.Println("val:", logoNode.Node.Val)
		if !(logoNode.Node.Left == nil && logoNode.Node.Right == nil) && (logoNode.NodeSum*logoNode.NodeSum >= sum*sum) {
			if len(stack)-1 == 0 {
				return false
			}
			logoNode = stack[len(stack)-1] //回溯到上一个兄弟节点,并且要删掉
			stack = stack[0 : len(stack)-1]
			continue
		}
		if logoNode.Node.Left == nil && logoNode.Node.Right == nil { //判断叶子节点
			if logoNode.NodeSum == sum {
				return true
			}
			if len(stack)-1 == 0 {
				return false
			}
			logoNode = stack[len(stack)-1] //回溯到上一个兄弟节点,并且要删掉
			stack = stack[0 : len(stack)-1]
			continue
		}
		if logoNode.Node.Left != nil && logoNode.Node.Right != nil {
			if logoNode.IsRead == 2 {
				logoNode = keepNode{
					Node:    logoNode.Node.Right,
					NodeSum: logoNode.NodeSum + logoNode.Node.Right.Val,
					IsRead:  1,
				}
				continue
			}
			logoNode.IsRead = 2
			stack = append(stack, logoNode)
			logoNode = keepNode{
				Node:    logoNode.Node.Left,
				NodeSum: logoNode.NodeSum + logoNode.Node.Left.Val,
				IsRead:  1,
			}
			continue
		}
		if logoNode.Node.Left != nil {
			logoNode = keepNode{
				Node:    logoNode.Node.Left,
				NodeSum: logoNode.NodeSum + logoNode.Node.Left.Val,
				IsRead:  1,
			}
			continue
		}
		if logoNode.Node.Right != nil {
			stack = append(stack, logoNode)
			logoNode = keepNode{
				Node:    logoNode.Node.Right,
				NodeSum: logoNode.NodeSum + logoNode.Node.Right.Val,
				IsRead:  1,
			}
		}
		if len(stack) == 0 {
			break
		}
	}
	return false
}

func treeLoad() {
	root := getRoot()
	logo := wideOrder(root, 22)
	log.Println(logo)
}
func getRoot() *TreeNode {
	t9 := &TreeNode{
		Val:   1,
		Left:  nil,
		Right: nil,
	}
	t8 := &TreeNode{
		Val:   2,
		Left:  nil,
		Right: nil,
	}
	t7 := &TreeNode{
		Val:   7,
		Left:  nil,
		Right: nil,
	}
	t6 := &TreeNode{
		Val:   4,
		Left:  nil,
		Right: t9,
	}
	t5 := &TreeNode{
		Val:   13,
		Left:  nil,
		Right: nil,
	}
	t4 := &TreeNode{
		Val:   11,
		Left:  t7,
		Right: t8,
	}
	t3 := &TreeNode{
		Val:   8,
		Left:  t5,
		Right: t6,
	}
	t2 := &TreeNode{
		Val:   4,
		Left:  t4,
		Right: nil,
	}
	t1 := &TreeNode{
		Val:   5,
		Left:  t2,
		Right: t3,
	}
	return t1
}
