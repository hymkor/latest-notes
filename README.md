latest-notes
============
(**English** / [Japanese](./README_ja.md))

`latest-notes` is a personal utility tool that extracts only the latest version section
from release note text files (`./release_note*.md`) and generates text suitable for the
*description* field of a GitHub Release.

It is intended to be used together with the GitHub CLI in a Makefile, so that the release
process can be completed with a single command.

Example Makefile snippet:

```make
ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=nul
else
    NUL=/dev/null
endif
NAME:=$(notdir $(CURDIR))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)

release:
	go run github.com/hymkor/latest-notes@latest | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)
````

## Behavior

1. Reads all files in the current directory that match `release_note*.md`.

2. Files whose names contain `ja` are treated as Japanese release notes; all others are
   treated as non-Japanese.

3. A line that consists only of a semantic version number
   (`^v\d+\.\d+\.\d+$`) is treated as a section header.
   For each file, the *first* such section is considered the latest version, and its
   contents are extracted.

4. Lines that immediately follow the version header *without an empty line* are ignored,
   as they typically contain underline markers or release dates.
   Instead, a generated heading such as:

   ```
   ### Changes in <version> (English)
   ```

   is prepended to the output.
   For Japanese files, `(Japanese)` is used instead of `(English)`.
