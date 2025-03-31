package slice

import "errors"

func Delete[T AnyType](slice []T, index int) ([]T, error) {
	if index < 0 || index >= len(slice) {
		return nil, errors.New("下标错误")
	}
	copy(slice[index:], slice[index+1:])
	return slice[:len(slice)-1], nil
}

type AnyType interface {
	~int8 | ~int | ~int32 | ~int64 | ~uint8 | ~uint | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}
