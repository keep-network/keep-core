import React from "react"
import CreateTokenGrantForm from "../components/CreateTokenGrantForm"
import PageWrapper from "../components/PageWrapper"
import Tile from "../components/Tile"
import { keepTokenApproveAndCall } from "../actions/web3"
import { ContractsLoaded } from "../contracts"
import { tokenGrantsService } from "../services/token-grants.service"
import { fromTokenUnit } from "../utils/token.utils"
import { connect, useSelector } from "react-redux"

const CreateTokenGrantPage = ({ createTokenGrant }) => {
  const keepTokenBalance = useSelector((state) => state.keepTokenBalance)

  const submitAction = async (values, meta) => {
    const { grantContract } = await ContractsLoaded
    const extraData = await tokenGrantsService.getCreateTokenGrantExtraData(
      values
    )
    const amount = fromTokenUnit(values.amount).toString()
    const tokenAddress = grantContract.options.address

    createTokenGrant(
      {
        amount,
        tokenAddress,
        extraData,
      },
      meta
    )
  }

  return (
    <PageWrapper title="Create Token Grant">
      <Tile title="Create Grant" className="rewards-history tile flex column">
        <CreateTokenGrantForm
          keepBalance={keepTokenBalance.value}
          submitAction={submitAction}
        />
      </Tile>
    </PageWrapper>
  )
}

const mapDispatchToProps = {
  createTokenGrant: keepTokenApproveAndCall,
}

export default connect(null, mapDispatchToProps)(CreateTokenGrantPage)
