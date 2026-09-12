# go-rule-checker ignore format

The checker automatically reads `.go-rule-checker-ignore` from the working directory when present.

Supported records:

```text
path generated/**
rule DOC third_party/**
rule NAME legacy.go
```

- `path <pattern>` ignores every rule for matching paths.
- `rule <RULE> <pattern>` ignores one rule for matching paths.
- `#` starts a comment.
- `*` matches within one path segment.
- a pattern ending in `/**` matches that directory recursively.

Known rule codes: `NAME`, `DOC`, `SIZE`, `FILESIZE`, `MAIN`.
