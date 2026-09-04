package main

import (
	"fmt"
)

func main() {
	snailMap := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	fmt.Println(Snail(snailMap))
}

func Snail(snaipMap [][]int) []int {

	dir := map[string][]int{
		"r": {0, 1},
		"b": {-1, 0},
		"l": {0, -1},
		"t": {1, 0},
	}

	//   dirArr := [][]int{{0, 1}, {-1, 0}, {0, -1}, {1, 0}}

	dim := len(snaipMap)

	visited := make([][]bool, dim)

	for i := range visited {
		visited[i] = make([]bool, dim)
	}

	res := make([]int, 0, dim*dim)
	// (y, x)
	place := [2]int{0, 0}
	currDir := "r"
	for true {
		res = append(res, snaipMap[place[0]][place[1]])

		newDir, isChanged := checkDir(place[0], place[1], dim, currDir, visited)
		if isChanged {
			currDir = newDir
		}
		// make step
		place[0] += dir[currDir][0]
		place[1] += dir[currDir][1]
		visited[place[0]][place[1]] = true

		fmt.Println(visited)
	}

	return []int{}
}

func checkDir(y, x, dim int, currDir string, visited [][]bool) (string, bool) {
	dim--
	if currDir == "r" && (x == dim || visited[y][x]) {
		return "b", true
	}

	if currDir == "b" && (y == dim || visited[y][x]) {
		return "l", true
	}

	if currDir == "l" && (x == 0 || visited[y][x]) {
		return "t", true
	}

	if currDir == "t" && (y == 0 || visited[y][x]) {
		return "r", true
	}
	return "", false
}
