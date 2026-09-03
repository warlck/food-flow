// Package opt provides optional and nullable types for partial JSON DTO updates.
package opt

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// NullInt represents an optional integer field that distinguishes between:
// - omitted (Present: false, Value: nil)
// - explicit value (Present: true, Value: &val)
// - explicit null (Present: true, Value: nil)
type NullInt struct {
	Present bool
	Value   *int
}

// NewNullInt constructs a NullInt set to an explicit integer value.
func NewNullInt(v int) NullInt {
	return NullInt{
		Present: true,
		Value:   &v,
	}
}

// NewNullIntNull constructs a NullInt set explicitly to null.
func NewNullIntNull() NullInt {
	return NullInt{
		Present: true,
		Value:   nil,
	}
}

// UnmarshalJSON implements custom JSON unmarshaling to distinguish between
// explicit null and values.
func (n *NullInt) UnmarshalJSON(data []byte) error {
	n.Present = true
	if bytes.Equal(data, []byte("null")) {
		n.Value = nil
		return nil
	}

	var val int
	if err := json.Unmarshal(data, &val); err != nil {
		return fmt.Errorf("invalid int: %w", err)
	}

	n.Value = &val
	return nil
}

// MarshalJSON implements custom JSON marshaling.
func (n NullInt) MarshalJSON() ([]byte, error) {
	if !n.Present || n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.Value)
}

// Equal implements equality comparison for tests and cmp.Diff.
func (n NullInt) Equal(other NullInt) bool {
	if n.Present != other.Present {
		return false
	}
	if n.Value == nil && other.Value == nil {
		return true
	}
	if n.Value != nil && other.Value != nil {
		return *n.Value == *other.Value
	}
	return false
}
