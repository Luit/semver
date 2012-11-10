// Package semver provides utilities for working with the Semantic Versioning
// Specification (SemVer) 2.0.0-rc1
package semver

import (
	"bytes"
	"errors"
	"regexp"
	"strconv"
)

// Version is a Semantic Version number, represented as three unsigned integer
// values Major, Minor and Patch, and two slices of []byte PreRelease and
// Build. Each of the []byte values of PreRelease and Build are dot-separated
// in the string representation.
type Version struct {
	Major, Minor, Patch uint
	PreRelease, Build   [][]byte
}

var versionExp = regexp.MustCompile(`([0-9])\.([0-9])\.([0-9])(\-[0-9A-Za-z\-\.]+)?(\+[0-9A-Za-z\-\.]+)?`)
var preReleaseExp = regexp.MustCompile(`[-.][0-9A-Za-z-]*`)
var buildExp = regexp.MustCompile(`[+.][0-9A-Za-z-]*`)

// Turn a version string into a Version. Will return an error if it can't
// recognize a proper (semantic) version number in this string.
func Parse(version string) (v Version, err error) {
	// TODO: rewrite properly, preferably without regexp
	matches := versionExp.FindSubmatch([]byte(version))
	if len(matches) != 6 {
		return v, errors.New("Regexp didn't match as expected")
	}
	major, err := strconv.ParseUint(string(matches[1]), 10, 64) // TODO: hard-code 64 here? Decide!
	if err != nil {
		return
	}
	v.Major = uint(major)

	minor, err := strconv.ParseUint(string(matches[2]), 10, 64) // TODO: hard-code 64 here? Decide!
	if err != nil {
		return
	}
	v.Minor = uint(minor)

	patch, err := strconv.ParseUint(string(matches[3]), 10, 64) // TODO: hard-code 64 here? Decide!
	if err != nil {
		return
	}
	v.Patch = uint(patch)

	v.PreRelease = preReleaseExp.FindAll(matches[4], -1)
	for n := range v.PreRelease {
		v.PreRelease[n] = v.PreRelease[n][1:] // Strip dot and dash
	}
	v.Build = buildExp.FindAll(matches[5], -1)
	for n := range v.Build {
		v.Build[n] = v.Build[n][1:] // Strip dot and plus
	}

	return
}

// TODO: could be unpredictible with hex identifiers that look decimal; is
// this the best way?
func chunkCompare(i, j []byte) int {
	var ii, ij uint64 // goto can't jump over declaration

	ii, err := strconv.ParseUint(string(i), 10, 64) // TODO: hard-code 64 here? Decide!
	if err != nil {
		goto compareBytewise
	}

	ij, err = strconv.ParseUint(string(j), 10, 64) // TODO: hard-code 64 here? Decide!
	if err != nil {
		goto compareBytewise
	}

	if ii == ij {
		return 0
	}
	if ii < ij {
		return -1
	}
	return 1

compareBytewise:
	return bytes.Compare(i, j)
}

func preReleaseCompare(i, j [][]byte) int {
	li, lj := len(i), len(j)
	if li == 0 && lj != 0 {
		return 1 // No PreRelease means i is greater version
	}
	if li != 0 && lj == 0 {
		return -1
	}

	for n := range j {
		if n == li {
			return -1 // j is more specific, so i is lesser version
		}
		c := chunkCompare(i[n], j[n])
		if c == 0 {
			continue // undecided, look at next chunk
		}
		return c
	}

	if li == lj {
		return 0
	}
	return 1 // i is more specific, so i is greater version
}

func buildCompare(i, j [][]byte) int {
	li, lj := len(i), len(j)
	if li == 0 && lj != 0 {
		return -1 // No build means i is lesser version
	}
	if li != 0 && lj == 0 {
		return 1
	}

	for n := range j {
		if n == li {
			return -1 // j is more specific, so i is lesser version
		}
		c := chunkCompare(i[n], j[n])
		if c == 0 {
			continue // undecided, look at next chunk
		}
		return c
	}

	if li == lj {
		return 0
	}
	return 1 // i is more specific, so i is greater version
}

// Less tests precedence of Version i over Version j.
func Less(i, j Version) bool {
	if i.Major < j.Major {
		return true
	}
	if i.Minor < j.Minor {
		return true
	}
	if i.Patch < j.Patch {
		return true
	}

	switch preReleaseCompare(i.PreRelease, j.PreRelease) {
	case -1:
		return true
	case 1:
		return false
	}
	return buildCompare(i.Build, j.Build) == -1
}

// Returns the string representation of Version v.
func (v Version) String() string {
	var pre []byte
	for _, p := range v.PreRelease {
		pre = append(pre, '.')
		pre = append(pre, p...)
	}
	if len(pre) > 0 {
		pre[0] = '-'
	}
	var build []byte
	for _, b := range v.Build {
		build = append(build, '.')
		build = append(build, b...)
	}
	if len(build) > 0 {
		build[0] = '+'
	}
	return strconv.FormatUint(uint64(v.Major), 10) + "." +
		strconv.FormatUint(uint64(v.Minor), 10) + "." +
		strconv.FormatUint(uint64(v.Patch), 10) +
		string(pre) + string(build)
}

// Returns the next major Version.
func (v Version) NextMajor() Version {
	return Version{
		Major: v.Major + 1,
		Minor: 0,
		Patch: 0,
	}
}

// Returns the next minor Version.
func (v Version) NextMinor() Version {
	return Version{
		Major: v.Major,
		Minor: v.Minor + 1,
		Patch: 0,
	}
}

// Returns the next patch Version.
func (v Version) NextPatch() Version {
	return Version{
		Major: v.Major,
		Minor: v.Minor,
		Patch: v.Patch + 1,
	}
}
