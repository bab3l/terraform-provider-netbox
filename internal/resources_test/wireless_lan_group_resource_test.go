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

func TestWirelessLANGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewWirelessLANGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil WirelessLANGroup resource")

	}

}

func TestWirelessLANGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewWirelessLANGroupResource()

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

	optionalAttrs := []string{"description", "parent", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestWirelessLANGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewWirelessLANGroupResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_wireless_lan_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestWirelessLANGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewWirelessLANGroupResource().(*resources.WirelessLANGroupResource)

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

func TestAccWirelessLANGroupResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-wlan-group")

	slug := testutil.RandomSlug("tf-test-wlan-group")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccWirelessLANGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_wireless_lan_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_wireless_lan_group.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccWirelessLANGroupResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-wlan-group-full")

	slug := testutil.RandomSlug("tf-test-wlan-group-full")

	description := "Test wireless LAN group with all fields"

	updatedDescription := "Updated wireless LAN group description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccWirelessLANGroupResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_wireless_lan_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", description),
				),
			},

			{

				Config: testAccWirelessLANGroupResourceConfig_full(name, slug, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_wireless_lan_group.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccWirelessLANGroupResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`































resource "netbox_wireless_lan_group" "test" {































  name = %q































  slug = %q































}































`, name, slug)

}

func testAccWirelessLANGroupResourceConfig_full(name, slug, description string) string {

	return fmt.Sprintf(`































resource "netbox_wireless_lan_group" "test" {































  name        = %q































  slug        = %q































  description = %q































}































`, name, slug, description)

}
