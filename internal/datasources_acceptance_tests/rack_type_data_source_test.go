package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRackTypeDataSource_basic(t *testing.T) {

	t.Parallel()

	mfrName := testutil.RandomName("mfr")
	mfrSlug := testutil.RandomSlug("mfr")
	model := testutil.RandomName("rack-type")
	slug := testutil.RandomSlug("rack-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfrSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckManufacturerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeDataSourceConfig(mfrName, mfrSlug, model, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("data.netbox_rack_type.test", "slug", slug),
				),
			},
		},
	})
}

func testAccRackTypeDataSourceConfig(mfrName, mfrSlug, model, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = "%s"
  slug = "%s"
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = "%s"
  slug         = "%s"
  width        = 19
  u_height     = 42
  form_factor  = "2-post-frame"
}

data "netbox_rack_type" "test" {
  id = netbox_rack_type.test.id
}
`, mfrName, mfrSlug, model, slug)
}
