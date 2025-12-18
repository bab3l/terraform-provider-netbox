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
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestPrefixResource(t *testing.T) {

	t.Parallel()

	r := resources.NewPrefixResource()

	if r == nil {

		t.Fatal("Expected non-nil Prefix resource")

	}

}

func TestPrefixResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewPrefixResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"prefix"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected required attribute %s to exist in schema", attr)

		}

	}

	computedAttrs := []string{"id"}

	for _, attr := range computedAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected computed attribute %s to exist in schema", attr)

		}

	}

	optionalAttrs := []string{"status", "site", "vrf", "tenant", "vlan", "role", "is_pool", "mark_utilized", "description", "comments"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestPrefixResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewPrefixResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_prefix"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestPrefixResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewPrefixResource().(*resources.PrefixResource)

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

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwresource.ConfigureResponse{}

	r.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with incorrect provider data")

	}

}

func TestAccPrefixResource_basic(t *testing.T) {

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckPrefixDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccPrefixResourceConfig_basic(prefix),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},
		},
	})

}

func TestAccPrefixResource_full(t *testing.T) {

	prefix := testutil.RandomIPv4Prefix()

	description := "Test prefix with all fields"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckPrefixDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccPrefixResourceConfig_full(prefix, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_prefix.test", "description", description),

					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "is_pool", "false"),
				),
			},
		},
	})

}

func TestAccPrefixResource_withVRF(t *testing.T) {

	prefix := testutil.RandomIPv4Prefix()

	vrfName := testutil.RandomName("tf-test-vrf")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterPrefixCleanup(prefix)

	cleanup.RegisterVRFCleanup(vrfName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.ComposeCheckDestroy(

			testutil.CheckPrefixDestroy,

			testutil.CheckVRFDestroy,
		),

		Steps: []resource.TestStep{

			{

				Config: testAccPrefixResourceConfig_withVRF(prefix, vrfName),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "vrf"),
				),
			},
		},
	})

}

func TestAccPrefixResource_ipv6(t *testing.T) {

	prefix := testutil.RandomIPv6Prefix()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckPrefixDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccPrefixResourceConfig_basic(prefix),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},
		},
	})

}

func TestAccPrefixResource_update(t *testing.T) {

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckPrefixDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccPrefixResourceConfig_basic(prefix),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},

			{

				Config: testAccPrefixResourceConfig_full(prefix, "Updated description"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_prefix.test", "description", "Updated description"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "status", "active"),
				),
			},
		},
	})

}

func testAccPrefixResourceConfig_basic(prefix string) string {

	return fmt.Sprintf(`































































resource "netbox_prefix" "test" {































































  prefix = %q































































}































































`, prefix)

}

func testAccPrefixResourceConfig_full(prefix, description string) string {

	return fmt.Sprintf(`































































resource "netbox_prefix" "test" {































































  prefix      = %q































































  description = %q































































  status      = "active"































































  is_pool     = false































































}































































`, prefix, description)

}

func testAccPrefixResourceConfig_withVRF(prefix, vrfName string) string {

	return fmt.Sprintf(`































































resource "netbox_vrf" "test" {































































  name = %q































































}































































































































resource "netbox_prefix" "test" {































































  prefix = %q































































  vrf    = netbox_vrf.test.name































































}































































`, vrfName, prefix)

}

func TestAccPrefixResource_import(t *testing.T) {

	prefix := testutil.RandomIPv4Prefix()

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterPrefixCleanup(prefix)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckPrefixDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccPrefixResourceConfig_basic(prefix),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_prefix.test", "id"),

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),
				),
			},

			{

				ResourceName: "netbox_prefix.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccConsistency_Prefix(t *testing.T) {

	t.Parallel()

	prefix := "10.0.0.0/24"

	siteName := testutil.RandomName("site")

	siteSlug := testutil.RandomSlug("site")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	vlanName := testutil.RandomName("vlan")

	vlanVid := 100

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccPrefixConsistencyConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_prefix.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_prefix.test", "site", siteName),

					resource.TestCheckResourceAttr("netbox_prefix.test", "tenant", tenantName),

					resource.TestCheckResourceAttr("netbox_prefix.test", "vlan", vlanName),
				),
			},

			{

				// Verify no drift

				PlanOnly: true,

				Config: testAccPrefixConsistencyConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid),
			},
		},
	})

}

func testAccPrefixConsistencyConfig(prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName string, vlanVid int) string {

	return fmt.Sprintf(`

resource "netbox_site" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}



resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}



resource "netbox_vlan" "test" {

  name = "%[6]s"

  vid  = %[7]d

  site = netbox_site.test.id

}



resource "netbox_prefix" "test" {

  prefix = "%[1]s"

  site = netbox_site.test.name

  tenant = netbox_tenant.test.name

  vlan = netbox_vlan.test.name

}

`, prefix, siteName, siteSlug, tenantName, tenantSlug, vlanName, vlanVid)

}
