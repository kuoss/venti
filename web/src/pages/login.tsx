import { useContext } from 'react'
import { useRouter } from 'next/router'
import { AuthContext } from '../lib/auth'

function LoginButton() {
  const router = useRouter()
  const { setAuthenticated } = useContext(AuthContext)
  function handleClick() {
    setAuthenticated(true)
    router.push('/')
  }
  return (
    <button onClick={handleClick}>
      Log in
    </button>
  )
}


export default function Login() {
  return (
    <div>
      <div className="login">
        <LoginButton />
      </div>
      {/* <style jsx>{`
        .login {
          max-width: 21rem;
          margin: 0 auto;
          padding: 1rem;
          border: 1px solid #ccc;
          border-radius: 4px;
        }
      `}</style> */}
    </div>
  )
}