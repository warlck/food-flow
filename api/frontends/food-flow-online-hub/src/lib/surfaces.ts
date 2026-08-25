export type Surface = 'marketing' | 'app';

/**
 * Routes rendered on the dark marketing surface. Everything else is 'app'.
 * Matched by exact path or prefix.
 */
const MARKETING_EXACT = ['/', '/contact', '/support', '/faq', '/privacy', '/terms'];

export function getSurface(pathname: string): Surface {
  if (MARKETING_EXACT.includes(pathname)) return 'marketing';
  if (pathname.startsWith('/track-order') || pathname.startsWith('/order-tracking')) {
    return 'marketing';
  }
  return 'app';
}
