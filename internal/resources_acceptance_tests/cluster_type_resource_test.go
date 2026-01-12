package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for cluster type resource are in resources_acceptance_tests_customfields package

func TestAccClusterTypeResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-type")
	slug := testutil.RandomSlug("tf-test-cluster-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccClusterTypeResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("ct-id")
	slug := testutil.RandomSlug("ct-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccClusterTypeResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-type-full")
	slug := testutil.RandomSlug("tf-test-cluster-type-full")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", description),
				),
			},
		},
	})
}

func TestAccClusterTypeResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-type-update")
	slug := testutil.RandomSlug("tf-test-cluster-type-update")
	updatedName := testutil.RandomName("tf-test-cluster-type-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
				),
			},
			{
				Config: testAccClusterTypeResourceConfig_full(updatedName, slug, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", "Updated description"),
				),
			},
		},
	})
}

func TestAccClusterTypeResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-type-import")
	slug := testutil.RandomSlug("tf-test-cluster-type-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_cluster_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_ClusterType_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-type-lit")
	slug := testutil.RandomSlug("tf-test-cluster-type-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", description),
				),
			},
			{
				Config:   testAccClusterTypeResourceConfig_full(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
				),
			},
		},
	})
}

func testAccClusterTypeResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccClusterTypeResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccClusterTypeResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-cluster-type-del")
	slug := testutil.GenerateSlug(name)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VirtualizationAPI.VirtualizationClusterTypesList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find cluster_type for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationClusterTypesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete cluster_type: %v", err)
					}
					t.Logf("Successfully externally deleted cluster_type with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccClusterTypeResource_removeDescription(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-type-desc")
	slug := testutil.RandomSlug("tf-test-cluster-type-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_cluster_type",
		BaseConfig: func() string {
			return testAccClusterTypeResourceConfig_basic(name, slug)
		},
		ConfigWithFields: func() string {
			return testAccClusterTypeResourceConfig_full(
				name,
				slug,
				"Test description",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
		},
		CheckDestroy: testutil.CheckClusterTypeDestroy,
	})
}
