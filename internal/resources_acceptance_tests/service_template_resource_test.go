package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.0", "80"),
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				Config: testAccServiceTemplateResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "2"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Updated description"),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_updated(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_service_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccServiceTemplateResourceConfig_updated(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceTemplateResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")
	updatedDescription := "Updated service template description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTagCleanup(tagSlug1)
	cleanup.RegisterTagCleanup(tagSlug2)
	cleanup.RegisterCustomFieldCleanup(cfName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_full(name, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "ports.#", "3"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Test description"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "comments", "Test comments"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_full(name, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				PlanOnly: true,
			},
			{
				Config: testAccServiceTemplateResourceConfig_fullUpdate(name, tagName1, tagSlug1, tagName2, tagSlug2, cfName, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_service_template.test", "comments", "Updated comments"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "custom_fields.0.value", "updated_value"),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_fullUpdate(name, tagName1, tagSlug1, tagName2, tagSlug2, cfName, updatedDescription),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceTemplateResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-tags")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccServiceTemplateResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccServiceTemplateResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag3-%s", tag3Slug),
						"slug": tag3Slug,
					}),
				),
			},
			{
				Config: testAccServiceTemplateResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccServiceTemplateResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-tag-order")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
			{
				Config: testAccServiceTemplateResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag1-%s", tag1Slug),
						"slug": tag1Slug,
					}),
					resource.TestCheckTypeSetElemNestedAttrs("netbox_service_template.test", "tags.*", map[string]string{
						"name": fmt.Sprintf("Tag2-%s", tag2Slug),
						"slug": tag2Slug,
					}),
				),
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}
`, name)
}

func testAccServiceTemplateResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "udp"
  ports       = [53, 123]
  description = "Updated description"
}
`, name+"-updated")
}

func testAccServiceTemplateResourceConfig_full(name, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[2]q
	slug = %[3]q
}

resource "netbox_tag" "tag2" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[6]q
	object_types = ["ipam.servicetemplate"]
	type         = "text"
}

resource "netbox_service_template" "test" {
	name        = %[1]q
  protocol    = "tcp"
  ports       = [80, 443, 8080]
  description = "Test description"
  comments    = "Test comments"

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
`, name, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccServiceTemplateResourceConfig_fullUpdate(name, tagName1, tagSlug1, tagName2, tagSlug2, cfName, description string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
	name = %[2]q
	slug = %[3]q
}

resource "netbox_tag" "tag2" {
	name = %[4]q
	slug = %[5]q
}

resource "netbox_custom_field" "test_field" {
	name         = %[6]q
	object_types = ["ipam.servicetemplate"]
	type         = "text"
}

resource "netbox_service_template" "test" {
	name        = %[1]q
	protocol    = "udp"
	ports       = [53, 123]
	description = %[7]q
	comments    = "Updated comments"

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
`, name, tagName1, tagSlug1, tagName2, tagSlug2, cfName, description)
}

func TestAccServiceTemplateResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-upd")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description1),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", testutil.Description1),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description1),
				PlanOnly: true,
			},
			{
				Config: testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", testutil.Description2),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_withDescription(name, testutil.Description2),
				PlanOnly: true,
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_withDescription(name, description string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %q
  protocol    = "tcp"
  ports       = [80]
  description = %q
}
`, name, description)
}

func TestAccServiceTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-ext")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
				),
			},
			{
				Config:   testAccServiceTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamServiceTemplatesList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find service template for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamServiceTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete service template: %v", err)
					}
					t.Logf("Successfully externally deleted service template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
func TestAccConsistency_ServiceTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_service_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
				),
			},
			{
				Config:   testAccServiceTemplateConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccServiceTemplateConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80, 443]
}
`, name)
}
func TestAccServiceTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("service-template-opt")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %[1]q
  protocol    = "udp"
  ports       = [80]
  description = "Description"
  comments    = "Comments"
  tags        = []
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "udp"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "description", "Description"),
					resource.TestCheckResourceAttr("netbox_service_template.test", "comments", "Comments"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name        = %[1]q
  protocol    = "udp"
  ports       = [80]
  description = "Description"
  comments    = "Comments"
  tags        = []
}
`, name),
				PlanOnly: true,
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %[1]q
  ports    = [80]
}
`, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_service_template.test", "protocol", "tcp"), // Default protocol is tcp
					resource.TestCheckNoResourceAttr("netbox_service_template.test", "description"),
					resource.TestCheckNoResourceAttr("netbox_service_template.test", "comments"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %[1]q
  ports    = [80]
}
`, name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccServiceTemplateResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_service_template",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_service_template" "test" {
  # name missing
  ports = [80]
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_ports": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_service_template" "test" {
  name = "test-service"
  # ports missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}

func testAccServiceTemplateResourceConfig_tags(name, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
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

resource "netbox_service_template" "test" {
  name     = %[1]q
  protocol = "tcp"
  ports    = [80]
  %[5]s
}
`, name, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccServiceTemplateResourceConfig_tagsOrder(name, tag1Slug, tag2Slug, tagCase string) string {
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

resource "netbox_service_template" "test" {
  name     = %[1]q
  protocol = "tcp"
  ports    = [80]
  %[4]s
}
`, name, tag1Slug, tag2Slug, tagsConfig)
}
