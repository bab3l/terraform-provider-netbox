# Terraform Provider Netbox - Resource Implementation TODO

This document tracks all potential resources and data sources that can be implemented based on the go-netbox API client.

## Implementation Checklist

For each new resource type, complete the following steps:

### 1. Code Implementation
- [ ] **Resource Implementation** (`internal/resources/<name>_resource.go`)
  - [ ] Define resource model struct with all attributes
  - [ ] Implement `Schema()` with proper types, descriptions, validators
  - [ ] Implement `Create()` - map Terraform state to API request
  - [ ] Implement `Read()` - refresh state from API
  - [ ] Implement `Update()` - handle attribute changes
  - [ ] Implement `Delete()` - remove resource from Netbox
  - [ ] Implement `ImportState()` - support `terraform import`
  - [ ] Handle null vs empty string for optional fields
  - [ ] Handle nested objects (tags, custom_fields) if applicable

- [ ] **Data Source Implementation** (`internal/datasources/<name>_data_source.go`)
  - [ ] Define data source model struct
  - [ ] Implement `Schema()` with lookup attributes (id, name, slug)
  - [ ] Implement `Read()` with support for multiple lookup methods
  - [ ] Ensure consistency with resource attribute names/types

- [ ] **Register in Provider** (`internal/provider/provider.go`)
  - [ ] Add resource to `Resources()` function
  - [ ] Add data source to `DataSources()` function

### 2. Unit Tests
- [ ] **Resource Unit Tests** (`internal/resources_test/<name>_resource_test.go`)
  - [ ] Test basic CRUD operations
  - [ ] Test with optional fields omitted
  - [ ] Test with all fields populated
  - [ ] Test import functionality
  - [ ] Test error handling for invalid inputs

- [ ] **Data Source Unit Tests** (`internal/datasources_test/<name>_data_source_test.go`)
  - [ ] Test lookup by ID
  - [ ] Test lookup by name
  - [ ] Test lookup by slug (if applicable)
  - [ ] Test not found error handling

### 3. Acceptance Tests
- [ ] **Resource Acceptance Test** (in resource test file)
  - [ ] Use `testutil.RandomName()` / `testutil.RandomSlug()` for unique names
  - [ ] Test create and read back
  - [ ] Test update in place
  - [ ] Test destroy (use `CheckDestroy` function)
  - [ ] Add cleanup with `testutil.CleanupResource()`

- [ ] **Data Source Acceptance Test** (in data source test file)
  - [ ] Create prerequisite resource first
  - [ ] Test all lookup methods return consistent data

### 4. Terraform Integration Tests
- [ ] **Resource Test** (`test/terraform/resources/<name>/`)
  - [ ] Create `main.tf` with example resource configuration
  - [ ] Create `outputs.tf` with validation outputs
  - [ ] Include `*_valid` boolean outputs for automated verification
  - [ ] Test relationships with dependent resources if applicable

- [ ] **Data Source Test** (`test/terraform/data-sources/<name>/`)
  - [ ] Create `main.tf` with resource + data source lookups
  - [ ] Create `outputs.tf` with `all_ids_match` and other validations
  - [ ] Test all lookup methods (id, name, slug)

- [ ] **Update Test Script** (`scripts/run-terraform-tests.ps1`)
  - [ ] Add new resource/data-source to `$testOrder` array

### 5. Documentation
- [ ] **Example Files** (`examples/resources/<name>/resource.tf`)
  - [ ] Create realistic usage example
  - [ ] Include comments explaining each attribute

- [ ] **Example Files** (`examples/data-sources/<name>/data-source.tf`)
  - [ ] Create example showing all lookup methods

- [ ] **Generate Documentation**
  ```powershell
  # Run from terraform-provider-netbox directory
  # Uses locally cloned terraform-plugin-docs repo
  c:\GitRoot\terraform-plugin-docs\tfplugindocs.exe generate --provider-dir=. --rendered-website-dir=docs
  ```
  > **Note:** Documentation is auto-generated from Go code schemas and templates in `templates/`. 
  > See `docs/DOCUMENTATION-GENERATION.md` for details.

- [ ] **Review Generated Docs** (`docs/resources/<name>.md`, `docs/data-sources/<name>.md`)
  - [ ] Verify descriptions are clear and complete
  - [ ] Verify examples render correctly

### 6. Validation & Testing
- [ ] **Build Provider**
  ```powershell
  go build .
  ```

- [ ] **Run Unit Tests**
  ```powershell
  go test ./internal/resources/... ./internal/datasources/... -v
  ```

- [ ] **Run Acceptance Tests** (requires running Netbox)
  ```powershell
  $env:TF_ACC = "1"
  $env:NETBOX_SERVER_URL = "http://localhost:8000"
  $env:NETBOX_API_TOKEN = "your-token"
  go test ./... -v -run "TestAcc"
  ```

- [ ] **Run Terraform Integration Tests**
  ```powershell
  .\scripts\run-terraform-tests.ps1
  ```

- [ ] **Verify All Tests Pass**
  - [ ] Unit tests: PASS
  - [ ] Acceptance tests: PASS
  - [ ] Terraform integration tests: PASS

### 7. Code Quality
- [ ] **Format Code**
  ```powershell
  go fmt ./...
  ```

