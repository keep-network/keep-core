package subscription

import (
	"context"
	"math/big"
	"sync"

	"github.com/keep-network/keep-core/pkg/beacon/relay/event"
)

// The guarantee we want:
// - IF an unsubscribe occurs
//   BEFORE handlers are called,
//   NO handlers should run.   ---> lock for unsubscribed flag
// - IF an unsubscribe occurs
//   AFTER handlers are called or
//   WHILE handlers are being called,
//   ALL handlers should finish running before we call
//   ANY unsubscribe handling. ---> WaitGroup for handlers

// Subscriber is a central location for:
//  - Event handler functions (+ associated context?).
//  - Event channels (+ associated context?).
//  - Managing connection to the server.

// Remaining: how do we handle contexts?

// DKGResultSubmissionSubscriber creates subscriptions to DKGResultSubmission
// events, either in an event-handler style or in a channel style. The
// subscriptions in turn can be managed as EventSubscriptions, allowing for
// unsubscribes. Subscriptions can also be cancelled via context if one is
// provided.
type DKGResultSubmissionSubscriber interface {
	OnEvent(dkgResultSubmissionHandlerFunc) EventSubscription
	OnEventContext(dkgResultSubmissionHandlerFunc, *context.Context) EventSubscription
	Pipe(chan<- event.DKGResultSubmission) EventSubscription
	PipeContext(chan<- event.DKGResultSubmission, *context.Context) EventSubscription
}

type dkgResultSubmissionHandlerFunc func(
	requestID *big.Int,
	memberIndex uint32,
	groupPublicKey []byte,
	blockNumber uint64,
)

type dkgResultSubmissionHandler struct {
	callback dkgResultSubmissionHandlerFunc
	context  *context.Context
}

type dkgResultSubmissionChannels struct {
	events  chan<- event.DKGResultSubmission
	errors  chan<- error
	context *context.Context
}

type dkgResultSubmissionSubscriber struct {
	context    context.Context
	cancelFunc context.CancelFunc

	subscriptionMutex sync.Mutex     // guards subscription management
	handlingWaitGroup sync.WaitGroup // guards handler execution/unsubscribe handler
	unsubscribed      bool

	subscriptionID   int
	callbackHandlers map[int]dkgResultSubmissionHandler
	channelHandlers  map[int]dkgResultSubmissionChannels
}

// NewDKGResultSubmissionSubscriber does some stuff.
func NewDKGResultSubmissionSubscriber(c context.Context, events <-chan *event.DKGResultSubmission, errors <-chan error) DKGResultSubmissionSubscriber {
	subscriberContext, subscriberCancel := context.WithCancel(c)
	subscriber := &dkgResultSubmissionSubscriber{
		context:    subscriberContext,
		cancelFunc: subscriberCancel,
	}

	// FIXME Should be go subscriber.eventLoop() or similar.
	go (func() {
		for {
			select {
			case event, closed := <-events:
				// closed == "subscribed", in that when the events channel is
				// closed it means we've been unsubscribed from the event
				// stream.
				subscriber.handleEvent(*event, closed)

			case err := <-errors:
				subscriber.handleFailure(err)
				return

			case <-subscriberContext.Done():
				if err := subscriberContext.Err(); err != nil {
					subscriber.handleFailure(err)
				} else {
					subscriber.handleUnsubscribe()
				}

				return
			}
		}
	})()

	// Wire up the go-ethereum watcher, which supports receiving an additional
	// context parameter to manage lifecycle.
	//
	// return ec.keepGroupContract.WatchDkgResultPublishedEvent(
	// 	func(requestID *big.Int, groupPubKey []byte, blockNumber uint64) {
	// 		handler(&event.DKGResultSubmission{
	// 			RequestID:      requestID,
	// 			GroupPublicKey: groupPubKey,
	// 			BlockNumber:    blockNumber,
	// 		})
	// 	},
	// 	func(err error) error {
	// 		return fmt.Errorf(
	// 			"watch DKG result published failed with [%v]",
	// 			err,
	// 		)
	// 	},
	// )

	return &dkgResultSubmissionSubscriber{}
}

