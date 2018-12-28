package subscription

import "sync"

// Subscription is returned by an event source as a result of operation
// subscribing to the given type of event. It allows to unsubscribe from
// the event stream at any point, by calling the `Unsubscribe` method.
type Subscription interface {
	Unsubscribe()
}

// NewSubscription is used by an event source to create a `Subscription`.
// It accepts a callback function that is called as a result of
// `Unsubscribe` operation on the `Subscription`. The callback function
// executes operation on the event source required to unsubscribe from the
// event stream.
func NewSubscription(doUnsubscribe func()) Subscription {
	return &subscription{
		mutex:             sync.Mutex{},
		unsubscribed:      false,
		doUnsubscribeFunc: doUnsubscribe,
	}
}

type subscription struct {
	mutex        sync.Mutex
	unsubscribed bool

	doUnsubscribeFunc func()
}

func (s *subscription) Unsubscribe() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.unsubscribed {
		return
	}

	s.doUnsubscribeFunc()
	s.unsubscribed = true
}
