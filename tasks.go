package main

import (
	"fmt"
	"strings"
)

func run() {
	// fmt.Println(lenOflongetstSubstring("eceba", 2))
	// fmt.Println(lenOflongetstSubstring("eceba", 3))
	// fmt.Println(lenOflongetstSubstring("eceba", 4))
	// fmt.Println(lenOflongetstSubstring("eedba", 1))
	// fmt.Println(lenOflongetstSubstring("edeba", 1))
	// fmt.Println(lenOflongetstSubstring("eeeba", 1))
	// fmt.Println(lenOflongetstSubstring("eeeba", 0))
	// fmt.Println(lenOflongetstSubstring("abbceeecc", 2))
	// fmt.Println(lenOflongetstSubstring("", 1))
	// fmt.Println(lenOflongetstSubstring("abaaac", 2))
	// fmt.Println(lenOflongetstSubstring("babc", 2))
	// fmt.Println(lenOflongetstSubstring("", 0))
	moveRightZeroes([]int{0, 8, 0, 9, 1, 2, 3, 0, 0, 0, 0, 4, 5, 6, 0})
}

func lenOflongetstSubstring(s string, k int) int {
	ptr1 := strings.NewReader(s)
	var max, current int
	mapChars := make(map[rune]int)
	for _, c := range s {
		mapChars[c]++
		for cntIter := 0; len(mapChars) > k; cntIter++ {
			r, _, _ := ptr1.ReadRune()
			if mapChars[r] > 1 {
				mapChars[r]--
			} else {
				delete(mapChars, r)
			}
			current--
			if cntIter >= 20 {
				break
			}
		}
		current++
		if current > max {
			max = current
		}
	}
	return max
}

func moveRightZeroes(s []int) {
	fmt.Println(s)
	var cntSwap int
	for idx, num := range s {
		if num != 0 {
			s[idx-cntSwap] = num
			s[idx] = 0
		} else {
			cntSwap++
		}
	}
	fmt.Println(s)
}
