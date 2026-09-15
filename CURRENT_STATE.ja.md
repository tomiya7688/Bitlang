# Bitlang 現在の実装状態

[English](CURRENT_STATE.md) | 日本語

AI 支援開発およびコントリビューター向けの、現在の実装状態を短くまとめた文書です。

## 実装済み
- Go bootstrap module と CUI entrypoint
- 明示的な staged pipeline skeleton
- source text と lexical token の表現
- source location を保持する lexer
- 大文字小文字を区別しない identifier canonicalization の基礎
- duplicate name detection を持つ symbol table
- artifact / stage / transition / orchestration に分割された pipeline 責務
- canonical name / symbol / duplicate error / symbol table に分割された名前管理責務
- Go 実装用 coding rules と responsibility registry
- 簡潔な検出結果、file/function size check、ignore に対応した `go-rule-checker`
- プロジェクト管理の `bitlang-ci` validation gate
- AI context entrypoint、change routing、validation routing

## 現在の Pipeline
- Source: 一部実装
- Preprocessed: 表現・処理とも未実装
- Compiled: stage slot のみ
- TreeObject: stage slot のみ
- VM Assembly: stage slot のみ
- VM / translators: 計画段階

## 既知の一時的制約
- identifier canonicalization は bootstrap の暫定基準として現在 Go の lowercase 挙動を使用
- lexer の identifier character rule は暫定的な ASCII rule
- comment と grammar 固有の multi-character operator は未確定

## 直近の実装
1. concrete Preprocessed artifact representation を導入
2. 未定義の構文を勝手に追加せず、最小限の preprocessor boundary を追加
3. 再利用可能な compiler-core API から stage result を公開
4. compiler-core の stage API が安定してから CUI を拡張

## Source of Truth に関する注意
このファイルは索引です。言語意味論の正本は specification / design source と compiler test に置きます。
