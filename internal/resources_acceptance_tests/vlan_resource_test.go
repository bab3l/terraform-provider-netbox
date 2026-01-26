package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANResource_basic(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vlan")
	vid := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceConfig_basic(name, vid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),
				),
			},
		},
	})
}

func TestAccVLANResource_full(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vlan-full")
	vid := testutil.RandomVID()
	description := "Test VLAN with all fields"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceConfig_full(name, vid, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),
					resource.TestCheckResourceAttr("netbox_vlan.test", "description", description),
					resource.TestCheckResourceAttr("netbox_vlan.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccVLANResource_withGroup(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vlan-grp")
	vid := testutil.RandomVID()
	groupName := testutil.RandomName("tf-test-vlangrp")
	groupSlug := testutil.GenerateSlug("tf-test-vlangrp")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)
	cleanup.RegisterVLANGroupCleanup(groupSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
			testutil.CheckVLANGroupDestroy,
		),
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceConfig_withGroup(name, vid, groupName, groupSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "group"),
				),
			},
		},
	})
}

func TestAccVLANResource_update(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vlan-upd")
	updatedName := testutil.RandomName("tf-test-vlan-updated")
	vid := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceConfig_basic(name, vid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),
				),
			},
			{
				Config: testAccVLANResourceConfig_full(updatedName, vid, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", updatedName),
					resource.TestCheckResourceAttr("netbox_vlan.test", "description", "Updated description"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "status", "active"),
				),
			},
		},
	})
}

func TestAccVLANResource_import(t *testing.T) {
	t.Parallel()

	name := "test-vlan-" + testutil.GenerateSlug("vlan")
	vid := int32(100)
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterVLANGroupCleanup(groupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANConsistencyConfig(name, int(vid), siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),
			},
			{
				ResourceName:      "netbox_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"site",
					"group",
					"tenant",
					"role",
				},
				Check: resource.ComposeTestCheckFunc(
					testutil.ReferenceFieldCheck("netbox_vlan.test", "site"),
					testutil.ReferenceFieldCheck("netbox_vlan.test", "group"),
					testutil.ReferenceFieldCheck("netbox_vlan.test", "tenant"),
					testutil.ReferenceFieldCheck("netbox_vlan.test", "role"),
				),
			},
			{
				Config:   testAccVLANConsistencyConfig(name, int(vid), siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccVLANResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("tf-test-vlan-id")
	vid := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceConfig_basic(name, vid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),
				),
			},
		},
	})
}

func testAccVLANResourceConfig_basic(name string, vid int32) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
}
`, name, vid)
}

func testAccVLANResourceConfig_full(name string, vid int32, description string) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
  name        = %q
  vid         = %d
  description = %q
  status      = "active"
}
`, name, vid, description)
}

func testAccVLANResourceConfig_withGroup(name string, vid int32, groupName, groupSlug string) string {
	return fmt.Sprintf(`
resource "netbox_vlan_group" "test" {
  name = %q
  slug = %q
}

resource "netbox_vlan" "test" {
  name  = %q
  vid   = %d
  group = netbox_vlan_group.test.id
}
`, groupName, groupSlug, name, vid)
}

func TestAccConsistency_VLAN(t *testing.T) {
	t.Parallel()

	vlanName := testutil.RandomName("vlan")
	vlanVid := int(testutil.RandomVID())
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterVLANGroupCleanup(groupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccVLANConsistencyConfig(vlanName, vlanVid, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttrPair("netbox_vlan.test", "site", "netbox_site.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_vlan.test", "group", "netbox_vlan_group.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_vlan.test", "tenant", "netbox_tenant.test", "id"),
					resource.TestCheckResourceAttrPair("netbox_vlan.test", "role", "netbox_role.test", "id"),
				),
			},
			{
				PlanOnly: true,

				Config: testAccVLANConsistencyConfig(vlanName, vlanVid, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),
			},
		},
	})
}

func testAccVLANConsistencyConfig(vlanName string, vlanVid int, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_vlan_group" "test" {
  name = "%[5]s"
  slug = "%[6]s"
  scope_type = "dcim.site"
  scope_id = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = "%[7]s"
  slug = "%[8]s"
}

