package acceptance

// The changelog entry fixtures the changelog checks' requirement tests share
// with the repo's own.

// EntryFileText is an entry file under stamp with three typed bullets, the
// first one item, below the front matter older entries carry.
func EntryFileText(stamp, item string) string {
	return "---\nrelatedConfigurations: ['../.fitnessrc.json']\n---\n\n" + BareEntryText(stamp, item)
}

// BareEntryText is EntryFileText without front matter, which is optional.
func BareEntryText(stamp, item string) string {
	return "### " + stamp + "\n\n- " + item + "\n- Test: cover it.\n- Docs: describe it.\n"
}

// WithEntries is files over happyRepo holding each named entry file.
func WithEntries(entries map[string]string) map[string]string {
	files := map[string]string{}
	for name, body := range entries {
		files["changelog/"+name+".md"] = body
	}
	return files
}

// UpstreamChangelog is a CHANGELOG.md as an upstream project keeps it:
// undated, typed sections of numbered items under a version heading.
const UpstreamChangelog = "# Changelog\n\n## v2.0.0\n\n### Added\n\n1. Catalog-level publishing\n2. Plugin hooks\n\n### Removed\n1. Old lint tools\n"
