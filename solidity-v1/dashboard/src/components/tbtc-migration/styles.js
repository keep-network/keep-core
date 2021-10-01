const swapBoxStyle = {
  padding: "0.5rem 0.75rem",
  borderRadius: "0.5rem",
  height: "40px",
  marginRight: "0.5rem",
}

const styles = {
  swapBox: swapBoxStyle,
  v1: {
    color: "white",
    backgroundColor: "black",
    textAlign: "center",
    ...swapBoxStyle,
  },
  v2: {
    color: "black",
    backgroundColor: "white",
    border: "1px solid black",
    textAlign: "center",
    ...swapBoxStyle,
  },
  boxWrapper: {
    padding: "1rem",
    borderRadius: "0.5rem",
  },
}

export default styles