- [ ] **Run Linter**
  ```powershell
  go vet ./...
  ```

- [ ] **Check for Errors**
  ```powershell
  # In VS Code, check Problems panel for any diagnostics
  ```

### 8. Final Review
- [ ] **Review API Coverage**
  - [ ] All required fields are implemented
  - [ ] All optional fields are implemented (or documented as future work)
  - [ ] Nested objects handled correctly

- [ ] **Review Error Messages**
  - [ ] Error messages are clear and actionable
  - [ ] Include context (resource type, ID, operation)

- [ ] **Update RESOURCE-IMPLEMENTATION-TODO.md**
  - [ ] Mark resource as ✅ Implemented
  - [ ] Update summary counts

---

## Implementation Patterns & Standards

This section documents the standard patterns used across all resource and data source implementations. Follow these patterns to ensure consistency.

### Package Structure

```
internal/
├── resources/           # Resource implementations
│   └── <name>_resource.go
├── datasources/         # Data source implementations
│   └── <name>_data_source.go
├── resources_test/      # Resource unit tests
├── datasources_test/    # Data source unit tests
├── netboxlookup/        # Lookup helpers for resolving references
│   └── lookup.go
├── utils/               # Shared utilities
│   └── common.go
└── validators/          # Custom validators
    └── validators.go
```

### Import Structure

Standard imports for a resource file:

```go
package resources

import (
    "context"
    "fmt"

    "github.com/bab3l/go-netbox"
    "github.com/hashicorp/terraform-plugin-framework/diag"           // If needed for helper functions
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"

    "github.com/bab3l/terraform-provider-netbox/internal/netboxlookup"  // If resolving references
    nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"  // Schema attribute helpers
    "github.com/bab3l/terraform-provider-netbox/internal/utils"
)
```

Standard imports for a data source file:

```go
package datasources

import (
    "context"
    "fmt"
    "net/http"

    "github.com/bab3l/go-netbox"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"

    nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"  // Schema attribute helpers
    "github.com/bab3l/terraform-provider-netbox/internal/utils"
)
```

### Resource Struct Pattern

```go
// Interface assertions - ensure we implement required interfaces
var _ resource.Resource = &ExampleResource{}
var _ resource.ResourceWithImportState = &ExampleResource{}

func NewExampleResource() resource.Resource {
    return &ExampleResource{}
}

// ExampleResource defines the resource implementation.
type ExampleResource struct {
    client *netbox.APIClient
}

// ExampleResourceModel describes the resource data model.
type ExampleResourceModel struct {
    ID           types.String `tfsdk:"id"`
    Name         types.String `tfsdk:"name"`
    Slug         types.String `tfsdk:"slug"`
    // Optional string fields
    Description  types.String `tfsdk:"description"`
    // Reference fields (use string with ID value)
    Parent       types.String `tfsdk:"parent"`
    Site         types.String `tfsdk:"site"`
    Tenant       types.String `tfsdk:"tenant"`
    // Status fields (enum strings)
    Status       types.String `tfsdk:"status"`
    // Numeric fields
    UHeight      types.Int64  `tfsdk:"u_height"`
    Weight       types.Float64 `tfsdk:"weight"`
    // Boolean fields
    DescUnits    types.Bool   `tfsdk:"desc_units"`
    // Nested objects
    Tags         types.Set    `tfsdk:"tags"`
    CustomFields types.Set    `tfsdk:"custom_fields"`
}
```

### Schema Patterns

#### Required String Attributes

```go
"name": schema.StringAttribute{
    MarkdownDescription: "Full name of the resource.",
    Required:            true,
    Validators: []validator.String{
        stringvalidator.LengthBetween(1, 100),
    },
},
"slug": schema.StringAttribute{
    MarkdownDescription: "URL-friendly identifier. Must be unique.",
    Required:            true,
    Validators: []validator.String{
        stringvalidator.LengthBetween(1, 100),
        validators.ValidSlug(),
    },
},
```

#### Optional String Attributes

```go
"description": schema.StringAttribute{
    MarkdownDescription: "Detailed description of the resource.",
    Optional:            true,
    Validators: []validator.String{
        stringvalidator.LengthAtMost(200),
    },
},
```

#### Computed-Only Attributes (ID)

```go
"id": schema.StringAttribute{
    Computed:            true,
    MarkdownDescription: "Unique identifier (assigned by Netbox).",
},
```

#### Optional with Computed (has default value)

```go
"status": schema.StringAttribute{
    MarkdownDescription: "Operational status. Valid values: `planned`, `staging`, `active`, `decommissioning`, `retired`.",
    Optional:            true,
    Computed:            true,  // Netbox sets default
    Validators: []validator.String{
        stringvalidator.OneOf(
            "planned",
            "staging", 
            "active",
            "decommissioning",
            "retired",
        ),
    },
},
```

#### Reference Attributes (Foreign Keys)

```go
// For hierarchical self-reference (parent of same type)
"parent": schema.StringAttribute{
    MarkdownDescription: "ID of the parent resource.",
    Optional:            true,
    Validators: []validator.String{
        stringvalidator.RegexMatches(
            validators.IntegerRegex(),
            "must be a valid integer ID",
        ),
    },
},

// For references to other resource types (use name/slug lookup)
"site": schema.StringAttribute{
    MarkdownDescription: "Name, slug, or ID of the site.",
    Required:            true,  // or Optional for nullable references
},
```

