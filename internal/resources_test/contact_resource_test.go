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

func TestContactResource(t *testing.T) {

	t.Parallel()

	r := resources.NewContactResource()

	if r == nil {

		t.Fatal("Expected non-nil contact resource")

	}

}

func TestContactResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewContactResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check required attributes

	requiredAttrs := []string{"name"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	// Check computed attributes

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	// Check optional attributes

	optionalAttrs := []string{"group", "title", "phone", "email", "address", "link", "description", "comments", "tags"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestContactResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewContactResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_contact"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestContactResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewContactResource().(*resources.ContactResource)

	// Test with nil provider data (should not error)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct provider data type

	client := &netbox.APIClient{}

	configureRequest = fwresource.ConfigureRequest{

		ProviderData: client,
	}

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with wrong provider data type

	configureRequest = fwresource.ConfigureRequest{

		ProviderData: "wrong type",
	}

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with wrong provider data type")

	}

}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccContactResource_basic(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			// Create and Read testing

			{

				Config: testAccContactResource(randomName, "john.doe@example.com", "+1-555-0100"),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttr("netbox_contact.test", "email", "john.doe@example.com"),

					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0100"),

					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
				),
			},

			// ImportState testing

			{

				ResourceName: "netbox_contact.test",

				ImportState: true,

				ImportStateVerify: true,
			},

			// Update and Read testing

			{

				Config: testAccContactResource(randomName, "jane.doe@example.com", "+1-555-0200"),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttr("netbox_contact.test", "email", "jane.doe@example.com"),

					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0200"),
				),
			},

			// Delete testing automatically occurs in TestCase

		},
	})

}

func TestAccContactResource_full(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-full")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			// Create with all fields

			{

				Config: testAccContactResourceFull(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttr("netbox_contact.test", "title", "Network Engineer"),

					resource.TestCheckResourceAttr("netbox_contact.test", "phone", "+1-555-0100"),

					resource.TestCheckResourceAttr("netbox_contact.test", "email", "engineer@example.com"),

					resource.TestCheckResourceAttr("netbox_contact.test", "address", "123 Main Street, City, Country"),

					resource.TestCheckResourceAttr("netbox_contact.test", "link", "https://example.com/profile"),

					resource.TestCheckResourceAttr("netbox_contact.test", "description", "Test contact description"),

					resource.TestCheckResourceAttr("netbox_contact.test", "comments", "Test contact comments"),

					resource.TestCheckResourceAttrSet("netbox_contact.test", "id"),
				),
			},
		},
	})

}

func testAccContactResource(name, email, phone string) string {

	return fmt.Sprintf(`







resource "netbox_contact" "test" {







  name  = %[1]q







  email = %[2]q







  phone = %[3]q







}







`, name, email, phone)

}

func testAccContactResourceFull(name string) string {

	return fmt.Sprintf(`







resource "netbox_contact" "test" {







  name        = %[1]q







  title       = "Network Engineer"







  phone       = "+1-555-0100"







  email       = "engineer@example.com"







  address     = "123 Main Street, City, Country"







  link        = "https://example.com/profile"







  description = "Test contact description"







  comments    = "Test contact comments"







}







`, name)

}
