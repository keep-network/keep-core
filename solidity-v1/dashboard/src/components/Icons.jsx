import React from "react"
import { colors } from "../constants/colors"
import { ReactComponent as ArrowRight } from "../static/svg/arrow-right.svg"
import { ReactComponent as Operations } from "../static/svg/operations.svg"
import { ReactComponent as Rewards } from "../static/svg/rewards.svg"
import { ReactComponent as Authorizer } from "../static/svg/authorizer.svg"
import { ReactComponent as KeepToken } from "../static/svg/tokens.svg"
import { ReactComponent as GrantContextIcon } from "../static/svg/grant-context-icon.svg"
import { ReactComponent as MoneyWalletOpen } from "../static/svg/money-wallet-open.svg"
import { ReactComponent as KeepGreenOutline } from "../static/svg/keep-green-form-icon.svg"
import { ReactComponent as AuthorizerFormIcon } from "../static/svg/authorizer-form-icon.svg"
import { ReactComponent as OperatorFormIcon } from "../static/svg/operator-form-icon.svg"
import { ReactComponent as BeneficiaryFormIcon } from "../static/svg/beneficiary-form-icon.svg"
import { ReactComponent as DashedLine } from "../static/svg/dashed-line.svg"
import { ReactComponent as KeepOutline } from "../static/svg/keep-green-outline.svg"
import { ReactComponent as LedgerDevice } from "../static/svg/ledger-device.svg"
import { ReactComponent as TrezorDevice } from "../static/svg/trezor-device.svg"
import { ReactComponent as TBTC } from "../static/svg/tbtc.svg"
import { ReactComponent as KeepBlackGreen } from "../static/svg/keep-token.svg"
import { ReactComponent as Filter } from "../static/svg/filter-icon.svg"
import { ReactComponent as Load } from "../static/svg/load.svg"
import { ReactComponent as DocumentWithBg } from "../static/svg/document-bg-green.svg"
import { ReactComponent as DelegationDiagram } from "../static/svg/delegation-diagram.svg"
import { ReactComponent as Diamond } from "../static/svg/diamond.svg"
import { ReactComponent as ETH } from "../static/svg/eth.svg"
import { ReactComponent as KeepLoadingIndicator } from "../static/svg/keep-loading-indicator.svg"
import { ReactComponent as KEEPTower } from "../static/svg/keep-tower.svg"
import { ReactComponent as CarDashboardSpeed } from "../static/svg/car-dashboard-speed.svg"
import { ReactComponent as Fees } from "../static/svg/fees.svg"
import { ReactComponent as UserFriendly } from "../static/svg/user-friendly.svg"
import { ReactComponent as Alert } from "../static/svg/alert.svg"
import { ReactComponent as Success } from "../static/svg/success.svg"
import { ReactComponent as Beacon } from "../static/svg/beacon.svg"
import { ReactComponent as Authorize } from "../static/svg/authorize.svg"
import { ReactComponent as Home } from "../static/svg/home.svg"
import { ReactComponent as Question } from "../static/svg/question.svg"
import { ReactComponent as FeesVector } from "../static/svg/fees-vector.svg"
import { ReactComponent as Time } from "../static/svg/time.svg"
import { ReactComponent as KeepDashboardLogo } from "../static/svg/token-dashboard-logo.svg"
import { ReactComponent as NetworkStatusIndicator } from "../static/svg/network-status-indicator.svg"
import { ReactComponent as MetaMask } from "../static/svg/metamask.svg"
import { ReactComponent as Tally } from "../static/svg/tally.svg"
import { ReactComponent as Trezor } from "../static/svg/trezor.svg"
import { ReactComponent as Ledger } from "../static/svg/ledger.svg"
import { ReactComponent as Add } from "../static/svg/add.svg"
import { ReactComponent as Subtract } from "../static/svg/subtract.svg"
import { ReactComponent as ArrowDown } from "../static/svg/arrow-down.svg"
export { ReactComponent as Warning } from "../static/svg/warning.svg"
export { ReactComponent as Wallet } from "../static/svg/wallet.svg"
export { ReactComponent as Grant } from "../static/svg/grant.svg"
export { ReactComponent as Calendar } from "../static/svg/calendar.svg"
export { ReactComponent as Plus } from "../static/svg/plus.svg"
export { ReactComponent as StakeDrop } from "../static/svg/stakedrop.svg"
export { ReactComponent as SwordOperations } from "../static/svg/sword-operations.svg"
export { ReactComponent as MoreInfo } from "../static/svg/more-info.svg"
export { ReactComponent as EthToken } from "../static/svg/eth_token.svg"
export { ReactComponent as KeepOnlyPool } from "../static/svg/keep-only-pool.svg"
export { ReactComponent as BalancerLogo } from "../static/svg/balancer-logo.svg"
export { ReactComponent as UniswapLogo } from "../static/svg/uniswap-logo.svg"
export { ReactComponent as WalletConnect } from "../static/svg/wallet-connect.svg"
export { ReactComponent as Explore } from "../static/svg/explore.svg"
export { ReactComponent as Dashboard } from "../static/svg/dashboard.svg"
export { ReactComponent as BrowserWindow } from "../static/svg/browser-window.svg"
export { ReactComponent as Decentralize } from "../static/svg/decentralize.svg"
export { ReactComponent as CoveragePool } from "../static/svg/coverage-pool.svg"
export { ReactComponent as SaddleWhite } from "../static/svg/saddle-logo.svg"
export { ReactComponent as Swap } from "../static/svg/swap.svg"
export { ReactComponent as ChevronRight } from "../static/svg/chevron-right.svg"
export { ReactComponent as ChevronUp } from "../static/svg/chevron-up.svg"
export { ReactComponent as ChevronDown } from "../static/svg/chevron-down.svg"
export { ReactComponent as Refresh } from "../static/svg/refresh.svg"
export { ReactComponent as CovPoolsHowItWorksDiagram } from "../static/svg/cov-pools-how-it-works-diagram.svg"
export { ReactComponent as TTokenSymbol } from "../static/svg/t-token-symbol.svg"
export { ReactComponent as KeepTUpgrade } from "../static/svg/keep-t-upgrade-logo.svg"
export { ReactComponent as Star } from "../static/svg/star.svg"
export { ReactComponent as Bell } from "../static/svg/bell.svg"
export { ReactComponent as MBTC } from "../static/svg/mBTC.svg"
export { ReactComponent as TBTC_V2 } from "../static/svg/tbtc_v2.svg"
export { ReactComponent as ArrowsRight } from "../static/svg/arrows-right.svg"
export { ReactComponent as TLogo } from "../static/svg/t-logo.svg"
export { ReactComponent as AlertFill } from "../static/svg/alert-fill.svg"
export { ReactComponent as ArrowTopRight } from "../static/svg/arrow-top-right.svg"
export { ReactComponent as QuestionFill } from "../static/svg/question-fill.svg"
export { ReactComponent as EarnThresholdTokens } from "../static/svg/earn-threshold-tokens.svg"
export { ReactComponent as Money } from "../static/svg/money.svg"
export { ReactComponent as CovKeep } from "../static/svg/keep-symbol.svg"
export { ReactComponent as LegacyDappIllustration } from "../static/svg/legacy-dapp-modal.svg"

