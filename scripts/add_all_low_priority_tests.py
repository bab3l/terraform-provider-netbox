#!/usr/bin/env python3
"""
Generate all remaining low-priority LiteralNames consistency tests.
This adds tests for template and component resources.
"""

import os

# Template pattern tests (use device_type attribute)
TEMPLATE_TESTS = {
    "front_port": {
        "file": "front_port_template_resource_test.go",
        "resource": "netbox_front_port_template",
        "test_func": "TestAccConsistency_FrontPortTemplate_LiteralNames",
        "config_func": "testAccFrontPortTemplateConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "8p8c"\n  rear_port = "rear-port"\n  rear_port_position = 1',
        "needs_rear_port_template": True,
    },
    "interface": {
        "file": "interface_template_resource_test.go",
        "resource": "netbox_interface_template",
        "test_func": "TestAccConsistency_InterfaceTemplate_LiteralNames",
        "config_func": "testAccInterfaceTemplateConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "1000base-t"',
    },
    "module_bay": {
        "file": "module_bay_template_resource_test.go",
        "resource": "netbox_module_bay_template",
        "test_func": "TestAccConsistency_ModuleBayTemplate_LiteralNames",
        "config_func": "testAccModuleBayTemplateConsistencyLiteralNamesConfig",
        "extra_attrs": '',
    },
    "power_outlet": {
        "file": "power_outlet_template_resource_test.go",
        "resource": "netbox_power_outlet_template",
        "test_func": "TestAccConsistency_PowerOutletTemplate_LiteralNames",
        "config_func": "testAccPowerOutletTemplateConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "iec-60320-c13"',
    },
    "power_port": {
        "file": "power_port_template_resource_test.go",
        "resource": "netbox_power_port_template",
        "test_func": "TestAccConsistency_PowerPortTemplate_LiteralNames",
        "config_func": "testAccPowerPortTemplateConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "iec-60320-c14"',
    },
    "rear_port": {
        "file": "rear_port_template_resource_test.go",
        "resource": "netbox_rear_port_template",
        "test_func": "TestAccConsistency_RearPortTemplate_LiteralNames",
        "config_func": "testAccRearPortTemplateConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "8p8c"\n  positions = 1',
    },
}

# Component pattern tests (use device attribute)
COMPONENT_TESTS = {
    "console_port": {
        "file": "console_port_resource_test.go",
        "resource": "netbox_console_port",
        "test_func": "TestAccConsistency_ConsolePort_LiteralNames",
        "config_func": "testAccConsolePortConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "rj-45"',
    },
    "console_server_port": {
        "file": "console_server_port_resource_test.go",
        "resource": "netbox_console_server_port",
        "test_func": "TestAccConsistency_ConsoleServerPort_LiteralNames",
        "config_func": "testAccConsoleServerPortConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "rj-45"',
    },
    "device_bay": {
        "file": "device_bay_resource_test.go",
        "resource": "netbox_device_bay",
        "test_func": "TestAccConsistency_DeviceBay_LiteralNames",
        "config_func": "testAccDeviceBayConsistencyLiteralNamesConfig",
        "extra_attrs": '',
    },
    "front_port": {
        "file": "front_port_resource_test.go",
        "resource": "netbox_front_port",
        "test_func": "TestAccConsistency_FrontPort_LiteralNames",
        "config_func": "testAccFrontPortConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "8p8c"\n  rear_port = netbox_rear_port.rear.id\n  rear_port_position = 1',
        "needs_rear_port": True,
    },
    "module_bay": {
        "file": "module_bay_resource_test.go",
        "resource": "netbox_module_bay",
        "test_func": "TestAccConsistency_ModuleBay_LiteralNames",
        "config_func": "testAccModuleBayConsistencyLiteralNamesConfig",
        "extra_attrs": '',
    },
    "power_outlet": {
        "file": "power_outlet_resource_test.go",
        "resource": "netbox_power_outlet",
        "test_func": "TestAccConsistency_PowerOutlet_LiteralNames",
        "config_func": "testAccPowerOutletConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "iec-60320-c13"',
    },
    "power_port": {
        "file": "power_port_resource_test.go",
        "resource": "netbox_power_port",
        "test_func": "TestAccConsistency_PowerPort_LiteralNames",
        "config_func": "testAccPowerPortConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "iec-60320-c14"',
    },
    "rear_port": {
        "file": "rear_port_resource_test.go",
        "resource": "netbox_rear_port",
        "test_func": "TestAccConsistency_RearPort_LiteralNames",
        "config_func": "testAccRearPortConsistencyLiteralNamesConfig",
        "extra_attrs": '\n  type = "8p8c"\n  positions = 1',
    },
}

