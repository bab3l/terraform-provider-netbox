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

func TestTenantResource(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantResource()

	if r == nil {

		t.Fatal("Expected non-nil tenant resource")

	}

}

func TestTenantResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantResource()

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

	optionalAttrs := []string{"group", "description", "tags", "custom_fields"}

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

func TestTenantResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_tenant"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestTenantResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewTenantResource().(*resources.TenantResource)

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

	configureRequest.ProviderData = testutil.InvalidProviderData

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccTenantResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-tenant")

	slug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup to ensure resource is deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccTenantResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tenant-full")

	slug := testutil.RandomSlug("tf-test-tenant-full")

	description := "Test tenant with all fields"

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_tenant.test", "description", description),
				),
			},
		},
	})

}

func TestAccTenantResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-tenant-update")

	slug := testutil.RandomSlug("tf-test-tenant-upd")

	updatedName := testutil.RandomName("tf-test-tenant-updated")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),
				),
			},

			{

				Config: testAccTenantResourceConfig_basic(updatedName, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_tenant.test", "id"),

					resource.TestCheckResourceAttr("netbox_tenant.test", "name", updatedName),
				),
			},
		},
	})

}

// testAccTenantResourceConfig_basic returns a basic test configuration.

func testAccTenantResourceConfig_basic(name, slug string) string {

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







resource "netbox_tenant" "test" {



  name = %q



  slug = %q



}







`, name, slug)

}

// testAccTenantResourceConfig_full returns a test configuration with all fields.

func testAccTenantResourceConfig_full(name, slug, description string) string {

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







resource "netbox_tenant" "test" {



  name        = %q



  slug        = %q



  description = %q



}







`, name, slug, description)

}

func TestAccTenantResource_import(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-tenant-import")

	slug := testutil.RandomSlug("tf-test-tenant-imp")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterTenantCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckTenantDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantResourceConfig_import(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tenant.test", "name", name),

					resource.TestCheckResourceAttr("netbox_tenant.test", "slug", slug),
				),
			},

			{

				ResourceName: "netbox_tenant.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func testAccTenantResourceConfig_import(name, slug string) string {

	return fmt.Sprintf(`







resource "netbox_tenant" "test" {



  name = %q



  slug = %q



}







`, name, slug)

}

func TestAccConsistency_Tenant(t *testing.T) {

	t.Parallel()

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	groupName := testutil.RandomName("group")

	groupSlug := testutil.RandomSlug("group")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_tenant.test", "name", tenantName),

					resource.TestCheckResourceAttr("netbox_tenant.test", "group", groupName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug),
			},
		},
	})

}

func testAccTenantConsistencyConfig(tenantName, tenantSlug, groupName, groupSlug string) string {

	return fmt.Sprintf(`







resource "netbox_tenant_group" "test" {



  name = "%[3]s"



  slug = "%[4]s"



}







resource "netbox_tenant" "test" {



  name = "%[1]s"



  slug = "%[2]s"



  group = netbox_tenant_group.test.name



}







`, tenantName, tenantSlug, groupName, groupSlug)

}
