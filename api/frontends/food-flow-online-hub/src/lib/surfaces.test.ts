import { describe, it, expect } from 'vitest';
import { getSurface } from './surfaces';

describe('getSurface', () => {
  it('maps marketing routes to marketing surface', () => {
    const marketingRoutes = ['/', '/contact', '/support', '/faq', '/privacy', '/terms'];
    for (const route of marketingRoutes) {
      expect(getSurface(route)).toBe('marketing');
    }
  });

  it('maps transactional and app routes to app surface', () => {
    const appRoutes = [
      '/track-order',
      '/track-order/ord_123',
      '/order-tracking',
      '/order-tracking/ord_123',
      '/restaurant/rest_123',
      '/mobile-restaurant/rest_123',
      '/cart',
      '/mobile-cart',
      '/checkout',
      '/mobile-checkout',
      '/order-confirmation/ord_123',
      '/some-unknown-route'
    ];
    for (const route of appRoutes) {
      expect(getSurface(route)).toBe('app');
    }
  });
});
