package event

import "testing"

func TestUnsubscribe(t *testing.T) {
	unsubscribed := false

	subscription := NewSubscription(func() {
		unsubscribed = true
	})

	if unsubscribed {
		t.Fatalf("should not be unsubscribed at this point")
	}

	subscription.Unsubscribe()

	if !unsubscribed {
		t.Fatalf("should be unsubscribed at this point")
	}
}
