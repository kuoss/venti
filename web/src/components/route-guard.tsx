import { ReactNode, useContext, useEffect } from 'react';
import { useRouter } from 'next/router';

import Layout from "./layout";
import { AuthContext } from '../lib/auth';

export default function RouteGuard({ children }: { children: ReactNode }) {
  const router = useRouter();
  const { authenticated } = useContext(AuthContext);
  useEffect(() => {
    if (router.pathname != '/login' && !authenticated) {
      router.push('/login')
    }
  })
  return router.pathname == '/login'
    ? <div>{children}</div>
    : (authenticated ? <Layout>{children}</Layout> : <div>Loading...</div>)
}
