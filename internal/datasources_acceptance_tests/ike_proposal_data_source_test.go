package datasources_acceptance_tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/bab3l/go-netbox"
	"github.com/bab3l/terraform-provider-netbox/internal/datasources"
	"github.com/bab3l/terraform-provider-netbox/internal/provider"
	"github.com/bab3l/terraform-provider-netbox/internal/testutil"
	fwdatasource "github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestIKEProposalDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewIKEProposalDataSource()

	if d == nil {

		t.Fatal("Expected non-nil IKE proposal data source")

	}

}

func TestIKEProposalDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewIKEProposalDataSource()

	schemaRequest := fwdatasource.SchemaRequest{}

	schemaResponse := &fwdatasource.SchemaResponse{}

	d.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	// Check that key attributes exist

	requiredAttrs := []string{"id", "name", "authentication_method", "encryption_algorithm", "authentication_algorithm", "group", "description", "comments", "tags", "sa_lifetime"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)

		}

	}

}

func TestIKEProposalDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewIKEProposalDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ike_proposal"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestIKEProposalDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewIKEProposalDataSource().(*datasources.IKEProposalDataSource)

	configureRequest := fwdatasource.ConfigureRequest{

		ProviderData: nil,
	}

	configureResponse := &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with nil provider data, got: %+v", configureResponse.Diagnostics)

	}

	client := &netbox.APIClient{}

	configureRequest.ProviderData = client

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if configureResponse.Diagnostics.HasError() {

		t.Errorf("Expected no error with correct provider data, got: %+v", configureResponse.Diagnostics)

	}

	configureRequest.ProviderData = testutil.InvalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with invalid provider data")

	}

}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccIKEProposalDataSource_byID(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("tf-test-ike-proposal-ds")

	cleanup.RegisterIKEProposalCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIKEProposalDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ike_proposal.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "authentication_method", "preshared-keys"),

					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "encryption_algorithm", "aes-256-cbc"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEProposalDestroy,
		),
	})

}

func TestAccIKEProposalDataSource_byName(t *testing.T) {

	t.Parallel()

	testutil.TestAccPreCheck(t)

	cleanup := testutil.NewCleanupResource(t)

	randomName := testutil.RandomName("tf-test-ike-proposal-ds")

	cleanup.RegisterIKEProposalCleanup(randomName)

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIKEProposalDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ike_proposal.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ike_proposal.test", "authentication_method", "preshared-keys"),
				),
			},
		},
		CheckDestroy: testutil.ComposeCheckDestroy(
			testutil.CheckIKEProposalDestroy,
		),
	})

}

func testAccIKEProposalDataSourceByID(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_proposal" "test" {

  name                     = %[1]q

  authentication_method    = "preshared-keys"

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"

  group                    = 14

}

data "netbox_ike_proposal" "test" {

  id = netbox_ike_proposal.test.id

}

`, name)

}

func testAccIKEProposalDataSourceByName(name string) string {

	return fmt.Sprintf(`

resource "netbox_ike_proposal" "test" {

  name                     = %[1]q

  authentication_method    = "preshared-keys"

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"

  group                    = 14

}

data "netbox_ike_proposal" "test" {

  name = netbox_ike_proposal.test.name

}

`, name)

}
