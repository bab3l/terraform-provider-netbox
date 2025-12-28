package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccConfigContextResource_basic(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-config-context")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigContextResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_context.test", "data", "{\"foo\":\"bar\"}"),
				),
			},

			{

				ResourceName: "netbox_config_context.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConsistency_ConfigContext_LiteralNames(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("tf-test-config-context-lit")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigContextCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccConfigContextConsistencyLiteralNamesConfig(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),

					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),

					resource.TestCheckResourceAttr("netbox_config_context.test", "data", "{\"foo\":\"bar\"}"),
				),
			},

			{

				Config: testAccConfigContextConsistencyLiteralNamesConfig(name),

				PlanOnly: true,

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
				),
			},
		},
	})

}

func testAccConfigContextConsistencyLiteralNamesConfig(name string) string {

	return fmt.Sprintf(`

resource "netbox_config_context" "test" {

  name = %q

  data = "{\"foo\":\"bar\"}"

}

`, name)

}

func TestAccConfigContextResource_IDPreservation(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("config-context")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_context.test", "name", name),
				),
			},
		},
	})
}

func TestAccConfigContextResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-ctx-upd")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_withDescription(name, testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccConfigContextResourceConfig_withDescription(name, testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_context.test", "description", testutil.Description2),
				),
			},
		},
	})
}

func TestAccConfigContextResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-ctx-extdel")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigContextResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_context.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find config context by name
					items, _, err := client.ExtrasAPI.ExtrasConfigContextsList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list config contexts: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Config context not found with name: %s", name)
					}

					// Delete the config context
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasConfigContextsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config context: %v", err)
					}

					t.Logf("Successfully externally deleted config context with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccConfigContextResourceConfig_basic(name string) string {

	return fmt.Sprintf(`

resource "netbox_config_context" "test" {

  name = %q

  data = "{\"foo\":\"bar\"}"

}

`, name)

}

func testAccConfigContextResourceConfig_withDescription(name string, description string) string {
	return fmt.Sprintf(`
resource "netbox_config_context" "test" {
  name        = %[1]q
  data        = "{\"key\": \"value\"}"
  description = %[2]q
}
`, name, description)
}
