package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccIKEProposalDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ike-proposal-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEProposalCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalDataSourceByID(randomName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "name", randomName),
				),
			},
		},
	})
}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccIKEProposalDataSource_byID(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ike-proposal-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEProposalCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalDataSourceByID(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "authentication_method", "preshared-keys"),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "encryption_algorithm", "aes-256-cbc"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEProposalDestroy,
		),
	})
}

func TestAccIKEProposalDataSource_byName(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ike-proposal-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEProposalCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalDataSourceByName(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "authentication_method", "preshared-keys"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEProposalDestroy,
		),
	})
}

func testAccIKEProposalDataSourceByID(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %[1]q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

data "netbox_ike_proposal" "test" {
  id = netbox_ike_proposal.test.id
}
`, name)
}

func testAccIKEProposalDataSourceByName(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %[1]q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

data "netbox_ike_proposal" "test" {
  name = netbox_ike_proposal.test.name
}
`, name)
}
