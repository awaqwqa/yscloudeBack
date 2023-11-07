package utils

import "fmt"

func RemoveIndex[T any](s []T, index int) ([]T, error) {
	if index < 0 || index >= len(s) {
		return nil, fmt.Errorf("index out of range")
	}
	return append(s[:index], s[index+1:]...), nil
}
