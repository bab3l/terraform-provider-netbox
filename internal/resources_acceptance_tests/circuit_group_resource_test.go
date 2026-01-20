package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitGroupResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group")
	slug := testutil.RandomSlug("tf-test-cg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-full")
	slug := testutil.RandomSlug("tf-test-cg-full")
	description := testutil.Description1

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", description),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-upd")
	slug := testutil.RandomSlug("tf-test-cg-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
				),
			},
			{
				Config: testAccCircuitGroupResourceConfig_full(name, slug, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_import(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-imp")
	slug := testutil.RandomSlug("tf-test-cg-imp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
			},
			{
				ResourceName:      "netbox_circuit_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccCircuitGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_CircuitGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-lit")
	slug := testutil.RandomSlug("tf-test-cg-lit")
	description := testutil.Description1

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckCircuitGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", description),
				),
			},
			{
				Config:   testAccCircuitGroupConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
				),
			},
		},
	})
}

func testAccCircuitGroupConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
}
`, name, slug, description)
}

func testAccCircuitGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccCircuitGroupResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
}
`, name, slug, description)
}

func TestAccCircuitGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-ext-del")
	slug := testutil.RandomSlug("circuit-group-ext-del")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List circuit groups filtered by slug
					items, _, err := client.CircuitsAPI.CircuitsCircuitGroupsList(context.Background()).SlugIc([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find circuit group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsCircuitGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete circuit group: %v", err)
					}
					t.Logf("Successfully externally deleted circuit group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccCircuitGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-circuit-group-opt")
	slug := testutil.RandomSlug("circuit-group-opt")
	tenantName := testutil.RandomName("tf-test-tenant-cg")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-cg")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckCircuitGroupDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create circuit group with tenant
			{
				Config: testAccCircuitGroupResourceConfig_withTenant(name, slug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "tenant"),
				),
			},
			// Step 2: Remove tenant (should set it to null)
			{
				Config: testAccCircuitGroupResourceConfig_withoutTenant(name, slug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_circuit_group.test", "tenant"),
				),
			},
			// Step 3: Re-add tenant (verify it can be set again)
			{
				Config: testAccCircuitGroupResourceConfig_withTenant(name, slug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "tenant"),
				),
			},
		},
	})
}

func testAccCircuitGroupResourceConfig_withTenant(name, slug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit_group" "test" {
  name   = %[1]q
  slug   = %[2]q
  tenant = netbox_tenant.test.id
}
`, name, slug, tenantName, tenantSlug)
}

func testAccCircuitGroupResourceConfig_withoutTenant(name, slug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug, tenantName, tenantSlug)
}

func TestAccCircuitGroupResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("circuit-group-tag")
	slug := testutil.RandomSlug("circuit-group-tag")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Name := testutil.RandomName("tag3")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	testutil.RunTagLifecycleTest(t, testutil.TagLifecycleTestConfig{
		ResourceName: "netbox_circuit_group",
		ConfigWithoutTags: func() string {
			return testAccCircuitGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "none")
		},
		ConfigWithTags: func() string {
			return testAccCircuitGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag1_tag2")
		},
		ConfigWithDifferentTags: func() string {
			return testAccCircuitGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, "tag2_tag3")
		},
		ExpectedTagCount:          2,
		ExpectedDifferentTagCount: 2,
		CheckDestroy:              testutil.CheckCircuitGroupDestroy,
	})
}

func TestAccCircuitGroupResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("circuit-group-tagord")
	slug := testutil.RandomSlug("circuit-group-tagord")
	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitGroupCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	testutil.RunTagOrderTest(t, testutil.TagOrderTestConfig{
		ResourceName: "netbox_circuit_group",
		ConfigWithTagsOrderA: func() string {
			return testAccCircuitGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, true)
		},
		ConfigWithTagsOrderB: func() string {
			return testAccCircuitGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, false)
		},
		ExpectedTagCount: 2,
		CheckDestroy:     testutil.CheckCircuitGroupDestroy,
	})
}

func testAccCircuitGroupResourceConfig_tagLifecycle(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug, tagSet string) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_tag" "tag2" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_tag" "tag3" {
  name = %[7]q
  slug = %[8]q
}
`, name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug, tag3Name, tag3Slug)

	//nolint:goconst // tagSet values are test-specific identifiers
	switch tagSet {
	case "tag1_tag2":
		return baseConfig + fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, name, slug)
	case "tag2_tag3":
		return baseConfig + fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
	tags = [netbox_tag.tag2.slug, netbox_tag.tag3.slug]
}
`, name, slug)
	default: // "none"
		return baseConfig + fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
  tags = []
}
`, name, slug)
	}
}

func testAccCircuitGroupResourceConfig_tagOrder(name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug string, tag1First bool) string {
	baseConfig := fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_tag" "tag2" {
  name = %[5]q
  slug = %[6]q
}
`, name, slug, tag1Name, tag1Slug, tag2Name, tag2Slug)

	if tag1First {
		return baseConfig + fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]
}
`, name, slug)
	}

	return baseConfig + fmt.Sprintf(`
resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
	tags = [netbox_tag.tag2.slug, netbox_tag.tag1.slug]
}
`, name, slug)
}

func TestAccCircuitGroupResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_circuit_group",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_circuit_group" "test" {
  # name missing
  slug = "test-group"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_circuit_group" "test" {
  name = "Test Group"
  # slug missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
