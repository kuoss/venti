import { ReactNode, useEffect } from 'react'
import { useRouter } from 'next/router'

import Layout from './layout'
import { useAuth } from '../lib/auth'

export default function RouteGuard({ children }: { children: ReactNode }) {
  const router = useRouter()
  const auth = useAuth()
  useEffect(() => {
    if (router.pathname !== '/login' && !auth.authenticated) {
      router.push('/login')
    }
  })
  return router.pathname === '/login' ? (
    <div>{children}</div>
  ) : auth.authenticated ? (
    <Layout>{children}</Layout>
  ) : (
    <div>Loading...</div>
  )
}
