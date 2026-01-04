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

const defaultTemplateCode = "hostname {{ device.name }}"

func TestAccConfigTemplateResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("config-tmpl")
	templateCode := defaultTemplateCode

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateResourceConfig_basic(name, templateCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),
				),
			},
			{
				ResourceName:      "netbox_config_template.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config:   testAccConfigTemplateResourceConfig_basic(name, templateCode),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConfigTemplateResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("cfg-tmpl-id")
	templateCode := defaultTemplateCode

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateResourceConfig_basic(name, templateCode),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),
				),
			},
			{
				Config:   testAccConfigTemplateResourceConfig_basic(name, templateCode),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConfigTemplateResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("config-tmpl")
	templateCode := defaultTemplateCode
	description := testutil.RandomName("description")
	updatedName := testutil.RandomName("config-tmpl-updated")
	updatedTemplateCode := "hostname {{ device.name }}\ninterface {{ interface.name }}"
	updatedDescription := "Updated test config template"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)
	cleanup.RegisterConfigTemplateCleanup(updatedName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateResourceConfig_full(name, templateCode, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),
					resource.TestCheckResourceAttr("netbox_config_template.test", "description", description),
				),
			},
			{
				Config: testAccConfigTemplateResourceConfig_full(updatedName, updatedTemplateCode, updatedDescription),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_template.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", updatedTemplateCode),
					resource.TestCheckResourceAttr("netbox_config_template.test", "description", updatedDescription),
				),
			},
			{
				Config:   testAccConfigTemplateResourceConfig_full(updatedName, updatedTemplateCode, updatedDescription),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConsistency_ConfigTemplate_LiteralNames(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("config-tmpl-lit")
	templateCode := defaultTemplateCode
	description := "Test template"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateResourceConfig_full(name, templateCode, description),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_config_template.test", "name", name),
					resource.TestCheckResourceAttr("netbox_config_template.test", "template_code", templateCode),
					resource.TestCheckResourceAttr("netbox_config_template.test", "description", description),
				),
			},
			{
				Config:   testAccConfigTemplateResourceConfig_full(name, templateCode, description),
				PlanOnly: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),
				),
			},
		},
	})
}

func TestAccConfigTemplateResource_update(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-tmpl-upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateResourceConfig_full(name, "{{ device.name }}", testutil.Description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_template.test", "description", testutil.Description1),
				),
			},
			{
				Config: testAccConfigTemplateResourceConfig_full(name, "{{ device.name }}", testutil.Description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_config_template.test", "description", testutil.Description2),
				),
			},
			{
				Config:   testAccConfigTemplateResourceConfig_full(name, "{{ device.name }}", testutil.Description2),
				PlanOnly: true,
			},
		},
	})
}

func TestAccConfigTemplateResource_externalDeletion(t *testing.T) {
	t.Parallel()
	testutil.TestAccPreCheck(t)

	name := testutil.RandomName("tf-test-tmpl-extdel")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterConfigTemplateCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTemplateResourceConfig_basic(name, defaultTemplateCode),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_config_template.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}

					// Find config template by name
					items, _, err := client.ExtrasAPI.ExtrasConfigTemplatesList(context.Background()).Name([]string{name}).Execute()
					if err != nil {
						t.Fatalf("Failed to list config templates: %v", err)
					}
					if items == nil || len(items.Results) == 0 {
						t.Fatalf("Config template not found with name: %s", name)
					}

					// Delete the config template
					itemID := items.Results[0].Id
					_, err = client.ExtrasAPI.ExtrasConfigTemplatesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to delete config template: %v", err)
					}

					t.Logf("Successfully externally deleted config template with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccConfigTemplateResourceConfig_basic(name, templateCode string) string {
	return fmt.Sprintf(`
resource "netbox_config_template" "test" {
  name          = %q
  template_code = %q
}
`, name, templateCode)
}

func testAccConfigTemplateResourceConfig_full(name, templateCode, description string) string {
	return fmt.Sprintf(`
resource "netbox_config_template" "test" {
  name          = %q
  template_code = %q
  description   = %q
}
`, name, templateCode, description)
}
