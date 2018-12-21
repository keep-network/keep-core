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

func TestDoNotUnsubscribeTwice(t *testing.T) {
	unsubscribeCount := 0

	subscription := NewSubscription(func() {
		unsubscribeCount = unsubscribeCount + 1
	})

	subscription.Unsubscribe()
	subscription.Unsubscribe()
	subscription.Unsubscribe()

	if unsubscribeCount != 1 {
		t.Fatalf("unsubscribe handler should be called only once")
	}
}
