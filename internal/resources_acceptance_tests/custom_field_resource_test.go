package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// NOTE: Some tests in this file create custom fields with required=true.
// These tests must NOT run in parallel (t.Parallel removed) because required
// custom fields can interfere with other acceptance tests that depend on
// predictable resource state.

func TestAccCustomFieldResource_basic(t *testing.T) {
	t.Parallel()

	// Custom field names can only contain alphanumeric characters and underscores
	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "text"),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomFieldResourceConfig_basic(name),
				PlanOnly: true,
			},
			{
				ResourceName:      "netbox_custom_field.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCustomFieldResource_full(t *testing.T) {
	t.Parallel()

	// Custom field names can only contain alphanumeric characters and underscores
	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))
	description := testutil.RandomName("description")
	updatedDescription := "Updated custom field description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldResourceConfig_full(name, description, "Custom Label", "Group A", 2000, "exact", "if-set", "no", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "integer"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", description),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "label", "Custom Label"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "group_name", "Group A"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "required", "false"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "search_weight", "2000"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "filter_logic", "exact"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "ui_visible", "if-set"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "ui_editable", "no"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "is_cloneable", "true"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "validation_minimum", "1"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "validation_maximum", "100"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "weight", "50"),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomFieldResourceConfig_full(name, description, "Custom Label", "Group A", 2000, "exact", "if-set", "no", true),
				PlanOnly: true,
			},
			{
				Config: testAccCustomFieldResourceConfig_full(name, updatedDescription, "Updated Label", "Group B", 3000, "loose", "always", "yes", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "label", "Updated Label"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "group_name", "Group B"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "search_weight", "3000"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "filter_logic", "loose"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "ui_visible", "always"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "ui_editable", "yes"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "is_cloneable", "false"),
				),
			},
			{
				// Verify no changes after update
				Config:   testAccCustomFieldResourceConfig_full(name, updatedDescription, "Updated Label", "Group B", 3000, "loose", "always", "yes", false),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_CustomField_LiteralNames(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "text"),
				),
			},
			{
				Config:   testAccCustomFieldResourceConfig_basic(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),
				),
			},
		},
	})
}

func TestAccCustomFieldResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomFieldResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCustomFieldResource_DescriptionUpdate(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldResourceConfig_withDescription(name, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", testutil.Description1),
				),
			},
			{
				// Verify no changes after create
				Config:   testAccCustomFieldResourceConfig_withDescription(name, testutil.Description1),
				PlanOnly: true,
			},
			{
				Config: testAccCustomFieldResourceConfig_withDescription(name, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", testutil.Description2),
				),
			},
			{
				// Verify no changes after update
				Config:   testAccCustomFieldResourceConfig_withDescription(name, testutil.Description2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCustomFieldResource_externalDeletion(t *testing.T) {
	t.Parallel()

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find custom field by name
					items, _, err := client.ExtrasAPI.ExtrasCustomFieldsList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list custom fields: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Custom field not found with name: %s", name)
					}

					// Delete the custom field
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasCustomFieldsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete custom field: %v", err)
					}

					t.Logf("Successfully externally deleted custom field with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCustomFieldResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  type         = "text"
  object_types = ["dcim.site"]
}
`, name)
}

func TestAccCustomFieldResource_digitStartingName(t *testing.T) {
	// This test cannot use t.Parallel() because it creates a custom field that will be used by other tests.
	// Testing custom field names that start with digits validates the regex fix that allows this pattern.

	name := fmt.Sprintf("%s_%s", "4me", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldResourceConfig_digitStartingName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "text"),
				),
			},
		},
	})
}

func testAccCustomFieldResourceConfig_full(name, description, label, groupName string, searchWeight int, filterLogic, uiVisible, uiEditable string, isCloneable bool) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name               = %q
  type               = "integer"
  object_types       = ["dcim.site", "dcim.device"]
  description        = %q
  label              = %q
  group_name         = %q
  required           = false
  search_weight      = %d
  filter_logic       = %q
  ui_visible         = %q
  ui_editable        = %q
  is_cloneable       = %t
  validation_minimum = 1
  validation_maximum = 100
  weight             = 50
}
`, name, description, label, groupName, searchWeight, filterLogic, uiVisible, uiEditable, isCloneable)
}

func testAccCustomFieldResourceConfig_digitStartingName(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  type         = "text"
  object_types = ["dcim.site"]
  required     = false
}
`, name)
}

func testAccCustomFieldResourceConfig_withDescription(name string, description string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %[1]q
  type         = "text"
  object_types = ["dcim.site"]
  description  = %[2]q
}
`, name, description)
}
