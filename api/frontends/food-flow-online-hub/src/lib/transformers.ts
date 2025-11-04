import { ApiRestaurantDetails, ApiCategory, ApiMenuItem } from '@/lib/api';
import { MenuItem, Restaurant } from '@/types';

/**
 * Transform API restaurant data to frontend Restaurant type
 */
export function transformApiRestaurant(apiData: ApiRestaurantDetails): Restaurant {
  return {
    id: apiData.id,
    name: apiData.name,
    description: apiData.description,
    logo: apiData.imageUrl,
    coverImage: apiData.imageUrl,
    address: apiData.address,
    phone: apiData.phone,
    email: apiData.email,
    openingHours: {
      monday: { open: "10:00", close: "22:00" },
      tuesday: { open: "10:00", close: "22:00" },
      wednesday: { open: "10:00", close: "22:00" },
      thursday: { open: "10:00", close: "22:00" },
      friday: { open: "10:00", close: "23:00" },
      saturday: { open: "11:00", close: "23:00" },
      sunday: { open: "11:00", close: "22:00" },
    },
    deliveryFee: 2.50,
    minimumOrder: 20.00,
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
 * Transform API menu items to frontend MenuItem type
 */
export function transformApiMenuItems(
  apiData: ApiRestaurantDetails
): { items: MenuItem[]; categories: string[] } {
  const items: MenuItem[] = [];
  const categoriesSet = new Set<string>();

  apiData.categories.forEach((category: ApiCategory) => {
    if (category.enabled) {
      categoriesSet.add(category.name);
      
      category.mentuItems.forEach((apiItem: ApiMenuItem) => {
        items.push({
          id: apiItem.id,
          name: apiItem.name,
          description: apiItem.description,
          price: apiItem.price,
          image: apiItem.imageUrl,
          category: category.name,
          available: apiItem.available,
          preparationTime: 15, // Default value since not in API
          restaurantId: apiData.id,
          tags: [],
        });
      });
    }
  });

  return {
    items,
    categories: Array.from(categoriesSet),
  };
}
