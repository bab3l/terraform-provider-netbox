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

func TestRackTypeResource(t *testing.T) {
	t.Parallel()
	r := resources.NewRackTypeResource()
	if r == nil {
		t.Fatal("Expected non-nil RackType resource")
	}
}

func TestRackTypeResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewRackTypeResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}
	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"manufacturer", "model", "slug"}
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

	optionalAttrs := []string{"description", "form_factor", "width", "u_height", "starting_unit", "desc_units", "outer_width", "outer_depth", "outer_unit", "weight", "max_weight", "weight_unit", "mounting_depth", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestRackTypeResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewRackTypeResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}
	r.Metadata(context.Background(), metadataRequest, metadataResponse)
	expected := "netbox_rack_type"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestRackTypeResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewRackTypeResource().(*resources.RackTypeResource)
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

func TestAccRackTypeResource_basic(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg")
	mfgSlug := testutil.RandomSlug("tf-test-mfg")
	model := testutil.RandomName("tf-test-rack-type")
	slug := testutil.RandomSlug("tf-test-rack-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),
				),
			},
			{
				ResourceName:            "netbox_rack_type.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"manufacturer"},
			},
		},
	})
}

func TestAccRackTypeResource_full(t *testing.T) {
	mfgName := testutil.RandomName("tf-test-mfg-full")
	mfgSlug := testutil.RandomSlug("tf-test-mfg-full")
	model := testutil.RandomName("tf-test-rack-type-full")
	slug := testutil.RandomSlug("tf-test-rack-type-full")
	description := "Test rack type with all fields"
	updatedDescription := "Updated rack type description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterManufacturerCleanup(mfgSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, description, 42, 19),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_rack_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", description),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "u_height", "42"),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "width", "19"),
				),
			},
			{
				Config: testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, updatedDescription, 48, 19),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "u_height", "48"),
				),
			},
		},
	})
}

func testAccRackTypeResourceConfig_basic(mfgName, mfgSlug, model, slug string) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  form_factor  = "4-post-cabinet"
}
`, mfgName, mfgSlug, model, slug)
}

func testAccRackTypeResourceConfig_full(mfgName, mfgSlug, model, slug, description string, uHeight, width int) string {
	return fmt.Sprintf(`
resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_type" "test" {
  manufacturer = netbox_manufacturer.test.id
  model        = %q
  slug         = %q
  description  = %q
  u_height     = %d
  width        = %d
  form_factor  = "4-post-cabinet"
}
`, mfgName, mfgSlug, model, slug, description, uHeight, width)
}

// TestAccConsistency_RackType_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_RackType_LiteralNames(t *testing.T) {
	t.Parallel()
	mfgName := testutil.RandomName("manufacturer")
	mfgSlug := testutil.RandomSlug("manufacturer")
	model := testutil.RandomName("rack-type")
	slug := testutil.RandomSlug("rack-type")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_rack_type.test", "model", model),
					resource.TestCheckResourceAttr("netbox_rack_type.test", "manufacturer", mfgName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug),
			},
		},
	})
}

func testAccRackTypeConsistencyLiteralNamesConfig(mfgName, mfgSlug, model, slug string) string {
	return fmt.Sprintf(`

resource "netbox_manufacturer" "test" {
  name = %q
  slug = %q
}

resource "netbox_rack_type" "test" {
  # Use literal string name to mimic existing user state
  manufacturer = %q
  model        = %q
  slug         = %q
  u_height     = 42
  width        = 19
  form_factor  = "4-post-cabinet"

  depends_on = [netbox_manufacturer.test]
}

`, mfgName, mfgSlug, mfgName, model, slug)
}
