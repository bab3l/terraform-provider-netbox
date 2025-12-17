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

func TestIPSecProfileResource(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProfileResource()
	if r == nil {
		t.Fatal("Expected non-nil IPSecProfile resource")
	}
}

func TestIPSecProfileResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProfileResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "mode", "ike_policy", "ipsec_policy"}
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

	optionalAttrs := []string{"description", "comments", "tags"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestIPSecProfileResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProfileResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ipsec_profile"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestIPSecProfileResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewIPSecProfileResource()

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

func TestAccIPSecProfileResource_basic(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ipsec-profile")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(ikePolicyName)
	cleanup.RegisterIPSecPolicyCleanup(ipsecPolicyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProfileDestroy,
			testutil.CheckIKEPolicyDestroy,
			testutil.CheckIPSecPolicyDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProfileResourceConfig_basic(name, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "esp"),
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "ike_policy"),
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "ipsec_policy"),
				),
			},
		},
	})
}

func TestAccIPSecProfileResource_full(t *testing.T) {
	// Generate unique names
	name := testutil.RandomName("tf-test-ipsec-profile-full")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile-full")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile-full")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(ikePolicyName)
	cleanup.RegisterIPSecPolicyCleanup(ipsecPolicyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProfileDestroy,
			testutil.CheckIKEPolicyDestroy,
			testutil.CheckIPSecPolicyDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProfileResourceConfig_full(name, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "ah"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Test IPSec profile with full options"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "comments", "Test comments for IPSec profile"),
				),
			},
		},
	})
}

func TestAccIPSecProfileResource_update(t *testing.T) {
	// Generate unique name
	name := testutil.RandomName("tf-test-ipsec-profile-update")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile-upd")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile-upd")

	// Register cleanup
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(ikePolicyName)
	cleanup.RegisterIPSecPolicyCleanup(ipsecPolicyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProfileDestroy,
			testutil.CheckIKEPolicyDestroy,
			testutil.CheckIPSecPolicyDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProfileResourceConfig_basic(name, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "esp"),
				),
			},
			{
				Config: testAccIPSecProfileResourceConfig_updated(name, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "ah"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccIPSecProfileResourceConfig_basic(name, ikePolicyName, ipsecPolicyName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = 2
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}

resource "netbox_ipsec_profile" "test" {
  name         = %q
  mode         = "esp"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
}
`, ikePolicyName, ipsecPolicyName, name)
}

func testAccIPSecProfileResourceConfig_full(name, ikePolicyName, ipsecPolicyName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = 2
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}

resource "netbox_ipsec_profile" "test" {
  name         = %q
  mode         = "ah"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
  description  = "Test IPSec profile with full options"
  comments     = "Test comments for IPSec profile"
}
`, ikePolicyName, ipsecPolicyName, name)
}

func testAccIPSecProfileResourceConfig_updated(name, ikePolicyName, ipsecPolicyName string) string {
	return fmt.Sprintf(`
resource "netbox_ike_policy" "test" {
  name    = %q
  version = 2
}

resource "netbox_ipsec_policy" "test" {
  name = %q
}

resource "netbox_ipsec_profile" "test" {
  name         = %q
  mode         = "ah"
  ike_policy   = netbox_ike_policy.test.id
  ipsec_policy = netbox_ipsec_policy.test.id
  description  = "Updated description"
}
`, ikePolicyName, ipsecPolicyName, name)
}

func TestAccIPSecProfileResource_import(t *testing.T) {
	// Generate unique name to avoid conflicts between test runs
	name := testutil.RandomName("tf-test-ipsec-profile")
	ikePolicyName := testutil.RandomName("tf-test-ike-policy-for-profile")
	ipsecPolicyName := testutil.RandomName("tf-test-ipsec-policy-for-profile")

	// Register cleanup to ensure resources are deleted even if test fails
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPSecProfileCleanup(name)
	cleanup.RegisterIKEPolicyCleanup(ikePolicyName)
	cleanup.RegisterIPSecPolicyCleanup(ipsecPolicyName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIPSecProfileDestroy,
			testutil.CheckIKEPolicyDestroy,
			testutil.CheckIPSecPolicyDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecProfileResourceConfig_basic(name, ikePolicyName, ipsecPolicyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "id"),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "name", name),
					resource.TestCheckResourceAttr("netbox_ipsec_profile.test", "mode", "esp"),
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "ike_policy"),
					resource.TestCheckResourceAttrSet("netbox_ipsec_profile.test", "ipsec_policy"),
				),
			},
			{
				ResourceName:      "netbox_ipsec_profile.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
