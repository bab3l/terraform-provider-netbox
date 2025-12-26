package datasources_acceptance_tests

import (
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccScriptDataSource_IDPreservation(t *testing.T) {
	t.Parallel()
	t.Skip("Scripts cannot be created via API and require filesystem access on the server. Skipping test as we cannot ensure a script exists.")
}

func TestAccScriptDataSource_basic(t *testing.T) {

	t.Parallel()
	t.Skip("Scripts cannot be created via API and require filesystem access on the server. Skipping test as we cannot ensure a script exists.")
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
