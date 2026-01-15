// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// CheckTenantGroupDestroy verifies that a tenant group has been destroyed.

func CheckTenantGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_tenant_group" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("tenant group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckTenantDestroy verifies that a tenant has been destroyed.

func CheckTenantDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_tenant" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("tenant with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckContactGroupDestroy verifies that a contact group has been destroyed.

func CheckContactGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_contact_group" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.TenancyAPI.TenancyContactGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("contact group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckContactDestroy verifies that a contact has been destroyed.

func CheckContactDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_contact" {

			continue

		}

		email := rs.Primary.Attributes["email"]

		if email == "" {

			continue

		}

		list, resp, err := client.TenancyAPI.TenancyContactsList(ctx).Email([]string{email}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("contact with email %s still exists (ID: %d)", email, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckContactRoleDestroy verifies that a contact role has been destroyed.

func CheckContactRoleDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_contact_role" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.TenancyAPI.TenancyContactRolesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("contact role with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckVRFDestroy verifies that a VRF has been destroyed.

func CheckVRFDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_vrf" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamVrfsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("VRF with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckVLANGroupDestroy verifies that a VLAN group has been destroyed.

func CheckVLANGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_vlan_group" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamVlanGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("VLAN group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckVLANDestroy verifies that a VLAN has been destroyed.

func CheckVLANDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_vlan" {

			continue

		}

		id := rs.Primary.ID

		if id != "" {

			var idInt int32

			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {

				_, resp, err := client.IpamAPI.IpamVlansRetrieve(ctx, idInt).Execute()

				if err == nil && resp.StatusCode == http.StatusOK {

					return fmt.Errorf("VLAN with ID %s still exists", id)

				}

			}

		}

	}

	return nil

}

// CheckPrefixDestroy verifies that a prefix has been destroyed.

func CheckPrefixDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_prefix" {

			continue

		}

		id := rs.Primary.ID

		if id != "" {

			var idInt int32

			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {

				_, resp, err := client.IpamAPI.IpamPrefixesRetrieve(ctx, idInt).Execute()

				if err == nil && resp.StatusCode == http.StatusOK {

					return fmt.Errorf("prefix with ID %s still exists", id)

				}

			}

		}

	}

	return nil

}

// CheckIPAddressDestroy verifies that an IP address has been destroyed.

func CheckIPAddressDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_ip_address" {

			continue

		}

		id := rs.Primary.ID

		if id != "" {

			var idInt int32

			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {

				_, resp, err := client.IpamAPI.IpamIpAddressesRetrieve(ctx, idInt).Execute()

				if err == nil && resp.StatusCode == http.StatusOK {

					return fmt.Errorf("IP address with ID %s still exists", id)

				}

			}

		}

	}

	return nil

}

// CheckASNRangeDestroy verifies that an ASN range has been destroyed.

func CheckASNRangeDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_asn_range" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamAsnRangesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("ASN range with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckRIRDestroy verifies that an RIR has been destroyed.

func CheckRIRDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_rir" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamRirsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("RIR with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckRouteTargetDestroy verifies that a route target has been destroyed.

func CheckRouteTargetDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_route_target" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamRouteTargetsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("route target with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckFHRPGroupDestroy verifies that an FHRP group has been destroyed.

func CheckFHRPGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_fhrp_group" {

			continue

		}

		idStr := rs.Primary.ID

		if idStr == "" {

			continue

		}

		id, parseErr := strconv.Atoi(idStr)

		if parseErr != nil {

			continue

		}

		_, resp, err := client.IpamAPI.IpamFhrpGroupsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("FHRP group with ID %d still exists", id)

		}

	}

	return nil

}

// CheckRoleDestroy verifies that a role has been destroyed.

func CheckRoleDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_role" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamRolesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("role with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckServiceDestroy verifies that a service has been destroyed.

func CheckServiceDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_service" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamServicesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("service with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckServiceTemplateDestroy verifies that a service template has been destroyed.

func CheckServiceTemplateDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_service_template" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.IpamAPI.IpamServiceTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("service template with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckProviderDestroy verifies that a circuit provider has been destroyed.

func CheckProviderDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_provider" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.CircuitsAPI.CircuitsProvidersList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("provider with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckCircuitTypeDestroy verifies that a circuit type has been destroyed.

func CheckCircuitTypeDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_circuit_type" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.CircuitsAPI.CircuitsCircuitTypesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("circuit type with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckCircuitDestroy verifies that a circuit has been destroyed.

func CheckCircuitDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_circuit" {

			continue

		}

		cid := rs.Primary.Attributes["cid"]

		if cid == "" {

			continue

		}

		list, resp, err := client.CircuitsAPI.CircuitsCircuitsList(ctx).Cid([]string{cid}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("circuit with CID %s still exists (ID: %d)", cid, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckCircuitGroupDestroy verifies that a circuit group has been destroyed.

func CheckCircuitGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_circuit_group" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("circuit group with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckCircuitGroupAssignmentDestroy verifies that a circuit group assignment has been destroyed.

func CheckCircuitGroupAssignmentDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_circuit_group_assignment" {

			continue

		}

		idStr := rs.Primary.ID

		if idStr == "" {

			continue

		}

		id, parseErr := strconv.Atoi(idStr)

		if parseErr != nil {

			continue

		}

		_, resp, err := client.CircuitsAPI.CircuitsCircuitGroupAssignmentsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("circuit group assignment with ID %d still exists", id)

		}

	}

	return nil

}