#### Numeric Attributes

```go
"u_height": schema.Int64Attribute{
    MarkdownDescription: "Height in rack units.",
    Optional:            true,
    Computed:            true,  // If has default
    Validators: []validator.Int64{
        int64validator.Between(1, 100),
    },
},
"weight": schema.Float64Attribute{
    MarkdownDescription: "Weight of the resource.",
    Optional:            true,
},
```

#### Boolean Attributes

```go
"desc_units": schema.BoolAttribute{
    MarkdownDescription: "If true, units are numbered top-to-bottom.",
    Optional:            true,
    Computed:            true,  // If has default
},
```

#### Tags Nested Attribute

```go
"tags": schema.SetNestedAttribute{
    MarkdownDescription: "Tags assigned to this resource.",
    Optional:            true,
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of the existing tag.",
                Required:            true,
                Validators: []validator.String{
                    stringvalidator.LengthBetween(1, 100),
                },
            },
            "slug": schema.StringAttribute{
                MarkdownDescription: "Slug of the existing tag.",
                Required:            true,
                Validators: []validator.String{
                    stringvalidator.LengthBetween(1, 100),
                    validators.ValidSlug(),
                },
            },
        },
    },
},
```

#### Custom Fields Nested Attribute

```go
"custom_fields": schema.SetNestedAttribute{
    MarkdownDescription: "Custom fields assigned to this resource.",
    Optional:            true,
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{
                MarkdownDescription: "Name of the custom field.",
                Required:            true,
                Validators: []validator.String{
                    stringvalidator.LengthBetween(1, 50),
                    validators.ValidCustomFieldName(),
                },
            },
            "type": schema.StringAttribute{
                MarkdownDescription: "Type of the custom field.",
                Required:            true,
                Validators: []validator.String{
                    validators.ValidCustomFieldType(),
                },
            },
            "value": schema.StringAttribute{
                MarkdownDescription: "Value of the custom field.",
                Required:            true,
                Validators: []validator.String{
                    stringvalidator.LengthAtMost(1000),
                    validators.SimpleValidCustomFieldValue(),
                },
            },
        },
    },
},
```

### nbschema Attribute Helpers

The `internal/schema` package (imported as `nbschema`) provides factory functions for common schema attributes. This ensures consistency across all resources and data sources and reduces boilerplate code.

#### Import

```go
import (
    nbschema "github.com/bab3l/terraform-provider-netbox/internal/schema"
)
```

#### Resource Attribute Helpers

```go
// Common attributes
"id":            nbschema.IDAttribute(),                    // Computed ID
"name":          nbschema.NameAttribute("resource name"),   // Required name with validation
"slug":          nbschema.SlugAttribute("resource name"),   // Required slug with validation
"description":   nbschema.DescriptionAttribute(),           // Optional description
"comments":      nbschema.CommentsAttribute(),              // Optional comments (longer text)

// Reference attributes (for foreign keys)
"parent":        nbschema.ReferenceAttribute("parent resource", false),     // Optional parent
"site":          nbschema.RequiredReferenceAttribute("site"),               // Required reference
"tenant":        nbschema.ReferenceAttribute("tenant", false),              // Optional reference
"manufacturer":  nbschema.IDOnlyReferenceAttribute("manufacturer", true),   // ID-only required ref

// Specialized attributes
"color":         nbschema.ColorAttribute(),                 // Optional hex color
"status":        nbschema.StatusAttribute("resource", []string{"active", "planned", "retired"}),
"serial":        nbschema.SerialAttribute(),                // Optional serial number
"asset_tag":     nbschema.AssetTagAttribute(),              // Optional asset tag
"facility":      nbschema.FacilityAttribute(),              // Optional facility ID
"model":         nbschema.ModelAttribute("resource", 100),  // Model name with max length

// Boolean with default
"vm_role":       nbschema.BoolAttributeWithDefault("VM role description", false),

// Nested objects
"tags":          nbschema.TagsAttribute(),                  // Optional tags set
"custom_fields": nbschema.CustomFieldsAttribute(),          // Optional custom fields set
```

#### Data Source Attribute Helpers (DS* prefix)

```go
// Lookup fields (optional input, computed output)
"id":            nbschema.DSIDAttribute("resource"),        // Optional/Computed ID lookup
"name":          nbschema.DSNameAttribute("resource"),      // Optional/Computed name lookup
"slug":          nbschema.DSSlugAttribute("resource"),      // Optional/Computed slug lookup

// Computed string attributes
"description":   nbschema.DSComputedStringAttribute("Description of the resource."),
"status":        nbschema.DSComputedStringAttribute("Status of the resource."),
"parent":        nbschema.DSComputedStringAttribute("Name of the parent."),
"parent_id":     nbschema.DSComputedStringAttribute("ID of the parent."),

// Computed numeric attributes
"u_height":      nbschema.DSComputedInt64Attribute("Height in rack units."),
"weight":        nbschema.DSComputedFloat64Attribute("Weight of the resource."),

// Computed boolean attributes
"vm_role":       nbschema.DSComputedBoolAttribute("Whether this is a VM role."),

// Nested objects
"tags":          nbschema.DSTagsAttribute(),                // Computed tags
"custom_fields": nbschema.DSCustomFieldsAttribute(),        // Computed custom fields
```

