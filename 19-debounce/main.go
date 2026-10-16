package main

import (
	"fmt"
	"time"
)

func debounce(delay int, fn func(any)) func(any) {
	var t *time.Timer

	return func(a any) {
		if t != nil {
			t.Stop()
		}
		t = time.AfterFunc(time.Millisecond*time.Duration(delay), func() {
			fn(a)
		})
	}
}

func main() {
	printLog := func(val any) {
		fmt.Printf("[%s] Executed with value: %v\n", time.Now().Format("15:04:05.000"), val)
	}

	debouncedPrint := debounce(500, printLog)

	fmt.Printf("[%s] Starting rapid calls...\n", time.Now().Format("15:04:05.000"))

	debouncedPrint("A")
	time.Sleep(100 * time.Millisecond) // Not enough time has passed
	debouncedPrint("B")
	time.Sleep(100 * time.Millisecond) // Not enough time has passed
	debouncedPrint("C")                // This is the last call, it should execute after 500ms

	time.Sleep(1 * time.Second)

}
