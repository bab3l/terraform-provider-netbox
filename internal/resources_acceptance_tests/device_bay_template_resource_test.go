package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceBayTemplateResource_basic(t *testing.T) {
	name := testutil.RandomName("tf-test-dbt")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})
}

func TestAccDeviceBayTemplateResource_full(t *testing.T) {
	name := testutil.RandomName("tf-test-dbt-full")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", "Test Label"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", "Test device bay template with full options"),
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},
		},
	})
}

func TestAccDeviceBayTemplateResource_update(t *testing.T) {
	name := testutil.RandomName("tf-test-dbt-upd")
	manufacturerName := testutil.RandomName("tf-test-mfr")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")
	deviceTypeName := testutil.RandomName("tf-test-dt")
	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
				),
			},
			{
				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", "Updated Label"),
					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccDeviceBayTemplateResourceConfig_basic(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.slug
  subdevice_role = "parent"
}

resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}

func testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.slug
  subdevice_role = "parent"
}

resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  label       = "Test Label"
  description = "Test device bay template with full options"
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}

func testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model          = %q
  slug           = %q
  manufacturer   = netbox_manufacturer.test.slug
  subdevice_role = "parent"
}

resource "netbox_device_bay_template" "test" {
  device_type = netbox_device_type.test.id
  name        = %q
  label       = "Updated Label"
  description = "Updated description"
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name)
}
