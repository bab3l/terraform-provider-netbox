package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigContextResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-config-context")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),
					resource.TestCheckResourceAttr("netbox_config_context.test", "data", "{\"foo\":\"bar\"}"),
				),
			},
			{
				ResourceName:      "netbox_config_context.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccConfigContextResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConfigContextResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-config-context-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated config context description"
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTagCleanup(tagSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_full(name, description, siteName, siteSlug, tenantName, tenantSlug, tagName, tagSlug, 500, true, "{\"ntp_servers\":[\"10.0.0.1\",\"10.0.0.2\"]}"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", description),
					resource.TestCheckResourceAttr("netbox_config_context.test", "weight", "500"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "is_active", "true"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "data", "{\"ntp_servers\":[\"10.0.0.1\",\"10.0.0.2\"]}"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "sites.#", "1"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "tenants.#", "1"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "tags.#", "1"),
				),
			},
			{
				Config: testAccConfigContextResourceConfig_full(name, updatedDescription, siteName, siteSlug, tenantName, tenantSlug, tagName, tagSlug, 2000, false, "{\"dns_servers\":[\"8.8.8.8\",\"8.8.4.4\"]}"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_config_context.test", "weight", "2000"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "is_active", "false"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "data", "{\"dns_servers\":[\"8.8.8.8\",\"8.8.4.4\"]}"),
				),
			},
			{
				Config:   testAccConfigContextResourceConfig_full(name, updatedDescription, siteName, siteSlug, tenantName, tenantSlug, tagName, tagSlug, 2000, false, "{\"dns_servers\":[\"8.8.8.8\",\"8.8.4.4\"]}"),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_ConfigContext_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-config-context-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),
					resource.TestCheckResourceAttr("netbox_config_context.test", "data", "{\"foo\":\"bar\"}"),
				),
			},
			{
				Config:   testAccConfigContextResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
				),
			},
		},
	})
}

func TestAccConfigContextResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("config-context")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),
				),
			},
			{
				Config:   testAccConfigContextResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConfigContextResource_update(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-ctx-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_withDescription(name, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccConfigContextResourceConfig_withDescription(name, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", testutil.Description2),
				),
			},
			{
				Config:   testAccConfigContextResourceConfig_withDescription(name, testutil.Description2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConfigContextResource_externalDeletion(t *testing.T) {
	t.Parallel()

	testutil.TestAccPreCheck(t)
	name := testutil.RandomName("tf-test-ctx-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find config context by name
					items, _, err := client.ExtrasAPI.ExtrasConfigContextsList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list config contexts: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Config context not found with name: %s", name)
					}

					// Delete the config context
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasConfigContextsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config context: %v", err)
					}

					t.Logf("Successfully externally deleted config context with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccConfigContextResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccConfigContextResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-ctx-rem")
	description := testutil.RandomName("description")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	tagSlug := testutil.RandomSlug("tf-test-tag")
	regionName := testutil.RandomName("tf-test-region")
	regionSlug := testutil.RandomSlug("tf-test-region")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug)
	cleanup.RegisterRegionCleanup(regionSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_config_context",
		BaseConfig: func() string {
			return testAccConfigContextResourceConfig_removeOptionalFields_base(name)
		},
		ConfigWithFields: func() string {
			return testAccConfigContextResourceConfig_removeOptionalFields_withFields(
				name, description, siteName, siteSlug, tagSlug,
				regionName, regionSlug, tenantName, tenantSlug,
			)
		},
		OptionalFields: map[string]string{
			"description": description,
			// Note: weight and is_active have computed defaults and cannot be truly cleared
		},
		RequiredFields: map[string]string{
			"name": name,
			// Note: data field has JSON normalization, so we don't check its exact value
		},
		CheckDestroy: nil, // No CheckConfigContextDestroy function available
	})
}

func testAccConfigContextResourceConfig_removeOptionalFields_base(name string) string {
	return fmt.Sprintf(`
resource "netbox_config_context" "test" {
  name = %q
  data = "{\"key\":\"value\"}"
}
`, name)
}

func testAccConfigContextResourceConfig_removeOptionalFields_withFields(
	name, description, siteName, siteSlug, tagSlug,
	regionName, regionSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_region" "test" {
  name = %[6]q
  slug = %[7]q
}

resource "netbox_site" "test" {
  name   = %[3]q
  slug   = %[4]q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %[8]q
  slug = %[9]q
}

resource "netbox_tag" "test" {
  name = %[5]q
  slug = %[5]q
}

resource "netbox_config_context" "test" {
  name        = %[1]q
  data        = "{\"key\":\"value\"}"
  description = %[2]q
  regions     = [netbox_region.test.id]
  sites       = [netbox_site.test.id]
  tenants     = [netbox_tenant.test.id]
  tags        = [netbox_tag.test.slug]
}
`, name, description, siteName, siteSlug, tagSlug, regionName, regionSlug, tenantName, tenantSlug)
}

func testAccConfigContextResourceConfig_basic(name string) string {

	return fmt.Sprintf(`
resource "netbox_config_context" "test" {
  name = %q
  data = "{\"foo\":\"bar\"}"
}
`, name)
}

func testAccConfigContextResourceConfig_withDescription(name string, description string) string {
	return fmt.Sprintf(`
resource "netbox_config_context" "test" {
  name        = %[1]q
  data        = "{\"key\": \"value\"}"
  description = %[2]q
}
`, name, description)
}

func testAccConfigContextResourceConfig_full(name, description, siteName, siteSlug, tenantName, tenantSlug, tagName, tagSlug string, weight int, isActive bool, data string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_config_context" "test" {
  name        = %q
  description = %q
  weight      = %d
  is_active   = %t
  data        = %q
  sites       = [netbox_site.test.id]
  tenants     = [netbox_tenant.test.id]
  tags        = [netbox_tag.test.slug]
}
`, siteName, siteSlug, tenantName, tenantSlug, tagName, tagSlug, name, description, weight, isActive, data)
}