// CheckIKEProposalDestroy verifies that an IKE proposal has been destroyed.

func CheckIKEProposalDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_ike_proposal" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnIkeProposalsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("IKE proposal with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckIKEPolicyDestroy verifies that an IKE policy has been destroyed.

func CheckIKEPolicyDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_ike_policy" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnIkePoliciesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("IKE policy with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckIPSecProposalDestroy verifies that an IPSec proposal has been destroyed.

func CheckIPSecProposalDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_ipsec_proposal" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnIpsecProposalsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("IPSec proposal with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckIPSecPolicyDestroy verifies that an IPSec policy has been destroyed.

func CheckIPSecPolicyDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_ipsec_policy" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnIpsecPoliciesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("IPSec policy with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckIPSecProfileDestroy verifies that an IPSec profile has been destroyed.

func CheckIPSecProfileDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_ipsec_profile" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnIpsecProfilesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("IPSec profile with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckTunnelGroupDestroy verifies that a tunnel group has been destroyed.

func CheckTunnelGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_tunnel_group" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnTunnelGroupsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("tunnel group with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckTunnelDestroy verifies that a tunnel has been destroyed.

func CheckTunnelDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_tunnel" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnTunnelsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("tunnel with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckTunnelTerminationDestroy verifies that a tunnel termination has been destroyed.

func CheckTunnelTerminationDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_tunnel_termination" {

			continue

		}

		idStr := rs.Primary.ID

		if idStr == "" {

			continue

		}

		id, parseErr := strconv.Atoi(idStr)

		if parseErr != nil {

			continue

		}

		_, resp, err := client.VpnAPI.VpnTunnelTerminationsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("tunnel termination with ID %d still exists", id)

		}

	}

	return nil

}

// CheckL2VPNDestroy verifies that an L2VPN has been destroyed.

func CheckL2VPNDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_l2vpn" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VpnAPI.VpnL2vpnsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("L2VPN with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckClusterTypeDestroy verifies that a cluster type has been destroyed.

func CheckClusterTypeDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_cluster_type" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.VirtualizationAPI.VirtualizationClusterTypesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("cluster type with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckClusterDestroy verifies that a cluster has been destroyed.

func CheckClusterDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_cluster" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VirtualizationAPI.VirtualizationClustersList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("cluster with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckVirtualMachineDestroy verifies that a virtual machine has been destroyed.

func CheckVirtualMachineDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_virtual_machine" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("virtual machine with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckVMInterfaceDestroy verifies that a VM interface has been destroyed.

func CheckVMInterfaceDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_vm_interface" {

			continue

		}

		id := rs.Primary.ID

		if id != "" {

			var idInt int32

			if _, parseErr := fmt.Sscanf(id, "%d", &idInt); parseErr == nil {

				_, resp, err := client.VirtualizationAPI.VirtualizationInterfacesRetrieve(ctx, idInt).Execute()

				if err == nil && resp.StatusCode == http.StatusOK {

					return fmt.Errorf("VM interface with ID %s still exists", id)

				}

			}

		}

	}

	return nil

}

// CheckVirtualDiskDestroy verifies that a virtual disk has been destroyed.

func CheckVirtualDiskDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_virtual_disk" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.VirtualizationAPI.VirtualizationVirtualDisksList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list.Count > 0 {

			return fmt.Errorf("virtual disk with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckClusterGroupDestroy verifies that a cluster group has been destroyed.

func CheckClusterGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_cluster_group" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.VirtualizationAPI.VirtualizationClusterGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("cluster group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckWirelessLinkDestroy verifies that a wireless link has been destroyed.

func CheckWirelessLinkDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_wireless_link" {

			continue

		}

		idStr := rs.Primary.ID

		if idStr == "" {

			continue

		}

		id, parseErr := strconv.Atoi(idStr)

		if parseErr != nil {

			continue

		}

		_, resp, err := client.WirelessAPI.WirelessWirelessLinksRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("wireless link with ID %d still exists", id)

		}

	}

	return nil

}

// CheckWirelessLANDestroy verifies that a wireless LAN has been destroyed.

func CheckWirelessLANDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_wireless_lan" {

			continue

		}

		ssid := rs.Primary.Attributes["ssid"]

		if ssid == "" {

			continue

		}

		list, resp, err := client.WirelessAPI.WirelessWirelessLansList(ctx).Ssid([]string{ssid}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("wireless LAN with SSID %s still exists (ID: %d)", ssid, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckWirelessLANGroupDestroy verifies that a wireless LAN group has been destroyed.

func CheckWirelessLANGroupDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_wireless_lan_group" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.WirelessAPI.WirelessWirelessLanGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("wireless LAN group with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckJournalEntryDestroy verifies that a journal entry has been destroyed.

func CheckJournalEntryDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_journal_entry" {

			continue

		}

		idStr := rs.Primary.ID

		if idStr == "" {

			continue

		}

		id, parseErr := strconv.Atoi(idStr)

		if parseErr != nil {

			continue

		}

		_, resp, err := client.ExtrasAPI.ExtrasJournalEntriesRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("journal entry with ID %d still exists", id)

		}

	}

	return nil

}

