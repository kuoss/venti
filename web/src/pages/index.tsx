import { Inter } from 'next/font/google'
import { useContext } from 'react';
import { AuthContext } from '../lib/auth';

const inter = Inter({ subsets: ['latin'] })

function DisplayAuthenticated() {
  const { authenticated } = useContext(AuthContext);
  return <p>{authenticated ? 'true' : 'false'}</p>;
}


export default function Home() {
  return (
    <main className={`flex min-v-screen flex-col items-center justify-between ${inter.className}`}    >
      <div>hello world</div>
      <DisplayAuthenticated />
    </main>
  )
}
