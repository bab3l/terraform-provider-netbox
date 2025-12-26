package datasources_acceptance_tests

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderAccountDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-id")
	providerSlug := testutil.RandomSlug("tf-test-prov-id")
	accountName := testutil.RandomName("tf-test-acct-id")
	accountNumber := testutil.RandomName("account-id")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountDataSourceConfig(providerName, providerSlug, accountName, accountNumber),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.netbox_provider_account.test", "id"),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "name", accountName),
				),
			},
		},
	})
}

func TestAccProviderAccountDataSource_basic(t *testing.T) {

	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-prov")
	accountName := testutil.RandomName("tf-test-acct")
	accountNumber := testutil.RandomName("account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountDataSourceConfig(providerName, providerSlug, accountName, accountNumber),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "name", accountName),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "account", accountNumber),
				),
			},
		},
	})
}

func TestAccProviderAccountDataSource_byAccount(t *testing.T) {

	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-prov")
	accountName := testutil.RandomName("tf-test-acct")
	accountNumber := testutil.RandomName("account")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountDataSourceConfigByAccount(providerName, providerSlug, accountName, accountNumber),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "name", accountName),
					resource.TestCheckResourceAttr("data.netbox_provider_account.test", "account", accountNumber),
				),
			},
		},
	})
}

func testAccProviderAccountDataSourceConfig(providerName, providerSlug, accountName, accountNumber string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_account" "test" {
  circuit_provider = netbox_provider.test.id
  account          = %[4]q
  name             = %[3]q
}

data "netbox_provider_account" "test" {
  id = netbox_provider_account.test.id
}
`, providerName, providerSlug, accountName, accountNumber)
}

func testAccProviderAccountDataSourceConfigByAccount(providerName, providerSlug, accountName, accountNumber string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_account" "test" {
  circuit_provider = netbox_provider.test.id
  account          = %[4]q
  name             = %[3]q
}

data "netbox_provider_account" "test" {
  account = netbox_provider_account.test.account
}
`, providerName, providerSlug, accountName, accountNumber)
}
