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

func TestIPSecPolicyResource(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecPolicyResource()
	if r == nil {
		t.Fatal("Expected non-nil IPSecPolicy resource")
	}
}

func TestIPSecPolicyResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecPolicyResource()
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

	optionalAttrs := []string{"description", "comments", "tags", "proposals", "pfs_group"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestIPSecPolicyResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecPolicyResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ipsec_policy"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestIPSecPolicyResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecPolicyResource()

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

func TestAccIPSecPolicyResource_basic(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ipsec-policy")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
		},
	})
}

func TestAccIPSecPolicyResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-ipsec-policy-full")
	proposalName := testutil.RandomName("tf-test-ipsec-proposal-for-policy")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)
	cleanup.RegisterIPSecProposalCleanup(proposalName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecPolicyDestroy,
			testutil.CheckIPSecProposalDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyResourceConfig_full(name, proposalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "14"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "description", "Test IPSec policy with full options"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "comments", "Test comments for IPSec policy"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "proposals.#", "1"),
				),
			},
		},
	})
}

func TestAccIPSecPolicyResource_update(t *testing.T) {
	// Generate unique name
	name := testutil.RandomName("tf-test-ipsec-policy-update")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
			{
				Config: testAccIPSecPolicyResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "pfs_group", "19"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccIPSecPolicyResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_policy" "test" {
  name = %q
}
`, name)
}

func testAccIPSecPolicyResourceConfig_full(name, proposalName string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_proposal" "test" {
  name                     = %q
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
}

resource "netbox_ipsec_policy" "test" {
  name        = %q
  proposals   = [netbox_ipsec_proposal.test.id]
  pfs_group   = 14
  description = "Test IPSec policy with full options"
  comments    = "Test comments for IPSec policy"
}
`, proposalName, name)
}

func testAccIPSecPolicyResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_ipsec_policy" "test" {
  name        = %q
  pfs_group   = 19
  description = "Updated description"
}
`, name)
}

func TestAccIPSecPolicyResource_import(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ipsec-policy")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIPSecPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_policy.test", "name", name),
				),
			},
			{
				ResourceName:      "netbox_ipsec_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
