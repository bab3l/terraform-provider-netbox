package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestFHRPGroupAssignmentResource(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource()

	if r == nil {

		t.Fatal("Expected non-nil FHRP Group Assignment resource")

	}

}

func TestFHRPGroupAssignmentResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Required attributes

	requiredAttrs := []string{"group_id", "interface_type", "interface_id", "priority"}

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

}

func TestFHRPGroupAssignmentResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_fhrp_group_assignment"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestFHRPGroupAssignmentResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewFHRPGroupAssignmentResource().(*resources.FHRPGroupAssignmentResource)

	// Test with nil provider data

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Fatalf("Configure with nil provider data should not error: %+v", configureResponse.Diagnostics)

	}

	// Test with valid API client

	configureRequest = fwresource.ConfigureRequest{

		ProviderData: netbox.NewAPIClient(netbox.NewConfiguration()),
	}

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Fatalf("Configure with valid provider data should not error: %+v", configureResponse.Diagnostics)

	}

}

func TestAccFHRPGroupAssignmentResource_basic(t *testing.T) {

	name := acctest.RandomWithPrefix("test-fhrp-assign")

	resource.ParallelTest(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccFHRPGroupAssignmentResourceConfig_basic(name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "interface_type", "dcim.interface"),

					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "100"),

					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "id"),

					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "group_id"),

					resource.TestCheckResourceAttrSet("netbox_fhrp_group_assignment.test", "interface_id"),
				),
			},

			// Test update

			{

				Config: testAccFHRPGroupAssignmentResourceConfig_updated(name),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_fhrp_group_assignment.test", "priority", "200"),
				),
			},

			// Test import

			{

				ResourceName: "netbox_fhrp_group_assignment.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"group_id", "interface_id"},
			},
		},
	})

}

func testAccFHRPGroupAssignmentResourceConfig_basic(name string) string {

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

`, name, name, name, name, name, name, name, name, name)

}

func testAccFHRPGroupAssignmentResourceConfig_updated(name string) string {

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

  priority       = 200

}

`, name, name, name, name, name, name, name, name, name)

}
