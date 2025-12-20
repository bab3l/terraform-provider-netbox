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

func TestRackReservationResource(t *testing.T) {

	t.Parallel()

	r := resources.NewRackReservationResource()

	if r == nil {

		t.Fatal("Expected non-nil RackReservation resource")

	}

}

func TestRackReservationResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewRackReservationResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"rack", "units", "user"}

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

	optionalAttrs := []string{"tenant", "description", "comments", "tags", "custom_fields"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestRackReservationResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewRackReservationResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_rack_reservation"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestRackReservationResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewRackReservationResource().(*resources.RackReservationResource)

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

func TestAccRackReservationResource_basic(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	rackName := testutil.RandomName("tf-test-rack")

	description := "Test rack reservation"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackReservationDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description),

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "units.#", "2"),
				),
			},

			// ImportState test

			{

				ResourceName: "netbox_rack_reservation.test",

				ImportState: true,

				ImportStateVerify: true,
			},
		},
	})

}

func TestAccRackReservationResource_update(t *testing.T) {

	siteName := testutil.RandomName("tf-test-site")

	siteSlug := testutil.RandomSlug("tf-test-site")

	rackName := testutil.RandomName("tf-test-rack")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterSiteCleanup(siteSlug)

	cleanup.RegisterRackCleanup(rackName)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckRackReservationDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description1),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_rack_reservation.test", "id"),

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description1),
				),
			},

			{

				Config: testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description2),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttr("netbox_rack_reservation.test", "description", description2),
				),
			},
		},
	})

}

func testAccRackReservationResourceConfig_basic(siteName, siteSlug, rackName, description string) string {

	return fmt.Sprintf(`

provider "netbox" {}



resource "netbox_site" "test" {

  name   = %[1]q

  slug   = %[2]q

  status = "active"

}



resource "netbox_rack" "test" {

  name     = %[3]q

  site     = netbox_site.test.id

  status   = "active"

  u_height = 42

}



data "netbox_user" "admin" {

  username = "admin"

}



resource "netbox_rack_reservation" "test" {

  rack        = netbox_rack.test.id

  units       = [1, 2]

  user        = data.netbox_user.admin.id

  description = %[4]q

}

`, siteName, siteSlug, rackName, description)

}
