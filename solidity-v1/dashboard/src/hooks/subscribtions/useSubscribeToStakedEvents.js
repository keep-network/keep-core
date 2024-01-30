import { useDispatch, useSelector } from "react-redux"
import { getEventsFromTransaction } from "../../utils/ethereum.utils"
import { isSameEthAddress } from "../../utils/general.utils"
import moment from "moment"
import {
  Keep,
  KeepExplorerMode,
  createManagedGrantContractInstance,
} from "../../contracts"
import { useSubscribeToExplorerModeContractEvent } from "../useSubscribeToExplorerModeContractEvent"
import { add } from "../../utils/arithmetics.utils"
import { ADD_STAKE_TO_THRESHOLD_AUTH_DATA } from "../../actions"
import {
  STAKING_PORT_BACKER_CONTRACT_NAME,
  TOKEN_GRANT_CONTRACT_NAME,
  TOKEN_STAKING_CONTRACT_NAME,
  TOKEN_STAKING_ESCROW_CONTRACT_NAME,
} from "../../lib/keep/contracts"

export const useSubscribeToStakedEvents = () => {
  const { initializationPeriod } = useSelector((state) => state.staking)
  const {
    [TOKEN_STAKING_CONTRACT_NAME]: { instance: stakingContract },
    [TOKEN_GRANT_CONTRACT_NAME]: { instance: grantContract },
    [TOKEN_STAKING_ESCROW_CONTRACT_NAME]: { instance: tokenStakingEscrow },
    [STAKING_PORT_BACKER_CONTRACT_NAME]: {
      instance: stakingPortBackerContract,
    },
    web3,
  } = KeepExplorerMode
  const dispatch = useDispatch()
  const yourAddress = web3?.lib?.eth?.defaultAccount

  useSubscribeToExplorerModeContractEvent(
    TOKEN_STAKING_CONTRACT_NAME,
    "StakeDelegated",
    async (event) => {
      try {
        const eventsToCheck = [
          [stakingContract, "OperatorStaked"],
          [grantContract, "TokenGrantStaked"],
          [tokenStakingEscrow, "DepositRedelegated"],
          [stakingPortBackerContract, "StakeCopied"],
        ]

        const {
          transactionHash,
          returnValues: { owner, operator },
        } = event

        const emittedEvents = await getEventsFromTransaction(
          eventsToCheck,
          transactionHash
        )

        let isAddressedToCurrentAccount = isSameEthAddress(owner, yourAddress)
        // The `OperatorStaked` is always emitted with the `StakeDelegated` event.
        const { authorizer, beneficiary, value } = emittedEvents.OperatorStaked
        const delegation = {
          createdAt: moment().unix(),
          operatorAddress: operator,
          authorizerAddress: authorizer,
          beneficiary,
          amount: value,
          isInInitializationPeriod: true,
          initializationOverAt: moment
            .unix(moment().unix())
            .add(initializationPeriod, "seconds"),
        }

        if (emittedEvents.StakeCopied) {
          const { owner } = emittedEvents.StakeCopied
          delegation.isCopiedStake = true
          isAddressedToCurrentAccount = isSameEthAddress(owner, yourAddress)

          // Check if the copied delegation is from grant.
          if (isAddressedToCurrentAccount) {
            try {
              const { grantId } = await grantContract.methods
                .getGrantStakeDetails(operator)
                .call()

              delegation.isFromGrant = true
              delegation.grantId = grantId
            } catch (error) {
              delegation.isFromGrant = false
            }
          }
        }

        if (
          (emittedEvents.TokenGrantStaked ||
            emittedEvents.DepositRedelegated) &&
          !isAddressedToCurrentAccount
        ) {
          // If the `TokenGrantStaked` or `DepositRedelegated` event exists, it means that a delegation is from grant.
          const { grantId } =
            emittedEvents.TokenGrantStaked || emittedEvents.DepositRedelegated
          delegation.grantId = grantId
          delegation.isFromGrant = true
          const { grantee } = await grantContract.methods
            .getGrant(grantId)
            .call()

          isAddressedToCurrentAccount = isSameEthAddress(grantee, yourAddress)

          if (!isAddressedToCurrentAccount) {
            // check if current address is a grantee in the managed grant
            try {
              const managedGrantContractInstance =
                createManagedGrantContractInstance(web3, grantee)
              const granteeAddressInManagedGrant =
                await managedGrantContractInstance.methods.grantee().call()
              delegation.managedGrantContractInstance =
                managedGrantContractInstance
              delegation.isManagedGrant = true

              // compere a current address with a grantee address from the ManagedGrant contract
              isAddressedToCurrentAccount = isSameEthAddress(
                yourAddress,
                granteeAddressInManagedGrant
              )
            } catch (error) {
              isAddressedToCurrentAccount = false
            }
          }
        }

        if (!isAddressedToCurrentAccount) {
          return
        }

        if (!delegation.isCopiedStake) {
          if (!delegation.isFromGrant) {
            dispatch({
              type: "staking/update_owned_delegated_tokens_balance",
              payload: { operation: add, value },
            })
          } else {
            dispatch({
              type: "token-grant/grant_staked",
              payload: {
                grantId: delegation.grantId,
                value,
              },
            })
          }
        }

        dispatch({ type: "staking/add_delegation", payload: delegation })
        if (isSameEthAddress(yourAddress, authorizer)) {
          dispatch({
            type: ADD_STAKE_TO_THRESHOLD_AUTH_DATA,
            payload: {
              ...delegation,
              owner: yourAddress,
              operatorContractAddress: Keep.thresholdStakingContract.address,
            },
          })
        }
      } catch (error) {
        console.log("error broo")
        console.error(
          `Failed subscribing to StakeDelegated event in Explorer Mode contract`,
          error
        )
      }
    }
  )
}
