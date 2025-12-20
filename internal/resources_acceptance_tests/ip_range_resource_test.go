package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", "192.0.2.10"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", "192.0.2.20"),
				),
			},
		},
	})
}

func TestAccIPRangeResource_full(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", "10.0.0.1"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", "10.0.0.254"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", "Test IP range"),
				),
			},
		},
	})
}

func TestAccIPRangeResource_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", "192.0.2.10"),
				),
			},
			{
				Config: testAccIPRangeResourceConfig_full(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", "Test IP range"),
				),
			},
		},
	})
}

func TestAccIPRangeResource_import(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ip_range.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccIPRangeResourceConfig_basic() string {
	return `
resource "netbox_ip_range" "test" {
  start_address = "192.0.2.10"
  end_address   = "192.0.2.20"
}
`
}

func testAccIPRangeResourceConfig_full() string {
	return `
resource "netbox_ip_range" "test" {
  start_address = "10.0.0.1"
  end_address   = "10.0.0.254"
  status        = "active"
  description   = "Test IP range"
}
`
}
