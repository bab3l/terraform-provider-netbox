package resources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVirtualMachineResource_platform_name_persistence(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("cluster-type")
	clusterTypeSlug := testutil.RandomSlug("cluster-type")
	clusterName := testutil.RandomName("cluster")
	platformName := testutil.RandomName("platform")
	platformSlug := testutil.RandomSlug("platform")
	vmName := testutil.RandomName("vm")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualMachineResourceConfig_platform_name(clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "name", vmName),
					resource.TestCheckResourceAttr("netbox_virtual_machine.test", "platform", platformName),
				),
			},
			{
				// Verify no drift
				PlanOnly: true,
				Config:   testAccVirtualMachineResourceConfig_platform_name(clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName),
			},
		},
	})
}

func testAccVirtualMachineResourceConfig_platform_name(clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_cluster" "test" {
  name = "%[3]s"
  type = netbox_cluster_type.test.id
}

resource "netbox_platform" "test" {
  name = "%[4]s"
  slug = "%[5]s"
}

resource "netbox_virtual_machine" "test" {
  name = "%[6]s"
  cluster = netbox_cluster.test.id
  platform = netbox_platform.test.name
}
`, clusterTypeName, clusterTypeSlug, clusterName, platformName, platformSlug, vmName)
}
