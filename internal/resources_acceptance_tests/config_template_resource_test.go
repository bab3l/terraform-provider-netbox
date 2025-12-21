package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigTemplateResource_basic(t *testing.T) {

	name := testutil.RandomName("config-tmpl")

	templateCode := "hostname {{ device.name }}"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigTemplateResourceConfig_basic(name, templateCode),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),
				),
			},

			{

				ResourceName: "netbox_config_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConfigTemplateResource_full(t *testing.T) {

	name := testutil.RandomName("config-tmpl")

	templateCode := "hostname {{ device.name }}"

	description := "Test config template"

	updatedName := testutil.RandomName("config-tmpl-updated")

	updatedTemplateCode := "hostname {{ device.name }}\ninterface {{ interface.name }}"

	updatedDescription := "Updated test config template"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigTemplateResourceConfig_full(name, templateCode, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),

					resource.TestCheckResourceAttr("netbox_config_template.test", "description", description),
				),
			},

			{

				Config: testAccConfigTemplateResourceConfig_full(updatedName, updatedTemplateCode, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_config_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", updatedTemplateCode),

					resource.TestCheckResourceAttr("netbox_config_template.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccConfigTemplateResourceConfig_basic(name, templateCode string) string {

	return fmt.Sprintf(`

resource "netbox_config_template" "test" {

  name          = %q

  template_code = %q

}

`, name, templateCode)

}

func testAccConfigTemplateResourceConfig_full(name, templateCode, description string) string {

	return fmt.Sprintf(`

resource "netbox_config_template" "test" {

  name          = %q

  template_code = %q

  description   = %q

}

`, name, templateCode, description)

}
