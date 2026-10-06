package main

import "fmt"

type DS struct {
	arr []int
}

func (d *DS) even() []int {
	res := make([]int, 0, len(d.arr))

	for _, v := range d.arr {

		if v%2 == 0 {
			res = append(res, 0)
		} else {

			res = append(res, 1)
		}
	}
	return res
}

func main() {
	ds := DS{
		arr: []int{-22, -1, 0, 12, 22},
	}

	fmt.Println(ds.even())

}
