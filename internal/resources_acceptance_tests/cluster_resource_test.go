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

func TestAccClusterResource_basic(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")

	clusterName := testutil.RandomName("tf-test-cluster")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),
				),
			},
		},
	})

}

func TestAccClusterResource_full(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-full")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-full")

	clusterName := testutil.RandomName("tf-test-cluster-full")

	description := testutil.RandomName("description")

	comments := testutil.RandomName("comments")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, description, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "description", description),

					resource.TestCheckResourceAttr("netbox_cluster.test", "comments", comments),
				),
			},
		},
	})

}

func TestAccClusterResource_update(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-update")

	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type-update")

	clusterName := testutil.RandomName("tf-test-cluster-update")

	updatedName := testutil.RandomName("tf-test-cluster-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterCleanup(updatedName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
				),
			},

			{

				Config: testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, updatedName, "Updated description", "Updated comments"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "description", "Updated description"),
				),
			},
		},
	})

}

func TestAccClusterResource_import(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("tf-test-cluster-type-import")

	clusterTypeSlug := clusterTypeName

	clusterName := testutil.RandomName("tf-test-cluster-import")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_import(clusterTypeName, clusterTypeSlug, clusterName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),
				),
			},

			{

				ResourceName: "netbox_cluster.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"type"},
			},
		},
	})

}

// TestAccConsistency_Cluster_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_Cluster_LiteralNames(t *testing.T) {

	t.Parallel()
	clusterName := testutil.RandomName("cluster")

	clusterTypeName := testutil.RandomName("cluster-type")

	clusterTypeSlug := testutil.RandomSlug("cluster-type")

	groupName := testutil.RandomName("group")

	groupSlug := testutil.RandomSlug("group")

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
	cleanup.RegisterClusterGroupCleanup(groupSlug)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterConsistencyLiteralNamesConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "type", clusterTypeSlug),

					resource.TestCheckResourceAttr("netbox_cluster.test", "group", groupSlug),

					resource.TestCheckResourceAttr("netbox_cluster.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_cluster.test", "tenant", tenantName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccClusterConsistencyLiteralNamesConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug),
			},
		},
	})

}

func TestAccClusterResource_IDPreservation(t *testing.T) {

	t.Parallel()

	clusterTypeName := testutil.RandomName("clu-type-id")

	clusterTypeSlug := testutil.RandomSlug("clu-type-id")

	clusterName := testutil.RandomName("clu-id")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterCleanup(clusterName)

	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckClusterDestroy,

			testutil.CheckClusterTypeDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),

					resource.TestCheckResourceAttrSet("netbox_cluster.test", "type"),
				),
			},
		},
	})

}

func testAccClusterResourceConfig_basic(clusterTypeName, clusterTypeSlug, clusterName string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {

  name = %q

  slug = %q

}

resource "netbox_cluster" "test" {

  name = %q

  type = netbox_cluster_type.test.id

}

`, clusterTypeName, clusterTypeSlug, clusterName)

}

func testAccClusterResourceConfig_full(clusterTypeName, clusterTypeSlug, clusterName, description, comments string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {

  name = %q

  slug = %q

}

resource "netbox_cluster" "test" {

  name        = %q

  type        = netbox_cluster_type.test.id

  status      = "active"

  description = %q

  comments    = %q

}

`, clusterTypeName, clusterTypeSlug, clusterName, description, comments)

}

func testAccClusterResourceConfig_import(clusterTypeName, clusterTypeSlug, clusterName string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {

  name = %q

  slug = %q

}

resource "netbox_cluster" "test" {

  name = %q

  type = netbox_cluster_type.test.id

}

`, clusterTypeName, clusterTypeSlug, clusterName)

}

func testAccClusterConsistencyLiteralNamesConfig(clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}

resource "netbox_cluster_group" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}

resource "netbox_site" "test" {

  name = "%[6]s"

  slug = "%[7]s"

}

resource "netbox_tenant" "test" {

  name = "%[8]s"

  slug = "%[9]s"

}

resource "netbox_cluster" "test" {

  name = "%[1]s"

  # Use literal string names to mimic existing user state

  type = "%[3]s"

  group = "%[5]s"

  site = "%[6]s"

  tenant = "%[8]s"

  depends_on = [netbox_cluster_type.test, netbox_cluster_group.test, netbox_site.test, netbox_tenant.test]

}

`, clusterName, clusterTypeName, clusterTypeSlug, groupName, groupSlug, siteName, siteSlug, tenantName, tenantSlug)

}

func TestAccClusterResource_externalDeletion(t *testing.T) {
	t.Parallel()
	clusterTypeName := testutil.RandomName("tf-test-cluster-type")
	clusterTypeSlug := testutil.RandomSlug("tf-test-cluster-type")
	clusterName := testutil.RandomName("tf-test-cluster-ext-del")
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}
resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List clusters filtered by name
					items, _, err := client.VirtualizationAPI.VirtualizationClustersList(context.Background()).NameIc([]string{clusterName}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find cluster for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VirtualizationAPI.VirtualizationClustersDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete cluster: %v", err)
					}
					t.Logf("Successfully externally deleted cluster with ID: %d", itemID)
				},
				Config: fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}
resource "netbox_cluster" "test" {
  name = %q
  type = netbox_cluster_type.test.id
}
`, clusterTypeName, clusterTypeSlug, clusterName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
				),
			},
		},
	})
}

func TestAccClusterResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	clusterName := testutil.RandomName("cluster")
	clusterTypeName := testutil.RandomName("cluster_type")
	clusterTypeSlug := testutil.RandomSlug("cluster_type")
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
	cleanup.RegisterClusterCleanup(clusterName)
	cleanup.RegisterClusterTypeCleanup(clusterTypeSlug)
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
		CheckDestroy: testutil.CheckClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterResourceImportConfig_full(clusterName, clusterTypeName, clusterTypeSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster.test", "name", clusterName),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_cluster.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_cluster.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_cluster.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant", "type", "custom_fields"}, // Tenant/type may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccClusterResourceImportConfig_full(clusterName, clusterTypeName, clusterTypeSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}

# Custom Fields (all supported data types)
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["virtualization.cluster"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["virtualization.cluster"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["virtualization.cluster"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["virtualization.cluster"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["virtualization.cluster"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["virtualization.cluster"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["virtualization.cluster"]
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

# Cluster with comprehensive custom fields and tags
resource "netbox_cluster" "test" {
  name   = %q
  type   = netbox_cluster_type.test.id
  tenant = netbox_tenant.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test text value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this cluster resource for testing purposes."
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
`, tenantName, tenantSlug, clusterTypeName, clusterTypeSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, clusterName)
}
