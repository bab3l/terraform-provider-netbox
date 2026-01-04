package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for cluster group resource are in resources_acceptance_tests_customfields package

func TestAccClusterGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group")
	slug := testutil.RandomSlug("tf-test-cluster-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_cluster_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_ClusterGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group-lit")
	slug := testutil.RandomSlug("tf-test-cluster-group-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccClusterGroupConsistencyLiteralNamesConfig(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
				),
			},
		},
	})
}

func TestAccClusterGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-group-id")
	slug := testutil.RandomSlug("tf-test-cluster-group-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
				),
			},
		},
	})
}

func TestAccClusterGroupResource_update(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	name := testutil.RandomName("tf-test-cluster-group")
	slug := testutil.RandomSlug("tf-test-cluster-group")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				Config: testAccClusterGroupResourceConfig_basic(name+"-updated", slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
		},
	})
}

func testAccClusterGroupConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccClusterGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccClusterGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-cluster-group-del")
	slug := testutil.GenerateSlug(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VirtualizationAPI.VirtualizationClusterGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find cluster_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationClusterGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete cluster_group: %v", err)
					}
					t.Logf("Successfully externally deleted cluster_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
