package ftm_color

import "fmt"

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

func PrintError(errorMessage string) {
	fmt.Println(Red, errorMessage, Reset)
}

func PrintMessage(color string, message string) {
	fmt.Println(color, message, Reset)
}