resource "netbox_role" "test" {
  name = "%[9]s"
  slug = "%[10]s"
}

resource "netbox_vlan" "test" {
  name = "%[1]s"
  vid  = %[2]d
	site = netbox_site.test.id
  group = netbox_vlan_group.test.id
	tenant = netbox_tenant.test.id
  role = netbox_role.test.id
}
`, vlanName, vlanVid, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug)
}

func TestAccVLANResource_optionalRoleNoUpdate(t *testing.T) {
	t.Parallel()

	siteName := testutil.RandomName("tf-test-site-vlan-role")
	siteSlug := testutil.RandomSlug("tf-test-site-vlan-role")
	vlanName := testutil.RandomName("tf-test-vlan-role")
	vlanVid := testutil.RandomVID()
	description1 := "Initial description"
	description2 := "Updated description"

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vlanVid)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
			testutil.CheckSiteDestroy,
		),
		Steps: []resource.TestStep{
			{
				// Create VLAN without role
				Config: testAccVLANOptionalRoleConfig(siteName, siteSlug, vlanName, vlanVid, description1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVid)),
					resource.TestCheckResourceAttr("netbox_vlan.test", "description", description1),
					resource.TestCheckNoResourceAttr("netbox_vlan.test", "role"),
				),
			},
			{
				// Update description (not role) - role should remain empty/null
				Config: testAccVLANOptionalRoleConfig(siteName, siteSlug, vlanName, vlanVid, description2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vlanVid)),
					resource.TestCheckResourceAttr("netbox_vlan.test", "description", description2),
					resource.TestCheckNoResourceAttr("netbox_vlan.test", "role"),
				),
			},
		},
	})
}

func testAccVLANOptionalRoleConfig(siteName, siteSlug, vlanName string, vlanVid int32, description string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_vlan" "test" {
  name        = %q
  vid         = %d
  site        = netbox_site.test.id
  description = %q
  # role intentionally omitted to test optional attribute handling
}
`, siteName, siteSlug, vlanName, vlanVid, description)
}

func TestAccConsistency_VLAN_LiteralNames(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-vlan-lit")
	vid := testutil.RandomVID()
	description := testutil.RandomName("description")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANConsistencyLiteralNamesConfig(name, vid, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", name),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", vid)),
					resource.TestCheckResourceAttr("netbox_vlan.test", "description", description),
				),
			},
			{
				Config:   testAccVLANConsistencyLiteralNamesConfig(name, vid, description),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
				),
			},
		},
	})
}

func testAccVLANConsistencyLiteralNamesConfig(name string, vid int32, description string) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
  name        = %q
  vid         = %d
  description = %q
}
`, name, vid, description)
}

func TestAccVLANResource_externalDeletion(t *testing.T) {
	t.Parallel()
	name := testutil.RandomName("tf-test-vlan-ext-del")
	vid := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceConfig_basic(name, vid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamVlansList(context.Background()).NameIc([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find VLAN for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamVlansDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete VLAN: %v", err)
					}
					t.Logf("Successfully externally deleted VLAN with ID: %d", itemID)
				},
				Config: testAccVLANResourceConfig_basic(name, vid),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
				),
			},
		},
	})
}

// NOTE: Custom field tests for VLAN resource are in resources_acceptance_tests_customfields package

// TestAccVlanResource_StatusOptionalField tests comprehensive scenarios for VLAN status.
// This validates that Optional+Computed fields work correctly across all scenarios.
func TestAccVlanResource_StatusOptionalField(t *testing.T) {
	t.Parallel()

	// Generate unique names for this test run
	vlanName := testutil.RandomName("tf-test-vlan-status")
	vid := testutil.RandomVID()

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_vlan",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "deprecated",
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckVLANDestroy,
			testutil.CheckVLANGroupDestroy,
			testutil.CheckSiteDestroy,
		),
		BaseConfig: func() string {
			return testAccVLANResourceConfig_statusBase(vlanName, vid)
		},
		WithFieldConfig: func(value string) string {
			return testAccVLANResourceConfig_statusWithField(vlanName, vid, value)
		},
	})
}

func testAccVLANResourceConfig_statusBase(name string, vid int32) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
	name = %[1]q
	vid  = %[2]d
	# status field intentionally omitted - should get default "active"
}
`, name, vid)
}

