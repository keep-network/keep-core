package gjkr

import (
	"fmt"
	"math"
)

//func (pm *PublishingMember) PublishResult(result *result.Result, t_dkg int) (*event.PublishedResult, error) {
//	chainRelay := pm.protocolConfig.chain.ThresholdRelay()
//	// while not resultPublished:
//	for !chainRelay.IsResultPublished(result) {
//		publisherID, err := pm.determinePublisherID() // j
//		if err != nil {
//			return nil, err
//		}
//		//   if j >= i:
//		if publisherID >= pm.ID {
//			errors := make(chan error)
//			eventPublish := make(chan *event.PublishedResult)
//			// broadcast(correctResult)
//			chainRelay.SubmitResult(
//				pm.ID,
//				result,
//			).OnComplete(func(publish *event.PublishedResult, err error) {
//				eventPublish <- publish
//				errors <- err
//			})
//			return <-eventPublish, <-errors
//		}
//	}
//	return nil, nil
//}

func (pm *PublishingMember) determinePublisherID() (int, error) {
	t_dkg := pm.protocolConfig.expectedProtocolDuration // t_dkg
	t_step := pm.protocolConfig.blockStep

	blockCounter, err := pm.protocolConfig.chain.BlockCounter()
	if err != nil {
		return 0, err
	}
	//   T_now = getCurrentBlockHeight()
	t_now, err := blockCounter.CurrentBlock()
	if err != nil {
		return 0, err
	}

	// # using T_init from phase 1
	t_init := pm.protocolConfig.initialBlockHeight
	//   T_elapsed = T_now - T_init
	t_elapsed := t_now - t_init

	// # determine highest index j eligible to submit
	// if T_elapsed <= T_dkg:
	var playerIndex int
	if t_elapsed <= t_dkg {
		// j = 1
		playerIndex = 0
		//   else:
	} else {
		//     T_over = T_elapsed - T_dkg
		t_over := t_elapsed - t_dkg
		//     j = 1 + ceiling(T_over / T_step)
		playerIndex = int(math.Ceil(float64(t_over / t_step)))
	}
	if playerIndex > pm.group.groupSize {
		panic(fmt.Errorf("player index %d out of group size", playerIndex))
	}
	j := pm.group.MemberIDs()[playerIndex]
	return j, nil
}
