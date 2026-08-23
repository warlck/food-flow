import React, { createContext, useContext, useState, useEffect } from 'react';
import { CartItem, MenuItem, OrderType, SelectedAddon } from '@/types';
import { toast } from '@/components/ui/use-toast';
import { ValidatePromoResponse, orderService } from '@/services/orderService';

// Generate a unique ID for cart items
const generateCartItemId = () => `cart-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

interface CartContextType {
  items: CartItem[];
  orderType: OrderType;
  restaurantId: string | null;
  appliedPromo: ValidatePromoResponse | null;
  addToCart: (item: MenuItem, quantity?: number, selectedAddons?: SelectedAddon[], specialInstructions?: string) => void;
  updateQuantity: (cartItemId: string, quantity: number) => void;
  removeFromCart: (cartItemId: string) => void;
  clearCart: () => void;
  getTotalItems: () => number;
  getTotalPrice: () => number;
  hasItems: () => boolean;
  updateSpecialInstructions: (cartItemId: string, instructions: string) => void;
  setOrderType: (type: OrderType) => void;
  setRestaurantId: (id: string) => void;
  applyPromoCode: (code: string) => Promise<{ success: boolean; message: string }>;
  removePromoCode: () => void;
}

const CartContext = createContext<CartContextType | undefined>(undefined);

// Helper to calculate total price for an item including addons
const calculateItemTotal = (item: CartItem): number => {
  let itemTotal = item.menuItem.price * item.quantity;
  if (item.selectedAddons) {
    item.selectedAddons.forEach((selectedAddon) => {
      itemTotal += selectedAddon.addon.price * selectedAddon.quantity * item.quantity;
    });
  }
  return itemTotal;
};

// Helper to calculate total price for all cart items
const calculateCartTotal = (cartItems: CartItem[]): number => {
  return cartItems.reduce((total, item) => total + calculateItemTotal(item), 0);
};

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [items, setItems] = useState<CartItem[]>([]);
  const [orderType, setOrderType] = useState<OrderType>('delivery');
  const [restaurantId, setRestaurantId] = useState<string | null>(null);
  const [appliedPromo, setAppliedPromo] = useState<ValidatePromoResponse | null>(null);

  // Load cart from localStorage on initial render
  useEffect(() => {
    const savedCart = localStorage.getItem('foodFlowCart');
    const savedOrderType = localStorage.getItem('foodFlowOrderType');
    const savedRestaurantId = localStorage.getItem('foodFlowRestaurantId');
    
    if (savedRestaurantId) {
      setRestaurantId(savedRestaurantId);
    }

    if (savedCart) {
      try {
        const parsedCart = JSON.parse(savedCart);
        // Ensure all items have cartItemId (for backward compatibility)
        const cartWithIds = parsedCart.map((item: CartItem) => ({
          ...item,
          cartItemId: item.cartItemId || generateCartItemId()
        }));
        setItems(cartWithIds);
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
  }, []);

  const setRestaurantIdHandler = (id: string) => {
    if (!id || id === restaurantId) return;

    // Reset the whole cart if restaurantId differs from the loaded restaurant
    if (restaurantId && restaurantId !== id && items.length > 0) {
      setItems([]);
      setAppliedPromo(null);
      toast({
        description: "Cart reset for the selected restaurant",
        variant: "default",
      });
    }

    setRestaurantId(id);
  };

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

  // Re-validate applied promo code whenever items or restaurantId changes
  const promoCode = appliedPromo?.code;

  useEffect(() => {
    if (!promoCode) return;

    if (items.length === 0) {
      setAppliedPromo(null);
      return;
    }

    const subtotal = calculateCartTotal(items);
    let isCurrent = true;

    orderService
      .validatePromoCode({
        promoCode,
        restaurantId: restaurantId || undefined,
        subtotal,
      })
      .then((res) => {
        if (!isCurrent) return;
        if (res.valid) {
          setAppliedPromo(res);
        } else {
          setAppliedPromo(null);
          toast({
            description: `Promo code ${promoCode} is no longer applicable: ${res.reason}`,
            variant: 'destructive',
          });
        }
      })
      .catch(() => {
        if (!isCurrent) return;
        setAppliedPromo(null);
      });

    return () => {
      isCurrent = false;
    };
  }, [items, restaurantId, promoCode]);

  const addToCart = (menuItem: MenuItem, quantity: number = 1, selectedAddons?: SelectedAddon[], specialInstructions?: string) => {
    const isDifferentRestaurant = Boolean(
      menuItem.restaurantId && restaurantId && menuItem.restaurantId !== restaurantId
    );

    if (isDifferentRestaurant) {
      setAppliedPromo(null);
      toast({
        description: "Cleared items from previous restaurant",
        variant: "default",
      });
    }

    if (menuItem.restaurantId && menuItem.restaurantId !== restaurantId) {
      setRestaurantId(menuItem.restaurantId);
    }

    setItems(prevItems => {
      const currentItems = isDifferentRestaurant ? [] : prevItems;

      // When addons or special instructions are provided, always add as new item
      if ((selectedAddons && selectedAddons.length > 0) || (specialInstructions && specialInstructions.length > 0)) {
        toast({
          description: `Added ${menuItem.name} to cart`,
          variant: "default",
        });
        return [...currentItems, { 
          cartItemId: generateCartItemId(),
          menuItem, 
          quantity, 
          selectedAddons, 
          specialInstructions 
        }];
      }

      // Check if item is already in cart (without addons and without special instructions)
      const existingItemIndex = currentItems.findIndex(
        item => item.menuItem.id === menuItem.id && 
        (!item.selectedAddons || item.selectedAddons.length === 0) && 
        (!item.specialInstructions || item.specialInstructions.length === 0)
      );

      if (existingItemIndex >= 0) {
        // Update existing item quantity
        const updatedItems = [...currentItems];
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
        return [...currentItems, { 
          cartItemId: generateCartItemId(),
          menuItem, 
          quantity 
        }];
      }
    });
  };

  const updateQuantity = (cartItemId: string, quantity: number) => {
    if (quantity <= 0) {
      removeFromCart(cartItemId);
      return;
    }

    setItems(prevItems =>
      prevItems.map(item =>
        item.cartItemId === cartItemId ? { ...item, quantity } : item
      )
    );
  };

  const removeFromCart = (cartItemId: string) => {
    setItems(prevItems => {
      const itemToRemove = prevItems.find(item => item.cartItemId === cartItemId);
      if (itemToRemove) {
        toast({
          description: `Removed ${itemToRemove.menuItem.name} from cart`,
          variant: "default",
        });
      }
      return prevItems.filter(item => item.cartItemId !== cartItemId);
    });
  };

  const clearCart = () => {
    setItems([]);
    setAppliedPromo(null);
    toast({
      description: "Cart cleared",
      variant: "default",
    });
  };

  const getTotalItems = () => {
    return items.reduce((total, item) => total + item.quantity, 0);
  };

  const getTotalPrice = () => {
    return calculateCartTotal(items);
  };

  const hasItems = () => {
    return items.length > 0;
  };

  const updateSpecialInstructions = (cartItemId: string, instructions: string) => {
    setItems(prevItems =>
      prevItems.map(item =>
        item.cartItemId === cartItemId
          ? { ...item, specialInstructions: instructions }
          : item
      )
    );
  };

  const applyPromoCode = async (code: string): Promise<{ success: boolean; message: string }> => {
    if (!code.trim()) {
      return { success: false, message: 'Please enter a promo code' };
    }
    try {
      const subtotal = getTotalPrice();
      const res = await orderService.validatePromoCode({
        promoCode: code,
        restaurantId: restaurantId || undefined,
        subtotal,
      });

      if (res.valid) {
        setAppliedPromo(res);
        toast({
          description: `Promo code ${res.code} applied! Saved $${res.discountAmount.toFixed(2)}`,
          variant: "default",
        });
        return { success: true, message: `Saved $${res.discountAmount.toFixed(2)}` };
      } else {
        setAppliedPromo(null);
        return { success: false, message: res.reason || 'Invalid promo code' };
      }
    } catch (err: unknown) {
      setAppliedPromo(null);
      const message = err instanceof Error ? err.message : 'Failed to validate promo code';
      return { success: false, message };
    }
  };

  const removePromoCode = () => {
    setAppliedPromo(null);
    toast({
      description: 'Promo code removed',
      variant: 'default',
    });
  };

  return (
    <CartContext.Provider
      value={{
        items,
        orderType,
        restaurantId,
        appliedPromo,
        addToCart,
        updateQuantity,
        removeFromCart,
        clearCart,
        getTotalItems,
        getTotalPrice,
        hasItems,
        updateSpecialInstructions,
        setOrderType,
        setRestaurantId: setRestaurantIdHandler,
        applyPromoCode,
        removePromoCode,
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
