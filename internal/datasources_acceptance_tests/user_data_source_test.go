package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource_IDPreservation(t *testing.T) {

	t.Parallel()
	username := "admin"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_byUsername(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_user.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_user.test", "username", username),
				),
			},
		},
	})
}

func TestAccUserDataSource_byUsername(t *testing.T) {

	t.Parallel()
	// Use the admin user that's always present in NetBox
	username := "admin"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_byUsername(username),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_user.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_user.test", "username", username),
				),
			},
		},
	})
}

func TestAccUserDataSource_byID(t *testing.T) {

	t.Parallel()
	// First get admin user to obtain the ID, then test by ID
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig_byID(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_user.by_id", "id"),
					resource.TestCheckResourceAttr("data.netbox_user.by_id", "username", "admin"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig_byUsername(username string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

data "netbox_user" "test" {
  username = %q
}
`, username)
}

func testAccUserDataSourceConfig_byID() string {
	return `
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {}

# First get the admin user by username to get their ID
data "netbox_user" "admin" {
  username = "admin"
}

# Then look them up by ID
data "netbox_user" "by_id" {
  id = data.netbox_user.admin.id
}
`
}
