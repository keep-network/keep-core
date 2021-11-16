import React, { useEffect } from "react"
import { MigrationPortalForm } from "../../components/tbtc-migration"
import { TBTC } from "../../utils/token.utils"
import { useModal } from "../../hooks/useModal"
import { useWeb3Address } from "../../components/WithWeb3Context"
import { useDispatch, useSelector } from "react-redux"
import { tbtcV2Migration } from "../../actions"
import { MODAL_TYPES } from "../../constants/constants"

const TokenUpgradePortalPage = () => {
  const { openModal } = useModal()
  const address = useWeb3Address()
  const dispatch = useDispatch()

  const { tbtcV1Balance, tbtcV2Balance, unmintFee } = useSelector(
    (state) => state.tbtcV2Migration
  )

  useEffect(() => {
    dispatch(tbtcV2Migration.fetchDataRequest(address))
  }, [address, dispatch])

  const onSubmitMigrationForm = async (values, awaitingPromise) => {
    const { amount, from, to } = values
    const _amount = TBTC.fromTokenUnit(amount).toString()

    openModal(MODAL_TYPES.ConfirmTBTCMigration, {
      from,
      to,
      amount: _amount,
      fee: unmintFee,
    })
  }

  return (
    <section className="tbtc-migration-portal">
      <section className="tbtc-migration-portal__form-wrapper">
        <h3 className="text-grey-70 mb-1">Migration Portal</h3>
        <MigrationPortalForm
          mintingFee={unmintFee}
          tbtcV1Balance={tbtcV1Balance}
          tbtcV2Balance={tbtcV2Balance}
          onSubmit={onSubmitMigrationForm}
        />
      </section>
    </section>
  )
}

TokenUpgradePortalPage.route = {
  title: "Token Upgrade Portal",
  path: "/tbtc-migration/portal",
  exact: true,
}

export default TokenUpgradePortalPage
