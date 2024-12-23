package main

import (
	"api/pkg/bootstrap"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var debug bool

func main() {
	if len(os.Args) < 2 {
		errorHandler("insufficient command line arguments.\nuse 'serve' or 'migrate' command")
	}
	args := os.Args[2:]
	flag.BoolVar(&debug, "debug", false, "debug mode")
	err := flag.CommandLine.Parse(args)
	if err != nil {
		errorHandler(err.Error())
	}
	if debug {
		log.SetFlags(0)
	}
	switch command := os.Args[1]; command {
	case "serve":
		bootstrap.Serve(debug)
	case "migrate":
		bootstrap.Migrate(debug)
	default:
		errorHandler(
			fmt.Sprintf("unknown command: %s\nuse 'serve'or 'migrate' command", command),
		)
	}

}

func errorHandler(msg string) {
	msg = strings.Replace(msg, "\n", "\n\t", -1)
	err := fmt.Errorf("error:\n\t%s", msg)
	log.Println(err.Error())
	os.Exit(1)
}
