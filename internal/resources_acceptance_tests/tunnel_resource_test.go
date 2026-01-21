package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tunnel")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "gre"),
				),
			},
			{
				Config:   testAccTunnelResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-full")
	description := testutil.RandomName("description")
	updatedDescription := testutil.RandomName("description-updated")
	comments := testutil.RandomName("comments")
	updatedComments := testutil.RandomName("comments-updated")
	groupName := testutil.RandomName("tf-test-tunnel-group")
	groupSlug := testutil.RandomSlug("tf-test-tunnel-group")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTunnelGroupCleanup(groupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_full(name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "wireguard"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", description),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "comments", comments),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tunnel_id", "12345"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "2"),
				),
			},
			{
				Config:   testAccTunnelResourceConfig_full(name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
			{
				Config: testAccTunnelResourceConfig_fullUpdate(name, updatedDescription, updatedComments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "comments", updatedComments),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tunnel_id", "54321"),
				),
			},
			{
				Config:   testAccTunnelResourceConfig_fullUpdate(name, updatedDescription, updatedComments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTunnelResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTunnelResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccTunnelResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccTunnelResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTunnelResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccTunnelResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-upd")
	updatedDescription := testutil.RandomName("description-updated")
	updatedComments := testutil.RandomName("comments-updated")
	groupName := testutil.RandomName("tf-test-tunnel-group")
	groupSlug := testutil.RandomSlug("tf-test-tunnel-group")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTunnelGroupCleanup(groupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "active"),
				),
			},
			{
				Config:   testAccTunnelResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				Config: testAccTunnelResourceConfig_fullUpdate(name, updatedDescription, updatedComments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "planned"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "comments", updatedComments),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tunnel_id", "54321"),
				),
			},
			{
				Config:   testAccTunnelResourceConfig_fullUpdate(name, updatedDescription, updatedComments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-imp")
	description := testutil.RandomName("tf-test-tunnel-desc")
	comments := testutil.RandomName("tf-test-tunnel-comments")
	groupName := testutil.RandomName("tf-test-tunnel-group")
	groupSlug := testutil.RandomSlug("tf-test-tunnel-group")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTunnelGroupCleanup(groupName)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_full(name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
			},
			{
				ResourceName:            "netbox_tunnel.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tags"},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_tunnel.test", "group"),
					testutil.ReferenceFieldCheck("netbox_tunnel.test", "tenant"),
					testutil.ReferenceFieldCheck("netbox_tunnel.test", "ipsec_profile"),
				),
			},
			{
				Config:   testAccTunnelResourceConfig_full(name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnTunnelsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tunnel for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnTunnelsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tunnel: %v", err)
					}
					t.Logf("Successfully externally deleted tunnel with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccConsistency_Tunnel_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tunnel")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelConsistencyLiteralNamesConfig(name, tagName1, tagSlug1, tagName2, tagSlug2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel.test", "tags.#", "2"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTunnelConsistencyLiteralNamesConfig(name, tagName1, tagSlug1, tagName2, tagSlug2),
			},
		},
	})
}

func testAccTunnelResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
  name          = %[1]q
  status        = "active"
  encapsulation = "gre"
}
`, name)
}

func testAccTunnelResourceConfig_full(name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tenant" "test" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_tag" "tag1" {
	name = %[8]q
	slug = %[9]q
}

resource "netbox_tag" "tag2" {
	name = %[10]q
	slug = %[11]q
}

resource "netbox_tunnel" "test" {
	name          = %[1]q
	status        = "planned"
	encapsulation = "wireguard"
	description   = %[2]q
	comments      = %[3]q
	tunnel_id     = 12345
	group         = netbox_tunnel_group.test.id
	tenant        = netbox_tenant.test.id

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2)
}

func testAccTunnelResourceConfig_fullUpdate(name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tenant" "test" {
	name = %[6]q
	slug = %[7]q
}

resource "netbox_tag" "tag1" {
	name = %[8]q
	slug = %[9]q
}

resource "netbox_tag" "tag2" {
	name = %[10]q
	slug = %[11]q
}

resource "netbox_tunnel" "test" {
	name          = %[1]q
	status        = "planned"
	encapsulation = "wireguard"
	description   = %[2]q
	comments      = %[3]q
	tunnel_id     = 54321
	group         = netbox_tunnel_group.test.id
	tenant        = netbox_tenant.test.id

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, description, comments, groupName, groupSlug, tenantName, tenantSlug, tagName1, tagSlug1, tagName2, tagSlug2)
}

func testAccTunnelConsistencyLiteralNamesConfig(name, tagName1, tagSlug1, tagName2, tagSlug2 string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[2]q
	slug = %[3]q
}

resource "netbox_tag" "tag2" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_tunnel" "test" {
	name          = %[1]q
	status        = "active"
	encapsulation = "gre"

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]
}
`, name, tagName1, tagSlug1, tagName2, tagSlug2)
}

