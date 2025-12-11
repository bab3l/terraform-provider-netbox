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
- ✅ Implemented (resource + data source + all tests)
- 🔶 Partial (implemented but missing some tests)
- ⬜ Not implemented

---

## Summary

| Category | Total Resources | Implemented | Unit Tests | TF Integration Tests | Notes |
|----------|----------------|-------------|------------|---------------------|-------|
| DCIM (Data Center Infrastructure) | ~35 | 31 | 31 | 31 | Core complete |
| Tenancy | 6 | 6 | 6 | 6 | Complete |
| IPAM (IP Address Management) | ~14 | 11 | 11 | 11 | route_target, asn_range added |
| Virtualization | 6 | 6 | 6 | 6 | Complete |
| Circuits | 8 | 8 | 8 | 8 | Complete |
| VPN | ~10 | 10 | 10 | 10 | L2VPN complete |
| Wireless | 3 | 3 | 3 | 3 | Complete |
| Extras | ~14 | 5 | 5 | 5 | Most extras not started |
| Users | 4 | 0 | 0 | 0 | Not started |
| Core | 1 | 0 | 0 | 0 | Not started |
| **TOTAL** | **~101** | **84** | **84** | **84** | **83% implemented** |

### Implementation Status by Type

| Type | Implemented | Unit Tests | TF Integration Tests |
|------|-------------|------------|---------------------|
| Resources | 83 | 83 ✅ | In progress (see tables) |
| Data Sources | 83 | 83 ✅ | In progress (see tables) |

**Resource TF integration tests added for rack_reservation, virtual_device_context, module_bay_template, inventory_item_template, front_port, front_port_template, and rear_port_template; matching data source TF coverage added for front_port, front_port_template, and rear_port_template. Remaining gaps are still marked ⬜ below.**

---

## DCIM (Data Center Infrastructure Management)

### Infrastructure
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_site` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_site_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_region` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_location` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Racks
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_rack` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_rack_role` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_rack_type` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_rack_reservation` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Devices
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_manufacturer` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_platform` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_device_type` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_device_role` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_device` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_device_bay` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_device_bay_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_virtual_chassis` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_virtual_device_context` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Modules
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_module` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_module_type` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_module_bay` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_module_bay_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Interfaces & Ports
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_interface` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_interface_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_console_port` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_console_port_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_console_server_port` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_console_server_port_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_front_port` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Resource & data source integration tests added |
| `netbox_front_port_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Resource & data source integration tests added |
| `netbox_rear_port` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_rear_port_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Resource & data source integration tests added |

### Power
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_power_panel` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_power_feed` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_power_port` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_power_port_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_power_outlet` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_power_outlet_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Cabling
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_cable` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_cable_termination` | N/A | N/A | N/A | N/A | N/A | N/A | DEPRECATED: Use `netbox_cable` with embedded terminations |

### Inventory
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_inventory_item` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_inventory_item_role` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_inventory_item_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

---

## Tenancy

| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_tenant` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_tenant_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_contact` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_contact_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_contact_role` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_contact_assignment` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete (note: role is required by NetBox API) |

---

## IPAM (IP Address Management)

### Core IPAM
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_vrf` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_prefix` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_ip_address` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_ip_range` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_aggregate` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### IPAM Organization
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_rir` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_role` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_route_target` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### VLANs
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_vlan` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_vlan_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### ASNs
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_asn` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_asn_range` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Services
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_service` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_service_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### FHRP
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_fhrp_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Completed |
| `netbox_fhrp_group_assignment` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

---

## Virtualization

| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_cluster` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_cluster_type` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_cluster_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_virtual_machine` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_vm_interface` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_virtual_disk` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

---

## Circuits

| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_provider` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_provider_account` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_provider_network` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_circuit` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_circuit_type` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_circuit_termination` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_circuit_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_circuit_group_assignment` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

---

## VPN

### IPSec
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_ike_policy` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_ike_proposal` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_ipsec_policy` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_ipsec_profile` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_ipsec_proposal` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Tunnels
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_tunnel` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_tunnel_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_tunnel_termination` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### L2VPN
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_l2vpn` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_l2vpn_termination` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

