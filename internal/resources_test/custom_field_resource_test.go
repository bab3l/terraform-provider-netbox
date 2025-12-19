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
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestCustomFieldResource(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource()

	if r == nil {

		t.Fatal("Expected non-nil CustomField resource")

	}

}

func TestCustomFieldResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"object_types", "type", "name"}

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

	optionalAttrs := []string{"label", "group_name", "description", "required", "search_weight", "filter_logic", "ui_visible", "ui_editable", "is_cloneable", "default", "weight", "validation_minimum", "validation_maximum", "validation_regex", "choice_set", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestCustomFieldResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_custom_field"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestCustomFieldResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewCustomFieldResource().(*resources.CustomFieldResource)

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

func TestAccCustomFieldResource_basic(t *testing.T) {

	// Custom field names can only contain alphanumeric characters and underscores

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldResourceConfig_basic(name),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "text"),
				),
			},

			{

				ResourceName: "netbox_custom_field.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccCustomFieldResource_full(t *testing.T) {

	// Custom field names can only contain alphanumeric characters and underscores

	name := fmt.Sprintf("tf_test_%s", acctest.RandString(8))

	description := "Test custom field with all fields"

	updatedDescription := "Updated custom field description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccCustomFieldResourceConfig_full(name, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_custom_field.test", "id"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "name", name),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "type", "integer"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", description),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "required", "true"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "validation_minimum", "1"),

					resource.TestCheckResourceAttr("netbox_custom_field.test", "validation_maximum", "100"),
				),
			},

			{

				Config: testAccCustomFieldResourceConfig_full(name, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_custom_field.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccCustomFieldResourceConfig_basic(name string) string {

	return fmt.Sprintf(`



resource "netbox_custom_field" "test" {

  name         = %q

  type         = "text"



  object_types = ["dcim.site"]

}



`, name)

}

func testAccCustomFieldResourceConfig_full(name, description string) string {

	return fmt.Sprintf(`



resource "netbox_custom_field" "test" {

  name               = %q

  type               = "integer"



  object_types       = ["dcim.site", "dcim.device"]

  description        = %q



  required           = true



  validation_minimum = 1



  validation_maximum = 100

  weight             = 50

}



`, name, description)

}
