package uuip

import (
	"net/netip"
	"testing"

	"github.com/google/uuid"
)

func TestUUIDToIPv6AndBack(t *testing.T) {
	for i := 0; i < 10; i++ {
		originalUUID := uuid.New()

		ipv6 := UUIDToIPv6(originalUUID)

		decodedUUID, err := IPv6ToUUID(ipv6)
		if err != nil {
			t.Fatalf("IPv6ToUUID failed: %v", err)
		}

		if originalUUID != decodedUUID {
			t.Errorf("Round-trip failed: expected %s, got %s", originalUUID, decodedUUID)
		}
	}
}

func TestIPv6ToUUID_InvalidAddress(t *testing.T) {
	ipv4Addr := netip.MustParseAddr("192.0.2.1")

	_, err := IPv6ToUUID(ipv4Addr)
	if err == nil {
		t.Error("Expected error for non-IPv6 address, got nil")
	}
}

func TestUUIDToIPv6_InvalidLength(t *testing.T) {
	// Invalid UUID (less than 16 bytes)
	var shortUUID uuid.UUID
	copy(shortUUID[:], []byte("short"))

	// Should still pass because uuid.UUID is always 16 bytes (even if garbage)
	ipv6 := UUIDToIPv6(shortUUID)

	// Check it's a valid IPv6 address
	if !ipv6.Is6() {
		t.Errorf("Expected IPv6 address, got: %v", ipv6)
	}
}
