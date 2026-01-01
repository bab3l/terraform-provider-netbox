#!/usr/bin/env python3
"""
Script to find potential "mode bug" patterns across Terraform provider resources.

The bug pattern is:
1. Field is Optional (not Computed) in schema
2. mapToState function always sets field from API response
3. User doesn't specify field in config
4. Results in "Provider produced inconsistent result after apply" crash

This script analyzes Go resource files to find this pattern.
"""

import os
import re
from pathlib import Path

def find_resources():
    """Find all resource files."""
    resources_dir = Path("internal/resources")
    return list(resources_dir.glob("*_resource.go"))

def analyze_resource(filepath):
    """Analyze a resource file for potential bugs."""
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()

    # Find optional fields in schema
    optional_fields = []
    schema_pattern = r'"(\w+)": schema\.(\w+)Attribute\{[^}]*Optional: true[^}]*\}'

    for match in re.finditer(schema_pattern, content, re.DOTALL):
        field_name = match.group(1)
        field_type = match.group(2)
        optional_fields.append((field_name, field_type))

    # Find mapToState function and check for always-set patterns
    bugs = []
    map_func_pattern = r'func.*mapToState\([^{]*\{.*?\n}'

    for field_name, field_type in optional_fields:
        # Look for patterns like: data.FieldName = types.StringValue(...)
        # without corresponding null handling
        always_set_pattern = rf'data\.{field_name.title()} = types\.{field_type}Value\('
        null_pattern = rf'data\.{field_name.title()} = types\.{field_type}Null\(\)'

        has_always_set = re.search(always_set_pattern, content, re.IGNORECASE)
        has_null_handling = re.search(null_pattern, content, re.IGNORECASE)

        if has_always_set and not has_null_handling:
            bugs.append({
                'field': field_name,
                'type': field_type,
                'pattern': 'always_set_no_null',
                'file': str(filepath)
            })

        # Also look for conditional patterns that might be missing user config check
        conditional_pattern = rf'if.*\w+\.{field_name.title()}\(\) \{{[^}}]*data\.{field_name.title()} = types\.{field_type}Value'
        if re.search(conditional_pattern, content, re.IGNORECASE | re.DOTALL):
            # Check if it has proper user config check
            user_check_pattern = rf'!data\.{field_name.title()}\.IsNull\(\)|data\.{field_name.title()}\.IsUnknown\(\)'
            if not re.search(user_check_pattern, content, re.IGNORECASE):
                bugs.append({
                    'field': field_name,
                    'type': field_type,
                    'pattern': 'conditional_missing_user_check',
                    'file': str(filepath)
                })

    return bugs

def main():
    print("Analyzing Terraform provider resources for optional field bugs...")
    print("=" * 70)

    total_bugs = []
    resources = find_resources()

    for resource_file in resources:
        bugs = analyze_resource(resource_file)
        if bugs:
            print(f"\n📁 {resource_file.name}:")
            for bug in bugs:
                print(f"  🐛 Field: {bug['field']} ({bug['type']}) - Pattern: {bug['pattern']}")
                total_bugs.append(bug)

    print(f"\n" + "=" * 70)
    print(f"Found {len(total_bugs)} potential bugs across {len(resources)} resources")

    # Group by pattern type
    patterns = {}
    for bug in total_bugs:
        pattern = bug['pattern']
        if pattern not in patterns:
            patterns[pattern] = []
        patterns[pattern].append(bug)

    print(f"\nBreakdown by pattern:")
    for pattern, bugs in patterns.items():
        print(f"  {pattern}: {len(bugs)} instances")

    return total_bugs

if __name__ == "__main__":
    bugs = main()
