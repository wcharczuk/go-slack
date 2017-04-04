package slack

import (
	"crypto/rand"
	"fmt"
)

// UUID represents a unique id.
type UUID []byte

// ToFullString returns the full string representation.
func (uuid UUID) ToFullString() string {
	b := []byte(uuid)
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// ToShortString returns the short string representation.
func (uuid UUID) ToShortString() string {
	b := []byte(uuid)
	return fmt.Sprintf("%x", b[:])
}

// Version returns the version.
func (uuid UUID) Version() byte {
	return uuid[6] >> 4
}

// UUIDv4 returns a new UUID version 4.
func UUIDv4() UUID {
	uuid := make([]byte, 16)
	rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80 // Variant is 10
	return uuid
}
