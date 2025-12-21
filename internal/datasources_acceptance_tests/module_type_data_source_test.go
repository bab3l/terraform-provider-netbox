package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleTypeDataSource_basic(t *testing.T) {

	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module_type.test", "model", "Test Module Type"),
					resource.TestCheckResourceAttrSet("data.netbox_module_type.test", "manufacturer_id"),
				),
			},
		},
	})
}

const testAccModuleTypeDataSourceConfig = `
resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "Test Module Type"
}

data "netbox_module_type" "test" {
  id = netbox_module_type.test.id
}
`
