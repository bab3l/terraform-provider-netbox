package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccDeviceResource_fieldConsistency tests that optional fields like airflow, tags, and custom_fields
// maintain consistency between plan and apply, addressing the "Provider produced inconsistent result after apply" bug.
func TestAccDeviceResource_fieldConsistency(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts
	deviceName := testutil.RandomName("tf-test-device-consistency")
	manufacturerName := testutil.RandomName("tf-test-manufacturer")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeModel := testutil.RandomName("tf-test-device-type")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")
	deviceRoleName := testutil.RandomName("tf-test-device-role")
	deviceRoleSlug := testutil.RandomSlug("tf-test-dr")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Test creating device without specifying optional fields like airflow, tags, custom_fields
				Config: testAccDeviceResourceConfig_minimalOptionalFields(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel,
					deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					resource.TestCheckResourceAttrSet("netbox_device.test", "device_type"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "role"),
					resource.TestCheckResourceAttrSet("netbox_device.test", "site"),
					resource.TestCheckResourceAttr("netbox_device.test", "status", "active"),
					// These fields should be handled consistently (null or computed default)
					// The exact values don't matter as much as consistency
				),
			},
			{
				// Test with empty sets for tags and custom_fields - should remain empty, not become null
				Config: testAccDeviceResourceConfig_emptySets(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel,
					deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "name", deviceName),
					// tags and custom_fields should remain as empty sets, not null
					resource.TestCheckResourceAttr("netbox_device.test", "tags.#", "0"),
					resource.TestCheckResourceAttr("netbox_device.test", "custom_fields.#", "0"),
				),
			},
			{
				// Test with explicit airflow value
				Config: testAccDeviceResourceConfig_explicitAirflow(
					deviceName, manufacturerName, manufacturerSlug, deviceTypeModel,
					deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device.test", "id"),
					resource.TestCheckResourceAttr("netbox_device.test", "airflow", "front-to-rear"),
				),
			},
		},
	})
}

func testAccDeviceResourceConfig_minimalOptionalFields(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  # Deliberately omitting airflow, tags, custom_fields to test consistency
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceResourceConfig_emptySets(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name         = %[9]q
  device_type  = netbox_device_type.test.id
  role         = netbox_device_role.test.id
  site         = netbox_site.test.id
  tags         = []
  custom_fields = []
  # Testing empty sets should remain empty, not become null
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}

func testAccDeviceResourceConfig_explicitAirflow(deviceName, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %[3]q
  slug         = %[4]q
}

resource "netbox_device_role" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_site" "test" {
  name   = %[7]q
  slug   = %[8]q
  status = "active"
}

resource "netbox_device" "test" {
  name        = %[9]q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
  airflow     = "front-to-rear"
}
`, manufacturerName, manufacturerSlug, deviceTypeModel, deviceTypeSlug, deviceRoleName, deviceRoleSlug, siteName, siteSlug, deviceName)
}
