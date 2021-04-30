import { useEffect, useState } from "react"

const useCurrentWidthLowerThanThreshold = (widthThreshold) => {
  const [width, setWidth] = useState(window.innerWidth)

  const handleWindowSizeChange = () => {
    setWidth(window.innerWidth)
  }

  useEffect(() => {
    window.addEventListener("resize", handleWindowSizeChange)
    return () => {
      window.removeEventListener("resize", handleWindowSizeChange)
    }
  }, [])

  return width < widthThreshold
}

export default useCurrentWidthLowerThanThreshold
