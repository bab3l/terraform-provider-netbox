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

func TestL2VPNTerminationResource(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNTerminationResource()
	if r == nil {
		t.Fatal("Expected non-nil L2VPN Termination resource")
	}
}

func TestL2VPNTerminationResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNTerminationResource()
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
	requiredAttrs := []string{"l2vpn", "assigned_object_type", "assigned_object_id"}
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
}

func TestL2VPNTerminationResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNTerminationResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_l2vpn_termination"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestL2VPNTerminationResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewL2VPNTerminationResource().(*resources.L2VPNTerminationResource)

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

func TestAccL2VPNTerminationResource_basic(t *testing.T) {
	name := acctest.RandomWithPrefix("test-l2vpn-term")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNTerminationResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn_termination.test", "assigned_object_type", "ipam.vlan"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "l2vpn"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn_termination.test", "assigned_object_id"),
				),
			},
			// Test import
			{
				ResourceName:      "netbox_l2vpn_termination.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccL2VPNTerminationResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = "%s"
  slug = "%s"
  type = "vxlan"
}

resource "netbox_vlan" "test" {
  name    = "%s-vlan"
  vid     = 100
}

resource "netbox_l2vpn_termination" "test" {
  l2vpn                = netbox_l2vpn.test.id
  assigned_object_type = "ipam.vlan"
  assigned_object_id   = netbox_vlan.test.id
}
`, name, name, name)
}
