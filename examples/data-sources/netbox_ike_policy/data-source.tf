data "netbox_ike_policy" "test" {
  name = "test-ike-policy"
}

output "example" {
  value = data.netbox_ike_policy.test.id
}
