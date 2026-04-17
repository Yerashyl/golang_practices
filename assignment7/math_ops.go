package assignment7

import "errors"

// Divide receives two integer numbers and proceeds arithmetical operations.
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Subtract receives two integer numbers and proceeds arithmetical operations.
func Subtract(a, b int) int {
	return a - b
}
