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

func TestAccVirtualChassisResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-vc")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualChassisResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_virtual_chassis.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccVirtualChassisResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-vc-full")

	description := testutil.RandomName("description")

	updatedDescription := testutil.RandomName("description")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualChassisResourceConfig_full(name, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", "test-domain"),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", description),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "comments", "Test comments"),
				),
			},

			{

				Config: testAccVirtualChassisResourceConfig_full(name, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccVirtualChassisResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vc-id")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualChassisResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),
					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
				),
			},
		},
	})

}

func testAccVirtualChassisResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_virtual_chassis" "test" {

  name = %q

}

`, name)

}

func testAccVirtualChassisResourceConfig_full(name, description string) string {

	return fmt.Sprintf(`

resource "netbox_virtual_chassis" "test" {

  name        = %q

  domain      = "test-domain"

  description = %q

  comments    = "Test comments"

}

`, name, description)

}
