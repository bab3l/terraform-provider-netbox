resource "netbox_vrf" "test" {
  name        = "Test VRF"
  rd          = "65000:1"
  description = "Customer VRF for multi-tenant network"
  comments    = "Isolated routing table for customer traffic"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "customer_name"
      value = "Acme Corporation"
    },
    {
      name  = "vrf_type"
      value = "customer"
    },
    {
      name  = "import_policy"
      value = "65000:100"
    },
    {
      name  = "export_policy"
      value = "65000:200"
    }
  ]

  tags = [
    "customer-vrf",
    "production"
  ]
}
