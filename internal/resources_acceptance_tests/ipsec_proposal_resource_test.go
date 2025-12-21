package resources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSECProposalResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProposalResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
				),
			},
		},
	})

}

func TestAccIPSECProposalResource_full(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop-full")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProposalResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-128-cbc"),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha256"),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "description", "Test IPsec proposal"),
				),
			},
		},
	})

}

func TestAccIPSECProposalResource_update(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop-update")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProposalResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
				),
			},

			{

				Config: testAccIPSECProposalResourceConfig_full(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-128-cbc"),

					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "description", "Test IPsec proposal"),
				),
			},
		},
	})

}

func TestAccIPSECProposalResource_import(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-ipsec-prop")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPSECProposalResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
				),
			},

			{

				ResourceName: "netbox_ipsec_proposal.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccIPSECProposalResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_proposal" "test" {

  name = %q

}

`, name)

}

func testAccIPSECProposalResourceConfig_full(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_proposal" "test" {

  name                     = %q

  encryption_algorithm     = "aes-128-cbc"

  authentication_algorithm = "hmac-sha256"

  description              = "Test IPsec proposal"

}

`, name)

}
