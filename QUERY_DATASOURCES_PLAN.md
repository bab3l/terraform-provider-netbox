# Query (Plural) Data Sources Plan

## Goal
Implement AWS-style “query” data sources for NetBox with a consistent UX:

- Per-resource plural data sources (e.g., `netbox_devices`, `netbox_virtual_machines`)
- AWS-like filter block:
  ```hcl
  filter {
    name   = "<filter-key>"
    values = ["v1", "v2"]
  }
  ```
  - Multiple `filter` blocks are ANDed.
  - Values within a `filter` block are ORed.

These plural data sources should return stable identifiers (at minimum `ids`) and optionally convenience outputs like `names`.

## Constraints / Reality of the NetBox Client
The provider uses the generated `go-netbox` OpenAPI client which exposes *typed* filter setters (e.g., `.Name([]string{...})`, `.NameIc([]string{...})`, etc.).

This makes a fully-generic “arbitrary query param passthrough” difficult without bypassing the typed request builder.

**Plan:** keep the AWS-style *interface* (`filter { name, values }`) but initially support a curated subset of filter keys per data source (mapping each supported key to a typed `go-netbox` request method). Unsupported filter names return a clear diagnostic.

We can expand the supported filter set over time and/or later introduce an advanced “raw query params” escape hatch (only if needed).

## Proposed Data Sources
### Phase 1 (high value)
- `netbox_devices` (DCIM)
- `netbox_virtual_machines` (Virtualization)

### Phase 2 (also high value)
- `netbox_ip_addresses` (IPAM)
- `netbox_prefixes` (IPAM)
- `netbox_interfaces` (DCIM)

## Trial Implementation Target
Implement `netbox_devices` first as the template for the pattern.

## Progress (current)
- Implemented plural/query datasources:
  - `netbox_devices`
  - `netbox_virtual_machines`
- Shared helper added: `ExpandQueryFilters` for normalized/merged filters.
- Custom field filtering supported for deterministic discovery:
  - `custom_field` (existence)
  - `custom_field_value` (value match; `field=value`)
- Outputs implemented:
  - `ids`, `names`
  - computed object list output (minimal `{id, name}`): `devices` / `virtual_machines`
- Tests implemented and passing (including tag + multi-filter coverage):
  - unit tests under `internal/datasources_unit_tests/`
  - acceptance tests under `internal/datasources_acceptance_tests/`
  - customfields acceptance tests under `internal/datasources_acceptance_tests_customfields/`

## Documentation Plan (required)
This repo uses `terraform-plugin-docs`, so documentation needs to be kept in sync with schema changes.

For each new plural datasource:
- Add examples under `examples/data-sources/<data source name>/` (so docs generation can include usable snippets).
- Regenerate docs once after all implementation + example generation work is complete (avoid running generation repeatedly mid-stream):
  ```
  go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-dir=. --rendered-website-dir=docs
  ```

For `netbox_devices`, examples should cover:
- Exact match (`name`)
- Case-insensitive match (`name__ic`)
- OR semantics (multiple values in one filter)
- AND semantics (multiple filter blocks)
- Tag filtering (`tag`)
- Custom field existence/value (`custom_field`, `custom_field_value`)
- Free-text search (`q`) (NetBox server behavior varies by version/config)

### `netbox_devices` schema (trial)
Inputs:
- `filter` (optional) repeated block:
  - `name` (required)
  - `values` (required list)

Outputs:
- `ids` (computed list of string IDs)
- `names` (computed list of names; best-effort)

Behavior:
- If no filters are provided, return a diagnostic (to avoid accidental “list the world”).
  - Possible future: `allow_unsafe_query = true` to override.

Supported filter names (initial):
- `name` -> `.Name(values)`
- `name__ic` -> `.NameIc(values)`
- `serial` -> `.Serial(values)`
- `status` -> `.Status(values)` (if available in go-netbox)
- `site` -> `.Site(values)` or `.SiteId(...)` depending on client
- `tag` -> `.Tag(values)` (if available)
- `custom_field` -> client-side filter on returned `custom_fields` (requires API includes `custom_fields` in results)
- `custom_field_value` -> client-side filter on returned `custom_fields` using `field=value` values

Additional implementation notes:
- `q` can be supported as NetBox free-text search:
  - Mapping: `q` -> `.Q(value)`
  - Constraint: exactly one value
- `custom_field` / `custom_field_value` are implemented client-side after fetching results, since the generated client does not expose typed `cf_*` query params.
- Pagination is required for list endpoints; iterate pages until `count` results are collected.

(Exact availability depends on generated client; we will align to what `ApiDcimDevicesListRequest` provides.)

## Implementation Notes
- Use terraform-plugin-framework.
- Implement as a normal data source under `internal/datasources/`.
- Register in provider `DataSources()`.
- Add an acceptance test that:
  - creates a device resource with a unique name
  - queries `data "netbox_devices"` with `filter { name = "name" values = [<device-name>] }`
  - asserts the returned `ids` includes exactly one item and matches the created resource ID.

## Open Questions
- Should the plural data sources also expose a computed list of objects (like AWS Framework data sources often do), e.g. `devices = [{id,name,...}]`?
  - Decision: yes. Expose a computed list of objects per plural datasource (minimal shape: `{id, name}`), in addition to `ids` and `names`.
- Should we allow “no filters” when NetBox instance is small?
  - For trial, default to safe behavior: require at least one filter.

## Helper Functions / Reuse Notes
To avoid filter parsing duplication across plural datasources:

- Introduce a shared filter model and expansion helper:
  - `QueryFilterModel`
  - `ExpandQueryFilters(ctx, filters)`
- Normalization behavior:
  - filter names are trimmed + lowercased
  - filter values are trimmed; empty values are discarded
  - repeated filter blocks for the same name are merged

Each plural datasource still needs a small mapping layer from normalized filter names to typed go-netbox request builder methods.
