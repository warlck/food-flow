import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { CategoryRail } from './Admin';
import type { AdminCategory } from '@/lib/admin-api';

const sampleCategories: AdminCategory[] = [
  {
    id: 'cat-1',
    name: 'Kebab Roll',
    description: 'Rolls',
    restaurantId: 'rest-1',
    enabled: true,
    rank: 10,
    dateCreated: '2026-01-01T00:00:00Z',
    dateUpdated: '2026-01-01T00:00:00Z',
  },
  {
    id: 'cat-2',
    name: 'Pide & Pizza',
    description: 'Pide',
    restaurantId: 'rest-1',
    enabled: true,
    rank: 20,
    dateCreated: '2026-01-01T00:00:00Z',
    dateUpdated: '2026-01-01T00:00:00Z',
  },
];

describe('CategoryRail Component', () => {
  it('renders "+ Category" button and calls onAdd when clicked', () => {
    const onAdd = vi.fn();
    const onSelect = vi.fn();
    const onAddItem = vi.fn();
    const onEdit = vi.fn();

    render(
      <CategoryRail
        categories={sampleCategories}
        counts={new Map([['cat-1', 5], ['cat-2', 3]])}
        total={8}
        selected="all"
        onSelect={onSelect}
        onAdd={onAdd}
        onAddItem={onAddItem}
        onEdit={onEdit}
      />
    );

    // Assert "+ Category" button exists and has correct accessible label
    const categoryButton = screen.getByRole('button', { name: /add category/i });
    expect(categoryButton).toBeInTheDocument();
    expect(categoryButton).toHaveTextContent(/Category/i);

    // Click "+ Category" button
    fireEvent.click(categoryButton);
    expect(onAdd).toHaveBeenCalledTimes(1);

    // Assert that "+ Item" button in the categories section header does NOT exist
    const itemButton = screen.queryByRole('button', { name: /^add menu item$/i });
    expect(itemButton).not.toBeInTheDocument();
  });

  it('renders category items and triggers onSelect when clicked', () => {
    const onSelect = vi.fn();

    render(
      <CategoryRail
        categories={sampleCategories}
        counts={new Map([['cat-1', 5], ['cat-2', 3]])}
        total={8}
        selected="cat-1"
        onSelect={onSelect}
        onAdd={vi.fn()}
        onAddItem={vi.fn()}
        onEdit={vi.fn()}
      />
    );

    const cat2Button = screen.getByText('Pide & Pizza');
    fireEvent.click(cat2Button);
    expect(onSelect).toHaveBeenCalledWith('cat-2');
  });

  it('renders category rank badges when rank is present', () => {
    render(
      <CategoryRail
        categories={sampleCategories}
        counts={new Map([['cat-1', 5], ['cat-2', 3]])}
        total={8}
        selected="all"
        onSelect={vi.fn()}
        onAdd={vi.fn()}
        onAddItem={vi.fn()}
        onEdit={vi.fn()}
      />
    );

    expect(screen.getByText('#10')).toBeInTheDocument();
    expect(screen.getByText('#20')).toBeInTheDocument();
  });

  it('triggers onMove when Move up / Move down buttons are clicked', () => {
    const onMove = vi.fn();

    render(
      <CategoryRail
        categories={sampleCategories}
        counts={new Map([['cat-1', 5], ['cat-2', 3]])}
        total={8}
        selected="all"
        onSelect={vi.fn()}
        onAdd={vi.fn()}
        onAddItem={vi.fn()}
        onEdit={vi.fn()}
        onMove={onMove}
      />
    );

    const moveDownButtons = screen.getAllByTitle('Move down');
    fireEvent.click(moveDownButtons[0]);
    expect(onMove).toHaveBeenCalledWith(0, 'down');

    const moveUpButtons = screen.getAllByTitle('Move up');
    fireEvent.click(moveUpButtons[1]);
    expect(onMove).toHaveBeenCalledWith(1, 'up');
  });
});

