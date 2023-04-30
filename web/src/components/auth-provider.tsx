import { ReactNode } from 'react';
import { AuthContextProvider } from './auth-context-provider';

export default function AuthProvider({ children }:{children: ReactNode}) {
  return (
    <AuthContextProvider defaultAuthenticated={false}>
      { children }
    </AuthContextProvider>
  );
}