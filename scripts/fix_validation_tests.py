#!/usr/bin/env python3
"""Fix validation test syntax to use proper Go struct format."""

import re
import sys

def fix_validation_test(content):
    """Convert inline Config strings to functions and rename ErrorPattern to ExpectedError."""

    # Pattern to match Config field with inline string
    pattern = r'Config:\s+`([^`]+)`,'

    def replace_config(match):
        config_content = match.group(1)
        return f'''Config: func() string {{
\t\t\t\treturn `{config_content}`
\t\t\t\t}},'''

    # Replace Config fields
    content = re.sub(pattern, replace_config, content, flags=re.DOTALL)

    # Replace ErrorPattern with ExpectedError
    content = content.replace('ErrorPattern:', 'ExpectedError:')

    return content

def main():
    files = [
        r'internal\resources_acceptance_tests\rack_resource_test.go',
        r'internal\resources_acceptance_tests\device_resource_test.go',
        r'internal\resources_acceptance_tests\interface_resource_test.go',
        r'internal\resources_acceptance_tests\virtual_machine_resource_test.go',
        r'internal\resources_acceptance_tests\tenant_resource_test.go',
    ]

    for filepath in files:
        try:
            with open(filepath, 'r', encoding='utf-8') as f:
                content = f.read()

            fixed_content = fix_validation_test(content)

            with open(filepath, 'w', encoding='utf-8', newline='\n') as f:
                f.write(fixed_content)

            print(f"✓ Fixed {filepath}")
        except Exception as e:
            print(f"✗ Error fixing {filepath}: {e}", file=sys.stderr)
            return 1

    print("\nAll files fixed successfully!")
    return 0

if __name__ == '__main__':
    sys.exit(main())
