package datasources_acceptance_tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCustomFieldDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	// Custom field names only allow alphanumeric and underscores
	name := strings.ReplaceAll(testutil.RandomName("tf_test_cf_ds_id"), "-", "_")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_custom_field.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_custom_field.test", "name", name),
				),
			},
		},
	})
}

func TestAccCustomFieldDataSource_basic(t *testing.T) {
	t.Parallel()

	// Custom field names only allow alphanumeric and underscores
	name := strings.ReplaceAll(testutil.RandomName("tf_test_cf_ds"), "-", "_")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_field.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_custom_field.test", "type", "text"),
				),
			},
		},
	})
}

func TestAccCustomFieldDataSource_byName(t *testing.T) {
	t.Parallel()

	// Custom field names only allow alphanumeric and underscores
	name := strings.ReplaceAll(testutil.RandomName("tf_test_cf_ds"), "-", "_")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldDataSourceConfigByName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_field.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_custom_field.test", "type", "text"),
				),
			},
		},
	})
}

func testAccCustomFieldDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = "%s"
  type         = "text"
  object_types = ["dcim.site"]
}

data "netbox_custom_field" "test" {
  id = netbox_custom_field.test.id
}
`, name)
}

func testAccCustomFieldDataSourceConfigByName(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "test" {
  name         = "%s"
  type         = "text"
  object_types = ["dcim.site"]
}

data "netbox_custom_field" "test" {
  name = netbox_custom_field.test.name
}
`, name)
}
