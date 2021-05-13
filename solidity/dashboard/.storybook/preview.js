import "./../src/css/app.css"
import React from "react";
export const parameters = {
  layout: 'centered',
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
}

// export const decorators = [
//   (Story) => (
//     <div style={{ display: "flex" }}>
//       <Story />
//     </div>
//   ),
// ];
