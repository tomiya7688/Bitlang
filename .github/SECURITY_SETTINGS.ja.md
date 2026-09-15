# リポジトリのセキュリティ設定

[English](SECURITY_SETTINGS.md) | 日本語

この文書は、リポジトリ内のファイルだけでは強制できない GitHub 側のセキュリティ設定を記録します。

この公開リポジトリで GitHub が提供している場合、リポジトリ管理者は次を有効に維持してください。

- `main` への変更は Pull Request 経由を必須にする
- merge 前に `Strict gate` と `Security gate` の成功を必須にする
- `main` への force push と branch deletion を禁止する
- 実用上可能な範囲で、merge 前に conversation の解決を必須にする
- default `GITHUB_TOKEN` permission を read-only にし、必要な job だけへ write permission を付与する
- Dependabot alerts と Dependabot security updates を有効にする
- Secret Scanning と Push Protection を有効にする
- private vulnerability reporting / repository Security Advisories を有効にする
- CodeQL code scanning を有効に保ち、新しい alert を確認する
- 設定が利用できる場合、信頼されていない初回 fork contributor の workflow 実行前に承認を必要とする
- allowed-actions policy を確認し、実用上可能な範囲で third-party Actions を制限する

リポジトリ内の workflow では、GitHub Actions を immutable commit SHA に固定し、job ごとに最小権限 permission を使用します。

これらの設定変更は security-sensitive な構成変更として扱い、`SECURITY.ja.md` に照らしてレビューしてください。
