package uuip

import (
	"errors"
	"net/netip"

	"github.com/google/uuid"
)

// UUIDToIPv6 converts a UUID to a netip.Addr (IPv6).
func UUIDToIPv6(u uuid.UUID) netip.Addr {
	return netip.AddrFrom16(u)
}

// IPv6ToUUID converts a netip.Addr (IPv6) to a UUID.
func IPv6ToUUID(addr netip.Addr) (uuid.UUID, error) {
	if !addr.Is6() {
		return uuid.Nil, errors.New("address is not IPv6")
	}
	addrBytes := addr.As16()
	return uuid.FromBytes(addrBytes[:])
}
