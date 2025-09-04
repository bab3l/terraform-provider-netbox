package testutil

import (
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
}
