// Package bitcoin defines types and interfaces required to work with the
// Bitcoin chain. This package is meant to reflect the part of the Bitcoin
// protocol to an extent required by the client applications. That is,
// this package can only hold the components specific to the Bitcoin domain
// and must remain free of application-specific items. Third-party helper
// libraries are allowed for use within this package though no external
// type should leak outside.
package bitcoin

// CompactSizeUint is a documentation type that is supposed to capture the
// details of the Bitcoin's CompactSize Unsigned Integer as described in:
// https://developer.bitcoin.org/reference/transactions.html#compactsize-unsigned-integers
type CompactSizeUint uint64

// ByteOrder represents the byte order used by the Bitcoin byte arrays. The
// Bitcoin ecosystem is not totally consistent in this regard and different
// byte orders are used depending on the purpose.
type ByteOrder int

const (
	// InternalByteOrder represents the internal byte order used by the Bitcoin
	// protocol. This is the primary byte order that is suitable for the
	// use cases related with the protocol logic and cryptography. Byte arrays
	// using this byte order should be converted to numbers according to
	// the little-endian sequence.
	InternalByteOrder = iota

	// ReversedByteOrder represents the "human" byte order. This is the
	// byte order that is typically used by the third party services like
	// block explorers or Bitcoin chain clients. Byte arrays using this byte
	// order should be converted to numbers according to the big-endian
	// sequence.
	ReversedByteOrder
)
