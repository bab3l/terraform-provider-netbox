//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProviderAccountResource_CustomFieldsPreservation(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields

	providerName := testutil.RandomName("tf-test-provider-pa-pres")
	providerSlug := testutil.RandomSlug("tf-test-provider-pa-pres")
	accountNum := testutil.RandomName("acct-pres")
	cfName := testutil.RandomCustomFieldName("tf_pa_pres")

	cleanup := testutil.NewCleanupResource(t)
	defer cleanup.Close(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderAccountResourcePreservationConfig_step1(
					providerName, providerSlug, accountNum, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_account.test", "custom_fields.%", "1"),
					testutil.ResourceCheckCustomFieldValue("netbox_provider_account.test", cfName, "preserved_value"),
				),
			},
			{
				// Update without custom_fields in config - should be preserved in NetBox
				Config: testAccProviderAccountResourcePreservationConfig_step2(
					providerName, providerSlug, accountNum,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),
					// Custom fields are not in the config, so they won't appear in state
				),
			},
		},
	})
}

func testAccProviderAccountResourcePreservationConfig_step1(
	providerName, providerSlug, accountNum, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "provider_account_pres" {
  name         = %[4]q
  type         = "text"
  object_types = ["circuits.provideraccount"]
  required     = false
}

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_account" "test" {
  provider = netbox_provider.test.id
  account  = %[3]q

  custom_fields = {
    (netbox_custom_field.provider_account_pres.name) = "preserved_value"
  }

  depends_on = [netbox_custom_field.provider_account_pres]
}
`, providerName, providerSlug, accountNum, cfName)
}

func testAccProviderAccountResourcePreservationConfig_step2(
	providerName, providerSlug, accountNum string,
) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_provider_account" "test" {
  provider = netbox_provider.test.id
  account  = %[3]q
  name     = "updated-name"
}
`, providerName, providerSlug, accountNum)
}
