package datasources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestContactAssignmentDataSource(t *testing.T) {

	t.Parallel()

	ds := datasources.NewContactAssignmentDataSource()

	if ds == nil {

		t.Fatal("Expected non-nil contact assignment data source")

	}

}

func TestContactAssignmentDataSourceSchema(t *testing.T) {

	t.Parallel()

	ds := datasources.NewContactAssignmentDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	ds.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check lookup attributes

	lookupAttrs := []string{"id"}

	for _, attr := range lookupAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected lookup attribute %s to exist in schema", attr)

		}

	}

	// Check computed attributes

	computedAttrs := []string{"object_type", "object_id", "contact_id", "contact_name", "role_id", "role_name", "priority", "priority_name", "tags"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

}

func TestContactAssignmentDataSourceMetadata(t *testing.T) {

	t.Parallel()

	ds := datasources.NewContactAssignmentDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	ds.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_contact_assignment"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestContactAssignmentDataSourceConfigure(t *testing.T) {

	t.Parallel()

	ds := datasources.NewContactAssignmentDataSource().(*datasources.ContactAssignmentDataSource)

	// Test with nil provider data (should not error)

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	ds.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct provider data type

	client := &netbox.APIClient{}

	configureRequest = fwdatasource.ConfigureRequest{

		ProviderData: client,
	}

	configureResponse = &fwdatasource.ConfigureResponse{}

	ds.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with wrong provider data type

	configureRequest = fwdatasource.ConfigureRequest{

		ProviderData: "wrong type",
	}

	configureResponse = &fwdatasource.ConfigureResponse{}

	ds.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with wrong provider data type")

	}

}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccContactAssignmentDataSource_basic(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")

	randomSlug := testutil.RandomSlug("test-ca-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			// Create resource and read via data source

			{

				Config: testAccContactAssignmentDataSourceConfig(randomName, randomSlug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrPair(

						"data.netbox_contact_assignment.test", "id",

						"netbox_contact_assignment.test", "id"),

					resource.TestCheckResourceAttrPair(

						"data.netbox_contact_assignment.test", "contact_id",

						"netbox_contact_assignment.test", "contact_id"),

					resource.TestCheckResourceAttr(

						"data.netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},
		},
	})

}

func testAccContactAssignmentDataSourceConfig(name, slug string) string {

	return fmt.Sprintf(`







resource "netbox_site" "test" {



  name   = "%s-site"



  slug   = "%s-site"



  status = "active"



}







resource "netbox_contact" "test" {



  name  = "%s-contact"







  email = "test@example.com"



}







resource "netbox_contact_role" "test" {



  name = "%s-role"



  slug = "%s-role"



}







resource "netbox_contact_assignment" "test" {



  object_type = "dcim.site"



  object_id   = netbox_site.test.id







  contact_id  = netbox_contact.test.id







  role_id     = netbox_contact_role.test.id



}







data "netbox_contact_assignment" "test" {







  id = netbox_contact_assignment.test.id



}







`, name, slug, name, name, slug)

}
