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

func TestVirtualChassisResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualChassisResource()

	if r == nil {

		t.Fatal("Expected non-nil VirtualChassis resource")

	}

}

func TestVirtualChassisResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualChassisResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"name"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	computedAttrs := []string{"id", "member_count"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	optionalAttrs := []string{"domain", "master", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestVirtualChassisResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualChassisResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_virtual_chassis"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestVirtualChassisResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVirtualChassisResource().(*resources.VirtualChassisResource)

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

func TestAccVirtualChassisResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-vc")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualChassisResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),
				),
			},

			{

				ResourceName: "netbox_virtual_chassis.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccVirtualChassisResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-vc-full")

	description := "Test virtual chassis with all fields"

	updatedDescription := "Updated virtual chassis description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccVirtualChassisResourceConfig_full(name, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_virtual_chassis.test", "id"),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "name", name),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "domain", "test-domain"),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", description),

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "comments", "Test comments"),
				),
			},

			{

				Config: testAccVirtualChassisResourceConfig_full(name, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_virtual_chassis.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccVirtualChassisResourceConfig_basic(name string) string {

	return fmt.Sprintf(`



resource "netbox_virtual_chassis" "test" {



  name = %q



}



`, name)

}

func testAccVirtualChassisResourceConfig_full(name, description string) string {

	return fmt.Sprintf(`



resource "netbox_virtual_chassis" "test" {



  name        = %q



  domain      = "test-domain"



  description = %q



  comments    = "Test comments"



}



`, name, description)

}
