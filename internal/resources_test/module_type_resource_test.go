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

func TestModuleTypeResource(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleTypeResource()

	if r == nil {

		t.Fatal("Expected non-nil ModuleType resource")

	}

}

func TestModuleTypeResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleTypeResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"manufacturer", "model"}

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

	optionalAttrs := []string{"part_number", "airflow", "weight", "weight_unit", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestModuleTypeResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleTypeResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_module_type"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestModuleTypeResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewModuleTypeResource().(*resources.ModuleTypeResource)

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

func TestAccModuleTypeResource_basic(t *testing.T) {

	mfgName := testutil.RandomName("tf-test-mfg")

	mfgSlug := testutil.RandomSlug("tf-test-mfg")

	model := testutil.RandomName("tf-test-module-type")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),
				),
			},

			{

				ResourceName: "netbox_module_type.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccModuleTypeResource_full(t *testing.T) {

	mfgName := testutil.RandomName("tf-test-mfg-full")

	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")

	model := testutil.RandomName("tf-test-module-type-full")

	description := "Test module type with all fields"

	updatedDescription := "Updated module type description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_module_type.test", "id"),

					resource.TestCheckResourceAttr("netbox_module_type.test", "model", model),

					resource.TestCheckResourceAttr("netbox_module_type.test", "part_number", "MT-001"),

					resource.TestCheckResourceAttr("netbox_module_type.test", "description", description),
				),
			},

			{

				Config: testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_module_type.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccModuleTypeResourceConfig_basic(mfgName, mfgSlug, model string) string {

	return fmt.Sprintf(`



resource "netbox_manufacturer" "test" {



  name = %q



  slug = %q



}







resource "netbox_module_type" "test" {



  manufacturer = netbox_manufacturer.test.id



  model        = %q



}



`, mfgName, mfgSlug, model)

}

func testAccModuleTypeResourceConfig_full(mfgName, mfgSlug, model, description string) string {

	return fmt.Sprintf(`



resource "netbox_manufacturer" "test" {



  name = %q



  slug = %q



}







resource "netbox_module_type" "test" {



  manufacturer = netbox_manufacturer.test.id



  model        = %q



  part_number  = "MT-001"



  description  = %q



}



`, mfgName, mfgSlug, model, description)

}
