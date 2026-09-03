import { ApiRestaurantDetails, ApiCategory, ApiMenuItem, ApiAddon, ApiModifierGroup } from '@/lib/api';
import { MenuItem, Restaurant, Addon, ModifierGroup } from '@/types';

/**
 * Transform API restaurant data to frontend Restaurant type
 */
export function transformApiRestaurant(apiData: ApiRestaurantDetails): Restaurant {
  const defaultHours = {
    monday: { open: "10:00", close: "22:00", isClosed: false },
    tuesday: { open: "10:00", close: "22:00", isClosed: false },
    wednesday: { open: "10:00", close: "22:00", isClosed: false },
    thursday: { open: "10:00", close: "22:00", isClosed: false },
    friday: { open: "10:00", close: "23:00", isClosed: false },
    saturday: { open: "11:00", close: "23:00", isClosed: false },
    sunday: { open: "11:00", close: "22:00", isClosed: false },
  };

  const openingHours = apiData.operatingHours && Object.keys(apiData.operatingHours).length > 0
    ? Object.entries(apiData.operatingHours).reduce((acc, [day, sched]) => {
        acc[day] = { open: sched.open, close: sched.close, isClosed: sched.isClosed };
        return acc;
      }, {} as Record<string, { open: string; close: string; isClosed?: boolean }>)
    : defaultHours;

  return {
    id: apiData.id,
    name: apiData.name,
    description: apiData.description,
    logo: apiData.logoUrl || apiData.imageUrl,
    coverImage: apiData.imageUrl,
    address: apiData.address,
    phone: apiData.phone,
    email: apiData.email,
    enabled: apiData.enabled,
    openingHours,
    latitude: apiData.latitude,
    longitude: apiData.longitude,
    maxDeliveryDistanceKm: apiData.maxDeliveryDistanceKm,
    taxRate: apiData.taxRate,
    deliveryFee: 2.50,
    minimumOrder: apiData.minSpend ?? 0,
    estimatedDeliveryTime: {
      min: 25,
      max: 45,
    },
    estimatedPickupTime: {
      min: 15,
      max: 25,
    },
    rating: 4.6,
  };
}

/**
 * Transform API menu items to frontend MenuItem type.
 * Note: The backend sales-service `/v1/restaurants/:id/details` orders menu items,
 * modifier groups, options, and addons by `rank` (ranked items first, then tiebreakers).
 * Iterating in order preserves the display ordering on the storefront.
 */
export function transformApiMenuItems(
  apiData: ApiRestaurantDetails
): { items: MenuItem[]; categories: string[] } {
  const items: MenuItem[] = [];
  const categoriesSet = new Set<string>();

  apiData.categories.forEach((category: ApiCategory) => {
    if (category.enabled) {
      categoriesSet.add(category.name);
      
      const categoryItems = category.menuItems ?? [];
      categoryItems.forEach((apiItem: ApiMenuItem) => {
        // Transform modifier groups if they exist
        const modifierGroups: ModifierGroup[] = apiItem.modifierGroups?.map((mg: ApiModifierGroup) => ({
          id: mg.id,
          name: mg.name,
          description: mg.description,
          minSelections: mg.minSelections,
          maxSelections: mg.maxSelections,
          available: mg.available,
          rank: mg.rank,
          options: (mg.options ?? []).map((opt) => ({
            id: opt.id,
            name: opt.name,
            description: opt.description,
            priceDelta: opt.priceDelta,
            available: opt.available,
            rank: opt.rank,
          })),
        })) || [];

        // Transform addons if they exist (ordered by rank from backend)
        const addons: Addon[] = apiItem.addons?.map((apiAddon: ApiAddon) => ({
          id: apiAddon.id,
          addonId: apiAddon.addonId || apiAddon.id,
          name: apiAddon.name,
          description: apiAddon.description,
          price: apiAddon.price,
          available: apiAddon.available,
          maxQuantity: apiAddon.maxQuantity,
          rank: apiAddon.rank,
        })) || [];

        items.push({
          id: apiItem.id,
          name: apiItem.name,
          description: apiItem.description,
          price: apiItem.price,
          image: apiItem.imageUrl,
          category: category.name,
          available: apiItem.available,
          orderable: apiItem.orderable ?? apiItem.available,
          preparationTime: 15,
          restaurantId: apiData.id,
          rank: apiItem.rank,
          tags: [],
          modifierGroups: modifierGroups.length > 0 ? modifierGroups : undefined,
          addons: addons.length > 0 ? addons : undefined,
        });
      });
    }
  });

  return {
    items,
    categories: Array.from(categoriesSet),
  };
}

export interface MinOrderablePrice {
  /** Base price plus the minimum available delta from each active required group. */
  amount: number;
  /** True when the quoted amount is only a floor, so the card must say "From". */
  fromPrefix: boolean;
}

/**
 * Minimum orderable price for a menu-item card (spec §10.2): the base price
 * plus the minimum available delta from each active required group. Optional
 * groups contribute zero. `fromPrefix` is true when any active group offers
 * more than one available option with different deltas, or when optional
 * selections can increase the price; otherwise the amount is exact.
 */
export function minOrderablePrice(item: MenuItem): MinOrderablePrice {
  let amount = item.price;
  let fromPrefix = false;

  for (const group of item.modifierGroups ?? []) {
    if (!group.available) continue;
    const availableDeltas = group.options
      .filter((opt) => opt.available)
      .map((opt) => opt.priceDelta);
    if (availableDeltas.length === 0) continue;

    if (group.minSelections > 0) {
      // Required group: contributes its cheapest available delta.
      amount += Math.min(...availableDeltas);
      if (availableDeltas.length > 1 && new Set(availableDeltas).size > 1) {
        fromPrefix = true;
      }
    } else if (availableDeltas.some((delta) => delta > 0)) {
      // Optional group: contributes zero, but selections can raise the price.
      fromPrefix = true;
    }
  }

  return { amount, fromPrefix };
}
