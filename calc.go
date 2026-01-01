package main

import (
	"errors"
	"math"
)

func calculate(a, b float64, op string) (float64, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		return a / b, nil
	case "^":
		return math.Pow(a, b), nil
	}
	return 0, errors.New("Неверно введено число или операция")
}
