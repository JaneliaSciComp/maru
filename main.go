package main

import "jape/cmd"

//go:generate gitfs ./...

func main() {
	cmd.Execute()
}
