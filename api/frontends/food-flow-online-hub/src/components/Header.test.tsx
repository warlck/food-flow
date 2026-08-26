import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import React from 'react';
import Header from './Header';
import * as reactRouter from 'react-router-dom';
import * as CartContextModule from '@/context/CartContext';

vi.mock('react-router-dom', () => ({
  useParams: vi.fn(() => ({})),
  useLocation: vi.fn(() => ({ pathname: '/' })),
  Link: ({ children, to, ...props }: React.AnchorHTMLAttributes<HTMLAnchorElement> & { to: string }) => (
    <a href={to} {...props}>
      {children}
    </a>
  ),
}));

vi.mock('@/context/CartContext', () => ({
  useCart: vi.fn(),
}));

describe('Header Menu & Cart Visibility Matrix', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  const renderHeader = (pathname: string, cartItemsCount = 0, cartRestaurantId: string | null = null) => {
    vi.mocked(reactRouter.useLocation).mockReturnValue({
      pathname,
      search: '',
      hash: '',
      state: null,
      key: 'test',
    });

    vi.mocked(CartContextModule.useCart).mockReturnValue({
      restaurantId: cartRestaurantId,
      hasItems: () => cartItemsCount > 0,
      getTotalItems: () => cartItemsCount,
      items: [],
    } as unknown as ReturnType<typeof CartContextModule.useCart>);

    return render(<Header />);
  };

  it('hides BOTH Menu and Cart on landing page (/) even if user had previous cart', () => {
    renderHeader('/', 2, 'rest_123');

    expect(screen.queryByText('Menu')).toBeNull();
    expect(document.querySelector('a[href="/cart"]')).toBeNull();
    expect(screen.getByText('Track Order')).toBeDefined();
    expect(screen.getByText('Partner with us')).toBeDefined();
  });

  it('hides Menu, Cart, and Track Order links on /track-order when cart is empty', () => {
    renderHeader('/track-order', 0, 'stale_id');

    expect(screen.queryByText('Menu')).toBeNull();
    expect(document.querySelector('a[href="/cart"]')).toBeNull();
    expect(screen.queryByText('Track Order')).toBeNull();
  });

  it('hides Menu, Cart, and Track Order links on /track-order/:id even when cart has items', () => {
    renderHeader('/track-order/ord_123', 2, 'rest_123');

    expect(screen.queryByText('Menu')).toBeNull();
    expect(document.querySelector('a[href="/cart"]')).toBeNull();
    expect(screen.queryByText('Track Order')).toBeNull();
  });

  it('shows Cart and Track Order on /restaurant/:id and hides Menu on self restaurant page', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({ restaurantId: 'rest_page_123' });
    renderHeader('/restaurant/rest_page_123', 1, 'rest_page_123');

    expect(screen.queryByText('Menu')).toBeNull(); // self-link suppressed
    expect(document.querySelector('a[href="/cart"]')).not.toBeNull();
    expect(screen.getByText('Track Order')).toBeDefined();
  });

  it('shows Menu, Cart, and Track Order on /cart (app surface) when cart has items', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});
    renderHeader('/cart', 2, 'cart_rest_123');

    expect(screen.getByText('Menu')).toBeDefined();
    expect(document.querySelector('a[href="/cart"]')).not.toBeNull();
    expect(screen.getByText('Track Order')).toBeDefined();
  });
});
