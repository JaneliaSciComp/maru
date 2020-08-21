package utils

import (
	"bufio"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"os"
	"os/exec"
	"strings"

	"crypto/sha256"
	Aurora "github.com/logrusorgru/aurora"
	"io"
)

const MaruVersion = "0.1.0"
const DockerFilePath = "Dockerfile"
var Debug = false

type ColorFunc func(arg interface{}) Aurora.Value

// PrintDebug - prints an debug message if debug is turned on
func PrintDebug(format string, a ...interface{}) {
	if Debug {
		fmt.Println(Aurora.Sprintf(Aurora.Yellow(format), a...))
	}
}

// PrintHint - prints a hint message in a darker than normal color so that it's readable but not eye-catching
func PrintHint(format string, a ...interface{}) {
	print(Aurora.White, format, a...)
}

func PrintMessage(format string, a ...interface{}) {
	print(nil, format, a...)
}

func PrintInfo(format string, a ...interface{}) {
	print(Aurora.BrightBlue, format, a...)
}

// PrintSuccess - prints an error message
func PrintSuccess(format string, a ...interface{}) {
	print(Aurora.BrightGreen, "\u2714 "+format, a...)
}

// PrintError - prints an error message
func PrintError(format string, a ...interface{}) {
	print(Aurora.BrightRed, "\u2718 "+format, a...)
}

// PrintFatal - prints an error message and exits with code 2
func PrintFatal(format string, a ...interface{}) {
	print(Aurora.BrightRed, "\u2714 "+format, a...)
	os.Exit(2)
}

// Print a message with a default color, and optional code highlighting.
// Highlighting is applied to any string between carots ^like this^.
func print(colorFunc ColorFunc, format string, a ...interface{}) {

	finalString := fmt.Sprintf(format, a...)
	// TODO: replace all backticks in the codebase with carrots to not conflict with multiline strings
	fixedString := strings.ReplaceAll(finalString, "`", "^")
	parts := strings.Split(fixedString, "^")

	for i, part := range parts {
		if i % 2 == 0 {
			if colorFunc != nil {
				// Use the color function if available
				fmt.Print(colorFunc(part))
			} else {
				// Otherwise, no formatting
				fmt.Print(part)
			}
		} else {
			// Format as code
			fmt.Print(Aurora.BrightMagenta(part))
		}
	}

	fmt.Println()
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

// Ask the user one question and get input. Deals with Ctrl-C interruptions and other errors.
func Ask(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) {
	err := survey.AskOne(prompt, response)
	if err == terminal.InterruptErr {
		fmt.Println("interrupted")
		os.Exit(0)
	} else if err != nil {
		PrintFatal("%s", err)
	}
}

// Ask the user one question requiring a string input.
func AskForString(message string, defaultValue string) string {
	value := defaultValue
	prompt := &survey.Input{
		Message: message,
		Default: value,
	}
	Ask(prompt, &value)
	return value
}

// Ask the user one question requiring a yes/no input.
func AskForBool(msg string, defaultValue bool) bool {
	value := defaultValue
	prompt := &survey.Confirm{
		Message: msg,
		Default: value,
	}
	Ask(prompt, &value)
	return value
}

func getChecksum(s string) string {

	h := sha256.New()
	if _, err := io.Copy(h, strings.NewReader(s)); err != nil {
		PrintFatal("%s",err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// Read the checksum encoded in the Dockerfile
func GetChecksumFromDockerfile() string {

	f, err := os.Open(DockerFilePath)
	if err != nil {
		PrintFatal("%s", err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	first := ""
	if scanner.Scan() {
		first = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		PrintFatal("%s", err)
	}
	return strings.TrimSpace(strings.Replace(first, "# ", "", 1))
}

func GetChecksumOfConfig() string {

	// Extract just the "config" section of the configuration file
	f2, err := os.Open(ConfFile)
	if err != nil {
		PrintFatal("%s",err)
	}
	defer f2.Close()
	scanner2 := bufio.NewScanner(f2)
	var sb strings.Builder
	rec := false
	for scanner2.Scan() {
		line := scanner2.Text()
		if strings.TrimSpace(line) == "config:" {
			rec = true // start recording
		} else if rec && strings.TrimLeft(line, " \t") == line {
			rec = false // stop recording when we reach another top level attribute
		}
		if rec {
			sb.WriteString(line)
			sb.WriteRune('\n')
		}
	}
	if err := scanner2.Err(); err != nil {
		PrintFatal("%s", err)
	}

	// Compute the current checksum
	return getChecksum(sb.String())
}

func TestChecksum(newChecksum string) bool {
	existingChecksum := GetChecksumFromDockerfile()
	return newChecksum == existingChecksum
}
