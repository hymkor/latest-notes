package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"path/filepath"
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
					fmt.Println("(Enlish)")
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

func mains() error {
	filenames, err := filepath.Glob("release_note*.md")
	if err != nil {
		return err
	}
	for _, fname := range filenames {
		if err := main1(fname); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
