package tecdsa

import (
	"github.com/keep-network/keep-core/pkg/internal/pbutils"
	"github.com/keep-network/keep-core/pkg/internal/tecdsatest"
	"reflect"
	"testing"
)

func TestPrivateKeyShareMarshalling(t *testing.T) {
	testData, err := tecdsatest.LoadPrivateKeyShareTestFixtures(1)
	if err != nil {
		t.Fatalf("failed to load test data: [%v]", err)
	}

	privateKeyShare := NewPrivateKeyShare(testData[0])

	unmarshaled := &PrivateKeyShare{}

	if err := pbutils.RoundTrip(privateKeyShare, unmarshaled); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(privateKeyShare, unmarshaled) {
		t.Fatal("unexpected content of unmarshaled private key share")
	}
}
