latest-notes
============

(**English** / [Japanese](./README_ja.md))

`latest-notes` is a utility tool designed to streamline the Go release process. It performs two main tasks:

1. **Release Notes Extraction**: Extracts the latest version section from `release_note*.md` to generate a description for GitHub Releases.
2. **Version Source Generation**: Generates a Go source file (e.g., `version.go`) containing the version string, ensuring that `-version` works correctly even when installed via `go install`.

It is intended to be used with the **GitHub CLI (`gh`)** in a `Makefile` to automate the entire release flow.

Supported Release Note Format
-----------------------------

The tool expects the following Markdown structure for `release_note*.md`:

```markdown
v0.0.4
------
Feb 14, 2026  <-- This line (and underlines) are ignored

- Added feature A (#1)
- Fixed bug B (#2)

v0.0.3
------
...
```

### Why this format?

* **Header Matching**: The tool looks for a line matching `^v\d+\.\d+\.\d+$`.
* **Smart Skip**: It automatically skips lines immediately following the header (like dates or `---` underlines) if there is no empty line, ensuring the output starts directly from the actual changes.

Example Makefile Snippet
------------------------

```make
ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=nul
else
    NUL=/dev/null
endif

NAME:=$(notdir $(CURDIR))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)

# 1. Generate version.go before building/releasing
# This ensures 'go install' users see the version with a "-goinstall" suffix.
bump:
	go run github.com/hymkor/latest-notes@latest -gosrc main -suffix "-goinstall" > version.go

# 2. Build and Release
release: bump
	go build -ldflags "-s -w -X main.version=$(VERSION)"
	go run github.com/hymkor/latest-notes@latest | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)
```

## Features

### 1. Smart Extraction

* Reads files matching `release_note*.md`.
* Files containing `ja` in the filename are treated as **Japanese**; others as **English**.
* Identifies version headers using a semantic version pattern (default: `^v\d+\.\d+\.\d+$`).
* Automatically ignores underline markers or date lines immediately following the header to provide clean output.

### 2. Go Source Generation (`-gosrc`)

When the `-gosrc <package_name>` flag is used, the tool outputs a Go source file instead of the release description. This is useful for hardcoding the version into your binary, allowing `go install` to report a valid version.

## Options

| Option | Description |
| --- | --- |
| `-gosrc <name>` | Output Go source code with the specified package name. |
| `-suffix <string>` | Append a suffix to the version (e.g., `-goinstall`). Only works with `-gosrc`. |
| `-pattern <regex>` | Custom regex to identify version headers (default: `^v\d+\.\d+\.\d+$`). |