const Badge = ({ height, width }) => (
  <svg
    width={width}
    height={height}
    viewBox="0 0 55 55"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M11.6155 29.7993L2.65625 40.7482L10.9375 42.4045L14.25 52.342L22.9597 38.5222"
      stroke="black"
      strokeWidth="4"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M43.3846 29.7993L52.3438 40.7482L44.0626 42.4045L40.7501 52.342L32.0404 38.5222"
      stroke="black"
      strokeWidth="4"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M27.5 39.0918C37.5619 39.0918 45.7187 30.935 45.7187 20.873C45.7187 10.8111 37.5619 2.6543 27.5 2.6543C17.4381 2.6543 9.28125 10.8111 9.28125 20.873C9.28125 30.935 17.4381 39.0918 27.5 39.0918Z"
      stroke="black"
      strokeWidth="4"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M30.825 22.5293C30.825 24.3587 29.3419 25.8418 27.5125 25.8418C25.683 25.8418 24.2 24.3587 24.2 22.5293V19.2168C24.2 17.3874 25.683 15.9043 27.5125 15.9043C29.3419 15.9043 30.825 17.3874 30.825 19.2168V22.5293Z"
      stroke="black"
      strokeWidth="4"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
)

const Cross = ({ height, width, ...restProps }) => (
  <svg
    height={height}
    width={width}
    {...restProps}
    viewBox="0 0 20 20"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M1 19L19 1"
      stroke={restProps.color}
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M19 19L1 1"
      stroke={restProps.color}
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
)

