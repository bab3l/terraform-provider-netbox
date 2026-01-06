// Package testutil provides utilities for acceptance testing of the Netbox provider.

package testutil

import (
	"context"
	"net/http"
	"time"
)

// ============================================================================

// TENANCY API CLEANUPS

// ============================================================================

// RegisterTenantGroupCleanup registers a cleanup function that will delete

// a tenant group by slug after the test completes.

func (c *CleanupResource) RegisterTenantGroupCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyTenantGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tenant groups with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: tenant group with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.TenancyAPI.TenancyTenantGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete tenant group %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted tenant group %d (slug: %s)", id, slug)

		}

	})

}

// RegisterTenantCleanup registers a cleanup function that will delete

// a tenant by slug after the test completes.

func (c *CleanupResource) RegisterTenantCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyTenantsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tenants with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: tenant with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.TenancyAPI.TenancyTenantsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete tenant %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted tenant %d (slug: %s)", id, slug)

		}

	})

}

// RegisterContactGroupCleanup registers a cleanup function that will delete

// a contact group by slug after the test completes.

func (c *CleanupResource) RegisterContactGroupCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyContactGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list contact groups with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: contact group with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.TenancyAPI.TenancyContactGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete contact group %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted contact group %d (slug: %s)", id, slug)

		}

	})

}

// RegisterContactCleanup registers a cleanup function that will delete

// a contact by email after the test completes.

func (c *CleanupResource) RegisterContactCleanup(email string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyContactsList(ctx).Email([]string{email}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list contacts with email %s: %v", email, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: contact with email %s not found (already deleted)", email)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.TenancyAPI.TenancyContactsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete contact %d (email: %s): %v", id, email, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted contact %d (email: %s)", id, email)

		}

	})

}

// RegisterContactRoleCleanup registers a cleanup function that will delete

// a contact role by slug after the test completes.

func (c *CleanupResource) RegisterContactRoleCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.TenancyAPI.TenancyContactRolesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list contact roles with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: contact role with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.TenancyAPI.TenancyContactRolesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete contact role %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted contact role %d (slug: %s)", id, slug)

		}

	})

}

// ============================================================================

// IPAM API CLEANUPS - Part 1

// ============================================================================

// RegisterVRFCleanup registers a cleanup function that will delete

// a VRF by name after the test completes.

func (c *CleanupResource) RegisterVRFCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamVrfsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list VRFs with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: VRF with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamVrfsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VRF %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VRF %d (name: %s)", id, name)

		}

	})

}

// RegisterVLANGroupCleanup registers a cleanup function that will delete

// a VLAN group by slug after the test completes.

func (c *CleanupResource) RegisterVLANGroupCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamVlanGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list VLAN groups with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: VLAN group with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamVlanGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VLAN group %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VLAN group %d (slug: %s)", id, slug)

		}

	})

}

// RegisterVLANCleanup registers a cleanup function that will delete

// a VLAN by vid after the test completes.

func (c *CleanupResource) RegisterVLANCleanup(vid int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamVlansList(ctx).Vid([]int32{vid}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list VLANs with VID %d: %v", vid, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: VLAN with VID %d not found (already deleted)", vid)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamVlansDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VLAN %d (VID: %d): %v", id, vid, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VLAN %d (VID: %d)", id, vid)

		}

	})

}

// RegisterPrefixCleanup registers a cleanup function that will delete

// a prefix by CIDR after the test completes.

func (c *CleanupResource) RegisterPrefixCleanup(prefix string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamPrefixesList(ctx).Prefix([]string{prefix}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list prefixes with CIDR %s: %v", prefix, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: prefix with CIDR %s not found (already deleted)", prefix)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamPrefixesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete prefix %d (CIDR: %s): %v", id, prefix, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted prefix %d (CIDR: %s)", id, prefix)

		}

	})

}

// RegisterIPAddressCleanup registers a cleanup function that will delete

// an IP address by address after the test completes.

func (c *CleanupResource) RegisterIPAddressCleanup(address string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamIpAddressesList(ctx).Address([]string{address}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IP addresses with address %s: %v", address, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: IP address %s not found (already deleted)", address)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamIpAddressesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IP address %d (address: %s): %v", id, address, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IP address %d (address: %s)", id, address)

		}

	})

}

