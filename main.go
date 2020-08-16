package main

import "maru/cmd"

//go:generate gitfs ./...

func main() {
	cmd.Execute()
}
