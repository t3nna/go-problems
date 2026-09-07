package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	s := Solution([]int{-6, -3, -2, -1, 0, 1, 3, 4, 5, 7, 8, 9, 10, 11, 14, 15, 17, 18, 20})
	fmt.Println(s) //"-6,-3-1,3-5,7-11,14,15,17-20"
}

func Solution(list []int) string {

	res := make([][2]int, 0)

	i := 0
	for i < len(list) {
		start := i
		j := i + 1

		for j < len(list) && list[j] == list[j-1]+1 {
			j++
		}
		end := j - 1

		if start == end {
			res = append(res, [2]int{list[start], list[start]})
		} else {
			res = append(res, [2]int{list[start], list[end]})
		}
		i = j
	}

	fmt.Println(res)
	out := make([]string, 0)

	for _, val := range res {
		if val[0] == val[1] {
			out = append(out, strconv.Itoa(val[0]))
		} else {
			out = append(out, fmt.Sprint(strconv.Itoa(val[0]), "-", strconv.Itoa(val[1])))
		}
	}

	return strings.Join(out, ",")
}
