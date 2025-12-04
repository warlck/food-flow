import React, { createContext, useContext, useState, useEffect } from 'react';
import { CartItem, MenuItem, OrderType, SelectedAddon } from '@/types';
import { toast } from '@/components/ui/use-toast';

interface CartContextType {
  items: CartItem[];
  orderType: OrderType;
  restaurantId: string | null;
  addToCart: (item: MenuItem, quantity?: number, selectedAddons?: SelectedAddon[]) => void;
  updateQuantity: (itemId: string, quantity: number) => void;
  removeFromCart: (itemId: string) => void;
  clearCart: () => void;
  getTotalItems: () => number;
  getTotalPrice: () => number;
  hasItems: () => boolean;
  updateSpecialInstructions: (itemId: string, instructions: string) => void;
  setOrderType: (type: OrderType) => void;
  setRestaurantId: (id: string) => void;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [items, setItems] = useState<CartItem[]>([]);
  const [orderType, setOrderType] = useState<OrderType>('delivery');
  const [restaurantId, setRestaurantId] = useState<string | null>(null);

  // Load cart from localStorage on initial render
  useEffect(() => {
    const savedCart = localStorage.getItem('foodFlowCart');
    const savedOrderType = localStorage.getItem('foodFlowOrderType');
    const savedRestaurantId = localStorage.getItem('foodFlowRestaurantId');
    
    if (savedCart) {
      try {
        setItems(JSON.parse(savedCart));
      } catch (error) {
        console.error('Failed to parse saved cart:', error);
        localStorage.removeItem('foodFlowCart');
      }
    }
    
    if (savedOrderType) {
      try {
        setOrderType(savedOrderType as OrderType);
      } catch (error) {
        console.error('Failed to parse saved order type:', error);
        localStorage.removeItem('foodFlowOrderType');
      }
    }
    
    if (savedRestaurantId) {
      setRestaurantId(savedRestaurantId);
    }
  }, []);

  // Save cart and order type to localStorage whenever they change
  useEffect(() => {
    localStorage.setItem('foodFlowCart', JSON.stringify(items));
  }, [items]);
  
  useEffect(() => {
    localStorage.setItem('foodFlowOrderType', orderType);
  }, [orderType]);
  
  useEffect(() => {
    if (restaurantId) {
      localStorage.setItem('foodFlowRestaurantId', restaurantId);
    }
  }, [restaurantId]);

  const addToCart = (menuItem: MenuItem, quantity: number = 1, selectedAddons?: SelectedAddon[]) => {
    setItems(prevItems => {
      // When addons are provided, always add as new item (different addon selections = different items)
      if (selectedAddons && selectedAddons.length > 0) {
        toast({
          description: `Added ${menuItem.name} to cart`,
          variant: "default",
        });
        return [...prevItems, { menuItem, quantity, selectedAddons }];
      }

      // Check if item is already in cart (without addons)
      const existingItemIndex = prevItems.findIndex(
        item => item.menuItem.id === menuItem.id && (!item.selectedAddons || item.selectedAddons.length === 0)
      );

      if (existingItemIndex >= 0) {
        // Update existing item quantity
        const updatedItems = [...prevItems];
        updatedItems[existingItemIndex] = {
          ...updatedItems[existingItemIndex],
          quantity: updatedItems[existingItemIndex].quantity + quantity
        };
        toast({
          description: `Updated ${menuItem.name} quantity in cart`,
          variant: "default",
        });
        return updatedItems;
      } else {
        // Add new item to cart
        toast({
          description: `Added ${menuItem.name} to cart`,
          variant: "default",
        });
        return [...prevItems, { menuItem, quantity }];
      }
    });
  };

  const updateQuantity = (itemId: string, quantity: number) => {
    if (quantity <= 0) {
      removeFromCart(itemId);
      return;
    }

    setItems(prevItems =>
      prevItems.map(item =>
        item.menuItem.id === itemId ? { ...item, quantity } : item
      )
    );
  };

  const removeFromCart = (itemId: string) => {
    setItems(prevItems => {
      const itemToRemove = prevItems.find(item => item.menuItem.id === itemId);
      if (itemToRemove) {
        toast({
          description: `Removed ${itemToRemove.menuItem.name} from cart`,
          variant: "default",
        });
      }
      return prevItems.filter(item => item.menuItem.id !== itemId);
    });
  };

  const clearCart = () => {
    setItems([]);
    toast({
      description: "Cart cleared",
      variant: "default",
    });
  };

  const getTotalItems = () => {
    return items.reduce((total, item) => total + item.quantity, 0);
  };

  const getTotalPrice = () => {
    return items.reduce((total, item) => {
      let itemTotal = item.menuItem.price * item.quantity;
      
      // Add addon prices
      if (item.selectedAddons) {
        item.selectedAddons.forEach(selectedAddon => {
          itemTotal += selectedAddon.addon.price * selectedAddon.quantity * item.quantity;
        });
      }
      
      return total + itemTotal;
    }, 0);
  };

  const hasItems = () => {
    return items.length > 0;
  };

  const updateSpecialInstructions = (itemId: string, instructions: string) => {
    setItems(prevItems =>
      prevItems.map(item =>
        item.menuItem.id === itemId
          ? { ...item, specialInstructions: instructions }
          : item
      )
    );
  };

  return (
    <CartContext.Provider
      value={{
        items,
        orderType,
        restaurantId,
        addToCart,
        updateQuantity,
        removeFromCart,
        clearCart,
        getTotalItems,
        getTotalPrice,
        hasItems,
        updateSpecialInstructions,
        setOrderType,
        setRestaurantId
      }}
    >
      {children}
    </CartContext.Provider>
  );
};

export const useCart = (): CartContextType => {
  const context = useContext(CartContext);
  if (context === undefined) {
    throw new Error('useCart must be used within a CartProvider');
  }
  return context;
};
