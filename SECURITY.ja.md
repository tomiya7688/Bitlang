# セキュリティポリシー

[English](SECURITY.md) | 日本語

Bitlang では、コンパイラの正しさ、ビルドの完全性、コントリビューター経由の supply-chain 安全性をセキュリティ上の問題として扱います。

## サポート対象

Bitlang は現在活発な開発段階です。安定版 release が公開されるまでは、セキュリティ修正は現在の `main` branch を対象に行います。過去の commit や開発 snapshot は個別にはサポートしません。

## 脆弱性の報告

exploit の詳細や機密性のある再現手順を公開 Issue に投稿しないでください。

利用可能な場合は、このリポジトリの GitHub private vulnerability reporting / Security Advisories を使用してください。private reporting が利用できない場合は、管理者へ private な報告経路を求めるための最小限の公開 Issue のみ作成し、その Issue に脆弱性の詳細、秘密情報、exploit code、影響を受けた credential を書かないでください。

有用な報告には次の情報を含めてください。

- 影響を受ける component と revision
- 影響と現実的な攻撃前提条件
- private に共有可能な最小限の再現情報
- compiler host、生成プログラム、VM、build environment、CI、repository credential のどこへ影響し得るか
- 分かる範囲の mitigation 案

実際の credential や無関係な private data を報告に含めないでください。

## セキュリティ対象範囲

以下を含む領域を security-sensitive として扱います。

- hostile input に対する parser、lexer、preprocessor、compiler、translator の処理
- command injection、path traversal、安全でない temporary file、意図しない file access
- crafted source による panic、resource exhaustion、denial of service
- unsafe または意味的に異なる挙動を黙って生成する compiler transformation
- 実装が進んだ段階での VM isolation と host interaction
- dependency vulnerability、dependency confusion、supply-chain risk
- GitHub Actions workflow injection、過大な token permission、未固定の Action
- source、Git history、log、test fixture、build artifact からの secret leak
- release / build provenance と改ざんリスク

通常の correctness bug が自動的に脆弱性になるわけではありません。ただし trust boundary を越えたり、悪用可能な出力を生成したりする compiler bug は security-sensitive として扱います。

## リポジトリのセキュリティ gate

通常の Strict CI とは別に Security CI を使用します。主な検証は次の通りです。

- `security-extended` query suite を使用した CodeQL
- 到達可能な Go 脆弱性を検査する `govulncheck`
- `gosec` による静的 security analysis
- Gitleaks による secret scanning
- Pull Request の dependency 変更に対する dependency review
- `zizmor` による GitHub Actions と supply-chain risk の検査
- `actionlint` による workflow validation
- GitHub Actions の immutable commit SHA 固定
- Go module と GitHub Actions に対する Dependabot monitoring

Security check は fail-closed を原則とします。scanner 自体の失敗を「検査成功」として扱いません。

## コントリビューター向けセキュリティ要件

次の規則を守ってください。

- 明確な必要性とレビュー可能な理由なしに dependency を追加しない
- 理由を説明せず GitHub Actions permission を広げない
- GitHub Actions は immutable commit SHA に固定し、人間向け version は comment に残す
- Git write が明示的に必要な job 以外では checkout credential を保持しない
- credential、token、private key、production secret、実在する personal data を commit しない
- security-sensitive な bug fix には可能な限り regression test を追加する
- untrusted Pull Request のコード実行に `pull_request_target` を使用しない
- CI を通すためだけに security check を弱めたり bypass したりしない。原因を修正するか、限定的な例外を文書化してレビューする

## 公開と修正

Security report はまず private に検証します。利用者が update できる前に不必要な詳細を公開しないようにし、fix または mitigation の提供後に公開 disclosure を行えます。希望があり適切な場合は reporter へ credit を付与します。

現在、固定された応答時間 SLA は設定していませんが、信頼できる security report は高優先度で扱います。
