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

func TestClusterTypeResource(t *testing.T) {
	t.Parallel()
	r := resources.NewClusterTypeResource()
	if r == nil {
		t.Fatal("Expected non-nil Cluster Type resource")
	}
}

func TestClusterTypeResourceSchema(t *testing.T) {
	t.Parallel()
	r := resources.NewClusterTypeResource()
	schemaRequest := fwresource.SchemaRequest{}
	schemaResponse := &fwresource.SchemaResponse{}
	r.Schema(context.Background(), schemaRequest, schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)
	}

	if schemaResponse.Schema.Attributes == nil {
		t.Fatal("Expected schema to have attributes")
	}

	requiredAttrs := []string{"name", "slug"}
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

	optionalAttrs := []string{"description"}
	for _, attr := range optionalAttrs {
		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
			t.Errorf("Expected optional attribute %s to exist in schema", attr)
		}
	}
}

func TestClusterTypeResourceMetadata(t *testing.T) {
	t.Parallel()
	r := resources.NewClusterTypeResource()
	metadataRequest := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	metadataResponse := &fwresource.MetadataResponse{}
	r.Metadata(context.Background(), metadataRequest, metadataResponse)
	expected := "netbox_cluster_type"

	if metadataResponse.TypeName != expected {
		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestClusterTypeResourceConfigure(t *testing.T) {
	t.Parallel()
	r := resources.NewClusterTypeResource().(*resources.ClusterTypeResource)
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

func TestAccClusterTypeResource_basic(t *testing.T) {
	name := testutil.RandomName("tf-test-cluster-type")
	slug := testutil.RandomSlug("tf-test-cluster-type")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},
		},
	})
}

func TestAccClusterTypeResource_full(t *testing.T) {
	name := testutil.RandomName("tf-test-cluster-type-full")
	slug := testutil.RandomSlug("tf-test-cluster-type-full")
	description := "Test cluster type with all fields"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_full(name, slug, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", description),
				),
			},
		},
	})
}

func TestAccClusterTypeResource_update(t *testing.T) {
	name := testutil.RandomName("tf-test-cluster-type-update")
	slug := testutil.RandomSlug("tf-test-cluster-type-update")
	updatedName := testutil.RandomName("tf-test-cluster-type-updated")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
				),
			},
			{
				Config: testAccClusterTypeResourceConfig_full(updatedName, slug, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "description", "Updated description"),
				),
			},
		},
	})
}

func testAccClusterTypeResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}

func testAccClusterTypeResourceConfig_full(name, slug, description string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_type" "test" {
  name        = %q
  slug        = %q
  description = %q
}
`, name, slug, description)
}

func TestAccClusterTypeResource_import(t *testing.T) {
	name := testutil.RandomName("tf-test-cluster-type-import")
	slug := testutil.RandomSlug("tf-test-cluster-type-import")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterTypeCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterTypeResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_type.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_type.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_cluster_type.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
