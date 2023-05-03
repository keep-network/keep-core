import { useLocalStorage as useRehooksLocalStorage } from "@rehooks/local-storage"

export const useLocalStorage = (key, defaultValue) => {
  return useRehooksLocalStorage(key, defaultValue)
}
