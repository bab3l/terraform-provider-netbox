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

output "site_status" {
  value = data.netbox_site.by_id.status
}

output "site_region" {
  value = data.netbox_site.by_slug.region
}

output "site_description" {
  value = data.netbox_site.by_name.description
}
