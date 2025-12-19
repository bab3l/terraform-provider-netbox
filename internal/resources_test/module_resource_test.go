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

func TestModuleResource(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleResource()
	if r == nil {
		t.Fatal("Expected non-nil Module resource")
	}
}

func TestModuleResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}
	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"device", "module_bay", "module_type"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	computedAttrs := []string{"id", "status"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	optionalAttrs := []string{"serial", "asset_tag", "description", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestModuleResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}
	r.Metadata(context.Background(), metadataRequest, metadataResponse)
	expected := "netbox_module"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestModuleResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewModuleResource().(*resources.ModuleResource)
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

func TestAccModuleResource_basic(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	dtModel := testutil.RandomName("tf-test-dt")
	dtSlug := testutil.RandomSlug("tf-test-dt")
	roleName := testutil.RandomName("tf-test-role")
	roleSlug := testutil.RandomSlug("tf-test-role")
	deviceName := testutil.RandomName("tf-test-device")
	bayName := testutil.RandomName("tf-test-mbay")
	mtModel := testutil.RandomName("tf-test-mt")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "device"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "module_bay"),
					resource.TestCheckResourceAttrSet("netbox_module.test", "module_type"),
				),
			},
			{
				ResourceName:            "netbox_module.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"device", "module_bay", "module_type"},
			},
		},
	})
}

func TestAccModuleResource_full(t *testing.T) {
	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	mfgName := testutil.RandomName("tf-test-mfg-full")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")
	dtModel := testutil.RandomName("tf-test-dt-full")
	dtSlug := testutil.RandomSlug("tf-test-dt-full")
	roleName := testutil.RandomName("tf-test-role-full")
	roleSlug := testutil.RandomSlug("tf-test-role-full")
	deviceName := testutil.RandomName("tf-test-device-full")
	bayName := testutil.RandomName("tf-test-mbay-full")
	mtModel := testutil.RandomName("tf-test-mt-full")
	description := "Test module with all fields"
	updatedDescription := "Updated module description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterManufacturerCleanup(mfgSlug)
	cleanup.RegisterDeviceTypeCleanup(dtSlug)
	cleanup.RegisterDeviceRoleCleanup(roleSlug)
	cleanup.RegisterDeviceCleanup(deviceName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccModuleResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_module.test", "id"),
					resource.TestCheckResourceAttr("netbox_module.test", "serial", "SN123456"),
					resource.TestCheckResourceAttr("netbox_module.test", "description", description),
					resource.TestCheckResourceAttr("netbox_module.test", "status", "active"),
				),
			},
			{
				Config: testAccModuleResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, updatedDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "description", updatedDescription),
				),
			},
		},
	})
}

func testAccModuleResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel)
}

func testAccModuleResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name  = %q
  slug  = %q
  color = "aa1409"
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

resource "netbox_module" "test" {
  device      = netbox_device.test.id
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  status      = "active"
  serial      = "SN123456"
  description = %q
}
`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, description)
}

// TestAccConsistency_Module_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_Module_LiteralNames(t *testing.T) {
	t.Parallel()
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	dtModel := testutil.RandomName("device-type")
	dtSlug := testutil.RandomSlug("device-type")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")
	deviceName := testutil.RandomName("device")
	bayName := testutil.RandomName("bay")
	mtModel := testutil.RandomName("module-type")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccModuleConsistencyLiteralNamesConfig(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_module.test", "device", deviceName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccModuleConsistencyLiteralNamesConfig(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel),
			},
		},
	})
}

func testAccModuleConsistencyLiteralNamesConfig(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel string) string {
	return fmt.Sprintf(`

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
}

resource "netbox_device_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_device" "test" {
  name        = %q
  device_type = netbox_device_type.test.id
  role        = netbox_device_role.test.id
  site        = netbox_site.test.id
}

resource "netbox_module_bay" "test" {
  device = netbox_device.test.id
  name   = %q
}

resource "netbox_module_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
}

resource "netbox_module" "test" {
  # Use literal string name to mimic existing user state
  device      = %q
  module_bay  = netbox_module_bay.test.id
  module_type = netbox_module_type.test.id
  status      = "active"
  serial      = "SN123456"

  depends_on = [netbox_device.test]
}

`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, mtModel, deviceName)
}
