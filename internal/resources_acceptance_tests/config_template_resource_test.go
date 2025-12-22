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

const defaultTemplateCode = "hostname {{ device.name }}"

func TestAccConfigTemplateResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("config-tmpl")

	templateCode := defaultTemplateCode

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)

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

	t.Parallel()

	name := testutil.RandomName("config-tmpl")

	templateCode := defaultTemplateCode

	description := testutil.RandomName("description")

	updatedName := testutil.RandomName("config-tmpl-updated")

	updatedTemplateCode := "hostname {{ device.name }}\ninterface {{ interface.name }}"

	updatedDescription := "Updated test config template"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)
	cleanup.RegisterConfigTemplateCleanup(updatedName)

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

func TestAccConsistency_ConfigTemplate_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("config-tmpl-lit")

	templateCode := defaultTemplateCode

	description := "Test template"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigTemplateConsistencyLiteralNamesConfig(name, templateCode, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),

					resource.TestCheckResourceAttr("netbox_config_template.test", "description", description),
				),
			},

			{

				Config: testAccConfigTemplateConsistencyLiteralNamesConfig(name, templateCode, description),

				PlanOnly: true,

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),
				),
			},
		},
	})

}

func testAccConfigTemplateConsistencyLiteralNamesConfig(name, templateCode, description string) string {

	return fmt.Sprintf(`

resource "netbox_config_template" "test" {

  name          = %q

  template_code = %q

  description   = %q

}

`, name, templateCode, description)

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
