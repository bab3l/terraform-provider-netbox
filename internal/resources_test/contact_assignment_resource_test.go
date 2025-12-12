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

func TestContactAssignmentResource(t *testing.T) {
	t.Parallel()

	r := resources.NewContactAssignmentResource()
	if r == nil {
		t.Fatal("Expected non-nil contact assignment resource")
	}
}

func TestContactAssignmentResourceSchema(t *testing.T) {
	t.Parallel()

	r := resources.NewContactAssignmentResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	// Check required attributes
	requiredAttrs := []string{"object_type", "object_id", "contact_id"}
	for _, attr := range requiredAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected required attribute %s to exist in schema", attr)
		}
	}

	// Check computed attributes
	computedAttrs := []string{"id"}
	for _, attr := range computedAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected computed attribute %s to exist in schema", attr)
		}
	}

	// Check optional attributes
	optionalAttrs := []string{"role_id", "priority", "tags", "custom_fields"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestContactAssignmentResourceMetadata(t *testing.T) {
	t.Parallel()

	r := resources.NewContactAssignmentResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_contact_assignment"
	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestContactAssignmentResourceConfigure(t *testing.T) {
	t.Parallel()

	r := resources.NewContactAssignmentResource().(*resources.ContactAssignmentResource)

	// Test with nil provider data (should not error)
	configureRequest := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with correct provider data type
	client := &netbox.APIClient{}
	configureRequest = fwresource.ConfigureRequest{
		ProviderData: client,
	}
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {
		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)
	}

	// Test with wrong provider data type
	configureRequest = fwresource.ConfigureRequest{
		ProviderData: "wrong type",
	}
	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {
		t.Error("Expected error with wrong provider data type")
	}
}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables
func TestAccContactAssignmentResource_basic(t *testing.T) {
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-assign")
	randomSlug := testutil.RandomSlug("test-ca")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccContactAssignmentResourceBasic(randomName, randomSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "contact_id"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "object_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "netbox_contact_assignment.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccContactAssignmentResource_withRole(t *testing.T) {
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-assign")
	randomSlug := testutil.RandomSlug("test-ca")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			// Create with role and priority
			{
				Config: testAccContactAssignmentResourceWithRole(randomName, randomSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "primary"),
					resource.TestCheckResourceAttrSet("netbox_contact_assignment.test", "role_id"),
				),
			},
		},
	})
}

func TestAccContactAssignmentResource_update(t *testing.T) {
	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("test-contact-assign")
	randomSlug := testutil.RandomSlug("test-ca")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			// Create without priority
			{
				Config: testAccContactAssignmentResourceBasic(randomName, randomSlug),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "object_type", "dcim.site"),
				),
			},
			// Update to add priority
			{
				Config: testAccContactAssignmentResourceWithPriority(randomName, randomSlug, "secondary"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_assignment.test", "priority", "secondary"),
				),
			},
		},
	})
}

func testAccContactAssignmentResourceBasic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "test@example.com"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
}
`, name, slug, name, name, slug)
}

func testAccContactAssignmentResourceWithRole(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "test@example.com"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "primary"
}
`, name, slug, name, name, slug)
}

func testAccContactAssignmentResourceWithPriority(name, slug, priority string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name   = "%s-site"
  slug   = "%s-site"
  status = "active"
}

resource "netbox_contact" "test" {
  name  = "%s-contact"
  email = "test@example.com"
}

resource "netbox_contact_role" "test" {
  name = "%s-role"
  slug = "%s-role"
}

resource "netbox_contact_assignment" "test" {
  object_type = "dcim.site"
  object_id   = netbox_site.test.id
  contact_id  = netbox_contact.test.id
  role_id     = netbox_contact_role.test.id
  priority    = "%s"
}
`, name, slug, name, name, slug, priority)
}
