# Export Template Data Source Integration Test
# Tests the netbox_export_template data source for looking up existing export templates

terraform {
  required_providers {
    netbox = {
      source = "bab3l/netbox"
    }
  }
}

provider "netbox" {
  # Uses NETBOX_SERVER_URL and NETBOX_API_TOKEN environment variables
}

# Create export templates to look up
resource "netbox_export_template" "test_csv" {
  name           = "Test Export Template CSV DS"
  object_types   = ["dcim.site"]
  template_code  = "name,slug\n{% for site in queryset %}{{ site.name }},{{ site.slug }}\n{% endfor %}"
  description    = "CSV export template for data source testing"
  mime_type      = "text/csv"
  file_extension = "csv"
}

resource "netbox_export_template" "test_json" {
  name           = "Test Export Template JSON DS"
  object_types   = ["dcim.device"]
  template_code  = "[{% for device in queryset %}{\"name\": \"{{ device.name }}\"}{% if not loop.last %},{% endif %}{% endfor %}]"
  description    = "JSON export template for data source testing"
  mime_type      = "application/json"
  file_extension = "json"
}

# Look up export template by ID
data "netbox_export_template" "by_id" {
  id = netbox_export_template.test_csv.id
}

# Outputs for verification
output "by_id_name" {
  value = data.netbox_export_template.by_id.name
}

output "by_id_object_types" {
  value = data.netbox_export_template.by_id.object_types
}

output "by_id_description" {
  value = data.netbox_export_template.by_id.description
}

output "by_id_mime_type" {
  value = data.netbox_export_template.by_id.mime_type
}

output "by_id_file_extension" {
  value = data.netbox_export_template.by_id.file_extension
}

# Validation outputs
output "id_lookup_matches" {
  value = data.netbox_export_template.by_id.id == netbox_export_template.test_csv.id
}

output "name_matches" {
  value = data.netbox_export_template.by_id.name == netbox_export_template.test_csv.name
}

output "all_lookups_valid" {
  value = data.netbox_export_template.by_id.id == netbox_export_template.test_csv.id
}
