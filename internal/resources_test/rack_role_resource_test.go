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

func TestRackRoleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource()

	if r == nil {

		t.Fatal("Expected non-nil rack role resource")
	}
}

func TestRackRoleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource()

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

	requiredAttrs := []string{"name", "slug"}

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

	// Optional attributes

	optionalAttrs := []string{"color", "description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestRackRoleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_rack_role"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestRackRoleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRackRoleResource().(*resources.RackRoleResource)

	// Test with nil provider data (should not error)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct provider data type

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with incorrect provider data type

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")
	}
}

func TestAccRackRoleResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	rackRoleName := testutil.RandomName("tf-test-rack-role")

	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),
				),
			},
		},
	})
}

func TestAccRackRoleResource_full(t *testing.T) {

	// Generate unique names

	rackRoleName := testutil.RandomName("tf-test-rack-role-full")

	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-f")

	description := "Test rack role with all fields"

	color := "ff5722"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, description, color),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "slug", rackRoleSlug),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", description),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", color),
				),
			},
		},
	})
}

func TestAccRackRoleResource_update(t *testing.T) {

	// Generate unique names

	rackRoleName := testutil.RandomName("tf-test-rack-role-upd")

	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-u")

	updatedDescription := description2

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),
				),
			},

			{

				Config: testAccRackRoleResourceConfig_full(rackRoleName, rackRoleSlug, updatedDescription, "00bcd4"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack_role.test", "name", rackRoleName),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_rack_role.test", "color", "00bcd4"),
				),
			},
		},
	})
}

func TestAccRackRoleResource_import(t *testing.T) {

	// Generate unique names

	rackRoleName := testutil.RandomName("tf-test-rack-role-imp")

	rackRoleSlug := testutil.RandomSlug("tf-test-rack-role-i")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterRackRoleCleanup(rackRoleSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackRoleDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackRoleResourceConfig_basic(rackRoleName, rackRoleSlug),
			},

			{

				ResourceName: "netbox_rack_role.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func testAccRackRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_rack_role" "test" {
  name = %[1]q
  slug = %[2]q
}

`, name, slug)
}

func testAccRackRoleResourceConfig_full(name, slug, description, color string) string {

	return fmt.Sprintf(`

resource "netbox_rack_role" "test" {
  name        = %[1]q
  slug        = %[2]q
  description = %[3]q
  color       = %[4]q
}

`, name, slug, description, color)
}
