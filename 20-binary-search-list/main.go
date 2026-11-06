package main

import "fmt"

func bsList(a []int, needle int) int {
	lo := 0
	hi := len(a)

	idx := -1

	for lo < hi {
		mid := lo + (hi-lo)/2
		val := a[mid]

		if val == needle {

			idx = mid
			break
		} else if val > needle {
			hi = mid
		} else {
			lo = mid + 1
		}
	}

	return idx
}

func main() {
	arr := []int{1, 2, 3, 6, 9, 11, 12, 16}
	fmt.Println(bsList(arr, 11))
}
