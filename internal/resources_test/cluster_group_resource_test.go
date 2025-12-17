package resources_test

import (
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestClusterGroupResource_Metadata(t *testing.T) {
	r := resources.NewClusterGroupResource()
	req := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &fwresource.MetadataResponse{}
	r.Metadata(nil, req, resp)

	if resp.TypeName != "netbox_cluster_group" {
		t.Errorf("expected TypeName 'netbox_cluster_group', got '%s'", resp.TypeName)
	}
}

func TestClusterGroupResource_Schema(t *testing.T) {
	r := resources.NewClusterGroupResource()
	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	r.Schema(nil, req, resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("expected schema attributes, got nil")
	}

	requiredAttrs := []string{"name", "slug"}
	optionalAttrs := []string{"description", "tags", "custom_fields"}
	computedAttrs := []string{"id"}

	for _, attr := range requiredAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected required attribute '%s' in schema", attr)
		}
	}

	for _, attr := range optionalAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected optional attribute '%s' in schema", attr)
		}
	}

	for _, attr := range computedAttrs {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("expected computed attribute '%s' in schema", attr)
		}
	}
}

func TestClusterGroupResource_SchemaDescription(t *testing.T) {
	r := resources.NewClusterGroupResource()
	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	r.Schema(nil, req, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Error("expected schema to have a description")
	}
}

func TestClusterGroupResource_Configure(t *testing.T) {
	r := resources.NewClusterGroupResource().(*resources.ClusterGroupResource)
	req := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &fwresource.ConfigureResponse{}
	r.Configure(nil, req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no errors with nil provider data, got: %v", resp.Diagnostics)
	}
}

func TestAccClusterGroupResource_basic(t *testing.T) {
	name := testutil.RandomName("tf-test-cluster-group")
	slug := testutil.RandomSlug("tf-test-cluster-group")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterClusterGroupCleanup(slug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckClusterGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccClusterGroupResourceConfig_basic(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_cluster_group.test", "id"),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_cluster_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_cluster_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccClusterGroupResourceConfig_basic(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_cluster_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}
