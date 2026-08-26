import { useParams } from 'react-router-dom';
import { useCart } from '@/context/CartContext';

export interface RestaurantContext {
  restaurantId: string | null;
}

/**
 * Resolves the restaurant context for the current view.
 * 1. Route param (/restaurant/:id)
 * 2. Active non-empty cart
 */
export function useRestaurantContext(): RestaurantContext {
  const { restaurantId: routeId } = useParams<{ restaurantId?: string }>();
  const { restaurantId: cartRestaurantId, hasItems } = useCart();

  if (routeId) return { restaurantId: routeId };
  if (cartRestaurantId && hasItems && hasItems()) {
    return { restaurantId: cartRestaurantId };
  }
  return { restaurantId: null };
}
