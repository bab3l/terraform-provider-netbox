resource "netbox_ip_address" "test_v4" {
  address     = "10.0.0.1/24"
  status      = "active"
  dns_name    = "test.example.com"
  description = "Primary web server IP"

  # Partial custom fields management
  # Only specified custom fields are managed, others preserved
  custom_fields = [
    {
      name  = "owner_team"
      value = "web-ops"
    },
    {
      name  = "purpose"
      value = "load-balancer-vip"
    },
    {
      name  = "monitoring_enabled"
      value = "true"
    }
  ]

  tags = [
    "production",
    "web-tier"
  ]
}

resource "netbox_ip_address" "test_v6" {
  address     = "2001:db8::1/64"
  status      = "active"
  description = "IPv6 primary address"

  # Partial custom fields management
  custom_fields = [
    {
      name  = "ipv6_deployment_phase"
      value = "production"
    }
  ]

  tags = [
    "ipv6",
    "production"
  ]
}

resource "netbox_vrf" "test" {
  name = "Test VRF"
  rd   = "65000:1"
}

resource "netbox_ip_address" "test_vrf" {
  address     = "192.168.1.1/24"
  vrf         = netbox_vrf.test.name
  status      = "active"
  description = "VRF gateway address"

  # Partial custom fields management
  custom_fields = [
    {
      name  = "gateway_role"
      value = "default"
    },
    {
      name  = "vrf_priority"
      value = "high"
    }
  ]

  tags = [
    "gateway",
    "vrf-test"
  ]
}

# Optional: seed owned custom fields during import
import {
  to = netbox_ip_address.test_v4
  id = "123"

  identity = {
    custom_fields = [
      "owner_team:text",
      "purpose:text",
      "monitoring_enabled:boolean",
    ]
  }
}

# Optional: seed owned custom fields during import
import {
  to = netbox_ip_address.test_v6
  id = "124"

  identity = {
    custom_fields = [
      "ipv6_deployment_phase:text",
    ]
  }
}

# Optional: seed owned custom fields during import
import {
  to = netbox_ip_address.test_vrf
  id = "125"

  identity = {
    custom_fields = [
      "gateway_role:text",
      "vrf_priority:text",
    ]
  }
}
