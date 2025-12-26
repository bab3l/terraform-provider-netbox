package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceTypeResource_basic(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("tf-test-device-type")

	slug := testutil.RandomSlug("tf-test-dt")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "1"),
				),
			},
		},
	})

}

func TestAccDeviceTypeResource_IDPreservation(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("dt-id")
	slug := testutil.GenerateSlug(model)
	manufacturerName := testutil.RandomName("mfr-dt")
	manufacturerSlug := testutil.GenerateSlug(manufacturerName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckDeviceTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "1"),
				),
			},
		},
	})

}

func TestAccDeviceTypeResource_full(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("tf-test-device-type-full")

	slug := testutil.RandomSlug("tf-test-dt-full")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "manufacturer"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "part_number", "TEST-PART-001"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "2"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "is_full_depth", "true"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "airflow", "front-to-rear"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Test device type with full options"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "comments", "Test comments for device type"),
				),
			},
		},
	})

}

func TestAccDeviceTypeResource_update(t *testing.T) {

	t.Parallel()

	model := testutil.RandomName("tf-test-device-type-update")

	slug := testutil.RandomSlug("tf-test-dt-upd")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	updatedModel := testutil.RandomName("tf-test-device-type-updated")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
				),
			},

			{

				Config: testAccDeviceTypeResourceConfig_updated(updatedModel, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", updatedModel),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "4"),
				),
			},
		},
	})

}

func testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

}

`, manufacturerName, manufacturerSlug, model, slug)

}

func testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer      = netbox_manufacturer.test.id

  model             = %q

  slug              = %q

  part_number       = "TEST-PART-001"

  u_height          = 2

  is_full_depth     = true

  airflow           = "front-to-rear"

  description       = "Test device type with full options"

  comments          = "Test comments for device type"

}

`, manufacturerName, manufacturerSlug, model, slug)

}

func testAccDeviceTypeResourceConfig_updated(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

  u_height     = 4

  description  = "Updated description"

}

`, manufacturerName, manufacturerSlug, model, slug)

}

func TestAccConsistency_DeviceType_LiteralNames(t *testing.T) {
	t.Parallel()
	model := testutil.RandomName("tf-test-device-type-lit")
	slug := testutil.RandomSlug("tf-test-dt-lit")
	manufacturerName := testutil.RandomName("tf-test-manufacturer-lit")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.ComposeCheckDestroy(testutil.CheckDeviceTypeDestroy, testutil.CheckManufacturerDestroy),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_device_type.test", "description", description),
				),
			},
			{
				Config:   testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceTypeConsistencyLiteralNamesConfig(model, slug, manufacturerName, manufacturerSlug, description string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  description  = %q
}
`, manufacturerName, manufacturerSlug, model, slug, description)
}
