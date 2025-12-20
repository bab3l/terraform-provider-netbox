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

func TestRIRResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource()

	if r == nil {

		t.Fatal("Expected non-nil RIR resource")

	}

}

func TestRIRResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource()

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

	optionalAttrs := []string{"is_private", "description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestRIRResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_rir"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestRIRResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRIRResource().(*resources.RIRResource)

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

func TestAccRIRResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-rir")

	slug := testutil.RandomSlug("tf-test-rir")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRIRResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rir.test", "id"),

					resource.TestCheckResourceAttr("netbox_rir.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rir.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_rir.test", "is_private", "false"),
				),
			},

			{

				ResourceName: "netbox_rir.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccRIRResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-rir-full")

	slug := testutil.RandomSlug("tf-test-rir-full")

	description := "Test RIR with all fields"

	updatedDescription := "Updated RIR description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRIRResourceConfig_full(name, slug, description, true),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rir.test", "id"),

					resource.TestCheckResourceAttr("netbox_rir.test", "name", name),

					resource.TestCheckResourceAttr("netbox_rir.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_rir.test", "description", description),

					resource.TestCheckResourceAttr("netbox_rir.test", "is_private", "true"),
				),
			},

			{

				Config: testAccRIRResourceConfig_full(name, slug, updatedDescription, false),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rir.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_rir.test", "is_private", "false"),
				),
			},
		},
	})

}

func testAccRIRResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccRIRResourceConfig_full(name, slug, description string, isPrivate bool) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {

  name        = %q

  slug        = %q

  description = %q

  is_private  = %t

}

`, name, slug, description, isPrivate)

}
