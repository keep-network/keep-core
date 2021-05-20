import "./../src/css/app.css"
import "./storybookStyleFix.css"
import React from "react";
import { BrowserRouter } from "react-router-dom";
export const parameters = {
  actions: { argTypesRegex: "^on[A-Z].*" },
  controls: {
    matchers: {
      color: /(background|color)$/i,
      date: /Date$/,
    },
  },
}

export const decorators = [
  (Story) => (
    <>
      <div></div>
      <BrowserRouter>
        <div style={{display: "flex", justifyContent: "center", alignItems: "center"}}>
          <Story />
        </div>
      </BrowserRouter>
    </>
  ),
];