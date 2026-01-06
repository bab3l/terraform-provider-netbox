package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteTargetDataSource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique name - route targets have 21 char max, use format like "65000:400-<random>"
	name := fmt.Sprintf("65000:400-%s", testutil.RandomSlug("ds")[:8])

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_route_target.test", "name", name),
				),
			},
		},
	})
}

func TestAccRouteTargetDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("65000:401-%s", testutil.RandomSlug("ds-id")[:8])

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_route_target.test", "name", name),
				),
			},
		},
	})
}

func TestAccRouteTargetDataSource_byID(t *testing.T) {
	t.Parallel()

	// Generate unique name - route targets have 21 char max, use format like "65000:400-<random>"
	name := fmt.Sprintf("65000:400-%s", testutil.RandomSlug("ds")[:8])

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetDataSourceConfigByID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_route_target.test", "name", name),
				),
			},
		},
	})
}

func testAccRouteTargetDataSourceConfig(name string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_route_target" "test" {
  name = %q
}

data "netbox_route_target" "test" {
  name = netbox_route_target.test.name
}
`, name)
}

func testAccRouteTargetDataSourceConfigByID(name string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_route_target" "test" {
  name = %q
}

data "netbox_route_target" "test" {
  id = netbox_route_target.test.id
}
`, name)
}
