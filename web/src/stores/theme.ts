import { defineStore } from 'pinia';

export const useThemeStore = defineStore('theme', {
  state: () => ({
    dark: false,
  }),
  actions: {
    init() {
      const isDark = localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches);
      this.setDark(isDark);
    },
    setDark(b: boolean) {
      this.dark = b;
      if (b) {
        document.documentElement.classList.add('dark');
        document.documentElement.style.setProperty('color-scheme', 'dark');
        localStorage.theme = 'dark';
      } else {
        document.documentElement.classList.remove('dark');
        document.documentElement.style.setProperty('color-scheme', 'normal');
        localStorage.theme = 'light';
      }
    },
  },
});