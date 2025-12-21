package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInventoryItemResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),

					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "device"),
				),
			},
		},
	})

}

func TestAccInventoryItemResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "label", "Inventory Label"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "serial", "SN-12345"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "asset_tag", "ASSET-001"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "description", "Test inventory item"),
				),
			},
		},
	})

}

func TestAccInventoryItemResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "name", name),
				),
			},

			{

				Config: testAccInventoryItemResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "label", "Inventory Label"),

					resource.TestCheckResourceAttr("netbox_inventory_item.test", "serial", "SN-12345"),
				),
			},
		},
	})

}

func TestAccInventoryItemResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-inv-item")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item.test", "id"),
				),
			},

			{

				ResourceName: "netbox_inventory_item.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})

}

func testAccInventoryItemResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_inventory_item" "test" {

  device = netbox_device.test.id

  name   = %q

}

`, testAccInventoryItemResourcePrereqs(name), name)

}

func testAccInventoryItemResourceConfig_full(name string) string {

	return fmt.Sprintf(`

%s

resource "netbox_inventory_item" "test" {

  device      = netbox_device.test.id

  name        = %q

  label       = "Inventory Label"

  serial      = "SN-12345"

  asset_tag   = "ASSET-001"

  description = "Test inventory item"

}

`, testAccInventoryItemResourcePrereqs(name), name)

}

func testAccInventoryItemResourcePrereqs(name string) string {

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
