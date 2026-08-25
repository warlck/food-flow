import { createContext, useContext } from 'react';

export interface ActiveRestaurantContextType {
  activeRestaurantId: string | null;
  setActiveRestaurantId: (id: string | null) => void;
}

export const ActiveRestaurantContext = createContext<ActiveRestaurantContextType>({
  activeRestaurantId: null,
  setActiveRestaurantId: () => {},
});

export function useActiveRestaurant(): ActiveRestaurantContextType {
  return useContext(ActiveRestaurantContext);
}
