package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceRoleResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-device-role")

	slug := testutil.RandomSlug("tf-test-dr")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccDeviceRoleResource_IDPreservation(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("dr-id")

	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		CheckDestroy: testutil.CheckDeviceRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccDeviceRoleResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-device-role-full")

	slug := testutil.RandomSlug("tf-test-dr-full")

	description := testutil.RandomName("description")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceRoleResourceConfig_full(name, slug, description, "aa1409", false),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_role.test", "description", description),

					resource.TestCheckResourceAttr("netbox_device_role.test", "color", "aa1409"),
				),
			},
		},
	})

}

func TestAccDeviceRoleResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-device-role-update")

	slug := testutil.RandomSlug("tf-test-dr-upd")

	updatedName := testutil.RandomName("tf-test-device-role-updated")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
				),
			},

			{

				Config: testAccDeviceRoleResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", updatedName),
				),
			},
		},
	})

}

func testAccDeviceRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_device_role" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccDeviceRoleResourceConfig_full(name, slug, description, color string, vmRole bool) string {

	return fmt.Sprintf(`

resource "netbox_device_role" "test" {

  name        = %q

  slug        = %q

  description = %q

  color       = %q

}

`, name, slug, description, color)

}

func TestAccConsistency_DeviceRole_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-device-role-lit")
	slug := testutil.RandomSlug("tf-test-dr-lit")
	description := testutil.RandomName("description")
	color := testutil.Color

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceRoleConsistencyLiteralNamesConfig(name, slug, description, color),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_device_role.test", "description", description),
					resource.TestCheckResourceAttr("netbox_device_role.test", "color", color),
				),
			},
			{
				Config:   testAccDeviceRoleConsistencyLiteralNamesConfig(name, slug, description, color),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),
				),
			},
		},
	})
}

func testAccDeviceRoleConsistencyLiteralNamesConfig(name, slug, description, color string) string {
	return fmt.Sprintf(`
resource "netbox_device_role" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
}
`, name, slug, description, color)
}
