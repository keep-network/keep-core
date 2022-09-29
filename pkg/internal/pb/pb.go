// Package pb defines interfaces used for self-marshaling and self-unmarshaling
// of objects defined in our codebase.
package pb

// Marshaler is the interface representing objects that can marshal themselves.
type Marshaler interface {
	Marshal() ([]byte, error)
}

// Unmarshaler is the interface representing objects that can
// unmarshal themselves.
type Unmarshaler interface {
	Unmarshal([]byte) error
}
