module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  darkMode: 'class',
  plugins: [
    require('nightwind'),
    require('tailwind-scrollbar'),
  ],
  theme: {
    extend: {
      // animation: {
      //   'reverse-spin': 'reverse-spin 3s linear infinite'
      // },
      // keyframes: {
      //   'reverse-spin': {
      //     from: {
      //       transform: 'rotate(360deg)'
      //     },
      //   }
      // }
    },
  },
}
