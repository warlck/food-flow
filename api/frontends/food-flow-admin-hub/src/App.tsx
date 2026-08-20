import Admin from '@/pages/Admin';
import Login from '@/pages/Login';
import { Toaster } from '@/components/ui/sonner';
import { AuthProvider, useAuth } from '@/context/AuthContext';

function AuthGate() {
  const { token } = useAuth();
  return token ? <Admin /> : <Login />;
}

export default function App() {
  return (
    <AuthProvider>
      <AuthGate />
      <Toaster position="top-right" richColors />
    </AuthProvider>
  );
}
