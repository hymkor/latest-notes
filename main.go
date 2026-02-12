package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var rxVersion = regexp.MustCompile(`^v\d+\.\d+\.\d$`)

func main1(fname string) error {
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
		if rxVersion.MatchString(line) {
			if version == "" {
				fmt.Printf("### Changes in %s ", line)
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

func mains(args []string) error {
	if len(args) <= 0 {
		args = []string{"release_note*.md"}
	}
	for _, arg1 := range args {
		filenames, err := filepath.Glob(arg1)
		if err != nil {
			return err
		}
		for _, fname := range filenames {
			if err := main1(fname); err != nil {
				return err
			}
		}
	}
	return nil
}

func main() {
	if err := mains(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