// RegisterIPRangeCleanup registers a cleanup function that will delete

// an IP range by start address after the test completes.

func (c *CleanupResource) RegisterIPRangeCleanup(startAddress string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamIpRangesList(ctx).StartAddress([]string{startAddress}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IP ranges with start address %s: %v", startAddress, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: IP range with start address %s not found (already deleted)", startAddress)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamIpRangesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IP range %d (start address: %s): %v", id, startAddress, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IP range %d (start address: %s)", id, startAddress)

		}

	})

}

// RegisterASNRangeCleanup registers a cleanup function that will delete

// an ASN range by name after the test completes.

func (c *CleanupResource) RegisterASNRangeCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamAsnRangesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list ASN ranges with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: ASN range with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamAsnRangesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete ASN range %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted ASN range %d (name: %s)", id, name)

		}

	})

}

// RegisterRIRCleanup registers a cleanup function that will delete

// an RIR by slug after the test completes.

func (c *CleanupResource) RegisterRIRCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamRirsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list RIRs with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: RIR with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamRirsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete RIR %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted RIR %d (slug: %s)", id, slug)

		}

	})

}

// RegisterAggregateCleanup registers a cleanup function that will delete

// an aggregate by prefix after the test completes.

func (c *CleanupResource) RegisterAggregateCleanup(prefix string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamAggregatesList(ctx).Prefix(prefix).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list aggregates with prefix %s: %v", prefix, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: aggregate with prefix %s not found (already deleted)", prefix)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamAggregatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete aggregate %d (prefix: %s): %v", id, prefix, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted aggregate %d (prefix: %s)", id, prefix)

		}

	})

}

// RegisterRouteTargetCleanup registers a cleanup function that will delete

// a route target by name after the test completes.

func (c *CleanupResource) RegisterRouteTargetCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamRouteTargetsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list route targets with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: route target with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamRouteTargetsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete route target %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted route target %d (name: %s)", id, name)

		}

	})

}

// RegisterFHRPGroupCleanup registers a cleanup function that will delete

// an FHRP group by protocol and group_id after the test completes.

func (c *CleanupResource) RegisterFHRPGroupCleanup(protocol string, groupID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamFhrpGroupsList(ctx).Protocol([]string{protocol}).GroupId([]int32{groupID}).Execute()

		if err != nil || resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: FHRP group with protocol %s and group_id %d not found: %v", protocol, groupID, err)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamFhrpGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete FHRP group %d (protocol: %s, group_id: %d): %v", id, protocol, groupID, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted FHRP group %d (protocol: %s, group_id: %d)", id, protocol, groupID)

		}

	})

}

// RegisterFHRPGroupAssignmentCleanup registers a cleanup function that will delete

// FHRP group assignments after the test completes.

func (c *CleanupResource) RegisterFHRPGroupAssignmentCleanup(name string) {

	c.t.Cleanup(func() {

		// FHRP group assignments are cleaned up via cascade deletion

		// through interface/device cleanup, so no explicit cleanup needed

		c.t.Logf("Cleanup: FHRP group assignments for %s will be cleaned up via cascade deletion", name)

	})

}

// RegisterRoleCleanup registers a cleanup function that will delete

// a role by slug after the test completes.

func (c *CleanupResource) RegisterRoleCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamRolesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list roles with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: role with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamRolesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete role %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted role %d (slug: %s)", id, slug)

		}

	})

}

// RegisterServiceCleanup registers a cleanup function that will delete

// a service by name and device ID after the test completes.

func (c *CleanupResource) RegisterServiceCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamServicesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list services with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: service with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamServicesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete service %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted service %d (name: %s)", id, name)

		}

	})

}

// RegisterServiceTemplateCleanup registers a cleanup function that will delete

// a service template by name after the test completes.

func (c *CleanupResource) RegisterServiceTemplateCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.IpamAPI.IpamServiceTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list service templates with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: service template with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.IpamAPI.IpamServiceTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete service template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted service template %d (name: %s)", id, name)

		}

	})

}

// ============================================================================

// CIRCUITS API CLEANUPS

// ============================================================================

// RegisterProviderCleanup registers a cleanup function that will delete

// a circuit provider by slug after the test completes.

