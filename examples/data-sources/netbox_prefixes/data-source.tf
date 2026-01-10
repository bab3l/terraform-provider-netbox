# Example for the plural/query prefixes data source.
#
# Notes:
# - Multiple `filter` blocks are ANDed together.
# - Multiple values inside one filter block are ORed together.
# - The datasource returns `ids`, `cidrs`, and `prefixes` (list of `{id,prefix}` objects).

resource "netbox_prefix" "example" {
  prefix = "10.10.0.0/24"
  status = "active"
}

data "netbox_prefixes" "by_prefix" {
  filter {
    name   = "prefix"
    values = [netbox_prefix.example.prefix]
  }
}

output "prefix_ids" {
  value = data.netbox_prefixes.by_prefix.ids
}

output "prefix_cidrs" {
  value = data.netbox_prefixes.by_prefix.cidrs
}

output "prefix_objects" {
  value = data.netbox_prefixes.by_prefix.prefixes
}
