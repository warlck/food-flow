import { useParams } from 'react-router-dom';
import { useCart } from '@/context/CartContext';
import { useActiveRestaurant } from '@/context/RestaurantContext';

export interface RestaurantContext {
  /** Restaurant to navigate to, or null when there is no legitimate context. */
  restaurantId: string | null;
  /** Where the value came from — useful for debugging and tests. */
  source: 'route' | 'order' | 'cart' | 'none';
}

/**
 * Resolves the restaurant the current view belongs to.
 *
 * Precedence (highest first):
 *   1. route      - the URL names a restaurant (/restaurant/:id)
 *   2. order      - an entity on screen declares its restaurant (passed in or active in context)
 *   3. cart       - the user has a NON-EMPTY cart, i.e. an in-progress order
 *   4. none       - no restaurant context; restaurant nav must be hidden
 *
 * A persisted restaurantId with an EMPTY cart is deliberately NOT a context.
 * That is the stale-localStorage case that causes the Menu link to leak.
 */
export function useRestaurantContext(orderRestaurantId?: string | null): RestaurantContext {
  const { restaurantId: routeId } = useParams<{ restaurantId?: string }>();
  const { restaurantId: cartRestaurantId, hasItems } = useCart();
  const { activeRestaurantId } = useActiveRestaurant();

  const effectiveOrderId = orderRestaurantId ?? activeRestaurantId;

  if (routeId) return { restaurantId: routeId, source: 'route' };
  if (effectiveOrderId) return { restaurantId: effectiveOrderId, source: 'order' };
  if (cartRestaurantId && hasItems && hasItems()) {
    return { restaurantId: cartRestaurantId, source: 'cart' };
  }
  return { restaurantId: null, source: 'none' };
}