#### Example Resource Schema Using nbschema

```go
func (r *ExampleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Manages an example resource in Netbox.",
        Attributes: map[string]schema.Attribute{
            "id":           nbschema.IDAttribute(),
            "name":         nbschema.NameAttribute("example"),
            "slug":         nbschema.SlugAttribute("example"),
            "description":  nbschema.DescriptionAttribute(),
            "site":         nbschema.RequiredReferenceAttribute("site"),
            "tenant":       nbschema.ReferenceAttribute("tenant", false),
            "status":       nbschema.StatusAttribute("example", []string{"active", "planned"}),
            "tags":         nbschema.TagsAttribute(),
            "custom_fields": nbschema.CustomFieldsAttribute(),
        },
    }
}
```

#### Example Data Source Schema Using nbschema

```go
func (d *ExampleDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "Use this data source to get information about an example in Netbox.",
        Attributes: map[string]schema.Attribute{
            "id":           nbschema.DSIDAttribute("example"),
            "name":         nbschema.DSNameAttribute("example"),
            "slug":         nbschema.DSSlugAttribute("example"),
            "description":  nbschema.DSComputedStringAttribute("Description of the example."),
            "site":         nbschema.DSComputedStringAttribute("Name of the site."),
            "site_id":      nbschema.DSComputedStringAttribute("ID of the site."),
            "status":       nbschema.DSComputedStringAttribute("Status of the example."),
            "tags":         nbschema.DSTagsAttribute(),
            "custom_fields": nbschema.DSCustomFieldsAttribute(),
        },
    }
}
```

### API Request Patterns

#### Creating API Requests

```go
// Use the Writable*Request type for create/update
request := netbox.NewWritableExampleRequest(data.Name.ValueString(), data.Slug.ValueString())

// For optional string fields - use pointer
if !data.Description.IsNull() && !data.Description.IsUnknown() {
    desc := data.Description.ValueString()
    request.Description = &desc
}

// For nullable integer references (parent ID, etc.)
if !data.Parent.IsNull() && !data.Parent.IsUnknown() {
    var parentIDInt int32
    if _, err := fmt.Sscanf(data.Parent.ValueString(), "%d", &parentIDInt); err != nil {
        resp.Diagnostics.AddError("Invalid Parent ID", 
            fmt.Sprintf("Parent ID must be a number, got: %s", data.Parent.ValueString()))
        return
    }
    request.Parent = *netbox.NewNullableInt32(&parentIDInt)
}

// For references using lookup (returns Brief*Request)
if !data.Site.IsNull() && !data.Site.IsUnknown() {
    siteRef, diags := netboxlookup.LookupSiteBrief(ctx, r.client, data.Site.ValueString())
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    request.Site = *siteRef  // For required fields
    // or
    request.Site = *netbox.NewNullableBriefSiteRequest(siteRef)  // For optional fields
}

// For enum/status fields
if !data.Status.IsNull() && !data.Status.IsUnknown() {
    status := netbox.LocationStatusValue(data.Status.ValueString())
    request.Status = &status
}

// For tags
if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
    tags, diags := utils.TagModelsToNestedTagRequests(ctx, data.Tags)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    request.Tags = tags
}

// For custom fields
if !data.CustomFields.IsNull() && !data.CustomFields.IsUnknown() {
    var customFieldModels []utils.CustomFieldModel
    diags := data.CustomFields.ElementsAs(ctx, &customFieldModels, false)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    request.CustomFields = utils.CustomFieldModelsToMap(customFieldModels)
}
```

#### Calling the API

```go
// Create
result, httpResp, err := r.client.DcimAPI.DcimExamplesCreate(ctx).
    WritableExampleRequest(*request).Execute()

// Read  
result, httpResp, err := r.client.DcimAPI.DcimExamplesRetrieve(ctx, idInt).Execute()

// Update
result, httpResp, err := r.client.DcimAPI.DcimExamplesUpdate(ctx, idInt).
    WritableExampleRequest(*request).Execute()

// Delete
httpResp, err := r.client.DcimAPI.DcimExamplesDestroy(ctx, idInt).Execute()
```

### Mapping Response to State

