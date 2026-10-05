package main

import "fmt"

func main() {
	s1 := make([]int, 3, 6)
	s1 = append(s1, 1, 2, 3, 4, 5)

	s2 := s1[1:3]
	s3 := append([]int(nil), 24)
	s4 := append([]int(nil), s1...)

	fmt.Println(s3)

	fmt.Printf("%v, %v\n", len(s2), cap(s2))

	fmt.Println(s4)

	s5 := append([]int(nil), s4...)
	fmt.Println(s5)
	fmt.Println("Trying Zed")
}
