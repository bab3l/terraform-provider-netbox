//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccPlatformDataSource_customFields(t *testing.T) {
	platformName := testutil.RandomName("tf-test-platform-ds-cf")
	platformSlug := testutil.GenerateSlug(platformName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_platform_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPlatformDataSourceConfig_withCustomFields(platformName, platformSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_platform.test", "name", platformName),
					// Platform resource doesn't support custom_fields, so we just verify the datasource doesn't error
					// and returns the custom_fields attribute (even if empty)
					resource.TestCheckResourceAttrSet("data.netbox_platform.test", "custom_fields.#"),
				),
			},
		},
	})
}

func testAccPlatformDataSourceConfig_withCustomFields(platformName, platformSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["dcim.platform"]
  type         = "text"
}

resource "netbox_platform" "test" {
  name = %q
  slug = %q
}

data "netbox_platform" "test" {
  name = netbox_platform.test.name

  depends_on = [netbox_platform.test]
}
`, customFieldName, platformName, platformSlug)
}
