package pedersen

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestGenerateAndValidateCommitment(t *testing.T) {
	vss, err := initializeVSS()
	if err != nil {
		t.Fatalf("VSS initialization failed error [%v]", err)
	}

	committedValue := "eeyore"

	var tests = map[string]struct {
		verificationValue     string
		modifyDecommitmentKey func(key *DecommitmentKey)
		modifyCommitment      func(commitment *Commitment)
		modifyVSS             func(parameters *VSS)
		expectedResult        bool
	}{
		"positive validation": {
			verificationValue: committedValue,
			expectedResult:    true,
		},
		"negative validation - verification value is not the committed one": {
			verificationValue: "pooh",
			expectedResult:    false,
		},
		"negative validation - incorrect decommitment key `t`": {
			verificationValue: committedValue,
			modifyDecommitmentKey: func(key *DecommitmentKey) {
				key.t, _ = generateNewRandom(key.t, vss.q)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `commitment`": {
			verificationValue: committedValue,
			modifyCommitment: func(commitment *Commitment) {
				commitment.commitment, _ = generateNewRandom(commitment.commitment, vss.q)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `g`": {
			verificationValue: committedValue,
			modifyVSS: func(vss *VSS) {
				vss.G, _ = generateNewRandom(vss.G, vss.q)
			},
			expectedResult: false,
		},
		"negative validation - incorrect `h`": {
			verificationValue: committedValue,
			modifyVSS: func(vss *VSS) {
				vss.h, _ = generateNewRandom(vss.h, vss.q)
			},
			expectedResult: false,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			commitment, decommitmentKey, err := vss.CommitmentTo(crand.Reader, []byte(committedValue))
			if err != nil {
				t.Fatalf("generation error [%v]", err)
			}

			if test.modifyCommitment != nil {
				test.modifyCommitment(commitment)
			}

			if test.modifyDecommitmentKey != nil {
				test.modifyDecommitmentKey(decommitmentKey)
			}

			if test.modifyVSS != nil {
				test.modifyVSS(vss)
			}

			result := commitment.Verify(decommitmentKey, []byte(test.verificationValue))

			if result != test.expectedResult {
				t.Fatalf(
					"expected: %v\nactual: %v\n",
					test.expectedResult,
					result,
				)
			}
		})
	}
}

func TestNewVSSpqValidation(t *testing.T) {
	var tests = map[string]struct {
		p             *big.Int
		q             *big.Int
		expectedError error
	}{
		"positive validation": {
			p:             big.NewInt(11),
			q:             big.NewInt(5),
			expectedError: nil,
		},
		"negative validation - p not prime": {
			p:             big.NewInt(16),
			q:             big.NewInt(7),
			expectedError: fmt.Errorf("p and q have to be primes"),
		},
		"negative validation - q not prime": {
			p:             big.NewInt(17),
			q:             big.NewInt(4),
			expectedError: fmt.Errorf("p and q have to be primes"),
		},
		"negative validation - incorrect p and q values": {
			p:             big.NewInt(19),
			q:             big.NewInt(3),
			expectedError: fmt.Errorf("incorrect p and q values"),
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			_, err := GenerateVSS(crand.Reader, test.p, test.q)

			if !reflect.DeepEqual(err, test.expectedError) {
				t.Fatalf("actual error: %v\nexpected error: %v", err, test.expectedError)
			}
		})
	}
}