func testAccVLANResourceConfig_statusWithField(name string, vid int32, status string) string {
	return fmt.Sprintf(`
resource "netbox_vlan" "test" {
	name   = %[1]q
	vid    = %[2]d
	status = %[3]q
}
`, name, vid, status)
}

// TestAccVLANResource_removeOptionalFields tests that optional nullable fields
// can be successfully removed from the configuration without causing inconsistent state.
// This verifies the bugfix for: "Provider produced inconsistent result after apply".
func TestAccVLANResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	vlanName := testutil.RandomName("tf-test-vlan-remove")
	vid := testutil.RandomVID()
	siteName := testutil.RandomName("test-site-vlan")
	siteSlug := testutil.RandomSlug("test-site-vlan")
	groupName := testutil.RandomName("test-group-vlan")
	groupSlug := testutil.RandomSlug("test-group-vlan")
	tenantName := testutil.RandomName("test-tenant-vlan")
	tenantSlug := testutil.RandomSlug("test-tenant-vlan")
	roleName := testutil.RandomName("test-role-vlan")
	roleSlug := testutil.RandomSlug("test-role-vlan")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterVLANGroupCleanup(groupSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		CheckDestroy:             testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			// Step 1: Create VLAN with all nullable fields
			{
				Config: testAccVLANResourceConfig_withAllFields(vlanName, vid, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "group"),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "role"),
				),
			},
			// Step 2: Remove all nullable fields - should set them to null
			{
				Config: testAccVLANResourceConfig_withoutFields(vlanName, vid, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckNoResourceAttr("netbox_vlan.test", "site"),
					resource.TestCheckNoResourceAttr("netbox_vlan.test", "group"),
					resource.TestCheckNoResourceAttr("netbox_vlan.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_vlan.test", "role"),
				),
			},
			// Step 3: Re-add all fields - verify they can be set again
			{
				Config: testAccVLANResourceConfig_withAllFields(vlanName, vid, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "site"),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "group"),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "role"),
				),
			},
		},
	})
}

func testAccVLANResourceConfig_withAllFields(vlanName string, vid int32, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_vlan_group" "test" {
  name       = %q
  slug       = %q
  scope_type = "dcim.site"
  scope_id   = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_vlan" "test" {
  name   = %q
  vid    = %d
  site   = netbox_site.test.id
  group  = netbox_vlan_group.test.id
  tenant = netbox_tenant.test.id
  role   = netbox_role.test.id
}
`, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug, vlanName, vid)
}

func testAccVLANResourceConfig_withoutFields(vlanName string, vid int32, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_vlan_group" "test" {
  name       = %q
  slug       = %q
  scope_type = "dcim.site"
  scope_id   = netbox_site.test.id
}

resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

resource "netbox_role" "test" {
  name = %q
  slug = %q
}

resource "netbox_vlan" "test" {
  name = %q
  vid  = %d
}
`, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug, vlanName, vid)
}

func TestAccVLANResource_validationErrors(t *testing.T) {
	t.Parallel()

	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_vlan",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_vid": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  name = "Test VLAN"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_name": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  vid = 100
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"vid_too_low": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  name = "Test VLAN"
  vid  = 0
}
`
				},
				ExpectedError: testutil.ErrPatternRange,
			},
			"vid_too_high": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  name = "Test VLAN"
  vid  = 5000
}
`
				},
				ExpectedError: testutil.ErrPatternRange,
			},
			"invalid_status": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  name   = "Test VLAN"
  vid    = 100
  status = "invalid_status"
}
`
				},
				ExpectedError: testutil.ErrPatternInvalidEnum,
			},
			"invalid_site_reference": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  name = "Test VLAN"
  vid  = 100
  site = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_group_reference": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  name  = "Test VLAN"
  vid   = 100
  group = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
			"invalid_tenant_reference": {
				Config: func() string {
					return `
resource "netbox_vlan" "test" {
  name   = "Test VLAN"
  vid    = 100
  tenant = "99999999"
}
`
				},
				ExpectedError: testutil.ErrPatternNotFound,
			},
		},
		CheckDestroy: testutil.CheckVLANDestroy,
	})
}
