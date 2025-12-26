package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRIRDataSource_IDPreservation(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("rir-id")
	rirSlug := testutil.RandomSlug("rir-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRIRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRIRDataSourceConfig(rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_rir.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_rir.test", "name", rirName),
				),
			},
		},
	})
}

func TestAccRIRDataSource_basic(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRIRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRIRDataSourceConfig(rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rir.test", "name", rirName),
					resource.TestCheckResourceAttr("data.netbox_rir.test", "slug", rirSlug),
				),
			},
		},
	})
}

func TestAccRIRDataSource_byName(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRIRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRIRDataSourceConfigByName(rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rir.test", "name", rirName),
					resource.TestCheckResourceAttr("data.netbox_rir.test", "slug", rirSlug),
				),
			},
		},
	})
}

func TestAccRIRDataSource_bySlug(t *testing.T) {

	t.Parallel()

	rirName := testutil.RandomName("rir")
	rirSlug := testutil.RandomSlug("rir")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckRIRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRIRDataSourceConfigBySlug(rirName, rirSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_rir.test", "name", rirName),
					resource.TestCheckResourceAttr("data.netbox_rir.test", "slug", rirSlug),
				),
			},
		},
	})
}

func testAccRIRDataSourceConfig(rirName, rirSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_rir" "test" {
  id = netbox_rir.test.id
}
`, rirName, rirSlug)
}

func testAccRIRDataSourceConfigByName(rirName, rirSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_rir" "test" {
  name = netbox_rir.test.name
}
`, rirName, rirSlug)
}

func testAccRIRDataSourceConfigBySlug(rirName, rirSlug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = "%s"
  slug = "%s"
}

data "netbox_rir" "test" {
  slug = netbox_rir.test.slug
}
`, rirName, rirSlug)
}
