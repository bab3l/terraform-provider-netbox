resource "netbox_route_target" "import" {
  name = "65000:100"
}

resource "netbox_route_target" "export" {
  name = "65000:200"
}

resource "netbox_vrf" "test" {
  name           = "Test VRF"
  rd             = "65000:1"
  description    = "Customer VRF for multi-tenant network"
  comments       = "Isolated routing table for customer traffic"
  import_targets = [netbox_route_target.import.id]
  export_targets = [netbox_route_target.export.id]

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

# Optional: seed owned custom fields during import
import {
  to = netbox_vrf.test
  id = "123"

  identity = {
    custom_fields = [
      "customer_name:text",
      "vrf_type:text",
      "import_policy:text",
      "export_policy:text",
    ]
  }
}
