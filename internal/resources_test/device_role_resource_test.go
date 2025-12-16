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

func TestDeviceRoleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceRoleResource()

	if r == nil {

		t.Fatal("Expected non-nil device role resource")

	}

}

func TestDeviceRoleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceRoleResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"name", "slug"}

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

	optionalAttrs := []string{"color", "vm_role", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestDeviceRoleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceRoleResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_device_role"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestDeviceRoleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewDeviceRoleResource().(*resources.DeviceRoleResource)

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

func TestAccDeviceRoleResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-device-role")

	slug := testutil.RandomSlug("tf-test-dr")

	// Register cleanup to ensure resource is deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckDeviceRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_role.test", "vm_role", "true"),
				),
			},
		},
	})

}

func TestAccDeviceRoleResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-device-role-full")

	slug := testutil.RandomSlug("tf-test-dr-full")

	description := "Test device role with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckDeviceRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceRoleResourceConfig_full(name, slug, description, "aa1409", false),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_device_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_device_role.test", "description", description),

					resource.TestCheckResourceAttr("netbox_device_role.test", "color", "aa1409"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "vm_role", "false"),
				),
			},
		},
	})

}

func TestAccDeviceRoleResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-device-role-update")

	slug := testutil.RandomSlug("tf-test-dr-upd")

	updatedName := testutil.RandomName("tf-test-device-role-updated")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterDeviceRoleCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckDeviceRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccDeviceRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", name),
				),
			},

			{

				Config: testAccDeviceRoleResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_device_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_device_role.test", "name", updatedName),
				),
			},
		},
	})

}

// testAccDeviceRoleResourceConfig_basic returns a basic test configuration.

func testAccDeviceRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`































terraform {































  required_providers {































    netbox = {































      source = "bab3l/netbox"































      version = ">= 0.1.0"































    }































  }































}































































provider "netbox" {}































































resource "netbox_device_role" "test" {































  name = %q































  slug = %q































}































`, name, slug)

}

// testAccDeviceRoleResourceConfig_full returns a test configuration with all fields.

func testAccDeviceRoleResourceConfig_full(name, slug, description, color string, vmRole bool) string {

	return fmt.Sprintf(`































terraform {































  required_providers {































    netbox = {































      source = "bab3l/netbox"































      version = ">= 0.1.0"































    }































  }































}































































provider "netbox" {}































































resource "netbox_device_role" "test" {































  name        = %q































  slug        = %q































  description = %q































  color       = %q































  vm_role     = %t































}































`, name, slug, description, color, vmRole)

}
