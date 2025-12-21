#!/usr/bin/env python3
"""
Generate LiteralNames consistency tests for remaining low-priority resources.
This script appends properly formatted Go test code to existing test files.
"""

import os

# Template for device bay template test (needs subdevice_role)
DEVICE_BAY_TEMPLATE_TEST = '''

// TestAccConsistency_DeviceBayTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_DeviceBayTemplate_LiteralNames(t *testing.T) {
\tt.Parallel()
\tmanufacturerName := testutil.RandomName("manufacturer")
\tmanufacturerSlug := testutil.RandomSlug("manufacturer")
\tdeviceTypeName := testutil.RandomName("device-type")
\tdeviceTypeSlug := testutil.RandomSlug("device-type")
\tbayName := testutil.RandomName("bay")

\tresource.Test(t, resource.TestCase{
\t\tPreCheck:                 func() { testutil.TestAccPreCheck(t) },
\t\tProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
\t\tSteps: []resource.TestStep{
\t\t\t{
\t\t\t\tConfig: testAccDeviceBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, bayName),
\t\t\t\tCheck: resource.ComposeTestCheckFunc(
\t\t\t\t\tresource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", bayName),
\t\t\t\t\tresource.TestCheckResourceAttr("netbox_device_bay_template.test", "device_type", deviceTypeSlug),
\t\t\t\t),
\t\t\t},
\t\t\t{
\t\t\t\t// Critical: Verify no drift when refreshing state
\t\t\t\tPlanOnly: true,
\t\t\t\tConfig:   testAccDeviceBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, bayName),
\t\t\t},
\t\t},
\t})
}

func testAccDeviceBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, bayName string) string {
\treturn fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.id
  subdevice_role = "parent"
}

resource "netbox_device_bay_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q

  depends_on = [netbox_device_type.test]
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, bayName)
}
'''

# Template for console port template test
CONSOLE_PORT_TEMPLATE_TEST = '''

// TestAccConsistency_ConsolePortTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ConsolePortTemplate_LiteralNames(t *testing.T) {
\tt.Parallel()
\tmanufacturerName := testutil.RandomName("manufacturer")
\tmanufacturerSlug := testutil.RandomSlug("manufacturer")
\tdeviceTypeName := testutil.RandomName("device-type")
\tdeviceTypeSlug := testutil.RandomSlug("device-type")
\tportName := testutil.RandomName("port")

\tresource.Test(t, resource.TestCase{
\t\tPreCheck:                 func() { testutil.TestAccPreCheck(t) },
\t\tProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
\t\tSteps: []resource.TestStep{
\t\t\t{
\t\t\t\tConfig: testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
\t\t\t\tCheck: resource.ComposeTestCheckFunc(
\t\t\t\t\tresource.TestCheckResourceAttr("netbox_console_port_template.test", "name", portName),
\t\t\t\t\tresource.TestCheckResourceAttr("netbox_console_port_template.test", "device_type", deviceTypeSlug),
\t\t\t\t),
\t\t\t},
\t\t\t{
\t\t\t\t// Critical: Verify no drift when refreshing state
\t\t\t\tPlanOnly: true,
\t\t\t\tConfig:   testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
\t\t\t},
\t\t},
\t})
}

func testAccConsolePortTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {
\treturn fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_console_port_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q
  type = "rj-45"

  depends_on = [netbox_device_type.test]
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, portName)
}
'''

def append_test(filepath, test_code):
    """Append test code to a file."""
    with open(filepath, 'a', encoding='utf-8') as f:
        f.write(test_code)
    print(f"✓ Added test to {os.path.basename(filepath)}")

def main():
    base_path = r"c:\GitRoot\terraform-provider-netbox\internal\resources_test"

    tests = [
        (os.path.join(base_path, "device_bay_template_resource_test.go"), DEVICE_BAY_TEMPLATE_TEST),
        (os.path.join(base_path, "console_port_template_resource_test.go"), CONSOLE_PORT_TEMPLATE_TEST),
    ]

    for filepath, test_code in tests:
        if os.path.exists(filepath):
            append_test(filepath, test_code)
        else:
            print(f"✗ File not found: {filepath}")

    print("\nDone! Added 2 template tests.")

if __name__ == "__main__":
    main()
