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

func TestAccCircuitTerminationResource_basic(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	circuitTypeName := testutil.RandomName("tf-test-ct")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct")
	circuitCID := testutil.RandomName("tf-test-circuit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				ResourceName:            "netbox_circuit_termination.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"circuit", "site"},
			},
		},
	})
}

func TestAccCircuitTerminationResource_full(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-full")
	providerSlug := testutil.RandomSlug("tf-test-provider-full")
	circuitTypeName := testutil.RandomName("tf-test-ct-full")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-full")
	circuitCID := testutil.RandomName("tf-test-circuit-full")
	siteName := testutil.RandomName("tf-test-site-full")
	siteSlug := testutil.RandomSlug("tf-test-site-full")
	description := testutil.RandomName("description")
	updatedDescription := "Updated circuit termination description"
	tagName := testutil.RandomName("tf-test-tag")
	tagSlug := testutil.RandomSlug("tf-test-tag")
	customFieldName := testutil.RandomCustomFieldName("tf_test_cf")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description, tagName, tagSlug, customFieldName, 1000000, 512000, "XCON-123", "PP1-Port5", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", description),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "1000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "upstream_speed", "512000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "xconnect_id", "XCON-123"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "pp_info", "PP1-Port5"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "mark_connected", "true"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "tags.#", "1"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "custom_fields.#", "1"),
				),
			},
			{
				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, updatedDescription, tagName, tagSlug, customFieldName, 10000000, 5000000, "XCON-456", "PP2-Port10", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", updatedDescription),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "10000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "upstream_speed", "5000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "xconnect_id", "XCON-456"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "pp_info", "PP2-Port10"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "mark_connected", "false"),
				),
			},
		},
	})
}

func TestAccCircuitTerminationResource_update(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider-upd")
	providerSlug := testutil.RandomSlug("tf-test-provider-upd")
	circuitTypeName := testutil.RandomName("tf-test-ct-upd")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-upd")
	circuitCID := testutil.RandomName("tf-test-circuit-upd")
	siteName := testutil.RandomName("tf-test-site-upd")
	siteSlug := testutil.RandomSlug("tf-test-site-upd")
	tagName := testutil.RandomName("tf-test-tag-upd")
	tagSlug := testutil.RandomSlug("tf-test-tag-upd")
	customFieldName := testutil.RandomCustomFieldName("tf_test_cf_upd")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTagCleanup(tagSlug)
	cleanup.RegisterCustomFieldCleanup(customFieldName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				Config: testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, testutil.Description2, tagName, tagSlug, customFieldName, 10000000, 5000000, "XCON-UPDATE", "PP-UPDATE", true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "description", testutil.Description2),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "port_speed", "10000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "upstream_speed", "5000000"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "xconnect_id", "XCON-UPDATE"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "pp_info", "PP-UPDATE"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "mark_connected", "true"),
				),
			},
		},
	})
}

func TestAccCircuitTerminationResource_IDPreservation(t *testing.T) {

	t.Parallel()
	providerName := testutil.RandomName("ct-prov-id")
	providerSlug := testutil.RandomSlug("ct-prov-id")
	circuitTypeName := testutil.RandomName("ct-type-id")
	circuitTypeSlug := testutil.RandomSlug("ct-type-id")
	circuitCID := testutil.RandomName("ct-ckt-id")
	siteName := testutil.RandomName("ct-site-id")
	siteSlug := testutil.RandomSlug("ct-site-id")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
		},
	})
}

func testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.id
  term_side = "A"
  site      = netbox_site.test.id
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug)
}

func testAccCircuitTerminationResourceConfig_full(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, description, tagName, tagSlug, customFieldName string, portSpeed, upstreamSpeed int, xconnectID, ppInfo string, markConnected bool) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

resource "netbox_tag" "test" {
  name = %q
  slug = %q
}

resource "netbox_custom_field" "test" {
  name         = %q
  object_types = ["circuits.circuittermination"]
  type         = "text"
}

resource "netbox_circuit_termination" "test" {
  circuit        = netbox_circuit.test.id
  term_side      = "A"
  site           = netbox_site.test.id
  port_speed     = %d
  upstream_speed = %d
  xconnect_id    = %q
  pp_info        = %q
  mark_connected = %t
  description    = %q
  tags = [
    {
      name = netbox_tag.test.name
      slug = netbox_tag.test.slug
    }
  ]
  custom_fields = [
    {
      name  = netbox_custom_field.test.name
      type  = "text"
      value = "test-value"
    }
  ]
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tagName, tagSlug, customFieldName, portSpeed, upstreamSpeed, xconnectID, ppInfo, markConnected, description)
}

