import { defineStore } from 'pinia'
import type { ToastServiceMethods } from 'primevue'

export const useSystemStore = defineStore('system', () => {
  let toast: ToastServiceMethods | null = null

  const setToast = (t: ToastServiceMethods) => {
    toast = t
  }

  return { toast, setToast }
})
