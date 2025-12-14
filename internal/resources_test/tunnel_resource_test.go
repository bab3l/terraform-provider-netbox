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

func TestTunnelResource(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelResource()

	if r == nil {

		t.Fatal("Expected non-nil Tunnel resource")

	}

}

func TestTunnelResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelResource()

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

	requiredAttrs := []string{"name", "status", "encapsulation"}

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

	optionalAttrs := []string{"group", "ipsec_profile", "tenant", "tunnel_id", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestTunnelResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tunnel"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestTunnelResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewTunnelResource()

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

func TestAccTunnelResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-tunnel")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "gre"),
				),
			},
		},
	})

}

func TestAccTunnelResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tunnel-full")

	description := "Test tunnel with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelResourceConfig_full(name, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tunnel.test", "id"),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "planned"),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "encapsulation", "wireguard"),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", description),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "tunnel_id", "12345"),
				),
			},
		},
	})

}

func TestAccTunnelResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tunnel-upd")

	updatedDescription := description2

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "active"),
				),
			},

			{

				Config: testAccTunnelResourceConfig_full(name, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tunnel.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "status", "planned"),

					resource.TestCheckResourceAttr("netbox_tunnel.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccTunnelResource_import(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tunnel-imp")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTunnelCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTunnelDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTunnelResourceConfig_basic(name),
			},

			{

				ResourceName: "netbox_tunnel.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccTunnelResourceConfig_basic(name string) string {

	return fmt.Sprintf(`



resource "netbox_tunnel" "test" {



  name          = %[1]q



  status        = "active"



  encapsulation = "gre"



}



`, name)

}

func testAccTunnelResourceConfig_full(name, description string) string {

	return fmt.Sprintf(`



resource "netbox_tunnel" "test" {



  name          = %[1]q



  status        = "planned"



  encapsulation = "wireguard"



  description   = %[2]q



  tunnel_id     = 12345



}



`, name, description)

}
