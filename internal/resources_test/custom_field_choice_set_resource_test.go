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

func TestCustomFieldChoiceSetResource(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()
	if r == nil {
		t.Fatal("Expected non-nil custom field choice set resource")
	}
}

func TestCustomFieldChoiceSetResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	requiredAttrs := []string{"name", "extra_choices"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}
}

func TestCustomFieldChoiceSetResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_custom_field_choice_set"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCustomFieldChoiceSetResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomFieldChoiceSetResource()

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
	configureResponse = &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}

func TestAccCustomFieldChoiceSetResource_basic(t *testing.T) {
	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "extra_choices.#", "3"),
				),
			},
			{
				ResourceName:      "netbox_custom_field_choice_set.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCustomFieldChoiceSetResource_full(t *testing.T) {
	name := testutil.RandomName("cfcs")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "description", "Test choice set"),
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "order_alphabetically", "true"),
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "extra_choices.#", "3"),
				),
			},
		},
	})
}

func TestAccCustomFieldChoiceSetResource_update(t *testing.T) {
	name := testutil.RandomName("cfcs")
	updatedName := name + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomFieldChoiceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomFieldChoiceSetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", name),
				),
			},
			{
				Config: testAccCustomFieldChoiceSetResourceConfig_basic(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_field_choice_set.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccCustomFieldChoiceSetResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field_choice_set" "test" {
  name = "%s"
  extra_choices = [
    { value = "opt1", label = "Option 1" },
    { value = "opt2", label = "Option 2" },
    { value = "opt3", label = "Option 3" },
  ]
}
`, name)
}

func testAccCustomFieldChoiceSetResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_field_choice_set" "test" {
  name                 = "%s"
  description          = "Test choice set"
  order_alphabetically = true
  extra_choices = [
    { value = "critical", label = "Critical" },
    { value = "high",     label = "High" },
    { value = "low",      label = "Low" },
  ]
}
`, name)
}
