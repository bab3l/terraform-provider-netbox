package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestL2VPNResource(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNResource()
	if r == nil {
		t.Fatal("Expected non-nil L2VPN resource")
	}
}

func TestL2VPNResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required attributes
	requiredAttrs := []string{"name", "slug", "type"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	// Optional attributes
	optionalAttrs := []string{"identifier", "tenant", "description", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestL2VPNResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_l2vpn"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestL2VPNResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNResource().(*resources.L2VPNResource)

	// Test with nil provider data
	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Fatalf("Configure with nil provider data should not error: %+v", configureResponse.Diagnostics)
	}

	// Test with valid API client
	configureRequest = fwresource.ConfigureRequest{
		ProviderData: netbox.NewAPIClient(netbox.NewConfiguration()),
	}
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Fatalf("Configure with valid provider data should not error: %+v", configureResponse.Diagnostics)
	}
}

func TestAccL2VPNResource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
			// Test update
			{
				Config: testAccL2VPNResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Updated description"),
				),
			},
			// Test import
			{
				ResourceName:      "netbox_l2vpn.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccL2VPNResource_full(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vpls"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "identifier", "12345"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Test L2VPN"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "comments", "Test comments"),
				),
			},
		},
	})
}

func testAccL2VPNResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}
`, name, name)
}

func testAccL2VPNResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = "Updated description"
}
`, name+"-updated", name)
}

func testAccL2VPNResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vpls"
  identifier  = 12345
  description = "Test L2VPN"
  comments    = "Test comments"
}
`, name, name)
}

// TestAccConsistency_L2VPN_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
