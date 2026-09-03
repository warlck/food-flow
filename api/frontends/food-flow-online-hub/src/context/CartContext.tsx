import React, { createContext, useContext, useState, useEffect } from 'react';
import { CartItem, MenuItem, OrderType, SelectedAddon, SelectedModifier } from '@/types';
import { toast } from '@/components/ui/use-toast';
import { ValidatePromoResponse, orderService } from '@/services/orderService';
import { parsePersistedCart, serializePersistedCart } from './cartStorage';

// Generate a unique ID for cart items
const generateCartItemId = () => `cart-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;

export interface CartContextType {
  items: CartItem[];
  orderType: OrderType;
  restaurantId: string | null;
  appliedPromo: ValidatePromoResponse | null;
  addToCart: (
    item: MenuItem,
    quantity?: number,
    selectedModifiers?: SelectedModifier[],
    selectedAddons?: SelectedAddon[],
    specialInstructions?: string
  ) => void;
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

// Helper to calculate total price for an item including modifiers and addons
export const calculateItemUnitPrice = (
  item: MenuItem,
  selectedModifiers?: SelectedModifier[],
  selectedAddons?: SelectedAddon[]
): number => {
  let unit = item.price;
  if (selectedModifiers) {
    selectedModifiers.forEach((m) => {
      unit += m.priceDelta;
    });
  }
  if (selectedAddons) {
    selectedAddons.forEach((a) => {
      unit += a.addon.price * a.quantity;
    });
  }
  return unit;
};

export const calculateItemTotal = (item: CartItem): number => {
  const unitPrice = item.unitPrice ?? calculateItemUnitPrice(item.menuItem, item.selectedModifiers, item.selectedAddons);
  return unitPrice * item.quantity;
};

// Helper to calculate total price for all cart items
export const calculateCartTotal = (cartItems: CartItem[]): number => {
  return cartItems.reduce((total, item) => total + calculateItemTotal(item), 0);
};

export const generateCustomizationKey = (
  menuItemId: string,
  modifiers?: SelectedModifier[],
  addons?: SelectedAddon[],
  instructions?: string
): string => {
  const sortedMods = [...(modifiers ?? [])]
    .sort((a, b) => a.modifierOptionId.localeCompare(b.modifierOptionId))
    .map((m) => `${m.modifierGroupId}:${m.modifierOptionId}`)
    .join('|');

  const sortedAddons = [...(addons ?? [])]
    .sort((a, b) => (a.addon.addonId || a.addon.id).localeCompare(b.addon.addonId || b.addon.id))
    .map((a) => `${a.addon.addonId || a.addon.id}:${a.quantity}`)
    .join('|');

  const instr = (instructions ?? '').trim();
  return `${menuItemId}__m[${sortedMods}]__a[${sortedAddons}]__i[${instr}]`;
};

export const CartProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [items, setItems] = useState<CartItem[]>([]);
  const [orderType, setOrderType] = useState<OrderType>('delivery');
  const [restaurantId, setRestaurantId] = useState<string | null>(null);
  const [appliedPromo, setAppliedPromo] = useState<ValidatePromoResponse | null>(null);
  const [cartStorageLoaded, setCartStorageLoaded] = useState(false);

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
        setItems(parsePersistedCart(savedCart));
      } catch (error) {
        console.error('Failed to parse saved cart:', error);
        localStorage.removeItem('foodFlowCart');
        toast({
          description: 'Your cart was refreshed because the menu changed',
          variant: 'default',
        });
      }
    }

    setCartStorageLoaded(true);
    
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
    if (!cartStorageLoaded) return;
    localStorage.setItem('foodFlowCart', serializePersistedCart(items));
  }, [cartStorageLoaded, items]);
  
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

  const addToCart = (
    menuItem: MenuItem,
    quantity: number = 1,
    selectedModifiers?: SelectedModifier[],
    selectedAddons?: SelectedAddon[],
    specialInstructions?: string
  ) => {
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

    setItems((prevItems) => {
      const currentItems = isDifferentRestaurant ? [] : prevItems;
      const customKey = generateCustomizationKey(
        menuItem.id,
        selectedModifiers,
        selectedAddons,
        specialInstructions
      );

      const existingIdx = currentItems.findIndex((item) => {
        const itemKey = generateCustomizationKey(
          item.menuItem.id,
          item.selectedModifiers,
          item.selectedAddons,
          item.specialInstructions
        );
        return itemKey === customKey;
      });

      const unitPrice = calculateItemUnitPrice(menuItem, selectedModifiers, selectedAddons);

      if (existingIdx >= 0) {
        const updatedItems = [...currentItems];
        updatedItems[existingIdx] = {
          ...updatedItems[existingIdx],
          quantity: updatedItems[existingIdx].quantity + quantity,
          unitPrice,
        };
        toast({
          description: `Updated ${menuItem.name} quantity in cart`,
          variant: "default",
        });
        return updatedItems;
      } else {
        toast({
          description: `Added ${menuItem.name} to cart`,
          variant: "default",
        });
        return [
          ...currentItems,
          {
            cartItemId: generateCartItemId(),
            menuItem,
            quantity,
            selectedModifiers: selectedModifiers && selectedModifiers.length > 0 ? selectedModifiers : undefined,
            selectedAddons: selectedAddons && selectedAddons.length > 0 ? selectedAddons : undefined,
            specialInstructions: specialInstructions && specialInstructions.trim() !== '' ? specialInstructions.trim() : undefined,
            unitPrice,
          },
        ];
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
