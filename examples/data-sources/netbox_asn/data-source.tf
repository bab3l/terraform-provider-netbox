data "netbox_asn" "by_asn" {
  asn = 65001
}

data "netbox_asn" "by_id" {
  id = 123
}

output "by_asn" {
  value = data.netbox_asn.by_asn.id
}

output "by_id" {
  value = data.netbox_asn.by_id.asn
}
