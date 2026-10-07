package main

import (
	"errors"
	"fmt"
)

var ErrDivideByZero = errors.New("cannot divide by zero")

type RequestError struct {
	StatusCode int
	Message    string
}

func (e *RequestError) Error() string {
	return fmt.Sprintf("request failed with status %d: %s", e.StatusCode, e.Message)
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivideByZero
	}
	return a / b, nil
}

func fetchResource(url string) error {
	if url == "" {
		return &RequestError{StatusCode: 400, Message: "empty url"}
	}
	return nil
}

func process() error {
	_, err := divide(10, 0)
	if err != nil {
		return fmt.Errorf("process failed: %w", err)
	}
	return nil
}

func main() {
	result, err := divide(10, 2)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}

	_, err = divide(10, 0)
	if err != nil {
		fmt.Println("Error:", err)
	}

	err = process()
	if errors.Is(err, ErrDivideByZero) {
		fmt.Println("Detected divide by zero:", err)
	}

	err = fetchResource("")
	var reqErr *RequestError
	if errors.As(err, &reqErr) {
		fmt.Printf("Custom error: status=%d, msg=%s\n", reqErr.StatusCode, reqErr.Message)
	}
}