```go
// Required fields
data.ID = types.StringValue(fmt.Sprintf("%d", result.GetId()))
data.Name = types.StringValue(result.GetName())
data.Slug = types.StringValue(result.GetSlug())

// Optional string fields
if result.HasDescription() && result.GetDescription() != "" {
    data.Description = types.StringValue(result.GetDescription())
} else {
    data.Description = types.StringNull()
}

// Nested object references (parent, site, tenant, etc.)
if result.HasParent() && result.GetParent().Id != 0 {
    parent := result.GetParent()
    data.Parent = types.StringValue(fmt.Sprintf("%d", parent.GetId()))
} else {
    data.Parent = types.StringNull()
}

// Status/enum fields
if result.HasStatus() {
    status := result.GetStatus()
    data.Status = types.StringValue(string(status.GetValue()))
} else {
    data.Status = types.StringNull()
}

// Numeric fields with defaults
if result.HasUHeight() {
    data.UHeight = types.Int64Value(int64(result.GetUHeight()))
} else {
    data.UHeight = types.Int64Null()
}

// Boolean fields
if result.HasDescUnits() {
    data.DescUnits = types.BoolValue(result.GetDescUnits())
} else {
    data.DescUnits = types.BoolNull()
}

// Tags
if result.HasTags() {
    tags := utils.NestedTagsToTagModels(result.GetTags())
    tagsValue, diags := types.SetValueFrom(ctx, utils.GetTagsAttributeType().ElemType, tags)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }
    data.Tags = tagsValue
} else {
    data.Tags = types.SetNull(utils.GetTagsAttributeType().ElemType)
}

// Custom fields (requires existing state for type info)
if result.HasCustomFields() {
    var existingModels []utils.CustomFieldModel
    if !data.CustomFields.IsNull() {
        diags := data.CustomFields.ElementsAs(ctx, &existingModels, false)
        resp.Diagnostics.Append(diags...)
    }
    customFields := utils.MapToCustomFieldModels(result.GetCustomFields(), existingModels)
    customFieldsValue, diags := types.SetValueFrom(ctx, 
        utils.GetCustomFieldsAttributeType().ElemType, customFields)
    resp.Diagnostics.Append(diags...)
    data.CustomFields = customFieldsValue
} else {
    data.CustomFields = types.SetNull(utils.GetCustomFieldsAttributeType().ElemType)
}
```

### Netbox Lookup Helpers

Add lookup functions to `internal/netboxlookup/lookup.go` for each referenced type:

```go
// LookupExampleBrief returns a BriefExampleRequest from an ID or slug
func LookupExampleBrief(ctx context.Context, client *netbox.APIClient, value string) (*netbox.BriefExampleRequest, diag.Diagnostics) {
    var id int32
    if _, err := fmt.Sscanf(value, "%d", &id); err == nil {
        // Lookup by ID
        resource, resp, err := client.DcimAPI.DcimExamplesRetrieve(ctx, id).Execute()
        if err != nil || resp.StatusCode != 200 {
            return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
                "Example lookup failed", err.Error())}
        }
        return &netbox.BriefExampleRequest{
            Name: resource.GetName(),
            Slug: resource.GetSlug(),
        }, nil
    }
    // Lookup by slug
    list, resp, err := client.DcimAPI.DcimExamplesList(ctx).Slug([]string{value}).Execute()
    if err != nil || resp.StatusCode != 200 {
        return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
            "Example lookup failed", 
            fmt.Sprintf("Could not find example with slug '%s': %v", value, err))}
    }
    if list != nil && len(list.Results) > 0 {
        resource := list.Results[0]
        return &netbox.BriefExampleRequest{
            Name: resource.GetName(),
            Slug: resource.GetSlug(),
        }, nil
    }
    return nil, diag.Diagnostics{diag.NewErrorDiagnostic(
        "Example lookup failed", 
        fmt.Sprintf("No example found with slug '%s'", value))}
}
```

### Data Source Patterns

Data sources follow a similar pattern but only implement Read:

```go
var _ datasource.DataSource = &ExampleDataSource{}

type ExampleDataSourceModel struct {
    ID           types.String `tfsdk:"id"`       // Optional for lookup
    Name         types.String `tfsdk:"name"`     // Optional for lookup
    Slug         types.String `tfsdk:"slug"`     // Optional for lookup
    // All other attributes are Computed only
    Description  types.String `tfsdk:"description"`
    // ...
}

func (d *ExampleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ExampleDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    
    var result *netbox.Example
    var httpResp *http.Response
    var err error

    // Lookup by ID
    if !data.ID.IsNull() {
        var idInt int32
        if _, parseErr := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); parseErr != nil {
            resp.Diagnostics.AddError("Invalid ID", "ID must be a number")
            return
        }
        result, httpResp, err = d.client.DcimAPI.DcimExamplesRetrieve(ctx, idInt).Execute()
    } else if !data.Slug.IsNull() {
        // Lookup by slug
        list, listResp, listErr := d.client.DcimAPI.DcimExamplesList(ctx).
            Slug([]string{data.Slug.ValueString()}).Execute()
        httpResp = listResp
        err = listErr
        if err == nil && len(list.GetResults()) > 0 {
            result = &list.GetResults()[0]
        } else if len(list.GetResults()) == 0 {
            resp.Diagnostics.AddError("Not Found", 
                fmt.Sprintf("No example found with slug: %s", data.Slug.ValueString()))
            return
        }
    } else if !data.Name.IsNull() {
        // Lookup by name (may return multiple - warn user)
        list, listResp, listErr := d.client.DcimAPI.DcimExamplesList(ctx).
            Name([]string{data.Name.ValueString()}).Execute()
        httpResp = listResp
        err = listErr
        if err == nil {
            if len(list.GetResults()) == 0 {
                resp.Diagnostics.AddError("Not Found", "...")
                return
            }
            if len(list.GetResults()) > 1 {
                resp.Diagnostics.AddError("Multiple Found", "...")
                return
            }
            result = &list.GetResults()[0]
        }
    } else {
        resp.Diagnostics.AddError("Missing Identifier", 
            "Either 'id', 'slug', or 'name' must be specified")
        return
    }
    
    // Handle errors and map to state...
}
```

