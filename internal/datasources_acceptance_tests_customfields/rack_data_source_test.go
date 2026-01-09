//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackDataSource_customFields(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-ds-cf")
	siteSlug := testutil.GenerateSlug(siteName)
	rackName := testutil.RandomName("tf-test-rack-ds-cf")
	customFieldName := testutil.RandomCustomFieldName("tf_test_rack_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackDataSourceConfig_withCustomFields(siteName, siteSlug, rackName, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rack.test", "name", rackName),
					resource.TestCheckResourceAttr("data.netbox_rack.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_rack.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_rack.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_rack.test", "custom_fields.0.value", "datasource-test-value"),
				),
			},
		},
	})
}

func testAccRackDataSourceConfig_withCustomFields(siteName, siteSlug, rackName, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.rack"]
  type         = "text"
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_rack" "test" {
  name = %q
  site = netbox_site.test.name

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "datasource-test-value"
    }
  ]
}

data "netbox_rack" "test" {
  name = netbox_rack.test.name

  depends_on = [netbox_rack.test]
}
`, customFieldName, siteName, siteSlug, rackName)
}
