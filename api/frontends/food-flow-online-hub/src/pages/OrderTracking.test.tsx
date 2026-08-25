import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import OrderTracking from './OrderTracking';
import { RestaurantContextProvider } from '@/context/RestaurantContextProvider';
import { useActiveRestaurant } from '@/context/RestaurantContext';
import { useRestaurantContext } from '@/hooks/useRestaurantContext';
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

// Test harness that also displays the current restaurant context
const TestTrackerHarness: React.FC = () => {
  const { activeRestaurantId } = useActiveRestaurant();
  const context = useRestaurantContext();

  return (
    <div>
      <div data-testid="active-id">{activeRestaurantId || 'none'}</div>
      <div data-testid="context-source">{context.source}</div>
      <div data-testid="context-restaurant-id">{context.restaurantId || 'none'}</div>
      <OrderTracking />
    </div>
  );
};

describe('OrderTracking Context Lifecycle', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('publishes restaurantId when orderId is present and clears it when orderId becomes absent', async () => {
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

    // Initial mount at /track-order/order-123
    vi.mocked(reactRouter.useParams).mockReturnValue({ orderId: 'order-123' });

    const { rerender } = render(
      <RestaurantContextProvider>
        <TestTrackerHarness />
      </RestaurantContextProvider>
    );

    // Verify order was fetched and restaurantId was published
    await waitFor(() => {
      expect(screen.getByTestId('active-id').textContent).toBe('restaurant-abc');
      expect(screen.getByTestId('context-source').textContent).toBe('order');
      expect(screen.getByTestId('context-restaurant-id').textContent).toBe('restaurant-abc');
    });

    // Simulate clicking "Track Order" in header -> orderId goes present -> absent
    vi.mocked(reactRouter.useParams).mockReturnValue({});

    rerender(
      <RestaurantContextProvider>
        <TestTrackerHarness />
      </RestaurantContextProvider>
    );

    // Verify activeRestaurantId is immediately cleared and source returns to 'none'
    await waitFor(() => {
      expect(screen.getByTestId('active-id').textContent).toBe('none');
      expect(screen.getByTestId('context-source').textContent).toBe('none');
      expect(screen.getByTestId('context-restaurant-id').textContent).toBe('none');
    });
  });
});
