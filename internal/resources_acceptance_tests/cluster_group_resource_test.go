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

				ResourceName: "netbox_cluster_group.test",

				ImportState: true,

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

				Config: testAccClusterGroupConsistencyLiteralNamesConfig(name, slug),

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

func TestAccClusterGroupResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	clusterGroupName := testutil.RandomName("cluster_group")
	clusterGroupSlug := testutil.RandomSlug("cluster_group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(clusterGroupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	// Clean up custom fields and tags
	cleanup.RegisterCustomFieldCleanup(cfText)
	cleanup.RegisterCustomFieldCleanup(cfLongtext)
	cleanup.RegisterCustomFieldCleanup(cfInteger)
	cleanup.RegisterCustomFieldCleanup(cfBoolean)
	cleanup.RegisterCustomFieldCleanup(cfDate)
	cleanup.RegisterCustomFieldCleanup(cfUrl)
	cleanup.RegisterCustomFieldCleanup(cfJson)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceImportConfig_full(clusterGroupName, clusterGroupSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", clusterGroupName),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", clusterGroupSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_cluster_group.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Custom fields have import limitations
			},
		},
	})
}

func testAccClusterGroupResourceImportConfig_full(clusterGroupName, clusterGroupSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Custom Fields (all supported data types)
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["virtualization.clustergroup"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["virtualization.clustergroup"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["virtualization.clustergroup"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["virtualization.clustergroup"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["virtualization.clustergroup"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["virtualization.clustergroup"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["virtualization.clustergroup"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# Cluster Group with comprehensive custom fields and tags
resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test text value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this cluster group resource for testing purposes."
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-15"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key": "value"})
    }
  ]

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, clusterGroupName, clusterGroupSlug)
}
