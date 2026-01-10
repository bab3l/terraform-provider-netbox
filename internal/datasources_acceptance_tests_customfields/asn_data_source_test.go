//go:build customfields
// +build customfields

package datasources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccASNDataSource_customFields(t *testing.T) {
	asnNumber := int64(acctest.RandIntRange(64512, 64711))
	rirName := testutil.RandomName("tf-test-rir-ds-cf")
	rirSlug := testutil.GenerateSlug(rirName)
	customFieldName := testutil.RandomCustomFieldName("tf_test_asn_ds_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRIRCleanup(rirSlug)
	cleanup.RegisterASNCleanup(asnNumber)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create ASN with custom field and verify datasource returns it
			{
				Config: testAccASNDataSourceConfig_withCustomFields(asnNumber, rirName, rirSlug, customFieldName),
				Check: resource.ComposeTestCheckFunc(
					// Verify datasource returns the custom field
					resource.TestCheckResourceAttr("data.netbox_asn.test", "custom_fields.#", "1"),
					resource.TestCheckResourceAttr("data.netbox_asn.test", "custom_fields.0.name", customFieldName),
					resource.TestCheckResourceAttr("data.netbox_asn.test", "custom_fields.0.type", "text"),
					resource.TestCheckResourceAttr("data.netbox_asn.test", "custom_fields.0.value", "asn-datasource-test"),
				),
			},
		},
	})
}

func testAccASNDataSourceConfig_withCustomFields(asnNumber int64, rirName, rirSlug, customFieldName string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["ipam.asn"]
  type         = "text"
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.slug

  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "asn-datasource-test"
    }
  ]
}

data "netbox_asn" "test" {
  asn = netbox_asn.test.asn

  depends_on = [netbox_asn.test]
}
`, rirName, rirSlug, customFieldName, asnNumber)
}
