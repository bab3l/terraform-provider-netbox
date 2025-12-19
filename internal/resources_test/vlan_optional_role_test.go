package resources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestAccVLANResource_optionalRoleUpdate tests that when role is not specified in config,

// it remains null in state even if the API returns a role value.

// This prevents "inconsistent result after apply" errors.

func TestAccVLANResource_optionalRoleUpdate(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	vlanName := testutil.RandomName("vlan")

	vlanVid := testutil.RandomVID()

	description1 := "Initial description"

	description2 := "Updated description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				// Create VLAN without role

				Config: testAccVLANOptionalRoleConfig(siteName, siteSlug, vlanName, vlanVid, description1),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),

					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVid)),

					resource.TestCheckResourceAttr("netbox_vlan.test", "description", description1),

					resource.TestCheckNoResourceAttr("netbox_vlan.test", "role"),
				),
			},

			{

				// Update description (not role) - role should remain empty/null

				Config: testAccVLANOptionalRoleConfig(siteName, siteSlug, vlanName, vlanVid, description2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),

					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVid)),

					resource.TestCheckResourceAttr("netbox_vlan.test", "description", description2),

					resource.TestCheckNoResourceAttr("netbox_vlan.test", "role"),
				),
			},
		},
	})

}

func testAccVLANOptionalRoleConfig(siteName, siteSlug, vlanName string, vlanVid int32, description string) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = %[1]q

  slug = %[2]q

}



resource "netbox_vlan" "test" {

  name        = %[3]q

  vid         = %[4]d

  site        = netbox_site.test.id

  description = %[5]q

  # role intentionally omitted to test optional attribute handling

}

`, siteName, siteSlug, vlanName, vlanVid, description)

}
