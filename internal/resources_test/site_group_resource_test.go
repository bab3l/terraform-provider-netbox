package resources_test

import (
"context"
"testing"

"github.com/bab3l/go-netbox"
"github.com/bab3l/terraform-provider-netbox/internal/provider"
"github.com/bab3l/terraform-provider-netbox/internal/resources"
"github.com/hashicorp/terraform-plugin-framework/providerserver"
fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
"github.com/hashicorp/terraform-plugin-go/tfprotov6"
"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestSiteGroupResource(t *testing.T) {
r := resources.NewSiteGroupResource()
if r == nil {
t.Fatal("Site group resource should not be nil")
}
}

func TestSiteGroupResourceSchema(t *testing.T) {
ctx := context.Background()
r := resources.NewSiteGroupResource()

schemaReq := fwresource.SchemaRequest{}
schemaResp := &fwresource.SchemaResponse{}

r.Schema(ctx, schemaReq, schemaResp)

if schemaResp.Diagnostics.HasError() {
t.Errorf("Schema should not have errors: %v", schemaResp.Diagnostics.Errors())
}

requiredAttributes := []string{"id", "name", "slug"}
for _, attr := range requiredAttributes {
if _, exists := schemaResp.Schema.Attributes[attr]; !exists {
t.Errorf("Site group resource schema should include %s attribute", attr)
}
}

optionalAttributes := []string{"parent", "description", "tags", "custom_fields"}
for _, attr := range optionalAttributes {
if _, exists := schemaResp.Schema.Attributes[attr]; !exists {
t.Errorf("Site group resource schema should include %s attribute", attr)
}
}
}

func TestSiteGroupResourceMetadata(t *testing.T) {
ctx := context.Background()
r := resources.NewSiteGroupResource()

metadataReq := fwresource.MetadataRequest{
ProviderTypeName: "netbox",
}
metadataResp := &fwresource.MetadataResponse{}

r.Metadata(ctx, metadataReq, metadataResp)

expectedTypeName := "netbox_site_group"
if metadataResp.TypeName != expectedTypeName {
t.Errorf("Expected type name %s, got %s", expectedTypeName, metadataResp.TypeName)
}
}

func TestSiteGroupResourceConfigure(t *testing.T) {
ctx := context.Background()
r := resources.NewSiteGroupResource().(*resources.SiteGroupResource)

configureReq := fwresource.ConfigureRequest{
ProviderData: nil,
}
configureResp := &fwresource.ConfigureResponse{}

r.Configure(ctx, configureReq, configureResp)

if configureResp.Diagnostics.HasError() {
t.Error("Configure should not error with nil provider data")
}

client := &netbox.APIClient{}
configureReq.ProviderData = client
configureResp = &fwresource.ConfigureResponse{}

r.Configure(ctx, configureReq, configureResp)

if configureResp.Diagnostics.HasError() {
t.Errorf("Configure should not error with correct provider data: %v", configureResp.Diagnostics.Errors())
}

configureReq.ProviderData = "invalid"
configureResp = &fwresource.ConfigureResponse{}

r.Configure(ctx, configureReq, configureResp)

if !configureResp.Diagnostics.HasError() {
t.Error("Configure should error with incorrect provider data type")
}
}

func TestAccSiteGroupResource_basic(t *testing.T) {
resource.Test(t, resource.TestCase{
ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
},
Steps: []resource.TestStep{
{
Config: `
terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
      version = ">= 0.1.0"
    }
  }
}

provider "netbox" {}

resource "netbox_site_group" "test" {
  name = "Test Site Group"
  slug = "test-site-group"
}
`,
Check: resource.ComposeTestCheckFunc(
resource.TestCheckResourceAttrSet("netbox_site_group.test", "id"),
resource.TestCheckResourceAttr("netbox_site_group.test", "name", "Test Site Group"),
resource.TestCheckResourceAttr("netbox_site_group.test", "slug", "test-site-group"),
),
},
},
})
}
