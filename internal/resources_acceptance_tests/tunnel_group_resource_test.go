package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTunnelGroupResource_basic(t *testing.T) {
	t.Parallel()

	// Generate unique names to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-tunnel-group")
	slug := testutil.RandomSlug("tf-test-tunnel-grp")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "slug", slug),
				),
			},
			{
				Config:   testAccTunnelGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelGroupResource_full(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-group-full")
	slug := testutil.RandomSlug("tf-test-tg-full")
	description := "Test tunnel group with all fields"
	updatedDescription := "Updated tunnel group description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", description),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config:   testAccTunnelGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
			{
				Config: testAccTunnelGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "custom_fields.0.value", "updated_value"),
				),
			},
			{
				Config:   testAccTunnelGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelGroupResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-group-tags")
	slug := testutil.RandomSlug("tf-test-tg-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTunnelGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTunnelGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccTunnelGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccTunnelGroupResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-group-tag-order")
	slug := testutil.RandomSlug("tf-test-tg-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccTunnelGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_tunnel_group.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccTunnelGroupResource_update(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-group-upd")
	slug := testutil.RandomSlug("tf-test-tg-upd")
	updatedDescription := testutil.Description2
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
				),
			},
			{
				Config:   testAccTunnelGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
			{
				Config: testAccTunnelGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", updatedDescription),
				),
			},
			{
				Config:   testAccTunnelGroupResourceConfig_fullUpdate(name, slug, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelGroupResource_import(t *testing.T) {
	t.Parallel()

	// Generate unique names
	name := testutil.RandomName("tf-test-tunnel-group-imp")
	slug := testutil.RandomSlug("tf-test-tg-imp")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
			},
			{
				ResourceName:      "netbox_tunnel_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccTunnelGroupResourceConfig_basic(name, slug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccTunnelGroupResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tunnel-group-extdel")
	slug := testutil.RandomSlug("tf-test-tg-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnTunnelGroupsList(context.Background()).Slug([]string{slug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find tunnel group for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnTunnelGroupsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete tunnel group: %v", err)
					}
					t.Logf("Successfully externally deleted tunnel group with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccConsistency_TunnelGroup_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tg")
	slug := testutil.RandomSlug("tg")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupConsistencyLiteralNamesConfig(name, slug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "custom_fields.#", "1"),
				),
			},
			{
				PlanOnly: true,
				Config:   testAccTunnelGroupConsistencyLiteralNamesConfig(name, slug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
			},
		},
	})
}

func testAccTunnelGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_tunnel_group" "test" {
  name = %[1]q
  slug = %[2]q
}
`, name, slug)
}

func testAccTunnelGroupResourceConfig_full(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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
	object_types = ["vpn.tunnelgroup"]
	type         = "text"
}

resource "netbox_tunnel_group" "test" {
	name        = %[1]q
	slug        = %[2]q
	description = %[3]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "test_value"
		}
	]
}
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccTunnelGroupResourceConfig_fullUpdate(name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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
	object_types = ["vpn.tunnelgroup"]
	type         = "text"
}

resource "netbox_tunnel_group" "test" {
	name        = %[1]q
	slug        = %[2]q
	description = %[3]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "updated_value"
		}
	]
}
`, name, slug, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccTunnelGroupResourceConfig_withoutDescription(name, slug, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_tag" "tag2" {
	name = %[5]q
	slug = %[6]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[7]q
	object_types = ["vpn.tunnelgroup"]
	type         = "text"
}

resource "netbox_tunnel_group" "test" {
	name = %[1]q
	slug = %[2]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "test_value"
		}
	]
}
`, name, slug, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccTunnelGroupConsistencyLiteralNamesConfig(name, slug, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[3]q
	slug = %[4]q
}

resource "netbox_tag" "tag2" {
	name = %[5]q
	slug = %[6]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[7]q
	object_types = ["vpn.tunnelgroup"]
	type         = "text"
}

resource "netbox_tunnel_group" "test" {
	name = %[1]q
	slug = %[2]q

	tags = [
		netbox_tag.tag1.slug,
		netbox_tag.tag2.slug
	]

	custom_fields = [
		{
			name  = netbox_custom_field.test_field.name
			type  = "text"
			value = "test_value"
		}
	]
}
`, name, slug, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func TestAccTunnelGroupResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-tg-rem")
	slug := testutil.RandomSlug("tf-test-tg-rem")
	const testDescription = "Test Description"
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTunnelGroupCleanup(name)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckTunnelGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelGroupResourceConfig_full(name, slug, testDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", testDescription),
				),
			},
			{
				Config: testAccTunnelGroupResourceConfig_withoutDescription(name, slug, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_tunnel_group.test", "description"),
				),
			},
		},
	})
}

func TestAccTunnelGroupResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_tunnel_group",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_tunnel_group" "test" {
  # name missing
  slug = "test-tunnel-group"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_slug": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_tunnel_group" "test" {
  name = "test-tunnel-group"
  # slug missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func testAccTunnelGroupResourceConfig_tags(name, slug, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[5]s"
  slug = %[5]q
}

resource "netbox_tunnel_group" "test" {
  name = %[1]q
  slug = %[2]q
  %[6]s
}
`, name, slug, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccTunnelGroupResourceConfig_tagsOrder(name, slug, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_tunnel_group" "test" {
  name = %[1]q
  slug = %[2]q
  %[5]s
}
`, name, slug, tag1Slug, tag2Slug, tagsConfig)
}
