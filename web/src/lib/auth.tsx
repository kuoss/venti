import {
  createContext,
  FC,
  Dispatch,
  ReactNode,
  SetStateAction,
  useContext,
  useMemo,
  useState,
} from "react";

type AuthContextType = {
  authenticated: boolean;
  setAuthenticated: Dispatch<SetStateAction<boolean>>;
};

export const AuthContext = createContext<AuthContextType>({
  authenticated: false,
  setAuthenticated: () => { },
});

export const useAuth = () => useContext(AuthContext);

type AuthContextProviderProps = {
  defaultAuthenticated: boolean;
  children: ReactNode;
};

const AuthContextProvider: FC<AuthContextProviderProps> = ({ defaultAuthenticated, children }) => {
  const [authenticated, setAuthenticated] = useState(defaultAuthenticated);
  const contextValue = useMemo(() => ({ authenticated, setAuthenticated }), [authenticated]);
  return <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
}

function AuthProvider({ children }: { children: ReactNode }) {
  return (
    <AuthContextProvider defaultAuthenticated={false}>
      {children}
    </AuthContextProvider>
  );
}

export default AuthProvider;
