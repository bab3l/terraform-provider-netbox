data "netbox_l2vpn" "by_id" {
  id = "1"
}

data "netbox_l2vpn" "by_name" {
  name = "Corporate-L2VPN"
}

data "netbox_l2vpn" "by_slug" {
  slug = "corporate-l2vpn"
}

output "l2vpn_id" {
  value = data.netbox_l2vpn.by_id.id
}

output "l2vpn_name" {
  value = data.netbox_l2vpn.by_name.name
}

output "l2vpn_type" {
  value = data.netbox_l2vpn.by_slug.type
}
