package resources_acceptance_tests

// Common tag configuration string constants used across multiple test files
// to satisfy goconst linter requirements.
const (
	// Case identifiers for tag test configurations.
	caseTag1Tag2    = "tag1,tag2"
	caseTag1Uscore2 = "tag1_tag2"
	caseTag2Uscore1 = "tag2_tag1"
	caseTag3        = "tag3"

	// Tag configuration strings for nested tag format (Phase 2 conversion needed).
	tagsEmpty        = "tags = []"
	tagsSingleNested = `tags = [
    { name = netbox_tag.tag3.name, slug = netbox_tag.tag3.slug }
  ]`
	tagsDoubleNested = `tags = [
    { name = netbox_tag.tag1.name, slug = netbox_tag.tag1.slug },
    { name = netbox_tag.tag2.name, slug = netbox_tag.tag2.slug }
  ]`
	tagsDoubleNestedReversed = `tags = [
    { name = netbox_tag.tag2.name, slug = netbox_tag.tag2.slug },
    { name = netbox_tag.tag1.name, slug = netbox_tag.tag1.slug }
  ]`

	// Tag configuration strings for slug list format.
	tagsSingleSlug         = "tags = [netbox_tag.tag3.slug]"
	tagsDoubleSlug         = "tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]"
	tagsDoubleSlugReversed = "tags = [netbox_tag.tag2.slug, netbox_tag.tag1.slug]"
)