func TestAccCircuitTerminationResource_import(t *testing.T) {
	t.Parallel()

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	circuitTypeName := testutil.RandomName("tf-test-ct")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct")
	circuitCID := testutil.RandomName("tf-test-circuit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				ResourceName:            "netbox_circuit_termination.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"circuit", "site"},
			},
			{
				Config:   testAccCircuitTerminationResourceConfig_basic(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				PlanOnly: true,
			},
		},
	})
}

func TestAccCircuitTerminationResource_importWithCustomFieldsAndTags(t *testing.T) {
	// NOTE: t.Parallel() intentionally omitted - this test creates/deletes global custom fields
	// that would affect other tests of the same resource type running in parallel.

	providerName := testutil.RandomName("tf-test-provider")
	providerSlug := testutil.RandomSlug("tf-test-provider")
	circuitTypeName := testutil.RandomName("tf-test-ct")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct")
	circuitCID := testutil.RandomName("tf-test-circuit")
	siteName := testutil.RandomName("tf-test-site")
	siteSlug := testutil.RandomSlug("tf-test-site")
	tenantName := testutil.RandomName("tf-test-tenant")
	tenantSlug := testutil.RandomSlug("tf-test-tenant")

	// Generate all random values for custom fields and tags
	textValue := testutil.RandomName("text-value")
	longtextValue := testutil.RandomName("longtext-value") + "\nThis is a multiline text field for comprehensive testing."
	intValue := 42
	boolValue := true
	dateValue := testutil.RandomDate()
	urlValue := testutil.RandomURL("test-url")
	jsonValue := testutil.RandomJSON()

	tag1Name := testutil.RandomName("tag1")
	tag1Slug := testutil.RandomSlug("tag1")
	tag2Name := testutil.RandomName("tag2")
	tag2Slug := testutil.RandomSlug("tag2")

	cfText := testutil.RandomCustomFieldName("tf_text")
	cfLongtext := testutil.RandomCustomFieldName("tf_longtext")
	cfInteger := testutil.RandomCustomFieldName("tf_integer")
	cfBoolean := testutil.RandomCustomFieldName("tf_boolean")
	cfDate := testutil.RandomCustomFieldName("tf_date")
	cfURL := testutil.RandomCustomFieldName("tf_url")
	cfJSON := testutil.RandomCustomFieldName("tf_json")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterProviderCleanup(providerSlug)
	cleanup.RegisterCircuitTypeCleanup(circuitTypeSlug)
	cleanup.RegisterCircuitCleanup(circuitCID)
	cleanup.RegisterSiteCleanup(siteSlug)
	cleanup.RegisterTenantCleanup(tenantSlug)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationResourceImportConfig_full(
					providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tenantName, tenantSlug,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue,
					tag1Name, tag1Slug, tag2Name, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "term_side", "A"),
				),
			},
			{
				ResourceName:            "netbox_circuit_termination.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"circuit", "site", "custom_fields", "tags"},
			},
			{
				Config: testAccCircuitTerminationResourceImportConfig_full(
					providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tenantName, tenantSlug,
					textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue,
					tag1Name, tag1Slug, tag2Name, tag2Slug,
					cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
				),
				PlanOnly: true,
			},
		},
	})
}

