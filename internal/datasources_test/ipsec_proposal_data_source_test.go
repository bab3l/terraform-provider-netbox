package datasources_test

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

func TestIPSecProposalDataSource(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProposalDataSource()

	if d == nil {

		t.Fatal("Expected non-nil IPSec proposal data source")
	}
}

func TestIPSecProposalDataSourceSchema(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProposalDataSource()

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

	requiredAttrs := []string{"id", "name", "encryption_algorithm", "authentication_algorithm", "description", "comments", "tags", "sa_lifetime_seconds", "sa_lifetime_data"}

	for _, attr := range requiredAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected attribute %s to exist in schema", attr)
		}
	}
}

func TestIPSecProposalDataSourceMetadata(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProposalDataSource()

	metadataRequest := fwdatasource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwdatasource.MetadataResponse{}

	d.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_ipsec_proposal"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)
	}
}

func TestIPSecProposalDataSourceConfigure(t *testing.T) {

	t.Parallel()

	d := datasources.NewIPSecProposalDataSource().(*datasources.IPSecProposalDataSource)

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

	configureRequest.ProviderData = invalidProviderData

	configureResponse = &fwdatasource.ConfigureResponse{}

	d.Configure(context.Background(), configureRequest, configureResponse)

	if !configureResponse.Diagnostics.HasError() {

		t.Error("Expected error with invalid provider data")
	}
}

// Acceptance tests require NETBOX_URL and NETBOX_API_TOKEN environment variables.

func TestAccIPSecProposalDataSource_byID(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-proposal-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecProposalDataSourceByID(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "encryption_algorithm", "aes-256-cbc"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "authentication_algorithm", "hmac-sha256"),
				),
			},
		},
	})
}

func TestAccIPSecProposalDataSource_byName(t *testing.T) {

	testutil.TestAccPreCheck(t)

	randomName := testutil.RandomName("tf-test-ipsec-proposal-ds")

	resource.Test(t, resource.TestCase{

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		Steps: []resource.TestStep{

			{

				Config: testAccIPSecProposalDataSourceByName(randomName),

				Check: resource.ComposeAggregateTestCheckFunc(

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "name", randomName),

					resource.TestCheckResourceAttrSet("data.netbox_ipsec_proposal.test", "id"),

					resource.TestCheckResourceAttr("data.netbox_ipsec_proposal.test", "encryption_algorithm", "aes-256-cbc"),
				),
			},
		},
	})
}

func testAccIPSecProposalDataSourceByID(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_proposal" "test" {
  name                     = %[1]q

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"
}

data "netbox_ipsec_proposal" "test" {

  id = netbox_ipsec_proposal.test.id
}

`, name)
}

func testAccIPSecProposalDataSourceByName(name string) string {

	return fmt.Sprintf(`

resource "netbox_ipsec_proposal" "test" {
  name                     = %[1]q

  encryption_algorithm     = "aes-256-cbc"

  authentication_algorithm = "hmac-sha256"
}

data "netbox_ipsec_proposal" "test" {
  name = netbox_ipsec_proposal.test.name
}

`, name)
}
