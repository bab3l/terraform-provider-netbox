package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServiceTemplateDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("test-svc-tmpl-ds")

	cleanup.RegisterServiceTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckServiceTemplateDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "protocol", "tcp"),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "ports.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "ports.0", "80"),
				),
			},
		},
	})
}

func TestAccServiceTemplateDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("test-svc-tmpl-ds")

	cleanup.RegisterServiceTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckServiceTemplateDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceTemplateDataSourceConfig_byName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_service_template.test", "protocol", "tcp"),
				),
			},
		},
	})
}

func testAccServiceTemplateDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}

data "netbox_service_template" "test" {
  id = netbox_service_template.test.id
}
`, name)
}

func testAccServiceTemplateDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_service_template" "test" {
  name     = %q
  protocol = "tcp"
  ports    = [80]
}

data "netbox_service_template" "test" {
  name = netbox_service_template.test.name

  depends_on = [netbox_service_template.test]
}
`, name)
}
