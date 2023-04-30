import '@/styles/globals.css'

import type { AppProps } from 'next/app'
import { ThemeProvider } from "next-themes"
import AuthProvider from '../lib/auth'
import RouteGuard from '../components/route-guard';

export default function App({ Component, pageProps }: AppProps) {
  return (
    <ThemeProvider attribute="class">
      <AuthProvider>
        <RouteGuard>
          <Component {...pageProps} />
        </RouteGuard>
      </AuthProvider>
    </ThemeProvider>
  );
}
