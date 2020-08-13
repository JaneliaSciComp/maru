package utils

import (
	"fmt"
	"os"

	Aurora "github.com/logrusorgru/aurora"
)

// PrintMessage - print a normal message
func PrintMessage(format string, a ...interface{}) {
	fmt.Println(Aurora.Sprintf(Aurora.White(format), a...))
}

// PrintInfo - print an info message
func PrintInfo(format string, a ...interface{}) {
	fmt.Println(Aurora.Sprintf(Aurora.BrightBlue(format), a...))
}

// PrintSuccess - prints an error message
func PrintSuccess(format string, a ...interface{}) {
	fmt.Println(Aurora.Sprintf(Aurora.BrightGreen(format), a...))
}

// PrintError - prints an error message
func PrintError(format string, a ...interface{}) {
	fmt.Println(Aurora.Sprintf(Aurora.BrightRed("[Error] "+format), a...))
}

// PrintFatal - prints an error message and exits with code 2
func PrintFatal(format string, a ...interface{}) {
	fmt.Println(Aurora.Sprintf(Aurora.BrightRed("[Error] "+format), a...))
	os.Exit(2)
}

// Mkdir - Make the given path
func Mkdir(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		e := os.Mkdir(path, 0777)
		if e == nil {
			return true
		}
		PrintFatal("%s", e)
	}
	return false
}

// FileExists - returns true if the given file exists
func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
