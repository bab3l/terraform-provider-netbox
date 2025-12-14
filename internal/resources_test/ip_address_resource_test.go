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

func TestIPAddressResource(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource()

	if r == nil {

		t.Fatal("Expected non-nil IP Address resource")

	}

}

func TestIPAddressResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"address"}

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

	optionalAttrs := []string{"status", "vrf", "tenant", "role", "dns_name", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestIPAddressResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ip_address"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestIPAddressResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewIPAddressResource().(*resources.IPAddressResource)

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

func TestAccIPAddressResource_basic(t *testing.T) {

	address := testutil.RandomIPv4Address()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterIPAddressCleanup(address)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckIPAddressDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(address),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
				),
			},
		},
	})

}

func TestAccIPAddressResource_full(t *testing.T) {

	address := testutil.RandomIPv4Address()

	description := "Test IP address with all fields"

	dnsName := "test.example.com"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterIPAddressCleanup(address)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckIPAddressDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_full(address, description, dnsName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", description),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", dnsName),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "active"),
				),
			},
		},
	})

}

func TestAccIPAddressResource_withVRF(t *testing.T) {

	address := testutil.RandomIPv4Address()

	vrfName := testutil.RandomName("tf-test-vrf")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterIPAddressCleanup(address)

	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckIPAddressDestroy,

			testutil.CheckVRFDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_withVRF(address, vrfName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "vrf"),
				),
			},
		},
	})

}

func TestAccIPAddressResource_ipv6(t *testing.T) {

	address := testutil.RandomIPv6Address()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterIPAddressCleanup(address)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckIPAddressDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(address),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
				),
			},
		},
	})

}

func TestAccIPAddressResource_update(t *testing.T) {

	address := testutil.RandomIPv4Address()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterIPAddressCleanup(address)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckIPAddressDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccIPAddressResourceConfig_basic(address),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),
				),
			},

			{

				Config: testAccIPAddressResourceConfig_full(address, "Updated description", "updated.example.com"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_address.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "address", address),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "dns_name", "updated.example.com"),

					resource.TestCheckResourceAttr("netbox_ip_address.test", "status", "active"),
				),
			},
		},
	})

}

func testAccIPAddressResourceConfig_basic(address string) string {

	return fmt.Sprintf(`



resource "netbox_ip_address" "test" {



  address = %q



}



`, address)

}

func testAccIPAddressResourceConfig_full(address, description, dnsName string) string {

	return fmt.Sprintf(`



resource "netbox_ip_address" "test" {



  address     = %q



  description = %q



  dns_name    = %q



  status      = "active"



}



`, address, description, dnsName)

}

func testAccIPAddressResourceConfig_withVRF(address, vrfName string) string {

	return fmt.Sprintf(`



resource "netbox_vrf" "test" {



  name = %q



}







resource "netbox_ip_address" "test" {



  address = %q



  vrf     = netbox_vrf.test.name



}



`, vrfName, address)

}
