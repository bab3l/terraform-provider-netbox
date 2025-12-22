package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceBayTemplateResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-dbt")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-full")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	label := testutil.RandomName("label")

	description := testutil.RandomName("description")

	updatedLabel := testutil.RandomName("label-upd")

	updatedDescription := testutil.RandomName("description-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", description),

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "device_type"),
				),
			},

			{

				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedLabel, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccDeviceBayTemplateResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-upd")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	updatedLabel := testutil.RandomName("label-upd")

	updatedDescription := testutil.RandomName("description-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

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

				Config: testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedLabel, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", updatedDescription),
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

func testAccDeviceBayTemplateResourceConfig_full(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %[1]q

  slug = %[2]q

}

resource "netbox_device_type" "test" {

  model          = %[3]q

  slug           = %[4]q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}

resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %[5]q

  label       = %[6]q

  description = %[7]q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, description)

}

func TestAccConsistency_DeviceBayTemplate_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-dbt-lit")

	manufacturerName := testutil.RandomName("tf-test-mfr")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	deviceTypeName := testutil.RandomName("tf-test-dt")

	deviceTypeSlug := testutil.RandomSlug("tf-test-dt")

	label := testutil.RandomName("label")

	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)
	cleanup.RegisterDeviceTypeCleanup(deviceTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceBayTemplateConsistencyLiteralNamesConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_device_bay_template.test", "description", description),
				),
			},

			{

				Config: testAccDeviceBayTemplateConsistencyLiteralNamesConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay_template.test", "id"),
				),
			},
		},
	})

}

func testAccDeviceBayTemplateConsistencyLiteralNamesConfig(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %[1]q

  slug = %[2]q

}

resource "netbox_device_type" "test" {

  model          = %[3]q

  slug           = %[4]q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}

resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %[5]q

  label       = %[6]q

  description = %[7]q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, description)

}

func testAccDeviceBayTemplateResourceConfig_updated(name, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, label, description string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %[1]q

  slug = %[2]q

}

resource "netbox_device_type" "test" {

  model          = %[3]q

  slug           = %[4]q

  manufacturer   = netbox_manufacturer.test.slug

  subdevice_role = "parent"

}

resource "netbox_device_bay_template" "test" {

  device_type = netbox_device_type.test.id

  name        = %[5]q

  label       = %[6]q

  description = %[7]q

}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, label, description)

}
