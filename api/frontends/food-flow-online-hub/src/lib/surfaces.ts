export type Surface = 'marketing' | 'app';

/**
 * Routes rendered on the dark marketing surface. Everything else is 'app'.
 * Matched by exact path.
 */
const MARKETING_ROUTES = ['/', '/contact', '/support', '/faq', '/privacy', '/terms'];

export function getSurface(pathname: string): Surface {
  return MARKETING_ROUTES.includes(pathname) ? 'marketing' : 'app';
}
