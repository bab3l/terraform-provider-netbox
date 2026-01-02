package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVLANResource_basic(t *testing.T) {

	t.Parallel()
	name := testutil.RandomName("tf-test-vlan")
	vid := testutil.RandomVID()
	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANDestroy,

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
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANDestroy,

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
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

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
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANDestroy,

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

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceConfig_basic(name, vid),
			},
			{
				ResourceName:      "netbox_vlan.test",
				ImportState:       true,
				ImportStateVerify: true,
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVLANDestroy,
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
	vlanVid := 100
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")
	groupName := testutil.RandomName("group")
	groupSlug := testutil.RandomSlug("group")
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")
	roleName := testutil.RandomName("role")
	roleSlug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{
			{
				Config: testAccVLANConsistencyConfig(vlanName, vlanVid, siteName, siteSlug, groupName, groupSlug, tenantName, tenantSlug, roleName, roleSlug),

				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("netbox_vlan.test", "site", siteName),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "group"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "tenant", tenantName),
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "role"),
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
  site = netbox_site.test.name
  group = netbox_vlan_group.test.id
  tenant = netbox_tenant.test.name
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVLANDestroy,
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
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

func TestAccVLANResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

	vlanName := testutil.RandomName("vlan")
	vid := testutil.RandomVID()
	tenantName := testutil.RandomName("tenant")
	tenantSlug := testutil.RandomSlug("tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterVLANCleanup(vid)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		CheckDestroy: testutil.CheckVLANDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVLANResourceImportConfig_full(vlanName, int(vid), tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_vlan.test", "id"),
					resource.TestCheckResourceAttr("netbox_vlan.test", "name", vlanName),
					resource.TestCheckResourceAttr("netbox_vlan.test", "vid", fmt.Sprintf("%d", int(vid))),
					// Verify custom fields are applied
					resource.TestCheckResourceAttr("netbox_vlan.test", "custom_fields.#", "7"),
					// Verify tags are applied
					resource.TestCheckResourceAttr("netbox_vlan.test", "tags.#", "2"),
				),
			},
			{
				ResourceName:            "netbox_vlan.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"tenant", "custom_fields"}, // Tenant may have lookup inconsistencies, custom fields have import limitations
			},
		},
	})
}

func testAccVLANResourceImportConfig_full(vlanName string, vid int, tenantName, tenantSlug string) string {
	// Custom field names with underscore format
	cfText := testutil.RandomCustomFieldName("cf_text")
	cfLongtext := testutil.RandomCustomFieldName("cf_longtext")
	cfInteger := testutil.RandomCustomFieldName("cf_integer")
	cfBoolean := testutil.RandomCustomFieldName("cf_boolean")
	cfDate := testutil.RandomCustomFieldName("cf_date")
	cfUrl := testutil.RandomCustomFieldName("cf_url")
	cfJson := testutil.RandomCustomFieldName("cf_json")

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
}

# Custom Fields (all supported data types)
resource "netbox_custom_field" "cf_text" {
  name         = %q
  type         = "text"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_longtext" {
  name         = %q
  type         = "longtext"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_integer" {
  name         = %q
  type         = "integer"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_boolean" {
  name         = %q
  type         = "boolean"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_date" {
  name         = %q
  type         = "date"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_url" {
  name         = %q
  type         = "url"
  object_types = ["ipam.vlan"]
}

resource "netbox_custom_field" "cf_json" {
  name         = %q
  type         = "json"
  object_types = ["ipam.vlan"]
}

# Tags
resource "netbox_tag" "tag1" {
  name = %q
  slug = %q
}

resource "netbox_tag" "tag2" {
  name = %q
  slug = %q
}

# VLAN with comprehensive custom fields and tags
resource "netbox_vlan" "test" {
  name   = %q
  vid    = %d
  tenant = netbox_tenant.test.id

  custom_fields = [
    {
      name  = netbox_custom_field.cf_text.name
      type  = "text"
      value = "test text value"
    },
    {
      name  = netbox_custom_field.cf_longtext.name
      type  = "longtext"
      value = "This is a much longer text value that spans multiple lines and contains more detailed information about this VLAN resource for testing purposes."
    },
    {
      name  = netbox_custom_field.cf_integer.name
      type  = "integer"
      value = "42"
    },
    {
      name  = netbox_custom_field.cf_boolean.name
      type  = "boolean"
      value = "true"
    },
    {
      name  = netbox_custom_field.cf_date.name
      type  = "date"
      value = "2023-01-15"
    },
    {
      name  = netbox_custom_field.cf_url.name
      type  = "url"
      value = "https://example.com"
    },
    {
      name  = netbox_custom_field.cf_json.name
      type  = "json"
      value = jsonencode({"key": "value"})
    }
  ]

  tags = [
    {
      name = netbox_tag.tag1.name
      slug = netbox_tag.tag1.slug
    },
    {
      name = netbox_tag.tag2.name
      slug = netbox_tag.tag2.slug
    }
  ]
}
`, tenantName, tenantSlug, cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfUrl, cfJson, tag1, tag1Slug, tag2, tag2Slug, vlanName, vid)
}
