Release notes
=============

v0.0.5
------
Feb 18, 2026

- Integrated `cmd/bump/main.go` functionality into the main `latest-notes` tool.
- Added `-gosrc` option to output Go source code containing the version constant (e.g., `var version = "v0.0.0(+ -suffix)"`), replacing the need for the separate `bump` tool.

v0.0.4
------
Feb 14, 2026

- Add `cmd/bump/main.go` (#2)
- Changed header: `### Changes in ...` → `## Changes in ...` (#3)

v0.0.3
------
Feb 12, 2026

- Add `-pattern` option

v0.0.2
------
Feb 12, 2026

- Enable to change filenames of release notes to read with arguments

v0.0.1
------
Feb 10, 2026

- The initial version
