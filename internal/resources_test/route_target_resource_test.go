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

func TestRouteTargetResource(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()
	if r == nil {
		t.Fatal("Expected non-nil RouteTarget resource")
	}
}

func TestRouteTargetResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name"}
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

	optionalAttrs := []string{"tenant", "description", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestRouteTargetResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_route_target"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestRouteTargetResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewRouteTargetResource()

	// Type assert to access Configure method
	configurable, ok := r.(fwresource.ResourceWithConfigure)
	if !ok {
		t.Fatal("Resource does not implement ResourceWithConfigure")
	}

	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	client := &netbox.APIClient{}
	configureRequest.ProviderData = client

	configurable.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with valid client, got: %+v", configureResponse.Diagnostics)
	}
}

// Acceptance Tests

func TestAccRouteTargetResource_basic(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("65000:100")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
		},
	})
}

func TestAccRouteTargetResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("65000:200")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckRouteTargetDestroy,
			testutil.CheckTenantDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Test route target with full options"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "comments", "Test comments for route target"),
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "tenant"),
				),
			},
		},
	})
}

func TestAccRouteTargetResource_update(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("65000:300")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterRouteTargetCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckRouteTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccRouteTargetResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
				),
			},
			{
				Config: testAccRouteTargetResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_route_target.test", "id"),
					resource.TestCheckResourceAttr("netbox_route_target.test", "name", name),
					resource.TestCheckResourceAttr("netbox_route_target.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccRouteTargetResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "test" {
  name = %q
}
`, name)
}

func testAccRouteTargetResourceConfig_full(name, tenantName, tenantSlug string) string {
	return fmt.Sprintf(`
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_route_target" "test" {
  name        = %q
  description = "Test route target with full options"
  comments    = "Test comments for route target"
  tenant      = netbox_tenant.test.id
}
`, tenantName, tenantSlug, name)
}

func testAccRouteTargetResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_route_target" "test" {
  name        = %q
  description = "Updated description"
}
`, name)
}
