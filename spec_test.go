package semver

import (
	"testing"
)

// Semantic Versioning Specification (SemVer) 2.0.0

// 1. Software using Semantic Versioning MUST declare a public API. This API
// could be declared in the code itself or exist strictly in documentation.
// However it is done, it should be precise and comprehensive.

// No tests writable

// 2. A normal version number MUST take the form X.Y.Z where X, Y, and Z are
// non-negative integers, and MUST NOT contain leading zeroes. X is the major
// version, Y is the minor version, and Z is the patch version. Each element
// MUST increase numerically. For instance: 1.9.0 -> 1.10.0 -> 1.11.0.

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

// 3. Once a versioned package has been released, the contents of that version
// MUST NOT be modified. Any modifications MUST be released as a new version.

// No tests writable

// 4. Major version zero (0.y.z) is for initial development. Anything may
// change at any time. The public API should not be considered stable.

// No tests writable

// 5. Version 1.0.0 defines the public API. The way in which the version
// number is incremented after this release is dependent on this public API
// and how it changes.

// No tests writable

// 6. Patch version Z (x.y.Z | x > 0) MUST be incremented if only backwards
// compatible bug fixes are introduced. A bug fix is defined as an internal
// change that fixes incorrect behavior.

// No tests writable

// 7. Minor version Y (x.Y.z | x > 0) MUST be incremented if new, backwards
// compatible functionality is introduced to the public API. It MUST be
// incremented if any public API functionality is marked as deprecated. It MAY
// be incremented if substantial new functionality or improvements are
// introduced within the private code. It MAY include patch level changes.
// Patch version MUST be reset to 0 when minor version is incremented.

// No tests writable

// 8. Major version X (X.y.z | X > 0) MUST be incremented if any backwards
// incompatible changes are introduced to the public API. It MAY include minor
// and patch level changes. Patch and minor version MUST be reset to 0 when
// major version is incremented.

// No tests writable

// 9. A pre-release version MAY be denoted by appending a hyphen and a series
// of dot separated identifiers immediately following the patch version.
// Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-].
// Identifiers MUST NOT be empty. Numeric identifiers MUST NOT include leading
// zeroes. Pre-release versions have a lower precedence than the associated
// normal version. A pre-release version indicates that the version is
// unstable and might not satisfy the intended compatibility requirements as
// denoted by its associated normal version. Examples: 1.0.0-alpha,
// 1.0.0-alpha.1, 1.0.0-0.3.7, 1.0.0-x.7.z.92.

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

// 10. Build metadata MAY be denoted by appending a plus sign and a series of
// dot separated identifiers immediately following the patch or pre-release
// version. Identifiers MUST comprise only ASCII alphanumerics and hyphen [0
// -9A-Za-z-]. Identifiers MUST NOT be empty. Build metadata SHOULD be ignored
// when determining version precedence. Thus two versions that differ only in
// the build metadata, have the same precedence. Examples: 1.0.0-alpha+001,
// 1.0.0+20130313144700, 1.0.0-beta+exp.sha.5114f85.

func TestBuildPrecedence(t *testing.T) {
	for _, vs := range []string{
		"1.0.0-alpha+001",
		"1.0.0+20130313144700",
		"1.0.0-beta+exp.sha.5114f85",
	} {
		v, err := Parse(vs)
		if err != nil {
			t.Errorf(`Couldn't parse version string "%s": %s`,
				vs, err)
			continue
		}
		normal := Version{v.Major, v.Minor, v.Patch, v.PreRelease, nil}
		if Less(normal, v) {
			t.Errorf(`Failed assertion: !("%s" < "%s)"`,
				normal, v)
		}
		if Less(v, normal) {
			t.Errorf(`Failed assertion: !("%s" < "%s)"`,
				v, normal)
		}
	}
}

func TestEqualBuildPrecedence(t *testing.T) {
	tests := []string{
		"1.0.0+ab.cd",
		"1.0.0+ef.gh",
		"1.0.0+1.99",
		"1.0.0+1.3",
		"1.0.0+ab.15.0.1",
		"1.0.0+asdf",
		"1.0.0+9001",
	}

	for i, vs := range tests {
		v, err := Parse(vs)
		if err != nil {
			t.Errorf(`Couldn't parse version string "%s": %s`,
				vs, err)
			continue
		}
		for _, v2s := range tests[i+1:] {
			v2, err := Parse(v2s)
			if err != nil {
				t.Errorf(`Couldn't parse version string "%s"`+
					`: %s`, v2s, err)
				continue
			}
			if Less(v, v2) {
				t.Errorf(`Failed assertion: !("%s" < "%s")`,
					v, v2)
			}
			if Less(v2, v) {
				t.Errorf(`Failed assertion: !("%s" < "%s")`,
					v2, v)
			}
		}
	}
}

// 11. Precedence refers to how versions are compared to each other when
// ordered. Precedence MUST be calculated by separating the version into
// major, minor, patch and pre-release identifiers in that order (Build
// metadata does not figure into precedence). Precedence is determined by the
// first difference when comparing each of these identifiers from left to
// right as follows: Major, minor, and patch versions are always compared
// numerically. Example: 1.0.0 < 2.0.0 < 2.1.0 < 2.1.1. When major, minor, and
// patch are equal, a pre-release version has lower precedence than a normal
// version. Example: 1.0.0-alpha < 1.0.0. Precedence for two pre-release
// versions with the same major, minor, and patch version MUST be determined
// by comparing each dot separated identifier from left to right until a
// difference is found as follows: identifiers consisting of only digits are
// compared numerically and identifiers with letters or hyphens are compared
// lexically in ASCII sort order. Numeric identifiers always have lower
// precedence than non-numeric identifiers. A larger set of pre-release fields
// has a higher precedence than a smaller set, if all of the preceding
// identifiers are equal. Example: 1.0.0-alpha < 1.0.0-alpha.1 <
// 1.0.0-alpha.beta < 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 <
// 1.0.0.

func TestSortOrder(t *testing.T) {
	versionstrings := []string{
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0",
		"2.0.0",
		"2.1.0",
		"2.1.1",
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
