package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIKEPolicyResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEPolicyResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},
		},
	})

}

func TestAccIKEPolicyResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-full")

	proposalName := testutil.RandomName("tf-test-ike-proposal-for-policy")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEPolicyResourceConfig_full(name, proposalName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "1"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "mode", "aggressive"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", "Test IKE policy with full options"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "comments", "Test comments for IKE policy"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "proposals.#", "1"),
				),
			},
		},
	})

}

func TestAccIKEPolicyResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEPolicyResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},

			{

				Config: testAccIKEPolicyResourceConfig_updated(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", "Updated description"),
				),
			},
		},
	})

}

func TestAccIKEPolicyResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-policy")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEPolicyResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_ike_policy.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccIKEPolicyResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-ike-policy-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
				),
			},
		},
	})
}

func testAccIKEPolicyResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_policy" "test" {

  name    = %q

  version = 2

}

`, name)

}

func testAccIKEPolicyResourceConfig_full(name, proposalName string) string {

	return fmt.Sprintf(`

resource "netbox_ike_proposal" "test" {

  name                     = %q

  authentication_method    = "preshared-keys"

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"

  group                    = 14

}

resource "netbox_ike_policy" "test" {

  name        = %q

  version     = 1

  mode        = "aggressive"

  proposals   = [netbox_ike_proposal.test.id]

  description = "Test IKE policy with full options"

  comments    = "Test comments for IKE policy"

}

`, proposalName, name)

}

func testAccIKEPolicyResourceConfig_updated(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_policy" "test" {

  name        = %q

  version     = 2

  description = "Updated description"

}

`, name)

}

func TestAccConsistency_IKEPolicy_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-ike-policy-lit")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},
			{
				Config:   testAccIKEPolicyConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
				),
			},
		},
	})
}

func testAccIKEPolicyConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = 2
}
`, name)
}
