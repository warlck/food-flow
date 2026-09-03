import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { CartProvider, useCart } from './CartContext';
import { CART_STORAGE_VERSION } from './cartStorage';
import { toast } from '@/components/ui/use-toast';

vi.mock('@/components/ui/use-toast', () => ({
  toast: vi.fn(),
}));

const CartCount = () => {
  const { items } = useCart();
  return <span data-testid="cart-count">{items.length}</span>;
};

describe('CartProvider storage migration', () => {
  beforeEach(() => {
    localStorage.clear();
    vi.clearAllMocks();
    vi.spyOn(console, 'error').mockImplementation(() => undefined);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it.each([
    ['legacy unversioned cart', JSON.stringify([{ cartItemId: 'old-line' }])],
    ['future-version cart', JSON.stringify({ version: 2, items: [] })],
    ['malformed cart', '{'],
  ])('clears a %s and notifies the user once', async (_name, savedCart) => {
    localStorage.setItem('foodFlowCart', savedCart);

    render(
      <CartProvider>
        <CartCount />
      </CartProvider>,
    );

    await waitFor(() => {
      expect(localStorage.getItem('foodFlowCart')).toBe(
        JSON.stringify({ version: CART_STORAGE_VERSION, items: [] }),
      );
    });

    expect(screen.getByTestId('cart-count').textContent).toBe('0');
    expect(toast).toHaveBeenCalledTimes(1);
    expect(toast).toHaveBeenCalledWith({
      description: 'Your cart was refreshed because the menu changed',
      variant: 'default',
    });
  });
});
