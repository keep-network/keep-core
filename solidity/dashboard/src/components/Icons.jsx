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
import { ReactComponent as KeepDashboardLogo } from "../static/svg/token-dashboard-logo.svg"

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

const Tooltip = ({ backgroundColor, color }) => (
  <svg
    width="15"
    height="16"
    viewBox="0 0 15 16"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
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
  color: colors.black,
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

const Ledger = () => (
  <svg
    className="wallet-icon ledger"
    width="40"
    height="40"
    viewBox="0 0 40 40"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <circle cx="20" cy="20" r="19.5" stroke="#4C4C4C" />
    <path
      fill="#4C4C4C"
      fillRule="evenodd"
      clipRule="evenodd"
      d="M20 39C30.4934 39 39 30.4934 39 20C39 9.50659 30.4934 1 20 1C9.50659 1 1 9.50659 1 20C1 30.4934 9.50659 39 20 39ZM20 40C31.0457 40 40 31.0457 40 20C40 8.9543 31.0457 0 20 0C8.9543 0 0 8.9543 0 20C0 31.0457 8.9543 40 20 40Z"
    />
    <path
      fill="#4C4C4C"
      d="M17.5105 21.4237C20.973 21.4237 24.4282 21.4237 27.9784 21.4237C27.9784 18.7185 28.0417 16.0498 27.9516 13.3884C27.9078 12.1417 26.8851 11.0874 25.6822 11.0484C22.9794 10.9583 20.2742 11.0216 17.5105 11.0216C17.5105 14.5279 17.5105 17.9247 17.5105 21.4237ZM11.0457 17.5351C11.0457 18.8524 11.0457 20.1502 11.0457 21.4627C12.3947 21.4627 13.6901 21.4627 14.9757 21.4627C14.9757 20.0942 14.9757 18.8378 14.9757 17.5351C13.6219 17.5351 12.3801 17.5351 11.0457 17.5351ZM21.5355 28.003C21.5355 26.8391 21.5793 25.7896 21.5063 24.7499C21.4868 24.4698 21.1313 23.9999 20.9049 23.9829C19.797 23.9025 18.6793 23.9463 17.5154 23.9463C17.5154 25.3903 17.5154 26.6613 17.5154 28.003C18.8425 28.003 20.1111 28.003 21.5355 28.003ZM15.0805 11.102C11.9564 10.7148 10.6269 12.0565 11.0896 15.1148C12.1877 15.1148 13.3127 15.1513 14.4328 15.0855C14.6616 15.0709 15.0366 14.6789 15.0512 14.4452C15.1194 13.3567 15.0805 12.2659 15.0805 11.102ZM15.0805 27.9543C15.0805 26.7879 15.117 25.6971 15.0512 24.6135C15.0366 24.3797 14.6592 23.9902 14.4303 23.9755C13.3102 23.9098 12.1853 23.9463 11.0871 23.9463C10.6342 27.0168 11.9613 28.3536 15.0805 27.9543ZM27.881 23.9463C26.7877 23.9463 25.6676 23.9171 24.55 23.9707C24.3308 23.9804 23.9534 24.3043 23.9437 24.4966C23.8852 25.6484 23.9145 26.805 23.9145 27.9348C26.9971 28.3219 28.4045 26.8926 27.881 23.9463Z"
    />
  </svg>
)

const MetaMask = () => (
  <svg
    className="wallet-icon metamask"
    width="40"
    height="40"
    viewBox="0 0 40 40"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <circle cx="20" cy="20" r="19.5" stroke="#4C4C4C" />
    <path
      d="M19.9908 26.394C20.4375 26.394 20.8809 26.4005 21.3275 26.394C21.5166 26.394 21.6634 26.3221 21.634 26.0967C21.5982 25.7896 21.559 25.4824 21.5101 25.1753C21.471 24.927 21.233 24.7114 20.9885 24.7081C20.3201 24.7016 19.6517 24.7016 18.9833 24.7081C18.7584 24.7114 18.4845 24.9368 18.4617 25.1263C18.4258 25.4106 18.3899 25.6948 18.3671 25.9791C18.341 26.3091 18.416 26.3907 18.7518 26.3973C19.1627 26.4038 19.5767 26.3973 19.9908 26.3973V26.394Z"
      fill="#4C4C4C"
    />
    <path
      d="M15.7423 21.1728C16.6189 21.4638 17.5067 21.7212 18.3945 21.9562C18.7072 22.0401 18.9305 21.8107 18.9305 21.4918C18.925 21.4638 18.9305 21.419 18.9138 21.3855C18.6346 20.8147 18.3778 20.2272 18.0651 19.6733C18.0037 19.5614 17.7077 19.4886 17.5737 19.539C16.9428 19.7852 16.323 20.0594 15.72 20.3671C15.357 20.5573 15.3682 21.0497 15.7423 21.1728Z"
      fill="#4C4C4C"
    />
    <path
      d="M24.2578 21.1728C23.3812 21.4638 22.4934 21.7212 21.6056 21.9562C21.2929 22.0401 21.0696 21.8107 21.0696 21.4918C21.0752 21.4638 21.0696 21.419 21.0863 21.3855C21.3655 20.8147 21.6224 20.2272 21.935 19.6733C21.9965 19.5614 22.2924 19.4886 22.4264 19.539C23.0573 19.7852 23.6771 20.0594 24.2801 20.3671C24.6431 20.5573 24.6319 21.0497 24.2578 21.1728Z"
      fill="#4C4C4C"
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M28.7799 11.0553C28.9315 11.1238 29.0495 11.2498 29.1079 11.4055L29.9602 13.6783C30.0031 13.7927 30.0115 13.9171 29.9845 14.0362L28.6027 20.1159L29.9679 24.2114C30.0091 24.335 30.0107 24.4684 29.9724 24.5929L28.836 28.2861C28.7359 28.6113 28.3941 28.7967 28.0669 28.7032L24.3759 27.6486L21.7789 29.4666C21.6738 29.5401 21.5487 29.5795 21.4205 29.5795H18.8636C18.7354 29.5795 18.6103 29.5401 18.5052 29.4666L15.9027 27.6448L11.9224 28.7062C11.5978 28.7927 11.2628 28.6072 11.164 28.2861L10.0276 24.5929C9.98933 24.4684 9.99088 24.335 10.0321 24.2114L11.3973 20.1159L10.0155 14.0362C9.98629 13.9075 9.99854 13.7729 10.0505 13.6515L10.9028 11.6629C11.0323 11.3608 11.3736 11.2106 11.6837 11.3192L17.2653 13.2727H22.4418L28.3002 11.0409C28.4557 10.9817 28.6283 10.9869 28.7799 11.0553ZM28.1581 12.4327L22.7793 14.4818C22.7083 14.5089 22.6329 14.5227 22.5568 14.5227H17.1591C17.0888 14.5227 17.019 14.5109 16.9526 14.4876L11.8218 12.6918L11.2795 13.9573L12.6549 20.0092C12.6802 20.1207 12.6745 20.2369 12.6384 20.3454L11.2813 24.4167L12.1818 27.3433L15.8617 26.362C16.0399 26.3145 16.23 26.3481 16.3811 26.4539L19.0606 28.3295H21.2234L23.903 26.4539C24.0572 26.3459 24.252 26.3132 24.4331 26.365L27.8214 27.3331L28.7187 24.4167L27.3616 20.3454C27.3255 20.2369 27.3198 20.1207 27.3451 20.0092L28.724 13.9418L28.1581 12.4327Z"
      fill="#4C4C4C"
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M20 39C30.4934 39 39 30.4934 39 20C39 9.50659 30.4934 1 20 1C9.50659 1 1 9.50659 1 20C1 30.4934 9.50659 39 20 39ZM20 40C31.0457 40 40 31.0457 40 20C40 8.9543 31.0457 0 20 0C8.9543 0 0 8.9543 0 20C0 31.0457 8.9543 40 20 40Z"
      fill="#4C4C4C"
    />
  </svg>
)

const Trezor = () => (
  <svg
    className="wallet-icon trezor"
    width="40"
    height="40"
    viewBox="0 0 40 40"
    fill="none"
    xmlns="http://www.w3.org/2000/svg"
  >
    <circle cx="20" cy="20" r="19.5" stroke="#4C4C4C" />
    <path
      d="M14.6813 15.3724C14.7227 14.436 14.7084 13.5538 14.8098 12.6845C15.1495 9.7996 17.6347 7.82398 20.6324 8.01241C23.3103 8.18085 25.5015 10.5462 25.5315 13.3055C25.5386 13.9892 25.5329 14.673 25.5329 15.341C26.1153 15.4995 26.6435 15.578 27.106 15.795C27.3429 15.9063 27.6099 16.266 27.6127 16.5158C27.6455 19.8204 27.6298 23.1264 27.6398 26.4325C27.6413 26.8864 27.4457 27.1248 27.0332 27.2704C24.862 28.0398 22.6979 28.832 20.5225 29.59C20.2798 29.6742 19.9558 29.6657 19.7117 29.58C17.6062 28.842 15.5149 28.0569 13.4066 27.3289C12.8413 27.1333 12.6386 26.8507 12.6443 26.2469C12.6714 23.0836 12.6757 19.9203 12.64 16.7585C12.6329 16.1047 12.8655 15.7921 13.4694 15.6679C13.8619 15.588 14.2473 15.4795 14.6813 15.3724ZM24.6493 21.5834C24.6493 20.6341 24.635 19.6834 24.655 18.7341C24.6636 18.3102 24.5379 18.0618 24.0883 17.9976C21.4632 17.6293 18.8395 17.6178 16.213 17.9847C15.7676 18.0475 15.6291 18.2716 15.6334 18.7027C15.6477 20.6027 15.6477 22.5012 15.6305 24.4012C15.6263 24.7923 15.749 24.9993 16.1302 25.1292C17.3421 25.5417 18.5383 25.9999 19.7502 26.411C19.9758 26.4881 20.2712 26.4953 20.4954 26.4182C21.7301 25.9985 22.9492 25.5332 24.1811 25.1078C24.5694 24.9736 24.665 24.7423 24.655 24.3598C24.6336 23.4362 24.6493 22.5098 24.6493 21.5834ZM17.285 15.0655C19.2363 15.0655 21.0549 15.0655 22.9949 15.0655C22.9392 14.1834 22.9734 13.3212 22.8136 12.4989C22.5709 11.2599 21.5603 10.5904 20.1656 10.5861C18.7496 10.5833 17.7261 11.2285 17.4791 12.4504C17.305 13.2955 17.3421 14.1819 17.285 15.0655Z"
      fill="#4C4C4C"
    />
    <path
      fillRule="evenodd"
      clipRule="evenodd"
      d="M20 39C30.4934 39 39 30.4934 39 20C39 9.50659 30.4934 1 20 1C9.50659 1 1 9.50659 1 20C1 30.4934 9.50659 39 20 39ZM20 40C31.0457 40 40 31.0457 40 20C40 8.9543 31.0457 0 20 0C8.9543 0 0 8.9543 0 20C0 31.0457 8.9543 40 20 40Z"
      fill="#4C4C4C"
    />
  </svg>
)

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
  KeepDashboardLogo,
}
