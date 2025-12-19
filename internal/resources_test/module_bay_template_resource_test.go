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

func TestModuleBayTemplateResource(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleBayTemplateResource()
	if r == nil {
		t.Fatal("Expected non-nil ModuleBayTemplate resource")
	}
}

func TestModuleBayTemplateResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleBayTemplateResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}
	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name"}
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

	// Either device_type or module_type is required
	optionalAttrs := []string{"device_type", "module_type", "label", "position", "description"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestModuleBayTemplateResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleBayTemplateResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}
	r.Metadata(context.Background(), metadataRequest, metadataResponse)
	expected := "netbox_module_bay_template"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestModuleBayTemplateResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleBayTemplateResource().(*resources.ModuleBayTemplateResource)
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

func TestAccModuleBayTemplateResource_basic(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	templateName := testutil.RandomName("tf-test-mbt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckModuleBayTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", templateName),
				),
			},
			// ImportState test
			{
				ResourceName:            "netbox_module_bay_template.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device_type"},
			},
		},
	})
}

func TestAccModuleBayTemplateResource_full(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	templateName := testutil.RandomName("tf-test-mbt")
	label := "Bay 1"
	position := "Front"
	description := "Test module bay template"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckModuleBayTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, label, position, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", templateName),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "label", label),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "position", position),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "description", description),
				),
			},
		},
	})
}

func TestAccModuleBayTemplateResource_update(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	templateName := testutil.RandomName("tf-test-mbt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckModuleBayTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, "", "", description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module_bay_template.test", "id"),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "description", description1),
				),
			},
			{
				Config: testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, "", "", description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "description", description2),
				),
			},
		},
	})
}

func testAccModuleBayTemplateResourceConfig_basic(mfgName, mfgSlug, dtModel, dtSlug, templateName string) string {
	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_module_bay_template" "test" {
  name        = %[5]q
  device_type = netbox_device_type.test.id
}
`, mfgName, mfgSlug, dtModel, dtSlug, templateName)
}

func testAccModuleBayTemplateResourceConfig_full(mfgName, mfgSlug, dtModel, dtSlug, templateName, label, position, description string) string {
	labelAttr := ""
	if label != "" {
		labelAttr = fmt.Sprintf(`label       = %q`, label)
	}
	positionAttr := ""
	if position != "" {
		positionAttr = fmt.Sprintf(`position    = %q`, position)
	}
	descAttr := ""
	if description != "" {
		descAttr = fmt.Sprintf(`description = %q`, description)
	}

	return fmt.Sprintf(`
provider "netbox" {}

resource "netbox_manufacturer" "test" {
  name = %[1]q
  slug = %[2]q
}

resource "netbox_device_type" "test" {
  model        = %[3]q
  slug         = %[4]q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_module_bay_template" "test" {
  name        = %[5]q
  device_type = netbox_device_type.test.id
  %[6]s
  %[7]s
  %[8]s
}
`, mfgName, mfgSlug, dtModel, dtSlug, templateName, labelAttr, positionAttr, descAttr)
}

// TestAccConsistency_ModuleBayTemplate_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_ModuleBayTemplate_LiteralNames(t *testing.T) {
	t.Parallel()
	manufacturerName := testutil.RandomName("manufacturer")
	manufacturerSlug := testutil.RandomSlug("manufacturer")
	deviceTypeName := testutil.RandomName("device-type")
	deviceTypeSlug := testutil.RandomSlug("device-type")
	resourceName := testutil.RandomName("module_bay")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "name", resourceName),
					resource.TestCheckResourceAttr("netbox_module_bay_template.test", "device_type", deviceTypeSlug),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccModuleBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName),
			},
		},
	})
}

func testAccModuleBayTemplateConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, resourceName string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.id
}

resource "netbox_module_bay_template" "test" {
  # Use literal string slug to mimic existing user state
  device_type = %q
  name = %q

  depends_on = [netbox_device_type.test]
}
`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceTypeSlug, resourceName)
}
