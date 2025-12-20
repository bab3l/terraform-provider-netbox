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

func TestWebhookDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewWebhookDataSource()

	if d == nil {

		t.Fatal("Expected non-nil webhook data source")

	}

}

func TestWebhookDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewWebhookDataSource()

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

	requiredAttrs := []string{"id", "name", "payload_url", "http_method", "http_content_type", "description", "additional_headers", "body_template", "ssl_verification", "ca_file_path", "tags"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestWebhookDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewWebhookDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_webhook"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestWebhookDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewWebhookDataSource().(*datasources.WebhookDataSource)

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

func TestAccWebhookDataSource_byID(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-webhook-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccWebhookDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_webhook.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_webhook.test", "id"),
				),
			},
		},
	})

}

func TestAccWebhookDataSource_byName(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-webhook-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccWebhookDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_webhook.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_webhook.test", "id"),
				),
			},
		},
	})

}

func testAccWebhookDataSourceByID(name string) string {

	return fmt.Sprintf(`







resource "netbox_webhook" "test" {



  name        = %[1]q







  payload_url = "https://example.com/webhook"



}







data "netbox_webhook" "test" {







  id = netbox_webhook.test.id



}







`, name)

}

func testAccWebhookDataSourceByName(name string) string {

	return fmt.Sprintf(`







resource "netbox_webhook" "test" {



  name        = %[1]q







  payload_url = "https://example.com/webhook"



}







data "netbox_webhook" "test" {



  name = netbox_webhook.test.name



}







`, name)

}
