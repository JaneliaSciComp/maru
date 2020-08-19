package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	Aurora "github.com/logrusorgru/aurora"
)

const MaruVersion = "0.1.0"

// PrintMessage - print a normal message
func PrintMessage(format string, a ...interface{}) {

	finalString := fmt.Sprintf(format, a...)
	// TODO: replace all backticks in the codebase with carrots to not conflict with multiline strings
	fixedString := strings.ReplaceAll(finalString, "`", "^")
	parts := strings.Split(fixedString, "^")

	for i, s := range parts {
		if i % 2 == 0 {
			fmt.Print(s)
		} else {
			fmt.Print(Aurora.BrightCyan(s))
		}
	}

	fmt.Println()
}

// PrintInfo - print an info message
func PrintInfo(format string, a ...interface{}) {

	finalString := fmt.Sprintf(format, a...)
	// TODO: replace all backticks in the codebase with carrots to not conflict with multiline strings
	fixedString := strings.ReplaceAll(finalString, "`", "^")
	parts := strings.Split(fixedString, "^")

	for i, s := range parts {
		if i % 2 == 0 {
			fmt.Print(Aurora.BrightBlue(s))
		} else {
			fmt.Print(Aurora.BrightCyan(s))
		}
	}

	fmt.Println()
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

func RunCommand(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}