package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExportTemplateDataSource_byID(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("test-export-tmpl-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckExportTemplateDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_export_template.test", "object_types.#", "1"),
					resource.TestCheckResourceAttrSet("data.netbox_export_template.test", "template_code"),
				),
			},
		},
	})
}

func testAccExportTemplateDataSourceConfig_byID(name string) string {
	return `
resource "netbox_export_template" "test" {
  name          = "` + name + `"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }}\n{% endfor %}"
}

data "netbox_export_template" "test" {
  id = netbox_export_template.test.id
}
`
}
func TestAccExportTemplateDataSource_byName(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("test-export-tmpl-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckExportTemplateDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateDataSourceConfig_byName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_export_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_export_template.test", "object_types.#", "1"),
					resource.TestCheckResourceAttrSet("data.netbox_export_template.test", "template_code"),
				),
			},
		},
	})
}

func TestAccExportTemplateDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("test-export-tmpl-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckExportTemplateDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccExportTemplateDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_export_template.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_export_template.test", "name", name),
				),
			},
		},
	})
}

func testAccExportTemplateDataSourceConfig_byName(name string) string {
	return `
resource "netbox_export_template" "test" {
  name          = "` + name + `"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }}\n{% endfor %}"
}

data "netbox_export_template" "test" {
  name = netbox_export_template.test.name
}
`
}
