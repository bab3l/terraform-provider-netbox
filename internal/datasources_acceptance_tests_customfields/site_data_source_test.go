//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSiteDataSource_customFields(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.GenerateSlug(siteName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create site with custom field and verify datasource returns it
			{
				Config: testAccSiteDataSourceConfig_withCustomFields(siteName, siteSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_site.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_site.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_site.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_site.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccSiteDataSourceConfig_withCustomFields(name, slug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.site"]
  type         = "text"
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_site" "test" {
  slug = netbox_site.test.slug

  depends_on = [netbox_site.test]
}
`, customFieldName, name, slug)
}
