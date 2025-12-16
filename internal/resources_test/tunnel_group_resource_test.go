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

func TestTunnelGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil TunnelGroup resource")

	}

}

func TestTunnelGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelGroupResource()

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

	optionalAttrs := []string{"description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestTunnelGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelGroupResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tunnel_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestTunnelGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelGroupResource()

	// Type assert to access Configure method

	configurable, ok := r.(fwresource.ResourceWithConfigure)

	if !ok {

		t.Fatal("Resource does not implement ResourceWithConfigure")

	}

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)

	}

}

// Acceptance Tests

func TestAccTunnelGroupResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-tunnel-group")

	slug := testutil.RandomSlug("tf-test-tunnel-grp")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccTunnelGroupResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tunnel-group-full")

	slug := testutil.RandomSlug("tf-test-tg-full")

	description := "Test tunnel group with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tunnel_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", description),
				),
			},
		},
	})

}

func TestAccTunnelGroupResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tunnel-group-upd")

	slug := testutil.RandomSlug("tf-test-tg-upd")

	updatedDescription := description2

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),
				),
			},

			{

				Config: testAccTunnelGroupResourceConfig_full(name, slug, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tunnel_group.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccTunnelGroupResource_import(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tunnel-group-imp")

	slug := testutil.RandomSlug("tf-test-tg-imp")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelGroupCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelGroupResourceConfig_basic(name, slug),
			},

			{

				ResourceName: "netbox_tunnel_group.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccTunnelGroupResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`































resource "netbox_tunnel_group" "test" {































  name = %[1]q































  slug = %[2]q































}































`, name, slug)

}

func testAccTunnelGroupResourceConfig_full(name, slug, description string) string {

	return fmt.Sprintf(`































resource "netbox_tunnel_group" "test" {































  name        = %[1]q































  slug        = %[2]q































  description = %[3]q































}































`, name, slug, description)

}