func (c *CleanupResource) RegisterProviderCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsProvidersList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list providers with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: provider with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsProvidersDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete provider %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted provider %d (slug: %s)", id, slug)

		}

	})

}

// RegisterCircuitTypeCleanup registers a cleanup function that will delete

// a circuit type by slug after the test completes.

func (c *CleanupResource) RegisterCircuitTypeCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsCircuitTypesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuit types with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: circuit type with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsCircuitTypesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete circuit type %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted circuit type %d (slug: %s)", id, slug)

		}

	})

}

// RegisterCircuitCleanup registers a cleanup function that will delete

// a circuit by CID after the test completes.

func (c *CleanupResource) RegisterCircuitCleanup(cid string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsCircuitsList(ctx).Cid([]string{cid}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuits with CID %s: %v", cid, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: circuit with CID %s not found (already deleted)", cid)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsCircuitsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete circuit %d (CID: %s): %v", id, cid, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted circuit %d (CID: %s)", id, cid)

		}

	})

}

// RegisterCircuitGroupCleanup registers a cleanup function that will delete

// a circuit group by name after the test completes.

func (c *CleanupResource) RegisterCircuitGroupCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuit groups with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: circuit group with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.CircuitsAPI.CircuitsCircuitGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete circuit group %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted circuit group %d (name: %s)", id, name)

		}

	})

}

// RegisterCircuitGroupAssignmentCleanup registers a cleanup function that will delete

// a circuit group assignment by ID after the test completes.

func (c *CleanupResource) RegisterCircuitGroupAssignmentCleanup(groupName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		// First find the circuit group by name

		groupList, resp, err := c.client.CircuitsAPI.CircuitsCircuitGroupsList(ctx).Name([]string{groupName}).Execute()

		if err != nil || resp.StatusCode != http.StatusOK || groupList.Count == 0 {

			c.t.Logf("Cleanup: circuit group with name %s not found, cannot cleanup assignments", groupName)

			return

		}

		groupID := groupList.Results[0].GetId()

		// Find all assignments for this group

		assignmentList, _, err := c.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsList(ctx).GroupId([]int32{groupID}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list circuit group assignments for group %d: %v", groupID, err)

			return

		}

		for _, assignment := range assignmentList.Results {

			_, err := c.client.CircuitsAPI.CircuitsCircuitGroupAssignmentsDestroy(ctx, assignment.GetId()).Execute()

			if err != nil {

				c.t.Logf("Cleanup: failed to delete circuit group assignment %d: %v", assignment.GetId(), err)

			} else {

				c.t.Logf("Cleanup: successfully deleted circuit group assignment %d", assignment.GetId())

			}

		}

	})

}

// ============================================================================

// VPN API CLEANUPS

// ============================================================================

// RegisterIKEProposalCleanup registers a cleanup function that will delete

// an IKE proposal by name after the test completes.

func (c *CleanupResource) RegisterIKEProposalCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIkeProposalsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IKE proposals with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: IKE proposal with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIkeProposalsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IKE proposal %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IKE proposal %d (name: %s)", id, name)

		}

	})

}

// RegisterIKEPolicyCleanup registers a cleanup function that will delete

// an IKE policy by name after the test completes.

func (c *CleanupResource) RegisterIKEPolicyCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIkePoliciesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IKE policies with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: IKE policy with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIkePoliciesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IKE policy %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IKE policy %d (name: %s)", id, name)

		}

	})

}

// RegisterIPSecProposalCleanup registers a cleanup function that will delete

// an IPSec proposal by name after the test completes.

func (c *CleanupResource) RegisterIPSecProposalCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIpsecProposalsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IPSec proposals with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: IPSec proposal with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIpsecProposalsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IPSec proposal %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IPSec proposal %d (name: %s)", id, name)

		}

	})

}

// RegisterIPSecPolicyCleanup registers a cleanup function that will delete

// an IPSec policy by name after the test completes.

func (c *CleanupResource) RegisterIPSecPolicyCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIpsecPoliciesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IPSec policies with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: IPSec policy with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIpsecPoliciesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IPSec policy %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IPSec policy %d (name: %s)", id, name)

		}

	})

}

// RegisterIPSecProfileCleanup registers a cleanup function that will delete

// an IPSec profile by name after the test completes.