def generate_template_test(name, info):
    """Generate a template test (uses device_type attribute)."""
    resource_name = name.replace("_", " ").title().replace(" ", "")
    needs_rear = info.get("needs_rear_port_template", False)

    rear_port_resource = ''
    if needs_rear:
        rear_port_resource = '''

resource "netbox_rear_port_template" "rear" {
  device_type = netbox_device_type.test.id
  name        = "rear-port"
  type        = "8p8c"
  positions   = 1
}
'''

    return f'''

// {info["test_func"]} tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func {info["test_func"]}(t *testing.T) {{
\tt.Parallel()
\tmanufacturerName := testutil.RandomName("manufacturer")
\tmanufacturerSlug := testutil.RandomSlug("manufacturer")
\tdeviceTypeName := testutil.RandomName("device-type")
\tdeviceTypeSlug := testutil.RandomSlug("device-type")
\tresourceName := testutil.RandomName("{name}")

\tresource.Test(t, resource.TestCase{{
\t\tPreCheck:                 func() {{ testutil.TestAccPreCheck(t) }},
\t\tProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
\t\tSteps: []resource.TestStep{{
\t\t\t{{
\t\t\t\tConfig: {info["config_func"]}(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
\t\t\t\tCheck: resource.ComposeTestCheckFunc(
\t\t\t\t\tresource.TestCheckResourceAttr("{info["resource"]}.test", "name", resourceName),
\t\t\t\t\tresource.TestCheckResourceAttr("{info["resource"]}.test", "device_type", deviceTypeSlug),
\t\t\t\t),
\t\t\t}},
\t\t\t{{
\t\t\t\t// Critical: Verify no drift when refreshing state
\t\t\t\tPlanOnly: true,
\t\t\t\tConfig:   {info["config_func"]}(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
\t\t\t}},
\t\t}},
\t}})
}}

func {info["config_func"]}(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName string) string {{
\treturn fmt.Sprintf(`

resource "netbox_manufacturer" "test" {{
  name = %q
  slug = %q
}}

resource "netbox_device_type" "test" {{
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}}{rear_port_resource}

resource "{info["resource"]}" "test" {{
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q{info["extra_attrs"]}

  depends_on = [netbox_device_type.test]
}}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, resourceName)
}}
'''

