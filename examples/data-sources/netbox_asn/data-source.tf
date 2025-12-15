data "netbox_asn" "test" {
  asn = 65001
}

output "example" {
  value = data.netbox_asn.test.id
}
