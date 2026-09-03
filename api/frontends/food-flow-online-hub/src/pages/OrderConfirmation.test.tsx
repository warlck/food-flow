import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import React from 'react';
import OrderConfirmation from './OrderConfirmation';
import * as reactRouter from 'react-router-dom';
import { orderService, type Order } from '@/services/orderService';

vi.mock('react-router-dom', () => ({
  useParams: vi.fn(),
  useNavigate: vi.fn(() => vi.fn()),
  useSearchParams: vi.fn(() => [new URLSearchParams()]),
}));

vi.mock('@/services/orderService', () => ({
  orderService: {
    getOrder: vi.fn(),
  },
}));

vi.mock('@/components/Layout', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="layout">{children}</div>,
}));

const orderWithSnapshots: Order = {
  id: 'order-123',
  restaurantId: 'restaurant-abc',
  customerName: 'Test Customer',
  customerEmail: 'customer@test.example',
  customerPhone: '555-0100',
  orderType: 'delivery',
  orderStatus: 'confirmed',
  paymentStatus: 'paid',
  paymentMethod: 'creditCard',
  subtotal: 25.0,
  deliveryFee: 2.5,
  tax: 2.0,
  total: 29.5,
  items: [
    {
      id: 'order-item-1',
      menuItemId: 'menu-item-1',
      menuItemName: 'Kebab Roll',
      menuItemPrice: 11.0,
      quantity: 2,
      modifiers: [
        {
          id: 'modifier-1',
          modifierGroupId: 'group-protein',
          modifierGroupName: 'Choose a protein',
          modifierOptionId: 'option-beef',
          modifierOptionName: 'Beef',
          priceDelta: 1.0,
        },
      ],
      addons: [
        {
          id: 'addon-1',
          addonId: 'addon-garlic-sauce',
          addonName: 'Garlic Sauce',
          addonPrice: 1.5,
          quantity: 1,
        },
      ],
    },
  ],
  dateCreated: new Date().toISOString(),
} as unknown as Order;

describe('OrderConfirmation receipt snapshots (spec §10.5)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(reactRouter.useParams).mockReturnValue({ orderId: 'order-123' });
  });

  it('displays modifier group/option names and deltas from the order snapshots', async () => {
    vi.mocked(orderService.getOrder).mockResolvedValue(orderWithSnapshots);

    render(<OrderConfirmation />);

    await waitFor(() => {
      expect(screen.getByText(/Order Confirmed/i)).toBeDefined();
    });

    // Item line renders.
    expect(screen.getByText('2x Kebab Roll')).toBeDefined();
    // Modifier snapshot renders with its group and per-line delta (1.00 x 2).
    expect(screen.getByText('+ Beef (Choose a protein)')).toBeDefined();
    expect(screen.getByText('+$2.00')).toBeDefined();
    // Addon snapshot still renders (1.50 x 1 x 2).
    expect(screen.getByText('+ Garlic Sauce x1')).toBeDefined();
    expect(screen.getByText('+$3.00')).toBeDefined();
  });

  it('renders item lines without modifiers when the item has none', async () => {
    vi.mocked(orderService.getOrder).mockResolvedValue({
      ...orderWithSnapshots,
      items: [
        {
          id: 'order-item-2',
          menuItemId: 'menu-item-2',
          menuItemName: 'Garden Salad',
          menuItemPrice: 6.0,
          quantity: 1,
        },
      ],
    } as unknown as Order);

    render(<OrderConfirmation />);

    await waitFor(() => {
      expect(screen.getByText('1x Garden Salad')).toBeDefined();
    });
    expect(screen.queryByText(/Choose a protein/i)).toBeNull();
  });
});
