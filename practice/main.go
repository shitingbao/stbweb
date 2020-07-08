package main

import (
	"log"
)

func main() {

	t6 := &ListNode{
		Val: 2, Lo: 6,
		Next: nil,
	}
	t5 := &ListNode{
		Val: 2, Lo: 5,
		Next: t6,
	}
	t4 := &ListNode{
		Val: 2, Lo: 4,
		Next: t5,
	}

	t3 := &ListNode{
		Val: 1, Lo: 3,
		Next: nil,
	}
	t2 := &ListNode{
		Val: 1, Lo: 2,
		Next: t3,
	}
	t1 := &ListNode{
		Val: 1, Lo: 1,
		Next: t2,
	}
	node := sTof(t1, t4)
	for node != nil {
		log.Println(node.Val, ":", node.Lo)
		node = node.Next
	}
}
