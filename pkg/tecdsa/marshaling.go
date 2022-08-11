package tecdsa

import (
	"fmt"
	"github.com/binance-chain/tss-lib/crypto"
	"github.com/binance-chain/tss-lib/crypto/paillier"
	"github.com/binance-chain/tss-lib/ecdsa/keygen"
	"github.com/keep-network/keep-core/pkg/tecdsa/gen/pb"
	"math/big"
)

// ErrIncompatiblePublicKey indicates that the given public key is not
// compatible with the tECDSA, i.e. uses a different elliptic curve.
// Such a key cannot be marshalled and unmarshalled.
var ErrIncompatiblePublicKey = fmt.Errorf(
	"public key is not tECDSA compatible and will cause unmarshaling error",
)

// Marshal converts the PrivateKeyShare to a byte array.
func (pks *PrivateKeyShare) Marshal() ([]byte, error) {
	if pks.PublicKey().Curve.Params().Name != Curve.Params().Name {
		return nil, ErrIncompatiblePublicKey
	}

	localPreParams := &pb.LocalPartySaveData_LocalPreParams{
		PaillierSK: &pb.LocalPartySaveData_LocalPreParams_PrivateKey{
			PublicKey: pks.data.LocalPreParams.PaillierSK.PublicKey.N.Bytes(),
			LambdaN:   pks.data.LocalPreParams.PaillierSK.LambdaN.Bytes(),
			PhiN:      pks.data.LocalPreParams.PaillierSK.PhiN.Bytes(),
		},
		NTilde: pks.data.LocalPreParams.NTildei.Bytes(),
		H1I:    pks.data.LocalPreParams.H1i.Bytes(),
		H2I:    pks.data.LocalPreParams.H2i.Bytes(),
		Alpha:  pks.data.LocalPreParams.Alpha.Bytes(),
		Beta:   pks.data.LocalPreParams.Beta.Bytes(),
		P:      pks.data.LocalPreParams.P.Bytes(),
		Q:      pks.data.LocalPreParams.Q.Bytes(),
	}

	localSecrets := &pb.LocalPartySaveData_LocalSecrets{
		Xi:      pks.data.LocalSecrets.Xi.Bytes(),
		ShareID: pks.data.LocalSecrets.ShareID.Bytes(),
	}

	marshalBigIntSlice := func(bigInts []*big.Int) [][]byte {
		bytesSlice := make([][]byte, len(bigInts))
		for i, bigInt := range bigInts {
			bytesSlice[i] = bigInt.Bytes()
		}
		return bytesSlice
	}

	bigXj := make([]*pb.LocalPartySaveData_ECPoint, len(pks.data.BigXj))
	for i, bigX := range pks.data.BigXj {
		bigXj[i] = &pb.LocalPartySaveData_ECPoint{
			X: bigX.X().Bytes(),
			Y: bigX.Y().Bytes(),
		}
	}

	paillierPKs := make([][]byte, len(pks.data.PaillierPKs))
	for i, paillierPK := range pks.data.PaillierPKs {
		paillierPKs[i] = paillierPK.N.Bytes()
	}

	ecdsaPub := &pb.LocalPartySaveData_ECPoint{
		X: pks.data.ECDSAPub.X().Bytes(),
		Y: pks.data.ECDSAPub.Y().Bytes(),
	}

	return (&pb.PrivateKeyShare{
		Data: &pb.LocalPartySaveData{
			LocalPreParams: localPreParams,
			LocalSecrets:   localSecrets,
			Ks:             marshalBigIntSlice(pks.data.Ks),
			NTildej:        marshalBigIntSlice(pks.data.NTildej),
			H1J:            marshalBigIntSlice(pks.data.H1j),
			H2J:            marshalBigIntSlice(pks.data.H2j),
			BigXj:          bigXj,
			PaillierPKs:    paillierPKs,
			EcdsaPub:       ecdsaPub,
		},
	}).Marshal()
}

