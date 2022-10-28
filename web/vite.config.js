import { fileURLToPath, URL } from 'url'

import { defineConfig, loadEnv } from 'vite'
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
    // CASE            port    hmr.clientPort hmr.protocol
    // ingress-http    (80)    80             ws
    // ingress-https   (443)   443            wss
    // auto-forwarded  9090    (9090)         ws
    server: {
      port: 9090,
      hmr: {
        overlay: false,
        protocol: 'ws',
        host: process.env.VITE_SERVER_HMR_HOST,
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