### Error Handling Patterns

```go
// Standard API error with response body
if err != nil {
    resp.Diagnostics.AddError(
        "Error creating example",
        utils.FormatAPIError("create example", err, httpResp),
    )
    return
}

// HTTP status check
if httpResp.StatusCode != 201 {
    resp.Diagnostics.AddError(
        "Error creating example",
        fmt.Sprintf("Expected HTTP 201, got: %d", httpResp.StatusCode),
    )
    return
}

// Resource not found on Read (should remove from state)
if err != nil {
    if httpResp != nil && httpResp.StatusCode == 404 {
        resp.State.RemoveResource(ctx)
        return
    }
    resp.Diagnostics.AddError(...)
    return
}
```

### Import State

```go
func (r *ExampleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
    resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
```

### Provider Registration

In `internal/provider/provider.go`:

```go
func (p *NetboxProvider) Resources(ctx context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        resources.NewExampleResource,
        // ...
    }
}

func (p *NetboxProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        datasources.NewExampleDataSource,
        // ...
    }
}
```

### Utility Functions Reference

| Function | Location | Purpose |
|----------|----------|---------|
| `utils.TagModelsToNestedTagRequests(ctx, tagsSet)` | utils/common.go | Convert Terraform tags Set to API request format |
| `utils.NestedTagsToTagModels(tags)` | utils/common.go | Convert API tags to Terraform model |
| `utils.CustomFieldModelsToMap(models)` | utils/common.go | Convert custom fields to API map format |
| `utils.MapToCustomFieldModels(map, existing)` | utils/common.go | Convert API custom fields to Terraform model |
| `utils.GetTagsAttributeType()` | utils/common.go | Get the attribute type for tags Set |
| `utils.GetCustomFieldsAttributeType()` | utils/common.go | Get the attribute type for custom fields Set |
| `utils.FormatAPIError(op, err, resp)` | utils/common.go | Format API error with response body |
| `validators.ValidSlug()` | validators/validators.go | Validate slug format |
| `validators.IntegerRegex()` | validators/validators.go | Regex for integer validation |
| `validators.ValidCustomFieldName()` | validators/validators.go | Validate custom field name |
| `validators.ValidCustomFieldType()` | validators/validators.go | Validate custom field type |
| `netboxlookup.LookupSiteBrief(ctx, client, value)` | netboxlookup/lookup.go | Resolve site by ID or slug |
| `netboxlookup.LookupTenantBrief(ctx, client, value)` | netboxlookup/lookup.go | Resolve tenant by ID |
| `netboxlookup.LookupRegionBrief(ctx, client, value)` | netboxlookup/lookup.go | Resolve region by ID |
| `netboxlookup.LookupLocationBrief(ctx, client, value)` | netboxlookup/lookup.go | Resolve location by ID or slug |
| `netboxlookup.LookupRackRoleBrief(ctx, client, value)` | netboxlookup/lookup.go | Resolve rack role by ID or slug |
| `netboxlookup.LookupRackTypeBrief(ctx, client, value)` | netboxlookup/lookup.go | Resolve rack type by ID or model |

### Lessons Learned & Common Pitfalls

This section documents common mistakes and their solutions based on actual implementation experience.

#### 1. Nullable Types Require Pointers

**Problem:** Compile error when using `NewNullableInt32(value)` or `NewNullableFloat64(value)`.

**Solution:** These functions take **pointers**, not values:

```go
// ❌ WRONG
rackRequest.MaxWeight = *netbox.NewNullableInt32(maxWeight)

// ✅ CORRECT
rackRequest.MaxWeight = *netbox.NewNullableInt32(&maxWeight)
```

Same applies to:
- `NewNullableFloat64(&value)`
- `NewNullableInt32(&value)` 
- `NewNullableInt64(&value)`
- `NewNullableString(&value)`

#### 2. Enum Types Are Often Integer-Based

**Problem:** Trying to cast a string directly to an enum type like `PatchedWritableRackRequestWidth`.

**Solution:** Check the underlying type - many enums are `int32`, not `string`:

```go
// ❌ WRONG - Width is an int32 enum (10, 19, 21, 23)
widthValue := netbox.PatchedWritableRackRequestWidth(data.Width.ValueString())

// ✅ CORRECT - Parse to int32 first, then use constructor
var widthInt int32
if _, err := fmt.Sscanf(data.Width.ValueString(), "%d", &widthInt); err == nil {
    widthValue, err := netbox.NewPatchedWritableRackRequestWidthFromValue(widthInt)
    if err == nil {
        rackRequest.Width = widthValue
    }
}
```

Check the model file to determine the underlying type:
- `type PatchedWritableRackRequestWidth int32` - use `FromValue(int32)`
- `type PatchedWritableRackRequestStatus string` - can cast directly

#### 3. Use Correct Request Type

**Problem:** Using `RackRequest` instead of `WritableRackRequest`.

**Solution:** Always use the `Writable*Request` type for create/update operations:

```go
// ❌ WRONG - RackRequest may not exist or have different fields
request := netbox.RackRequest{...}

// ✅ CORRECT - Use WritableRackRequest
rackRequest := netbox.WritableRackRequest{
    Name: data.Name.ValueString(),
    Site: *siteRef,
}
```

