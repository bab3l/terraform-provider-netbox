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

func TestWirelessLANResource(t *testing.T) {
	t.Parallel()

	r := resources.NewWirelessLANResource()
	if r == nil {
		t.Fatal("Expected non-nil WirelessLAN resource")
	}
}

func TestWirelessLANResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewWirelessLANResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"ssid"}
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

	optionalAttrs := []string{"description", "group", "status", "vlan", "tenant", "auth_type", "auth_cipher", "auth_psk", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestWirelessLANResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewWirelessLANResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_wireless_lan"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestWirelessLANResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewWirelessLANResource().(*resources.WirelessLANResource)

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

func TestAccWirelessLANResource_basic(t *testing.T) {
	ssid := testutil.RandomName("tf-test-ssid")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANResourceConfig_basic(ssid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "status", "active"),
				),
			},
			{
				ResourceName:      "netbox_wireless_lan.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccWirelessLANResource_full(t *testing.T) {
	ssid := testutil.RandomName("tf-test-ssid-full")
	groupName := testutil.RandomName("tf-test-wlan-group")
	groupSlug := testutil.RandomSlug("tf-test-wlan-group")
	description := "Test wireless LAN with all fields"
	updatedDescription := "Updated wireless LAN description"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccWirelessLANResourceConfig_full(ssid, groupName, groupSlug, description, "active"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_wireless_lan.test", "id"),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "ssid", ssid),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", description),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "status", "active"),
				),
			},
			{
				Config: testAccWirelessLANResourceConfig_full(ssid, groupName, groupSlug, updatedDescription, "disabled"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_wireless_lan.test", "status", "disabled"),
				),
			},
		},
	})
}

func testAccWirelessLANResourceConfig_basic(ssid string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan" "test" {
  ssid = %q
}
`, ssid)
}

func testAccWirelessLANResourceConfig_full(ssid, groupName, groupSlug, description, status string) string {
	return fmt.Sprintf(`
resource "netbox_wireless_lan_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_wireless_lan" "test" {
  ssid        = %q
  group       = netbox_wireless_lan_group.test.id
  description = %q
  status      = %q
}
`, groupName, groupSlug, ssid, description, status)
}
