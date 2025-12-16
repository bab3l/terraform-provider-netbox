# Example: Basic route target
resource "netbox_route_target" "example" {
  name = "65000:100"
}

# Example: Route target with description
resource "netbox_route_target" "export" {
  name        = "65000:200"
  description = "Export route target for customer VRF"
  comments    = "Used for BGP VPN export"
}

# Example: Route target with tenant association
resource "netbox_route_target" "tenant_rt" {
  name        = "65001:100"
  tenant      = netbox_tenant.example.slug
  description = "Route target for tenant VRF"
}
