latest-notes
============
([English](./README.md) / **Japanese**)

リリースノートを記したテキストファイル(`./release_note*.md`)から、最新バージョン部分のみ切り出して、GitHub の Release の description に記述するテキストを生成する、個人的な運用ツールです。

Makefile で次のように GitHub Client と組合せて、リリース作業をワンコマンドで行うことを想定します。

```
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
```

1. カレントディレクトリにある `release_note*.md` にマッチするファイルを読み込む。  

2. そのうちファイル名に ja を含むものは日本語のものとみなす。  

3. セマンティックバージョンナンバーのみの行(`^v\d+\.\d+\.\d+$`)をセクションの見出しとみなし、各ファイルごとに最初に登場するセクションを最新バージョンのテキストとみなし、出力する。  
4. ただし、見出しから空行なしで続く行には、見出し用の下線やリリース日付などが含まれるので、これをカットする。そのかわり `### Changes in パージョンナンバー (English)` という見出しを出力には付加する（日本語版ではJapanese と書く）
