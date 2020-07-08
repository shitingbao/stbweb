package main

//ListNode node
type ListNode struct {
	Val  int
	Lo   int //这个字段用来测试输出
	Next *ListNode
}

//同步将一个有序森林变成有序树
func mergeKLists(lists []*ListNode) *ListNode {
	if len(lists) == 0 {
		return nil
	}
	node := lists[0]
	for i, v := range lists {
		if i == 0 {
			continue
		}
		node = sTof(node, v)
	}
	return node
}

//两条链表合并合并一条,输入两条有序链表，反馈一条整合的有序链表
func sTof(tnode, snode *ListNode) *ListNode {
	if tnode == nil {
		return snode
	}
	if snode == nil {
		return tnode
	}
	var first, seconde, node *ListNode
	if tnode.Val <= snode.Val {
		first = tnode
		seconde = snode
	} else {
		first = snode
		seconde = tnode
	}
	node = first
	for first.Next != nil {
		if seconde == nil {
			break
		}
		if seconde.Val > first.Next.Val {
			first = first.Next
			continue
		}
		//这里不能直接t=second赋值，因为是指针。直接地址传递了，改变t的时候连同seconde一起改了
		t := &ListNode{
			Val:  seconde.Val,
			Lo:   seconde.Lo,
			Next: seconde.Next,
		}
		t.Next = first.Next
		first.Next = t
		seconde = seconde.Next
		first = first.Next
	}
	if seconde != nil { //可能第二条链值比第一条最大的大，第一条链结束后，要把第二条没加上的连接到第一条末尾，比如135，24689，689就比起前面都大，就接到最后
		first.Next = seconde
	}
	return node
}
