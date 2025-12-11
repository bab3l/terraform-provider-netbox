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

func TestCustomLinkResource(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomLinkResource()
	if r == nil {
		t.Fatal("Expected non-nil custom link resource")
	}
}

func TestCustomLinkResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomLinkResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	requiredAttrs := []string{"name", "object_types", "link_text", "link_url"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}
}

func TestCustomLinkResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomLinkResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_custom_link"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCustomLinkResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewCustomLinkResource()

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

func TestAccCustomLinkResource_basic(t *testing.T) {
	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_text", "View in External System"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "link_url", "https://example.com/device/{{ object.name }}"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "object_types.#", "1"),
				),
			},
			{
				ResourceName:      "netbox_custom_link.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccCustomLinkResource_full(t *testing.T) {
	name := testutil.RandomName("cl")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "enabled", "true"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "weight", "50"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "group_name", "External Links"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "button_class", "blue"),
					resource.TestCheckResourceAttr("netbox_custom_link.test", "new_window", "true"),
				),
			},
		},
	})
}

func TestAccCustomLinkResource_update(t *testing.T) {
	name := testutil.RandomName("cl")
	updatedName := name + "-updated"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckCustomLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomLinkResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", name),
				),
			},
			{
				Config: testAccCustomLinkResourceConfig_basic(updatedName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_custom_link.test", "name", updatedName),
				),
			},
		},
	})
}

func testAccCustomLinkResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device"]
  link_text    = "View in External System"
  link_url     = "https://example.com/device/{{ object.name }}"
}
`, name)
}

func testAccCustomLinkResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_custom_link" "test" {
  name         = "%s"
  object_types = ["dcim.device", "dcim.site"]
  enabled      = true
  link_text    = "View Details"
  link_url     = "https://example.com/{{ object.name }}"
  weight       = 50
  group_name   = "External Links"
  button_class = "blue"
  new_window   = true
}
`, name)
}
