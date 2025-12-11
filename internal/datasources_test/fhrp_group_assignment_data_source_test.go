package datasources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFHRPGroupAssignmentDataSource(t *testing.T) {
	t.Parallel()

	d := datasources.NewFHRPGroupAssignmentDataSource()
	if d == nil {
		t.Fatal("Expected non-nil FHRP Group Assignment data source")
	}
}

func TestFHRPGroupAssignmentDataSourceSchema(t *testing.T) {
	t.Parallel()

	d := datasources.NewFHRPGroupAssignmentDataSource()
	schemaRequest := fwdatasource.SchemaRequest{}
	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Required attributes
	requiredAttrs := []string{"id"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Computed attributes
	computedAttrs := []string{"group_id", "group_name", "interface_type", "interface_id", "priority"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}
}

func TestFHRPGroupAssignmentDataSourceMetadata(t *testing.T) {
	t.Parallel()

	d := datasources.NewFHRPGroupAssignmentDataSource()
	metadataRequest := fwdatasource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_fhrp_group_assignment"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestFHRPGroupAssignmentDataSourceConfigure(t *testing.T) {
	t.Parallel()

	d := datasources.NewFHRPGroupAssignmentDataSource().(*datasources.FHRPGroupAssignmentDataSource)

	// Test with nil provider data
	configureRequest := fwdatasource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct client type
	client := &netbox.APIClient{}
	configureRequest.ProviderData = client
	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}
}

func TestAccFHRPGroupAssignmentDataSource_byID(t *testing.T) {
	name := acctest.RandomWithPrefix("test-fhrp-assign-ds")

	resource.ParallelTest(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccFHRPGroupAssignmentDataSourceConfig_byID(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),
					resource.TestCheckResourceAttr("data.netbox_fhrp_group_assignment.test", "priority", "100"),
					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group_assignment.test", "group_id"),
					resource.TestCheckResourceAttrSet("data.netbox_fhrp_group_assignment.test", "interface_id"),
				),
			},
		},
	})
}

func testAccFHRPGroupAssignmentDataSourceConfig_byID(name string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%s-site"
  slug = "%s-site"
}

resource "netbox_manufacturer" "test" {
  name = "%s-mfr"
  slug = "%s-mfr"
}

resource "netbox_device_type" "test" {
  model           = "%s-dt"
  slug            = "%s-dt"
  manufacturer_id = netbox_manufacturer.test.id
}

resource "netbox_device_role" "test" {
  name  = "%s-role"
  slug  = "%s-role"
  color = "ff0000"
}

resource "netbox_device" "test" {
  name           = "%s-device"
  site_id        = netbox_site.test.id
  device_type_id = netbox_device_type.test.id
  role_id        = netbox_device_role.test.id
}

resource "netbox_interface" "test" {
  name      = "eth0"
  device_id = netbox_device.test.id
  type      = "virtual"
}

resource "netbox_fhrp_group" "test" {
  protocol = "vrrp2"
  group_id = 1
}

resource "netbox_fhrp_group_assignment" "test" {
  group_id       = netbox_fhrp_group.test.id
  interface_type = "dcim.interface"
  interface_id   = netbox_interface.test.id
  priority       = 100
}

data "netbox_fhrp_group_assignment" "test" {
  id = netbox_fhrp_group_assignment.test.id
}
`, name, name, name, name, name, name, name, name, name)
}
