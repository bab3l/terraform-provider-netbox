package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRouteTargetResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:100")
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
		},
	})
}

func TestAccRouteTargetResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:200")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")
	updatedDescription := "Updated route target description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRouteTargetDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Test route target with full options"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "comments", "Test comments for route target"),
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
			{
				Config: testAccRouteTargetResourceConfig_fullUpdate(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_route_target.test", "comments", "Updated comments for route target"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "custom_fields.0.value", "updated_value"),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_fullUpdate(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName, updatedDescription),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRouteTargetResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:210")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRouteTargetResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRouteTargetResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccRouteTargetResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccRouteTargetResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:220")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccRouteTargetResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_route_target.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func TestAccRouteTargetResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:300")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				Config: testAccRouteTargetResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Updated description"),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_updated(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRouteTargetResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:100")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_route_target.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccRouteTargetResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccRouteTargetResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:500")
	tenantName := testutil.RandomName("tf-test-tenant-remove")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-remove")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRouteTargetDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create route target with tenant
			{
				Config: testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Test route target with full options"),
				),
			},
			// Step 2: Remove tenant and verify it's actually removed
			{
				Config: testAccRouteTargetResourceConfig_noTenant(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_route_target.test", "tenant"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Description after tenant removal"),
				),
			},
			// Step 3: Re-add tenant to verify it can be set again
			{
				Config: testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "tenant"),
				),
			},
		},
	})
}

func TestAccRouteTargetResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("65000:999")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamRouteTargetsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find route target for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamRouteTargetsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete route target: %v", err)
					}
					t.Logf("Successfully externally deleted route target with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccConsistency_RouteTarget_LiteralNames(t *testing.T) {
	t.Parallel()

	rtName := testutil.RandomName("65000:100")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", rtName),
					resource.TestCheckResourceAttr("netbox_route_target.test", "tenant", tenantName),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug),
			},
		},
	})
}

func testAccRouteTargetResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "test" {
  name = %q
}
`, name)
}

func testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[8]q
	object_types = ["ipam.routetarget"]
	type         = "text"
}

resource "netbox_tenant" "test" {
	name = %[2]q
	slug = %[3]q
}

resource "netbox_route_target" "test" {
	name        = %[1]q
  description = "Test route target with full options"
  comments    = "Test comments for route target"
  tenant      = netbox_tenant.test.id

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

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "test_value"
		}
	]
}
`, name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccRouteTargetResourceConfig_fullUpdate(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName, description string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[8]q
	object_types = ["ipam.routetarget"]
	type         = "text"
}

resource "netbox_tenant" "test" {
	name = %[2]q
	slug = %[3]q
}

resource "netbox_route_target" "test" {
	name        = %[1]q
	description = %[9]q
	comments    = "Updated comments for route target"
	tenant      = netbox_tenant.test.id

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

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "updated_value"
		}
	]
}
`, name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName, description)
}

func testAccRouteTargetResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "test" {
  name        = %q
  description = "Updated description"
}
`, name)
}

func testAccRouteTargetResourceConfig_noTenant(name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tag" "tag2" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[8]q
	object_types = ["ipam.routetarget"]
	type         = "text"
}

resource "netbox_tenant" "test" {
	name = %[2]q
	slug = %[3]q
}

resource "netbox_route_target" "test" {
	name        = %[1]q
  description = "Description after tenant removal"
  comments    = "Test comments for route target"
  # tenant intentionally omitted - should be null in state

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

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "test_value"
		}
	]
}
`, name, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = "%[2]s"
  slug = "%[3]s"
}

resource "netbox_route_target" "test" {
  name = "%[1]s"
  # Use literal string name to mimic existing user state
  tenant = "%[2]s"

  depends_on = [netbox_tenant.test]
}
`, rtName, tenantName, tenantSlug)
}

func TestAccRouteTargetResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_route_target",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_route_target" "test" {
  # name missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func testAccRouteTargetResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleNested
	case caseTag3:
		tagsConfig = tagsSingleNested
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[2]s"
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[4]s"
  slug = %[4]q
}

resource "netbox_route_target" "test" {
  name = %[1]q
  %[5]s
}
`, name, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccRouteTargetResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleNested
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleNestedReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[2]s"
  slug = %[2]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[3]s"
  slug = %[3]q
}

resource "netbox_route_target" "test" {
  name = %[1]q
  %[4]s
}
`, name, tag1Slug, tag2Slug, tagsConfig)
}
