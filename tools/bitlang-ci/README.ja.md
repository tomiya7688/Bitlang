# Bitlang CI

[English](README.md) | 日本語

`bitlang-ci` は Bitlang コンパイラ実装のためにプロジェクト自身が管理する、厳格な validation gate です。

コンパイラや変換処理の silent defect は、見かけ上有効な出力を生成したままプログラムの意味を変える可能性があります。そのため、通常のアプリケーション CI より意図的に厳しくしています。

## 使用方法

```sh
go run ./tools/bitlang-ci/cmd/bitlang-ci
```

## 現在の gate

以下をすべて実行し、1つでも失敗すれば全体を失敗とします。

- `gofmt` cleanliness
- Bitlang `go-rule-checker` self-check
- 日英ドキュメントペアの存在・相互リンク確認
- `go mod tidy -diff`
- `go vet ./...`
- `go test ./...`
- shuffle/repeat 付き test
- `go build ./...`
- temporary output directory への Bitlang CUI build
- `git diff --check`

GitHub Actions ではさらに次を実行します。

- Linux / Windows / macOS 上で完全な gate
- race detector
- coverage smoke test
- Pull Request での日英ドキュメント変更同期チェック
- Linux / Windows / macOS の amd64 / arm64 向け cross-build

## 失敗時の方針

Bitlang CI は fail-closed です。検査対象の欠落、checker 実行失敗、不正な checker 設定、panic、疑わしい非決定性、登録済み翻訳ペアの更新漏れ、その他の未検証状態を成功として扱ってはいけません。

将来は semantic invariant、fuzzing、property test、golden conversion test、differential test、意味上可能な round-trip test、cross-stage consistency check もプロジェクト管理の検証として追加していきます。
