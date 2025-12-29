package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigContextDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("cc-ds-id-context")

	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckConfigContextDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_config_context.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "description", "Test Config Context Description"),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "weight", "100"),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "is_active", "true"),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "data", "{\"foo\":\"bar\"}"),
				),
			},
		},
	})
}

func TestAccConfigContextDataSource_byID(t *testing.T) {

	t.Parallel()

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("test-config-context")

	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckConfigContextDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextDataSourceConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "description", "Test Config Context Description"),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "weight", "100"),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "is_active", "true"),
					resource.TestCheckResourceAttr("data.netbox_config_context.test", "data", "{\"foo\":\"bar\"}"),
				),
			},
		},
	})
}

func testAccConfigContextDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_config_context" "test" {
  name        = %q
  description = "Test Config Context Description"
  weight      = 100
  is_active   = true
  data        = "{\"foo\":\"bar\"}"
}

data "netbox_config_context" "test" {
  id = netbox_config_context.test.id
}
`, name)
}
