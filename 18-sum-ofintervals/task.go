package main

import (
	"cmp"
	"fmt"
	"slices"
)

func SumOfIntervals(intervals [][2]int) int {
	// sort intervals
	slices.SortFunc(intervals, func(a, b [2]int) int {
		return cmp.Compare(a[0], b[0])
	})
	// merge
	res := make([]*[2]int, 0, len(intervals))
	res = append(res, &intervals[0])
	for i := 1; i < len(intervals); i++ {
		ref := res[len(res)-1]
		curr := intervals[i]
		//fmt.Println("ref: ", ref, "curr: ", curr)

		if ref[1] > curr[0] {
			ref[1] = max(ref[1], curr[1])
		} else {
			res = append(res, &curr)
		}

	}
	length := 0
	for _, v := range res {
		if v != nil {
			fmt.Println(*v)
			length += v[1] - v[0]
		} else {
			fmt.Println("<nil>")
		}
	}
	return length
}

func main() {
	fmt.Println("res: ", SumOfIntervals([][2]int{{1, 4}, {7, 10}, {3, 5}}))
	SumOfIntervals([][2]int{{0, 20}, {-100_000_000, 10}, {30, 40}})
}
