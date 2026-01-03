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

func TestAccTenantResource_basic(t *testing.T) {

	t.Parallel()
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant")
	slug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup to ensure resource is deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccTenantResource_full(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-full")
	slug := testutil.RandomSlug("tf-test-tenant-full")
	description := testutil.RandomName("description")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant.test", "description", description),
				),
			},
		},
	})
}

func TestAccTenantResource_update(t *testing.T) {

	t.Parallel()
	// Generate unique names
	name := testutil.RandomName("tf-test-tenant-update")
	slug := testutil.RandomSlug("tf-test-tenant-upd")
	updatedName := testutil.RandomName("tf-test-tenant-updated")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
				),
			},
			{
				Config: testAccTenantResourceConfig_basic(updatedName, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", updatedName),
				),
			},
		},
	})
}

// testAccTenantResourceConfig_basic returns a basic test configuration.
func testAccTenantResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

// testAccTenantResourceConfig_full returns a test configuration with all fields.
func testAccTenantResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccTenantResource_import(t *testing.T) {

	t.Parallel()
	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tenant-import")
	slug := testutil.RandomSlug("tf-test-tenant-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_import(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_tenant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccTenantResourceConfig_import(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccTenantResourceConfig_import(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func TestAccConsistency_Tenant(t *testing.T) {

	t.Parallel()
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTenantGroupCleanup(groupSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", tenantName),
					resource.TestCheckResourceAttr("netbox_tenant.test", "group", groupName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug),
			},
		},
	})
}

func testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
  group = netbox_tenant_group.test.name
}
`, groupName, groupSlug, tenantName, tenantSlug)
}

func TestAccConsistency_Tenant_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-tenant-lit")
	slug := testutil.RandomSlug("tf-test-tenant-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tenant.test", "description", description),
				),
			},
			{
				Config:   testAccTenantConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
				),
			},
		},
	})
}

func TestAccTenantResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-id")
	slug := testutil.RandomSlug("tf-test-tenant-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
		},
	})
}

func testAccTenantConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccTenantResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tenant-del")
	slug := testutil.RandomSlug("tf-test-tenant-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckTenantDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTenantResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.TenancyAPI.TenancyTenantsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tenant for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.TenancyAPI.TenancyTenantsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tenant: %v", err)
					}
					t.Logf("Successfully externally deleted tenant with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTenantResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	tenantGroupName := testutil.RandomName("tenant_group")
	tenantGroupSlug := testutil.RandomSlug("tenant_group")

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
				Config: testAccTenantResourceImportConfig_full(tenantName, tenantSlug, tenantGroupName, tenantGroupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),
					resource.TestCheckResourceAttr("netbox_tenant.test", "name", tenantName),
					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", tenantSlug),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_tenant.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_tenant.test", "tags.#", "2"),
				),
			},
			{
				Config:                  testAccTenantResourceImportConfig_full(tenantName, tenantSlug, tenantGroupName, tenantGroupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				ResourceName:            "netbox_tenant.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"custom_fields", "tags", "group"},
			},
			{
				Config:   testAccTenantResourceImportConfig_full(tenantName, tenantSlug, tenantGroupName, tenantGroupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug),
				PlanOnly: true,
			},
		},
	})
}

func testAccTenantResourceImportConfig_full(tenantName, tenantSlug, tenantGroupName, tenantGroupSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug string) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant_group" "test" {
  name = %q
  slug = %q
}

# Custom Fields
resource "netbox_custom_field" "cf_text" {
  name        = %q
  type        = "text"
  object_types = ["tenancy.tenant"]
}

resource "netbox_custom_field" "cf_longtext" {
  name        = %q
  type        = "longtext"
  object_types = ["tenancy.tenant"]
}

resource "netbox_custom_field" "cf_integer" {
  name        = %q
  type        = "integer"
  object_types = ["tenancy.tenant"]
}

resource "netbox_custom_field" "cf_boolean" {
  name        = %q
  type        = "boolean"
  object_types = ["tenancy.tenant"]
}

resource "netbox_custom_field" "cf_date" {
  name        = %q
  type        = "date"
  object_types = ["tenancy.tenant"]
}

resource "netbox_custom_field" "cf_url" {
  name        = %q
  type        = "url"
  object_types = ["tenancy.tenant"]
}

resource "netbox_custom_field" "cf_json" {
  name        = %q
  type        = "json"
  object_types = ["tenancy.tenant"]
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
resource "netbox_tenant" "test" {
  name  = %q
  slug  = %q
  group = netbox_tenant_group.test.slug

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
		tenantGroupName, tenantGroupSlug,
		cfText,
		cfLongtext,
		cfInteger,
		cfBoolean,
		cfDate,
		cfUrl,
		cfJson,
		tag1, tag1Slug,
		tag2, tag2Slug,
		tenantName, tenantSlug,
	)
}
