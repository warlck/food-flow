import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import React from 'react';
import MenuItemDialog from './MenuItemDialog';
import type { MenuItem, ModifierGroup } from '@/types';
import * as CartContextModule from '@/context/CartContext';

vi.mock('@/context/CartContext', () => ({
  useCart: vi.fn(),
}));

const addToCart = vi.fn();

const baseItem: MenuItem = {
  id: 'item-1',
  name: 'Test Burger',
  description: 'Test description',
  price: 12,
  image: 'test.jpg',
  category: 'cat-1',
  available: true,
  orderable: true,
  preparationTime: 10,
  restaurantId: 'rest-1',
};

const activeRequiredGroup: ModifierGroup = {
  id: 'grp-active',
  name: 'Active Spice Level',
  minSelections: 1,
  maxSelections: 1,
  available: true,
  options: [
    { id: 'opt-mild', name: 'Mild Option', priceDelta: 0, available: true },
    { id: 'opt-hot', name: 'Hot Option', priceDelta: 1, available: true },
  ],
};

const suspendedRequiredGroup: ModifierGroup = {
  id: 'grp-suspended',
  name: 'Suspended Size Group',
  minSelections: 1,
  maxSelections: 1,
  available: false,
  options: [{ id: 'opt-large', name: 'Suspended Large Option', priceDelta: 2, available: true }],
};

describe('MenuItemDialog unavailable modifier group rendering', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(CartContextModule.useCart).mockReturnValue({
      addToCart,
    } as unknown as ReturnType<typeof CartContextModule.useCart>);
  });

  const renderDialog = (groups?: ModifierGroup[]) =>
    render(
      <MenuItemDialog
        item={{ ...baseItem, modifierGroups: groups }}
        isOpen
        onClose={vi.fn()}
      />
    );

  it('renders only available modifier groups and never shows the suspended group', () => {
    renderDialog([activeRequiredGroup, suspendedRequiredGroup]);

    expect(screen.getByText('Active Spice Level')).toBeDefined();
    expect(screen.queryByText('Suspended Size Group')).toBeNull();
    expect(screen.queryByText('Suspended Large Option')).toBeNull();
    // Only the active required group contributes a Required badge; the
    // suspended group must not render one.
    expect(screen.getAllByText('Required - Choose 1')).toHaveLength(1);
  });

  it('hides the modifier groups section entirely when every group is unavailable', () => {
    renderDialog([suspendedRequiredGroup]);

    expect(screen.queryByText('Suspended Size Group')).toBeNull();
    expect(screen.queryByText('Suspended Large Option')).toBeNull();
  });

  it('preselects only available single-select groups and submits their selection', () => {
    renderDialog([activeRequiredGroup, suspendedRequiredGroup]);

    const addButton = screen.getByRole('button', { name: /add to cart/i });
    expect(addButton).not.toHaveProperty('disabled', true);
    fireEvent.click(addButton);

    expect(addToCart).toHaveBeenCalledTimes(1);
    const [, , submittedModifiers] = addToCart.mock.calls[0];
    expect(submittedModifiers).toHaveLength(1);
    expect(submittedModifiers[0]).toMatchObject({
      modifierGroupId: 'grp-active',
      modifierOptionId: 'opt-mild',
    });
  });

  it('renders a None option for optional groups and clears selection when clicked', () => {
    const optionalGroup: ModifierGroup = {
      id: 'grp-opt',
      name: 'Extra Sauce',
      minSelections: 0,
      maxSelections: 1,
      available: true,
      options: [
        { id: 'opt-bbq', name: 'BBQ Sauce', priceDelta: 0.5, available: true },
      ],
    };

    renderDialog([optionalGroup]);

    expect(screen.getByText('None')).toBeDefined();
    expect(screen.getByText('No selection')).toBeDefined();

    // Click BBQ option
    fireEvent.click(screen.getByText('BBQ Sauce'));

    // Now click None
    fireEvent.click(screen.getByText('None'));

    const addButton = screen.getByRole('button', { name: /add to cart/i });
    fireEvent.click(addButton);

    expect(addToCart).toHaveBeenCalledTimes(1);
    const [, , submittedModifiers] = addToCart.mock.calls[0];
    expect(submittedModifiers).toBeUndefined();
  });
});
