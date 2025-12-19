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

func TestAggregateResource(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource()

	if r == nil {

		t.Fatal("Expected non-nil Aggregate resource")

	}

}

func TestAggregateResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"prefix", "rir"}

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

	optionalAttrs := []string{"tenant", "date_added", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestAggregateResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_aggregate"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestAggregateResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewAggregateResource().(*resources.AggregateResource)

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

func TestAccAggregateResource_basic(t *testing.T) {

	rirName := testutil.RandomName("tf-test-rir")

	rirSlug := testutil.RandomSlug("tf-test-rir")

	prefix := "192.0.2.0/24"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "rir"),
				),
			},

			{

				ResourceName: "netbox_aggregate.test",

				ImportState: true,

				ImportStateVerify: true,

				ImportStateVerifyIgnore: []string{"rir"},
			},
		},
	})

}

func TestAccAggregateResource_full(t *testing.T) {

	rirName := testutil.RandomName("tf-test-rir-full")

	rirSlug := testutil.RandomSlug("tf-test-rir-full")

	prefix := "198.51.100.0/24"

	description := "Test aggregate with all fields"

	updatedDescription := "Updated aggregate description"

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, description, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_aggregate.test", "id"),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", description),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "comments", comments),
				),
			},

			{

				Config: testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, updatedDescription, comments),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "description", updatedDescription),
				),
			},
		},
	})

}

func testAccAggregateResourceConfig_basic(rirName, rirSlug, prefix string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {

  name = %q

  slug = %q

}



resource "netbox_aggregate" "test" {

  prefix = %q

  rir    = netbox_rir.test.id

}



`, rirName, rirSlug, prefix)

}

func testAccAggregateResourceConfig_full(rirName, rirSlug, prefix, description, comments string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {

  name = %q

  slug = %q

}



resource "netbox_aggregate" "test" {

  prefix      = %q

  rir         = netbox_rir.test.id

  description = %q

  comments    = %q

}



`, rirName, rirSlug, prefix, description, comments)

}

func TestAccConsistency_Aggregate(t *testing.T) {

	t.Parallel()

	prefix := "192.168.0.0/16"

	rirName := testutil.RandomName("rir")

	rirSlug := testutil.RandomSlug("rir")

	tenantName := testutil.RandomName("tenant")

	tenantSlug := testutil.RandomSlug("tenant")

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: testutil.TestAccProtoV6ProviderFactories,

		Steps: []resource.TestStep{

			{

				Config: testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_aggregate.test", "prefix", prefix),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "rir", rirSlug),

					resource.TestCheckResourceAttr("netbox_aggregate.test", "tenant", tenantName),
				),
			},

			{

				PlanOnly: true,

				Config: testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug),
			},
		},
	})

}

func testAccAggregateConsistencyConfig(prefix, rirName, rirSlug, tenantName, tenantSlug string) string {

	return fmt.Sprintf(`



resource "netbox_rir" "test" {

  name = "%[2]s"

  slug = "%[3]s"

}



resource "netbox_tenant" "test" {

  name = "%[4]s"

  slug = "%[5]s"

}



resource "netbox_aggregate" "test" {

  prefix = "%[1]s"

  rir = netbox_rir.test.slug

  tenant = netbox_tenant.test.name

}



`, prefix, rirName, rirSlug, tenantName, tenantSlug)

}
