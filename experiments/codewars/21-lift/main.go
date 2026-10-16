package main

import (
	"fmt"
)

func main() {
	fmt.Printf("TheLift(): %v\n", TheLift([][]int{{}, {}, {0}}, 5))                                 //[]int{0, 2, 5, 0}
	fmt.Printf("TheLift(): %v\n", TheLift([][]int{{}, {}, {1, 1}, {}, {}, {}, {}}, 5))              //[]int{0, 2, 5, 0}
	fmt.Printf("TheLift(): %v\n", TheLift([][]int{{3}, {3}, {6, 5, 4, 1}, {}, {}, {1, 2}, {4}}, 5)) //[]int{0, 2, 5, 0}
	fmt.Printf("TheLift([][]int{{}, {0}, {}, {}, {2}, {3}, {}}, 5): %v\n", TheLift([][]int{{}, {0}, {}, {}, {2}, {3}, {}}, 5))
}

type LiftState struct {
	currFloor    int
	dir          int
	amountInside int
	capacity     int
	inside       map[int]int
	res          []int
	height       int
	house        [][]int
}

func TheLift(queues [][]int, capacity int) []int {
	if len(queues) < 1 {
		return []int{}
	}
	lift := LiftState{
		capacity:     capacity,
		currFloor:    0,
		dir:          1,
		amountInside: 0,
		inside:       make(map[int]int),
		res:          make([]int, 0, len(queues)),
		height:       len(queues),
		house:        queues,
	}

	lift.res = append(lift.res, lift.currFloor)

	for {
		if lift.shouldStopOnCurr(lift.house[lift.currFloor]) {
			lift.addStop(lift.currFloor)
			lift.disembarking()
			lift.house[lift.currFloor] = lift.boarding(lift.house[lift.currFloor])
		}
		if lift.isWorkLeft() {
			lift.currFloor += lift.dir
		} else {
			lift.dir = -lift.dir
			if lift.shouldStopOnCurr(lift.house[lift.currFloor]) {
				continue
			}
			if !lift.isWorkLeft() {
				lift.addStop(0)
				break
			}
		}
	}
	return lift.res
}

func (l *LiftState) addStop(floor int) {
	if len(l.res) > 0 && l.res[len(l.res)-1] != floor {
		l.res = append(l.res, floor)
	}
}

func (l *LiftState) disembarking() {
	l.amountInside -= l.inside[l.currFloor]
	delete(l.inside, l.currFloor)
}

func (l *LiftState) boarding(q []int) []int {
	restQ := make([]int, 0, len(q))

	for _, v := range q {
		m := l.isDirMatching(v)

		if m && l.capacity >= l.amountInside+1 {
			l.amountInside += 1
			if _, ok := l.inside[v]; !ok {
				l.inside[v] = 1
			} else {
				l.inside[v] += 1
			}
		} else {
			restQ = append(restQ, v)
		}

	}
	return restQ
}

func (l *LiftState) isDirMatching(p int) bool {
	isDirUp := l.dir > 0
	isPasUp := p-l.currFloor > 0

	if isDirUp && isPasUp {
		return true
	}
	if !isDirUp && !isPasUp {
		return true
	}
	return false
}

func (l *LiftState) shouldStopOnCurr(q []int) bool {
	isCandBoard := false
	isExit := l.inside[l.currFloor] > 0

loop:
	for i := 0; i < len(q); i++ {
		if l.isDirMatching(q[i]) {
			isCandBoard = true
			break loop
		}
	}

	return isCandBoard || isExit

}

func (l *LiftState) isWorkLeft() bool {
	isLeft := false

	if l.dir == 1 {

		for k := range l.inside {
			if k > l.currFloor {
				return true
			}
		}
		for i := l.currFloor + 1; i < l.height; i++ {
			// if l.shouldStop(l.house[i]) {
			curr := l.house[i]

			if len(curr) > 0 {

				return true
			}
			// }

		}
	} else {
		for k := range l.inside {
			if k < l.currFloor {
				return true
			}
		}
		for i := l.currFloor - 1; i >= 0; i-- {
			curr := l.house[i]

			if len(curr) > 0 {

				return true
			}
		}

	}
	return isLeft
}
