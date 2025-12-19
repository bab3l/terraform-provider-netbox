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

func TestRoleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource()

	if r == nil {

		t.Fatal("Expected non-nil Role resource")

	}

}

func TestRoleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource()

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

	optionalAttrs := []string{"weight", "description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestRoleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_role"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestRoleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRoleResource().(*resources.RoleResource)

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

func TestAccRoleResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-role")

	slug := testutil.RandomSlug("tf-test-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "1000"),
				),
			},

			{

				ResourceName: "netbox_role.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccRoleResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-role-full")

	slug := testutil.RandomSlug("tf-test-role-full")

	description := "Test IPAM role with all fields"

	updatedDescription := "Updated IPAM role description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccRoleResourceConfig_full(name, slug, description, 100),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_role.test", "description", description),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "100"),
				),
			},

			{

				Config: testAccRoleResourceConfig_full(name, slug, updatedDescription, 200),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_role.test", "description", updatedDescription),

					resource.TestCheckResourceAttr("netbox_role.test", "weight", "200"),
				),
			},
		},
	})

}

func testAccRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`



resource "netbox_role" "test" {

  name = %q

  slug = %q

}



`, name, slug)

}

func testAccRoleResourceConfig_full(name, slug, description string, weight int) string {

	return fmt.Sprintf(`



resource "netbox_role" "test" {

  name        = %q

  slug        = %q

  description = %q

  weight      = %d

}



`, name, slug, description, weight)

}
