import { useEffect } from 'react';
import { useLocation, useNavigationType } from 'react-router-dom';

/**
 * Resets scroll on PUSH/REPLACE navigations. POP (browser back/forward)
 * keeps the browser's restored position.
 */
export const ScrollToTop = () => {
  const { pathname, hash } = useLocation();
  const navigationType = useNavigationType();

  useEffect(() => {
    if (navigationType === 'POP') return;
    if (hash) return;

    window.scrollTo({ top: 0, behavior: 'auto' });
  }, [pathname, hash, navigationType]);

  return null;
};
