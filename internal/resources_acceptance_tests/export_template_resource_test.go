package resources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExportTemplateResource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-export-template")

	resource.ParallelTest(t, resource.TestCase{
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
				ResourceName:      "netbox_export_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccExportTemplateResource_full(t *testing.T) {
	name := acctest.RandomWithPrefix("test-export-template")

	resource.ParallelTest(t, resource.TestCase{
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
  template_code = "{% for site in queryset %}{{ site.name }},{{ site.slug }}\n{% endfor %}"
  description   = "Updated description"
}
`
}

func testAccExportTemplateResourceConfig_full(name string) string {
	return `
resource "netbox_export_template" "test" {
  name           = "` + name + `"
  object_types   = ["dcim.site", "dcim.device"]
  template_code  = "name,slug\n{% for obj in queryset %}{{ obj.name }},{{ obj.slug }}\n{% endfor %}"
  description    = "Test description"
  mime_type      = "text/csv"
  file_extension = "csv"
  as_attachment  = true
}
`
}
