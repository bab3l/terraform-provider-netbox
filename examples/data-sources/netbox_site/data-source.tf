# Look up site by ID
data "netbox_site" "by_id" {
  id = "1"
}

# Look up site by slug
data "netbox_site" "by_slug" {
  slug = "dc-east"
}

# Look up site by name
data "netbox_site" "by_name" {
  name = "Data Center East"
}

# Use site data in other resources
output "site_name" {
  value = data.netbox_site.by_id.name
}

output "site_status" {
  value = data.netbox_site.by_id.status
}

output "site_region" {
  value = data.netbox_site.by_slug.region
}

output "site_tenant" {
  value = data.netbox_site.by_id.tenant
}

output "site_facility" {
  value = data.netbox_site.by_id.facility
}

output "site_time_zone" {
  value = data.netbox_site.by_name.time_zone
}

# Access all custom fields
output "site_custom_fields" {
  value       = data.netbox_site.by_id.custom_fields
  description = "All custom fields defined in NetBox for this site"
}

# Access specific custom fields by name
output "site_building_code" {
  value       = try([for cf in data.netbox_site.by_id.custom_fields : cf.value if cf.name == "building_code"][0], null)
  description = "Example: accessing a text custom field"
}

output "site_datacenter_tier" {
  value       = try([for cf in data.netbox_site.by_id.custom_fields : cf.value if cf.name == "datacenter_tier"][0], null)
  description = "Example: accessing a select custom field"
}
