package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestConfigContextResource_Metadata(t *testing.T) {
	r := resources.NewConfigContextResource()

	req := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "netbox_config_context" {
		t.Errorf("Expected type name 'netbox_config_context', got '%s'", resp.TypeName)
	}
}

func TestConfigContextResource_Schema(t *testing.T) {
	r := resources.NewConfigContextResource()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	// Verify required attributes exist
	requiredAttrs := []string{"id", "name", "data"}
	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected attribute '%s' to exist in schema", attr)
		}
	}

	// Verify optional attributes exist
	optionalAttrs := []string{
		"description", "weight", "is_active",
		"regions", "site_groups", "sites", "locations",
		"device_types", "roles", "platforms",
		"cluster_types", "cluster_groups", "clusters",
		"tenant_groups", "tenants", "tags",
	}
	for _, attr := range optionalAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("Expected attribute '%s' to exist in schema", attr)
		}
	}
}

func TestConfigContextResource_SchemaDescription(t *testing.T) {
	r := resources.NewConfigContextResource()

	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Error("Expected schema to have a markdown description")
	}
}

func TestConfigContextResource_Configure(t *testing.T) {
	r := resources.NewConfigContextResource()

	// Verify the resource implements the configurable interface
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Skip("Resource does not implement ResourceWithConfigure")
	}

	// Test with nil provider data - should not error
	req := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("Configure with nil provider data should not error: %v", resp.Diagnostics)
	}
}

func TestAccConfigContextResource_basic(t *testing.T) {
	name := testutil.RandomName("tf-test-config-context")

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
				ResourceName:      "netbox_config_context.test",
				ImportState:       true,
				ImportStateVerify: true,
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
