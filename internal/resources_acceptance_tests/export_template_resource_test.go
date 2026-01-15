package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExportTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("export-template")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)
	cleanup.RegisterExportTemplateCleanup(name + "-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "object_types.#", "1"),
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
				),
			},
			// Test update
			{
				Config: testAccExportTemplateResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", "Updated description"),
				),
			},
			// PlanOnly: verify plan stability
			{
				Config:   testAccExportTemplateResourceConfig_updated(name),
				PlanOnly: true,
			},
			// Test import
			{
				ResourceName:            "netbox_export_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name"},
			},
		},
	})
}

func TestAccExportTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("exp-tmpl-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "object_types.#", "1"),
				),
			},
			// PlanOnly: verify plan stability
			{
				Config:   testAccExportTemplateResourceConfig_basic(name),
				PlanOnly: true,
			},
		},
	})
}

func TestAccExportTemplateResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("export-template")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "object_types.#", "2"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", "Test description"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "mime_type", "text/csv"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "file_extension", "csv"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "as_attachment", "true"),
				),
			},
			// PlanOnly: verify plan stability
			{
				Config:   testAccExportTemplateResourceConfig_full(name),
				PlanOnly: true,
			},
		},
	})
}

func testAccExportTemplateResourceConfig_basic(name string) string {
	return `
resource "netbox_export_template" "test" {
  name          = "` + name + `"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }}\n{% endfor %}"
}
`
}

func testAccExportTemplateResourceConfig_updated(name string) string {
	return `
resource "netbox_export_template" "test" {
  name          = "` + name + `-updated"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }},{{ site.id }}\n{% endfor %}"
  description   = "Updated description"
}
`
}

func testAccExportTemplateResourceConfig_full(name string) string {
	return `
resource "netbox_export_template" "test" {
  name           = "` + name + `"
  object_types   = ["dcim.site", "dcim.device"]
  template_code  = "name,slug\n{% for obj in queryset %}{{ obj.name }},{{ obj.id }}\n{% endfor %}"
  description    = "Test description"
  mime_type      = "text/csv"
  file_extension = "csv"
  as_attachment  = true
}
`
}

func TestAccConsistency_ExportTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-export-template-lit")
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_withDescription(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", description),
				),
			},
			{
				Config:   testAccExportTemplateResourceConfig_withDescription(name, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
				),
			},
		},
	})
}

func TestAccExportTemplateResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-expt-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_withDescription(name, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccExportTemplateResourceConfig_withDescription(name, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", testutil.Description2),
				),
			},
			// PlanOnly: verify plan stability
			{
				Config:   testAccExportTemplateResourceConfig_withDescription(name, testutil.Description2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccExportTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-expt-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find export template by name
					items, _, err := client.ExtrasAPI.ExtrasExportTemplatesList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list export templates: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Export template not found with name: %s", name)
					}

					// Delete the export template
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasExportTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete export template: %v", err)
					}

					t.Logf("Successfully externally deleted export template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccExportTemplateResource_removeOptionalFields tests that optional fields
// can be successfully removed from the configuration without causing inconsistent state.
func TestAccExportTemplateResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	const testDescription = "Test Description"

	name := testutil.RandomName("export-tmpl-remove")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateResourceConfig_withDescription(name, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", testDescription),
				),
			},
			{
				Config: testAccExportTemplateResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckNoResourceAttr("netbox_export_template.test", "description"),
				),
			},
		},
	})
}

func testAccExportTemplateResourceConfig_withDescription(name string, description string) string {
	return fmt.Sprintf(`
resource "netbox_export_template" "test" {
  name           = %q
  object_types   = ["dcim.site"]
  template_code  = "name,slug\n{%% for site in queryset %%}{{ site.name }},{{ site.id }}\n{%% endfor %%}"
  description    = %q
}
`, name, description)
}

func TestAccExportTemplateResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_export_template",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_name": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_export_template" "test" {
  # name missing
  object_types = ["dcim.device"]
  template_code = "{{ obj.name }}"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_object_types": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_export_template" "test" {
  name = "Test Template"
  # object_types missing
  template_code = "{{ obj.name }}"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_template_code": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_export_template" "test" {
  name = "Test Template"
  object_types = ["dcim.device"]
  # template_code missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
