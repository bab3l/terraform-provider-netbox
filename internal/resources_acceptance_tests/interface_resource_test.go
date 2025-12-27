package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),

					resource.TestCheckResourceAttrSet("netbox_interface.test", "device"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),

					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "false"),

					resource.TestCheckResourceAttr("netbox_interface.test", "mtu", "1500"),

					resource.TestCheckResourceAttr("netbox_interface.test", "mgmt_only", "true"),

					resource.TestCheckResourceAttr("netbox_interface.test", "description", "Test interface with full options"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "true"),
				),
			},

			{

				Config: testAccInterfaceResourceConfig_updated(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "enabled", "false"),

					resource.TestCheckResourceAttr("netbox_interface.test", "description", "Updated interface description"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_interface.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})

}

func TestAccConsistency_Interface(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-consistency")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_consistency_device_name(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})

}

func TestAccConsistency_Interface_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-interface-literal")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInterfaceResourceConfig_consistency_device_id(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),

					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),

					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})

}

func TestAccInterfaceResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-interface-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccInterfaceResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
					resource.TestCheckResourceAttr("netbox_interface.test", "type", "1000base-t"),
				),
			},
		},
	})

}

func testAccInterfaceResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device = netbox_device.test.id

  name   = %q

  type   = "1000base-t"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_full(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device       = netbox_device.test.id

  name         = %q

  type         = "1000base-t"

  enabled      = false

  mtu          = 1500

  mgmt_only    = true

  description  = "Test interface with full options"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_updated(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device      = netbox_device.test.id

  name        = %q

  type        = "1000base-t"

  enabled     = false

  description = "Updated interface description"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_consistency_device_name(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device = netbox_device.test.name

  name   = %q

  type   = "1000base-t"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourceConfig_consistency_device_id(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_interface" "test" {

  device = netbox_device.test.id

  name   = %q

  type   = "1000base-t"

}

`, testAccInterfaceResourcePrereqs(name), name)

}

func testAccInterfaceResourcePrereqs(name string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %q

  slug = %q

}

resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}

resource "netbox_device_type" "test" {

  manufacturer = netbox_manufacturer.test.id

  model        = %q

  slug         = %q

}

resource "netbox_device_role" "test" {

  name = %q

  slug = %q

}

resource "netbox_device" "test" {

  site        = netbox_site.test.id

  name        = %q

  device_type = netbox_device_type.test.id

  role        = netbox_device_role.test.id

  status      = "offline"

}

`, name+"-site", testutil.RandomSlug("site"), name+"-mfr", testutil.RandomSlug("mfr"), name+"-model", testutil.RandomSlug("device"), name+"-role", testutil.RandomSlug("role"), name+"-device")

}

func TestAccInterfaceResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)
	name := testutil.RandomName("tf-test-interface-ext-del")
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_site" "test" { name = "%[1]s-site"; slug = "%[2]s" }
resource "netbox_manufacturer" "test" { name = "%[1]s-mfr"; slug = "%[3]s" }
resource "netbox_device_type" "test" { model = "%[1]s-model"; slug = "%[4]s"; manufacturer_id = netbox_manufacturer.test.id }
resource "netbox_device_role" "test" { name = "%[1]s-role"; slug = "%[5]s" }
resource "netbox_device" "test" { site_id = netbox_site.test.id; name = "%[1]s-device"; device_type_id = netbox_device_type.test.id; role_id = netbox_device_role.test.id }
resource "netbox_interface" "test" { name = "%[1]s"; device_id = netbox_device.test.id; type = "1000base-t" }
`, name, testutil.RandomSlug("site"), testutil.RandomSlug("mfr"), testutil.RandomSlug("device"), testutil.RandomSlug("role")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
					resource.TestCheckResourceAttr("netbox_interface.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimInterfacesList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find interface for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimInterfacesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete interface: %v", err)
					}
					t.Logf("Successfully externally deleted interface with ID: %d", itemID)
				},
				Config: fmt.Sprintf(`
resource "netbox_site" "test" { name = "%[1]s-site"; slug = "%[2]s" }
resource "netbox_manufacturer" "test" { name = "%[1]s-mfr"; slug = "%[3]s" }
resource "netbox_device_type" "test" { model = "%[1]s-model"; slug = "%[4]s"; manufacturer_id = netbox_manufacturer.test.id }
resource "netbox_device_role" "test" { name = "%[1]s-role"; slug = "%[5]s" }
resource "netbox_device" "test" { site_id = netbox_site.test.id; name = "%[1]s-device"; device_type_id = netbox_device_type.test.id; role_id = netbox_device_role.test.id }
resource "netbox_interface" "test" { name = "%[1]s"; device_id = netbox_device.test.id; type = "1000base-t" }
`, name, testutil.RandomSlug("site"), testutil.RandomSlug("mfr"), testutil.RandomSlug("device"), testutil.RandomSlug("role")),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_interface.test", "id"),
				),
			},
		},
	})
}
