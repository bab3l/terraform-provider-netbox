package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

func TestAccConsoleServerPortTemplateResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	mfgName := testutil.RandomName("tf-test-mfg-rem")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-rem")
	dtModel := testutil.RandomName("tf-test-dt-rem")
	dtSlug := testutil.RandomSlug("tf-test-dt-rem")
	portName := testutil.RandomName("tf-test-cspt-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	testFields := map[string]string{
		"type": "rj-45",
	}

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_console_server_port_template",
		BaseConfig: func() string {
			return testAccConsoleServerPortTemplateResourceConfig_removeOptionalFields_ext_base(mfgName, mfgSlug, dtModel, dtSlug, portName)
		},
		ConfigWithFields: func() string {
			return testAccConsoleServerPortTemplateResourceConfig_removeOptionalFields_ext_withFields(mfgName, mfgSlug, dtModel, dtSlug, portName, testFields)
		},
		OptionalFields: testFields,
		RequiredFields: map[string]string{
			"name": portName,
		},
	})
}

func testAccConsoleServerPortTemplateResourceConfig_removeOptionalFields_ext_base(mfgName, mfgSlug, dtModel, dtSlug, portName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %[5]q
}
`, mfgName, mfgSlug, dtModel, dtSlug, portName)
}

func testAccConsoleServerPortTemplateResourceConfig_removeOptionalFields_ext_withFields(mfgName, mfgSlug, dtModel, dtSlug, portName string, fields map[string]string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_console_server_port_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %[5]q
  type        = %[6]q
}
`, mfgName, mfgSlug, dtModel, dtSlug, portName, fields["type"])
}