func testAccCircuitTerminationResourceImportConfig_full(
	providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tenantName, tenantSlug string,
	textValue, longtextValue string, intValue int, boolValue bool, dateValue, urlValue, jsonValue string,
	tag1Name, tag1Slug, tag2Name, tag2Slug string,
	cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON string,
) string {
	return fmt.Sprintf(`
# Dependencies
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
  status           = "active"
}

resource "netbox_site" "test" {
  name   = %q
  slug   = %q
  status = "active"
}

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

# Custom Fields for circuits.circuittermination object type
resource "netbox_custom_field" "test_text" {
  name         = %q
  label        = "Test Text CF"
  type         = "text"
  object_types = ["circuits.circuittermination"]
}

resource "netbox_custom_field" "test_longtext" {
  name         = %q
  label        = "Test Longtext CF"
  type         = "longtext"
  object_types = ["circuits.circuittermination"]
}

resource "netbox_custom_field" "test_integer" {
  name         = %q
  label        = "Test Integer CF"
  type         = "integer"
  object_types = ["circuits.circuittermination"]
}

resource "netbox_custom_field" "test_boolean" {
  name         = %q
  label        = "Test Boolean CF"
  type         = "boolean"
  object_types = ["circuits.circuittermination"]
}

resource "netbox_custom_field" "test_date" {
  name         = %q
  label        = "Test Date CF"
  type         = "date"
  object_types = ["circuits.circuittermination"]
}

resource "netbox_custom_field" "test_url" {
  name         = %q
  label        = "Test URL CF"
  type         = "url"
  object_types = ["circuits.circuittermination"]
}

resource "netbox_custom_field" "test_json" {
  name         = %q
  label        = "Test JSON CF"
  type         = "json"
  object_types = ["circuits.circuittermination"]
}

# Circuit Termination with comprehensive custom fields and tags
resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.id
  site      = netbox_site.test.id
  term_side = "A"

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
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug, tenantName, tenantSlug,
		tag1Name, tag1Slug, tag2Name, tag2Slug,
		cfText, cfLongtext, cfInteger, cfBoolean, cfDate, cfURL, cfJSON,
		textValue, longtextValue, intValue, boolValue, dateValue, urlValue, jsonValue)
}

// TestAccConsistency_CircuitTermination_LiteralNames tests that reference attributes specified as literal string names
// are preserved and do not cause drift when the API returns numeric IDs.
func TestAccConsistency_CircuitTermination_LiteralNames(t *testing.T) {

	t.Parallel()
	providerName := testutil.RandomName("provider")
	providerSlug := testutil.RandomSlug("provider")
	circuitTypeName := testutil.RandomName("circuit-type")
	circuitTypeSlug := testutil.RandomSlug("circuit-type")
	circuitCid := testutil.RandomName("CID")
	siteName := testutil.RandomName("site")
	siteSlug := testutil.RandomSlug("site")

	cleanup := testutil.NewCleanupResource(t)
	cleanup.RegisterCircuitCleanup(circuitCid)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCircuitTerminationConsistencyLiteralNamesConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "circuit", circuitCid),
					resource.TestCheckResourceAttr("netbox_circuit_termination.test", "site", siteName),
				),
			},
			{
				// Critical: Verify no drift when refreshing state
				PlanOnly: true,
				Config:   testAccCircuitTerminationConsistencyLiteralNamesConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug),
			},
		},
	})
}

func testAccCircuitTerminationConsistencyLiteralNamesConfig(providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug string) string {
	return fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = "%[1]s"
  slug = "%[2]s"
}

resource "netbox_circuit_type" "test" {
  name = "%[3]s"
  slug = "%[4]s"
}

resource "netbox_circuit" "test" {
  cid = "%[5]s"
  circuit_provider = netbox_provider.test.id
  type = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name = "%[6]s"
  slug = "%[7]s"
}

resource "netbox_circuit_termination" "test" {
  # Use literal string names to mimic existing user state
  circuit = "%[5]s"
  term_side = "A"
  site = "%[6]s"
  depends_on = [netbox_circuit.test, netbox_site.test]
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCid, siteName, siteSlug)
}

func TestAccCircuitTerminationResource_externalDeletion(t *testing.T) {
	t.Parallel()
	providerName := testutil.RandomName("tf-test-provider-ext-del")
	providerSlug := testutil.RandomSlug("tf-test-provider-ext-del")
	circuitTypeName := testutil.RandomName("tf-test-ct-ext-del")
	circuitTypeSlug := testutil.RandomSlug("tf-test-ct-ext-del")
	circuitCID := testutil.RandomName("tf-test-circuit-ext-del")
	siteName := testutil.RandomName("tf-test-site-ext-del")
	siteSlug := testutil.RandomSlug("tf-test-site-ext-del")
	resource.Test(t, resource.TestCase{
		PreCheck: func() { testutil.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){
			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "netbox_provider" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_type" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit" "test" {
  cid              = %q
  circuit_provider = netbox_provider.test.id
  type             = netbox_circuit_type.test.id
}

resource "netbox_site" "test" {
  name = %q
  slug = %q
}

resource "netbox_circuit_termination" "test" {
  circuit   = netbox_circuit.test.cid
  term_side = "A"
  site      = netbox_site.test.slug
}
`, providerName, providerSlug, circuitTypeName, circuitTypeSlug, circuitCID, siteName, siteSlug),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("netbox_circuit_termination.test", "id"),
				),
			},
			{
				PreConfig: func() {
					client, err := testutil.GetSharedClient()
					if err != nil {
						t.Fatalf("Failed to get shared client: %v", err)
					}
					// List terminations filtered by site slug
					items, _, err := client.CircuitsAPI.CircuitsCircuitTerminationsList(context.Background()).Site([]string{siteSlug}).Execute()
					if err != nil || items == nil || len(items.Results) == 0 {
						t.Fatalf("Failed to find circuit termination for external deletion: %v", err)
					}
					itemID := items.Results[0].Id
					_, err = client.CircuitsAPI.CircuitsCircuitTerminationsDestroy(context.Background(), itemID).Execute()
					if err != nil {
						t.Fatalf("Failed to externally delete circuit termination: %v", err)
					}
					t.Logf("Successfully externally deleted circuit termination with ID: %d", itemID)
				},
				RefreshState:       true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}
