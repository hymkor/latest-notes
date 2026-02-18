package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func printDescription(fname string, rx *regexp.Regexp) error {
	fd, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer fd.Close()

	sc := bufio.NewScanner(fd)
	var version string
	section := 0
	for sc.Scan() {
		line := sc.Text()
		if rx.MatchString(line) {
			if version == "" {
				fmt.Printf("## Changes in %s ", line)
				if strings.Contains(fname, "ja") {
					fmt.Println("(Japanese)")
				} else {
					fmt.Println("(English)")
				}
				version = line
			} else {
				return nil
			}
		} else if version != "" {
			if line == "" {
				section++
			}
			if section > 0 {
				fmt.Println(line)
			}
		}
	}
	return sc.Err()
}

var (
	flagPattern  = flag.String("pattern", `^v\d+\.\d+\.\d+$`, "Regex to identify version headers in the markdown")
	flagGoSource = flag.String("gosrc", "", "Generate Go source code with the version; specify the package name.\nIf not specified, the tool outputs the latest release description.")
	flagSuffix   = flag.String("suffix", "", "Suffix to append to the version string (e.g., \"-goinstall\").\nOnly effective when -gosrc is used.")
)

func mains(args []string) error {
	rxVersion, err := regexp.Compile(*flagPattern)
	if err != nil {
		return err
	}
	if len(args) <= 0 {
		args = []string{"release_note*.md"}
	}

	if *flagGoSource != "" {
		return bump(args, rxVersion)
	}
	for _, arg1 := range args {
		filenames, err := filepath.Glob(arg1)
		if err != nil {
			return err
		}
		for _, fname := range filenames {
			if err := printDescription(fname, rxVersion); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	fmt.Fprintf(os.Stderr, "%s %s-%s-%s\n",
		os.Args[0],
		version,
		runtime.GOOS,
		runtime.GOARCH)
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
