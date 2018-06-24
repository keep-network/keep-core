package chaintype

import "math/big"

// GroupPublicKey represents the data from the
// KeepRandomBeaconImplV1SubmitGroupPublicKeyEvent found in
// ../chain/gen/KeepRandomBeaconImplV1.go  This type allows for the
// conversion of tyeps and keeps the difference betwen this event
// (on the implementaiton cotract) and the matching event on the
// proxy contract from propogatting to the upper layers of the Go code.
type GroupPublicKey struct {
	GroupPublicKey        []byte
	RequestID             *big.Int
	ActivationBlockHeight *big.Int
}
