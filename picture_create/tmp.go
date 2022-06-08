package main

import (
	"fmt"
)

func main() {
	fmt.Println(subsets([]int{9, 0, 3, 5, 7}))
}

func subsets(nums []int) (res [][]int) {
	res = append(res, []int{})
	for _, val := range nums {
		for _, slice := range res {
			res = append(res, append(slice, val))
			fmt.Printf("%v\t%p\n", res[len(res)-1], &res[len(res)-1][0])
		}
	}
	return res
}
