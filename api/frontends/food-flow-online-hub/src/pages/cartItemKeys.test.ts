import { describe, expect, it } from 'vitest';
import cartDesktopSource from './CartDesktop.tsx?raw';
import restaurantSource from './Restaurant.tsx?raw';

describe('cart item list keys', () => {
  it.each([
    ['desktop cart', cartDesktopSource],
    ['restaurant cart', restaurantSource],
  ])('uses cart line identity in the %s', (_name, source) => {
    expect(source).not.toContain('key={item.menuItem.id}');
    expect(source).toContain('key={item.cartItemId}');
  });
});
