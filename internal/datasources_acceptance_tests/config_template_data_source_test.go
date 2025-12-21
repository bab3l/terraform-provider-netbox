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

// testAccConfigTemplateDataSourcePrereqs creates prerequisites for config template data source tests.

func testAccConfigTemplateDataSourcePrereqs(name, templateCode string) string {

	return fmt.Sprintf(`

resource "netbox_config_template" "test" {

  name          = %q

  template_code = %q

}

`, name, templateCode)

}

// testAccConfigTemplateDataSourceByID looks up a config template by ID.

func testAccConfigTemplateDataSourceByID(name, templateCode string) string {

	return testAccConfigTemplateDataSourcePrereqs(name, templateCode) + `

data "netbox_config_template" "test" {

  id = netbox_config_template.test.id

}

`

}

// testAccConfigTemplateDataSourceByName looks up a config template by name.

func testAccConfigTemplateDataSourceByName(name, templateCode string) string {

	return testAccConfigTemplateDataSourcePrereqs(name, templateCode) + fmt.Sprintf(`

data "netbox_config_template" "test" {

  name = %q

  depends_on = [netbox_config_template.test]

}

`, name)

}

func TestAccConfigTemplateDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("config-tmpl-ds")

	templateCode := "hostname {{ device.name }}"

	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckConfigTemplateDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccConfigTemplateDataSourceByID(name, templateCode),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_config_template.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_config_template.test", "template_code", templateCode),

					resource.TestCheckResourceAttrSet("data.netbox_config_template.test", "id"),
				),
			},
		},
	})

}

func TestAccConfigTemplateDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	name := testutil.RandomName("config-tmpl-ds")

	templateCode := "hostname {{ device.name }}"

	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckConfigTemplateDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccConfigTemplateDataSourceByName(name, templateCode),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_config_template.test", "name", name),

					resource.TestCheckResourceAttr("data.netbox_config_template.test", "template_code", templateCode),

					resource.TestCheckResourceAttrSet("data.netbox_config_template.test", "id"),
				),
			},
		},
	})

}
