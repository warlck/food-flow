import { describe, expect, it } from 'vitest';
import type { ApiRestaurantDetails } from './api';
import { transformApiMenuItems, minOrderablePrice } from './transformers';
import type { MenuItem, ModifierGroup } from '@/types';

const menuItem = {
  id: 'item-1',
  name: 'Chicken Kebab Roll',
  description: 'Freshly prepared',
  price: 11,
  imageUrl: 'item.jpg',
  available: true,
};

const restaurantDetails = (category: Record<string, unknown>): ApiRestaurantDetails =>
  ({
    categories: [
      {
        id: 'category-1',
        name: 'Kebab Roll',
        description: 'Rolled kebabs',
        enabled: true,
        ...category,
      },
    ],
  }) as ApiRestaurantDetails;

describe('transformApiMenuItems', () => {
  it('reads menu items from the canonical menuItems key', () => {
    const result = transformApiMenuItems(restaurantDetails({ menuItems: [menuItem] }));

    expect(result.items.map((item) => item.id)).toEqual(['item-1']);
    expect(result.categories).toEqual(['Kebab Roll']);
  });

  it('does not read the retired misspelled compatibility key', () => {
    const retiredKey = ['mentu', 'Items'].join('');
    const result = transformApiMenuItems(restaurantDetails({ [retiredKey]: [menuItem] }));

    expect(result.items).toEqual([]);
    expect(result.categories).toEqual(['Kebab Roll']);
  });
});

const group = (overrides: Partial<ModifierGroup>): ModifierGroup => ({
  id: 'grp-1',
  name: 'Group',
  minSelections: 1,
  maxSelections: 1,
  available: true,
  options: [],
  ...overrides,
});

const option = (id: string, priceDelta: number, available = true) => ({
  id,
  name: `Option ${id}`,
  priceDelta,
  available,
});

const cardItem = (modifierGroups?: ModifierGroup[]): MenuItem => ({
  id: 'item-1',
  name: 'Kebab Roll',
  description: 'Rolled kebab',
  price: 11.0,
  image: 'item.jpg',
  category: 'Mains',
  available: true,
  orderable: true,
  preparationTime: 15,
  restaurantId: 'rest-1',
  ...(modifierGroups ? { modifierGroups } : {}),
});

describe('minOrderablePrice (spec §10.2)', () => {
  it('quotes the exact base price with no prefix when there are no modifier groups', () => {
    expect(minOrderablePrice(cardItem())).toEqual({ amount: 11.0, fromPrefix: false });
  });

  it('adds the minimum available delta of each active required group exactly', () => {
    // Required protein group deltas: 0.00 / 1.00 / 2.50 -> contributes 0.00.
    const protein: ModifierGroup = group({
      id: 'grp-protein',
      minSelections: 1,
      maxSelections: 1,
      options: [
        option('chicken', 0.0),
        option('beef', 1.0),
        option('mix', 2.5),
      ],
    });
    // Required sauce group deltas: 0.50 / 0.50 -> contributes 0.50.
    const sauce: ModifierGroup = group({
      id: 'grp-sauce',
      minSelections: 1,
      maxSelections: 1,
      options: [option('garlic', 0.5), option('hot', 0.5)],
    });

    expect(minOrderablePrice(cardItem([protein, sauce]))).toEqual({ amount: 11.5, fromPrefix: true });
  });

  it('omits the prefix when a required group offers several options at the same delta', () => {
    const sameDelta: ModifierGroup = group({
      minSelections: 1,
      maxSelections: 1,
      options: [option('a', 1.0), option('b', 1.0)],
    });

    expect(minOrderablePrice(cardItem([sameDelta]))).toEqual({ amount: 12.0, fromPrefix: false });
  });

  it('omits the prefix for a required group with a single available option', () => {
    const single: ModifierGroup = group({
      minSelections: 1,
      maxSelections: 1,
      options: [option('only', 1.5), option('sold-out', 3.0, false)],
    });

    // Only the available option counts: 11.00 + 1.50, no prefix.
    expect(minOrderablePrice(cardItem([single]))).toEqual({ amount: 12.5, fromPrefix: false });
  });

  it('contributes zero from optional groups but prefixes when they can raise the price', () => {
    const optional: ModifierGroup = group({
      id: 'grp-extras',
      minSelections: 0,
      maxSelections: 3,
      options: [option('cheese', 2.0), option('sauce', 1.0)],
    });

    expect(minOrderablePrice(cardItem([optional]))).toEqual({ amount: 11.0, fromPrefix: true });
  });

  it('contributes zero and prefixes nothing when optional options cannot raise the price', () => {
    const optional: ModifierGroup = group({
      id: 'grp-free',
      minSelections: 0,
      maxSelections: 1,
      options: [option('none', 0.0), option('also-none', 0.0)],
    });

    expect(minOrderablePrice(cardItem([optional]))).toEqual({ amount: 11.0, fromPrefix: false });
  });

  it('ignores unavailable groups entirely', () => {
    const suspended: ModifierGroup = group({
      id: 'grp-suspended',
      available: false,
      options: [option('x', 5.0)],
    });

    expect(minOrderablePrice(cardItem([suspended]))).toEqual({ amount: 11.0, fromPrefix: false });
  });

  it('skips required groups with no available option (defensive; catalog rejects this state)', () => {
    const soldOut: ModifierGroup = group({
      id: 'grp-soldout',
      options: [option('x', 1.0, false)],
    });

    expect(minOrderablePrice(cardItem([soldOut]))).toEqual({ amount: 11.0, fromPrefix: false });
  });
});
