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

func TestVRFResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVRFResource()

	if r == nil {

		t.Fatal("Expected non-nil VRF resource")

	}

}

func TestVRFResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVRFResource()

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

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	optionalAttrs := []string{"rd", "tenant", "enforce_unique", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestVRFResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVRFResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_vrf"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestVRFResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVRFResource().(*resources.VRFResource)

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

func TestAccVRFResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-vrf")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVRFCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVRFDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVRFResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),

					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
				),
			},
		},
	})

}

func TestAccVRFResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-vrf-full")

	rd := "65000:100"

	description := "Test VRF with all fields"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVRFCleanup(name)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVRFDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVRFResourceConfig_full(name, rd, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),

					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),

					resource.TestCheckResourceAttr("netbox_vrf.test", "rd", rd),

					resource.TestCheckResourceAttr("netbox_vrf.test", "description", description),

					resource.TestCheckResourceAttr("netbox_vrf.test", "enforce_unique", "true"),
				),
			},
		},
	})

}

func TestAccVRFResource_update(t *testing.T) {

	name := testutil.RandomName("tf-test-vrf-update")

	updatedName := testutil.RandomName("tf-test-vrf-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVRFCleanup(name)

	cleanup.RegisterVRFCleanup(updatedName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVRFDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVRFResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),

					resource.TestCheckResourceAttr("netbox_vrf.test", "name", name),
				),
			},

			{

				Config: testAccVRFResourceConfig_full(updatedName, "65000:200", "Updated description"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vrf.test", "id"),

					resource.TestCheckResourceAttr("netbox_vrf.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_vrf.test", "rd", "65000:200"),

					resource.TestCheckResourceAttr("netbox_vrf.test", "description", "Updated description"),
				),
			},
		},
	})

}

func testAccVRFResourceConfig_basic(name string) string {

	return fmt.Sprintf(`







resource "netbox_vrf" "test" {







  name = %q







}







`, name)

}

func testAccVRFResourceConfig_full(name, rd, description string) string {

	return fmt.Sprintf(`







resource "netbox_vrf" "test" {







  name           = %q







  rd             = %q







  description    = %q







  enforce_unique = true







}







`, name, rd, description)

}
