import { fileURLToPath, URL } from 'node:url';

import { defineConfig, loadEnv } from 'vite';
import vue from '@vitejs/plugin-vue';
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from '@tailwindcss/vite';

export default defineConfig({
  clearScreen: false,
  plugins: [vue(), vueDevTools(), tailwindcss()],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    host: true,
    proxy: {
      // local
      '/api': { target: 'http://localhost:3030' },
      '/auth': { target: 'http://localhost:3030' },
    },
  },
});


// import { fileURLToPath, URL } from 'node:url';

// import { defineConfig, loadEnv } from 'vite';
// import vue from '@vitejs/plugin-vue';
// import vueDevTools from 'vite-plugin-vue-devtools'
// import tailwindcss from '@tailwindcss/vite';

// export default defineConfig(({ mode }) => {
//   const env = loadEnv(mode, process.cwd(), '');

//   let backendURL = 'http://localhost:3030';

//   // CODESPACE_URL=https://silver-space-goggles-xrgg4p6r6vh97q6.github.dev/
//   if (env.CODESPACE_URL) {
//     try {
//       const hostname = new URL(env.CODESPACE_URL).hostname;
//       const parts = hostname.split('-');
//       const hash = parts.pop();
//       backendURL = `https://${[...parts, 3030, hash].join('-')}.github.dev`;
//     } catch (e) {
//       console.warn(`Invalid CODESPACE_URL: ${env.CODESPACE_URL}`);
//     }
//   }

//   return {
//     clearScreen: false,
//     plugins: [vue(), vueDevTools(), tailwindcss()],
//     resolve: {
//       alias: {
//         '@': fileURLToPath(new URL('./src', import.meta.url)),
//       },
//     },
//     server: {
//       host: true,
//       proxy: {
//         '/api': { target: backendURL },
//         '/auth': { target: backendURL },
//       },
//     },
//   };
// });
