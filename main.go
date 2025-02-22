package main

import (
	"fmt"
	"libnf/examples"
	"os"
	"strconv"
)

func helpAndExit() {
	fmt.Println("Usage: <program> <option>")
	fmt.Println("Options:")
	fmt.Println("  1 - Run reader()")
	fmt.Println("  2 - Run writer()")
	os.Exit(1)
}

func main() {
	// Check if an argument is provided
	if len(os.Args) < 2 {
		helpAndExit()
	}

	// Parse the argument
	option, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Error: Argument must be a number.")
		os.Exit(1)
	}

	// Call the corresponding function based on the argument
	switch option {
	case 2:
		examples.Reader()
	case 1:
		examples.Writer()
	default:
		fmt.Println("Invalid option.")
		helpAndExit()
	}
}
