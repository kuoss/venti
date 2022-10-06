import { fileURLToPath, URL } from 'url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default ({ mode }) => {
  process.env = {...process.env, ...loadEnv(mode, process.cwd())}

  return defineConfig({
    build: {
      sourcemap: true,
    },
    clearScreen: false,
    plugins: [vue()],
    resolve: {
      alias: {
        '@': fileURLToPath(new URL('./src', import.meta.url)),
      },
    },
    server: {
      hmr: {
        overlay: false,
        protocol: 'wss',
        host: process.env.VITE_SERVER_HMR_HOST,
        clientPort: 443,
      },
      host: true,
      proxy: {
        '/user': {
          target: 'http://localhost:8080',
        },
        '/api': {
          target: 'http://localhost:8080',
        },
      },
    },
  })
}