// TestAccTunnelResource_StatusComprehensive tests comprehensive scenarios for tunnel status field.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccTunnelResource_StatusComprehensive(t *testing.T) {
	t.Parallel()

	// Generate unique names for this test run
	tunnelName := testutil.RandomName("tf-test-tunnel-status")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(tunnelName)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_tunnel",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "planned",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckTunnelDestroy,
			testutil.CheckTunnelGroupDestroy,
		),
		BaseConfig: func() string {
			return testAccTunnelResourceConfig_statusBase(tunnelName)
		},
		WithFieldConfig: func(value string) string {
			return testAccTunnelResourceConfig_statusWithField(tunnelName, value)
		},
	})
}

func testAccTunnelResourceConfig_statusBase(name string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
	name          = %[1]q
	encapsulation = "gre"
	# status field intentionally omitted - should get default "active"
}
`, name)
}

func testAccTunnelResourceConfig_statusWithField(name, status string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel" "test" {
	name          = %[1]q
	encapsulation = "gre"
	status        = %[2]q
}
`, name, status)
}

func TestAccTunnelResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-rem")
	groupName := testutil.RandomName("tf-test-tunnel-group")
	groupSlug := testutil.RandomSlug("tf-test-tunnel-group")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelCleanup(name)
	cleanup.RegisterTunnelGroupCleanup(groupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_tunnel",
		BaseConfig: func() string {
			return testAccTunnelResourceConfig_removeOptionalFields_base(
				name, groupName, groupSlug, tenantName, tenantSlug,
			)
		},
		ConfigWithFields: func() string {
			return testAccTunnelResourceConfig_removeOptionalFields_withFields(
				name, groupName, groupSlug, tenantName, tenantSlug,
			)
		},
		OptionalFields: map[string]string{
			"description": "Test Description",
			"comments":    "Test Comments",
			"tunnel_id":   "100",
			// Note: status has a default value and cannot be truly cleared
			// Note: ipsec_profile requires ipsec encapsulation type
		},
		RequiredFields: map[string]string{
			"name": name,
		},
		CheckDestroy: testutil.CheckTunnelDestroy,
	})
}

func testAccTunnelResourceConfig_removeOptionalFields_base(name, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"
}
`, name, groupName, groupSlug, tenantName, tenantSlug)
}

func testAccTunnelResourceConfig_removeOptionalFields_withFields(
	name, groupName, groupSlug, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name = %[2]q
  slug = %[3]q
}

resource "netbox_tenant" "test" {
  name = %[4]q
  slug = %[5]q
}

resource "netbox_tunnel" "test" {
  name          = %[1]q
  encapsulation = "gre"
  description   = "Test Description"
  comments      = "Test Comments"
  tunnel_id     = 100
  group         = netbox_tunnel_group.test.id
  tenant        = netbox_tenant.test.id
}
`, name, groupName, groupSlug, tenantName, tenantSlug)
}

func TestAccTunnelResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_tunnel",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_encapsulation": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_tunnel" "test" {
  name = "test-tunnel"
  # encapsulation missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func testAccTunnelResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

resource "netbox_tunnel" "test" {
  name          = %[1]q
  status        = "active"
  encapsulation = "gre"
  %[5]s
}
`, name, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccTunnelResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, tagCase string) string {
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

resource "netbox_tunnel" "test" {
  name          = %[1]q
  status        = "active"
  encapsulation = "gre"
  %[4]s
}
`, name, tag1Slug, tag2Slug, tagsConfig)
}
