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

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccIKEPolicyDataSource_byID(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ike-policy-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyDataSourceByID(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "version", "2"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEPolicyDestroy,
		),
	})
}

func TestAccIKEPolicyDataSource_byName(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ike-policy-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyDataSourceByName(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "version", "2"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEPolicyDestroy,
		),
	})
}

func testAccIKEPolicyDataSourceByID(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %[1]q
  version = 2
}

data "netbox_ike_policy" "test" {
  id = netbox_ike_policy.test.id
}
`, name)
}

func TestAccIKEPolicyDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ike-policy-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyDataSourceByID(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ike_policy.test", "name", randomName),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEPolicyDestroy,
		),
	})
}

func testAccIKEPolicyDataSourceByName(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %[1]q
  version = 2
}

data "netbox_ike_policy" "test" {
  name = netbox_ike_policy.test.name
}
`, name)
}
