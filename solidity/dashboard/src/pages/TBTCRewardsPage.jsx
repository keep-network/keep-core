// import React, { useMemo, useCallback } from "react"
// import PageWrapper from "../components/PageWrapper"
// import TBTCRewardsDataTable from "../components/TBTCRewardsDataTable"
// import { tbtcRewardsService } from "../services/tbtc-rewards.service"
// import { useWeb3Context } from "../components/WithWeb3Context"
// import { useFetchData } from "../hooks/useFetchData"
// import { add } from "../utils/arithmetics.utils"
// import { toTokenUnit } from "../utils/token.utils"
// import { findIndexAndObject } from "../utils/array.utils"

// const TBTCRewardsPage = () => {
//   const web3Context = useWeb3Context()
//   const { yourAddress } = web3Context
//   const [{ data }, updateRewardsData] = useFetchData(
//     tbtcRewardsService.fetchTBTCRewards,
//     [],
//     yourAddress
//   )

//   const totalAmount = useMemo(() => {
//     return toTokenUnit(
//       data.map((reward) => reward.amount).reduce(add, 0)
//     ).toString()
//   }, [data])

//   const fetchOperatorByDepositId = useCallback(
//     async (depositId) => {
//       const operators = await tbtcRewardsService.fetchBeneficiaryOperatorsFromDeposit(
//         web3Context,
//         yourAddress,
//         depositId
//       )
//       const updatedData = [...data]
//       operators.forEach((operator) => {
//         const { indexInArray, obj } = findIndexAndObject(
//           "depositTokenId",
//           depositId,
//           updatedData
//         )

//         if (indexInArray !== null) {
//           updatedData[indexInArray] = { ...obj, operatorAddress: operator }
//         }
//       })

//       updateRewardsData(updatedData)
//     },
//     [data, updateRewardsData, yourAddress, web3Context]
//   )

//   return (
//     <PageWrapper title="tBTC Rewards">
//       <section className="tile">
//         <h2 className="text-grey-70">Total Amount</h2>
//         <h1 className="text-primary">
//           {totalAmount}&nbsp;<span className="h3">TBTC</span>
//         </h1>
//         <TBTCRewardsDataTable
//           rewards={data}
//           fetchOperatorByDepositId={fetchOperatorByDepositId}
//         />
//       </section>
//     </PageWrapper>
//   )
// }

// export default React.memo(TBTCRewardsPage)
