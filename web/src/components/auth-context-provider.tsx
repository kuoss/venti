import { useState, ReactNode, FC, useRef, useEffect, useMemo } from 'react';
import { AuthContext } from '../lib/auth';

type AuthContextProviderProps = {
  defaultAuthenticated?: boolean;
  // onLogin?: () => void;
  // onLogout?: () => void;
  children: ReactNode;
};

// function usePrevious<T = any>(value: T) {
//   const ref = useRef<T>();
//   useEffect(() => {
//     ref.current = value;
//   }, [value]);
//   return ref.current;
// }

export const AuthContextProvider: FC<AuthContextProviderProps> = ({
  defaultAuthenticated = false,
  // onLogin,
  // onLogout,
  children,
}) => {
  const [authenticated, setAuthenticated] = useState(
    defaultAuthenticated
  );
  // const previousAuthenticated = usePrevious(authenticated);
  // useEffect(() => {
  //   if (!previousAuthenticated && authenticated) {
  //     onLogin && onLogin();
  //   }
  // }, [previousAuthenticated, authenticated, onLogin]);

  // useEffect(() => {
  //   if (previousAuthenticated && !authenticated) {
  //     onLogout && onLogout();
  //   }
  // }, [previousAuthenticated, authenticated, onLogout]);

  const contextValue = useMemo(
    () => ({
      authenticated,
      setAuthenticated,
    }),
    [authenticated]
  );

  return (
    <AuthContext.Provider value={contextValue}>{children}</AuthContext.Provider>
  );
}

