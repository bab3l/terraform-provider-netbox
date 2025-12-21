#!/usr/bin/env python3
"""
Split datasource acceptance tests from acceptance_test.go into individual files.
"""

import re
from pathlib import Path

test_file = Path("c:/GitRoot/terraform-provider-netbox/internal/datasources_acceptance_tests/acceptance_test.go")
output_dir = Path("c:/GitRoot/terraform-provider-netbox/internal/datasources_acceptance_tests")

# Read the content
content = test_file.read_text()

# Define test mappings: (test_name, datasource_name)
tests = [
    ("TestAccSiteDataSource_basic", "site"),
    ("TestAccTenantDataSource_basic", "tenant"),
    ("TestAccSiteGroupDataSource_basic", "site_group"),
    ("TestAccTenantGroupDataSource_basic", "tenant_group"),
    ("TestAccManufacturerDataSource_basic", "manufacturer"),
    ("TestAccPlatformDataSource_basic", "platform"),
    ("TestAccRegionDataSource_basic", "region"),
    ("TestAccLocationDataSource_basic", "location"),
    ("TestAccRackDataSource_basic", "rack"),
    ("TestAccRackRoleDataSource_basic", "rack_role"),
    ("TestAccDeviceRoleDataSource_basic", "device_role"),
    ("TestAccDeviceTypeDataSource_basic", "device_type"),
    ("TestAccRouteTargetDataSource_basic", "route_target"),
    ("TestAccVirtualDiskDataSource_basic", "virtual_disk"),
    ("TestAccASNRangeDataSource_basic", "asn_range"),
    ("TestAccDeviceBayTemplateDataSource_basic", "device_bay_template"),
]

imports = '''package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)
'''

# Extract each test function
for test_name, ds_name in tests:
    # Find the test function
    pattern = rf'(func {test_name}.*?(?=\nfunc |$))'
    match = re.search(pattern, content, re.DOTALL)

    if not match:
        print(f"ERROR: Could not find {test_name}")
        continue

    # Get the test function and its config function
    test_func = match.group(1).strip()

    # Find the corresponding config function
    config_name = f"testAcc{test_name[7:-6]}Config"  # Remove TestAcc and _basic
    config_pattern = rf'(func {config_name}.*?(?=\nfunc |$))'
    config_match = re.search(config_pattern, content, re.DOTALL)

    if not config_match:
        print(f"WARNING: Could not find {config_name}")
        config_func = ""
    else:
        config_func = config_match.group(1).strip()

    # Clean up excessive newlines
    test_func_clean = re.sub(r'\n\n+', '\n\n', test_func)
    config_func_clean = re.sub(r'\n\n+', '\n\n', config_func) if config_func else ""

    # Create the file content
    file_content = f"{imports}\n\n{test_func_clean}\n\n{config_func_clean}\n"

    # Write to file
    output_file = output_dir / f"{ds_name}_data_source_test.go"
    output_file.write_text(file_content)
    print(f"Created: {output_file.name}")

print("Done! All tests have been split into individual files.")
