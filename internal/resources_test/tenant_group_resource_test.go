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

func TestTenantGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil tenant group resource")

	}

}

func TestTenantGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantGroupResource()

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

	optionalAttrs := []string{"parent", "description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

}

func TestTenantGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantGroupResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tenant_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestTenantGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantGroupResource().(*resources.TenantGroupResource)

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

	if r.GetClient() != client {

		t.Error("Expected client to be set")

	}

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccTenantGroupResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-tenant-group")

	slug := testutil.RandomSlug("tf-test-tg")

	// Register cleanup to ensure resource is deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccTenantGroupResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tenant-group-full")

	slug := testutil.RandomSlug("tf-test-tg-full")

	description := "Test tenant group with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "description", description),
				),
			},
		},
	})

}

func TestAccTenantGroupResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tenant-group-update")

	slug := testutil.RandomSlug("tf-test-tg-upd")

	updatedName := testutil.RandomName("tf-test-tenant-group-updated")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),
				),
			},

			{

				Config: testAccTenantGroupResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", updatedName),
				),
			},
		},
	})

}

// testAccTenantGroupResourceConfig_basic returns a basic test configuration.

func testAccTenantGroupResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`



terraform {



  required_providers {



    netbox = {



      source = "bab3l/netbox"



      version = ">= 0.1.0"

    }

  }

}



provider "netbox" {}



resource "netbox_tenant_group" "test" {

  name = %q

  slug = %q

}



`, name, slug)

}

// testAccTenantGroupResourceConfig_full returns a test configuration with all fields.

func testAccTenantGroupResourceConfig_full(name, slug, description string) string {

	return fmt.Sprintf(`



terraform {



  required_providers {



    netbox = {



      source = "bab3l/netbox"



      version = ">= 0.1.0"

    }

  }

}



provider "netbox" {}



resource "netbox_tenant_group" "test" {

  name        = %q

  slug        = %q

  description = %q

}



`, name, slug, description)

}

func TestAccTenantGroupResource_import(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-tenant-group-import")

	slug := testutil.RandomSlug("tf-test-tenant-group-imp")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantGroupResourceConfig_import(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_tenant_group.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccTenantGroupResourceConfig_import(name, slug string) string {

	return fmt.Sprintf(`

resource "netbox_tenant_group" "test" {

  name = %q

  slug = %q

}

`, name, slug)

}
