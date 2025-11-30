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

func TestTenantGroupResource(t *testing.T) {
t.Parallel()

r := resources.NewTenantGroupResource()
if r == nil {
t.Fatal("Expected non-nil tenant group resource")
}
}

func TestTenantGroupResourceSchema(t *testing.T) {
t.Parallel()

r := resources.NewTenantGroupResource()
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

optionalAttrs := []string{"parent", "description", "tags", "custom_fields"}
for _, attr := range optionalAttrs {
if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
t.Errorf("Expected optional attribute %s to exist in schema", attr)
}
}

computedAttrs := []string{"id"}
for _, attr := range computedAttrs {
if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {
t.Errorf("Expected computed attribute %s to exist in schema", attr)
}
}
}

func TestTenantGroupResourceMetadata(t *testing.T) {
t.Parallel()

r := resources.NewTenantGroupResource()
metadataRequest := fwresource.MetadataRequest{
ProviderTypeName: "netbox",
}
metadataResponse := &fwresource.MetadataResponse{}

r.Metadata(context.Background(), metadataRequest, metadataResponse)

expected := "netbox_tenant_group"
if metadataResponse.TypeName != expected {
t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
}
}

func TestTenantGroupResourceConfigure(t *testing.T) {
t.Parallel()

r := resources.NewTenantGroupResource().(*resources.TenantGroupResource)

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

if r.GetClient() != client {
t.Error("Expected client to be set")
}

configureRequest.ProviderData = "invalid"
configureResponse = &fwresource.ConfigureResponse{}

r.Configure(context.Background(), configureRequest, configureResponse)

if !configureResponse.Diagnostics.HasError() {
t.Error("Expected error with incorrect provider data")
}
}

func TestAccTenantGroupResource_basic(t *testing.T) {
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

resource "netbox_tenant_group" "test" {
  name = "Test Tenant Group"
  slug = "test-tenant-group"
}
`,
Check: resource.ComposeTestCheckFunc(
resource.TestCheckResourceAttrSet("netbox_tenant_group.test", "id"),
resource.TestCheckResourceAttr("netbox_tenant_group.test", "name", "Test Tenant Group"),
resource.TestCheckResourceAttr("netbox_tenant_group.test", "slug", "test-tenant-group"),
),
},
},
})
}
