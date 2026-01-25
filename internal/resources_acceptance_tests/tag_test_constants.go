package resources_acceptance_tests

// Common tag configuration string constants used across multiple test files
// to satisfy goconst linter requirements.
const (
	// Case identifiers for tag test configurations.
	caseTag1Tag2    = "tag1,tag2"
	caseTag1Uscore2 = "tag1_tag2"
	caseTag2Uscore1 = "tag2_tag1"
	caseTag3        = "tag3"

	// Tag configuration strings for slug list format.
	tagsEmpty              = "tags = []"
	tagsSingleSlug         = "tags = [netbox_tag.tag3.slug]"
	tagsDoubleSlug         = "tags = [netbox_tag.tag1.slug, netbox_tag.tag2.slug]"
	tagsDoubleSlugReversed = "tags = [netbox_tag.tag2.slug, netbox_tag.tag1.slug]"
)
