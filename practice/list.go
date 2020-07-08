package main

//ListNodeData listnode
type ListNodeData struct {
	Val  int
	Next *ListNodeData
}

func addTwoNumbersChan(l1 *ListNodeData, l2 *ListNodeData) *ListNodeData {
	nm1 := make(chan int)
	nm2 := make(chan int)
	go getNodeNum(l1, nm1)
	go getNodeNum(l2, nm2)
	result := &ListNodeData{}
	resKeyres := result //复制头指针
	resKey := true
	num1 := 0
	num2 := 0
	resCount := 0
	for {
		if num1 != -1 {
			num1 = <-nm1
		}

		if num2 != -1 {
			num2 = <-nm2
		}
		if num2 == -1 && num1 == -1 {
			if resCount != 0 { //判断一个5+5这种多一位的情况
				rt := &ListNodeData{
					Val: 1,
				}
				resKeyres.Next = rt
			}
			break
		}
		numCount1, numCount2 := num1, num2
		if numCount1 == -1 {
			numCount1 = 0
		}

		if numCount2 == -1 {
			numCount2 = 0
		}
		count := numCount1 + numCount2 + resCount
		resVal := &ListNodeData{
			Val: count % 10,
		}
		if resKey {
			result = resVal
			resKeyres = resVal
			resKey = false
		} else {
			resKeyres.Next = resVal
			resKeyres = resKeyres.Next
		}

		if count > 9 {
			resCount = 1
		} else {
			resCount = 0
		}
	}

	return result
}

func addTwoNumbers(l1 *ListNodeData, l2 *ListNodeData) *ListNodeData {
	oneList := []int{}
	twoList := []int{}
	nm1 := make(chan int)
	nm2 := make(chan int)
	go getNodeNum(l1, nm1)
	go getNodeNum(l2, nm2)
	result := &ListNodeData{}
	resKeyres := result
	num1 := 0
	num2 := 0
	for {
		if num1 != -1 {
			num1 = <-nm1
		}
		if num2 != -1 {
			num2 = <-nm2
		}
		if num2 == -1 && num1 == -1 {
			break
		}
		if num1 != -1 {
			oneList = append(oneList, num1)
		}
		if num2 != -1 {
			twoList = append(twoList, num2)
		}
	}
	sumlen := len(oneList)
	if len(twoList) > sumlen {
		sumlen = len(twoList)
	}
	keyval := 0
	resultList := []int{}
	for i := 0; i < sumlen; i++ {
		tnum, qnum := 0, 0
		onelen := len(oneList) - 1 - i
		if onelen >= 0 {
			tnum = oneList[len(oneList)-1-i]
		}

		twolen := len(twoList) - 1 - i
		if twolen >= 0 {
			qnum = twoList[len(twoList)-1-i]
		}

		sumcount := tnum + qnum + keyval
		if sumcount > 9 {
			keyval = 1
		} else {
			keyval = 0
		}
		resultList = append(resultList, sumcount%10)
	}
	if keyval == 1 {
		resultList = append(resultList, 1)
	}
	for i := len(resultList) - 1; i >= 0; i-- {
		// log.Println(resultList[i])
		resVal := &ListNodeData{
			Val: resultList[i],
		}
		if i == len(resultList)-1 {
			result = resVal
			resKeyres = resVal
		} else {
			resKeyres.Next = resVal
			resKeyres = resKeyres.Next
		}
	}
	return result
}

func getNodeNum(node *ListNodeData, cnum chan int) {
	for node != nil {
		cnum <- node.Val
		node = node.Next
	}
	cnum <- -1
}
