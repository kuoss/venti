import {
  createContext,
  Dispatch,
  FunctionComponent,
  ReactNode,
  SetStateAction,
  useContext,
  useMemo,
  useState,
} from 'react'

// types
type AuthContextType = {
  authenticated: boolean
  username: string
  setAuthenticated: Dispatch<SetStateAction<boolean>>
  setUsername: Dispatch<SetStateAction<string>>
}

type AuthContextProviderProps = {
  defaultAuthenticated: boolean
  defaultUsername: string
  children: ReactNode
}

// consts
const AuthContext = createContext<AuthContextType>({
  authenticated: false,
  username: '',
  setAuthenticated: () => {},
  setUsername: () => {},
})

const AuthContextProvider: FunctionComponent<AuthContextProviderProps> = ({
  defaultAuthenticated,
  defaultUsername,
  children,
}) => {
  const [authenticated, setAuthenticated] = useState(defaultAuthenticated)
  const [username, setUsername] = useState(defaultUsername)
  const contextValue = useMemo(
    () => ({
      authenticated,
      username,
      setAuthenticated,
      setUsername,
    }),
    [authenticated, username, setAuthenticated, setUsername]
  )
  return <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
}

// export
export const useAuth = () => useContext(AuthContext)

// export default
export default function AuthProvider({ children }: { children: ReactNode }) {
  return (
    <AuthContextProvider defaultAuthenticated={false} defaultUsername={''}>
      {children}
    </AuthContextProvider>
  )
}
