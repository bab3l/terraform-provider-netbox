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

func TestDeviceTypeResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceTypeResource()

	if r == nil {

		t.Fatal("Expected non-nil device type resource")

	}

}

func TestDeviceTypeResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceTypeResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"manufacturer", "model", "slug"}

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

	optionalAttrs := []string{"default_platform", "part_number", "u_height", "exclude_from_utilization", "is_full_depth", "subdevice_role", "airflow", "weight", "weight_unit", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestDeviceTypeResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceTypeResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_device_type"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestDeviceTypeResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceTypeResource().(*resources.DeviceTypeResource)

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

func TestAccDeviceTypeResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	model := testutil.RandomName("tf-test-device-type")

	slug := testutil.RandomSlug("tf-test-dt")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceTypeCleanup(slug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_type.test", "manufacturer", manufacturerSlug),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "1"),
				),
			},
		},
	})

}

func TestAccDeviceTypeResource_full(t *testing.T) {

	// Generate unique names

	model := testutil.RandomName("tf-test-device-type-full")

	slug := testutil.RandomSlug("tf-test-dt-full")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceTypeCleanup(slug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_type.test", "manufacturer", manufacturerSlug),

					resource.TestCheckResourceAttr("netbox_device_type.test", "part_number", "TEST-PART-001"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "2"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "exclude_from_utilization", "false"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "is_full_depth", "true"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "airflow", "front-to-rear"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Test device type with full options"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "comments", "Test comments for device type"),
				),
			},
		},
	})

}

func TestAccDeviceTypeResource_update(t *testing.T) {

	// Generate unique names

	model := testutil.RandomName("tf-test-device-type-update")

	slug := testutil.RandomSlug("tf-test-dt-upd")

	manufacturerName := testutil.RandomName("tf-test-manufacturer")

	manufacturerSlug := testutil.RandomSlug("tf-test-mfr")

	updatedModel := testutil.RandomName("tf-test-device-type-updated")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceTypeCleanup(slug)

	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckDeviceTypeDestroy,

			testutil.CheckManufacturerDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
				),
			},

			{

				Config: testAccDeviceTypeResourceConfig_updated(updatedModel, slug, manufacturerName, manufacturerSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "model", updatedModel),

					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_type.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_device_type.test", "u_height", "4"),
				),
			},
		},
	})

}

// Test config helper functions

func testAccDeviceTypeResourceConfig_basic(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`































resource "netbox_manufacturer" "test" {































  name = %q































  slug = %q































}































































resource "netbox_device_type" "test" {































  manufacturer = netbox_manufacturer.test.slug































  model        = %q































  slug         = %q































}































`, manufacturerName, manufacturerSlug, model, slug)

}

func testAccDeviceTypeResourceConfig_full(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`































resource "netbox_manufacturer" "test" {































  name = %q































  slug = %q































}































































resource "netbox_device_type" "test" {































  manufacturer             = netbox_manufacturer.test.slug































  model                    = %q































  slug                     = %q































  part_number              = "TEST-PART-001"































  u_height                 = 2































  exclude_from_utilization = false































  is_full_depth            = true































  airflow                  = "front-to-rear"































  description              = "Test device type with full options"































  comments                 = "Test comments for device type"































}































`, manufacturerName, manufacturerSlug, model, slug)

}

func testAccDeviceTypeResourceConfig_updated(model, slug, manufacturerName, manufacturerSlug string) string {

	return fmt.Sprintf(`































resource "netbox_manufacturer" "test" {































  name = %q































  slug = %q































}































































resource "netbox_device_type" "test" {































  manufacturer = netbox_manufacturer.test.slug































  model        = %q































  slug         = %q































  u_height     = 4































  description  = "Updated description"































}































`, manufacturerName, manufacturerSlug, model, slug)

}

func TestAccDeviceTypeResource_import(t *testing.T) {
	model := testutil.RandomName("tf-test-dt-import")
	slug := testutil.RandomSlug("tf-test-dt-import")
	manufacturerName := testutil.RandomName("tf-test-mfr-import")
	manufacturerSlug := testutil.RandomSlug("tf-test-mfr-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterDeviceTypeCleanup(slug)
	cleanup.RegisterManufacturerCleanup(manufacturerSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckDeviceTypeDestroy,
			testutil.CheckManufacturerDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceTypeResourceConfig_import(model, slug, manufacturerName, manufacturerSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_device_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_device_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_device_type.test", "slug", slug),
				),
			},
			{
				ResourceName:            "netbox_device_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})
}

func testAccDeviceTypeResourceConfig_import(model, slug, manufacturerName, manufacturerSlug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_device_type" "test" {
  model        = %q
  slug         = %q
  manufacturer = netbox_manufacturer.test.slug
}
`, manufacturerName, manufacturerSlug, model, slug)
}
