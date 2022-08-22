package main

import (
	"flag"
	"fmt"
	"github.com/frakaft/gomodsearch/internal/gomodsearch"
	"os"
)

func main() {

	// Init args parser
	help := flag.Bool("help", false, "Show help")
	path := flag.String("path", "", "Project path")
	mod := flag.String("mod", "", "Project path")
	flag.Parse()

	// Help
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	// Missing args
	if *path == "" || *mod == "" {
		flag.Usage()
		os.Exit(2)
	}

	// Run modch
	err := gomodsearch.Run(path, mod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
