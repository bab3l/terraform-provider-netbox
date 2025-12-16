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

func TestVLANGroupResource(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANGroupResource()

	if r == nil {

		t.Fatal("Expected non-nil VLAN Group resource")

	}

}

func TestVLANGroupResourceSchema(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANGroupResource()

	schemaRequest := fwresource.SchemaRequest{}

	schemaResponse := &fwresource.SchemaResponse{}

	r.Schema(context.Background(), schemaRequest, schemaResponse)

	if schemaResponse.Diagnostics.HasError() {

		t.Fatalf("Schema method diagnostics: %+v", schemaResponse.Diagnostics)

	}

	if schemaResponse.Schema.Attributes == nil {

		t.Fatal("Expected schema to have attributes")

	}

	requiredAttrs := []string{"name", "slug"}

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

	optionalAttrs := []string{"scope_type", "scope_id", "description"}

	for _, attr := range optionalAttrs {

		if _, exists := schemaResponse.Schema.Attributes[attr]; !exists {

			t.Errorf("Expected optional attribute %s to exist in schema", attr)

		}

	}

}

func TestVLANGroupResourceMetadata(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANGroupResource()

	metadataRequest := fwresource.MetadataRequest{

		ProviderTypeName: "netbox",
	}

	metadataResponse := &fwresource.MetadataResponse{}

	r.Metadata(context.Background(), metadataRequest, metadataResponse)

	expected := "netbox_vlan_group"

	if metadataResponse.TypeName != expected {

		t.Errorf("Expected type name %s, got %s", expected, metadataResponse.TypeName)

	}

}

func TestVLANGroupResourceConfigure(t *testing.T) {

	t.Parallel()

	r := resources.NewVLANGroupResource().(*resources.VLANGroupResource)

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

func TestAccVLANGroupResource_basic(t *testing.T) {

	name := testutil.RandomName("tf-test-vlangrp")

	slug := testutil.GenerateSlug("tf-test-vlangrp")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVLANGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
				),
			},
		},
	})

}

func TestAccVLANGroupResource_full(t *testing.T) {

	name := testutil.RandomName("tf-test-vlangrp-full")

	slug := testutil.GenerateSlug("tf-test-vlangrp-full")

	description := "Test VLAN Group with all fields"

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVLANGroupResourceConfig_full(name, slug, description),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "description", description),
				),
			},
		},
	})

}

func TestAccVLANGroupResource_update(t *testing.T) {

	name := testutil.RandomName("tf-test-vlangrp-upd")

	slug := testutil.GenerateSlug("tf-test-vlangrp-upd")

	updatedName := testutil.RandomName("tf-test-vlangrp-updated")

	cleanup := testutil.NewCleanupResource(t)

	cleanup.RegisterVLANGroupCleanup(slug)

	resource.Test(t, resource.TestCase{

		PreCheck: func() { testutil.TestAccPreCheck(t) },

		ProtoV6ProviderFactories: map[string]func() (tfprotov6.ProviderServer, error){

			"netbox": providerserver.NewProtocol6WithError(provider.New("test")()),
		},

		CheckDestroy: testutil.CheckVLANGroupDestroy,

		Steps: []resource.TestStep{

			{

				Config: testAccVLANGroupResourceConfig_basic(name, slug),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", name),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),
				),
			},

			{

				Config: testAccVLANGroupResourceConfig_full(updatedName, slug, "Updated description"),

				Check: resource.ComposeTestCheckFunc(

					resource.TestCheckResourceAttrSet("netbox_vlan_group.test", "id"),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "name", updatedName),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "slug", slug),

					resource.TestCheckResourceAttr("netbox_vlan_group.test", "description", "Updated description"),
				),
			},
		},
	})

}

func testAccVLANGroupResourceConfig_basic(name, slug string) string {

	return fmt.Sprintf(`































resource "netbox_vlan_group" "test" {































  name = %q































  slug = %q































}































`, name, slug)

}

func testAccVLANGroupResourceConfig_full(name, slug, description string) string {

	return fmt.Sprintf(`































resource "netbox_vlan_group" "test" {































  name        = %q































  slug        = %q































  description = %q































}































`, name, slug, description)

}