func (c *CleanupResource) RegisterIPSecProfileCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnIpsecProfilesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list IPSec profiles with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: IPSec profile with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnIpsecProfilesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete IPSec profile %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted IPSec profile %d (name: %s)", id, name)

		}

	})

}

// RegisterTunnelGroupCleanup registers a cleanup function that will delete

// a tunnel group by name after the test completes.

func (c *CleanupResource) RegisterTunnelGroupCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnTunnelGroupsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnel groups with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: tunnel group with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnTunnelGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete tunnel group %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted tunnel group %d (name: %s)", id, name)

		}

	})

}

// RegisterTunnelCleanup registers a cleanup function that will delete

// a tunnel by name after the test completes.

func (c *CleanupResource) RegisterTunnelCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnTunnelsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnels with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: tunnel with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnTunnelsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete tunnel %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted tunnel %d (name: %s)", id, name)

		}

	})

}

// RegisterTunnelTerminationCleanup registers a cleanup function that will delete

// a tunnel termination by ID after the test completes.

func (c *CleanupResource) RegisterTunnelTerminationCleanup(tunnelName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		// First find the tunnel by name

		tunnelList, resp, err := c.client.VpnAPI.VpnTunnelsList(ctx).Name([]string{tunnelName}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnels with name %s: %v", tunnelName, err)

			return

		}

		if resp.StatusCode != http.StatusOK || tunnelList.Count == 0 {

			c.t.Logf("Cleanup: tunnel with name %s not found (already deleted)", tunnelName)

			return

		}

		tunnelID := tunnelList.Results[0].GetId()

		// List terminations for this tunnel

		termList, _, err := c.client.VpnAPI.VpnTunnelTerminationsList(ctx).TunnelId([]int32{tunnelID}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tunnel terminations for tunnel %d: %v", tunnelID, err)

			return

		}

		for _, term := range termList.Results {

			_, err := c.client.VpnAPI.VpnTunnelTerminationsDestroy(ctx, term.GetId()).Execute()

			if err != nil {

				c.t.Logf("Cleanup: failed to delete tunnel termination %d: %v", term.GetId(), err)

			} else {

				c.t.Logf("Cleanup: successfully deleted tunnel termination %d", term.GetId())

			}

		}

	})

}

// RegisterL2VPNCleanup registers a cleanup function that will delete

// an L2VPN by name after the test completes.

func (c *CleanupResource) RegisterL2VPNCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnL2vpnsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list L2VPNs with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: L2VPN with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnL2vpnsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete L2VPN %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted L2VPN %d (name: %s)", id, name)

		}

	})

}

// RegisterL2VPNTerminationCleanup registers a cleanup function that will delete

// an L2VPN termination by L2VPN ID after the test completes.

func (c *CleanupResource) RegisterL2VPNTerminationCleanup(l2vpnID int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VpnAPI.VpnL2vpnTerminationsList(ctx).L2vpnId([]int32{l2vpnID}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list L2VPN terminations with L2VPN ID %d: %v", l2vpnID, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list == nil || len(list.Results) == 0 {

			c.t.Logf("Cleanup: L2VPN termination with L2VPN ID %d not found (already deleted)", l2vpnID)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VpnAPI.VpnL2vpnTerminationsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete L2VPN termination %d (L2VPN ID: %d): %v", id, l2vpnID, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted L2VPN termination %d (L2VPN ID: %d)", id, l2vpnID)

		}

	})

}

// ============================================================================

// VIRTUALIZATION API CLEANUPS

// ============================================================================

// RegisterClusterTypeCleanup registers a cleanup function that will delete

// a cluster type by slug after the test completes.

func (c *CleanupResource) RegisterClusterTypeCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationClusterTypesList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list cluster types with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: cluster type with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationClusterTypesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete cluster type %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted cluster type %d (slug: %s)", id, slug)

		}

	})

}

// RegisterClusterCleanup registers a cleanup function that will delete

// a cluster by name after the test completes.

func (c *CleanupResource) RegisterClusterCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationClustersList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list clusters with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: cluster with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationClustersDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete cluster %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted cluster %d (name: %s)", id, name)

		}

	})

}

// RegisterVirtualMachineCleanup registers a cleanup function that will delete

// a virtual machine by name after the test completes.

