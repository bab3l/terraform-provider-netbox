//go:build customfields
// +build customfields

package resources_acceptance_tests_customfields

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccCircuitTerminationResource_CustomFieldsPreservation(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields

	providerName := testutil.RandomName("tf-test-provider-ct-pres")
	providerSlug := testutil.RandomSlug("tf-test-provider-ct-pres")
	circuitTypeName := testutil.RandomName("tf-test-ct-ct-pres")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-ct-pres")
	circuitCID := testutil.RandomName("tf-test-circuit-ct-pres")
	siteName := testutil.RandomName("tf-test-site-ct-pres")
	siteSlug := testutil.RandomSlug("tf-test-site-ct-pres")
	cfName := testutil.RandomCustomFieldName("tf_ct_pres")

	cleanup := testutil.NewCleanupResource(t)
	defer cleanup.Close(t)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourcePreservationConfig_step1(
					providerName, providerSlug, circuitTypeName, circuitTypeSlug,
					circuitCID, siteName, siteSlug, cfName,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "custom_fields.%", "1"),
					testutil.ResourceCheckCustomFieldValue("netbox_circuit_termination.test", cfName, "preserved_value"),
				),
			},
			{
				// Update without custom_fields in config - should be preserved in NetBox
				Config: testAccCircuitTerminationResourcePreservationConfig_step2(
					providerName, providerSlug, circuitTypeName, circuitTypeSlug,
					circuitCID, siteName, siteSlug,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					// Custom fields are not in the config, so they won't appear in state
				),
			},
		},
	})
}

func testAccCircuitTerminationResourcePreservationConfig_step1(
	providerName, providerSlug, circuitTypeName, circuitTypeSlug,
	circuitCID, siteName, siteSlug, cfName string,
) string {
	return fmt.Sprintf(`
resource "netbox_custom_field" "circuit_termination_pres" {
  name         = %[8]q
  type         = "text"
  object_types = ["dcim.circuittermination"]
  required     = false
}

resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit" "test" {
  cid     = %[5]q
  provider = netbox_provider.test.id
  type    = netbox_circuit_type.test.id
  status  = "planned"
}

resource "netbox_site" "test" {
  name   = %[6]q
  slug   = %[7]q
  status = "active"
}

resource "netbox_circuit_termination" "test" {
  circuit      = netbox_circuit.test.id
  termination  = "a"
  site         = netbox_site.test.id

  custom_fields = {
    (netbox_custom_field.circuit_termination_pres.name) = "preserved_value"
  }

  depends_on = [netbox_custom_field.circuit_termination_pres]
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, cfName)
}

func testAccCircuitTerminationResourcePreservationConfig_step2(
	providerName, providerSlug, circuitTypeName, circuitTypeSlug,
	circuitCID, siteName, siteSlug string,
) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_circuit_type" "test" {
  name = %[3]q
  slug = %[4]q
}

resource "netbox_circuit" "test" {
  cid     = %[5]q
  provider = netbox_provider.test.id
  type    = netbox_circuit_type.test.id
  status  = "active"
}

resource "netbox_site" "test" {
  name   = %[6]q
  slug   = %[7]q
  status = "active"
}

resource "netbox_circuit_termination" "test" {
  circuit      = netbox_circuit.test.id
  termination  = "a"
  site         = netbox_site.test.id
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug)
}
