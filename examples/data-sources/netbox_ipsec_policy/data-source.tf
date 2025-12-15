data "netbox_ipsec_policy" "test" {
  name = "test-ipsec-policy"
}

output "example" {
  value = data.netbox_ipsec_policy.test.id
}
