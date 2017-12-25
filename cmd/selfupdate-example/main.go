package main

import (
	"flag"
	"fmt"
	"os"
)

const version = "1.2.3"

func selfUpdate() error {
	return nil
	up, err := selfupdate.TryUpdate(version, "go-github-selfupdate", nil)
	if err != nil {
		return err
	}
	if up.Version == version {
		fmt.Println("Current binary is the latest version", version)
	} else {
		fmt.Println("Update successfully done to version", up.Version)
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: selfupdate-example [flags]")
	flag.PrintDefaults()
}

func main() {
	help := flag.Bool("help", false, "Show this help")
	ver := flag.Bool("version", false, "Show version")
	selfupdate := flag.Bool("selfupdate", false, "Try go-github-selfupdate via GitHub")

	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	if *ver {
		fmt.Println(version)
		os.Exit(0)
	}

	if *selfupdate {
		if err := selfUpdate(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	usage()
}
