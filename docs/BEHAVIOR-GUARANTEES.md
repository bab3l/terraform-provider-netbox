# Terraform Provider NetBox 1.0 Behavior Guarantees

This guide documents the user-visible behavior the provider intends to keep stable starting with `1.0.0`.

Scope:

- provider release line: `1.0.x`
- tested NetBox version: `v4.1.11`
- audience: users who want predictable plan, state, and import behavior

## Reference fields

Many resource attributes reference other NetBox objects such as tenants, sites, device roles, platforms, racks, and similar related records.

### Accepted input forms

For reference attributes that use the shared lookup helpers, the provider accepts:

- numeric IDs
- names, where the target lookup supports name resolution
- slugs, where the target lookup supports slug resolution

In practice, this means a configuration can often use a human-readable name or slug during authoring, while the provider resolves the value to the matching NetBox object during create and update operations.

### Resolution behavior

During create and update operations, the provider resolves configured reference values against NetBox before sending the API request.

Important implications:

- a valid ID is treated as already canonical
- a valid name or slug is resolved to the same target object as its numeric ID
- equivalent values should not cause churn just because one configuration uses a name or slug and another path returns an ID

### State behavior

Reference state is intentionally stability-oriented rather than strictly format-preserving.

The practical contract is:

- if an existing single-value reference already matches the API response by name, slug, or ID, the provider usually preserves the current configured form
- if a reference enters state through import, computed state, or another path where no user-chosen format is available, the provider commonly falls back to the numeric ID
- several list-style, set-style, computed, or explicitly normalized reference fields are stored as numeric IDs in state

As a result, users should treat numeric IDs as the canonical stable representation for references in Terraform state.

If long-term state stability matters more than authoring convenience, prefer IDs in configuration.

### Diff behavior

For reference fields using the shared equivalence logic, the provider suppresses plan noise when two different textual values resolve to the same NetBox object and the configuration is not forcing a different representation.

Examples of equivalent values include:

- `123`
- a matching object name
- a matching object slug

when they all resolve to the same NetBox record.

## `custom_fields` ownership model

Resource `custom_fields` use a filter-to-owned model.

That model exists so Terraform can manage selected custom fields without forcing ownership of every custom field present in NetBox.

### Mental model

- NetBox may contain more custom fields than Terraform is actively managing
- Terraform state for resources shows the custom fields Terraform currently owns
- unowned custom fields are preserved in NetBox during updates, but are not echoed back into managed resource state

### Configuration semantics

#### `custom_fields` omitted or `null`

Meaning:

- Terraform does not take ownership of custom fields for that resource instance in the current configuration
- existing NetBox custom field values are preserved during update operations
- resource state does not report those preserved fields as owned values

Operationally, this is the "preserve but do not manage" mode.

#### `custom_fields = []`

Meaning:

- Terraform explicitly requests removal of all custom fields on that object

This is the "clear everything" mode.

#### Non-empty `custom_fields`

Meaning:

- Terraform owns the named fields present in configuration
- those configured values are merged with existing NetBox values during update operations
- unowned fields already present in NetBox are preserved

If a previously managed field is removed from configuration:

- Terraform stops owning that field
- the field value is preserved in NetBox
- the field disappears from Terraform resource state because it is no longer owned

This is the core filter-to-owned behavior.

### Resource state versus data source state

Resource state and data source state intentionally differ:

- resources: return owned `custom_fields` only
- data sources: return all custom fields visible from the NetBox API response

That difference is intentional and part of the provider contract.

## Import behavior

### Standard import identifier

The standard import identifier is the resource numeric ID.

For resources using the shared import validation helper, non-numeric IDs are rejected when the resource requires a numeric import ID.

### Import identity with custom-field hints

Many resources also support import identity blocks with:

- `id`
- optional `custom_fields` hints in the form `name[:type]`

These hints seed identity information for custom-field-aware import flows without changing the ownership model described above.

### Imported `custom_fields`

On import there is no prior Terraform configuration to define ownership.

Therefore:

- the import operation itself does not claim ownership of every custom field present in NetBox
- imported resource state commonly starts with no owned `custom_fields`
- after import, adding `custom_fields` to configuration adopts management of those specific named fields

This behavior is intentional and matches the filter-to-owned model.

## Recommended usage

For the most predictable results:

- prefer numeric IDs for references when you want import-stable state
- use names or slugs when they improve readability and the target resource clearly supports lookup by those forms
- omit `custom_fields` when NetBox or another system should remain the owner
- use a non-empty `custom_fields` set only for the fields Terraform should actively manage
- use `custom_fields = []` only when the intent is to clear all custom fields on the object

## Notes

This guide describes the shared provider contract. Individual resources may still have resource-specific constraints based on what the NetBox API accepts for a given field.
