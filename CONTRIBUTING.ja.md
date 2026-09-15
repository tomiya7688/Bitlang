# Bitlang へのコントリビュート

[English](CONTRIBUTING.md) | 日本語

Bitlang に興味を持っていただきありがとうございます。

Bitlang は現在も設計・実装を進めている段階です。コントリビューションは大規模である必要はなく、参加前にコンパイラパイプライン全体を理解している必要もありません。

## 歓迎するコントリビューション

例えば次のような領域を歓迎します。

- Go bootstrap compiler の実装
- lexer、preprocessor、static analysis、pipeline
- テストと regression case
- compiler diagnostics と error handling
- 開発ツール
- ドキュメントとサンプル
- 仕様レビューと矛盾・曖昧さの報告
- 各 stage の成熟に応じた VM / assembly / translator の設計・実装

コードを書かない貢献も歓迎します。分かりにくい規則、矛盾する仕様、不足しているテスト、読みにくいドキュメントを見つけることも有用な作業です。

## 最初に確認するもの

初めて参加する場合は、[`good first issue`](https://github.com/tomiya7688/Bitlang/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) が付いた Issue を確認してください。

Issue の内容が曖昧な場合、大きな設計上の仮定を置く前に Issue 上で確認してください。

Bitlang は安定した言語意味論と実装詳細を意図的に分離しています。実装タスクを完了するためだけに新しい構文や意味論を作らないでください。必要な挙動が未定義なら、その曖昧さ自体を報告してください。

## リポジトリの構成

想定している lowering pipeline は次の通りです。

```text
Bitlang source
  -> preprocessor
  -> Bitlang preprocessed
  -> static analysis / advisor
  -> Bitlang compiled
  -> tree object
  -> Bitlang VM assembly
  -> VM or architecture translator
```

主な入口:

- `README.md` / `README.ja.md` — プロジェクト概要と基本的なビルド方法
- `CURRENT_STATE.md` / `CURRENT_STATE.ja.md` — 現在の実装状態の短い一覧
- `SECURITY.md` / `SECURITY.ja.md` — 脆弱性報告とセキュリティ要件
- `bitlang/LANGUAGE_SPEC.md` — 言語仕様
- `CODING_RULES.md` — 移植性とコンパイラ実装規約
- `GO_CODING_RULES.md` — Go bootstrap 実装専用の規約
- `FILE_RESPONSIBILITIES.md` — source file の責務索引

## 開発環境

現在の bootstrap 実装には Go を使用します。

テスト:

```sh
go test ./...
```

ローカルで完全な CI gate を実行:

```sh
go run ./tools/bitlang-ci/cmd/bitlang-ci
```

現在の CLI をビルド:

```sh
go build -o bitlang ./cmd/bitlang
```

変更を提出する前に関連テストが通ることを確認し、変更範囲は対象 Issue に集中させてください。

## コーディング上の期待事項

Go 実装を変更する場合:

- `GO_CODING_RULES.md` に従う
- リポジトリで定義された責務境界を維持する
- panic や silent failure より決定的な error を優先する
- bug fix には regression test を追加する
- 明確な理由がない依存関係を追加しない
- コンパイラの挙動を明示的かつ検査可能に保つ

Go 実装規約は Bitlang 言語仕様ではありません。bootstrap 実装上の都合を利用者向け言語要件にしないでください。

## セキュリティ上の期待事項

security-sensitive な変更や脆弱性報告を行う前に [SECURITY.ja.md](SECURITY.ja.md)（[English](SECURITY.md)）を確認してください。

- credential、token、private key、production secret、実在する personal data を commit しない
- 新規 dependency は明確な理由がある場合に限り、必要最小限にする
- GitHub Actions は immutable commit SHA に固定し、人間向け release version は comment に残す
- 明確な理由なしに workflow token permission を広げない
- authenticated Git write を意図して行う job 以外では checkout credential を保持しない
- untrusted Pull Request のコード実行に `pull_request_target` を使用しない
- security-sensitive な修正には可能な限り regression coverage を追加する
- CI を通すためだけに security check を無効化、弱体化、bypass しない

脆弱性の可能性がある内容を exploit 詳細付きで公開 Issue に投稿しないでください。`SECURITY.ja.md` の private reporting 手順に従ってください。

## ドキュメントの言語方針

コントリビューター向けドキュメントは英語版と日本語版を維持します。

- 基本の `.md` は英語版です。
- `.ja.md` は日本語版です。
- CI に登録された日英ペアを変更する場合、同じ Pull Request で両方を更新します。
- CI は登録済みペアの存在、相互リンク、片方だけが変更されていないことを検証します。
- 翻訳は意味を保持するものであり、別の仕様を新設するものではありません。

現在 CI が強制するペアは `tools/doc-pair-checker` で管理します。技術仕様の翻訳は内容のレビューが済んだものから登録します。未レビューの翻訳を別の source of truth として扱ってはいけません。

## 仕様変更と設計変更

小規模な実装修正は、通常は既存 Issue に対して直接提出できます。

Bitlang の意味論、中間表現、pipeline contract を変更する場合は、先に議論してください。言語挙動については仕様が source of truth であり、実装が暗黙に仕様を書き換えてはいけません。

仕様上の曖昧さを見つけた場合、複数の解釈と問題点を説明する Issue を作ること自体が有効なコントリビューションです。

## Pull Request

有用な Pull Request は原則として次を満たします。

1. 1つの明確な問題を解決する。
2. 関連 Issue がある場合は参照する。
3. 挙動を変更する場合はテストを追加または更新する。
4. 公開契約を変更する場合は登録済みの日英ドキュメントを両方更新する。
5. リポジトリの Strict CI と Security CI を通過する。

大規模で無関係な refactor は、可能な限り機能変更と分離してください。

## ライセンス

このリポジトリへのコントリビューションは、リポジトリの MIT License の下で提供されることに同意したものとして扱います。
