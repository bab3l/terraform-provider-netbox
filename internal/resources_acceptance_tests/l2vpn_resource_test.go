package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNResource_basic(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)
	cleanup.RegisterL2VPNCleanup(name + "-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Updated description"),
				),
			},
			{
				ResourceName:            "netbox_l2vpn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name"}, // display_name is computed and may differ after name changes
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_l2vpn.test", "tenant"),
				),
			},
			{
				Config:             testAccL2VPNResourceConfig_updated(name),
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccL2VPNResource_full(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vpls"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "identifier", "12345"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Test L2VPN"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "comments", "Test comments"),
				),
			},
		},
	})
}

func TestAccL2VPNResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccL2VPNResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_l2vpn.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

// NOTE: Custom field tests for l2vpn resource are in resources_acceptance_tests_customfields package

func testAccL2VPNResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}
`, name, name)
}

func testAccL2VPNResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = "Updated description"
}
`, name+"-updated", name)
}

func testAccL2VPNResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vpls"
  identifier  = 12345
  description = "Test L2VPN"
  comments    = "Test comments"
}
`, name, name)
}

func testAccL2VPNResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
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

resource "netbox_l2vpn" "test" {
  name = %[1]q
  slug = %[1]q
  type = "vxlan"
  %[5]s
}
`, name, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccL2VPNResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
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

resource "netbox_l2vpn" "test" {
  name = %[1]q
  slug = %[1]q
  type = "vxlan"
  %[4]s
}
`, name, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccConsistency_L2VPN_LiteralNames(t *testing.T) {
	t.Parallel()

	name := "test-l2vpn-lit"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				Config:   testAccL2VPNResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
		},
	})
}

func TestAccL2VPNResource_update(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_updateInitial(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", testutil.Description1),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_updateModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", testutil.Description2),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
		},
	})
}

func testAccL2VPNResourceConfig_updateInitial(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = %q
}
`, name, name, testutil.Description1)
}

func testAccL2VPNResourceConfig_updateModified(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = %q
}
`, name, name, testutil.Description2)
}

func TestAccL2VPNResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnL2vpnsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find l2vpn for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnL2vpnsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete l2vpn: %v", err)
					}
					t.Logf("Successfully externally deleted l2vpn with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccL2VPNResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn-opt")
	tenantName := testutil.RandomName("tf-test-tenant-l2vpn")
	tenantSlug := testutil.RandomSlug("tf-test-tenant-l2vpn")
	importTargetName := fmt.Sprintf("65000:%d", acctest.RandIntRange(1000, 9999))
	exportTargetName := fmt.Sprintf("65000:%d", acctest.RandIntRange(1000, 9999))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterL2VPNCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRouteTargetCleanup(importTargetName)
	cleanup.RegisterRouteTargetCleanup(exportTargetName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckL2VPNDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			// Step 1: Create L2VPN with tenant and identifier
			{
				Config: testAccL2VPNResourceConfig_withTenantAndTargets(name, tenantName, tenantSlug, importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "identifier"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "import_targets.#", "1"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "export_targets.#", "1"),
				),
			},
			// Step 2: Remove optional fields (should clear them)
			{
				Config: testAccL2VPNResourceConfig_requiredOnly(name, tenantName, tenantSlug, importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_l2vpn.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_l2vpn.test", "identifier"),
					resource.TestCheckNoResourceAttr("netbox_l2vpn.test", "import_targets"),
					resource.TestCheckNoResourceAttr("netbox_l2vpn.test", "export_targets"),
				),
			},
			// Step 3: Re-add tenant and identifier (verify they can be set again)
			{
				Config: testAccL2VPNResourceConfig_withTenantAndTargets(name, tenantName, tenantSlug, importTargetName, exportTargetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "identifier"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "import_targets.#", "1"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "export_targets.#", "1"),
				),
			},
		},
	})
}

func testAccL2VPNResourceConfig_withTenantAndTargets(name, tenantName, tenantSlug, importTargetName, exportTargetName string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_route_target" "import" {
  name = %[4]q
}

resource "netbox_route_target" "export" {
  name = %[5]q
}

resource "netbox_l2vpn" "test" {
  name       = %[1]q
  slug       = %[1]q
  type       = "vxlan"
  tenant     = netbox_tenant.test.id
  identifier = 12345
  import_targets = [netbox_route_target.import.id]
  export_targets = [netbox_route_target.export.id]
}
`, name, tenantName, tenantSlug, importTargetName, exportTargetName)
}

func testAccL2VPNResourceConfig_requiredOnly(name, tenantName, tenantSlug, importTargetName, exportTargetName string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_route_target" "import" {
  name = %[4]q
}

resource "netbox_route_target" "export" {
  name = %[5]q
}

resource "netbox_l2vpn" "test" {
  name = %[1]q
  slug = %[1]q
  type = "vxlan"
}
`, name, tenantName, tenantSlug, importTargetName, exportTargetName)
}

func TestAccL2VPNResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_l2vpn",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_l2vpn" "test" {
  # name missing
  slug = "test-l2vpn"
  type = "vxlan"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_l2vpn" "test" {
  name = "Test L2VPN"
  # slug missing
  type = "vxlan"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_type": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_l2vpn" "test" {
  name = "Test L2VPN"
  slug = "test-l2vpn"
  # type missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
