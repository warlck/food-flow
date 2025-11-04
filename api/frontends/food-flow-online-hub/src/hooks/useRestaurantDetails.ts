import { useQuery } from '@tanstack/react-query';
import { restaurantApi, ApiRestaurantDetails } from '@/lib/api';

export function useRestaurantDetails(restaurantId: string) {
  return useQuery<ApiRestaurantDetails, Error>({
    queryKey: ['restaurant', restaurantId],
    queryFn: () => restaurantApi.getRestaurantDetails(restaurantId),
    enabled: !!restaurantId,
    staleTime: 5 * 60 * 1000, // 5 minutes
    retry: 2,
  });
}
