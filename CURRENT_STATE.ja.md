# Bitlang 現在の実装状態

[English](CURRENT_STATE.md) | 日本語

AI 支援開発およびコントリビューター向けの、現在の実装状態を短くまとめた文書です。

## 実装済み
- Go bootstrap module と CUI entrypoint
- 明示的な staged pipeline skeleton
- source text と lexical token の表現
- source location と identifier の canonical form を保持する lexer および1-pass lexical preprocessing
- modifier、generic 構文、accessor、未解析 initializer/body token を保持する C# 風 class/member parser
- 明示状態の解決・検証を行う property axis specification
- specification に基づく declaration parser と source modifier expansion の接続口
- 大文字小文字を区別しない identifier canonicalization の基礎
- duplicate name detection を持つ symbol table
- artifact / stage / transition / orchestration に分割された pipeline 責務
- canonical name / symbol / duplicate error / symbol table に分割された名前管理責務
- Go 実装用 coding rules と responsibility registry
- 簡潔な検出結果、file/function size check、ignore に対応した `go-rule-checker`
- プロジェクト管理の `bitlang-ci` validation gate
- AI context entrypoint、change routing、validation routing

## 現在の Pipeline
- Source: source text、lexer、初期 class/member 構文解析を実装済み
- Explicit: 現行の完全明示宣言だけを型付き `ExplicitArtifact` へ解析できる。必須propertyを検証し、元ソース全文を保持する。データ駆動のsource既定値補完は `retention`（`Dynamic`）、`retention_domain`（`Process_retention`）、`instance_access`（`Instance_required`）、解決済み `lifetime` に基づく `finalization_trigger` に対応し、残りの既定値、legacy migration、class/member正規化は未接続・未実装
- Compiled: stage slot のみ
- TreeObject: stage slot のみ
- VM Assembly: stage slot のみ
- VM / translators: 計画段階

## 既知の一時的制約
- identifier canonicalization は bootstrap の暫定基準として現在 Go の lowercase 挙動を使用
- lexer の identifier character rule は暫定的な ASCII rule
- `LexicalPreprocess` は元ソース全文を保持しつつ1-passの字句処理のみを行う。`ApplySourcePropertyDefaults` は無条件既定値と、依存軸の解決順に沿った連鎖条件付き既定値を補完し、条件が未解決なら推測せず残す。ただしpreprocessing pipelineには未接続。宣言位置に依存する既定値、source sugar展開、legacy形式の移行、全propertyの補完、pipeline全体の反復 normalizationは未実装
- class bodyは対応delimiterを保ったtokenとして保持するが、expression/statementの意味解析は未実装
- generic constraint と初期grammar外のdeclaration形式は未対応
- grammar固有のmulti-character operatorは未確定

## 直近の実装
1. 残りのsource既定値と必要な宣言文脈を明確にし、既定値補完を正規化へ接続してから決定的なlegacy migrationを実装する
2. 反復 normalization と決定的な収束・循環検査を実装する
3. adapter固有の modifier 意味を推測せず、source declaration を正規化境界へ接続する
4. CUIを拡張する前に、再利用可能な compiler-core API から完全な stage result を公開する

## Source of Truth に関する注意
このファイルは索引です。言語意味論の正本は specification / design source と compiler test に置きます。
