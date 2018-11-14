package result

import (
	"bytes"
	"encoding/binary"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto/sha3"
)

type ResultType int

const (
	NoFaultFailure ResultType = iota
	FailureDQ
	PerfectSuccess
	SuccessIA
	SuccessDQ
	MixedSuccess
)

type Result struct {
	Type           ResultType `json:"resulttype"`
	GroupPublicKey *big.Int   `json:"pubkey"`
	HashValue      []byte     `json:"-"`
	Disqualified   []int      `json:"disqualified"`
	Inactive       []int      `json:"inactive"`
}

func (r *Result) Hash() []byte {
	// OLD: return []byte(fmt.Sprintf("%v", r))
	bb := SearilizeResultForHash(*r)
	hh := Keccak256(bb)
	return hh
}

func (r Result) Equal(s Result) bool {
	return bytes.Equal(r.HashValue, s.HashValue)
}

func SearilizeResultForHash(rs Result) []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, int64(rs.Type))

	binary.Write(&buf, binary.BigEndian, int64(len(rs.Disqualified)))
	for _, b := range rs.Disqualified {
		if b == 1 {
			buf.Write([]byte{'1'})
		} else {
			buf.Write([]byte{'0'})
		}
	}

	buf.Write([]byte(rs.GroupPublicKey.String()))

	binary.Write(&buf, binary.BigEndian, int64(len(rs.Inactive)))
	for _, b := range rs.Inactive {
		if b == 1 {
			buf.Write([]byte{'1'})
		} else {
			buf.Write([]byte{'0'})
		}
	}

	return buf.Bytes()
}

// Keccak256 use the Ethereum Keccak hashing fucntions to return a hash from a list of values.
func Keccak256(data ...[]byte) []byte {
	d := sha3.NewKeccak256()
	for _, b := range data {
		d.Write(b)
	}
	return d.Sum(nil)
}
