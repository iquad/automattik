package main

import (
	"flag"

	"bufio"
	"os"

	"github.com/iquad/automattik/config"
)

func main() {
	conf := config.NewConfig()
	// processFlags := conf.BindFlags()
	defaults := flag.Bool("print-defaults", false, "print the default configuration file to stdOut")
	flag.Parse()

	if *defaults {
		stdout := bufio.NewWriter(os.Stdout)
		conf.Pretty(stdout)
		stdout.Flush()
		os.Exit(0)
	} else {
		stdout := bufio.NewWriter(os.Stdout)
		conf.Pretty(stdout)
		stdout.Flush()
		os.Exit(0)
	}
}
