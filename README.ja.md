# Bitlang

[English](README.md) | 日本語

Bitlang は、明示的で検査可能な多段 lowering を中心に設計された、厳格なプログラミング言語および変換パイプラインです。

> **好きなコードを 好きな書き方で**

これは宣伝文句だけではなく、Bitlang の設計原則です。Bitlang 自体も多少冗長ながら書きやすく非常に安全な言語を目指しますが、それだけでなく、他言語や異なる記述スタイルから意味を保ったまま容易に Bitlang へ変換できることを重視します。preprocessor、language adapter、scope 単位の property default、namespace mount、完全明示化は、その自由さと安全性を両立するための仕組みです。

初期実装は **Go** で作られています。これにより、開発中のビルドをすぐに単体実行可能ファイルへでき、ランタイム依存を増やさずクロスコンパイルできます。

Go はあくまで bootstrap 実装言語です。コンパイラ中核は特定言語への依存を避ける方針で実装し、将来的には他言語で同等実装を作成でき、最終的には Bitlang 自身で Bitlang を実装できることを目標にしています。

コンパイラ実装の移植性と self-hosting に関する規約は [CODING_RULES.md](CODING_RULES.md) を参照してください。

## 想定パイプライン

```text
Bitlang source
  -> preprocessor
  -> Bitlang Explicit
  -> static analysis / advisor
  -> Bitlang Low
  -> tree object
  -> Bitlang VM assembly
  -> VM or architecture translator
```

各段階は明示的に表現されます。各 stage は1種類の artifact を受け取り、次の artifact を生成します。現在の core では、暗黙的な stage の飛ばしや並べ替えを許可しません。

中間 stage の境界は、将来的に実装言語をまたぐ契約として扱う予定です。将来 Bitlang で書かれたコンパイラも、Go 固有の実装詳細を引き継がず、同等の artifact を入出力できることを目標とします。

## Bootstrap と self-hosting の目標

想定している発展段階は次の通りです。

```text
Go bootstrap compiler
  -> 安定した Bitlang 意味論と stage format を実装
  -> コンパイラ部品を Bitlang で実装
  -> bootstrap compiler で Bitlang compiler をビルド
  -> Bitlang でビルドされた Bitlang でもう一度 Bitlang をビルド
  -> stage 出力と挙動の等価性を比較
```

Self-hosting は長期的なアーキテクチャ目標であり、不安定なコンパイラ部品を早期に書き直すこと自体は目的ではありません。

## 現在の実装

- Go module と CLI entrypoint
- 明示的な artifact / stage モデル
- 検証付き pipeline transition
- 大文字小文字を区別しない identifier canonicalization
- 診断用に元の綴りを保持する collision-aware symbol table
- 外部依存を必要としない unit tests
- 将来の別実装を考慮した portability-first coding rules

文字列リテラルと文字リテラルの内容は canonicalize しません。identifier canonicalization は symbol level の処理です。具体的な grammar が確定した後、source lexing は別責務として実装します。

Unicode identifier の normalization / case-folding はまだ言語契約に含めていません。現在は一時的かつ決定的な基準として Go standard library の lowercase mapping を使用しています。

## コントリビュート

コントリビューターを歓迎しています。コンパイラ全体を理解している必要はありません。実装修正、テスト、ドキュメント、開発ツール、仕様レビュー、バグ報告はいずれも有用なコントリビューションです。

最初に [CONTRIBUTING.ja.md](CONTRIBUTING.ja.md)（[English](CONTRIBUTING.md)）を読み、[good first issue](https://github.com/tomiya7688/Bitlang/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) を確認してください。

コントリビューター向けの主要ドキュメントは英語版と日本語版を維持します。登録済みの日英ペアについては、片方だけを更新した Pull Request を CI が拒否します。

## ビルド

ネイティブ実行ファイル:

```sh
go build -o bitlang ./cmd/bitlang
```

別の対応ホストから Windows x86-64 向けにビルド:

```sh
GOOS=windows GOARCH=amd64 go build -o bitlang.exe ./cmd/bitlang
```

Windows ARM64:

```sh
GOOS=windows GOARCH=arm64 go build -o bitlang-arm64.exe ./cmd/bitlang
```

## テスト

```sh
go test ./...
```

プロジェクト標準の完全な検証 gate:

```sh
go run ./tools/bitlang-ci/cmd/bitlang-ci
```

## 現在の CLI

```sh
./bitlang canonicalize MyVariable MYVARIABLE myvariable
```

## セキュリティ

セキュリティ上重要な報告は [SECURITY.ja.md](SECURITY.ja.md)（[English](SECURITY.md)）に従ってください。exploit の詳細、credential、機密性のある再現情報を公開 Issue に書かないでください。

リポジトリでは通常の Strict CI とは別に Security CI を実行し、CodeQL、到達可能な Go 脆弱性、Go の静的 security analysis、secret leak、dependency 変更、GitHub Actions の supply-chain risk を検査します。

## ライセンス

このリポジトリで管理される Bitlang および Bitlang 関連資産は、明示的に別記されない限り [MIT License](LICENSE) で提供されます。