func (c *CleanupResource) RegisterVirtualMachineCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list virtual machines with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: virtual machine with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationVirtualMachinesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete virtual machine %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted virtual machine %d (name: %s)", id, name)

		}

	})

}

// RegisterVMInterfaceCleanup registers a cleanup function that will delete

// a VM interface by name and virtual machine after the test completes.

func (c *CleanupResource) RegisterVMInterfaceCleanup(name string, vmName string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		// First check if the parent VM exists
		vmList, vmResp, vmErr := c.client.VirtualizationAPI.VirtualizationVirtualMachinesList(ctx).Name([]string{vmName}).Execute()
		if vmErr != nil || vmResp.StatusCode != http.StatusOK || vmList.Count == 0 {
			// VM doesn't exist, so interface is cascade-deleted
			c.t.Logf("Cleanup: VM interface with name %s not cleaned up (parent VM %s already deleted)", name, vmName)
			return
		}

		list, resp, err := c.client.VirtualizationAPI.VirtualizationInterfacesList(ctx).Name([]string{name}).VirtualMachine([]string{vmName}).Execute()

		if err != nil {
			// 404 means the interface doesn't exist
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				c.t.Logf("Cleanup: VM interface with name %s on VM %s not found (already deleted)", name, vmName)
				return
			}

			c.t.Logf("Cleanup: failed to list VM interfaces with name %s on VM %s: %v", name, vmName, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: VM interface with name %s on VM %s not found (already deleted)", name, vmName)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationInterfacesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete VM interface %d (name: %s, VM: %s): %v", id, name, vmName, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted VM interface %d (name: %s, VM: %s)", id, name, vmName)

		}

	})

}

// RegisterVirtualDiskCleanup registers a cleanup function that will delete

// a virtual disk by name after the test completes.

func (c *CleanupResource) RegisterVirtualDiskCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationVirtualDisksList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list virtual disks with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: virtual disk with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationVirtualDisksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete virtual disk %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted virtual disk %d (name: %s)", id, name)

		}

	})

}

// RegisterClusterGroupCleanup registers a cleanup function that will delete

// a cluster group by slug after the test completes.

func (c *CleanupResource) RegisterClusterGroupCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.VirtualizationAPI.VirtualizationClusterGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list cluster groups with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: cluster group with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.VirtualizationAPI.VirtualizationClusterGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete cluster group %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted cluster group %d (slug: %s)", id, slug)

		}

	})

}

// ============================================================================

// WIRELESS API CLEANUPS

// ============================================================================

// RegisterWirelessLinkCleanup registers a cleanup function that will delete

// a wireless link by ID after the test completes.

func (c *CleanupResource) RegisterWirelessLinkCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.WirelessAPI.WirelessWirelessLinksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete wireless link %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted wireless link %d", id)

		}

	})

}

// RegisterWirelessLANCleanup registers a cleanup function that will delete

// a wireless LAN by SSID after the test completes.

func (c *CleanupResource) RegisterWirelessLANCleanup(ssid string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.WirelessAPI.WirelessWirelessLansList(ctx).Ssid([]string{ssid}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list wireless LANs with SSID %s: %v", ssid, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: wireless LAN with SSID %s not found (already deleted)", ssid)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.WirelessAPI.WirelessWirelessLansDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete wireless LAN %d (SSID: %s): %v", id, ssid, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted wireless LAN %d (SSID: %s)", id, ssid)

		}

	})

}

// RegisterWirelessLANGroupCleanup registers a cleanup function that will delete

// a wireless LAN group by slug after the test completes.

func (c *CleanupResource) RegisterWirelessLANGroupCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.WirelessAPI.WirelessWirelessLanGroupsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list wireless LAN groups with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: wireless LAN group with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.WirelessAPI.WirelessWirelessLanGroupsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete wireless LAN group %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted wireless LAN group %d (slug: %s)", id, slug)

		}

	})

}

// ============================================================================

// EXTRAS API CLEANUPS

// ============================================================================

// RegisterJournalEntryCleanup registers a cleanup function that will delete

// a journal entry by ID after the test completes.

func (c *CleanupResource) RegisterJournalEntryCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.ExtrasAPI.ExtrasJournalEntriesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete journal entry %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted journal entry %d", id)

		}

	})

}

