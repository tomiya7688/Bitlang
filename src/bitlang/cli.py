"""Small CLI for the Bitlang reference implementation."""

from __future__ import annotations

import argparse

from .names import canonicalize_identifier


def _build_parser() -> argparse.ArgumentParser:
    parser = argparse.ArgumentParser(prog="bitlang")
    subcommands = parser.add_subparsers(dest="command", required=True)

    canonicalize = subcommands.add_parser(
        "canonicalize", help="show canonical forms of Bitlang identifiers"
    )
    canonicalize.add_argument("identifiers", nargs="+")
    return parser


def main(argv: list[str] | None = None) -> int:
    args = _build_parser().parse_args(argv)

    if args.command == "canonicalize":
        for identifier in args.identifiers:
            print(f"{identifier}\t{canonicalize_identifier(identifier)}")
        return 0

    raise AssertionError(f"unhandled command: {args.command}")


if __name__ == "__main__":
    raise SystemExit(main())
