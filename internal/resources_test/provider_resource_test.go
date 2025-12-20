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

func TestProviderResource(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource()

	if r == nil {

		t.Fatal("Expected non-nil Provider resource")

	}

}

func TestProviderResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource()

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

	optionalAttrs := []string{"description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestProviderResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_provider"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestProviderResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderResource().(*resources.ProviderResource)

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

func TestAccProviderResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-provider")

	slug := testutil.RandomSlug("tf-test-provider")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),

					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccProviderResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-provider-full")

	slug := testutil.RandomSlug("tf-test-provider-full")

	description := "Test circuit provider with all fields"

	comments := "Test comments for circuit provider"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_full(name, slug, description, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),

					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_provider.test", "description", description),

					resource.TestCheckResourceAttr("netbox_provider.test", "comments", comments),
				),
			},
		},
	})

}

func TestAccProviderResource_update(t *testing.T) {

	name := testutil.RandomName("tf-test-provider-update")

	slug := testutil.RandomSlug("tf-test-provider-update")

	updatedName := testutil.RandomName("tf-test-provider-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),
				),
			},

			{

				Config: testAccProviderResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_provider.test", "name", updatedName),
				),
			},
		},
	})

}

func testAccProviderResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}

func testAccProviderResourceConfig_full(name, slug, description, comments string) string {

	return fmt.Sprintf(`

resource "netbox_provider" "test" {

  name        = %q

  slug        = %q

  description = %q

  comments    = %q

}

`, name, slug, description, comments)

}

func TestAccProviderResource_import(t *testing.T) {

	name := testutil.RandomName("tf-test-provider")

	slug := testutil.RandomSlug("tf-test-provider")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckProviderDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccProviderResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider.test", "name", name),

					resource.TestCheckResourceAttr("netbox_provider.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_provider.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}
