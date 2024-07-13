package main

import (
	"fmt"
	"os"

	"github.com/Shresth72/tor/pkg/command"
)

func main() {
	cmd := os.Args[1]
	arg := os.Args[2]

	jsonOutput, err := command.ExecuteCommand(cmd, arg)
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("%s", jsonOutput)
}
