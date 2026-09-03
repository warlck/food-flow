import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, within } from '@testing-library/react';
import React from 'react';
import MobileRestaurant from './MobileRestaurant';
import type { ApiRestaurantDetails } from '@/lib/api';
import * as CartContextModule from '@/context/CartContext';
import * as useRestaurantDetailsModule from '@/hooks/useRestaurantDetails';

vi.mock('react-router-dom', () => ({
  useParams: vi.fn(() => ({ restaurantId: 'rest-1' })),
  useSearchParams: vi.fn(() => [new URLSearchParams()]),
}));

vi.mock('@/context/CartContext', () => ({
  useCart: vi.fn(),
}));

vi.mock('@/hooks/useRestaurantDetails', () => ({
  useRestaurantDetails: vi.fn(),
}));

const apiItem = (id: string, name: string, price: number) => ({
  id,
  name,
  description: `Description of ${name}`,
  price,
  imageUrl: `${id}.jpg`,
  available: true,
});

const apiData = {
  id: 'rest-1',
  name: 'Test Restaurant',
  description: 'Test kitchen',
  logoUrl: 'logo.jpg',
  imageUrl: 'cover.jpg',
  address: '1 Test Street',
  phone: '555-0100',
  email: 'orders@test.example',
  enabled: true,
  operatingHours: {},
  latitude: 0,
  longitude: 0,
  maxDeliveryDistanceKm: 10,
  taxRate: 0.08,
  minSpend: 0,
  categories: [
    {
      id: 'category-kebab',
      name: 'Kebab Roll',
      description: 'Rolled kebabs',
      enabled: true,
      menuItems: [
        apiItem('item-chicken', 'Chicken Kebab Roll', 11.0),
        apiItem('item-beef', 'Beef Kebab Roll', 12.0),
      ],
    },
    {
      id: 'category-sides',
      name: 'Sides',
      description: 'Sides',
      enabled: true,
      menuItems: [apiItem('item-salad', 'Garden Salad', 6.0)],
    },
  ],
} as unknown as ApiRestaurantDetails;

describe('MobileRestaurant per-item cards (spec §10.2)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(useRestaurantDetailsModule.useRestaurantDetails).mockReturnValue({
      data: apiData,
      isLoading: false,
      error: null,
    } as unknown as ReturnType<typeof useRestaurantDetailsModule.useRestaurantDetails>);
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      getTotalItems: () => 0,
      getTotalPrice: () => 0,
      setRestaurantId: vi.fn(),
    } as unknown as ReturnType<typeof CartContextModule.useCart>);
  });

  it('renders every sibling dish as its own card with its own name and price', () => {
    render(<MobileRestaurant />);

    // Both siblings are visible as individual cards (previously one
    // representative per category)...
    expect(screen.getByText('Chicken Kebab Roll')).toBeDefined();
    expect(screen.getByText('Beef Kebab Roll')).toBeDefined();
    // ...and each card quotes the dish's own price, not the category cheapest.
    expect(screen.getByText('$11.00')).toBeDefined();
    expect(screen.getByText('$12.00')).toBeDefined();
    expect(screen.getByText('$6.00')).toBeDefined();
  });

  it('shows all sibling dishes of a selected category, not only the top-ranked one', () => {
    render(<MobileRestaurant />);

    fireEvent.click(screen.getByRole('button', { name: 'Kebab Roll' }));

    expect(screen.getByText('Chicken Kebab Roll')).toBeDefined();
    expect(screen.getByText('Beef Kebab Roll')).toBeDefined();
    expect(screen.queryByText('Garden Salad')).toBeNull();
  });

  it('activating a card opens that exact dish in the item dialog', async () => {
    render(<MobileRestaurant />);

    fireEvent.click(screen.getByText('Beef Kebab Roll'));

    const dialog = await screen.findByRole('dialog');
    expect(within(dialog).getByText('Beef Kebab Roll')).toBeDefined();
    expect(within(dialog).getByText('$12.00')).toBeDefined();
  });
});
