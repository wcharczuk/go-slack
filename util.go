package slack

import "time"

// OptionalUInt8 is a convenience method.
func OptionalUInt8(value uint8) *uint8 {
	return &value
}

// OptionalByte is a convenience method.
func OptionalByte(value byte) *byte {
	return &value
}

// OptionalUInt16 is a convenience method.
func OptionalUInt16(value uint16) *uint16 {
	return &value
}

// OptionalUInt is a convenience method.
func OptionalUInt(value uint) *uint {
	return &value
}

// OptionalUInt32 is a convenience method.
func OptionalUInt32(value uint32) *uint32 {
	return &value
}

// OptionalUInt64 is a convenience method.
func OptionalUInt64(value uint64) *uint64 {
	return &value
}

// OptionalInt16 is a convenience method.
func OptionalInt16(value int16) *int16 {
	return &value
}

// OptionalInt is a convenience method.
func OptionalInt(value int) *int {
	return &value
}

// OptionalInt32 is a convenience method.
func OptionalInt32(value int32) *int32 {
	return &value
}

// OptionalInt64 is a convenience method.
func OptionalInt64(value int64) *int64 {
	return &value
}

// OptionalFloat32 is a convenience method.
func OptionalFloat32(value float32) *float32 {
	return &value
}

// OptionalFloat64 is a convenience method.
func OptionalFloat64(value float64) *float64 {
	return &value
}

// OptionalString is a convenience method.
func OptionalString(value string) *string {
	return &value
}

// OptionalBool is a convenience method.
func OptionalBool(value bool) *bool {
	return &value
}

// OptionalTime is a convenience method.
func OptionalTime(value time.Time) *time.Time {
	return &value
}

// OptionalTimestamp is a convenience method.
func OptionalTimestamp(value Timestamp) *Timestamp {
	return &value
}

// IsEmpty returns if a string is empty or not.
func IsEmpty(s string) bool {
	return len(s) == 0
}
