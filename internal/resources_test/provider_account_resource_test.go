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

func TestProviderAccountResource(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderAccountResource()

	if r == nil {

		t.Fatal("Expected non-nil ProviderAccount resource")

	}

}

func TestProviderAccountResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderAccountResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"circuit_provider", "account"}

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

	optionalAttrs := []string{"name", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestProviderAccountResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderAccountResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_provider_account"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestProviderAccountResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewProviderAccountResource().(*resources.ProviderAccountResource)

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

func TestAccProviderAccountResource_basic(t *testing.T) {

	providerName := testutil.RandomName("tf-test-provider")

	providerSlug := testutil.RandomSlug("tf-test-provider")

	accountID := testutil.RandomName("acct")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),
				),
			},

			{

				ResourceName: "netbox_provider_account.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccProviderAccountResource_full(t *testing.T) {

	providerName := testutil.RandomName("tf-test-provider-full")

	providerSlug := testutil.RandomSlug("tf-test-provider-full")

	accountID := testutil.RandomName("acct")

	accountName := testutil.RandomName("tf-test-acct")

	description := "Test provider account with all fields"

	updatedDescription := "Updated provider account description"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterProviderCleanup(providerSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_provider_account.test", "id"),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "account", accountID),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "name", accountName),

					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", description),
				),
			},

			{

				Config: testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, updatedDescription),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_provider_account.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccProviderAccountResourceConfig_basic(providerName, providerSlug, accountID string) string {

	return fmt.Sprintf(`







resource "netbox_provider" "test" {







  name = %q







  slug = %q







}















resource "netbox_provider_account" "test" {







  circuit_provider = netbox_provider.test.id







  account          = %q







}







`, providerName, providerSlug, accountID)

}

func testAccProviderAccountResourceConfig_full(providerName, providerSlug, accountID, accountName, description string) string {

	return fmt.Sprintf(`







resource "netbox_provider" "test" {







  name = %q







  slug = %q







}















resource "netbox_provider_account" "test" {







  circuit_provider = netbox_provider.test.id







  account          = %q







  name             = %q







  description      = %q







}







`, providerName, providerSlug, accountID, accountName, description)

}
