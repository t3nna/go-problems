package main

import (
	"fmt"
	"slices"
)

func main() {
	fmt.Printf("TheLift(): %v\n", TheLift([][]int{{3}, {3}, {6, 5, 4, 1}, {}, {}, {1, 2}, {4}}, 5)) //[]int{0, 2, 5, 0}
}

func TheLift(queues [][]int, capacity int) []int {

	for _, val := range queues {
		slices.Sort(val)
	}

	for _, val := range queues {
		fmt.Println(val)
	}

	res := make([]int, 0, len(queues))
	liftDir := 1
	currFloor := 0

	for {
		if liftDir > 0 && currFloor >= len(queues)-1 {
			if isLiftRequested(queues) {
				liftDir = changeDir(liftDir)
			} else {
				break
			}
		}
		if liftDir < 0 && currFloor == 0 {
			if isLiftRequested(queues) {
				liftDir = changeDir(liftDir)
			} else {
				break
			}

		}

		shouldMakeStop := false
		curr := queues[currFloor]
		newQueue := make([]int, 0, len(curr))
		for i := 0; i < len(queues[currFloor]); i++ {
			currItem := curr[i]
			fmt.Println(currItem, " ", isEligible(liftDir, currFloor, currItem))
			if isEligible(liftDir, currFloor, currItem) {
				// queues[currFloor] = append(curr[:i], curr[i+1:]...)
				shouldMakeStop = true
			} else {
				newQueue = append(newQueue, currItem)
			}
		}
		if shouldMakeStop {
			res = append(res, currFloor)
		}

		queues[currFloor] = newQueue

		currFloor += liftDir

	}
	return res
}

func changeDir(dir int) int {
	if dir > 0 {
		return -1
	}
	return 1
}

func isLiftRequested(queues [][]int) bool {
	for i := 0; i < len(queues); i++ {
		if len(queues[i]) != 0 {
			return true
		}
	}
	return false
}

func isEligible(liftDir, currFloor, desiredFloor int) bool {

	// going up
	if desiredFloor-currFloor > 0 {
		if liftDir > 0 {
			return true
		} else {
			return false
		}
	}
	//going down

	if liftDir < 0 {
		return true
	} else {
		return false
	}

}
