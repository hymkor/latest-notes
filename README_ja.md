latest-notes
============

([English](./README.md) / **Japanese**)

`latest-notes` は、Goプロジェクトのリリース作業を効率化するための個人的な運用ツールです。主に以下の2つの機能を提供します。

1. **リリースノートの抽出**: `release_note*.md` から最新バージョンのセクションのみを切り出し、GitHub Release の Description 用テキストを生成します。
2. **バージョン情報のソースコード生成**: バージョン文字列を埋め込んだ Go のソースファイル（`version.go` など）を生成します。これにより、`go install` でインストールされたバイナリでも正しいバージョンを表示できるようになります。

Makefile で **GitHub CLI (`gh`)** と組み合わせて使用することで、リリース作業をワンコマンドで完結させることを想定しています。

## サポートする形式

本ツールは、`release_note*.md` が以下の構造であることを想定しています。

```markdown
v0.0.4
------
Feb 14, 2026  <-- この行は無視されます

- Added feature A (#1)
- Fixed bug B (#2)

v0.0.3
------
...
```

### なぜこの形式か

* **ヘッダー照合**: `^v\d+\.\d+\.\d+$` にマッチする行を探します。
* **スマートスキップ**: 見出しの直後に空行なしで続く行（日付や下線）は、リリース内容ではないと判断して自動的に除外します。

Makefile での使用例
-------------------

```make
ifeq ($(OS),Windows_NT)
    SHELL=CMD.EXE
    NUL=nul
else
    NUL=/dev/null
endif

NAME:=$(notdir $(CURDIR))
VERSION:=$(shell git describe --tags 2>$(NUL) || echo v0.0.0)

# 1. ビルド・リリース前に version.go を更新する
# これにより go install したバイナリには "-goinstall" サフィックスが付与されます
bump:
	go run github.com/hymkor/latest-notes@latest -gosrc main -suffix "-goinstall" > version.go

# 2. ビルドおよび GitHub Release の作成
release: bump
	go build -ldflags "-s -w -X main.version=$(VERSION)"
	go run github.com/hymkor/latest-notes@latest | gh release create -d --notes-file - -t $(VERSION) $(VERSION) $(wildcard $(NAME)-$(VERSION)-*.zip)

```

主な機能
--------

### 1. リリースノートの自動抽出

* カレントディレクトリにある `release_note*.md` にマッチするファイルを読み込みます。
* ファイル名に `ja` を含むものは日本語、それ以外は英語のリリースノートとして扱います。
* セマンティックバージョンナンバーのみの行（デフォルト: `^v\d+\.\d+\.\d+$`）をセクションの見出しとして認識します。
* 各ファイルで最初に見つかったセクションを「最新バージョン」とみなして抽出します。
* 見出し直後の下線（Markdownのスタイル）やリリース日付などの行は自動的にカットし、クリーンな出力を生成します。

### 2. Goソースコード生成 (`-gosrc`)

`-gosrc <パッケージ名>` オプションを指定すると、リリースノートから取得したバージョン番号を定数として持つ Go のソースコードを出力します。

オプション
----------

| オプション | 説明 |
| --- | --- |
| `-gosrc <name>` | 指定したパッケージ名で Go のソースコードを出力します。 |
| `-suffix <string>` | バージョン文字列に付与する接尾辞（例: `-goinstall`）。`-gosrc` 指定時のみ有効。 |
| `-pattern <regex>` | バージョン見出しを特定するための正規表現（デフォルト: `^v\d+\.\d+\.\d+$`）。 |
