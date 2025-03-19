package main

import (
	"gojira/cmd"
)

func main() {
	cmd.Version = "1.0.0"
	cmd.Execute()
}
