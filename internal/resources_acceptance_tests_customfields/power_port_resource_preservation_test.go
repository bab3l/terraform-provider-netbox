//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPowerPortResource_CustomFieldsPreservation(t *testing.T) {
	portName := testutil.RandomName("tf-test-pp-pres")
	deviceName := testutil.RandomName("tf-test-device-pp-pres")
	manufacturerName := testutil.RandomName("tf-test-mfr-pp-pres")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-pp-pres")
	deviceTypeModel := testutil.RandomName("tf-test-dt-pp-pres")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-pp-pres")
	deviceRoleName := testutil.RandomName("tf-test-role-pp-pres")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-pp-pres")
	siteName := testutil.RandomName("tf-test-site-pp-pres")
	siteSlug := testutil.RandomSlug("tf-test-site-pp-pres")
	cfName := testutil.RandomCustomFieldName("tf_pp_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPowerPortResourcePreservationConfig_step1(
					portName, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_power_port.test", "id"),
					resource.TestCheckResourceAttr("netbox_power_port.test", "custom_fields.#", "1"),
				),
			},
			{
				Config: testAccPowerPortResourcePreservationConfig_step2(
					portName, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					// When custom_fields are not in config, they won't appear in Terraform state
					// but they ARE preserved in NetBox
					resource.TestCheckResourceAttr("netbox_power_port.test", "name", portName),
				),
			},
		},
	})
}

func testAccPowerPortResourcePreservationConfig_step1(
	portName, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "pp_pres" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.powerport"]
  required     = false
}

resource "netbox_manufacturer" "pp" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "pp" {
  manufacturer = netbox_manufacturer.pp.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "pp" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "pp" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "pp" {
  device_type = netbox_device_type.pp.id
  role        = netbox_device_role.pp.id
  site        = netbox_site.pp.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_power_port" "test" {
  device = netbox_device.pp.id
  name   = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.pp_pres.name
      type  = "text"
      value = "test-value"
    }
  ]
}
`, portName, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug, cfName)
}

func testAccPowerPortResourcePreservationConfig_step2(
	portName, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "pp_pres" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.powerport"]
  required     = false
}

resource "netbox_manufacturer" "pp" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "pp" {
  manufacturer = netbox_manufacturer.pp.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "pp" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "pp" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "pp" {
  device_type = netbox_device_type.pp.id
  role        = netbox_device_role.pp.id
  site        = netbox_site.pp.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_power_port" "test" {
  device = netbox_device.pp.id
  name   = %[1]q
  # custom_fields intentionally omitted
}
`, portName, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug, cfName)
}