Cross.defaultProps = {
  color: "#6D6D6D",
}

const Tooltip = ({ backgroundColor, color, className }) => (
  <svg
    width="15"
    height="16"
    viewBox="0 0 15 16"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
    className={className}
  >
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M13.117 12.3021H8.3426L4.32209 15.2888V12.3021H2.31183C1.75671 12.3021 1.3067 11.8564 1.3067 11.3066V1.35102C1.3067 0.801194 1.75671 0.355469 2.31183 0.355469H13.117C13.6721 0.355469 14.1221 0.801194 14.1221 1.35102V11.3066C14.1221 11.8564 13.6721 12.3021 13.117 12.3021Z"
      fill={backgroundColor}
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M7.51076 4.33833C7.3878 3.80981 6.88395 3.45686 6.33987 3.5181C5.79578 3.57934 5.38474 4.03526 5.38464 4.57763C5.3846 4.774 5.22385 4.93316 5.02559 4.93312C4.82734 4.93308 4.66665 4.77386 4.66669 4.5775C4.66686 3.67351 5.35196 2.9136 6.2588 2.81153C7.16565 2.70946 8.00542 3.29774 8.21038 4.17864C8.41533 5.05954 7.92019 5.95251 7.05971 6.25384C6.91623 6.30408 6.82035 6.43847 6.82043 6.5892V6.71074C6.82043 6.9071 6.65971 7.06629 6.46145 7.06629C6.26319 7.06629 6.10248 6.9071 6.10248 6.71074L6.10248 6.58957C6.10248 6.58951 6.10248 6.58963 6.10248 6.58957C6.10232 6.13743 6.38999 5.73412 6.82039 5.5834C7.33666 5.40261 7.63373 4.86684 7.51076 4.33833Z"
      fill={color}
    />
    <path
      d="M6.10248 6.58957L6.10248 6.71074C6.10248 6.9071 6.26319 7.06629 6.46145 7.06629C6.65971 7.06629 6.82043 6.9071 6.82043 6.71074V6.5892C6.82035 6.43847 6.91623 6.30408 7.05971 6.25384C7.92019 5.95251 8.41533 5.05954 8.21038 4.17864C8.00542 3.29774 7.16565 2.70946 6.2588 2.81153C5.35196 2.9136 4.66686 3.67351 4.66669 4.5775C4.66665 4.77386 4.82734 4.93308 5.02559 4.93312C5.22385 4.93316 5.3846 4.774 5.38464 4.57763C5.38474 4.03526 5.79578 3.57934 6.33987 3.5181C6.88395 3.45686 7.3878 3.80981 7.51076 4.33833C7.63373 4.86684 7.33666 5.40261 6.82039 5.5834C6.38999 5.73412 6.10232 6.13743 6.10248 6.58957ZM6.10248 6.58957C6.10248 6.58963 6.10248 6.58951 6.10248 6.58957Z"
      strokeWidth="0.25"
      strokeLinecap="round"
      strokeLinejoin="round"
      stroke={color}
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M5.92267 8.46667C5.92267 8.17213 6.16374 7.93335 6.46112 7.93335C6.75849 7.93335 6.99956 8.17213 6.99956 8.46667C6.99956 8.76121 6.75849 8.99999 6.46112 8.99999C6.16374 8.99999 5.92267 8.76121 5.92267 8.46667Z"
      fill={color}
      stroke={color}
      strokeWidth="0.1"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M1.3052 0.705146C0.977543 0.705146 0.711926 0.968233 0.711926 1.29277V10.6947C0.711926 11.0192 0.977543 11.2823 1.3052 11.2823H3.20367C3.40026 11.2823 3.55963 11.4402 3.55963 11.6349V13.7503L6.78703 11.3528C6.84864 11.3071 6.92358 11.2823 7.0006 11.2823H11.5095C11.8371 11.2823 12.1027 11.0192 12.1027 10.6947V1.29277C12.1027 0.968233 11.8371 0.705146 11.5095 0.705146H1.3052ZM0 1.29277C0 0.578791 0.584357 0 1.3052 0H11.5095C12.2303 0 12.8147 0.578791 12.8147 1.29277V10.6947C12.8147 11.4087 12.2303 11.9875 11.5095 11.9875H7.11926L3.41724 14.7375C3.30938 14.8177 3.16507 14.8306 3.04448 14.7708C2.92388 14.7111 2.8477 14.589 2.8477 14.4555V11.9875H1.3052C0.584357 11.9875 0 11.4087 0 10.6947V1.29277Z"
      fill={color}
    />
  </svg>
)