#### 4. Lookup Functions Return Diagnostics, Not Errors

**Problem:** Expecting lookup functions to return `error`.

**Solution:** Lookup functions return `diag.Diagnostics`:

```go
// ❌ WRONG
siteRequest, err := netboxlookup.LookupSiteBrief(ctx, r.client, data.Site.ValueInt64())
if err != nil { ... }

// ✅ CORRECT
siteRef, diags := netboxlookup.LookupSiteBrief(ctx, r.client, data.Site.ValueString())
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
```

#### 5. Use String Type for IDs in Terraform Schema

**Problem:** Using `types.Int64` for ID fields causes issues with import and state management.

**Solution:** Use `types.String` for all ID fields (including foreign keys):

```go
// ❌ WRONG
ID   types.Int64  `tfsdk:"id"`
Site types.Int64  `tfsdk:"site_id"`

// ✅ CORRECT
ID   types.String `tfsdk:"id"`
Site types.String `tfsdk:"site"`  // Store as "123" string, parse when needed
```

Convert during API operations:
```go
// Parse string ID to int32 for API call
var idInt int32
if _, err := fmt.Sscanf(data.ID.ValueString(), "%d", &idInt); err != nil {
    resp.Diagnostics.AddError("Invalid ID", err.Error())
    return
}
```

#### 6. Weight Unit Uses DeviceTypeWeightUnitValue (Not Request-Specific)

**Problem:** Looking for `PatchedWritableRackRequestWeightUnit` type.

**Solution:** Weight unit is shared across types:

```go
// ✅ CORRECT - Use DeviceTypeWeightUnitValue
weightUnitValue := netbox.DeviceTypeWeightUnitValue(data.WeightUnit.ValueString())
rackRequest.WeightUnit = &weightUnitValue
```

#### 7. Required Fields in Constructors

**Problem:** Using struct literal when constructor requires certain fields.

**Solution:** Check for `New*Request()` constructors that enforce required fields:

```go
// Struct literal works but misses validation
request := netbox.WritableRackRequest{Name: name, Site: site}

// Constructor ensures required fields
request := netbox.NewWritableRackRequest(name, site)
```

#### 8. Brief*Request vs Nullable Reference Types

**Problem:** Confusing when to use `Brief*Request` directly vs `NullableBrief*Request`.

**Solution:**
- **Required references**: Use `Brief*Request` directly
- **Optional references**: Wrap with `NewNullableBrief*Request()`

```go
// Required field - Site is mandatory for racks
rackRequest.Site = *siteRef  // BriefSiteRequest

// Optional field - Location is nullable
rackRequest.Location = *netbox.NewNullableBriefLocationRequest(locationRef)
```

---

## Legend
- ✅ Implemented (resource + data source)
- 🔶 Partial (resource only or data source only)
- ⬜ Not implemented

---

## Summary

| Category | Total | Implemented | Has Tests | Notes |
|----------|-------|-------------|-----------|-------|
| DCIM (Data Center Infrastructure) | 30 | 30 | 16 | All implemented, need more tests |
| Tenancy | 6 | 5 | 5 | contact_assignment missing |
| IPAM (IP Address Management) | 14 | 9 | 5 | Most implemented |
| Virtualization | 6 | 5 | 5 | virtual_disk missing |
| Circuits | 7 | 6 | 2 | circuit_group* missing |
| VPN | 9 | 0 | 0 | Not started |
| Wireless | 3 | 2 | 0 | wireless_link missing |
| Extras | 14 | 4 | 4 | config_template, custom_field added |
| Users | 4 | 0 | 0 | Not started |
| Core | 1 | 0 | 0 | Not started |
| **TOTAL** | **94** | **61** | **37** | 65% implemented |

### Test Coverage Summary

| Type | Implemented | Has Unit Tests | Has Acceptance Tests |
|------|-------------|----------------|---------------------|
| Resources | 66 | 41 | 41 |
| Data Sources | 66 | 28 | 28 |

**Resources missing tests (25):**
aggregate, asn, circuit_termination, console_port, console_server_port, custom_field, device_bay, inventory_item, inventory_item_role, ip_range, module, module_bay, module_type, power_feed, power_outlet, power_panel, power_port, provider_account, provider_network, rack_type, rir, role, service, virtual_chassis, wireless_lan, wireless_lan_group

**Data sources missing tests (38):**
aggregate, asn, circuit_termination, console_port, console_port_template, console_server_port, console_server_port_template, custom_field, device_bay, inventory_item, inventory_item_role, ip_range, module, module_bay, module_type, power_feed, power_outlet, power_outlet_template, power_panel, power_port, power_port_template, provider_account, provider_network, rack_type, rir, role, service, virtual_chassis, wireless_lan, wireless_lan_group

---

## DCIM (Data Center Infrastructure Management)

### Infrastructure
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_site` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_site_group` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_region` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_location` | ✅ | ✅ | ✅ | ✅ | Fully implemented |

