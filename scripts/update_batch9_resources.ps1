# Script to update Batch 9 resources with partial custom fields management pattern
# Updates: ike_proposal, ipsec_policy, ipsec_profile, ipsec_proposal

$resources = @(
    @{name="IKEProposal"; file="ike_proposal_resource.go"; var="ike"}
    @{name="IPSecPolicy"; file="ipsec_policy_resource.go"; var="ipsec"}
    @{name="IPSecProfile"; file="ipsec_profile_resource.go"; var="ipsec"}
    @{name="IPSecProposal"; file="ipsec_proposal_resource.go"; var="ipsec"}
)

foreach ($res in $resources) {
    $file = "internal\resources\$($res.file)"
    Write-Host "`nProcessing $($res.name)..." -ForegroundColor Cyan

    if (-not (Test-Path $file)) {
        Write-Host "  File not found: $file" -ForegroundColor Red
        continue
    }

    $content = Get-Content $file -Raw
    $original = $content

    # 1. Update Update() method signature
    $content = $content -replace `
        '(func \(r \*' + $res.name + 'Resource\) Update\(.*?\) \{[\r\n\s]+)var data ' + $res.name + 'ResourceModel([\r\n\s]+)\/\/ Read Terraform plan data', `
        "`$1var state, plan $($res.name)ResourceModel`$2// Read both state and plan for merge-aware custom fields handling`${2}resp.Diagnostics.Append(req.State.Get(ctx, &state)...)`$2// Read Terraform plan"

    # 2. Update Plan.Get in Update()
    $content = $content -replace `
        '(// Read Terraform plan data[^\n]*\n\s+)resp\.Diagnostics\.Append\(req\.Plan\.Get\(ctx, &data\)', `
        "`$1resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)"

    # 3. Replace data. with plan. in Update() method (ID, Name references)
    $content = $content -replace `
        '(func \(r \*' + $res.name + 'Resource\) Update\([^}]+)data\.ID\.ValueString\(\)', 'plan.ID.ValueString()'
    $content = $content -replace `
        '(func \(r \*' + $res.name + 'Resource\) Update\([^}]+)data\.Name\.ValueString\(\)', 'plan.Name.ValueString()'

    # 4. Update setOptionalFields call in Update() to pass state
    $content = $content -replace `
        '(\s+r\.setOptionalFields\(ctx, [^,]+, &)data(, &resp\.Diagnostics\))', `
        "`$1plan, &state`$2"

    # 5. Add filter-to-owned logic before final Set in Update()
    $pattern = '(\s+)(\/\/ Map response to model[\r\n\s]+r\.map' + $res.name + 'ToState\(ctx, ' + $res.var + ', &)data(, &resp\.Diagnostics\)[\r\n\s]+)(tflog\.Debug.+?[\r\n\s]+}[\r\n\s]+)(\/\/ Save updated data.+?[\r\n\s]+resp\.Diagnostics\.Append\(resp\.State\.Set\(ctx, &)data(\)\.\.\.)'

    $replacement = "`$1// Save the plan's custom fields before mapping (for filter-to-owned pattern)`${1}planCustomFields := plan.CustomFields`${1}`${1}// Map response to model`${1}r.map$($res.name)ToState(ctx, $($res.var), &plan, &resp.Diagnostics)`${1}if resp.Diagnostics.HasError() {`${1}`${1}return`${1}}`${1}`${1}// Apply filter-to-owned pattern for custom fields`${1}plan.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, planCustomFields, $($res.var).GetCustomFields(), &resp.Diagnostics)`${1}`${1}`$4`$5plan`$6"

    $content = $content -replace $pattern, $replacement

    # 6. Update setOptionalFields signature
    $content = $content -replace `
        '(func \(r \*' + $res.name + 'Resource\) setOptionalFields\(ctx context\.Context, [^,]+, )data (\*' + $res.name + 'ResourceModel, diags)', `
        "`${1}plan `$2, state *$($res.name)ResourceModel, diags"

    # 7. Replace ApplyCustomFields with merge-aware version in setOptionalFields
    $pattern = '(\s+)(\/\/ Set comments, tags, and custom fields[\r\n\s]+utils\.ApplyComments[^\n]+[\r\n\s]+utils\.ApplyTags[^\n]+[\r\n\s]+)utils\.ApplyCustomFields\(ctx, [^,]+, data\.CustomFields, diags\)'

    $replacement = "`$1// Set comments, tags, and custom fields with merge-aware helpers`${1}utils.ApplyComments(request, plan.Comments)`${1}utils.ApplyTags(ctx, request, plan.Tags, diags)`${1}// Apply custom fields with merge logic to preserve unmanaged fields`${1}if state != nil {`${1}`${1}utils.ApplyCustomFieldsWithMerge(ctx, request, plan.CustomFields, state.CustomFields, diags)`${1}} else {`${1}`${1}// During Create, no state exists yet`${1}`${1}utils.ApplyCustomFields(ctx, request, plan.CustomFields, diags)`${1}}"

    $content = $content -replace $pattern, $replacement

    # 8. Update Create() to pass nil state
    $content = $content -replace `
        '(\s+r\.setOptionalFields\(ctx, [^,]+, &data)(, &resp\.Diagnostics\))', `
        "`${1}, nil`$2"

    # 9. Update Read() to preserve null/empty custom fields
    $pattern = '(\s+)(\/\/ Map response to model[\r\n\s]+r\.map' + $res.name + 'ToState\(ctx, ' + $res.var + ', &data, &resp\.Diagnostics\)[\r\n\s]+)(\/\/ Save updated data)'

    $replacement = "`$1// Preserve original custom_fields value from state`${1}originalCustomFields := data.CustomFields`${1}`${1}`$2if resp.Diagnostics.HasError() {`${1}`${1}return`${1}}`${1}`${1}// If custom_fields was null or empty before, restore that state`${1}// This prevents drift when config doesn't declare custom_fields`${1}if originalCustomFields.IsNull() || (utils.IsSet(originalCustomFields) && len(originalCustomFields.Elements()) == 0) {`${1}`${1}data.CustomFields = originalCustomFields`${1}}`${1}`${1}`$3"

    $content = $content -replace $pattern, $replacement

    # 10. Update mapXXXToState to use filter-to-owned
    $content = $content -replace `
        'data\.CustomFields = utils\.PopulateCustomFieldsFromAPI\(ctx, ' + $res.var + '\.HasCustomFields\(\), ' + $res.var + '\.GetCustomFields\(\), data\.CustomFields, diags\)', `
        "data.CustomFields = utils.PopulateCustomFieldsFilteredToOwned(ctx, data.CustomFields, $($res.var).GetCustomFields(), diags)"

    # Write the updated content
    if ($content -ne $original) {
        Set-Content -Path $file -Value $content -NoNewline
        Write-Host "  Updated $($res.file)" -ForegroundColor Green
    } else {
        Write-Host "  No changes needed for $($res.file)" -ForegroundColor Yellow
    }
}

Write-Host "`nDone! Now testing compilation..." -ForegroundColor Cyan
foreach ($res in $resources) {
    $file = "internal\resources\$($res.file)"
    Write-Host "Testing $($res.file)..." -NoNewline
    $result = go build "./$file" 2>&1
    if ($LASTEXITCODE -eq 0) {
        Write-Host " OK" -ForegroundColor Green
    } else {
        Write-Host " FAILED" -ForegroundColor Red
        Write-Host $result
    }
}
