package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Custom field tests for virtual chassis resource are in resources_acceptance_tests_customfields package

func TestAccVirtualChassisResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vc")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
				),
			},
			{
				ResourceName:      "netbox_virtual_chassis.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccVirtualChassisResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vc-update")
	updatedName := testutil.RandomName("tf-test-vc-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisResourceConfig_forUpdate(name, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", "domain1.example.com"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccVirtualChassisResourceConfig_forUpdate(updatedName, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", "domain2.example.com"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccVirtualChassisResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vc-full")
	description := testutil.RandomName("description")
	updatedDescription := testutil.RandomName("description")
	tagName1 := testutil.RandomName("tag1")
	tagSlug1 := testutil.RandomSlug("tag1")
	tagName2 := testutil.RandomName("tag2")
	tagSlug2 := testutil.RandomSlug("tag2")
	cfName := testutil.RandomCustomFieldName("test_field")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisResourceConfig_full(name, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", "test-domain"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", description),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "comments", "Test comments"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "custom_fields.0.value", "test_value"),
				),
			},
			{
				Config: testAccVirtualChassisResourceConfig_fullUpdate(name, updatedDescription, tagName1, tagSlug1, tagName2, tagSlug2, cfName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "custom_fields.0.value", "updated_value"),
				),
			},
		},
	})
}

func TestAccVirtualChassisResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vc-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
				),
			},
		},
	})

}

func TestAccConsistency_VirtualChassis_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vc-lit")
	domain := "test-domain-lit.example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisConsistencyLiteralNamesConfig(name, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", domain),
				),
			},
			{
				Config:   testAccVirtualChassisConsistencyLiteralNamesConfig(name, domain),
				PlanOnly: true,
			},
		},
	})
}

func testAccVirtualChassisConsistencyLiteralNamesConfig(name, domain string) string {
	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
  name   = %q
  domain = %q
}
`, name, domain)
}

func testAccVirtualChassisResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
  name = %q
}
`, name)
}

func testAccVirtualChassisResourceConfig_forUpdate(name, description string) string {
	domain := "domain1.example.com"
	if description == testutil.Description2 {
		domain = "domain2.example.com"
	}

	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
  name        = %q
  domain      = %q
  description = %q
}
`, name, domain, description)
}

func testAccVirtualChassisResourceConfig_full(name, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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
  object_types = ["dcim.virtualchassis"]
  type         = "text"
}

resource "netbox_virtual_chassis" "test" {
  name        = %[1]q
  domain      = "test-domain"
  description = %[2]q
  comments    = "Test comments"

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "test_value"
    }
  ]
}
`, name, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)
}

func testAccVirtualChassisResourceConfig_fullUpdate(name, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName string) string {
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
  object_types = ["dcim.virtualchassis"]
  type         = "text"
}

resource "netbox_virtual_chassis" "test" {
  name        = %[1]q
  domain      = "test-domain"
  description = %[2]q
  comments    = "Test comments"

	tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]

  custom_fields = [
    {
      name  = netbox_custom_field.test_field.name
      type  = "text"
      value = "updated_value"
    }
  ]
}

`, name, description, tagName1, tagSlug1, tagName2, tagSlug2, cfName)

}

func TestAccVirtualChassisResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-vc-del")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.DcimAPI.DcimVirtualChassisList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find virtual_chassis for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.DcimAPI.DcimVirtualChassisDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete virtual_chassis: %v", err)
					}
					t.Logf("Successfully externally deleted virtual_chassis with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccVirtualChassisResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccVirtualChassisResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vc-rem")

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_virtual_chassis",
		BaseConfig: func() string {
			return testAccVirtualChassisResourceConfig_removeOptionalFields_base(name)
		},
		ConfigWithFields: func() string {
			return testAccVirtualChassisResourceConfig_removeOptionalFields_withFields(name)
		},
		OptionalFields: map[string]string{
			"domain":      "test-domain.example.com",
			"description": "Test Description",
			"comments":    "Test Comments",
			// Note: master requires a device to be created and is tested separately
			// Note: tags and custom_fields are generic metadata fields tested elsewhere
		},
		RequiredFields: map[string]string{
			"name": name,
		},
	})
}

func testAccVirtualChassisResourceConfig_removeOptionalFields_base(name string) string {
	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
  name = %q
}
`, name)
}

func testAccVirtualChassisResourceConfig_removeOptionalFields_withFields(name string) string {
	return fmt.Sprintf(`
resource "netbox_virtual_chassis" "test" {
  name        = %[1]q
  domain      = "test-domain.example.com"
  description = "Test Description"
  comments    = "Test Comments"
}
`, name)
}

func TestAccVirtualChassisResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_virtual_chassis",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_virtual_chassis" "test" {
  domain = "test.example.com"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
