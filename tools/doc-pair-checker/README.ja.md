# Documentation Pair Checker

[English](README.md) | 日本語

`doc-pair-checker` は、登録済みの英語版・日本語版ドキュメントペアを保護します。

常に両方のファイルが存在し、互いの言語版へリンクしていることを検証します。`--base` を指定した場合は Git diff も確認し、登録済みペアの片方だけが変更されている場合に失敗します。

## 使用方法

現在のリポジトリ状態を確認:

```sh
go run ./tools/doc-pair-checker/cmd/doc-pair-checker
```

Pull Request と同様に変更同期を確認:

```sh
go run ./tools/doc-pair-checker/cmd/doc-pair-checker --base <base-commit>
```

対象ペアは意図的に明示登録します。技術仕様は翻訳内容のレビューが済んでから登録し、未レビューの翻訳が暗黙に第二の仕様にならないようにします。
