package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccModuleTypeDataSource_basic(t *testing.T) {

	t.Parallel()

	mfgName := testutil.RandomName("mfg")
	mfgSlug := testutil.RandomSlug("mfg")
	modelName := testutil.RandomName("module-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeDataSourceConfig(mfgName, mfgSlug, modelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module_type.test", "model", modelName),
					resource.TestCheckResourceAttrSet("data.netbox_module_type.test", "manufacturer_id"),
				),
			},
		},
	})
}

func testAccModuleTypeDataSourceConfig(mfgName, mfgSlug, modelName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
}

data "netbox_module_type" "test" {
  id = netbox_module_type.test.id
}
`, mfgName, mfgSlug, modelName)
}
