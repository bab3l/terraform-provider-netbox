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

func TestTagResource(t *testing.T) {

	t.Parallel()

	r := resources.NewTagResource()

	if r == nil {

		t.Fatal("Expected non-nil tag resource")

	}

}

func TestTagResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewTagResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check required attributes

	requiredAttrs := []string{"name", "slug"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	// Check computed attributes

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	// Check optional attributes

	optionalAttrs := []string{"color", "description", "object_types"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestTagResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewTagResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tag"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestTagResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewTagResource().(*resources.TagResource)

	// Test with nil provider data (should not error)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	// Test with correct provider data

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

}

// testAccTagResourceBasic creates a basic tag with required fields only.

func testAccTagResourceBasic(name, slug string) string {

	return fmt.Sprintf(`































resource "netbox_tag" "test" {































  name = %q































  slug = %q































}































`, name, slug)

}

// testAccTagResourceFull creates a tag with all optional fields.

func testAccTagResourceFull(name, slug, color, description string) string {

	return fmt.Sprintf(`































resource "netbox_tag" "test" {































  name        = %q































  slug        = %q































  color       = %q































  description = %q































}































`, name, slug, color, description)

}

// testAccTagResourceWithObjectTypes creates a tag with object_types restriction.

func testAccTagResourceWithObjectTypes(name, slug string) string {

	return fmt.Sprintf(`































resource "netbox_tag" "test" {































  name         = %q































  slug         = %q































  object_types = ["dcim.device", "dcim.site"]































}































`, name, slug)

}

func TestAccTagResource_basic(t *testing.T) {

	name := testutil.RandomName("tag")

	slug := testutil.RandomSlug("tag")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTagResourceBasic(name, slug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),

					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_tag.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccTagResource_full(t *testing.T) {

	name := testutil.RandomName("tag")

	slug := testutil.RandomSlug("tag")

	color := "ff5722"

	description := "Test tag description"

	updatedName := testutil.RandomName("tag-updated")

	updatedSlug := testutil.RandomSlug("tag-updated")

	updatedColor := "2196f3"

	updatedDescription := "Updated test tag description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTagResourceFull(name, slug, color, description),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),

					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_tag.test", "color", color),

					resource.TestCheckResourceAttr("netbox_tag.test", "description", description),
				),
			},

			{

				Config: testAccTagResourceFull(updatedName, updatedSlug, updatedColor, updatedDescription),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tag.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_tag.test", "slug", updatedSlug),

					resource.TestCheckResourceAttr("netbox_tag.test", "color", updatedColor),

					resource.TestCheckResourceAttr("netbox_tag.test", "description", updatedDescription),
				),
			},
		},
	})

}

func TestAccTagResource_withObjectTypes(t *testing.T) {

	name := testutil.RandomName("tag")

	slug := testutil.RandomSlug("tag")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccTagResourceWithObjectTypes(name, slug),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tag.test", "id"),

					resource.TestCheckResourceAttr("netbox_tag.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tag.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_tag.test", "object_types.#", "2"),
				),
			},
		},
	})

}
