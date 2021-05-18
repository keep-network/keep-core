// TODO: Banner story

// import React from "react"
// import { Banner } from "../components/Banner"
// import { storiesOf } from "@storybook/react"
// import centered from "@storybook/addon-centered/react"
// import * as Icons from "../components/Icons";
//
// storiesOf("Banner", module).addDecorator(centered)
//
// export default {
//   title: "Banner",
//   component: Banner,
// }
//
// console.log('banner', Banner)
// console.log('banner.title', Banner.Title)
//
// const BannerTemplate = ({
//   incentivesRemoved,
//   bannerTitle,
//   bannerDescription,
//   link,
//   linkText,
//   ...props
// }) => (
//   <Banner
//     className={`liquidity__new-user-info ${
//       incentivesRemoved ? "liquidity__new-user-info--warning mt-2" : ""
//     }`}
//   >
//     <Banner.Icon
//       className={`liquidity__rewards-icon ${
//         incentivesRemoved ? "liquidity__rewards-icon--warning" : ""
//       }`}
//     />
//     <div className={"liquidity__new-user-info-text"}>
//       <Banner.Title
//         className={`liquidity-banner__title ${
//           incentivesRemoved ? "text-grey-60" : "text-white"
//         }`}
//       >
//         {bannerTitle}
//       </Banner.Title>
//       <Banner.Description
//         className={`liquidity-banner__info ${
//           incentivesRemoved ? "text-grey-60" : "text-white"
//         }`}
//       >
//         {bannerDescription}
//         &nbsp;
//         <a
//           target="_blank"
//           rel="noopener noreferrer"
//           href={link}
//           className={`text-link ${
//             incentivesRemoved ? "text-grey-60" : "text-white"
//           }`}
//         >
//           {linkText}
//         </a>
//       </Banner.Description>
//     </div>
//   </Banner>
// )
//
// export const Default = BannerTemplate.bind({})
// BannerTemplate.args = {
//   incentivesRemoved: true,
//   bannerTitle: "banner title",
//   bannerDescription: "banner description",
//   link: "https://google.com",
//   linkText: "google",
// }
//
//
