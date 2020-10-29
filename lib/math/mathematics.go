package math

import "fmt"

// 第一种写法
func quickSort(values []int, left, right int) {
	temp := values[left]
	p := left
	i, j := left, right

	for i <= j {
		for j >= p && values[j] >= temp {
			j--
		}
		if j >= p {
			values[p] = values[j]
			p = j
		}

		for i <= p && values[i] <= temp {
			i++
		}
		if i <= p {
			values[p] = values[i]
			p = i
		}
	}
	values[p] = temp
	if p-left > 1 {
		quickSort(values, left, p-1)
	}
	if right-p > 1 {
		quickSort(values, p+1, right)
	}
}

//QuickSort 1
func QuickSort(values []int) {
	if len(values) <= 1 {
		return
	}
	quickSort(values, 0, len(values)-1)
}

//Quick2Sort 2
func Quick2Sort(values []int) {
	if len(values) <= 1 {
		return
	}
	mid, i := values[0], 1
	head, tail := 0, len(values)-1
	for head < tail {
		fmt.Println(values)
		if values[i] > mid {
			values[i], values[tail] = values[tail], values[i]
			tail--
		} else {
			values[i], values[head] = values[head], values[i]
			head++
			i++
		}
	}
	values[head] = mid
	Quick2Sort(values[:head])
	Quick2Sort(values[head+1:])
}

//Quick3Sort 3
func Quick3Sort(a []int, left int, right int) {
	if left >= right {
		return
	}
	explodeIndex := left
	for i := left + 1; i <= right; i++ {

		if a[left] >= a[i] {

			//分割位定位++
			explodeIndex++
			a[i], a[explodeIndex] = a[explodeIndex], a[i]
		}
	}
	//起始位和分割位
	a[left], a[explodeIndex] = a[explodeIndex], a[left]
	Quick3Sort(a, left, explodeIndex-1)
	Quick3Sort(a, explodeIndex+1, right)

}
