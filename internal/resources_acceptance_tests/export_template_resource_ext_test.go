package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
)

// TestAccExportTemplateResource_removeOptionalFields_extended tests removing additional optional fields.
func TestAccExportTemplateResource_removeOptionalFields_extended(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("export-tmpl-rem")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterExportTemplateCleanup(name)

	testFields := map[string]string{
		"mime_type":      "text/csv",
		"file_extension": "csv",
	}

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_export_template",
		BaseConfig: func() string {
			return testAccExportTemplateResourceConfig_removeOptionalFields_base(name)
		},
		ConfigWithFields: func() string {
			return testAccExportTemplateResourceConfig_removeOptionalFields_withFields(name, testFields)
		},
		OptionalFields: testFields,
		RequiredFields: map[string]string{
			"name": name,
		},
	})
}

func testAccExportTemplateResourceConfig_removeOptionalFields_base(name string) string {
	return fmt.Sprintf(`
resource "netbox_export_template" "test" {
  name           = %[1]q
  object_types   = ["dcim.site"]
  template_code  = "name,slug\n{%% for site in queryset %%}{{ site.name }},{{ site.id }}\n{%% endfor %%}"
}
`, name)
}

func testAccExportTemplateResourceConfig_removeOptionalFields_withFields(name string, fields map[string]string) string {
	return fmt.Sprintf(`
resource "netbox_export_template" "test" {
  name           = %[1]q
  object_types   = ["dcim.site"]
  template_code  = "name,slug\n{%% for site in queryset %%}{{ site.name }},{{ site.id }}\n{%% endfor %%}"
  mime_type      = %[2]q
  file_extension = %[3]q
}
`, name, fields["mime_type"], fields["file_extension"])
}
