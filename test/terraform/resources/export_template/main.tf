# Export Template Resource Integration Test
# Tests the netbox_export_template resource with basic and complete configurations

terraform {
  required_version = ">= 1.0"
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Test 1: Basic export template with only required fields
resource "netbox_export_template" "basic" {
  name          = "Test Export Template Basic"
  object_types  = ["dcim.site"]
  template_code = "{% for site in queryset %}{{ site.name }}\n{% endfor %}"
}

# Test 2: Export template for devices
resource "netbox_export_template" "devices" {
  name          = "Test Export Template Devices"
  object_types  = ["dcim.device"]
  template_code = "name,site,status\n{% for device in queryset %}{{ device.name }},{{ device.site.name }},{{ device.status }}\n{% endfor %}"
  description   = "Export template for device inventory"
}

# Test 3: Export template with multiple object types
resource "netbox_export_template" "multi_type" {
  name          = "Test Export Template Multi Type"
  object_types  = ["dcim.site", "dcim.location"]
  template_code = "{% for obj in queryset %}{{ obj.name }},{{ obj.slug }}\n{% endfor %}"
  description   = "Export template for sites and locations"
}

# Test 4: Complete export template with all optional fields
resource "netbox_export_template" "complete" {
  name           = "Test Export Template Complete"
  object_types   = ["ipam.prefix"]
  template_code  = "prefix,vrf,status\n{% for prefix in queryset %}{{ prefix.prefix }},{{ prefix.vrf.name|default:'Global' }},{{ prefix.status }}\n{% endfor %}"
  description    = "Complete export template for IP prefixes"
  mime_type      = "text/csv"
  file_extension = "csv"
  as_attachment  = true
}

# Test 5: Export template as JSON
resource "netbox_export_template" "json" {
  name           = "Test Export Template JSON"
  object_types   = ["dcim.site"]
  template_code  = "[\n{% for site in queryset %}  {\"name\": \"{{ site.name }}\", \"slug\": \"{{ site.slug }}\"}{% if not loop.last %},{% endif %}\n{% endfor %}]"
  description    = "Export template outputting JSON"
  mime_type      = "application/json"
  file_extension = "json"
  as_attachment  = false
}
