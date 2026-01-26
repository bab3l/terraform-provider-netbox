package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIPRangeResource_basic(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(1, 50)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
		},
	})
}

func TestAccIPRangeResource_full(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(51, 100)
	thirdOctet := acctest.RandIntRange(51, 100)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)
	description := testutil.RandomName("ip-range-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_full(startAddress, endAddress, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "status", "active"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", description),
				),
			},
		},
	})
}

func TestAccIPRangeResource_update(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(101, 150)
	thirdOctet := acctest.RandIntRange(101, 150)
	startOctet2 := 10 + acctest.RandIntRange(1, 200)
	endOctet2 := startOctet2 + 10
	startAddress2 := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet2)
	endAddress2 := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet2)
	description := testutil.RandomName("ip-range-desc")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress2)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress2, endAddress2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress2),
				),
			},
			{
				Config: testAccIPRangeResourceConfig_full(startAddress2, endAddress2, description),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", description),
				),
			},
		},
	})
}

func TestAccIPRangeResource_import(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(151, 200)
	thirdOctet := acctest.RandIntRange(151, 200)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
			{
				ResourceName:      "netbox_ip_range.test",
				ImportState:       true,
				ImportStateVerify: true,
				// NetBox normalizes IP range endpoints to /32 on import; ignore to avoid false diffs.
				ImportStateVerifyIgnore: []string{
					"start_address",
					"end_address",
				},
			},
			{
				Config:   testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				PlanOnly: true,
			},
		},
	})
}

func TestAccIPRangeResource_tagLifecycle(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(51, 100)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")
	tag3Slug := testutil.RandomSlug("tag3")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)
	cleanup.RegisterTagCleanup(tag3Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_tags(startAddress, endAddress, tag1Slug, tag2Slug, tag3Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccIPRangeResourceConfig_tags(startAddress, endAddress, tag1Slug, tag2Slug, tag3Slug, caseTag1Uscore2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccIPRangeResourceConfig_tags(startAddress, endAddress, tag1Slug, tag2Slug, tag3Slug, caseTag3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "1"),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag3Slug),
				),
			},
			{
				Config: testAccIPRangeResourceConfig_tags(startAddress, endAddress, tag1Slug, tag2Slug, tag3Slug, tagsEmpty),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "0"),
				),
			},
		},
	})
}

func TestAccIPRangeResource_tagOrderInvariance(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(101, 150)
	thirdOctet := acctest.RandIntRange(51, 100)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	tag1Slug := testutil.RandomSlug("tag1")
	tag2Slug := testutil.RandomSlug("tag2")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)
	cleanup.RegisterTagCleanup(tag1Slug)
	cleanup.RegisterTagCleanup(tag2Slug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_tagsOrder(startAddress, endAddress, tag1Slug, tag2Slug, caseTag1Tag2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag2Slug),
				),
			},
			{
				Config: testAccIPRangeResourceConfig_tagsOrder(startAddress, endAddress, tag1Slug, tag2Slug, caseTag2Uscore1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_ip_range.test", "tags.#", "2"),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag1Slug),
					resource.TestCheckTypeSetElemAttr("netbox_ip_range.test", "tags.*", tag2Slug),
				),
			},
		},
	})
}

func TestAccIPRangeResource_externalDeletion(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(201, 250)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.IpamAPI.IpamIpRangesList(context.Background()).StartAddress([]string{startAddress}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find IP range for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.IpamAPI.IpamIpRangesDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete IP range: %v", err)
					}
					t.Logf("Successfully externally deleted IP range with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccIPRangeResourceConfig_basic(startAddress, endAddress string) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
}
`, startAddress, endAddress)
}

func testAccIPRangeResourceConfig_full(startAddress, endAddress, description string) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  status        = "active"
  description   = %[3]q
}
`, startAddress, endAddress, description)
}

func testAccIPRangeResourceConfig_tags(startAddress, endAddress, tag1Slug, tag2Slug, tag3Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag1Uscore2:
		tagsConfig = tagsDoubleSlug
	case caseTag3:
		tagsConfig = tagsSingleSlug
	case tagsEmpty:
		tagsConfig = tagsEmpty
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_tag" "tag3" {
  name = "Tag3-%[5]s"
  slug = %[5]q
}

resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  %[6]s
}
`, startAddress, endAddress, tag1Slug, tag2Slug, tag3Slug, tagsConfig)
}

func testAccIPRangeResourceConfig_tagsOrder(startAddress, endAddress, tag1Slug, tag2Slug, tagCase string) string {
	var tagsConfig string
	switch tagCase {
	case caseTag1Tag2:
		tagsConfig = tagsDoubleSlug
	case caseTag2Uscore1:
		tagsConfig = tagsDoubleSlugReversed
	}

	return fmt.Sprintf(`
resource "netbox_tag" "tag1" {
  name = "Tag1-%[3]s"
  slug = %[3]q
}

resource "netbox_tag" "tag2" {
  name = "Tag2-%[4]s"
  slug = %[4]q
}

resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  %[5]s
}
`, startAddress, endAddress, tag1Slug, tag2Slug, tagsConfig)
}

func TestAccConsistency_IPRange_LiteralNames(t *testing.T) {
	t.Parallel()

	startOctet := 50 + acctest.RandIntRange(1, 100)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("172.16.%d.10", startOctet)
	endAddress := fmt.Sprintf("172.16.%d.20", endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),
					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),
				),
			},
			{
				Config:   testAccIPRangeResourceConfig_basic(startAddress, endAddress),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),
				),
			},
		},
	})
}

// TestAccIPRangeResource_removeOptionalFields tests that removing previously set VRF, tenant, and role fields correctly sets them to null.
// This addresses the bug where removing a nullable reference field from the configuration would not clear it in NetBox,
// causing "Provider produced inconsistent result after apply" errors.
func TestAccIPRangeResource_removeOptionalFields(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(200, 254)
	thirdOctet := acctest.RandIntRange(1, 254)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	vrfName := testutil.RandomName("test-vrf-range-remove")
	vrfRD := fmt.Sprintf("65000:%d", acctest.RandIntRange(1000, 9999))
	tenantName := testutil.RandomName("test-tenant-range-remove")
	tenantSlug := testutil.GenerateSlug(tenantName)
	roleName := testutil.RandomName("test-role-range-remove")
	roleSlug := testutil.GenerateSlug(roleName)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)
	cleanup.RegisterVRFCleanup(vrfRD)
	cleanup.RegisterTenantCleanup(tenantSlug)
	cleanup.RegisterRoleCleanup(roleSlug)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Step 1: Create with VRF, tenant, and role
			{
				Config: testAccIPRangeResourceConfig_withAllFields(startAddress, endAddress, vrfName, vrfRD, tenantName, tenantSlug, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "vrf"),
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "role"),
				),
			},
			// Step 2: Remove VRF, tenant, and role - should set to null
			{
				Config: testAccIPRangeResourceConfig_withoutFields(startAddress, endAddress, vrfName, vrfRD, tenantName, tenantSlug, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("netbox_ip_range.test", "vrf"),
					resource.TestCheckNoResourceAttr("netbox_ip_range.test", "tenant"),
					resource.TestCheckNoResourceAttr("netbox_ip_range.test", "role"),
				),
			},
			// Step 3: Re-add VRF, tenant, and role - should work without errors
			{
				Config: testAccIPRangeResourceConfig_withAllFields(startAddress, endAddress, vrfName, vrfRD, tenantName, tenantSlug, roleName, roleSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "vrf"),
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "tenant"),
					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "role"),
				),
			},
		},
	})
}

func testAccIPRangeResourceConfig_withAllFields(startAddress, endAddress, vrfName, vrfRD, tenantName, tenantSlug, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %[3]q
  rd   = %[4]q
}

resource "netbox_tenant" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  vrf           = netbox_vrf.test.id
  tenant        = netbox_tenant.test.id
  role          = netbox_role.test.id
  status        = "active"
}
`, startAddress, endAddress, vrfName, vrfRD, tenantName, tenantSlug, roleName, roleSlug)
}