func (drss *dkgResultSubmissionSubscriber) handleEvent(event event.DKGResultSubmission, subscribed bool) {
	drss.subscriptionMutex.Lock()
	defer drss.subscriptionMutex.Unlock()

	if !subscribed {
		drss.unsubscribed = true
		go drss.runUnsubscribe() // waits for any running handlers, then closes out
		drss.subscriptionMutex.Unlock()
		return // stop watching for events
	}

	// Run handlers + manages waitgroup.
	drss.runSuccessHandlers(event)
}

func (drss *dkgResultSubmissionSubscriber) handleFailure(err error) {
	drss.subscriptionMutex.Lock()
	defer drss.subscriptionMutex.Unlock()

	// Comm out error.
	drss.runFailureHandlers(err)

	// force unsubscribe
	drss.runUnsubscribe()
}

func (drss *dkgResultSubmissionSubscriber) handleUnsubscribe() {
	drss.subscriptionMutex.Lock()
	defer drss.subscriptionMutex.Unlock()

	// force unsubscribe
	drss.runUnsubscribe()
}

func (drss *dkgResultSubmissionSubscriber) runUnsubscribe() {
	if !drss.unsubscribed {
		return
	}

	drss.handlingWaitGroup.Wait() // Wait for all handlers to complete.

	// Closes out the go-ethereum subscription.
	drss.cancelFunc()

	// Have optional cancelFunc for each channel handler container and call it?
	// This makes more sense than taking a context, which we can only observe,
	// not cancel.

	// Close event and error channels.
	for _, channelHandler := range drss.channelHandlers {
		// FIXME what to do with the channel-related context? Cancel requires cancel func.
		close(channelHandler.errors)
		close(channelHandler.events)
	}
}

func (drss *dkgResultSubmissionSubscriber) runSuccessHandlers(e event.DKGResultSubmission) {
	if drss.unsubscribed {
		return
	}

	for _, handler := range drss.callbackHandlers {
		go (func(waitGroup *sync.WaitGroup, handler dkgResultSubmissionHandlerFunc, e event.DKGResultSubmission) {
			waitGroup.Add(1)
			defer waitGroup.Done()

			handler(
				e.RequestID,
				e.MemberIndex,
				e.GroupPublicKey,
				e.BlockNumber,
			)
		})(&drss.handlingWaitGroup, handler.callback, e)
	}
}

func (drss *dkgResultSubmissionSubscriber) runFailureHandlers(err error) {
	if drss.unsubscribed {
		return
	}

	// No allowing for failure handlers at the moment.
}

func (drss *dkgResultSubmissionSubscriber) OnEvent(handler dkgResultSubmissionHandlerFunc) EventSubscription {
	return drss.OnEventContext(handler, nil)
}

func (drss *dkgResultSubmissionSubscriber) OnEventContext(handler dkgResultSubmissionHandlerFunc, context *context.Context) EventSubscription {
	drss.subscriptionMutex.Lock()
	defer drss.subscriptionMutex.Unlock()

	subscriptionID := drss.subscriptionID
	drss.callbackHandlers[subscriptionID] = dkgResultSubmissionHandler{handler, context}
	eventSubscription := NewEventSubscription(func() {
		drss.subscriptionMutex.Lock()
		defer drss.subscriptionMutex.Unlock()

		delete(drss.callbackHandlers, subscriptionID)
	})

	drss.subscriptionID++

	return eventSubscription
}

// Pipe ... Piped channels should not be closed until after the
// EventSubscription is cancelled, which will ensure the subscriber will not try
// to write to the channel.
func (drss *dkgResultSubmissionSubscriber) Pipe(channel chan<- event.DKGResultSubmission) EventSubscription {
	return drss.PipeContext(channel, nil)
}

func (drss *dkgResultSubmissionSubscriber) PipeContext(channel chan<- event.DKGResultSubmission, context *context.Context) EventSubscription {
	drss.subscriptionMutex.Lock()
	defer drss.subscriptionMutex.Unlock()

	subscriptionID := drss.subscriptionID
	drss.channelHandlers[subscriptionID] = dkgResultSubmissionChannels{
		events:  channel,
		errors:  nil,
		context: context,
	}
	eventSubscription := NewEventSubscription(func() {
		drss.subscriptionMutex.Lock()
		defer drss.subscriptionMutex.Unlock()

		delete(drss.channelHandlers, subscriptionID)
	})

	drss.subscriptionID++

	return eventSubscription
}
