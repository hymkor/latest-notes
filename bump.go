package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func findVersion1(fname string, rx *regexp.Regexp) (string, error) {
	fd, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	sc := bufio.NewScanner(fd)
	for sc.Scan() {
		line := sc.Text()
		if m := rx.FindStringSubmatch(line); len(m) > 0 {
			return m[len(m)-1], nil
		}
	}
	return "", sc.Err()
}

func findVersion(args []string, rx *regexp.Regexp) (string, error) {
	for _, arg := range args {
		notes, err := filepath.Glob(arg)
		if err != nil {
			notes = []string{arg}
		}
		for _, fname := range notes {
			version, err := findVersion1(fname, rx)
			if err != nil {
				return "", err
			}
			if version != "" {
				return version, nil
			}
		}
	}
	return "", errors.New("not found")
}

func makeGoSrc(packageName, version string) {
	version = strings.ReplaceAll(version, `\_`, `_`)
	fmt.Printf("package %s\n\nvar version = %q\n", *flagGoSource, version)
}

func bump(args []string, rx *regexp.Regexp) error {
	version, err := findVersion(args, rx)
	if err != nil {
		return err
	}
	if *flagSuffix != "" {
		version += *flagSuffix
	}
	makeGoSrc(*flagGoSource, version)
	return nil
}
