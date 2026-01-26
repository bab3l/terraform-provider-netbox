resource "netbox_site" "example" {
  name   = "Example Site"
  slug   = "example-site"
  status = "active"
}

resource "netbox_rir" "example" {
  name = "Example RIR"
  slug = "example-rir"
}

resource "netbox_asn" "example" {
  asn = 64512
  rir = netbox_rir.example.id
}

resource "netbox_site_asn_assignment" "example" {
  site = netbox_site.example.id
  asn  = netbox_asn.example.id
}
