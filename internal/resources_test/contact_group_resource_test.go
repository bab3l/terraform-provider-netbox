package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestContactGroupResource_Metadata(t *testing.T) {
	r := resources.NewContactGroupResource()
	req := fwresource.MetadataRequest{
		ProviderTypeName: "netbox",
	}
	resp := &fwresource.MetadataResponse{}
	r.Metadata(nil, req, resp)

	if resp.TypeName != "netbox_contact_group" {
		t.Errorf("expected TypeName 'netbox_contact_group', got '%s'", resp.TypeName)
	}
}

func TestContactGroupResource_Schema(t *testing.T) {
	r := resources.NewContactGroupResource()
	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	r.Schema(nil, req, resp)

	if resp.Schema.Attributes == nil {
		t.Fatal("expected schema attributes, got nil")
	}

	requiredAttrs := []string{"name", "slug"}
	optionalAttrs := []string{"parent", "description", "tags", "custom_fields"}
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

func TestContactGroupResource_SchemaDescription(t *testing.T) {
	r := resources.NewContactGroupResource()
	req := fwresource.SchemaRequest{}
	resp := &fwresource.SchemaResponse{}
	r.Schema(nil, req, resp)

	if resp.Schema.MarkdownDescription == "" {
		t.Error("expected schema to have a description")
	}
}

func TestContactGroupResource_Configure(t *testing.T) {
	r := resources.NewContactGroupResource().(*resources.ContactGroupResource)
	req := fwresource.ConfigureRequest{
		ProviderData: nil,
	}
	resp := &fwresource.ConfigureResponse{}
	r.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Errorf("expected no errors with nil provider data, got: %v", resp.Diagnostics)
	}
}

func TestAccContactGroupResource_basic(t *testing.T) {
	name := testutil.RandomName("test-contact-group")
	slug := testutil.GenerateSlug(name)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccContactGroupResourceConfig(name, slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name),
					resource.TestCheckResourceAttr("netbox_contact_group.test", "slug", slug),
				),
			},
			{
				ResourceName:      "netbox_contact_group.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccContactGroupResourceConfig(name+"-updated", slug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_contact_group.test", "name", name+"-updated"),
				),
			},
		},
	})
}

func testAccContactGroupResourceConfig(name, slug string) string {
	return fmt.Sprintf(`
resource "netbox_contact_group" "test" {
  name = %q
  slug = %q
}
`, name, slug)
}
