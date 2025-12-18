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

func TestRearPortTemplateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource()

	if r == nil {

		t.Fatal("Expected non-nil rear port template resource")
	}
}

func TestRearPortTemplateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")
	}

	// Check required attributes

	requiredAttrs := []string{"name", "type"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Check computed attributes

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	// Check optional attributes

	optionalAttrs := []string{"device_type", "module_type", "label", "color", "positions", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestRearPortTemplateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_rear_port_template"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestRearPortTemplateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRearPortTemplateResource().(*resources.RearPortTemplateResource)

	// Test with nil provider data (should not error)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct provider data

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with incorrect provider data

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")
	}
}

// Acceptance test configurations

func testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_rear_port_template" "test" {

  device_type = netbox_device_type.test.id
  name        = %q
  type        = %q
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType)
}

func testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, description string, positions int32) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_rear_port_template" "test" {

  device_type = netbox_device_type.test.id
  name        = %q
  type        = %q
  label       = %q
  color       = %q

  positions   = %d
  description = %q
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, positions, description)
}

func TestAccRearPortTemplateResource_basic(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	const name = "rear0"

	portType := "8p8c"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRearPortTemplateResourceBasic(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", portType),
				),
			},

			{

				ResourceName: "netbox_rear_port_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func TestAccRearPortTemplateResource_full(t *testing.T) {

	manufacturerName := testutil.RandomName("mfr")

	manufacturerSlug := testutil.RandomSlug("mfr")

	deviceTypeName := testutil.RandomName("dt")

	deviceTypeSlug := testutil.RandomSlug("dt")

	const name = "rear0"

	portType := "lc"

	label := "Rear Port 0"

	color := "aa1409"

	description := "Test rear port template"

	positions := int32(2)

	updatedName := "rear1"

	updatedLabel := "Rear Port 1"

	updatedDescription := "Updated rear port template"

	updatedPositions := int32(4)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, name, portType, label, color, description, positions),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rear_port_template.test", "id"),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "type", portType),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", label),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "color", color),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", fmt.Sprintf("%d", positions)),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", description),
				),
			},

			{

				Config: testAccRearPortTemplateResourceFull(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, updatedName, portType, updatedLabel, color, updatedDescription, updatedPositions),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "label", updatedLabel),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "positions", fmt.Sprintf("%d", updatedPositions)),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "description", updatedDescription),
				),
			},

			{

				ResourceName: "netbox_rear_port_template.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_RearPortTemplate(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("devicetype")

	deviceTypeSlug := testutil.RandomSlug("devicetype")

	const portName = "rear0"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "name", portName),

					resource.TestCheckResourceAttr("netbox_rear_port_template.test", "device_type", deviceTypeName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName),
			},
		},
	})
}

func testAccRearPortTemplateConsistencyConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName string) string {

	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_device_type" "test" {
  model = "%[3]s"
  slug = "%[4]s"
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_rear_port_template" "test" {

  device_type = netbox_device_type.test.model
  name = "%[5]s"
  type = "8p8c"

  positions = 1
}

`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, portName)
}
