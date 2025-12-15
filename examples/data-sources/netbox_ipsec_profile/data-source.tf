data "netbox_ipsec_profile" "test" {
  name = "test-ipsec-profile"
}

output "example" {
  value = data.netbox_ipsec_profile.test.id
}
