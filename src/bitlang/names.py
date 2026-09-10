"""Identifier handling for Bitlang.

Bitlang identifiers are case-insensitive.  This module deliberately operates on
already-tokenized identifiers; string and character literal contents must never
be passed through this canonicalizer.
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Generic, Iterator, TypeVar

T = TypeVar("T")


def canonicalize_identifier(name: str) -> str:
    """Return the canonical comparison form of a Bitlang identifier.

    Unicode ``casefold`` is used instead of ``lower`` because Bitlang's rule is
    semantic case-insensitivity rather than ASCII-only presentation folding.
    The original spelling should be retained separately for diagnostics.
    """

    if not isinstance(name, str):
        raise TypeError("identifier must be str")
    if not name:
        raise ValueError("identifier must not be empty")
    return name.casefold()


@dataclass(frozen=True, slots=True)
class CanonicalName:
    """An identifier with both diagnostic spelling and comparison spelling."""

    spelling: str
    canonical: str

    @classmethod
    def from_spelling(cls, spelling: str) -> "CanonicalName":
        return cls(spelling=spelling, canonical=canonicalize_identifier(spelling))


@dataclass(frozen=True, slots=True)
class Symbol(Generic[T]):
    name: CanonicalName
    value: T


class DuplicateSymbolError(ValueError):
    """Raised when two spellings resolve to the same Bitlang identifier."""

    def __init__(self, existing: str, incoming: str) -> None:
        super().__init__(
            f"duplicate symbol: {incoming!r} conflicts with existing {existing!r}"
        )
        self.existing = existing
        self.incoming = incoming


class SymbolTable(Generic[T]):
    """Case-insensitive symbol table that preserves declared spelling."""

    def __init__(self) -> None:
        self._symbols: dict[str, Symbol[T]] = {}

    def define(self, name: str, value: T) -> Symbol[T]:
        parsed = CanonicalName.from_spelling(name)
        previous = self._symbols.get(parsed.canonical)
        if previous is not None:
            raise DuplicateSymbolError(previous.name.spelling, parsed.spelling)

        symbol = Symbol(name=parsed, value=value)
        self._symbols[parsed.canonical] = symbol
        return symbol

    def get(self, name: str) -> T:
        canonical = canonicalize_identifier(name)
        try:
            return self._symbols[canonical].value
        except KeyError:
            raise KeyError(name) from None

    def find(self, name: str) -> Symbol[T] | None:
        return self._symbols.get(canonicalize_identifier(name))

    def __contains__(self, name: object) -> bool:
        if not isinstance(name, str) or not name:
            return False
        return name.casefold() in self._symbols

    def __len__(self) -> int:
        return len(self._symbols)

    def __iter__(self) -> Iterator[Symbol[T]]:
        return iter(self._symbols.values())
