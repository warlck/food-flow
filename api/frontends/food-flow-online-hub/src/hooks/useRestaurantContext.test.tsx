import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import React from 'react';
import { useRestaurantContext } from './useRestaurantContext';
import { ActiveRestaurantContext } from '@/context/RestaurantContext';
import * as reactRouter from 'react-router-dom';
import * as CartContextModule from '@/context/CartContext';

vi.mock('react-router-dom', () => ({
  useParams: vi.fn(),
  useLocation: vi.fn(() => ({ pathname: '/' })),
}));

vi.mock('@/context/CartContext', () => ({
  useCart: vi.fn(),
}));

describe('useRestaurantContext', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const wrapperWithActiveRestaurant = (activeRestaurantId: string | null = null) => {
    return ({ children }: { children: React.ReactNode }) => (
      <ActiveRestaurantContext.Provider value={{ activeRestaurantId, setActiveRestaurantId: vi.fn() }}>
        {children}
      </ActiveRestaurantContext.Provider>
    );
  };

  it('resolves none when cart is empty and no route or order context exists (fresh or stale localStorage)', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: 'stale_restaurant_id',
      hasItems: () => false,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext(), {
      wrapper: wrapperWithActiveRestaurant(null),
    });

    expect(result.current).toEqual({
      restaurantId: null,
      source: 'none',
    });
  });

  it('resolves cart context when user has items in cart', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: 'cart_restaurant_id',
      hasItems: () => true,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext(), {
      wrapper: wrapperWithActiveRestaurant(null),
    });

    expect(result.current).toEqual({
      restaurantId: 'cart_restaurant_id',
      source: 'cart',
    });
  });

  it('resolves order context when orderRestaurantId parameter is provided', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: null,
      hasItems: () => false,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext('order_restaurant_id'), {
      wrapper: wrapperWithActiveRestaurant(null),
    });

    expect(result.current).toEqual({
      restaurantId: 'order_restaurant_id',
      source: 'order',
    });
  });

  it('prioritizes order context over cart context when tracking an order while holding a cart for another restaurant', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: 'cart_restaurant_b',
      hasItems: () => true,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext('order_restaurant_a'), {
      wrapper: wrapperWithActiveRestaurant(null),
    });

    expect(result.current).toEqual({
      restaurantId: 'order_restaurant_a',
      source: 'order',
    });
  });

  it('resolves order context from ActiveRestaurantContext when passed as active in provider', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: null,
      hasItems: () => false,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext(), {
      wrapper: wrapperWithActiveRestaurant('active_published_restaurant_id'),
    });

    expect(result.current).toEqual({
      restaurantId: 'active_published_restaurant_id',
      source: 'order',
    });
  });

  it('resolves route context with highest precedence over order and cart', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({ restaurantId: 'route_restaurant_id' });
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: 'cart_restaurant_id',
      hasItems: () => true,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    const { result } = renderHook(() => useRestaurantContext('order_restaurant_id'), {
      wrapper: wrapperWithActiveRestaurant('active_restaurant_id'),
    });

    expect(result.current).toEqual({
      restaurantId: 'route_restaurant_id',
      source: 'route',
    });
  });
});
