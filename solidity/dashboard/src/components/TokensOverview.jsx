import React from 'react'
import TokenGrantOverview from './TokenGrantOverview'
import { useTokensPageContext } from '../contexts/TokensPageContext'
import OwnedTokensOverview from './OwnedTokensOverview.jsx'
import { add } from '../utils/arithmetics.utils'

const TokensOverview = (props) => {
  const {
    tokensContext,
    selectedGrant,
    keepTokenBalance,
    ownedTokensUndelegationsBalance,
    ownedTokensDelegationsBalance,
  } = useTokensPageContext()

  return tokensContext === 'granted' ?
    <TokenGrantOverview selectedGrant={selectedGrant} /> :
    <OwnedTokensOverview
      keepBalance={keepTokenBalance}
      stakedBalance={add(ownedTokensUndelegationsBalance, ownedTokensDelegationsBalance)}
    />
}

export default TokensOverview
