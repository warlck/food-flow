import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import React from 'react';
import MenuGrid from './MenuGrid';
import type { MenuItem } from '@/types';
import * as CartContextModule from '@/context/CartContext';

vi.mock('@/context/CartContext', () => ({
  useCart: vi.fn(),
}));

const addToCart = vi.fn();

const item = (id: string, name: string, category: string, price: number): MenuItem => ({
  id,
  name,
  description: `Description of ${name}`,
  price,
  image: `${id}.jpg`,
  category,
  available: true,
  orderable: true,
  preparationTime: 15,
  restaurantId: 'rest-1',
});

const chicken = item('item-chicken', 'Chicken Kebab Roll', 'Kebab Roll', 11.0);
const beef = item('item-beef', 'Beef Kebab Roll', 'Kebab Roll', 12.0);
const salad = item('item-salad', 'Garden Salad', 'Sides', 6.0);

const renderGrid = (items: MenuItem[]) =>
  render(<MenuGrid items={items} categories={['Kebab Roll', 'Sides']} onCartUpdate={vi.fn()} />);

describe('MenuGrid per-item cards (spec §10.2)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    // jsdom lacks matchMedia; the shadcn sidebar queries it on mount.
    if (!window.matchMedia) {
      window.matchMedia = ((query: string) => ({
        matches: false,
        media: query,
        onchange: null,
        addListener: vi.fn(),
        removeListener: vi.fn(),
        addEventListener: vi.fn(),
        removeEventListener: vi.fn(),
        dispatchEvent: vi.fn(),
      })) as unknown as typeof window.matchMedia;
    }
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      addToCart,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);
  });

  it('renders every sibling dish as its own card, not one representative per category', () => {
    renderGrid([chicken, beef, salad]);

    // Sibling dishes in the same category are both visible...
    expect(screen.getByText('Chicken Kebab Roll')).toBeDefined();
    expect(screen.getByText('Beef Kebab Roll')).toBeDefined();
    // ...and so are dishes from every other category.
    expect(screen.getByText('Garden Salad')).toBeDefined();
  });

  it('shows all sibling dishes of a selected category, not only the top-ranked one', () => {
    renderGrid([chicken, beef, salad]);

    fireEvent.click(screen.getByRole('button', { name: 'Kebab Roll' }));

    expect(screen.getByText('Chicken Kebab Roll')).toBeDefined();
    expect(screen.getByText('Beef Kebab Roll')).toBeDefined();
    expect(screen.queryByText('Garden Salad')).toBeNull();
  });

  it('search reveals every matching sibling instead of the first match per category', () => {
    renderGrid([chicken, beef, salad]);

    fireEvent.change(screen.getByPlaceholderText('Search menu items...'), {
      target: { value: 'roll' },
    });

    expect(screen.getByText('Chicken Kebab Roll')).toBeDefined();
    expect(screen.getByText('Beef Kebab Roll')).toBeDefined();
    expect(screen.queryByText('Garden Salad')).toBeNull();
  });

  it('prefixes the card price with From when a required group has differing option deltas', () => {
    const configurable: MenuItem = {
      ...chicken,
      modifierGroups: [
        {
          id: 'grp-protein',
          name: 'Choose a protein',
          minSelections: 1,
          maxSelections: 1,
          available: true,
          options: [
            { id: 'opt-mild', name: 'Mild', priceDelta: 0, available: true },
            { id: 'opt-hot', name: 'Hot', priceDelta: 1, available: true },
          ],
        },
      ],
    };

    renderGrid([configurable]);

    expect(screen.getByText('From $11.00')).toBeDefined();
  });
});