def generate_component_test(name, info):
    """Generate a component test (uses device attribute)."""
    resource_name = name.replace("_", " ").title().replace(" ", "")
    needs_rear = info.get("needs_rear_port", False)

    rear_port_resource = ''
    rear_port_extra_device_param = ''
    if needs_rear:
        rear_port_resource = '''

resource "netbox_rear_port" "rear" {
  device    = %q
  name      = "rear-port"
  type      = "8p8c"
  positions = 1

  depends_on = [netbox_device.test]
}
'''
        rear_port_arg_call = ''  # Not needed - use deviceName directly
        rear_port_arg_def = ''  # Not needed - use deviceName directly
        rear_port_extra_device_param = ', deviceName'  # Need extra deviceName for rear_port's device attribute
    else:
        rear_port_arg_call = ''
        rear_port_arg_def = ''
        rear_port_extra_device_param = ''

    return f'''

// {info["test_func"]} tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func {info["test_func"]}(t *testing.T) {{
\tt.Parallel()
\tmanufacturerName := testutil.RandomName("manufacturer")
\tmanufacturerSlug := testutil.RandomSlug("manufacturer")
\tdeviceTypeName := testutil.RandomName("device-type")
\tdeviceTypeSlug := testutil.RandomSlug("device-type")
\troleName := testutil.RandomName("role")
\troleSlug := testutil.RandomSlug("role")
\tsiteName := testutil.RandomName("site")
\tsiteSlug := testutil.RandomSlug("site")
\tdeviceName := testutil.RandomName("device")
\tresourceName := testutil.RandomName("{name}")

\tresource.Test(t, resource.TestCase{{
\t\tPreCheck:                 func() {{ testutil.TestAccPreCheck(t) }},
\t\tProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
\t\tSteps: []resource.TestStep{{
\t\t\t{{
\t\t\t\tConfig: {info["config_func"]}(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName{rear_port_arg_call}),
\t\t\t\tCheck: resource.ComposeTestCheckFunc(
\t\t\t\t\tresource.TestCheckResourceAttr("{info["resource"]}.test", "name", resourceName),
\t\t\t\t\tresource.TestCheckResourceAttr("{info["resource"]}.test", "device", deviceName),
\t\t\t\t),
\t\t\t}},
\t\t\t{{
\t\t\t\t// Critical: Verify no drift when refreshing state
\t\t\t\tPlanOnly: true,
\t\t\t\tConfig:   {info["config_func"]}(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName{rear_port_arg_call}),
\t\t\t}},
\t\t}},
\t}})
}}

func {info["config_func"]}(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName string{rear_port_arg_def}) string {{
\treturn fmt.Sprintf(`

resource "netbox_manufacturer" "test" {{
  name = %q
  slug = %q
}}

resource "netbox_device_type" "test" {{
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.id
  subdevice_role = "parent"  # Enable device bays
}}

resource "netbox_site" "test" {{
  name = %q
  slug = %q
}}

resource "netbox_device_role" "test" {{
  name = %q
  slug = %q
  color = "ff0000"
}}

resource "netbox_device" "test" {{
  name        = %q
  device_type = netbox_device_type.test.id
  site        = netbox_site.test.id
  role        = netbox_device_role.test.id
  status      = "active"
}}{rear_port_resource}

resource "{info["resource"]}" "test" {{
  # Use literal string name to mimic existing user state
  device = %q
  name = %q{info["extra_attrs"]}

  depends_on = [netbox_device.test]
}}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName{rear_port_extra_device_param}, deviceName, resourceName)
}}
'''

def append_test(filepath, test_code):
    """Append test code to a file."""
    with open(filepath, 'a', encoding='utf-8') as f:
        f.write(test_code)
    print(f"✓ Added test to {os.path.basename(filepath)}")

def main():
    base_path = r"c:\GitRoot\terraform-provider-netbox\internal\resources_test"

    added_count = 0

    # Add template tests
    print("Adding template tests...")
    for name, info in TEMPLATE_TESTS.items():
        filepath = os.path.join(base_path, info["file"])
        if os.path.exists(filepath):
            test_code = generate_template_test(name, info)
            append_test(filepath, test_code)
            added_count += 1
        else:
            print(f"✗ File not found: {filepath}")

    # Add component tests
    print("\nAdding component tests...")
    for name, info in COMPONENT_TESTS.items():
        filepath = os.path.join(base_path, info["file"])
        if os.path.exists(filepath):
            test_code = generate_component_test(name, info)
            append_test(filepath, test_code)
            added_count += 1
        else:
            print(f"✗ File not found: {filepath}")

    print(f"\nDone! Added {added_count} tests total.")
    print(f"  - {len(TEMPLATE_TESTS)} template tests")
    print(f"  - {len(COMPONENT_TESTS)} component tests")

if __name__ == "__main__":
    main()