// Unmarshal converts a byte array back to the PrivateKeyShare.
func (pks *PrivateKeyShare) Unmarshal(bytes []byte) error {
	pbPrivateKeyShare := pb.PrivateKeyShare{}
	if err := pbPrivateKeyShare.Unmarshal(bytes); err != nil {
		return fmt.Errorf("failed to unmarshal private key share: [%v]", err)
	}

	data := new(keygen.LocalPartySaveData)
	pbData := pbPrivateKeyShare.Data

	paillierSK := &paillier.PrivateKey{
		PublicKey: paillier.PublicKey{
			N: new(big.Int).SetBytes(pbData.GetLocalPreParams().GetPaillierSK().GetPublicKey()),
		},
		LambdaN: new(big.Int).SetBytes(pbData.GetLocalPreParams().GetPaillierSK().GetLambdaN()),
		PhiN:    new(big.Int).SetBytes(pbData.GetLocalPreParams().GetPaillierSK().GetPhiN()),
	}

	data.LocalPreParams = keygen.LocalPreParams{
		PaillierSK: paillierSK,
		NTildei:    new(big.Int).SetBytes(pbData.GetLocalPreParams().GetNTilde()),
		H1i:        new(big.Int).SetBytes(pbData.GetLocalPreParams().GetH1I()),
		H2i:        new(big.Int).SetBytes(pbData.GetLocalPreParams().GetH2I()),
		Alpha:      new(big.Int).SetBytes(pbData.GetLocalPreParams().GetAlpha()),
		Beta:       new(big.Int).SetBytes(pbData.GetLocalPreParams().GetBeta()),
		P:          new(big.Int).SetBytes(pbData.GetLocalPreParams().GetP()),
		Q:          new(big.Int).SetBytes(pbData.GetLocalPreParams().GetQ()),
	}

	data.LocalSecrets = keygen.LocalSecrets{
		Xi:      new(big.Int).SetBytes(pbData.GetLocalSecrets().GetXi()),
		ShareID: new(big.Int).SetBytes(pbData.GetLocalSecrets().GetShareID()),
	}

	unmarshalBigIntSlice := func(bytesSlice [][]byte) []*big.Int {
		bigIntSlice := make([]*big.Int, len(bytesSlice))
		for i, bytes := range bytesSlice {
			bigIntSlice[i] = new(big.Int).SetBytes(bytes)
		}
		return bigIntSlice
	}

	data.BigXj = make([]*crypto.ECPoint, len(pbData.GetBigXj()))
	for i, bigX := range pbData.GetBigXj() {
		decoded, err := crypto.NewECPoint(
			Curve,
			new(big.Int).SetBytes(bigX.X),
			new(big.Int).SetBytes(bigX.Y),
		)
		if err != nil {
			return fmt.Errorf("failed to decode BigXj: [%v]", err)
		}

		data.BigXj[i] = decoded
	}

	data.PaillierPKs = make([]*paillier.PublicKey, len(pbData.GetPaillierPKs()))
	for i, paillierPK := range pbData.GetPaillierPKs() {
		data.PaillierPKs[i] = &paillier.PublicKey{
			N: new(big.Int).SetBytes(paillierPK),
		}
	}

	decoded, err := crypto.NewECPoint(
		Curve,
		new(big.Int).SetBytes(pbData.GetEcdsaPub().GetX()),
		new(big.Int).SetBytes(pbData.GetEcdsaPub().GetY()),
	)
	if err != nil {
		return fmt.Errorf("failed to decode ECDSAPub: [%v]", err)
	}
	data.ECDSAPub = decoded

	data.Ks = unmarshalBigIntSlice(pbData.GetKs())
	data.NTildej = unmarshalBigIntSlice(pbData.GetNTildej())
	data.H1j = unmarshalBigIntSlice(pbData.GetH1J())
	data.H2j = unmarshalBigIntSlice(pbData.GetH2J())

	pks.data = *data

	return nil
}
