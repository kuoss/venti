import { defineStore } from 'pinia'

export const useConfigStore = defineStore('config', {
    state: () => ({
        dark: localStorage.theme == 'dark',
    }),
    actions: {
        setDark(b) {
            this.dark = b
            if (b) {
                document.documentElement.classList.add('dark')
                document.documentElement.style.setProperty('color-scheme', 'dark')
                localStorage.theme = 'dark'
            } else {
                document.documentElement.classList.remove('dark')
                document.documentElement.style.setProperty('color-scheme', 'normal')
                localStorage.theme = 'light'
            }
        }
    }
})