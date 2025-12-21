package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccInterfaceResource_basic(t *testing.T) {

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
			},
		},
	})

}

func TestAccConsistency_Interface(t *testing.T) {

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