// CheckCustomFieldChoiceSetDestroy verifies that a custom field choice set has been destroyed.

func CheckCustomFieldChoiceSetDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_custom_field_choice_set" {

			continue

		}

		idStr := rs.Primary.ID

		if idStr == "" {

			continue

		}

		id, parseErr := strconv.Atoi(idStr)

		if parseErr != nil {

			continue

		}

		_, resp, err := client.ExtrasAPI.ExtrasCustomFieldChoiceSetsRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("custom field choice set with ID %d still exists", id)

		}

	}

	return nil

}

// CheckCustomLinkDestroy verifies that a custom link has been destroyed.

func CheckCustomLinkDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_custom_link" {

			continue

		}

		idStr := rs.Primary.ID

		if idStr == "" {

			continue

		}

		id, parseErr := strconv.Atoi(idStr)

		if parseErr != nil {

			continue

		}

		_, resp, err := client.ExtrasAPI.ExtrasCustomLinksRetrieve(ctx, int32(id)).Execute() // #nosec G109,G115 -- test utility, ID from Terraform state is within int32 range

		if err == nil && resp.StatusCode == http.StatusOK {

			return fmt.Errorf("custom link with ID %d still exists", id)

		}

	}

	return nil

}

// CheckTagDestroy verifies that a tag has been destroyed.

func CheckTagDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_tag" {

			continue

		}

		slug := rs.Primary.Attributes["slug"]

		if slug == "" {

			continue

		}

		list, resp, err := client.ExtrasAPI.ExtrasTagsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("tag with slug %s still exists (ID: %d)", slug, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckWebhookDestroy verifies that a webhook has been destroyed.

func CheckWebhookDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_webhook" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.ExtrasAPI.ExtrasWebhooksList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("webhook with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckEventRuleDestroy verifies that an event rule has been destroyed.
func CheckEventRuleDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_event_rule" {
			continue
		}

		// Try to get the ID from state
		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			// If ID is not a number, skip this resource
			continue
		}

		// Try to retrieve the event rule by ID
		//nolint:gosec // ID from Terraform state is always a valid positive integer
		_, resp, err := client.ExtrasAPI.ExtrasEventRulesRetrieve(ctx, int32(id)).Execute()
		if err != nil {
			// If 404, the resource is destroyed (expected)
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				continue
			}
			// For other errors, continue checking other resources
			continue
		}

		// If we successfully retrieved the resource, it still exists
		if resp.StatusCode == http.StatusOK {
			return fmt.Errorf("event rule with ID %d still exists", id)
		}
	}

	return nil
}

// CheckNotificationGroupDestroy verifies that a notification group has been destroyed.
func CheckNotificationGroupDestroy(s *terraform.State) error {
	client, err := GetSharedClient()
	if err != nil {
		return fmt.Errorf("failed to get client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "netbox_notification_group" {
			continue
		}

		// Try to get the ID from state
		idStr := rs.Primary.ID
		if idStr == "" {
			continue
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			// If ID is not a number, skip this resource
			continue
		}

		// Try to retrieve the notification group by ID
		//nolint:gosec // ID from Terraform state is always a valid positive integer
		_, resp, err := client.ExtrasAPI.ExtrasNotificationGroupsRetrieve(ctx, int32(id)).Execute()
		if err != nil {
			// If 404, the resource is destroyed (expected)
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				continue
			}
			// For other errors, continue checking other resources
			continue
		}

		// If we successfully retrieved the resource, it still exists
		if resp.StatusCode == http.StatusOK {
			return fmt.Errorf("notification group with ID %d still exists", id)
		}
	}

	return nil
}

// CheckExportTemplateDestroy verifies that an export template has been destroyed.

func CheckExportTemplateDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_export_template" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.ExtrasAPI.ExtrasExportTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("export template with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckConfigTemplateDestroy verifies that a config template has been destroyed.

func CheckConfigTemplateDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_config_template" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.ExtrasAPI.ExtrasConfigTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("config template with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}

// CheckConfigContextDestroy verifies that a config context has been destroyed.

func CheckConfigContextDestroy(s *terraform.State) error {

	client, err := GetSharedClient()

	if err != nil {

		return fmt.Errorf("failed to get client: %w", err)

	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "netbox_config_context" {

			continue

		}

		name := rs.Primary.Attributes["name"]

		if name == "" {

			continue

		}

		list, resp, err := client.ExtrasAPI.ExtrasConfigContextsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			continue

		}

		if resp.StatusCode == http.StatusOK && list != nil && len(list.Results) > 0 {

			return fmt.Errorf("config context with name %s still exists (ID: %d)", name, list.Results[0].GetId())

		}

	}

	return nil

}
