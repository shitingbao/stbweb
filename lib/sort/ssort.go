package ssort

import (
	"log"
	"sort"
)

// 原理
// 使用了sort原生包
// 他是一个interface
// 只要实现下列三个方法的都可以进行排序

//SortList 将[]string定义为SortList类型
type SortList []string

// 实现sort.Interface接口的获取元素数量方法
func (m SortList) Len() int {
	return len(m)
}

// 实现sort.Interface接口的比较元素方法
func (m SortList) Less(i, j int) bool {
	return m[i] < m[j]
}

// 实现sort.Interface接口的交换元素方法
func (m SortList) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

//例子，只要是使用SortList该类型，就可以对其进行排序
//默认ascii码递增排序，96，97，98。。。
func example() {
	names := SortList{"app_key", "end_time", "method"}
	// 使用sort包进行排序
	sort.Sort(names) // 这个是原生方法
	log.Println("names:", names)
}