Tooltip.defaultProps = {
  backgroundColor: colors.primary,
  color: colors.grey70,
}

const KeepCircle = ({ color }) => (
  <svg
    width="58"
    height="58"
    viewBox="0 0 58 58"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M29 0.25C13.1271 0.25 0.25 13.1271 0.25 29C0.25 44.8729 13.1271 57.75 29 57.75C44.8729 57.75 57.75 44.8729 57.75 29C57.7258 13.1271 44.8487 0.25 29 0.25ZM41.0557 22.5735H39.2437L33.7595 28.9758L39.2437 35.3782H41.0557V40.5242H29.4107V35.4265H31.2227L27.6229 31.2227H26.1492V35.4265H28.2027V40.5242H16.9202V35.3782H19.3603V29V22.5735H16.9202V17.4275H19.6985V19.3361H21.1964V17.4275H23.9265V19.3361H25.4244V17.4275H28.1786V22.5011H26.125V26.7048H27.5987L31.1985 22.5011H29.3866V17.4275H41.0315V22.5735H41.0557Z"
      fill={color}
    />
  </svg>
)

KeepCircle.defaultProps = {
  color: colors.lightGrey,
}

const OK = ({ color }) => (
  <svg
    width="17"
    height="18"
    viewBox="0 0 17 18"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      d="M4.5 8.5L6.533 10.6465C6.73911 10.935 7.06732 11.1118 7.42167 11.125C7.77602 11.1382 8.11648 10.9864 8.3435 10.714L16.4375 1"
      stroke={color}
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      d="M12.5872 2.3064C9.27145 0.254412 4.95279 0.941643 2.43773 3.92149C-0.0773383 6.90135 -0.0294076 11.2741 2.55037 14.1981C5.13015 17.1221 9.46283 17.7145 12.7328 15.5903C16.0028 13.4661 17.2228 9.26677 15.6 5.72115"
      stroke={color}
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
)

OK.defaultProps = {
  color: colors.night,
}

const OKBadge = ({ bgColor, color }) => (
  <svg
    width="20"
    height="21"
    viewBox="0 0 20 21"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect fill={bgColor} y="1" width="20" height="19" rx="9.5" />
    <path
      stroke={color}
      d="M6.96484 9.67466L8.50359 11.2987C8.6596 11.5171 8.90802 11.6508 9.17622 11.6608C9.44442 11.6708 9.70211 11.5559 9.87394 11.3498L16.0002 4"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
    <path
      stroke={color}
      d="M13.0858 4.98857C10.5762 3.436 7.30743 3.95597 5.40381 6.21059C3.50019 8.46521 3.53647 11.7737 5.48907 13.9861C7.44167 16.1984 10.721 16.6467 13.196 15.0395C15.6711 13.4323 16.5945 10.2549 15.3662 7.57224"
      strokeLinecap="round"
      strokeLinejoin="round"
    />
  </svg>
)

