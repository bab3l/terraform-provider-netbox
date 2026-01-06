package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPSecPolicyDataSource_byID(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-policy-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyDataSourceByID(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ipsec_policy.test", "id"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecPolicyDestroy,
		),
	})
}

func TestAccIPSecPolicyDataSource_byName(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-policy-ds")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyDataSourceByName(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "name", randomName),
					resource.TestCheckResourceAttrSet("data.netbox_ipsec_policy.test", "id"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecPolicyDestroy,
		),
	})
}

func testAccIPSecPolicyDataSourceByID(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_policy" "test" {
  name = %[1]q
}

data "netbox_ipsec_policy" "test" {
  id = netbox_ipsec_policy.test.id
}
`, name)
}

func TestAccIPSecPolicyDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-policy-ds-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(randomName)

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyDataSourceByID(randomName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_ipsec_policy.test", "name", randomName),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecPolicyDestroy,
		),
	})
}

func testAccIPSecPolicyDataSourceByName(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_policy" "test" {
  name = %[1]q
}

data "netbox_ipsec_policy" "test" {
  name = netbox_ipsec_policy.test.name
}
`, name)
}
