import React, { useState } from 'react';
import { ActiveRestaurantContext } from './RestaurantContext';

export const RestaurantContextProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [activeRestaurantId, setActiveRestaurantId] = useState<string | null>(null);

  return (
    <ActiveRestaurantContext.Provider value={{ activeRestaurantId, setActiveRestaurantId }}>
      {children}
    </ActiveRestaurantContext.Provider>
  );
};