// RegisterCustomFieldChoiceSetCleanup registers a cleanup function that will delete

// a custom field choice set by ID after the test completes.

func (c *CleanupResource) RegisterCustomFieldChoiceSetCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete custom field choice set %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted custom field choice set %d", id)

		}

	})

}

// RegisterCustomLinkCleanup registers a cleanup function that will delete

// a custom link by ID after the test completes.

func (c *CleanupResource) RegisterCustomLinkCleanup(id int32) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		_, err := c.client.ExtrasAPI.ExtrasCustomLinksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete custom link %d: %v", id, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted custom link %d", id)

		}

	})

}

// RegisterTagCleanup registers a cleanup function that will delete

// a tag by slug after the test completes.

func (c *CleanupResource) RegisterTagCleanup(slug string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasTagsList(ctx).Slug([]string{slug}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list tags with slug %s: %v", slug, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: tag with slug %s not found (already deleted)", slug)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.ExtrasAPI.ExtrasTagsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete tag %d (slug: %s): %v", id, slug, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted tag %d (slug: %s)", id, slug)

		}

	})

}

// RegisterWebhookCleanup registers a cleanup function that will delete

// a webhook by name after the test completes.

func (c *CleanupResource) RegisterWebhookCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasWebhooksList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list webhooks with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: webhook with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.ExtrasAPI.ExtrasWebhooksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete webhook %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted webhook %d (name: %s)", id, name)

		}

	})

}

// RegisterExportTemplateCleanup registers a cleanup function that will delete

// an export template by name after the test completes.

func (c *CleanupResource) RegisterExportTemplateCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasExportTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list export templates with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: export template with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.ExtrasAPI.ExtrasExportTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete export template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted export template %d (name: %s)", id, name)

		}

	})

}

// RegisterConfigTemplateCleanup registers a cleanup function that will delete

// a config template by name after the test completes.

func (c *CleanupResource) RegisterConfigTemplateCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasConfigTemplatesList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list config templates with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: config template with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.ExtrasAPI.ExtrasConfigTemplatesDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete config template %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted config template %d (name: %s)", id, name)

		}

	})

}

// RegisterConfigContextCleanup registers a cleanup function that will delete

// a config context by name after the test completes.

func (c *CleanupResource) RegisterConfigContextCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasConfigContextsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list config contexts with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: config context with name %s not found (already deleted)", name)

			return

		}

	})

}

// RegisterCustomFieldCleanup registers a cleanup function that will delete

// a custom field by name after the test completes.

func (c *CleanupResource) RegisterCustomFieldCleanup(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasCustomFieldsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list custom fields with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: custom field with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.ExtrasAPI.ExtrasCustomFieldsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete custom field %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted custom field %d (name: %s)", id, name)

		}

	})

}

// RegisterCustomFieldChoiceSetCleanupByName registers a cleanup function that will delete

// a custom field choice set by name after the test completes.

func (c *CleanupResource) RegisterCustomFieldChoiceSetCleanupByName(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list custom field choice sets with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: custom field choice set with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.ExtrasAPI.ExtrasCustomFieldChoiceSetsDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete custom field choice set %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted custom field choice set %d (name: %s)", id, name)

		}

	})

}

// RegisterCustomLinkCleanupByName registers a cleanup function that will delete

// a custom link by name after the test completes.

func (c *CleanupResource) RegisterCustomLinkCleanupByName(name string) {

	c.t.Cleanup(func() {

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		list, resp, err := c.client.ExtrasAPI.ExtrasCustomLinksList(ctx).Name([]string{name}).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to list custom links with name %s: %v", name, err)

			return

		}

		if resp.StatusCode != http.StatusOK || list.Count == 0 {

			c.t.Logf("Cleanup: custom link with name %s not found (already deleted)", name)

			return

		}

		id := list.Results[0].GetId()

		_, err = c.client.ExtrasAPI.ExtrasCustomLinksDestroy(ctx, id).Execute()

		if err != nil {

			c.t.Logf("Cleanup: failed to delete custom link %d (name: %s): %v", id, name, err)

		} else {

			c.t.Logf("Cleanup: successfully deleted custom link %d (name: %s)", id, name)

		}

	})

}
