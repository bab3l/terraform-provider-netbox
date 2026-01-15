#!/usr/bin/env python3
"""
Pre-commit hook to check that go.mod doesn't contain local replace directives.
Local replace directives should only be used during development and should not be committed.
"""
import re
import sys
from pathlib import Path


def check_go_mod(filepath: Path) -> tuple[bool, list[str]]:
    """Check if go.mod contains local replace directives."""
    issues = []

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Pattern to match local replace directives
    # Example: replace github.com/some/module => ../local-path
    local_replace_pattern = re.compile(
        r'^\s*replace\s+\S+\s+=>\s+\.\./.*$',
        re.MULTILINE
    )

    matches = local_replace_pattern.findall(content)

    if matches:
        for match in matches:
            issues.append(f"Local replace directive found: {match.strip()}")

    return len(issues) == 0, issues


def main():
    """Main function to check go.mod files."""
    go_mod_path = Path('go.mod')

    if not go_mod_path.exists():
        print("ERROR: go.mod file not found")
        return 1

    passed, issues = check_go_mod(go_mod_path)

    if not passed:
        print("ERROR: go.mod contains local replace directives that should not be committed:")
        for issue in issues:
            print(f"  - {issue}")
        print("\nLocal replace directives (e.g., => ../local-path) are for development only.")
        print("Please remove them before committing.")
        return 1

    return 0


if __name__ == '__main__':
    sys.exit(main())
