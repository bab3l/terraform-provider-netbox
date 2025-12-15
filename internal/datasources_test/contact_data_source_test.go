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

func TestContactDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewContactDataSource()

	if d == nil {

		t.Fatal("Expected non-nil contact data source")

	}

}

func TestContactDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewContactDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check that key attributes exist

	requiredAttrs := []string{"id", "name", "group", "title", "phone", "email", "address", "link", "description", "comments", "tags"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestContactDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewContactDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_contact"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestContactDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewContactDataSource().(*datasources.ContactDataSource)

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with invalid provider data")

	}

}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccContactDataSource_byID(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccContactDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
				),
			},
		},
	})

}

func TestAccContactDataSource_byName(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccContactDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
				),
			},
		},
	})

}

func TestAccContactDataSource_byEmail(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-ds")

	email := fmt.Sprintf("%s@example.com", randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccContactDataSourceByEmail(randomName, email),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_contact.test", "name", randomName),

					resource.TestCheckResourceAttr("data.netbox_contact.test", "email", email),

					resource.TestCheckResourceAttrSet("data.netbox_contact.test", "id"),
				),
			},
		},
	})

}

func testAccContactDataSourceByID(name string) string {

	return fmt.Sprintf(`







resource "netbox_contact" "test" {







  name = %[1]q







}















data "netbox_contact" "test" {







  id = netbox_contact.test.id







}







`, name)

}

func testAccContactDataSourceByName(name string) string {

	return fmt.Sprintf(`







resource "netbox_contact" "test" {







  name = %[1]q







}















data "netbox_contact" "test" {







  name = netbox_contact.test.name







}







`, name)

}

func testAccContactDataSourceByEmail(name, email string) string {

	return fmt.Sprintf(`







resource "netbox_contact" "test" {







  name  = %[1]q







  email = %[2]q







}















data "netbox_contact" "test" {







  email = netbox_contact.test.email







}







`, name, email)

}
