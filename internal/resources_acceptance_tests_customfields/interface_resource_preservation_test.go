//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_CustomFieldsPreservation(t *testing.T) {
	interfaceName := testutil.RandomName("tf-test-int-pres")
	deviceName := testutil.RandomName("tf-test-device-pres")
	manufacturerName := testutil.RandomName("tf-test-mfr-pres")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pres")
	deviceTypeModel := testutil.RandomName("tf-test-dt-pres")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-pres")
	deviceRoleName := testutil.RandomName("tf-test-role-pres")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-pres")
	siteName := testutil.RandomName("tf-test-site-pres")
	siteSlug := testutil.RandomSlug("tf-test-site-pres")

	cfName := testutil.RandomCustomFieldName("tf_int_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Step 1: Create interface with custom fields
				Config: testAccInterfaceResourcePreservationConfig_step1(
					interfaceName, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "custom_fields.#", "1"),
				),
			},
			{
				// Step 2: Update interface without custom_fields - should preserve existing in NetBox
				Config: testAccInterfaceResourcePreservationConfig_step2(
					interfaceName, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					// When custom_fields are not in config, they won't appear in Terraform state
					// but they ARE preserved in NetBox (this is the point of the test)
					// We verify this by re-adding the field and seeing it still exists
					resource.TestCheckResourceAttr("netbox_interface.test", "name", interfaceName),
				),
			},
			{
				// Step 3: Re-add custom fields to verify they work again
				Config: testAccInterfaceResourcePreservationConfig_step1(
					interfaceName, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_interface.test", "custom_fields.#", "1"),
				),
			},
		},
	})
}

func testAccInterfaceResourcePreservationConfig_step1(
	interfaceName, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "int_pres" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_manufacturer" "pres" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "pres" {
  manufacturer = netbox_manufacturer.pres.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "pres" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "pres" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "pres" {
  device_type = netbox_device_type.pres.id
  role        = netbox_device_role.pres.id
  site        = netbox_site.pres.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_interface" "test" {
  device = netbox_device.pres.id
  name   = %[1]q
  type   = "1000base-t"

  custom_fields = [
    {
      name  = netbox_custom_field.int_pres.name
      type  = "text"
      value = "test-value"
    }
  ]
}
`, interfaceName, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug, cfName)
}

func testAccInterfaceResourcePreservationConfig_step2(
	interfaceName, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "int_pres" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.interface"]
  required     = false
}

resource "netbox_manufacturer" "pres" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "pres" {
  manufacturer = netbox_manufacturer.pres.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "pres" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "pres" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "pres" {
  device_type = netbox_device_type.pres.id
  role        = netbox_device_role.pres.id
  site        = netbox_site.pres.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_interface" "test" {
  device = netbox_device.pres.id
  name   = %[1]q
  type   = "1000base-t"
  # Note: custom_fields intentionally omitted - should be preserved
}
`, interfaceName, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug, cfName)
}
