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

func TestAccIPSecProfileDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-profile-ds-id")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile-ds-id")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProfileDataSourceByID(randomName, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "name", randomName),
				),
			},
		},
	})
}

func TestAccIPSecProfileDataSource_byID(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-profile-ds")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile-ds")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProfileDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProfileDataSourceByID(randomName, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
		},
	})
}

func TestAccIPSecProfileDataSource_byName(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-profile-ds")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile-ds")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProfileDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProfileDataSourceByName(randomName, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
		},
	})
}

func testAccIPSecProfileDataSourceByID(name, ikePolicyName, ipsecPolicyName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %[2]q
  version = 2
}

resource "netbox_ipsec_policy" "test" {
  name = %[3]q
}

resource "netbox_ipsec_profile" "test" {
  name         = %[1]q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}

data "netbox_ipsec_profile" "test" {
  id = netbox_ipsec_profile.test.id
}
`, name, ikePolicyName, ipsecPolicyName)
}

func testAccIPSecProfileDataSourceByName(name, ikePolicyName, ipsecPolicyName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %[2]q
  version = 2
}

resource "netbox_ipsec_policy" "test" {
  name = %[3]q
}

resource "netbox_ipsec_profile" "test" {
  name         = %[1]q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}

data "netbox_ipsec_profile" "test" {
  name = netbox_ipsec_profile.test.name
}
`, name, ikePolicyName, ipsecPolicyName)
}
