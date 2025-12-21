package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIKEProposalResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-proposal")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEProposalResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_method", "preshared-keys"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-256-cbc"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_algorithm", "hmac-sha256"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "group", "14"),
				),
			},
		},
	})

}

func TestAccIKEProposalResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-proposal-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEProposalResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_method", "certificates"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-128-gcm"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_algorithm", "hmac-sha512"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "group", "19"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "sa_lifetime", "28800"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "description", "Test IKE proposal with full options"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "comments", "Test comments for IKE proposal"),
				),
			},
		},
	})

}

func TestAccIKEProposalResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-proposal-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEProposalResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-256-cbc"),
				),
			},

			{

				Config: testAccIKEProposalResourceConfig_updated(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-128-cbc"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "description", "Updated description"),
				),
			},
		},
	})

}

func TestAccIKEProposalResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ike-proposal")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIKEProposalResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_ike_proposal.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccIKEProposalResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_proposal" "test" {

  name                     = %q

  authentication_method    = "preshared-keys"

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"

  group                    = 14

}

`, name)

}

func testAccIKEProposalResourceConfig_full(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_proposal" "test" {

  name                     = %q

  authentication_method    = "certificates"

  encryption_algorithm     = "aes-128-gcm"

  authentication_algorithm = "hmac-sha512"

  group                    = 19

  sa_lifetime              = 28800

  description              = "Test IKE proposal with full options"

  comments                 = "Test comments for IKE proposal"

}

`, name)

}

func testAccIKEProposalResourceConfig_updated(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_proposal" "test" {

  name                     = %q

  authentication_method    = "preshared-keys"

  encryption_algorithm     = "aes-128-cbc"

  authentication_algorithm = "hmac-sha256"

  group                    = 14

  description              = "Updated description"

}

`, name)

}
