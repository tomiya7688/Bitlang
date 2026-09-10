# Bitlang

Bitlang is a strict language and translation pipeline designed around explicit, inspectable lowering stages.

This repository currently contains an early **reference implementation**. The Python code is intentionally small and dependency-free so language semantics can be tested before committing to a production implementation language.

## Planned pipeline

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

Each stage is represented explicitly. A stage consumes one artifact kind and produces another; implicit skipping or reordering is rejected by the reference pipeline.

## Current implementation

- explicit artifact/stage model
- checked pipeline transitions
- case-insensitive identifier canonicalization using Unicode `casefold()`
- collision-aware symbol table preserving original spelling for diagnostics
- minimal CLI for inspecting identifier canonicalization
- dependency-free unit tests

String and character *contents* are not canonicalized. Identifier canonicalization is a symbol-level operation; source lexing will be implemented separately once the concrete grammar is fixed.

## Development

```sh
python -m unittest discover -s tests -v
python -m bitlang.cli canonicalize MyVariable MYVARIABLE myvariable
```

The reference implementation targets Python 3.11+.
