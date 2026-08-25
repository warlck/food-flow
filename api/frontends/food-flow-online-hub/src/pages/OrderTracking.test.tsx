import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import OrderTracking from './OrderTracking';
import * as reactRouter from 'react-router-dom';
import { orderService, type Order } from '@/services/orderService';

vi.mock('react-router-dom', () => ({
  useParams: vi.fn(),
  useNavigate: vi.fn(() => vi.fn()),
  useLocation: vi.fn(() => ({ pathname: '/track-order' })),
  Link: ({ children, to, ...props }: React.AnchorHTMLAttributes<HTMLAnchorElement> & { to: string }) => (
    <a href={to} {...props}>
      {children}
    </a>
  ),
}));

vi.mock('@/services/orderService', () => ({
  orderService: {
    getOrder: vi.fn(),
  },
}));

vi.mock('@/context/CartContext', () => ({
  useCart: vi.fn(() => ({
    restaurantId: null,
    hasItems: () => false,
    getTotalItems: () => 0,
    items: [],
  })),
}));

vi.mock('@/components/Layout', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="layout">{children}</div>,
}));

describe('OrderTracking Component', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('renders search form when orderId is absent', () => {
    vi.mocked(reactRouter.useParams).mockReturnValue({});

    render(<OrderTracking />);

    expect(screen.getByPlaceholderText(/Enter Order ID/i)).toBeDefined();
  });

  it('fetches and displays order details when orderId is present', async () => {
    const mockOrder: Order = {
      id: 'order-123',
      restaurantId: 'restaurant-abc',
      orderStatus: 'preparing',
      orderType: 'delivery',
      paymentMethod: 'creditCard',
      paymentStatus: 'paid',
      subtotal: 20,
      total: 25,
      deliveryFee: 5,
      tip: 0,
      items: [],
      dateCreated: new Date().toISOString(),
    };

    vi.mocked(orderService.getOrder).mockResolvedValue(mockOrder);
    vi.mocked(reactRouter.useParams).mockReturnValue({ orderId: 'order-123' });

    render(<OrderTracking />);

    await waitFor(() => {
      expect(screen.getByText(/ORDER-12/i)).toBeDefined();
    });
  });
});
