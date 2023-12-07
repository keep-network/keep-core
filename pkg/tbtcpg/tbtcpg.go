package tbtcpg

import (
	"fmt"

	"github.com/ipfs/go-log/v2"
	"github.com/keep-network/keep-core/pkg/bitcoin"
	"github.com/keep-network/keep-core/pkg/tbtc"
	"go.uber.org/zap"
	"golang.org/x/exp/slices"
)

var logger = log.Logger("keep-tbtcpg")

// ProposalTask encapsulates logic used to generate an action proposal
// of the given type.
type ProposalTask interface {
	// Run executes the task and returns a proposal, a boolean flag indicating
	// whether the proposal was generated and an error if any.
	Run(
		request *tbtc.CoordinationProposalRequest,
	) (tbtc.CoordinationProposal, bool, error)

	// ActionType returns the type of the action proposal.
	ActionType() tbtc.WalletActionType
}

// ProposalGenerator is a component responsible for generating coordination
// proposals for tbtc wallets.
type ProposalGenerator struct {
	tasks []ProposalTask
}

// NewProposalGenerator returns a new proposal generator.
func NewProposalGenerator(
	chain Chain,
	btcChain bitcoin.Chain,
) *ProposalGenerator {
	tasks := []ProposalTask{
		NewDepositSweepTask(chain, btcChain),
		NewRedemptionTask(chain, btcChain),
		NewHeartbeatTask(chain),
		NewMovingFundsTask(chain, btcChain),
		// TODO: Uncomment when moving funds support is implemented.
		// newMovedFundsSweepTask(),
	}

	return &ProposalGenerator{
		tasks: tasks,
	}
}

// Generate generates a coordination proposal based on the given checklist
// of possible wallet actions. The checklist is a list of actions that
// should be checked for the given coordination window. This function returns
// a proposal for the first action from the checklist that is valid for the
// given wallet's state. If none of the actions are valid, the generator
// returns a no-op proposal.
func (pg *ProposalGenerator) Generate(
	request *tbtc.CoordinationProposalRequest,
) (tbtc.CoordinationProposal, error) {
	walletLogger := logger.With(
		zap.String(
			"walletPKH",
			fmt.Sprintf("0x%x", request.WalletPublicKeyHash),
		),
	)

	walletLogger.Info(
		"starting proposal generation with tasks checklist [%v]",
		request.ActionsChecklist,
	)

	for _, action := range request.ActionsChecklist {
		walletLogger.Infof("starting proposal task [%s]", action)

		taskIndex := slices.IndexFunc(pg.tasks, func(task ProposalTask) bool {
			return task.ActionType() == action
		})

		if taskIndex < 0 {
			walletLogger.Warnf("proposal task [%s] is not supported", action)
			continue
		}

		proposal, ok, err := pg.tasks[taskIndex].Run(request)
		if err != nil {
			return nil, fmt.Errorf(
				"error while running proposal task [%s]: [%v]",
				action,
				err,
			)
		}

		if !ok {
			walletLogger.Infof(
				"proposal task [%s] completed without result",
				action,
			)
			continue
		}

		walletLogger.Infof(
			"proposal task [%s] completed with a result",
			action,
		)

		return proposal, nil
	}

	walletLogger.Infof(
		"all proposal tasks completed without result; " +
			"returning no-op proposal",
	)

	return &tbtc.NoopProposal{}, nil
}
