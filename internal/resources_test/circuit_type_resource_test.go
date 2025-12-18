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

func TestCircuitTypeResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTypeResource()

	if r == nil {

		t.Fatal("Expected non-nil CircuitType resource")
	}
}

func TestCircuitTypeResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTypeResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "slug"}

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

	optionalAttrs := []string{"description", "color", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestCircuitTypeResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTypeResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_circuit_type"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCircuitTypeResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitTypeResource().(*resources.CircuitTypeResource)

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

func TestAccCircuitTypeResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-circuit-type")

	slug := testutil.RandomSlug("tf-test-circuit-type")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitTypeResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-circuit-type-full")

	slug := testutil.RandomSlug("tf-test-circuit-type-full")

	description := "Test circuit type with all fields"

	const color = "aa1409"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTypeResourceConfig_full(name, slug, description, color),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "description", description),

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "color", color),
				),
			},
		},
	})
}

func TestAccCircuitTypeResource_update(t *testing.T) {

	name := testutil.RandomName("tf-test-circuit-type-update")

	slug := testutil.RandomSlug("tf-test-circuit-type-update")

	updatedName := testutil.RandomName("tf-test-circuit-type-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitTypeDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitTypeResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
				),
			},

			{

				Config: testAccCircuitTypeResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccCircuitTypeResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

`, name, slug)
}

func testAccCircuitTypeResourceConfig_full(name, slug, description, color string) string {

	return fmt.Sprintf(`

resource "netbox_circuit_type" "test" {
  name        = %q
  slug        = %q
  description = %q
  color       = %q
}

`, name, slug, description, color)
}

func TestAccCircuitTypeResource_import(t *testing.T) {
	name := testutil.RandomName("tf-test-circuit-type")
	slug := testutil.RandomSlug("tf-test-circuit-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCircuitTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_circuit_type.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_circuit_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
