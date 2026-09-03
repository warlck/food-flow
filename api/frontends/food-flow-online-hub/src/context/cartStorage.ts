import type { CartItem } from '@/types';

export const CART_STORAGE_VERSION = 1;

interface PersistedCart {
  version: typeof CART_STORAGE_VERSION;
  items: CartItem[];
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null;

const isCartItem = (value: unknown): value is CartItem => {
  if (!isRecord(value) || !isRecord(value.menuItem)) {
    return false;
  }

  return (
    typeof value.cartItemId === 'string' &&
    value.cartItemId.length > 0 &&
    typeof value.menuItem.id === 'string' &&
    value.menuItem.id.length > 0 &&
    typeof value.quantity === 'number' &&
    Number.isInteger(value.quantity) &&
    value.quantity > 0 &&
    (value.selectedModifiers === undefined || Array.isArray(value.selectedModifiers)) &&
    (value.selectedAddons === undefined || Array.isArray(value.selectedAddons))
  );
};

export const parsePersistedCart = (rawCart: string): CartItem[] => {
  const parsed: unknown = JSON.parse(rawCart);
  if (
    !isRecord(parsed) ||
    parsed.version !== CART_STORAGE_VERSION ||
    !Array.isArray(parsed.items) ||
    !parsed.items.every(isCartItem)
  ) {
    throw new Error('unsupported or malformed cart storage');
  }

  return parsed.items;
};

export const serializePersistedCart = (items: CartItem[]): string =>
  JSON.stringify({
    version: CART_STORAGE_VERSION,
    items,
  } satisfies PersistedCart);
