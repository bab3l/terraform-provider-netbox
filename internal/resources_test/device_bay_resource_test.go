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

func TestDeviceBayResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource()

	if r == nil {

		t.Fatal("Expected non-nil DeviceBay resource")

	}

}

func TestDeviceBayResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"device", "name"}

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

	optionalAttrs := []string{"label", "description", "installed_device", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestDeviceBayResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_device_bay"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestDeviceBayResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceBayResource().(*resources.DeviceBayResource)

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

func TestAccDeviceBayResource_basic(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	roleName := testutil.RandomName("tf-test-role")

	roleSlug := testutil.RandomSlug("tf-test-role")

	deviceName := testutil.RandomName("tf-test-device")

	bayName := testutil.RandomName("tf-test-bay")

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

				Config: testAccDeviceBayResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),
				),
			},

			{

				ResourceName: "netbox_device_bay.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccDeviceBayResource_full(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site-full")

	siteSlug := testutil.RandomSlug("tf-test-site-full")

	mfgName := testutil.RandomName("tf-test-mfg-full")

	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")

	dtModel := testutil.RandomName("tf-test-dt-full")

	dtSlug := testutil.RandomSlug("tf-test-dt-full")

	roleName := testutil.RandomName("tf-test-role-full")

	roleSlug := testutil.RandomSlug("tf-test-role-full")

	deviceName := testutil.RandomName("tf-test-device-full")

	bayName := testutil.RandomName("tf-test-bay-full")

	description := "Test device bay with all fields"

	updatedDescription := "Updated device bay description"

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

				Config: testAccDeviceBayResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_bay.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_bay.test", "name", bayName),

					resource.TestCheckResourceAttr("netbox_device_bay.test", "label", "Bay Label"),

					resource.TestCheckResourceAttr("netbox_device_bay.test", "description", description),
				),
			},

			{

				Config: testAccDeviceBayResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_device_bay.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccDeviceBayResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName string) string {

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



  manufacturer    = netbox_manufacturer.test.id



  model           = %q



  slug            = %q



  subdevice_role  = "parent"



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







resource "netbox_device_bay" "test" {



  device = netbox_device.test.id



  name   = %q



}



`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName)

}

func testAccDeviceBayResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description string) string {

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



  manufacturer   = netbox_manufacturer.test.id



  model          = %q



  slug           = %q



  subdevice_role = "parent"



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







resource "netbox_device_bay" "test" {



  device      = netbox_device.test.id



  name        = %q



  label       = "Bay Label"



  description = %q



}



`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, bayName, description)

}