func initializeVSS() (*VSS, error) {
	// Sets p and q to predefined fixed values, such that `p = 2q + 1`.
	// `p` is 4096-bit safe prime.
	pStr := "0xc8526644a9c4739683742b7003640b2023ca42cc018a42b02a551bb825c6828f86e2e216ea5d31004c433582a3fa720459efb42e091d73fb281810e1825691f0799811be62ae57f62ab00670edd35426d108d3b9c4fd008eddc67275a0489fe132e4c31bd7069ea7884cbb8f8f9255fe7b87fc0099f246776c340912df48f7945bc2bc0bc6814978d27b7af2ebc41f458ae795186db0fd7e6151bb8a7fe2b41370f7a2848ef75d3ec88f3439022c10e78b434c2f24b2f40bd02930e6c8aadef87b0dc87cdba07dcfa86884a168bd1381a4f48be12e5d98e41f954c37aec011cc683570e8890418756ed98ace8c8e59ae1df50962c1622fe66b5409f330cad6b7c68f2e884786d9807190b89ac4a3b3507e49b2dd3f33d765ad29e2015180c8cd0258dd8bdaab17be5d74871fec04c492240c6a2692b2c9a62c9adbaac34a333f135801ff948e8dfb6bbd6212a67950fb8edd628d05d19d1b94e9be7c52ed484831d50adaa29e71de197e351878f1c40ec67ee809e824124529e27bd5ecf3054f6784153f7db27ff0c87420bb2b2754ed363fc2ba8399d49d291f342173e7619183467a9694efa243e1d41b26c13b38ca0f43bb7c9050eb966461f28436583a9d13d2c1465b78184eae360f009505ccea288a053d111988d55c12befd882a857a530efac2c0592987cd83c39844a10e058739ab1c39006a3123e7fc887845675f"
	// `q` is 4095-bit Sophie Germain prime.
	qStr := "0x6429332254e239cb41ba15b801b2059011e5216600c52158152a8ddc12e34147c371710b752e988026219ac151fd39022cf7da17048eb9fd940c0870c12b48f83ccc08df31572bfb1558033876e9aa13688469dce27e80476ee3393ad0244ff09972618deb834f53c4265dc7c7c92aff3dc3fe004cf9233bb61a04896fa47bca2de15e05e340a4bc693dbd7975e20fa2c573ca8c36d87ebf30a8ddc53ff15a09b87bd142477bae9f64479a1c81160873c5a1a61792597a05e814987364556f7c3d86e43e6dd03ee7d4344250b45e89c0d27a45f0972ecc720fcaa61bd76008e6341ab87444820c3ab76cc56746472cd70efa84b160b117f335aa04f998656b5be347974423c36cc038c85c4d6251d9a83f24d96e9f99ebb2d694f100a8c06466812c6ec5ed558bdf2eba438ff602624912063513495964d3164d6dd561a5199f89ac00ffca4746fdb5deb109533ca87dc76eb14682e8ce8dca74df3e2976a42418ea856d514f38ef0cbf1a8c3c78e207633f7404f412092294f13deaf67982a7b3c20a9fbed93ff8643a105d9593aa769b1fe15d41ccea4e948f9a10b9f3b0c8c1a33d4b4a77d121f0ea0d93609d9c6507a1ddbe482875cb3230f9421b2c1d4e89e960a32dbc0c27571b07804a82e6751445029e888cc46aae095f7ec41542bd29877d61602c94c3e6c1e1cc22508702c39cd58e1c80351891f3fe443c22b3af"

	var result bool

	p, result := new(big.Int).SetString(pStr, 0)
	if !result {
		return nil, fmt.Errorf("failed to initialize p")
	}

	q, result := new(big.Int).SetString(qStr, 0)
	if !result {
		return nil, fmt.Errorf("failed to initialize q")
	}

	vss, err := GenerateVSS(crand.Reader, p, q)
	if err != nil {
		return nil, fmt.Errorf("vss creation failed [%v]", err)
	}
	return vss, nil
}

// generateNewRandom generates value different than `currentValue` in range (1, max)
func generateNewRandom(currentValue *big.Int, max *big.Int) (*big.Int, error) {
	for {
		x, err := crand.Int(crand.Reader, max)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random number [%v]", err)
		}
		if x.Cmp(currentValue) != 0 && x.Cmp(big.NewInt(1)) > 0 {
			return x, nil
		}
	}
}
