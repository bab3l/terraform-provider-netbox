package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPAddressResource_basic(t *testing.T) {

	t.Parallel()

	ip := fmt.Sprintf("192.0.%d.%d/24", 100+acctest.RandIntRange(0, 50), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(ip),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip),
				),
			},
		},
	})

}

func TestAccIPAddressResource_full(t *testing.T) {

	t.Parallel()

	ip := fmt.Sprintf("10.0.%d.%d/32", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_full(ip),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "test.example.com"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Test IP address"),
				),
			},
		},
	})

}

func TestAccIPAddressResource_update(t *testing.T) {

	t.Parallel()

	ip1 := fmt.Sprintf("172.16.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	ip2 := fmt.Sprintf("172.16.%d.%d/24", acctest.RandIntRange(0, 255), acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(ip1),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", ip1),
				),
			},

			{

				Config: testAccIPAddressResourceConfig_full(ip2),

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

	t.Parallel()

	ip := fmt.Sprintf("203.0.113.%d/32", acctest.RandIntRange(1, 254))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(ip),

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

func testAccIPAddressResourceConfig_basic(address string) string {

	return fmt.Sprintf(`resource "netbox_ip_address" "test" {

  address = %q

}

`, address)

}

func testAccIPAddressResourceConfig_full(address string) string {

	return fmt.Sprintf(`resource "netbox_ip_address" "test" {

  address     = %q

  status      = "active"

  dns_name    = "test.example.com"

  description = "Test IP address"

}

`, address)

}

func TestAccConsistency_IPAddress_LiteralNames(t *testing.T) {
	t.Parallel()
	address := "10.0.0.1/24"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPAddressConsistencyLiteralNamesConfig(address),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
				),
			},
			{
				Config:   testAccIPAddressConsistencyLiteralNamesConfig(address),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),
				),
			},
		},
	})
}

func testAccIPAddressConsistencyLiteralNamesConfig(address string) string {
	return fmt.Sprintf(`resource "netbox_ip_address" "test" {
  address = %q
}
`, address)
}