func testAccIPRangeResourceConfig_withoutFields(startAddress, endAddress, vrfName, vrfRD, tenantName, tenantSlug, roleName, roleSlug string) string {
	return fmt.Sprintf(`
resource "netbox_vrf" "test" {
  name = %[3]q
  rd   = %[4]q
}

resource "netbox_tenant" "test" {
  name = %[5]q
  slug = %[6]q
}

resource "netbox_role" "test" {
  name = %[7]q
  slug = %[8]q
}

resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  status        = "active"
  # vrf, tenant, and role removed - should set to null
}
`, startAddress, endAddress, vrfName, vrfRD, tenantName, tenantSlug, roleName, roleSlug)
}

func TestAccIPRangeResource_removeDescriptionAndComments(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(201, 254)
	thirdOctet := acctest.RandIntRange(201, 254)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	testutil.TestRemoveOptionalFields(t, testutil.MultiFieldOptionalTestConfig{
		ResourceName: "netbox_ip_range",
		BaseConfig: func() string {
			return testAccIPRangeResourceConfig_basic(startAddress, endAddress)
		},
		ConfigWithFields: func() string {
			return testAccIPRangeResourceConfig_withDescriptionAndComments(
				startAddress,
				endAddress,
				"Test description",
				"Test comments",
			)
		},
		OptionalFields: map[string]string{
			"description": "Test description",
			"comments":    "Test comments",
		},
		RequiredFields: map[string]string{
			"start_address": startAddress,
			"end_address":   endAddress,
		},
	})
}

func testAccIPRangeResourceConfig_withDescriptionAndComments(startAddress, endAddress, description, comments string) string {
	return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  status        = "active"
  description   = %[3]q
  comments      = %[4]q
}
`, startAddress, endAddress, description, comments)
}

func TestAccIPRangeResource_StatusOptionalComputed(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(1, 50)
	thirdOctet := acctest.RandIntRange(1, 50)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_ip_range",
		OptionalField:  "status",
		DefaultValue:   "active",
		FieldTestValue: "reserved",
		BaseConfig: func() string {
			return testAccIPRangeResourceConfig_basic(startAddress, endAddress)
		},
		WithFieldConfig: func(value string) string {
			return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address = %[1]q
  end_address   = %[2]q
  status        = %[3]q
}
`, startAddress, endAddress, value)
		},
	})
}

func TestAccIPRangeResource_MarkUtilizedOptionalComputed(t *testing.T) {
	t.Parallel()

	secondOctet := acctest.RandIntRange(51, 100)
	thirdOctet := acctest.RandIntRange(51, 100)
	startOctet := 10 + acctest.RandIntRange(1, 200)
	endOctet := startOctet + 10
	startAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, startOctet)
	endAddress := fmt.Sprintf("10.%d.%d.%d", secondOctet, thirdOctet, endOctet)

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterIPRangeCleanup(startAddress)

	testutil.RunOptionalComputedFieldTestSuite(t, testutil.OptionalComputedFieldTestConfig{
		ResourceName:   "netbox_ip_range",
		OptionalField:  "mark_utilized",
		DefaultValue:   "false",
		FieldTestValue: "true",
		BaseConfig: func() string {
			return testAccIPRangeResourceConfig_basic(startAddress, endAddress)
		},
		WithFieldConfig: func(value string) string {
			return fmt.Sprintf(`
resource "netbox_ip_range" "test" {
  start_address  = %[1]q
  end_address    = %[2]q
  mark_utilized  = %[3]s
}
`, startAddress, endAddress, value)
		},
	})
}

func TestAccIPRangeResource_validationErrors(t *testing.T) {
	testutil.RunMultiValidationErrorTest(t, testutil.MultiValidationErrorTestConfig{
		ResourceName: "netbox_ip_range",
		TestCases: map[string]testutil.ValidationErrorCase{
			"missing_start_address": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_ip_range" "test" {
  # start_address missing
  end_address = "192.168.1.20/24"
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
			"missing_end_address": {
				Config: func() string {
					return `
provider "netbox" {}

resource "netbox_ip_range" "test" {
  start_address = "192.168.1.10/24"
  # end_address missing
}
`
				},
				ExpectedError: testutil.ErrPatternRequired,
			},
		},
	})
}
