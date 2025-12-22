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

func TestAccConfigContextResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-config-context")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigContextResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_context.test", "data", "{\"foo\":\"bar\"}"),
				),
			},

			{

				ResourceName: "netbox_config_context.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccConfigContextResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_config_context" "test" {

  name = %q

  data = "{\"foo\":\"bar\"}"

}

`, name)

}
