package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPAddressResource_basic(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", "192.0.2.100/24"),
				),
			},
		},
	})

}

func TestAccIPAddressResource_full(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_full(),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", "10.0.0.50/32"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "test.example.com"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Test IP address"),
				),
			},
		},
	})

}

func TestAccIPAddressResource_update(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", "192.0.2.100/24"),
				),
			},

			{

				Config: testAccIPAddressResourceConfig_full(),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "test.example.com"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Test IP address"),
				),
			},
		},
	})

}

func TestAccIPAddressResource_import(t *testing.T) {

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
				),
			},

			{

				ResourceName: "netbox_ip_address.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccIPAddressResourceConfig_basic() string {

	return `resource "netbox_ip_address" "test" {

  address = "192.0.2.100/24"

}

`

}

func testAccIPAddressResourceConfig_full() string {

	return `resource "netbox_ip_address" "test" {

  address     = "10.0.0.50/32"

  status      = "active"

  dns_name    = "test.example.com"

  description = "Test IP address"

}

`

}
