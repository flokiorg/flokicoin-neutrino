package neutrino

import (
	"os"
	"regexp"
	"testing"
)

var changelogHeading = regexp.MustCompile(`(?m)^## \[([^\]]+)\]`)

// TestUserAgentVersionMatchesChangelog keeps UserAgentVersion honest.
//
// Consumers of this library advertise UserAgentVersion to their peers, so a
// stale value misrepresents every node built against it -- it had been pinned
// at 0.16.4 across several releases. A library has no main package and so no
// build-time version injection, which leaves CHANGELOG.md as the source of
// truth and this test as the thing that stops the two drifting apart.
func TestUserAgentVersionMatchesChangelog(t *testing.T) {
	data, err := os.ReadFile("CHANGELOG.md")
	if err != nil {
		t.Fatalf("reading CHANGELOG.md: %v", err)
	}

	m := changelogHeading.FindSubmatch(data)
	if m == nil {
		t.Fatal("no '## [X.Y.Z]' heading found in CHANGELOG.md")
	}

	if got, want := UserAgentVersion, string(m[1]); got != want {
		t.Errorf("UserAgentVersion = %q, but the topmost CHANGELOG.md heading "+
			"is %q; update UserAgentVersion in neutrino.go to match", got, want)
	}
}
