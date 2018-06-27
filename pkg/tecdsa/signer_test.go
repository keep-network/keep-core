package tecdsa

import (
	"testing"

	"github.com/ethereum/go-ethereum/crypto/secp256k1"
)

var publicParameters = &PublicParameters{
	groupSize: 10,
	threshold: 6,
	curve:     secp256k1.S256(),
}

func TestLocalSignerGenerateDsaKeyShare(t *testing.T) {
	group, err := newGroup(publicParameters)
	if err != nil {
		t.Fatal(err)
	}

	signer := group[0]

	dsaKeyShare, err := signer.generateDsaKeyShare()
	if err != nil {
		t.Fatal(err)
	}

	// TODO: probably it's not a cardinality... need to fix it
	curveCardinality := publicParameters.curve.Params().N

	if curveCardinality.Cmp(dsaKeyShare.xi) != 1 {
		t.Errorf("xi DSA key share must be less than Curve's cardinality")
	}

	if curveCardinality.Cmp(dsaKeyShare.xi) != 1 {
		t.Errorf("y.x DSA key share must be less than Curve's cardinality")
	}

	if curveCardinality.Cmp(dsaKeyShare.xi) != 1 {
		t.Errorf("y.y DSA key share must be less than Curve's cardinality")
	}
}
