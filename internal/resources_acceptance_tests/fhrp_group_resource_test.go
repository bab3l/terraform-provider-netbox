package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFHRPGroupResource_basic(t *testing.T) {
	protocol := "vrrp2"
	groupID := int32(acctest.RandIntRange(1, 254)) // nolint:gosec

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
		},
	})
}

func TestAccFHRPGroupResource_full(t *testing.T) {
	protocol := "hsrp"
	groupID := int32(acctest.RandIntRange(1, 254)) // nolint:gosec
	name := testutil.RandomName("tf-test-fhrp")
	description := "Test FHRP Group with all fields"
	authType := "plaintext"
	authKey := "secretkey123"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_full(protocol, groupID, name, description, authType, authKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", description),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_type", authType),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_key", authKey),
				),
			},
		},
	})
}

func TestAccFHRPGroupResource_update(t *testing.T) {
	protocol := "vrrp3"
	groupID := int32(acctest.RandIntRange(1, 254)) // nolint:gosec
	updatedName := testutil.RandomName("tf-test-fhrp-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
				),
			},
			{
				Config: testAccFHRPGroupResourceConfig_full(protocol, groupID, updatedName, "Updated description", "md5", "newsecret456"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_type", "md5"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "auth_key", "newsecret456"),
				),
			},
		},
	})
}

func TestAccFHRPGroupResource_import(t *testing.T) {
	protocol := "vrrp2"
	groupID := int32(acctest.RandIntRange(1, 254)) // nolint:gosec

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupResourceConfig_basic(protocol, groupID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
				),
			},
			{
				ResourceName:      "netbox_fhrp_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccFHRPGroupResourceConfig_basic(protocol string, groupID int32) string {
	return fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol = %q
  group_id = %d
}
`, protocol, groupID)
}

func testAccFHRPGroupResourceConfig_full(protocol string, groupID int32, name, description, authType, authKey string) string {
	return fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol    = %q
  group_id    = %d
  name        = %q
  description = %q
  auth_type   = %q
  auth_key    = %q
}
`, protocol, groupID, name, description, authType, authKey)
}
