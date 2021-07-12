import React from "react";

export const centeredWithFullWidth = (Story) => (
  <div style={{display: "flex", overflow: "auto", alignItems: "center", inset: "0px", position: "fixed"}}>
    <div style={{margin: "auto", maxHeight: "100%", width: "100%"}}>
      <Story />
    </div>
  </div>
)

export const whiteBackground = (Story) => (
  <div style={{backgroundColor: "white", padding: "30px", borderRadius: "5px"}}>
    <Story />
  </div>
)

export const blackBackground = (Story) => (
  <div style={{backgroundColor: "black", padding: "30px", borderRadius: "5px"}}>
    <Story />
  </div>
)