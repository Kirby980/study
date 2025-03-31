package main

import "fmt"

func main() {
	s := []int{1, 2, 3, 4, 5}
	s1 := s[:3]
	fmt.Println(len(s1), "------", cap(s1))
}
