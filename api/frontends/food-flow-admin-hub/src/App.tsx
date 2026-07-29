import Admin from '@/pages/Admin';
import { Toaster } from '@/components/ui/sonner';

export default function App() {
  return (
    <>
      <Admin />
      <Toaster position="top-right" richColors />
    </>
  );
}
