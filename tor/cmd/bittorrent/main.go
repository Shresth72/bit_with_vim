package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codecrafters-io/tor/pkg/decode"
)

func main() {
	command := os.Args[1]

	if command == "decode" {
		bencodedValue := os.Args[2]

		decoded, _, err := decode.DecodeBencode(bencodedValue, 0)
		if err != nil {
			fmt.Println(err)
			return
		}

		jsonOutput, _ := json.Marshal(decoded)
		fmt.Println(string(jsonOutput))
	} else {
		fmt.Println("Unknown command: " + command)
		os.Exit(1)
	}
}
