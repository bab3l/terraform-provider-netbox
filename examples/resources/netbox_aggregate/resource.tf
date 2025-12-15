resource "netbox_rir" "test" {
  name = "Test RIR"
  slug = "test-rir"
}

resource "netbox_aggregate" "test" {
  prefix = "10.0.0.0/8"
  rir    = netbox_rir.test.id
}
