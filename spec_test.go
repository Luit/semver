package semver

import (
	"testing"
)

// Semantic Versioning Specification (SemVer) 2.0.0-rc1

// 1. Software using Semantic Versioning MUST declare a public API. This API
// could be declared in the code itself or exist strictly in documentation.
// However it is done, it should be precise and comprehensive.

// No tests writable

// 2. A normal version number MUST take the form X.Y.Z where X, Y, and Z are
// non- negative integers. X is the major version, Y is the minor version, and
// Z is the patch version. Each element MUST increase numerically by
// increments of one. For instance: 1.9.0 -> 1.10.0 -> 1.11.0.

func TestSimpleIncrement(t *testing.T) {
	v := Version{
		Major: 1,
		Minor: 9,
	}
	if v.String() != "1.9.0" {
		t.Fatalf(`Expected "1.9.0", got "%s"`, v.String())
	}
	v = v.NextMinor()
	if v.String() != "1.10.0" {
		t.Fatalf(`Expected "1.10.0", got "%s"`, v.String())
	}
	v = v.NextMinor()
	if v.String() != "1.11.0" {
		t.Fatalf(`Expected "1.11.0", got "%s"`, v.String())
	}
}

// 3. When a major version number is incremented, the minor version and patch
// version MUST be reset to zero. When a minor version number is incremented,
// the patch version MUST be reset to zero. For instance: 1.1.3 -> 2.0.0 and
// 2.1.7 -> 2.2.0.

func TestZeroingIncrement(t *testing.T) {
	v := Version{
		Major: 1,
		Minor: 1,
		Patch: 3,
	}
	if v.String() != "1.1.3" {
		t.Fatalf(`Expected "1.1.3", got "%s"`, v.String())
	}
	v = v.NextMajor()
	if v.String() != "2.0.0" {
		t.Fatalf(`Expected "2.0.0", got "%s"`, v.String())
	}

	v, err := Parse("2.1.7")
	if err != nil {
		t.Fatalf(`Unable to parse version string "2.1.7": %s`, err)
	}
	if v.String() != "2.1.7" {
		t.Fatalf(`Expected "2.1.7", got "%s"`, v.String())
	}
	v = v.NextMinor()
	if v.String() != "2.2.0" {
		t.Fatalf(`Expected "2.2.0", got "%s"`, v.String())
	}

}

// 4. Once a versioned package has been released, the contents of that version
// MUST NOT be modified. Any modifications must be released as a new version.

// No tests writable

// 5. Major version zero (0.y.z) is for initial development. Anything may
// change at any time. The public API should not be considered stable.

// No tests writable

// 6. Version 1.0.0 defines the public API. The way in which the version
// number is incremented after this release is dependent on this public API
// and how it changes.

// No tests writable

// 7. Patch version Z (x.y.Z | x > 0) MUST be incremented if only backwards
// compatible bug fixes are introduced. A bug fix is defined as an internal
// change that fixes incorrect behavior.

// No tests writable

// 8. Minor version Y (x.Y.z | x > 0) MUST be incremented if new, backwards
// compatible functionality is introduced to the public API. It MUST be
// incremented if any public API functionality is marked as deprecated. It MAY
// be incremented if substantial new functionality or improvements are
// introduced within the private code. It MAY include patch level changes.
// Patch version MUST be reset to 0 when minor version is incremented.

// No tests writable

// 9. Major version X (X.y.z | X > 0) MUST be incremented if any backwards
// incompatible changes are introduced to the public API. It MAY include minor
// and patch level changes. Patch and minor version MUST be reset to 0 when
// major version is incremented.

// No tests writable

// 10. A pre-release version MAY be denoted by appending a dash and a series
// of dot separated identifiers immediately following the patch version.
// Identifiers MUST be comprised of only ASCII alphanumerics and dash [0-9A-
// Za-z-]. Pre-release versions satisfy but have a lower precedence than the
// associated normal version. Examples: 1.0.0-alpha, 1.0.0-alpha.1,
// 1.0.0-0.3.7, 1.0.0-x.7.z.92.

func TestPreReleasePrecedence(t *testing.T) {
	for _, vs := range []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-0.3.7",
		"1.0.0-x.7.z.92",
	} {
		v, err := Parse(vs)
		if err != nil {
			t.Errorf(`Couldn't parse version string "%s": %s`,
				vs, err)
			continue
		}
		normal := Version{v.Major, v.Minor, v.Patch, nil, nil}
		if !Less(v, normal) {
			t.Errorf(`Failed assertion: "%s" < "%s"`, v, normal)
		}
	}
}

// 11. A build version MAY be denoted by appending a plus sign and a series of
// dot separated identifiers immediately following the patch version or pre-
// release version. Identifiers MUST be comprised of only ASCII alphanumerics
// and dash [0-9A-Za-z-]. Build versions satisfy and have a higher precedence
// than the associated normal version. Examples: 1.0.0+build.1,
// 1.3.7+build.11.e0f985a.

func TestBuildPrecedence(t *testing.T) {
	for _, vs := range []string{
		"1.0.0+build.1",
		"1.3.7+build.11.e0f985a",
	} {
		v, err := Parse(vs)
		if err != nil {
			t.Errorf(`Couldn't parse version string "%s": %s`,
				vs, err)
			continue
		}
		normal := Version{v.Major, v.Minor, v.Patch, nil, nil}
		if !Less(normal, v) {
			t.Errorf(`Failed assertion: "%s" < "%s"`, normal, v)
		}
	}
}

// 12. Precedence MUST be calculated by separating the version into major,
// minor, patch, pre-release, and build identifiers in that order. Major,
// minor, and patch versions are always compared numerically. Pre-release and
// build version precedence MUST be determined by comparing each dot separated
// identifier as follows: identifiers consisting of only digits are compared
// numerically and identifiers with letters or dashes are compared lexically
// in ASCII sort order. Numeric identifiers always have lower precedence than
// non-numeric identifiers. Example: 1.0.0-alpha < 1.0.0-alpha.1 <
// 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0-rc.1+build.1 < 1.0.0 <
// 1.0.0+0.3.7 < 1.3.7+build < 1.3.7+build.2.b8f12d7 < 1.3.7+build.11.e0f985a.

func TestSortOrder(t *testing.T) {
	versionstrings := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0-rc.1+build.1",
		"1.0.0",
		"1.0.0+0.3.7",
		"1.3.7+build",
		"1.3.7+build.2.b8f12d7",
		"1.3.7+build.11.e0f985a",
	}

	versions := make([]Version, len(versionstrings))
	for n := range versions {
		v, err := Parse(versionstrings[n])
		if err != nil {
			t.Fatalf(`Couldn't parse version string "%s": %s`,
				versionstrings[n], err)
		}
		versions[n] = v
	}

	for n := range versions {
		if n == 0 {
			continue
		}

		if !Less(versions[n-1], versions[n]) {
			t.Errorf(`Failed assertion: "%s" < "%s"`,
				versions[n-1], versions[n])
		}
		if Less(versions[n], versions[n-1]) {
			t.Errorf(`Failed assertion: !("%s" < "%s")`,
				versions[n], versions[n-1])
		}
	}
}
