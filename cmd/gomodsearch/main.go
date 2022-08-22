package main

import (
	"fmt"
	"github.com/francogeller/gomodsearch/internal/gomodsearch"
	"os"
	"strings"
)

func main() {

	// Init args parser
	args := os.Args[1:]

	// Print help
	if len(args) == 1 && args[0] == "help" {
		printHelp()
		os.Exit(0)
	}

	// Print error
	if len(args) < 2 {
		printHelpError(&args)
		os.Exit(2)
	}

	// Run modch
	err := gomodsearch.Run(args[0], args[1:]...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func printHelp() {
	fmt.Fprintf(os.Stdout, "usage: gomodsearch <path> <module>[@version] ...\n")
}

func printHelpError(args *[]string) {
	if len(*args) > 0 {
		fmt.Fprintf(os.Stderr, "gomodsearch: Invalid arguments '%s'.\n", strings.Join(*args, " "))
	} else {
		fmt.Fprintf(os.Stderr, "gomodsearch: Missing arguments.\n")
	}
	fmt.Fprintf(os.Stderr, "run 'gomodsearch help' for usage.\n")
}
