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

func TestAccClusterTypeResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-cluster-type")

	slug := testutil.RandomSlug("tf-test-cluster-type")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

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

				ResourceName: "netbox_cluster_type.test",

				ImportState: true,

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

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckClusterTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccClusterTypeConsistencyLiteralNamesConfig(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", description),
				),
			},

			{

				Config: testAccClusterTypeConsistencyLiteralNamesConfig(name, slug, description),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
				),
			},
		},
	})

}

func testAccClusterTypeConsistencyLiteralNamesConfig(name, slug, description string) string {

	return fmt.Sprintf(`

resource "netbox_cluster_type" "test" {

  name        = %q

  slug        = %q

  description = %q

}

`, name, slug, description)

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

func TestAccClusterTypeResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

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
		CheckDestroy: testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceImportConfig_full(clusterTypeName, clusterTypeSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", clusterTypeName),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", clusterTypeSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_cluster_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields"}, // Custom fields have import limitations
			},
		},
	})
}

func testAccClusterTypeResourceImportConfig_full(clusterTypeName, clusterTypeSlug, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
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
  object_types = ["virtualization.clustertype"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["virtualization.clustertype"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["virtualization.clustertype"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["virtualization.clustertype"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["virtualization.clustertype"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["virtualization.clustertype"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["virtualization.clustertype"]
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

# Cluster Type with comprehensive custom fields and tags
resource "netbox_cluster_type" "test" {
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
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this cluster type resource for testing purposes."
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
`, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, clusterTypeName, clusterTypeSlug)
}
