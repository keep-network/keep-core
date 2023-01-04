package retransmission

import (
	"reflect"
	"testing"
)

func TestStandardStrategy(t *testing.T) {
	strategy := WithStandardStrategy()

	retransmitInvocations := make(map[int]bool)

	for i := 1; i <= 10; i++ {
		err := strategy.Tick(func() error {
			retransmitInvocations[i] = true
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	expectedRetransmitInvocations := map[int]bool{
		1:  true,
		2:  true,
		3:  true,
		4:  true,
		5:  true,
		6:  true,
		7:  true,
		8:  true,
		9:  true,
		10: true,
	}
	if !reflect.DeepEqual(expectedRetransmitInvocations, retransmitInvocations) {
		t.Errorf(
			"unexpected invocations\n"+
				"expected: [%v]\n"+
				"actual:   [%v]",
			expectedRetransmitInvocations,
			retransmitInvocations,
		)
	}
}

func TestBackoffStrategy(t *testing.T) {
	strategy := WithBackoffStrategy()

	retransmitInvocations := make(map[int]bool)

	for i := 1; i <= 100; i++ {
		err := strategy.Tick(func() error {
			retransmitInvocations[i] = true
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	expectedRetransmitInvocations := map[int]bool{
		1:  true,
		3:  true,
		6:  true,
		11: true,
		20: true,
		37: true,
		70: true,
	}
	if !reflect.DeepEqual(expectedRetransmitInvocations, retransmitInvocations) {
		t.Errorf(
			"unexpected invocations\n"+
				"expected: [%v]\n"+
				"actual:   [%v]",
			expectedRetransmitInvocations,
			retransmitInvocations,
		)
	}
}
