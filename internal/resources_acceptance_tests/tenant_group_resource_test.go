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

func TestAccTenantGroupResource_basic(t *testing.T) {

	t.Parallel()
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-group")
	slug := testutil.RandomSlug("tf-test-tg")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccTenantGroupResource_full(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-group-full")
	slug := testutil.RandomSlug("tf-test-tg-full")
	description := "Test tenant group with all fields"

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", description),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_full(name, slug, description),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_update(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-group-update")
	slug := testutil.RandomSlug("tf-test-tg-upd")
	updatedName := testutil.RandomName("tf-test-tenant-group-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccTenantGroupResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", updatedName),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(updatedName, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_import(t *testing.T) {

	t.Parallel()
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-group-import")
	slug := testutil.RandomSlug("tf-test-tenant-group-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_tenant_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccTenantGroupResourceConfig_import(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTenantGroupResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-id")
	slug := testutil.RandomSlug("tf-test-tg-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccTenantGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

// testAccTenantGroupResourceConfig_basic returns a basic test configuration.
func testAccTenantGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_TenantGroup_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-tenant-group-lit")
	slug := testutil.RandomSlug("tf-test-tenant-group-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", description),
				),
			},
			{
				Config:   testAccTenantGroupConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
				),
			},
		},
	})
}

func testAccTenantGroupConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
} // testAccTenantGroupResourceConfig_full returns a test configuration with all fields.
func testAccTenantGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func testAccTenantGroupResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccTenantGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-group-del")
	slug := testutil.RandomSlug("tf-test-tenant-group-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyTenantGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tenant_group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyTenantGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tenant_group: %v", err)
					}
					t.Logf("Successfully externally deleted tenant_group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTenantGroupResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	groupName := testutil.RandomName("tenant_group")
	groupSlug := testutil.RandomSlug("tenant_group")

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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", groupName),
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", groupSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_tenant_group.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccTenantGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_tenant_group.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags"},
			},
			{
				Config:   testAccTenantGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccTenantGroupResourceImportConfig_full(groupName, groupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["tenancy.tenantgroup"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["tenancy.tenantgroup"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["tenancy.tenantgroup"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["tenancy.tenantgroup"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["tenancy.tenantgroup"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["tenancy.tenantgroup"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["tenancy.tenantgroup"]
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

# Main Resource
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test-value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "test-longtext-value"
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
      value = "2023-01-01"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key" = "value"})
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
`,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		groupName, groupSlug,
	)
}
