package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExportTemplateResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("export-template")

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
		},
	})
}

func TestAccExportTemplateResource_full(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("export-template")

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

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateConsistencyLiteralNamesConfig(name, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_export_template.test", "description", description),
				),
			},
			{
				Config:   testAccExportTemplateConsistencyLiteralNamesConfig(name, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_export_template.test", "id"),
				),
			},
		},
	})
}

func testAccExportTemplateConsistencyLiteralNamesConfig(name, description string) string {
	return `
resource "netbox_export_template" "test" {
  name           = "` + name + `"
  object_types   = ["dcim.site"]
  template_code  = "name,slug\n{% for site in queryset %}{{ site.name }},{{ site.id }}\n{% endfor %}"
  description    = "` + description + `"
}
`
}
