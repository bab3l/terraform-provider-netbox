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

func TestASNRangeResource(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

	if r == nil {

		t.Fatal("Expected non-nil ASNRange resource")
	}
}

func TestASNRangeResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "slug", "rir", "start", "end"}

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

	optionalAttrs := []string{"tenant", "description", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestASNRangeResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_asn_range"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestASNRangeResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewASNRangeResource()

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

func TestAccASNRangeResource_basic(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-asn-range")

	slug := testutil.RandomSlug("tf-test-asn-range")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_full(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-asn-range-full")

	slug := testutil.RandomSlug("tf-test-asn-range-full")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	tenantName := testutil.RandomName("tf-test-tenant")

	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,

			testutil.CheckTenantDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "65000"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "65100"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Test ASN range with full options"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "tenant"),
				),
			},
		},
	})
}

func TestAccASNRangeResource_update(t *testing.T) {

	// Generate unique names

	name := testutil.RandomName("tf-test-asn-range-upd")

	slug := testutil.RandomSlug("tf-test-asn-range-upd")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),
				),
			},

			{

				Config: testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64700"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name  = %q
  slug  = %q
  rir   = netbox_rir.test.id
  start = "64512"
  end   = "64612"
}

`, rirName, rirSlug, name, slug)
}

func testAccASNRangeResourceConfig_full(name, slug, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name        = %q
  slug        = %q
  rir         = netbox_rir.test.id
  start       = "65000"
  end         = "65100"
  tenant      = netbox_tenant.test.id
  description = "Test ASN range with full options"
}

`, rirName, rirSlug, tenantName, tenantSlug, name, slug)
}

func testAccASNRangeResourceConfig_updated(name, slug, rirName, rirSlug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn_range" "test" {
  name        = %q
  slug        = %q
  rir         = netbox_rir.test.id
  start       = "64512"
  end         = "64700"
  description = "Updated description"
}

`, rirName, rirSlug, name, slug)
}

func TestAccASNRangeResource_import(t *testing.T) {

	// Generate unique names to avoid conflicts between test runs

	name := testutil.RandomName("tf-test-asn-range")

	slug := testutil.RandomSlug("tf-test-asn-range")

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	// Register cleanup to ensure resources are deleted even if test fails

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterASNRangeCleanup(name)

	cleanup.RegisterRIRCleanup(rirSlug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckASNRangeDestroy,

			testutil.CheckRIRDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeResourceConfig_basic(name, slug, rirName, rirSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", name),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "start", "64512"),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "end", "64612"),

					resource.TestCheckResourceAttrSet("netbox_asn_range.test", "rir"),
				),
			},

			{

				ResourceName: "netbox_asn_range.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})
}

func TestAccConsistency_ASNRange(t *testing.T) {

	t.Parallel()

	rangeName := testutil.RandomName("asn-range")

	rangeSlug := testutil.RandomSlug("asn-range")

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_asn_range.test", "name", rangeName),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_asn_range.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})
}

func testAccASNRangeConsistencyConfig(rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`

resource "netbox_rir" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_tenant" "test" {
  name = "%[5]s"
  slug = "%[6]s"
}

resource "netbox_asn_range" "test" {
  name = "%[1]s"
  slug = "%[2]s"
  rir = netbox_rir.test.slug
  tenant = netbox_tenant.test.name
  start = 65000
  end = 65100
}

`, rangeName, rangeSlug, rirName, rirSlug, tenantName, tenantSlug)
}
