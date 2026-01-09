//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocationDataSource_customFields(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.GenerateSlug(siteName)
	locationName := testutil.RandomName("tf-test-location-ds-cf")
	locationSlug := testutil.GenerateSlug(locationName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_location_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccLocationDataSourceConfig_withCustomFields(siteName, siteSlug, locationName, locationSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_location.test", "name", locationName),
					resource.TestCheckResourceAttr("data.netbox_location.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_location.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_location.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_location.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccLocationDataSourceConfig_withCustomFields(siteName, siteSlug, locationName, locationSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.location"]
  type         = "text"
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_location" "test" {
  name = %q
  slug = %q
  site = netbox_site.test.name

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_location" "test" {
  name = netbox_location.test.name

  depends_on = [netbox_location.test]
}
`, customFieldName, siteName, siteSlug, locationName, locationSlug)
}
