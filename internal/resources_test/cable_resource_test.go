package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCableResource(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource()
	if r == nil {
		t.Fatal("Expected non-nil cable resource")
	}
}

func TestCableResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required attributes - a_terminations and b_terminations are required
	requiredAttrs := []string{"a_terminations", "b_terminations"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	// Optional attributes
	optionalAttrs := []string{
		"type", "status", "tenant", "label", "color",
		"length", "length_unit", "description", "comments",
		"tags", "custom_fields",
	}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestCableResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_cable"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestCableResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewCableResource().(*resources.CableResource)

	// Test with nil provider data (should not error)
	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Fatalf("Configure with nil provider data should not error: %+v", configureResponse.Diagnostics)
	}
}

func TestAccCableResource_basic(t *testing.T) {
	siteName := testutil.RandomName("test-site-cable")
	siteSlug := testutil.GenerateSlug(siteName)
	deviceName := testutil.RandomName("test-device-cable")
	interfaceNameA := "eth2"
	interfaceNameB := "eth1"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCableResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cable.test", "status", "connected"),
					resource.TestCheckResourceAttr("netbox_cable.test", "type", "cat6"),
				),
			},
			{
				ResourceName:      "netbox_cable.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCableResourceConfig(siteName, siteSlug, deviceName, interfaceNameA, interfaceNameB string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = "Test Manufacturer"
  slug = "test-manufacturer"
}

resource "netbox_device_role" "test" {
  name = "Test Device Role"
  slug = "test-device-role"
}

resource "netbox_device_type" "test" {
  model = "Test Device Type"
  slug  = "test-device-type"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_device" "test_a" {
  name           = "%s-a"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_device" "test_b" {
  name           = "%s-b"
  device_type    = netbox_device_type.test.id
  role           = netbox_device_role.test.id
  site           = netbox_site.test.id
}

resource "netbox_interface" "test_a" {
  name      = %q
  device    = netbox_device.test_a.id
  type      = "1000base-t"
}

resource "netbox_interface" "test_b" {
  name      = %q
  device    = netbox_device.test_b.id
  type      = "1000base-t"
}

resource "netbox_cable" "test" {
  status = "connected"
  type   = "cat6"
  a_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test_a.id
    }
  ]
  b_terminations = [
    {
      object_type = "dcim.interface"
      object_id   = netbox_interface.test_b.id
    }
  ]
}
`, siteName, siteSlug, deviceName, deviceName, interfaceNameA, interfaceNameB)
}

// TestAccConsistency_Cable_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
