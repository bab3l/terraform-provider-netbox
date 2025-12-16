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

func TestInventoryItemRoleResource(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemRoleResource()

	if r == nil {

		t.Fatal("Expected non-nil InventoryItemRole resource")

	}

}

func TestInventoryItemRoleResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemRoleResource()

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

	optionalAttrs := []string{"color", "description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestInventoryItemRoleResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemRoleResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_inventory_item_role"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestInventoryItemRoleResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewInventoryItemRoleResource().(*resources.InventoryItemRoleResource)

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

func TestAccInventoryItemRoleResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-inv-role")

	slug := testutil.RandomSlug("tf-test-inv-role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemRoleResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_inventory_item_role.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccInventoryItemRoleResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-inv-role-full")

	slug := testutil.RandomSlug("tf-test-inv-role-full")

	description := "Test inventory item role with all fields"

	updatedDescription := "Updated inventory item role description"

	color := "aa1409"

	updatedColor := "2196f3"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccInventoryItemRoleResourceConfig_full(name, slug, color, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_inventory_item_role.test", "id"),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "name", name),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", color),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", description),
				),
			},

			{

				Config: testAccInventoryItemRoleResourceConfig_full(name, slug, updatedColor, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "color", updatedColor),

					resource.TestCheckResourceAttr("netbox_inventory_item_role.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccInventoryItemRoleResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`















resource "netbox_inventory_item_role" "test" {















  name = %q















  slug = %q















}















`, name, slug)

}

func testAccInventoryItemRoleResourceConfig_full(name, slug, color, description string) string {

	return fmt.Sprintf(`















resource "netbox_inventory_item_role" "test" {















  name        = %q















  slug        = %q















  color       = %q















  description = %q















}















`, name, slug, color, description)

}
