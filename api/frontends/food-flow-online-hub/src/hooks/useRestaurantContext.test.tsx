import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { useRestaurantContext } from './useRestaurantContext';
import * as reactRouter from 'react-router-dom';
import * as CartContextModule from '@/context/CartContext';

vi.mock('react-router-dom', () => ({
  useParams: vi.fn(),
}));

vi.mock('@/context/CartContext', () => ({
  useCart: vi.fn(),
}));

describe('useRestaurantContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('resolves null when cart is empty and no route context exists', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: 'stale_restaurant_id',
      hasItems: () => false,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext());

    expect(result.current).toEqual({
      restaurantId: null,
    });
  });

  it('resolves cart restaurant when user has items in cart', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: 'cart_restaurant_id',
      hasItems: () => true,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext());

    expect(result.current).toEqual({
      restaurantId: 'cart_restaurant_id',
    });
  });

  it('prioritizes route param over cart restaurant', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({ restaurantId: 'route_restaurant_id' });
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: 'cart_restaurant_id',
      hasItems: () => true,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext());

    expect(result.current).toEqual({
      restaurantId: 'route_restaurant_id',
    });
  });
});
