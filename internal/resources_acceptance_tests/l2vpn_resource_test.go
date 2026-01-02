package resources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccL2VPNResource_basic(t *testing.T) {

	t.Parallel()
	name := acctest.RandomWithPrefix("test-l2vpn")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_updated(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name+"-updated"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Updated description"),
				),
			},
			{
				ResourceName:            "netbox_l2vpn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name"}, // display_name is computed and may differ after name changes
			},
		},
	})
}

func TestAccL2VPNResource_full(t *testing.T) {

	t.Parallel()
	name := acctest.RandomWithPrefix("test-l2vpn")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_full(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vpls"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "identifier", "12345"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", "Test L2VPN"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "comments", "Test comments"),
				),
			},
		},
	})
}

func TestAccL2vpnResource_importWithCustomFieldsAndTags(t *testing.T) {
	t.Parallel()

	name := testutil.RandomName("test-l2vpn")
	slug := testutil.RandomSlug("test-l2vpn")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2vpnResourceImportConfig_full(name, slug, tenantName, tenantSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", slug),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				ResourceName:            "netbox_l2vpn.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"display_name", "tenant", "custom_fields", "tags"},
			},
		},
	})
}

func testAccL2vpnResourceImportConfig_full(name, slug, tenantName, tenantSlug string) string {
	// Generate test data for all custom field types
	textValue := testutil.RandomName("text-value")
	longtextValue := testutil.RandomName("longtext-value") + "\nThis is a multiline text field for comprehensive testing."
	intValue := 42 // Fixed value for reproducibility
	boolValue := true
	dateValue := testutil.RandomDate()
	urlValue := testutil.RandomURL("test-url")
	jsonValue := testutil.RandomJSON()

	// Tag names
	tag1 := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2 := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	// Custom field names
	cfText := testutil.RandomCustomFieldName("tf_text")
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfInteger := testutil.RandomCustomFieldName("tf_integer")
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfJSON := testutil.RandomCustomFieldName("tf_json")

	return fmt.Sprintf(`
# Dependencies
resource "netbox_tenant" "test" {
  name = %q
  slug = %q
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

# Custom Fields for vpn.l2vpn object type
resource "netbox_custom_field" "test_text" {
  name         = %q
  label        = "Test Text CF"
  type         = "text"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_longtext" {
  name         = %q
  label        = "Test Longtext CF"
  type         = "longtext"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_integer" {
  name         = %q
  label        = "Test Integer CF"
  type         = "integer"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_boolean" {
  name         = %q
  label        = "Test Boolean CF"
  type         = "boolean"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_date" {
  name         = %q
  label        = "Test Date CF"
  type         = "date"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_url" {
  name         = %q
  label        = "Test URL CF"
  type         = "url"
  object_types = ["vpn.l2vpn"]
}

resource "netbox_custom_field" "test_json" {
  name         = %q
  label        = "Test JSON CF"
  type         = "json"
  object_types = ["vpn.l2vpn"]
}

# L2VPN with comprehensive custom fields and tags
resource "netbox_l2vpn" "test" {
  name   = %q
  slug   = %q
  type   = "vxlan"
  tenant = netbox_tenant.test.id

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

  custom_fields = [
    {
      name  = netbox_custom_field.test_text.name
      type  = "text"
      value = %q
    },
    {
      name  = netbox_custom_field.test_longtext.name
      type  = "longtext"
      value = %q
    },
    {
      name  = netbox_custom_field.test_integer.name
      type  = "integer"
      value = "%d"
    },
    {
      name  = netbox_custom_field.test_boolean.name
      type  = "boolean"
      value = "%t"
    },
    {
      name  = netbox_custom_field.test_date.name
      type  = "date"
      value = %q
    },
    {
      name  = netbox_custom_field.test_url.name
      type  = "url"
      value = %q
    },
    {
      name  = netbox_custom_field.test_json.name
      type  = "json"
      value = %q
    },
  ]

  depends_on = [
    netbox_custom_field.test_text,
    netbox_custom_field.test_longtext,
    netbox_custom_field.test_integer,
    netbox_custom_field.test_boolean,
    netbox_custom_field.test_date,
    netbox_custom_field.test_url,
    netbox_custom_field.test_json,
  ]
}
`, tenantName, tenantSlug,
		tag1, tag1Slug, tag2, tag2Slug,
		cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
		name, slug, textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue)
}

func testAccL2VPNResourceConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}
`, name, name)
}

func testAccL2VPNResourceConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = "Updated description"
}
`, name+"-updated", name)
}

func testAccL2VPNResourceConfig_full(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vpls"
  identifier  = 12345
  description = "Test L2VPN"
  comments    = "Test comments"
}
`, name, name)
}

func TestAccL2VPNResource_IDPreservation(t *testing.T) {
	t.Parallel()

	name := acctest.RandomWithPrefix("test-l2vpn-id")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
		},
	})
}

func TestAccConsistency_L2VPN_LiteralNames(t *testing.T) {
	t.Parallel()
	name := "test-l2vpn-lit"

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNConsistencyLiteralNamesConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				Config:   testAccL2VPNConsistencyLiteralNamesConfig(name),
				PlanOnly: true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
		},
	})
}

func testAccL2VPNConsistencyLiteralNamesConfig(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name = %q
  slug = %q
  type = "vxlan"
}
`, name, name)
}

func TestAccL2VPNResource_update(t *testing.T) {
	t.Parallel()
	name := acctest.RandomWithPrefix("test-l2vpn")

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_updateInitial(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", testutil.Description1),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
			{
				Config: testAccL2VPNResourceConfig_updateModified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "description", testutil.Description2),
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
		},
	})
}

func testAccL2VPNResourceConfig_updateInitial(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = %q
}
`, name, name, testutil.Description1)
}

func testAccL2VPNResourceConfig_updateModified(name string) string {
	return fmt.Sprintf(`
resource "netbox_l2vpn" "test" {
  name        = %q
  slug        = %q
  type        = "vxlan"
  description = %q
}
`, name, name, testutil.Description2)
}

func TestAccL2VPNResource_external_deletion(t *testing.T) {
	t.Parallel()
	name := acctest.RandomWithPrefix("test-l2vpn")

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "name", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "slug", name),
					resource.TestCheckResourceAttr("netbox_l2vpn.test", "type", "vxlan"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					items, _, err := client.VpnAPI.VpnL2vpnsList(context.Background()).Name([]string{name}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find l2vpn for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.VpnAPI.VpnL2vpnsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete l2vpn: %v", err)
					}
					t.Logf("Successfully externally deleted l2vpn with ID: %d", itemID)
				},
				Config: testAccL2VPNResourceConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_l2vpn.test", "id"),
				),
			},
		},
	})
}
