package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Stremax-Team/stremax-lang/pkg/interpreter"
)

func main() {
	// Define command-line flags
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)
	runFile := runCmd.String("file", "", "Path to the Stremax-Lang file to run")

	// Check if a command was provided
	if len(os.Args) < 2 {
		fmt.Println("Expected 'run' subcommand")
		os.Exit(1)
	}

	// Parse the command
	switch os.Args[1] {
	case "run":
		runCmd.Parse(os.Args[2:])
		if *runFile == "" {
			fmt.Println("Please provide a file to run with -file flag")
			os.Exit(1)
		}
		runProgram(*runFile)
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func runProgram(filePath string) {
	// Read the file
	source, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}

	// Create an interpreter and run the program
	i := interpreter.New(string(source))
	err = i.Run()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}
