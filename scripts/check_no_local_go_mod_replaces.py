#!/usr/bin/env python3

from __future__ import annotations

import re
import sys
from pathlib import Path


_LOCAL_RHS_PREFIXES = (
    "./",
    "../",
    ".\\",
    "..\\",
    "/",  # absolute posix
    "~/",  # home shorthand (posix shells)
)


def _is_local_path(rhs: str) -> bool:
    rhs = rhs.strip()

    if rhs.startswith(_LOCAL_RHS_PREFIXES):
        return True

    # Windows absolute paths: C:\..., D:\...
    if re.match(r"^[A-Za-z]:\\", rhs):
        return True

    return False


def _find_replace_lines(go_mod_text: str) -> list[tuple[int, str]]:
    """Return (lineno, line) for replace mappings.

    Supports both:
      - replace a => ../b
      - replace (
          a => ../b
        )
    """
    matches: list[tuple[int, str]] = []
    in_replace_block = False
    for lineno, raw in enumerate(go_mod_text.splitlines(), start=1):
        stripped = raw.strip()

        if stripped.startswith("replace ("):
            in_replace_block = True
            continue

        if in_replace_block:
            if stripped == ")":
                in_replace_block = False
                continue
            if "=>" in stripped:
                matches.append((lineno, raw))
            continue

        if stripped.startswith("replace ") and "=>" in stripped:
            matches.append((lineno, raw))

    return matches


def _extract_rhs_from_replace_line(line: str) -> str | None:
    # Handles both single-line and block entries, e.g.:
    #   replace a => ../b
    #   replace a v1.2.3 => ../b
    #   a => ../b
    if "=>" not in line:
        return None
    return line.split("=>", 1)[1].strip()


def check_file(path: Path) -> list[str]:
    errors: list[str] = []
    try:
        text = path.read_text(encoding="utf-8")
    except Exception as exc:  # noqa: BLE001
        return [f"{path}: failed to read: {exc}"]

    for lineno, raw in _find_replace_lines(text):
        rhs = _extract_rhs_from_replace_line(raw)
        if rhs is None:
            continue

        # drop any trailing comments
        rhs = rhs.split("//", 1)[0].strip()

        if _is_local_path(rhs):
            errors.append(
                f"{path}:L{lineno}: local go.mod replace target not allowed: {raw.strip()}"
            )

    return errors


def main(argv: list[str]) -> int:
    files = [Path(a) for a in argv[1:] if a.strip()]
    if not files:
        files = [Path("go.mod")]

    all_errors: list[str] = []
    for f in files:
        if f.name != "go.mod":
            continue
        all_errors.extend(check_file(f))

    if all_errors:
        for err in all_errors:
            print(err)
        print(
            "\nFix: remove local filesystem replace directives from go.mod (they break CI).",
            file=sys.stderr,
        )
        return 1

    return 0


if __name__ == "__main__":
    raise SystemExit(main(sys.argv))
