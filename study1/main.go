package main

import (
	"errors"
	"fmt"
)

func main() {
	s1 := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(Delete(s1, 2))
	s2 := []int32{1, 2, 3, 4, 5, 6}
	fmt.Println(Delete(s2, 3))
	s3 := []string{"a", "b", "c", "d", "e", "f"}
	fmt.Println(Delete(s3, -1))
	fmt.Println(Delete(s3, 4))

}

func Delete[T anyType](slice []T, idex int) ([]T, error) {
	if idex < 0 || idex >= len(slice) {
		return nil, errors.New("下标错误")
	}
	prev := slice[:idex]
	next := slice[idex+1:]
	prev = append(prev, next...)
	return prev, nil
}

type anyType interface {
	~int8 | ~int | ~int32 | ~int64 | ~uint8 | ~uint | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}
