package datasources_test

import (
	"context"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestRackReservationDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewRackReservationDataSource()
	if d == nil {
		t.Fatal("Expected non-nil RackReservation data source")
	}
}

func TestRackReservationDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewRackReservationDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Lookup attributes
	lookupAttrs := []string{"id"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}

	// Output attributes
	outputAttrs := []string{"rack", "rack_id", "units", "user", "user_id", "tenant", "tenant_id", "description", "comments", "tags", "custom_fields"}
	for _, attr := range outputAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected output attribute %s to exist in schema", attr)
		}
	}
}

func TestRackReservationDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewRackReservationDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_rack_reservation"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestRackReservationDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewRackReservationDataSource().(*datasources.RackReservationDataSource)

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
}
