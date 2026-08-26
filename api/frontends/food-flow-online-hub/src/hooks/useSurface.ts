import { useLocation } from 'react-router-dom';
import { getSurface, type Surface } from '@/lib/surfaces';

export function useSurface(): Surface {
  const { pathname } = useLocation();
  return getSurface(pathname);
}
