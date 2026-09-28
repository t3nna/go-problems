package main

import "fmt"

func main() {
	fmt.Printf("FindReverseNumber(10): %v\n", FindReverseNumber(10000000))
}

func FindReverseNumber(n uint64) uint64 {
	count := uint64(0)
	i := uint64(0)
	for {
		if isReverse(i) {
			count++
		}
		if count == n {
			return i
		}
		i++
	}
}

func isReverse(n uint64) bool {
	m := n
	for m >= 10 {
		if !compareFirstAndLast(m) {
			return false
		}
		// remove fist and last digit
		m = m / 10

		if m < 10 {
			return true
		}

		m = RemoveFirstDigitUint64(m)

	}
	return true
}

func compareFirstAndLast(n uint64) bool {
	if n < 10 {
		return true
	}

	last := n % 10

	first := n
	for first >= 10 {
		first /= 10
	}
	return last == first
}

func RemoveFirstDigitUint64(n uint64) uint64 {
	switch {
	case n < 10:
		return 0
	case n < 100:
		return n % 10
	case n < 1000:
		return n % 100
	case n < 10000:
		return n % 1000
	case n < 100000:
		return n % 10000
	case n < 1000000:
		return n % 100000
	case n < 10000000:
		return n % 1000000
	case n < 100000000:
		return n % 10000000
	case n < 1000000000:
		return n % 100000000
	case n < 10000000000:
		return n % 1000000000
	case n < 100000000000:
		return n % 10000000000
	case n < 1000000000000:
		return n % 100000000000
	case n < 10000000000000:
		return n % 1000000000000
	case n < 100000000000000:
		return n % 10000000000000
	case n < 1000000000000000:
		return n % 100000000000000
	case n < 10000000000000000:
		return n % 1000000000000000
	case n < 100000000000000000:
		return n % 10000000000000000
	case n < 1000000000000000000:
		return n % 100000000000000000
	case n < 10000000000000000000:
		return n % 1000000000000000000
	default:
		return n % 10000000000000000000
	}
}
