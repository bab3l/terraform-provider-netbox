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

func TestASNResource(t *testing.T) {
	t.Parallel()

	r := resources.NewASNResource()
	if r == nil {
		t.Fatal("Expected non-nil ASN resource")
	}
}

func TestASNResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewASNResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"asn"}
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

	optionalAttrs := []string{"rir", "tenant", "description", "comments", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestASNResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewASNResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_asn"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestASNResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewASNResource().(*resources.ASNResource)

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

	configureRequest.ProviderData = "invalid"
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with incorrect provider data")
	}
}

func TestAccASNResource_basic(t *testing.T) {
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	// Generate a random ASN in the private range (64512-65534)
	asn := int64(acctest.RandIntRange(64512, 65534))

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_basic(rirName, rirSlug, asn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
				),
			},
			{
				ResourceName:      "netbox_asn.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccASNResource_full(t *testing.T) {
	rirName := testutil.RandomName("tf-test-rir")
	rirSlug := testutil.RandomSlug("tf-test-rir")
	// Generate a random ASN in the private range (64512-65534)
	asn := int64(acctest.RandIntRange(64512, 65534))
	description := "Test ASN with all fields"
	updatedDescription := "Updated ASN description"
	comments := "Test comments"

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, asn, description, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_asn.test", "id"),
					resource.TestCheckResourceAttr("netbox_asn.test", "asn", fmt.Sprintf("%d", asn)),
					resource.TestCheckResourceAttr("netbox_asn.test", "description", description),
					resource.TestCheckResourceAttr("netbox_asn.test", "comments", comments),
					resource.TestCheckResourceAttrSet("netbox_asn.test", "rir"),
				),
			},
			{
				Config: testAccASNResourceConfig_full(rirName, rirSlug, asn, updatedDescription, comments),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_asn.test", "description", updatedDescription),
				),
			},
		},
	})
}

func testAccASNResourceConfig_basic(rirName, rirSlug string, asn int64) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn = %d
  rir = netbox_rir.test.id
}
`, rirName, rirSlug, asn)
}

func testAccASNResourceConfig_full(rirName, rirSlug string, asn int64, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_rir" "test" {
  name = %q
  slug = %q
}

resource "netbox_asn" "test" {
  asn         = %d
  rir         = netbox_rir.test.id
  description = %q
  comments    = %q
}
`, rirName, rirSlug, asn, description, comments)
}
