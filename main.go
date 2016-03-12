package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <config.json> <port>", os.Args[0])
		os.Exit(1)
	}
	_, err := GetConfig(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to load config:", err)
		os.Exit(1)
	}
}
