import { describe, expect, it } from 'vitest';
import type { ApiRestaurantDetails } from './api';
import { transformApiMenuItems } from './transformers';

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
