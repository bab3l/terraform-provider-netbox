package datasources_acceptance_tests

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

func TestCustomFieldChoiceSetDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomFieldChoiceSetDataSource()
	if d == nil {
		t.Fatal("Expected non-nil custom field choice set data source")
	}
}

func TestCustomFieldChoiceSetDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomFieldChoiceSetDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	lookupAttrs := []string{"id", "name"}
	for _, attr := range lookupAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected lookup attribute %s to exist in schema", attr)
		}
	}
}

func TestCustomFieldChoiceSetDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomFieldChoiceSetDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_custom_field_choice_set"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCustomFieldChoiceSetDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewCustomFieldChoiceSetDataSource()

	// Type assert to access Configure method
	configurable, ok := d.(fwdatasource.DataSourceWithConfigure)
	if !ok {
		t.Fatal("Data source does not implement DataSourceWithConfigure")
	}

	configureRequest := fwdatasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwdatasource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwdatasource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}

func TestAccCustomFieldChoiceSetDataSource_byID(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetDataSourceConfig_byID(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_field_choice_set.test", "name", name),
					resource.TestCheckResourceAttr("data.netbox_custom_field_choice_set.test", "extra_choices.#", "3"),
				),
			},
		},
	})
}

func TestAccCustomFieldChoiceSetDataSource_byName(t *testing.T) {

	t.Parallel()

	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetDataSourceConfig_byName(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_custom_field_choice_set.test", "name", name),
					resource.TestCheckResourceAttrSet("data.netbox_custom_field_choice_set.test", "id"),
				),
			},
		},
	})
}

func testAccCustomFieldChoiceSetDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field_choice_set" "test" {
  name = "%s"
  extra_choices = [
    { value = "opt1", label = "Option 1" },
    { value = "opt2", label = "Option 2" },
    { value = "opt3", label = "Option 3" },
  ]
}

data "netbox_custom_field_choice_set" "test" {
  id = netbox_custom_field_choice_set.test.id
}
`, name)
}

func testAccCustomFieldChoiceSetDataSourceConfig_byName(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field_choice_set" "test" {
  name = "%s"
  extra_choices = [
    { value = "opt1", label = "Option 1" },
    { value = "opt2", label = "Option 2" },
    { value = "opt3", label = "Option 3" },
  ]
}

data "netbox_custom_field_choice_set" "test" {
  name = netbox_custom_field_choice_set.test.name
}
`, name)
}
