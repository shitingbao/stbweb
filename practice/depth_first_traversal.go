package main

import "log"

//深度优先遍历
func isBipartite(graph [][]int) bool {
	gh := make(map[int]int)
	isExist := make(map[int]bool)
	for i := range graph {
		if gh[i] == 0 {
			gh[i] = 1
			spk(i, 1, graph, gh, isExist)
		}

	}
	if isExist[0] {
		return false
	}
	return true
}

func spk(node, color int, graph [][]int, gh map[int]int, isExist map[int]bool) {
	log.Println(node)
	if color == 1 {
		color = 2
	} else {
		color = 1
	}
	for _, v := range graph[node] {
		if gh[v] != 0 && gh[v] == gh[node] {
			isExist[0] = true
			return
		}
		if gh[v] == 0 {
			gh[v] = color
			if len(gh) == len(graph) {
				return
			}
			spk(v, color, graph, gh, isExist)
		}
	}
}

func isBipartiteLoad() {
	// listL := [][]int{{1, 3}, {0, 2}, {1, 3}, {0, 2}}
	listL := [][]int{{4, 1}, {0, 2}, {1, 3}, {2, 4}, {3, 0}}
	log.Println(isBipartite(listL))
}
