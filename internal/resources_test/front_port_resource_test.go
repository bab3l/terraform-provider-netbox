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

func TestFrontPortResource(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortResource()

	if r == nil {

		t.Fatal("Expected non-nil FrontPort resource")

	}

}

func TestFrontPortResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"device", "name", "type", "rear_port"}

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

	optionalAttrs := []string{"label", "color", "rear_port_position", "description", "mark_connected", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestFrontPortResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_front_port"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestFrontPortResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewFrontPortResource().(*resources.FrontPortResource)

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

func TestAccFrontPortResource_basic(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	roleName := testutil.RandomName("tf-test-role")

	roleSlug := testutil.RandomSlug("tf-test-role")

	deviceName := testutil.RandomName("tf-test-device")

	rearPortName := testutil.RandomName("tf-test-rp")

	frontPortName := testutil.RandomName("tf-test-fp")

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

				Config: testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),

					resource.TestCheckResourceAttr("netbox_front_port.test", "name", frontPortName),

					resource.TestCheckResourceAttr("netbox_front_port.test", "type", "8p8c"),
				),
			},

			{

				ResourceName: "netbox_front_port.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccFrontPortResource_full(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	roleName := testutil.RandomName("tf-test-role")

	roleSlug := testutil.RandomSlug("tf-test-role")

	deviceName := testutil.RandomName("tf-test-device")

	rearPortName := testutil.RandomName("tf-test-rp")

	frontPortName := testutil.RandomName("tf-test-fp")

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

				Config: testAccFrontPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_front_port.test", "id"),

					resource.TestCheckResourceAttr("netbox_front_port.test", "name", frontPortName),

					resource.TestCheckResourceAttr("netbox_front_port.test", "type", "lc"),

					resource.TestCheckResourceAttr("netbox_front_port.test", "label", "Front Port Test"),

					resource.TestCheckResourceAttr("netbox_front_port.test", "color", "aa1409"),

					resource.TestCheckResourceAttr("netbox_front_port.test", "rear_port_position", "1"),

					resource.TestCheckResourceAttr("netbox_front_port.test", "description", "Test front port"),

					resource.TestCheckResourceAttr("netbox_front_port.test", "mark_connected", "true"),
				),
			},

			{

				ResourceName: "netbox_front_port.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccFrontPortResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName string) string {

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







  site        = netbox_site.test.id







  device_type = netbox_device_type.test.id







  role        = netbox_device_role.test.id







}















resource "netbox_rear_port" "test" {







  device    = netbox_device.test.id







  name      = %q







  type      = "8p8c"







  positions = 2







}















resource "netbox_front_port" "test" {







  device    = netbox_device.test.id







  name      = %q







  type      = "8p8c"







  rear_port = netbox_rear_port.test.id







}







`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName)

}

func testAccFrontPortResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName string) string {

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







  site        = netbox_site.test.id







  device_type = netbox_device_type.test.id







  role        = netbox_device_role.test.id







}















resource "netbox_rear_port" "test" {







  device    = netbox_device.test.id







  name      = %q







  type      = "lc"







  positions = 4







}















resource "netbox_front_port" "test" {







  device             = netbox_device.test.id







  name               = %q







  type               = "lc"







  rear_port          = netbox_rear_port.test.id







  rear_port_position = 1







  label              = "Front Port Test"







  color              = "aa1409"







  description        = "Test front port"







  mark_connected     = true







}







`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, rearPortName, frontPortName)

}
