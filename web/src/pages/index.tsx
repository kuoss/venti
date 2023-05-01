import { Inter } from 'next/font/google'
import { useAuth } from '../lib/auth'

const inter = Inter({ subsets: ['latin'] })

export default function Home() {
  const auth = useAuth()

  return (
    <main className={`flex min-v-screen flex-col items-center justify-between ${inter.className}`}>
      <div>
        Welcome, <b>{auth.username}</b>.
      </div>
    </main>
  )
}
