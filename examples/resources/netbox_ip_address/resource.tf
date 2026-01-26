resource "netbox_ip_address" "nat_inside" {
  address = "10.0.0.254/24"
}

resource "netbox_ip_address" "test_v4" {
  address     = "10.0.0.1/24"
  status      = "active"
  nat_inside  = netbox_ip_address.nat_inside.id
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
  vrf         = netbox_vrf.test.id
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

# IP address assigned to a device interface
resource "netbox_site" "ip_example" {
  name   = "IP Example Site"
  slug   = "ip-example-site"
  status = "active"
}

resource "netbox_manufacturer" "ip_example" {
  name = "IP Example Manufacturer"
  slug = "ip-example-manufacturer"
}

resource "netbox_device_type" "ip_example" {
  model        = "IP Example Device Type"
  slug         = "ip-example-device-type"
  manufacturer = netbox_manufacturer.ip_example.id
  u_height     = 1
}

resource "netbox_device_role" "ip_example" {
  name  = "IP Example Role"
  slug  = "ip-example-role"
  color = "ff0000"
}

resource "netbox_device" "ip_example" {
  name        = "ip-example-device"
  device_type = netbox_device_type.ip_example.id
  role        = netbox_device_role.ip_example.id
  site        = netbox_site.ip_example.id
  status      = "active"
}

resource "netbox_interface" "ip_example" {
  device = netbox_device.ip_example.id
  name   = "eth0"
  type   = "1000base-t"
}

resource "netbox_ip_address" "assigned_device" {
  address              = "10.10.0.10/24"
  status               = "active"
  assigned_object_type = "dcim.interface"
  assigned_object_id   = netbox_interface.ip_example.id
  description          = "Assigned to device interface"
}

# IP address assigned to a VM interface
resource "netbox_cluster_type" "ip_vm" {
  name = "IP VM Cluster Type"
  slug = "ip-vm-cluster-type"
}

resource "netbox_cluster" "ip_vm" {
  name = "IP VM Cluster"
  type = netbox_cluster_type.ip_vm.id
}

resource "netbox_virtual_machine" "ip_vm" {
  name    = "ip-example-vm"
  cluster = netbox_cluster.ip_vm.id
  status  = "active"
}

resource "netbox_vm_interface" "ip_vm" {
  name            = "eth0"
  virtual_machine = netbox_virtual_machine.ip_vm.id
}

resource "netbox_ip_address" "assigned_vm" {
  address              = "192.0.2.10/24"
  status               = "active"
  assigned_object_type = "virtualization.vminterface"
  assigned_object_id   = netbox_vm_interface.ip_vm.id
  description          = "Assigned to VM interface"
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
