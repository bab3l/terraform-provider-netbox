package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const scriptDataSourceAcceptanceWaiver = "Skipping for 1.0 readiness: NetBox scripts are filesystem-managed server fixtures and cannot be created through the API, so this datasource cannot yet be exercised in portable acceptance automation."

func TestAccScriptDataSource_IDPreservation(t *testing.T) {
	t.Parallel()

	t.Skip(scriptDataSourceAcceptanceWaiver)
}

func TestAccScriptDataSource_basic(t *testing.T) {
	t.Parallel()

	t.Skip(scriptDataSourceAcceptanceWaiver)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "netbox_script" "test" {
					name = "nonexistent"
				}`,
				ExpectError: nil, // We expect it to fail if we ran it, but we skip it.
			},
		},
	})
}

func TestAccScriptDataSource_byID(t *testing.T) {
	t.Parallel()

	t.Skip(scriptDataSourceAcceptanceWaiver)
}

func TestAccScriptDataSource_byName(t *testing.T) {
	t.Parallel()

	t.Skip(scriptDataSourceAcceptanceWaiver)
}
