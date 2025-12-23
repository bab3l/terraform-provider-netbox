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

func TestAccModuleTypeDataSource_byModelAndManufacturer(t *testing.T) {

	t.Parallel()

	mfgName1 := testutil.RandomName("mfg1")
	mfgSlug1 := testutil.RandomSlug("mfg1")
	mfgName2 := testutil.RandomName("mfg2")
	mfgSlug2 := testutil.RandomSlug("mfg2")
	modelName := testutil.RandomName("module-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug1)
	cleanup.RegisterManufacturerCleanup(mfgSlug2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccModuleTypeDataSourceConfigByModelAndManufacturer(mfgName1, mfgSlug1, mfgName2, mfgSlug2, modelName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_module_type.by_model", "model", modelName),
					resource.TestCheckResourceAttrSet("data.netbox_module_type.by_model", "manufacturer_id"),
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

func testAccModuleTypeDataSourceConfigByModelAndManufacturer(mfgName1, mfgSlug1, mfgName2, mfgSlug2, modelName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test1" {
  name = "%s"
  slug = "%s"
}

resource "netbox_manufacturer" "test2" {
  name = "%s"
  slug = "%s"
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test1.id
  model        = "%s"
}

data "netbox_module_type" "by_model" {
  model           = netbox_module_type.test.model
  manufacturer_id = netbox_manufacturer.test1.id
}
`, mfgName1, mfgSlug1, mfgName2, mfgSlug2, modelName)
}
