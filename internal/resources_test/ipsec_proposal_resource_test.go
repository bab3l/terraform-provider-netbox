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

func TestIPSecProposalResource(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProposalResource()
	if r == nil {
		t.Fatal("Expected non-nil IPSecProposal resource")
	}
}

func TestIPSecProposalResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProposalResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "encryption_algorithm", "authentication_algorithm"}
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

	optionalAttrs := []string{"description", "comments", "tags", "sa_lifetime_seconds", "sa_lifetime_data"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestIPSecProposalResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProposalResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ipsec_proposal"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestIPSecProposalResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProposalResource()

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

func TestAccIPSecProposalResource_basic(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ipsec-proposal")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIPSecProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-256-cbc"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha256"),
				),
			},
		},
	})
}

func TestAccIPSecProposalResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-ipsec-proposal-full")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIPSecProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProposalResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-128-gcm"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha512"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "sa_lifetime_seconds", "28800"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "sa_lifetime_data", "102400"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "description", "Test IPSec proposal with full options"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "comments", "Test comments for IPSec proposal"),
				),
			},
		},
	})
}

func TestAccIPSecProposalResource_update(t *testing.T) {
	// Generate unique name
	name := testutil.RandomName("tf-test-ipsec-proposal-update")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIPSecProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-256-cbc"),
				),
			},
			{
				Config: testAccIPSecProposalResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-128-cbc"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccIPSecProposalResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                     = %q
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
}
`, name)
}

func testAccIPSecProposalResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                     = %q
  encryption_algorithm     = "aes-128-gcm"
  authentication_algorithm = "hmac-sha512"
  sa_lifetime_seconds      = 28800
  sa_lifetime_data         = 102400
  description              = "Test IPSec proposal with full options"
  comments                 = "Test comments for IPSec proposal"
}
`, name)
}

func testAccIPSecProposalResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                     = %q
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
  description              = "Updated description"
}
`, name)
}

func TestAccIPSecProposalResource_import(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ipsec-proposal")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIPSecProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "encryption_algorithm", "aes-256-cbc"),
					resource.TestCheckResourceAttr("netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha256"),
				),
			},
			{
				ResourceName:      "netbox_ipsec_proposal.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