---

## Wireless

| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_wireless_lan` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_wireless_lan_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_wireless_link` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

---

## Extras (Customization & Automation)

### Tags & Custom Fields
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_tag` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_custom_field` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_custom_field_choice_set` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_custom_link` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Configuration & Templates
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_config_context` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_config_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_export_template` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |

### Automation
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_webhook` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_event_rule` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_script` | N/A | N/A | N/A | ✅ | ✅ | ✅ | Data source only (scripts are read-only in NetBox API) |

### Documentation
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_journal_entry` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_image_attachment` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Requires file upload - complex |
| `netbox_bookmark` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Requires user reference - complex |

### Notifications
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_notification` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_notification_group` | ✅ | ✅ | ✅ | ✅ | ✅ | ✅ | Complete |
| `netbox_subscription` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

### Filters
| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_saved_filter` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Users (Limited Scope)

| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_user` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_group` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_permission` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |
| `netbox_token` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Core

| Resource | Status | Unit Tests | TF Tests | Data Source | DS Unit Tests | DS TF Tests | Notes |
|----------|--------|------------|----------|-------------|---------------|-------------|-------|
| `netbox_data_source` | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | ⬜ | Not started |

---

## Recommended Implementation Order

### Phase 1: Core Infrastructure (High Priority) ✅ COMPLETE
1. ✅ `netbox_site` + `netbox_site_group` - Site hierarchy
2. ✅ `netbox_region` - Geographic hierarchy
3. ✅ `netbox_location` - Physical locations
4. ✅ `netbox_rack` + `netbox_rack_role` + `netbox_rack_type` - Rack infrastructure
5. ✅ `netbox_manufacturer` + `netbox_platform` - Hardware vendors
6. ✅ `netbox_device_type` - Device templates
7. ✅ `netbox_device_role` - Device classification
8. ✅ `netbox_device` - Physical devices
9. ✅ `netbox_interface` + `netbox_interface_template` - Network interfaces

### Phase 2: IPAM Essentials ✅ COMPLETE
10. ✅ `netbox_vrf` - Virtual routing
11. ✅ `netbox_prefix` - IP subnets
12. ✅ `netbox_ip_address` - IP addresses
13. ✅ `netbox_ip_range` - IP ranges
14. ✅ `netbox_vlan` + `netbox_vlan_group` - VLANs
15. ✅ `netbox_aggregate` + `netbox_rir` - Aggregates & RIRs
16. ✅ `netbox_role` - IPAM roles
17. ✅ `netbox_asn` - Autonomous System Numbers
18. ✅ `netbox_service` - Network services

### Phase 3: Virtualization ✅ COMPLETE
19. ✅ `netbox_cluster_type` + `netbox_cluster_group` - Cluster organization
20. ✅ `netbox_cluster` - VM clusters
21. ✅ `netbox_virtual_machine` - VMs
22. ✅ `netbox_vm_interface` - VM interfaces

### Phase 4: Circuits & Connectivity ✅ COMPLETE
23. ✅ `netbox_provider` + `netbox_provider_account` + `netbox_provider_network` - Circuit providers
24. ✅ `netbox_circuit` + `netbox_circuit_type` - WAN circuits
25. ✅ `netbox_circuit_termination` - Circuit terminations
26. ✅ `netbox_cable` - Physical cabling

### Phase 5: Tenancy & Contacts ✅ COMPLETE
27. ✅ `netbox_tenant` + `netbox_tenant_group` - Multi-tenancy
28. ✅ `netbox_contact` + `netbox_contact_group` + `netbox_contact_role` - Contact management

### Phase 6: Extras & Customization ✅ COMPLETE
29. ✅ `netbox_tag` - Tagging
30. ✅ `netbox_custom_field` - Custom fields
31. ✅ `netbox_config_context` + `netbox_config_template` - Configuration contexts
32. ✅ `netbox_webhook` - Automation hooks

### Phase 7: Advanced DCIM ✅ COMPLETE
33. ✅ `netbox_device_bay` + `netbox_virtual_chassis` - Chassis management
34. ✅ `netbox_module` + `netbox_module_type` + `netbox_module_bay` - Modular devices
35. ✅ `netbox_console_port` + `netbox_console_port_template` - Console ports
36. ✅ `netbox_console_server_port` + `netbox_console_server_port_template` - Console server ports
37. ✅ `netbox_power_panel` + `netbox_power_feed` - Power infrastructure
38. ✅ `netbox_power_port` + `netbox_power_port_template` - Power ports
39. ✅ `netbox_power_outlet` + `netbox_power_outlet_template` - Power outlets
40. ✅ `netbox_inventory_item` + `netbox_inventory_item_role` - Inventory management

### Phase 8: Wireless ✅ COMPLETE
41. ✅ `netbox_wireless_lan` + `netbox_wireless_lan_group` - Wireless networks

### Phase 9: VPN Resources ✅ COMPLETE
42. ✅ `netbox_ike_proposal` - IKE proposals for VPN
43. ✅ `netbox_ike_policy` - IKE policies for VPN
44. ✅ `netbox_ipsec_proposal` - IPSec proposals for VPN
45. ✅ `netbox_ipsec_policy` - IPSec policies for VPN
46. ✅ `netbox_ipsec_profile` - IPSec profiles for VPN
47. ✅ `netbox_tunnel_group` - VPN tunnel groups
48. ✅ `netbox_tunnel` - VPN tunnels
49. ✅ `netbox_tunnel_termination` - VPN tunnel terminations

### Phase 10: Recently Implemented ✅ COMPLETE
50. ✅ `netbox_route_target` - VRF route targets
51. ✅ `netbox_virtual_disk` - Virtual machine disks
52. ✅ `netbox_asn_range` - ASN ranges
53. ✅ `netbox_device_bay_template` - Device bay templates

### Phase 11: L2VPN ✅ COMPLETE
54. ✅ `netbox_l2vpn` - Layer 2 VPN
55. ✅ `netbox_l2vpn_termination` - L2VPN terminations

### Phase 12: Circuit Groups ✅ COMPLETE
56. ✅ `netbox_circuit_group` - Circuit grouping
57. ✅ `netbox_circuit_group_assignment` - Circuit group assignments

### Phase 13: Front/Rear Ports ✅ COMPLETE
58. ✅ `netbox_rear_port_template` - Rear port templates (device type definitions)
59. ✅ `netbox_front_port_template` - Front port templates (device type definitions)
60. ✅ `netbox_rear_port` - Device rear ports (physical patch panel connections)
61. ✅ `netbox_front_port` - Device front ports (mapped to rear ports)

### Phase 14: DCIM Templates & Infrastructure ✅ COMPLETE
62. ✅ `netbox_rack_reservation` - Rack unit reservations
63. ✅ `netbox_virtual_device_context` - Virtual device contexts (VDCs)
64. ✅ `netbox_module_bay_template` - Module bay templates for device types
65. ⬜ `netbox_cable_termination` - DEPRECATED: Use `netbox_cable` with embedded terminations (data source still available)
66. ✅ `netbox_inventory_item_template` - Inventory item templates

### Phase 15: Services, FHRP & Templates ✅ COMPLETE
67. ✅ `netbox_service_template` - Service templates (pre-defined service types)
68. ✅ `netbox_fhrp_group_assignment` - FHRP group to interface assignments
69. ✅ `netbox_export_template` - Export templates (Jinja2 data export)
70. ✅ `netbox_script` - Scripts data source (read-only, scripts managed via filesystem)

### Future Phases (Not Started)
- Contact assignments
- Wireless links
- Users/Groups/Permissions
- Notifications/Subscriptions
- Saved filters
- Image attachments
- Bookmarks
- Data sources

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
- **Resources:** Unit tests complete; Terraform integration tests are still in progress for some resources (see tables)
- **Data sources:** Unit tests complete; Terraform integration tests still needed where marked ⬜

---

_Last updated: December 11, 2025_
