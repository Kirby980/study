package slice_test

import (
	"fmt"
	"study/study1/slice"
	"testing"
)

func TestSlice(t *testing.T) {
	s1 := []int{1, 2, 3, 4, 5, 6}
	fmt.Println(slice.Delete(s1, 2))
	s2 := []int32{1, 2, 3, 4, 5, 6}
	fmt.Println(slice.Delete(s2, 3))
	s3 := []string{"a", "b", "c", "d", "e", "f"}
	fmt.Println(slice.Delete(s3, -1))
	fmt.Println(slice.Delete(s3, 4))
	s4 := []float64{1.1, 2.2, 3.3, 4.4, 5.5, 6.6}
	fmt.Println(slice.Delete(s4, 2))

}
