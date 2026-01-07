//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConsolePortResource_CustomFieldsPreservation(t *testing.T) {
	portName := testutil.RandomName("tf-test-cp-pres")
	deviceName := testutil.RandomName("tf-test-device-cp-pres")
	manufacturerName := testutil.RandomName("tf-test-mfr-cp-pres")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-cp-pres")
	deviceTypeModel := testutil.RandomName("tf-test-dt-cp-pres")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-cp-pres")
	deviceRoleName := testutil.RandomName("tf-test-role-cp-pres")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-cp-pres")
	siteName := testutil.RandomName("tf-test-site-cp-pres")
	siteSlug := testutil.RandomSlug("tf-test-site-cp-pres")
	cfName := testutil.RandomCustomFieldName("tf_cp_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConsolePortResourcePreservationConfig_step1(
					portName, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_console_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_console_port.test", "custom_fields.#", "1"),
				),
			},
			{
				Config: testAccConsolePortResourcePreservationConfig_step2(
					portName, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_console_port.test", "custom_fields.#", "1"),
				),
			},
		},
	})
}

func testAccConsolePortResourcePreservationConfig_step1(
	portName, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cp_pres" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.consoleport"]
  required     = false
}

resource "netbox_manufacturer" "cp" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "cp" {
  manufacturer = netbox_manufacturer.cp.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "cp" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "cp" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "cp" {
  device_type = netbox_device_type.cp.id
  role        = netbox_device_role.cp.id
  site        = netbox_site.cp.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_console_port" "test" {
  device = netbox_device.cp.id
  name   = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.cp_pres.name
      type  = "text"
      value = "test-value"
    }
  ]
}
`, portName, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug, cfName)
}

func testAccConsolePortResourcePreservationConfig_step2(
	portName, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug string,
) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "cp" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "cp" {
  manufacturer = netbox_manufacturer.cp.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "cp" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "cp" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "cp" {
  device_type = netbox_device_type.cp.id
  role        = netbox_device_role.cp.id
  site        = netbox_site.cp.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_console_port" "test" {
  device = netbox_device.cp.id
  name   = %[1]q
  # custom_fields intentionally omitted
}
`, portName, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug)
}
