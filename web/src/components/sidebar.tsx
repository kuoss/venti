import { useTheme } from 'next-themes'
import { useRouter } from 'next/router'
import Link from 'next/link'
import { useAuth } from '../lib/auth'

const ThemeButton = () => {
  const { theme, setTheme } = useTheme()
  return (
    <button
      onClick={() => (theme === 'dark' ? setTheme('light') : setTheme('dark'))}
      className="w-full bg-gray-300 hover:bg-gray-200 transition-all duration-100 text-white p-2"
    >
      Dark
    </button>
  )
}

const LogoutButton = () => {
  const router = useRouter()
  const auth = useAuth()

  function handleClick() {
    auth.setAuthenticated(false)
    router.push('/')
  }
  return (
    <button
      onClick={handleClick}
      className="w-full bg-gray-300 hover:bg-gray-200 transition-all duration-100 text-white p-2"
    >
      Log out
    </button>
  )
}

export default function Navbar() {
  return (
    <div className="h-screen grid grid-rows-[1fr_4fr_1fr]">
      <LogoutButton />
      <div className="p-3 text-center">Venti</div>
      <div>
        <div className="p-4">
          <Link href="/dashboards">dashboards</Link>
        </div>
        <div className="p-4">
          <Link href="/datasources">datasources</Link>
        </div>
      </div>
      <div>
        <ThemeButton />
      </div>
    </div>
  )
}
