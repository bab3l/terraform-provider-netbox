# Example: Basic route target
resource "netbox_route_target" "example" {
  name = "65000:100"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "rt_type"
      value = "import-export"
    }
  ]

  tags = [
    "route-target"
  ]
}

# Example: Route target with description
resource "netbox_route_target" "export" {
  name        = "65000:200"
  description = "Export route target for customer VRF"
  comments    = "Used for BGP VPN export"

  # Partial custom fields management
  custom_fields = [
    {
      name  = "rt_type"
      value = "export-only"
    },
    {
      name  = "customer_id"
      value = "CUST-12345"
    }
  ]

  tags = [
    "export",
    "customer"
  ]
}

# Example: Route target with tenant association
resource "netbox_route_target" "tenant_rt" {
  name        = "65001:100"
  tenant      = netbox_tenant.example.slug
  description = "Route target for tenant VRF"

  # Partial custom fields management
  custom_fields = [
    {
      name  = "rt_type"
      value = "tenant-specific"
    },
    {
      name  = "isolation_level"
      value = "high"
    }
  ]

  tags = [
    "tenant",
    "isolated"
  ]
}
