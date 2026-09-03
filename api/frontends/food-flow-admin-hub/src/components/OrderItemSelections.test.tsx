import React from 'react';
import { render, screen } from '@testing-library/react';
import { describe, expect, it } from 'vitest';
import { OrderItemSelections } from './OrderItemSelections';
import type { AdminOrderItem } from '@/lib/admin-api';

describe('OrderItemSelections', () => {
  it('renders snapshot names and prices without resolving historical catalog IDs', () => {
    const item: AdminOrderItem = {
      id: 'order-item-1',
      menuItemId: 'deleted-menu-item',
      menuItemName: 'Historical Burger',
      menuItemPrice: 10,
      quantity: 2,
      dateCreated: '2026-09-03T00:00:00Z',
      modifiers: [{
        id: 'snapshot-modifier',
        modifierGroupId: 'deleted-group',
        modifierGroupName: 'Historical Size',
        modifierOptionId: 'deleted-option',
        modifierOptionName: 'Historical Large',
        priceDelta: 0.5,
      }],
      addons: [{
        id: 'snapshot-addon',
        addonId: 'deleted-addon',
        addonName: 'Historical Sauce',
        addonPrice: 1.25,
        quantity: 2,
      }],
    };

    render(<OrderItemSelections item={item} />);

    expect(screen.getByText('+ Historical Large (Historical Size)')).toBeInTheDocument();
    expect(screen.getByText('+$1.00')).toBeInTheDocument();
    expect(screen.getByText('+ Historical Sauce ×2')).toBeInTheDocument();
    expect(screen.getByText('+$5.00')).toBeInTheDocument();
  });
});
