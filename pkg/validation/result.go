package validaiton

import (
	"bytes"
	"encoding/binary"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

/* Haskel Data Type

data Result = NoFaultFailure
            | FailureDQ { disqualified :: Array Bool }
            | PerfectSuccess { pubkey :: BlsPubkey }
            | SuccessIA { pubkey :: BlsPubkey
                        , inactive :: Array Bool }
            | SuccessDQ { pubkey :: BlsPubkey
                        , disqualified :: Array Bool }
            | MixedSuccess { pubkey :: BlsPubkey
                           , inactive :: Array Bool
                           , disqualified :: Array Bool }

class Result(NamedTuple):
    groupPubkey:  Optional[BlsPubkey]
    disqualified: Optional[BitArray]
    inactive:     Optional[BitArray]
*/

type ResultType int

const (
	NoFaultFailure ResultType = iota
	FailureDQ
	PerfectSuccess
	SuccessIA
	SuccessDQ
	MixedSuccess
)

// BlsPubkey is refered to as BeaconPubkey in the paper/deisng specification.
// This is actually A BLS key.  In the code it is called BlsPubkey.
type BlsPubkey struct {
}

type Result struct {
	TypeIs       ResultType
	Disqualified []bool
	Pubkey       BlsPubkey
	Inactive     []bool
}

func SearilizeResultForHash(rs Result) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, int64(rs.TypeIs))

	binary.Write(&buf, binary.BigEndian, int64(len(rs.Disqualified)))
	for _, b := range rs.Disqualified {
		if b {
			buf.Write([]byte{'1'})
		} else {
			buf.Write([]byte{'0'})
		}
	}

	// FIXME - TODO - put in Pubkey

	binary.Write(&buf, binary.BigEndian, int64(len(rs.Inactive)))
	for _, b := range rs.Inactive {
		if b {
			buf.Write([]byte{'1'})
		} else {
			buf.Write([]byte{'0'})
		}
	}

	return buf.Bytes()
}

func Hash(data ...[]byte) []byte {
	// this probably isn't the correct hash - but will work for testing at the moment.
	return Keccak256(data...)
}

// Keccak256 use the Ethereum Keccak hasing fucntions to return a hash from a list of values.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
