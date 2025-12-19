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

func TestPowerOutletResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource()

	if r == nil {

		t.Fatal("Expected non-nil PowerOutlet resource")

	}

}

func TestPowerOutletResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource()

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

	optionalAttrs := []string{"label", "type", "power_port", "feed_leg", "description", "mark_connected", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestPowerOutletResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_power_outlet"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestPowerOutletResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPowerOutletResource().(*resources.PowerOutletResource)

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

func TestAccPowerOutletResource_basic(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	dtModel := testutil.RandomName("tf-test-dt")

	dtSlug := testutil.RandomSlug("tf-test-dt")

	roleName := testutil.RandomName("tf-test-role")

	roleSlug := testutil.RandomSlug("tf-test-role")

	deviceName := testutil.RandomName("tf-test-device")

	powerOutletName := testutil.RandomName("tf-test-po")

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

				Config: testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),
				),
			},

			{

				ResourceName: "netbox_power_outlet.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"device"},
			},
		},
	})

}

func TestAccPowerOutletResource_full(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site-full")

	siteSlug := testutil.RandomSlug("tf-test-site-full")

	mfgName := testutil.RandomName("tf-test-mfg-full")

	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")

	dtModel := testutil.RandomName("tf-test-dt-full")

	dtSlug := testutil.RandomSlug("tf-test-dt-full")

	roleName := testutil.RandomName("tf-test-role-full")

	roleSlug := testutil.RandomSlug("tf-test-role-full")

	deviceName := testutil.RandomName("tf-test-device-full")

	powerOutletName := testutil.RandomName("tf-test-po-full")

	description := "Test power outlet with all fields"

	updatedDescription := "Updated power outlet description"

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

				Config: testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_power_outlet.test", "id"),

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", powerOutletName),

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "description", description),

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "type", "iec-60320-c13"),
				),
			},

			{

				Config: testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccPowerOutletResourceConfig_basic(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName string) string {

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















resource "netbox_power_outlet" "test" {







  device = netbox_device.test.id







  name   = %q







}















`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName)

}

func testAccPowerOutletResourceConfig_full(siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, description string) string {

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















resource "netbox_power_outlet" "test" {







  device      = netbox_device.test.id







  name        = %q







  type        = "iec-60320-c13"







  description = %q







}















`, siteName, siteSlug, mfgName, mfgSlug, dtModel, dtSlug, roleName, roleSlug, deviceName, powerOutletName, description)

}

func TestAccConsistency_PowerOutlet(t *testing.T) {

	t.Parallel()

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	deviceRoleName := testutil.RandomName("device-role")

	deviceRoleSlug := testutil.RandomSlug("device-role")

	deviceName := testutil.RandomName("device")

	powerOutletName := testutil.RandomName("power-outlet")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccPowerOutletConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "device", deviceName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccPowerOutletConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName),
			},
		},
	})

}

func testAccPowerOutletConsistencyConfig(siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName string) string {

	return fmt.Sprintf(`















resource "netbox_site" "test" {







  name = "%[1]s"







  slug = "%[2]s"







}















resource "netbox_manufacturer" "test" {







  name = "%[3]s"







  slug = "%[4]s"







}















resource "netbox_device_type" "test" {







  model = "%[5]s"







  slug = "%[6]s"







  manufacturer = netbox_manufacturer.test.id







}















resource "netbox_device_role" "test" {







  name = "%[7]s"







  slug = "%[8]s"







}















resource "netbox_device" "test" {







  name = "%[9]s"















  device_type = netbox_device_type.test.id







  role = netbox_device_role.test.id







  site = netbox_site.test.id







}















resource "netbox_power_outlet" "test" {







  device = netbox_device.test.name







  name = "%[10]s"







  type = "iec-60320-c13"







}















`, siteName, siteSlug, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, deviceRoleName, deviceRoleSlug, deviceName, powerOutletName)

}

// TestAccConsistency_PowerOutlet_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_PowerOutlet_LiteralNames(t *testing.T) {

	t.Parallel()

	manufacturerName := testutil.RandomName("manufacturer")

	manufacturerSlug := testutil.RandomSlug("manufacturer")

	deviceTypeName := testutil.RandomName("device-type")

	deviceTypeSlug := testutil.RandomSlug("device-type")

	roleName := testutil.RandomName("role")

	roleSlug := testutil.RandomSlug("role")

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	deviceName := testutil.RandomName("device")

	resourceName := testutil.RandomName("power_outlet")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccPowerOutletConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "name", resourceName),

					resource.TestCheckResourceAttr("netbox_power_outlet.test", "device", deviceName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccPowerOutletConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName),
			},
		},
	})

}

func testAccPowerOutletConsistencyLiteralNamesConfig(manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, resourceName string) string {

	return fmt.Sprintf(`



resource "netbox_manufacturer" "test" {

  name = %q

  slug = %q

}



resource "netbox_device_type" "test" {

  model          = %q

  slug           = %q

  manufacturer   = netbox_manufacturer.test.id

  subdevice_role = "parent"  # Enable device bays

}



resource "netbox_site" "test" {

  name = %q

  slug = %q

}



resource "netbox_device_role" "test" {

  name = %q

  slug = %q

  color = "ff0000"

}



resource "netbox_device" "test" {

  name        = %q

  device_type = netbox_device_type.test.id

  site        = netbox_site.test.id

  role        = netbox_device_role.test.id

  status      = "active"

}



resource "netbox_power_outlet" "test" {

  # Use literal string name to mimic existing user state

  device = %q

  name = %q

  type = "iec-60320-c13"



  depends_on = [netbox_device.test]

}



`, manufacturerName, manufacturerSlug, deviceTypeName, deviceTypeSlug, roleName, roleSlug, siteName, siteSlug, deviceName, deviceName, resourceName)

}
