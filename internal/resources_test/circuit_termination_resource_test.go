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

func TestCircuitTerminationResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTerminationResource()

	if r == nil {

		t.Fatal("Expected non-nil CircuitTermination resource")

	}

}

func TestCircuitTerminationResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTerminationResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"circuit", "term_side"}

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

	optionalAttrs := []string{"site", "provider_network", "port_speed", "upstream_speed", "xconnect_id", "pp_info", "description", "mark_connected", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestCircuitTerminationResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTerminationResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_circuit_termination"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestCircuitTerminationResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTerminationResource().(*resources.CircuitTerminationResource)

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

func TestAccCircuitTerminationResource_basic(t *testing.T) {

	providerName := testutil.RandomName("tf-test-provider")

	providerSlug := testutil.RandomSlug("tf-test-provider")

	circuitTypeName := testutil.RandomName("tf-test-ct")

	circuitTypeSlug := testutil.RandomSlug("tf-test-ct")

	circuitCID := testutil.RandomName("tf-test-circuit")

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)

	cleanup.RegisterCircuitCleanup(circuitCID)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},

			{

				ResourceName: "netbox_circuit_termination.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccCircuitTerminationResource_full(t *testing.T) {

	providerName := testutil.RandomName("tf-test-provider-full")

	providerSlug := testutil.RandomSlug("tf-test-provider-full")

	circuitTypeName := testutil.RandomName("tf-test-ct-full")

	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-full")

	circuitCID := testutil.RandomName("tf-test-circuit-full")

	siteName := testutil.RandomName("tf-test-site-full")

	siteSlug := testutil.RandomSlug("tf-test-site-full")

	description := "Test circuit termination with all fields"

	updatedDescription := "Updated circuit termination description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)

	cleanup.RegisterCircuitCleanup(circuitCID)

	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description, 1000000),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),

					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", description),

					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "1000000"),
				),
			},

			{

				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, updatedDescription, 10000000),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "10000000"),
				),
			},
		},
	})

}

func testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug string) string {

	return fmt.Sprintf(`































resource "netbox_provider" "test" {































  name = %q































  slug = %q































}































































resource "netbox_circuit_type" "test" {































  name = %q































  slug = %q































}































































resource "netbox_circuit" "test" {































  cid              = %q































  circuit_provider = netbox_provider.test.id































  type             = netbox_circuit_type.test.id































}































































resource "netbox_site" "test" {































  name   = %q































  slug   = %q































  status = "active"































}































































resource "netbox_circuit_termination" "test" {































  circuit   = netbox_circuit.test.id































  term_side = "A"































  site      = netbox_site.test.id































}































`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug)

}

func testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description string, portSpeed int) string {

	return fmt.Sprintf(`































resource "netbox_provider" "test" {































  name = %q































  slug = %q































}































































resource "netbox_circuit_type" "test" {































  name = %q































  slug = %q































}































































resource "netbox_circuit" "test" {































  cid              = %q































  circuit_provider = netbox_provider.test.id































  type             = netbox_circuit_type.test.id































}































































resource "netbox_site" "test" {































  name   = %q































  slug   = %q































  status = "active"































}































































resource "netbox_circuit_termination" "test" {































  circuit     = netbox_circuit.test.id































  term_side   = "A"































  site        = netbox_site.test.id































  port_speed  = %d































  description = %q































}































`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, portSpeed, description)

}

func TestAccCircuitTerminationResource_import(t *testing.T) {
	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	circuitTypeName := testutil.RandomName("tf-test-ct")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct")
	circuitCID := testutil.RandomName("tf-test-circuit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				ResourceName:      "netbox_circuit_termination.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
