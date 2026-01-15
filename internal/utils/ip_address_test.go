package utils

import "testing"

func TestNormalizeIPAddress_IPv6CIDR(t *testing.T) {
	in := "fd00:0f13:e5bf::ced6/128"
	want := "fd00:f13:e5bf::ced6/128"

	if got := NormalizeIPAddress(in); got != want {
		t.Fatalf("NormalizeIPAddress(%q) = %q, want %q", in, got, want)
	}
}

func TestNormalizeIPAddress_IPv6NoPrefix(t *testing.T) {
	in := "fd00:0f13:e5bf::ced6"
	want := "fd00:f13:e5bf::ced6"

	if got := NormalizeIPAddress(in); got != want {
		t.Fatalf("NormalizeIPAddress(%q) = %q, want %q", in, got, want)
	}
}

func TestNormalizeIPAddress_IPv4Unchanged(t *testing.T) {
	tests := []string{
		"203.0.113.179/32",
		"203.0.113.179",
	}

	for _, in := range tests {
		if got := NormalizeIPAddress(in); got != in {
			t.Fatalf("NormalizeIPAddress(%q) = %q, want unchanged", in, got)
		}
	}
}

func TestNormalizeIPAddress_InvalidUnchanged(t *testing.T) {
	in := "not-an-ip/64"
	if got := NormalizeIPAddress(in); got != in {
		t.Fatalf("NormalizeIPAddress(%q) = %q, want unchanged", in, got)
	}
}
