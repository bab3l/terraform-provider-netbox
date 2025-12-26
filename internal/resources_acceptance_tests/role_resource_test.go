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

func TestAccRoleResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-role")

	slug := testutil.RandomSlug("tf-test-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "1000"),
				),
			},

			{

				ResourceName: "netbox_role.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccRoleResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-role-full")

	slug := testutil.RandomSlug("tf-test-role-full")

	description := testutil.RandomName("description")

	updatedDescription := "Updated IPAM role description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRoleResourceConfig_full(name, slug, description, 100),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_role.test", "description", description),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "100"),
				),
			},

			{

				Config: testAccRoleResourceConfig_full(name, slug, updatedDescription, 200),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_role.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "200"),
				),
			},
		},
	})

}

func TestAccRoleResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-role-id")
	slug := testutil.RandomSlug("tf-test-role-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRoleResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
				),
			},
		},
	})
}

func testAccRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_role" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccRoleResourceConfig_full(name, slug, description string, weight int) string {

	return fmt.Sprintf(`

resource "netbox_role" "test" {

  name        = %q

  slug        = %q

  description = %q

  weight      = %d

}

`, name, slug, description, weight)

}
func TestAccConsistency_Role_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-role-lit")
	slug := testutil.RandomSlug("tf-test-role-lit")
	description := testutil.RandomName("description")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRoleConsistencyLiteralNamesConfig(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
					resource.TestCheckResourceAttr("netbox_role.test", "name", name),
					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_role.test", "description", description),
				),
			},
			{
				Config:   testAccRoleConsistencyLiteralNamesConfig(name, slug, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),
				),
			},
		},
	})
}

func testAccRoleConsistencyLiteralNamesConfig(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_role" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}
