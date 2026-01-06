package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccFHRPGroupDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	protocol := "vrrp2"
	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32
	name := testutil.RandomName("tf-test-fhrp-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckFHRPGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupDataSourceConfig_byID(protocol, groupID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "protocol", protocol),
				),
			},
		},
	})
}

func TestAccFHRPGroupDataSource_byID(t *testing.T) {
	t.Parallel()

	protocol := "vrrp2"
	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32
	name := testutil.RandomName("tf-test-fhrp-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckFHRPGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupDataSourceConfig_byID(protocol, groupID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "name", name),
				),
			},
		},
	})
}

func TestAccFHRPGroupDataSource_byProtocolAndGroupID(t *testing.T) {
	t.Parallel()

	protocol := "hsrp"
	groupID := int32(acctest.RandIntRange(1, 254)) // #nosec G115 -- test value range is 1-254, safe for int32
	name := testutil.RandomName("tf-test-fhrp-ds-lookup")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterFHRPGroupCleanup(protocol, groupID)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckFHRPGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupDataSourceConfig_byProtocolAndGroupID(protocol, groupID, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "protocol", protocol),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "group_id", fmt.Sprintf("%d", groupID)),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group.test", "name", name),
				),
			},
		},
	})
}

func testAccFHRPGroupDataSourceConfig_byID(protocol string, groupID int32, name string) string {
	return fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol = %q
  group_id = %d
  name     = %q
}

data "netbox_fhrp_group" "test" {
  id = netbox_fhrp_group.test.id
}
`, protocol, groupID, name)
}

func testAccFHRPGroupDataSourceConfig_byProtocolAndGroupID(protocol string, groupID int32, name string) string {
	return fmt.Sprintf(`
resource "netbox_fhrp_group" "test" {
  protocol = %q
  group_id = %d
  name     = %q
}

data "netbox_fhrp_group" "test" {
  protocol = netbox_fhrp_group.test.protocol
  group_id = netbox_fhrp_group.test.group_id
}
`, protocol, groupID, name)
}
