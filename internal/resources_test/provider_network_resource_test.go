package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestProviderNetworkResource(t *testing.T) {
	t.Parallel()

	r := resources.NewProviderNetworkResource()
	if r == nil {
		t.Fatal("Expected non-nil ProviderNetwork resource")
	}
}

func TestProviderNetworkResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewProviderNetworkResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"circuit_provider", "name"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	optionalAttrs := []string{"service_id", "description", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestProviderNetworkResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewProviderNetworkResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_provider_network"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestProviderNetworkResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewProviderNetworkResource().(*resources.ProviderNetworkResource)

	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	configureRequest.ProviderData = invalidProviderData
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with incorrect provider data")
	}
}

func TestAccProviderNetworkResource_basic(t *testing.T) {
	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	networkName := testutil.RandomName("tf-test-network")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
				),
			},
			{
				ResourceName:      "netbox_provider_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccProviderNetworkResource_full(t *testing.T) {
	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	networkName := testutil.RandomName("tf-test-network-full")
	serviceID := testutil.RandomName("svc")
	description := "Test provider network with all fields"
	updatedDescription := "Updated provider network description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_provider_network.test", "id"),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "name", networkName),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "service_id", serviceID),
					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", description),
				),
			},
			{
				Config: testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_provider_network.test", "description", updatedDescription),
				),
			},
		},
	})
}

func testAccProviderNetworkResourceConfig_basic(providerName, providerSlug, networkName string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_circuit_provider.test.id
  name             = %q
}
`, providerName, providerSlug, networkName)
}

func testAccProviderNetworkResourceConfig_full(providerName, providerSlug, networkName, serviceID, description string) string {
	return fmt.Sprintf(`
resource "netbox_circuit_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_provider_network" "test" {
  circuit_provider = netbox_circuit_provider.test.id
  name             = %q
  service_id       = %q
  description      = %q
}
`, providerName, providerSlug, networkName, serviceID, description)
}
