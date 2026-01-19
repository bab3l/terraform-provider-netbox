resource "netbox_site" "test" {
  name   = "Test Site"
  slug   = "test-site"
  status = "active"
}

resource "netbox_vlan" "test" {
  vid    = 100
  name   = "Test VLAN"
  site   = netbox_site.test.id
  status = "active"
}

resource "netbox_prefix" "test" {
  prefix      = "10.0.0.0/24"
  site        = netbox_site.test.slug
  vlan        = netbox_vlan.test.name
  status      = "active"
  description = "Primary datacenter subnet"
  comments    = "Main server network segment"

  # Partial custom fields management
  # Only specified custom fields are managed, others in NetBox preserved
  custom_fields = [
    {
      name  = "subnet_purpose"
      value = "servers"
    },
    {
      name  = "vlan_id"
      value = "100"
    },
    {
      name  = "dhcp_enabled"
      value = "false"
    },
    {
      name  = "gateway_ip"
      value = "10.0.0.1"
    }
  ]

  tags = [
    "production",
    "datacenter"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_prefix.test
  id = "123"

  identity = {
    custom_fields = [
      "subnet_purpose:text",
      "vlan_id:integer",
      "dhcp_enabled:boolean",
      "gateway_ip:text",
    ]
  }
}
