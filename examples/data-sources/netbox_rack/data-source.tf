# Lookup by ID
data "netbox_rack" "by_id" {
  id = "123"
}

output "by_id" {
  value = data.netbox_rack.by_id.name
}

# Lookup by name
data "netbox_rack" "by_name" {
  name = "RACK-A1"
}

output "by_name" {
  value = data.netbox_rack.by_name.site
}
