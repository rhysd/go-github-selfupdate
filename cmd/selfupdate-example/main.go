package main

import (
	"flag"
	"fmt"
	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"os"
)

const version = "1.2.3"

func selfUpdate() error {
	selfupdate.EnableLog()

	previous := semver.MustParse(version)
	latest, err := selfupdate.UpdateSelf(previous, "rhysd/go-github-selfupdate")
	if err != nil {
		return err
	}

	if previous.Equals(latest.Version) {
		fmt.Println("Current binary is the latest version", version)
	} else {
		fmt.Println("Update successfully done to version", latest.Version)
		fmt.Println("Release note:\n", latest.Description)
	}
	return nil
}

func usage() {
	fmt.Fprintln(os.Stderr, "Usage: selfupdate-example [flags]")
	flag.PrintDefaults()
}

func main() {
	help := flag.Bool("help", false, "Show this help")
	ver := flag.Bool("version", false, "Show version")
	update := flag.Bool("selfupdate", false, "Try go-github-selfupdate via GitHub")

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

	if *update {
		if err := selfUpdate(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	usage()
}
