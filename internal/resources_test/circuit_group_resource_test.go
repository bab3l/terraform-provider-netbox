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

func TestCircuitGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil CircuitGroup resource")
	}
}

func TestCircuitGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

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

	requiredAttrs := []string{"name", "slug"}

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

	optionalAttrs := []string{"description", "tenant", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestCircuitGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_circuit_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCircuitGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCircuitGroupResource()

	// Type assert to access Configure method

	configurable, ok := r.(fwresource.ResourceWithConfigure)

	if !ok {

		t.Fatal("Resource does not implement ResourceWithConfigure")
	}

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)
	}
}

// Acceptance Tests

func TestAccCircuitGroupResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-circuit-group")

	slug := testutil.RandomSlug("tf-test-cg")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-circuit-group-full")

	slug := testutil.RandomSlug("tf-test-cg-full")

	description := "Test circuit group with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitGroupResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_circuit_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", description),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-circuit-group-upd")

	slug := testutil.RandomSlug("tf-test-cg-upd")

	updatedDescription := description2

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),
				),
			},

			{

				Config: testAccCircuitGroupResourceConfig_full(name, slug, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_circuit_group.test", "description", updatedDescription),
				),
			},
		},
	})
}

func TestAccCircuitGroupResource_import(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-circuit-group-imp")

	slug := testutil.RandomSlug("tf-test-cg-imp")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterCircuitGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckCircuitGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccCircuitGroupResourceConfig_basic(name, slug),
			},

			{

				ResourceName: "netbox_circuit_group.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func testAccCircuitGroupResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_circuit_group" "test" {
  name = %[1]q
  slug = %[2]q
}

`, name, slug)
}

func testAccCircuitGroupResourceConfig_full(name, slug, description string) string {

	return fmt.Sprintf(`

resource "netbox_circuit_group" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
}

`, name, slug, description)
}
