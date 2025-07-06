package utils

import "fmt"

type ColorCode string

const (
	RED     ColorCode = "\033[31m"
	GREEN   ColorCode = "\033[32m"
	YELLOW  ColorCode = "\033[33m"
	BLUE    ColorCode = "\033[34m"
	MAGENTA ColorCode = "\033[35m"
	CYAN    ColorCode = "\033[36m"
	RESET   ColorCode = "\033[0m"
)

func GetColoredString[T interface{ ~int | ~int64 | ~string }](v T, color ColorCode) string {
	return fmt.Sprint(color, v, RESET)
}

func GetColoredStatusCode(code int) string {
	var color ColorCode

	switch {
	case code >= 500:
		color = RED
	case code >= 400:
		color = YELLOW
	case code >= 200:
		color = GREEN
	default:
		color = RESET
	}

	return GetColoredString(code, color)
}

func GetColoredHttpMethod(method string) string {
	var color ColorCode

	switch method {
	case "DELETE":
		color = RED
	default:
		color = MAGENTA
	}

	return GetColoredString(method, color)
}
