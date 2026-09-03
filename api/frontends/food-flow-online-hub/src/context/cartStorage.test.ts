import { describe, expect, it } from 'vitest';
import type { CartItem } from '@/types';
import {
  CART_STORAGE_VERSION,
  parsePersistedCart,
  serializePersistedCart,
} from './cartStorage';

const cartItem: CartItem = {
  cartItemId: 'cart-line-1',
  menuItem: {
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
  },
  quantity: 1,
  selectedModifiers: [],
  selectedAddons: [],
  unitPrice: 12,
};

describe('cart storage', () => {
  it('round-trips the supported version without inventing selections', () => {
    const serialized = serializePersistedCart([cartItem]);

    expect(JSON.parse(serialized)).toEqual({
      version: CART_STORAGE_VERSION,
      items: [cartItem],
    });
    expect(parsePersistedCart(serialized)).toEqual([cartItem]);
  });

  it.each([
    ['legacy unversioned cart', JSON.stringify([cartItem])],
    ['missing version', JSON.stringify({ items: [cartItem] })],
    ['future version', JSON.stringify({ version: 2, items: [cartItem] })],
    ['malformed version', JSON.stringify({ version: '1', items: [cartItem] })],
    ['malformed JSON', '{'],
    [
      'line without stable identity',
      JSON.stringify({ version: CART_STORAGE_VERSION, items: [{ ...cartItem, cartItemId: '' }] }),
    ],
  ])('rejects a %s', (_name, serialized) => {
    expect(() => parsePersistedCart(serialized)).toThrow();
  });
});
