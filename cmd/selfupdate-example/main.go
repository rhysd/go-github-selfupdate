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
	latest, err := selfupdate.TryUpdate(version, "go-github-selfupdate", nil)
	if err != nil {
		return err
	}

	previous := semver.Make(version)
	if previous.Equals(latest.Version) {
		fmt.Println("Current binary is the latest version", version)
	} else {
		fmt.Printf(
			`Update successfully done to version %v

Release note:
%s
`,
			latest.Version, latest.Description)
	}
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
