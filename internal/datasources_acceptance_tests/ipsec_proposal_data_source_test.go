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

func TestAccIPSecProposalDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("tf-test-ipsec-proposal-ds")

	cleanup.RegisterIPSecProposalCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecProposalDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "encryption_algorithm", "aes-256-cbc"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha256"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProposalDestroy,
		),
	})

}

func TestAccIPSecProposalDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("tf-test-ipsec-proposal-ds")

	cleanup.RegisterIPSecProposalCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecProposalDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "encryption_algorithm", "aes-256-cbc"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProposalDestroy,
		),
	})

}

func testAccIPSecProposalDataSourceByID(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_proposal" "test" {

  name                     = %[1]q

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"

}

data "netbox_ipsec_proposal" "test" {

  id = netbox_ipsec_proposal.test.id

}

`, name)

}

func testAccIPSecProposalDataSourceByName(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_proposal" "test" {

  name                     = %[1]q

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"

}

data "netbox_ipsec_proposal" "test" {

  name = netbox_ipsec_proposal.test.name

}

`, name)

}
