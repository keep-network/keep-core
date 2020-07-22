import { useContext } from "react"
import { ModalContext } from "../components/Modal"

export const useModal = () => {
  return useContext(ModalContext)
}