### Racks
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_rack` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_rack_role` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_rack_type` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_rack_reservation` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Devices
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_manufacturer` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_platform` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_device_type` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_device_role` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_device` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_device_bay` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_device_bay_template` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_virtual_chassis` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_virtual_device_context` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Modules
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_module` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_module_type` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_module_bay` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_module_bay_template` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Interfaces & Ports
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_interface` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_interface_template` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_console_port` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_console_port_template` | ✅ | ✅ | ✅ | ⬜ | DS needs tests |
| `netbox_console_server_port` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_console_server_port_template` | ✅ | ✅ | ✅ | ⬜ | DS needs tests |
| `netbox_front_port` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_front_port_template` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_rear_port` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_rear_port_template` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Power
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_power_panel` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_power_feed` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_power_port` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_power_port_template` | ✅ | ✅ | ✅ | ⬜ | DS needs tests |
| `netbox_power_outlet` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_power_outlet_template` | ✅ | ✅ | ✅ | ⬜ | DS needs tests |

### Cabling
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_cable` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_cable_termination` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Inventory
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_inventory_item` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_inventory_item_role` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_inventory_item_template` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
|----------|--------|----------|-------|
| `netbox_inventory_item` | ⬜ | Low | Device inventory tracking |
| `netbox_inventory_item_role` | ⬜ | Low | Inventory item categorization |
| `netbox_inventory_item_template` | ⬜ | Low | Inventory item templates |

---

## Tenancy

| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_tenant` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_tenant_group` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_contact` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_contact_group` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_contact_role` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_contact_assignment` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## IPAM (IP Address Management)

### Core IPAM
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_vrf` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_prefix` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_ip_address` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_ip_range` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_aggregate` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |

### IPAM Organization
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_rir` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_role` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_route_target` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### VLANs
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_vlan` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_vlan_group` | ✅ | ✅ | ✅ | ✅ | Fully implemented |

### ASNs
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_asn` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_asn_range` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Services
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_service` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_service_template` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### FHRP
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_fhrp_group` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_fhrp_group_assignment` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Virtualization

| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_cluster` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_cluster_type` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_cluster_group` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_virtual_machine` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_vm_interface` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_virtual_disk` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Circuits

| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_provider` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_provider_account` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_provider_network` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_circuit` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_circuit_type` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_circuit_termination` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_circuit_group` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_circuit_group_assignment` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## VPN

### IPSec
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_ike_policy` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_ike_proposal` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_ipsec_policy` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_ipsec_profile` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_ipsec_proposal` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Tunnels
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_tunnel` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_tunnel_group` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_tunnel_termination` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### L2VPN
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_l2vpn` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_l2vpn_termination` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Wireless

| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_wireless_lan` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_wireless_lan_group` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_wireless_link` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Extras (Customization & Automation)

### Tags & Custom Fields
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_tag` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_custom_field` | ✅ | ⬜ | ✅ | ⬜ | Needs tests |
| `netbox_custom_field_choice_set` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_custom_link` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Configuration & Templates
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_config_context` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_config_template` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_export_template` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Automation
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_webhook` | ✅ | ✅ | ✅ | ✅ | Fully implemented |
| `netbox_event_rule` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_script` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Documentation
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_journal_entry` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_image_attachment` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_bookmark` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Notifications
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_notification` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_notification_group` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_subscription` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Filters
| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_saved_filter` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Users (Limited Scope)

| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_user` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_group` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_permission` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_token` | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Core

| Resource | Status | Tests | Data Source | DS Tests | Notes |
|----------|--------|-------|-------------|----------|-------|
| `netbox_data_source` | ⬜ | ⬜ | ⬜ | ⬜ | Not started | |

---

## Recommended Implementation Order

### Phase 1: Core Infrastructure (High Priority) ✅
1. ✅ `netbox_region` - Geographic hierarchy
2. ✅ `netbox_location` - Physical locations
3. ✅ `netbox_rack` - Rack infrastructure
4. ✅ `netbox_device_type` - Device templates
5. ✅ `netbox_device_role` - Device classification
6. ✅ `netbox_device` - Physical devices
7. ✅ `netbox_interface` - Network interfaces

### Phase 2: IPAM Essentials
8. `netbox_vrf` - Virtual routing
9. `netbox_prefix` - IP subnets
10. `netbox_ip_address` - IP addresses
11. `netbox_vlan` - VLANs
12. `netbox_vlan_group` - VLAN organization

### Phase 3: Virtualization
13. `netbox_cluster_type` - Cluster types
14. `netbox_cluster` - VM clusters
15. `netbox_virtual_machine` - VMs
16. `netbox_vm_interface` - VM interfaces

### Phase 4: Circuits & Connectivity
17. `netbox_provider` - Circuit providers
18. `netbox_circuit` - WAN circuits
19. `netbox_circuit_type` - Circuit classifications
20. `netbox_cable` - Physical cabling

### Phase 5: Extras & Customization
21. `netbox_tag` - Tagging
22. `netbox_config_context` - Config contexts
23. `netbox_webhook` - Automation hooks
24. `netbox_contact` - Contact management

### Phase 6: Advanced Features
25. VPN resources
26. Wireless resources
27. Power management
28. Remaining DCIM templates

---

## Notes

- Each resource should have a corresponding data source
- All resources should support:
  - Tags (where applicable)
  - Custom fields
  - Standard CRUD operations
  - Import functionality
- Consider read-only data sources for computed/derived data (e.g., available IPs, available prefixes)
- Template resources (e.g., `*_template`) are lower priority but useful for device type management

---

_Last updated: December 6, 2025_
