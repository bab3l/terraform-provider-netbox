# Example: Basic config context
# Note: netbox_config_context resource only supports tags, not custom_fields
resource "netbox_config_context" "basic" {
  name = "basic-config"
  data = jsonencode({
    dns_servers = ["8.8.8.8", "8.8.4.4"]
    ntp_servers = ["time.google.com"]
  })
}

# Example: Config context with assignment criteria
resource "netbox_config_context" "site_specific" {
  name        = "dc1-config"
  description = "Configuration for DC1 site"
  weight      = 1500
  is_active   = true
  data = jsonencode({
    syslog_server  = "10.1.0.10"
    snmp_community = "public"
    timezone       = "America/New_York"
  })
  sites = [netbox_site.dc1.id]
}

# Example: Config context for specific device roles
resource "netbox_config_context" "router_config" {
  name        = "router-defaults"
  description = "Default configuration for all routers"
  weight      = 1000
  data = jsonencode({
    routing = {
      ospf = {
        enabled = true
        area    = "0.0.0.0"
      }
      bgp = {
        enabled = false
      }
    }
    interfaces = {
      mtu = 9000
    }
  })
  roles = [netbox_device_role.router.id]
}

# Example: Config context with multiple assignment criteria
resource "netbox_config_context" "production_servers" {
  name        = "production-server-config"
  description = "Configuration for all production servers"
  weight      = 2000
  data = jsonencode({
    monitoring = {
      enabled     = true
      interval    = 60
      alert_email = "ops@example.com"
    }
    backup = {
      enabled   = true
      schedule  = "0 2 * * *"
      retention = 30
    }
  })
  sites         = [netbox_site.dc1.id, netbox_site.dc2.id]
  roles         = [netbox_device_role.server.id]
  tenant_groups = [netbox_tenant_group.production.id]
}

# Example: Config context with tag-based assignment
resource "netbox_config_context" "high_security" {
  name        = "high-security-config"
  description = "Security hardening configuration for tagged devices"
  weight      = 3000
  data = jsonencode({
    security = {
      firewall_enabled = true
      ssh_port         = 2222
      password_policy  = "complex"
      audit_logging    = true
    }
  })
  tags = ["security-hardened", "pci-compliant"]
}