OKBadge.defaultProps = {
  bgColor: colors.bgSuccess,
  color: colors.success,
}

const PendingBadge = ({ bgColor, color }) => (
  <svg
    width="20"
    height="21"
    viewBox="0 0 20 21"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <rect fill={bgColor} y="1" width="20" height="19" rx="9.5" />
    <svg width="14" height="14" x="3" y="3.5" fill="none">
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M7 13C10.3137 13 13 10.3137 13 7C13 3.68629 10.3137 1 7 1C3.68629 1 1 3.68629 1 7C1 10.3137 3.68629 13 7 13Z"
        stroke={color}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M7 6.99979V4.85693V6.99979Z"
        fill="#4C4C4C"
      />
      <path
        d="M7 6.99979V4.85693"
        stroke={color}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M7 7L9.67829 9.67886L7 7Z"
        fill="#4C4C4C"
      />
      <path
        d="M7 7L9.67829 9.67886"
        stroke={color}
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  </svg>
)

PendingBadge.defaultProps = {
  bgColor: colors.bgPending,
  color: colors.pending,
}

const Coinbase = () => (
  <svg
    className="wallet-icon coinbase"
    width="40"
    height="40"
    viewBox="0 0 40 40"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <circle cx="20" cy="20" r="19.5" stroke="#4C4C4C" />
    <path
      className="c-letter"
      d="M26.8025 13.3536C26.8859 13.4267 26.9576 13.5135 27 13.5579C26.0816 14.4885 25.1824 15.4008 24.2552 16.3396C23.0833 15.1985 21.6334 14.7578 19.9982 15.0038C18.7942 15.1848 17.7959 15.7896 17.0327 16.74C15.498 18.6519 15.621 21.3653 17.2336 23.1221C18.8469 24.8795 22.0454 25.4145 24.2955 23.2662C24.3549 23.31 24.424 23.3482 24.4779 23.4022C25.2774 24.1983 26.0707 24.9991 26.8743 25.7897C27.0328 25.9461 27.0396 26.0507 26.8722 26.2037C24.7861 28.1081 22.3276 28.9335 19.5472 28.499C16.0563 27.953 13.6756 25.9366 12.486 22.6123C10.8529 18.0479 13.4857 12.8555 18.1349 11.4403C21.3648 10.4571 24.2634 11.1315 26.8025 13.3536Z"
      fill="#4C4C4C"
      stroke="#4C4C4C"
    />
  </svg>
)

const Saddle = ({ className }) => {
  return (
    <img
      style={{
        width: "auto",
        height: "1.8rem",
        backgroundColor: "white",
        borderRadius: "100%",
        border: "2px solid #3800D6",
        padding: ".1rem .35rem",
      }}
      className={className}
      src={require("../static/svg/Saddle_logomark_blue.png")}
      alt="Saddle Logo"
    />
  )
}

export {
  Badge,
  Cross,
  Tooltip,
  KeepCircle,
  OK,
  OKBadge,
  PendingBadge,
  ArrowRight,
  Ledger,
  Trezor,
  MetaMask,
  Tally,
  Coinbase,
  Authorizer,
  Operations,
  KeepToken,
  Rewards,
  GrantContextIcon,
  MoneyWalletOpen,
  KeepGreenOutline,
  DashedLine,
  AuthorizerFormIcon,
  BeneficiaryFormIcon,
  OperatorFormIcon,
  KeepOutline,
  LedgerDevice,
  TrezorDevice,
  TBTC,
  KeepBlackGreen,
  Filter,
  Load,
  DocumentWithBg,
  DelegationDiagram,
  Diamond,
  ETH,
  KeepLoadingIndicator,
  KEEPTower,
  CarDashboardSpeed,
  UserFriendly,
  Fees,
  Alert,
  Success,
  Beacon,
  Authorize,
  Question,
  Home,
  FeesVector,
  Time,
  KeepDashboardLogo,
  NetworkStatusIndicator,
  Saddle,
  Add,
  Subtract,
  ArrowDown,
}
