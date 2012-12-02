package semver

import (
	"testing"
)

func BenchmarkParseSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse("0.1.0")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParsePreRel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse("14.2.9-rc1.3")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseBuild(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse("14.2.9+build.11.e0f985a")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseFull(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse("14.2.9-rc1.3+build.11.e0f985a")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseBigInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Parse("2147483648.0.0") // = 1 << 32 - 1
		if err != nil {
			b.Fatal(err)
		}
	}
}
