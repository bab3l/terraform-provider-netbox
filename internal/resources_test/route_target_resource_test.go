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

func TestAccRouteTargetResource_import(t *testing.T) {

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

			{

				ResourceName: "netbox_route_target.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

// TestAccConsistency_RouteTarget_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_RouteTarget_LiteralNames(t *testing.T) {

	t.Parallel()

	rtName := testutil.RandomName("65000:100")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_route_target.test", "name", rtName),

					resource.TestCheckResourceAttr("netbox_route_target.test", "tenant", tenantName),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug),
			},
		},
	})

}

func testAccRouteTargetConsistencyLiteralNamesConfig(rtName, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`



resource "netbox_tenant" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}



resource "netbox_route_target" "test" {

  name = "%[1]s"

  # Use literal string name to mimic existing user state

  tenant = "%[2]s"



  depends_on = [netbox_tenant.test]

}



`, rtName, tenantName, tenantSlug)

}
