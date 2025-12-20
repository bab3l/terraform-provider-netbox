package resources_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/resources"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	fwresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIPRangeResource(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource()

	if r == nil {

		t.Fatal("Expected non-nil IPRange resource")

	}

}

func TestIPRangeResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"start_address", "end_address"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	computedAttrs := []string{"id", "size"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	optionalAttrs := []string{"vrf", "tenant", "status", "role", "description", "comments", "mark_utilized", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestIPRangeResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ip_range"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestIPRangeResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewIPRangeResource().(*resources.IPRangeResource)

	configureRequest := fwresource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = testutil.InvalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccIPRangeResource_basic(t *testing.T) {

	// Generate random octets for unique IP range

	second := acctest.RandIntRange(0, 255)

	third := acctest.RandIntRange(0, 255)

	startAddr := fmt.Sprintf("10.%d.%d.10/24", second, third)

	endAddr := fmt.Sprintf("10.%d.%d.20/24", second, third)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_basic(startAddr, endAddr),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddr),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddr),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "status", "active"),
				),
			},

			{

				ResourceName: "netbox_ip_range.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccIPRangeResource_full(t *testing.T) {

	// Generate random octets for unique IP range

	second := acctest.RandIntRange(0, 255)

	third := acctest.RandIntRange(0, 255)

	startAddr := fmt.Sprintf("10.%d.%d.10/24", second, third)

	endAddr := fmt.Sprintf("10.%d.%d.20/24", second, third)

	description := "Test IP range with all fields"

	updatedDescription := "Updated IP range description"

	const comments = "Test comments"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeResourceConfig_full(startAddr, endAddr, "active", description, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_ip_range.test", "id"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddr),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddr),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", description),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "comments", comments),
				),
			},

			{

				Config: testAccIPRangeResourceConfig_full(startAddr, endAddr, "reserved", updatedDescription, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_range.test", "status", "reserved"),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccIPRangeResourceConfig_basic(startAddr, endAddr string) string {

	return fmt.Sprintf(`

resource "netbox_ip_range" "test" {

  start_address = %q

  end_address   = %q

}

`, startAddr, endAddr)

}

func testAccIPRangeResourceConfig_full(startAddr, endAddr, status, description, comments string) string {

	return fmt.Sprintf(`

resource "netbox_ip_range" "test" {

  start_address = %q

  end_address   = %q

  status        = %q

  description   = %q

  comments      = %q

}

`, startAddr, endAddr, status, description, comments)

}

func TestAccConsistency_IPRange(t *testing.T) {

	t.Parallel()

	startAddress := "10.100.0.1/24"

	endAddress := "10.100.0.100/24"

	vrfName := testutil.RandomName("vrf")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	roleName := testutil.RandomName("role")

	roleSlug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeConsistencyConfig(startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "vrf", vrfName),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "tenant", tenantName),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "role", roleSlug),
				),
			},

			{

				PlanOnly: true,

				Config: testAccIPRangeConsistencyConfig(startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug),
			},
		},
	})

}

func testAccIPRangeConsistencyConfig(startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = "%[3]s"

}



resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}



resource "netbox_role" "test" {

  name = "%[6]s"

  slug = "%[7]s"

}



resource "netbox_ip_range" "test" {

  start_address = "%[1]s"

  end_address = "%[2]s"

  vrf = netbox_vrf.test.name

  tenant = netbox_tenant.test.name

  role = netbox_role.test.slug

}

`, startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug)

}

// TestAccConsistency_IPRange_LiteralNames tests that reference attributes specified as literal string names

// are preserved and do not cause drift when the API returns numeric IDs.

func TestAccConsistency_IPRange_LiteralNames(t *testing.T) {

	t.Parallel()

	startAddress := "10.10.0.1/32"

	endAddress := "10.10.0.254/32"

	vrfName := testutil.RandomName("vrf")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	roleName := testutil.RandomName("role")

	roleSlug := testutil.RandomSlug("role")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_ip_range.test", "start_address", startAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "end_address", endAddress),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "vrf", vrfName),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "tenant", tenantName),

					resource.TestCheckResourceAttr("netbox_ip_range.test", "role", roleSlug),
				),
			},

			{

				// Critical: Verify no drift when refreshing state

				PlanOnly: true,

				Config: testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug),
			},
		},
	})

}

func testAccIPRangeConsistencyLiteralNamesConfig(startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug string) string {

	return fmt.Sprintf(`

resource "netbox_vrf" "test" {

  name = "%[3]s"

}



resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}



resource "netbox_role" "test" {

  name = "%[6]s"

  slug = "%[7]s"

}



resource "netbox_ip_range" "test" {

  start_address = "%[1]s"

  end_address = "%[2]s"

  # Use literal string names to mimic existing user state

  vrf = "%[3]s"

  tenant = "%[4]s"

  role = "%[7]s"



  depends_on = [netbox_vrf.test, netbox_tenant.test, netbox_role.test]

}

`, startAddress, endAddress, vrfName, tenantName, tenantSlug, roleName, roleSlug)

}
