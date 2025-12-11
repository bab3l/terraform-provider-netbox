# TODO: Integration with go-netbox

This file outlines the next steps to fully integrate your go-netbox wrapper with the Terraform provider.

## 1. Understand go-netbox Structure

First, you need to determine the correct import path and client structure in your go-netbox library:

```bash
# Navigate to your go-netbox directory
cd ../go-netbox

# Check the available packages and their structure
find . -name "*.go" | head -10
grep -r "type.*Client" . --include="*.go"
grep -r "func.*New" . --include="*.go"
```

## 2. Update Provider Configuration

Once you know the correct import path and client constructor, update `internal/provider/provider.go`:

```go
import (
    // ... other imports
    "github.com/bab3l/go-netbox/netbox" // Update this import path
)

func (p *NetboxProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    // ... existing validation code ...

    // Replace the TODO section with actual client creation:
    client, err := netbox.NewClient(serverURL, apiToken)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create Netbox API Client",
            "An unexpected error occurred when creating the Netbox API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "Netbox Client Error: "+err.Error(),
        )
        return
    }

    if insecure {
        client.SetInsecureSkipVerify(true)
    }

    // Pass the actual client to resources and data sources
    resp.DataSourceData = client
    resp.ResourceData = client
}
```

## 3. Update Site Resource

Update `internal/resources/site_resource.go` to use the actual client:

```go
import (
    // ... other imports
    "github.com/bab3l/go-netbox/netbox" // Update this import
)

type SiteResource struct {
    client *netbox.Client // Replace map[string]interface{} with actual client type
}

func (r *SiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }

    client, ok := req.ProviderData.(*netbox.Client) // Update type assertion

    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected Resource Configure Type",
            fmt.Sprintf("Expected *netbox.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )
        return
    }

    r.client = client
}

func (r *SiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var data SiteResourceModel

    resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Use actual go-netbox API calls:
    siteRequest := &netbox.SiteRequest{
        Name: data.Name.ValueString(),
        Slug: data.Slug.ValueString(),
        // ... map other fields
    }

    site, err := r.client.CreateSite(ctx, siteRequest)
    if err != nil {
        resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create site, got error: %s", err))
        return
    }

    // Map response back to Terraform state:
    data.ID = types.StringValue(strconv.Itoa(site.ID))
    data.Name = types.StringValue(site.Name)
    // ... map other fields

    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Similarly update Read, Update, Delete methods...
```

## 4. Common Patterns for go-netbox Integration

### Error Handling
```go
if err != nil {
    resp.Diagnostics.AddError(
        "Client Error",
        fmt.Sprintf("Unable to %s site, got error: %s", operation, err),
    )
    return
}
```

### ID Conversion
```go
// String to int for API calls
siteID, err := strconv.Atoi(data.ID.ValueString())
if err != nil {
    resp.Diagnostics.AddError("Parse Error", fmt.Sprintf("Unable to parse site ID: %s", err))
    return
}

// Int to string for Terraform state
data.ID = types.StringValue(strconv.Itoa(site.ID))
```

### Null Value Handling
```go
// Check if optional field is set before including in API request
if !data.Description.IsNull() {
    request.Description = data.Description.ValueStringPointer()
}

// Set computed/optional fields from API response
if site.Description != nil {
    data.Description = types.StringValue(*site.Description)
} else {
    data.Description = types.StringNull()
}
```

## 5. Additional Resources to Implement

Based on common Netbox usage, consider implementing these resources next:

- `netbox_device_type`
- `netbox_device_role`
- `netbox_device`
- `netbox_ip_address`
- `netbox_prefix`
- `netbox_vlan`
- `netbox_tenant`
- `netbox_region`

## 6. Data Sources

Implement corresponding data sources for read-only access:

- `netbox_site` (data source)
- `netbox_device` (data source)
- etc.

## 7. Testing with Real Netbox Instance

Create acceptance tests that use a real Netbox instance:

```go
func TestAccSiteResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck:                 func() { testAccPreCheck(t) },
        ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: testAccSiteResourceConfig("test-site"),
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("netbox_site.test", "name", "test-site"),
                    resource.TestCheckResourceAttr("netbox_site.test", "slug", "test-site"),
                ),
            },
        },
    })
}
```

## 8. Documentation Generation

After implementing resources, generate documentation:

```bash
go generate ./...
```

This will create documentation in the `docs/` directory based on your resource schemas.
