package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClusterResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-status")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-status")
	clusterName := testutil.RandomName("tf-test-cluster-status")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterDestroy,
		Steps: []resource.TestStep{
			// Step 1: Set status to a non-default value
			{
				Config: testAccClusterResourceConfig_withStatus(clusterTypeName, clusterTypeSlug, clusterName, "offline"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "status", "offline"),
				),
			},
			// Step 2: Remove status from config - should revert to default "active" without drift
			{
				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "status", "active"),
				),
			},
			// Step 3: Ensure no perpetual diff
			{
				PlanOnly: true,
				Config:   testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),
			},
		},
	})
}

func testAccClusterResourceConfig_withStatus(clusterTypeName, clusterTypeSlug, clusterName, status string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster" "test" {
  name   = %q
  type   = netbox_cluster_type.test.id
  status = %q
}
`, clusterTypeName, clusterTypeSlug, clusterName, status)
}
