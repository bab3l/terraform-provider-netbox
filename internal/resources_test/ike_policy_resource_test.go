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

func TestIKEPolicyResource(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEPolicyResource()
	if r == nil {
		t.Fatal("Expected non-nil IKEPolicy resource")
	}
}

func TestIKEPolicyResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEPolicyResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "version", "mode"}
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

	optionalAttrs := []string{"description", "comments", "tags", "proposals", "preshared_key"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestIKEPolicyResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEPolicyResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ike_policy"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestIKEPolicyResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewIKEPolicyResource()

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

func TestAccIKEPolicyResource_basic(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ike-policy")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},
		},
	})
}

func TestAccIKEPolicyResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-ike-policy-full")
	proposalName := testutil.RandomName("tf-test-ike-proposal-for-policy")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)
	cleanup.RegisterIKEProposalCleanup(proposalName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEPolicyDestroy,
			testutil.CheckIKEProposalDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_full(name, proposalName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "1"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "mode", "aggressive"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", "Test IKE policy with full options"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "comments", "Test comments for IKE policy"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "proposals.#", "1"),
				),
			},
		},
	})
}

func TestAccIKEPolicyResource_update(t *testing.T) {
	// Generate unique name
	name := testutil.RandomName("tf-test-ike-policy-update")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
				),
			},
			{
				Config: testAccIKEPolicyResourceConfig_updated(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "version", "2"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccIKEPolicyResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = 2
}
`, name)
}

func testAccIKEPolicyResourceConfig_full(name, proposalName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_proposal" "test" {
  name                     = %q
  authentication_method    = "preshared-keys"
  encryption_algorithm     = "aes-256-cbc"
  authentication_algorithm = "hmac-sha256"
  group                    = 14
}

resource "netbox_ike_policy" "test" {
  name        = %q
  version     = 1
  mode        = "aggressive"
  proposals   = [netbox_ike_proposal.test.id]
  description = "Test IKE policy with full options"
  comments    = "Test comments for IKE policy"
}
`, proposalName, name)
}

func testAccIKEPolicyResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name        = %q
  version     = 2
  description = "Updated description"
}
`, name)
}

func TestAccIKEPolicyResource_import(t *testing.T) {
	name := testutil.RandomName("tf-test-ike-policy")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIKEPolicyCleanup(name)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckIKEPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ike_policy.test", "id"),
					resource.TestCheckResourceAttr("netbox_ike_policy.test", "name", name),
				),
			},
			{
				ResourceName:      "netbox_ike_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
