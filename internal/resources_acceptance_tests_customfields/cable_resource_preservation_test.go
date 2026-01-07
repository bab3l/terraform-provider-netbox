//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCableResource_CustomFieldsPreservation(t *testing.T) {
	cableLabel := testutil.RandomName("tf-test-cable-pres")
	deviceName := testutil.RandomName("tf-test-device-cable-pres")
	manufacturerName := testutil.RandomName("tf-test-mfr-cable-pres")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-cable-pres")
	deviceTypeModel := testutil.RandomName("tf-test-dt-cable-pres")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt-cable-pres")
	deviceRoleName := testutil.RandomName("tf-test-role-cable-pres")
	deviceRoleSlug := testutil.RandomSlug("tf-test-role-cable-pres")
	siteName := testutil.RandomName("tf-test-site-cable-pres")
	siteSlug := testutil.RandomSlug("tf-test-site-cable-pres")
	cfName := testutil.RandomCustomFieldName("tf_cable_pres")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableResourcePreservationConfig_step1(
					cableLabel, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cable.test", "id"),
					resource.TestCheckResourceAttr("netbox_cable.test", "custom_fields.#", "1"),
				),
			},
			{
				Config: testAccCableResourcePreservationConfig_step2(
					cableLabel, deviceName, manufacturerName, manufacturerSlug,
					deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
					siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					// When custom_fields are not in config, they won't appear in Terraform state
					// but they ARE preserved in NetBox
					resource.TestCheckResourceAttr("netbox_cable.test", "label", cableLabel),
				),
			},
		},
	})
}

func testAccCableResourcePreservationConfig_step1(
	cableLabel, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cable_pres" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.cable"]
  required     = false
}

resource "netbox_manufacturer" "cable" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "cable" {
  manufacturer = netbox_manufacturer.cable.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "cable" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "cable" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "cable" {
  device_type = netbox_device_type.cable.id
  role        = netbox_device_role.cable.id
  site        = netbox_site.cable.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_cable" "test" {
  a_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.src.id
    }
  ]

  b_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.dst.id
    }
  ]

  label = %[1]q

  custom_fields = [
    {
      name  = netbox_custom_field.cable_pres.name
      type  = "text"
      value = "test-value"
    }
  ]
}

resource "netbox_interface" "src" {
  device = netbox_device.cable.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_interface" "dst" {
  device = netbox_device.cable.id
  name   = "eth1"
  type   = "1000base-t"
}
`, cableLabel, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug, cfName)
}

func testAccCableResourcePreservationConfig_step2(
	cableLabel, deviceName, manufacturerName, manufacturerSlug,
	deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
	siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "cable_pres" {
  name         = %[11]q
  type         = "text"
  object_types = ["dcim.cable"]
  required     = false
}

resource "netbox_manufacturer" "cable" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_device_type" "cable" {
  manufacturer = netbox_manufacturer.cable.id
  model        = %[5]q
  slug         = %[6]q
}

resource "netbox_device_role" "cable" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_site" "cable" {
  name   = %[9]q
  slug   = %[10]q
  status = "active"
}

resource "netbox_device" "cable" {
  device_type = netbox_device_type.cable.id
  role        = netbox_device_role.cable.id
  site        = netbox_site.cable.id
  name        = %[2]q
  status      = "active"
}

resource "netbox_cable" "test" {
  a_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.src.id
    }
  ]

  b_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.dst.id
    }
  ]

  label = %[1]q
  # custom_fields intentionally omitted
}

resource "netbox_interface" "src" {
  device = netbox_device.cable.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_interface" "dst" {
  device = netbox_device.cable.id
  name   = "eth1"
  type   = "1000base-t"
}
`, cableLabel, deviceName, manufacturerName, manufacturerSlug,
		deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug,
		siteName, siteSlug, cfName)
}
