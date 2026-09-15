# Documentation Pair Checker

English | [日本語](README.ja.md)

`doc-pair-checker` protects registered English/Japanese documentation pairs.

It always verifies that both files exist and contain a link to the other language. With `--base`, it also checks the Git diff and fails when only one side of a registered pair changed.

## Usage

Repository-state check:

```sh
go run ./tools/doc-pair-checker/cmd/doc-pair-checker
```

Pull-request style synchronization check:

```sh
go run ./tools/doc-pair-checker/cmd/doc-pair-checker --base <base-commit>
```

Registered pairs are intentionally explicit. Technical specifications should be added only after their translation has been reviewed, so an unreviewed translation cannot silently become a second specification.
