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

func TestAccRIRResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-rir")

	slug := testutil.RandomSlug("tf-test-rir")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRIRResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rir.test", "id"),

					resource.TestCheckResourceAttr("netbox_rir.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rir.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_rir.test", "is_private", "false"),
				),
			},

			{

				ResourceName: "netbox_rir.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccRIRResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-rir-full")

	slug := testutil.RandomSlug("tf-test-rir-full")

	description := testutil.RandomName("description")

	updatedDescription := "Updated RIR description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRIRResourceConfig_full(name, slug, description, true),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rir.test", "id"),

					resource.TestCheckResourceAttr("netbox_rir.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rir.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_rir.test", "description", description),

					resource.TestCheckResourceAttr("netbox_rir.test", "is_private", "true"),
				),
			},

			{

				Config: testAccRIRResourceConfig_full(name, slug, updatedDescription, false),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_rir.test", "is_private", "false"),
				),
			},
		},
	})

}

func testAccRIRResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func TestAccConsistency_RIR_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-rir-lit")
	slug := testutil.RandomSlug("tf-test-rir-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRIRConsistencyLiteralNamesConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rir.test", "id"),
					resource.TestCheckResourceAttr("netbox_rir.test", "name", name),
					resource.TestCheckResourceAttr("netbox_rir.test", "slug", slug),
				),
			},
			{
				Config:   testAccRIRConsistencyLiteralNamesConfig(name, slug),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rir.test", "id"),
				),
			},
		},
	})
}

func testAccRIRConsistencyLiteralNamesConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccRIRResourceConfig_full(name, slug, description string, isPrivate bool) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name        = %q

  slug        = %q

  description = %q

  is_private  = %t

}

`, name, slug, description, isPrivate)

}
