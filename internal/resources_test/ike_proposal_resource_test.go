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

func TestIKEProposalResource(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEProposalResource()
	if r == nil {
		t.Fatal("Expected non-nil IKEProposal resource")
	}
}

func TestIKEProposalResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEProposalResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "authentication_method", "encryption_algorithm", "authentication_algorithm", "group"}
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

	optionalAttrs := []string{"description", "comments", "tags", "sa_lifetime"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestIKEProposalResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEProposalResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ike_proposal"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestIKEProposalResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEProposalResource()

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

func TestAccIKEProposalResource_basic(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ike-proposal")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIKEProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_method", "preshared-keys"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-256-cbc"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_algorithm", "hmac-sha256"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "group", "14"),
				),
			},
		},
	})
}

func TestAccIKEProposalResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-ike-proposal-full")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIKEProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalResourceConfig_full(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_method", "certificates"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-128-gcm"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "authentication_algorithm", "hmac-sha512"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "group", "19"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "sa_lifetime", "28800"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "description", "Test IKE proposal with full options"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "comments", "Test comments for IKE proposal"),
				),
			},
		},
	})
}

func TestAccIKEProposalResource_update(t *testing.T) {
	// Generate unique name
	name := testutil.RandomName("tf-test-ike-proposal-update")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIKEProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-256-cbc"),
				),
			},
			{
				Config: testAccIKEProposalResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "encryption_algorithm", "aes-128-cbc"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccIKEProposalResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}
`, name)
}

func testAccIKEProposalResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %q
  authentication_method    = "certificates"
  encryption_algorithm     = "aes-128-gcm"
  authentication_algorithm = "hmac-sha512"
  group                    = 19
  sa_lifetime              = 28800
  description              = "Test IKE proposal with full options"
  comments                 = "Test comments for IKE proposal"
}
`, name)
}

func testAccIKEProposalResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-128-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
  description              = "Updated description"
}
`, name)
}

func TestAccIKEProposalResource_import(t *testing.T) {
	name := testutil.RandomName("tf-test-ike-proposal")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEProposalCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIKEProposalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEProposalResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_proposal.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_proposal.test", "name", name),
				),
			},
			{
				ResourceName:      "netbox_ike_proposal.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
