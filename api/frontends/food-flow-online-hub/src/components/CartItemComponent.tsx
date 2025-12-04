import React, { useState, useMemo } from 'react';
import { Button } from '@/components/ui/button';
import { useCart } from '@/context/CartContext';
import { CartItem } from '@/types';
import { Trash2, Minus, Plus } from 'lucide-react';
import { Input } from '@/components/ui/input';

interface CartItemProps {
  item: CartItem;
}

const CartItemComponent: React.FC<CartItemProps> = ({ item }) => {
  const { updateQuantity, removeFromCart } = useCart();
  const { cartItemId, menuItem, quantity, selectedAddons } = item;

  // Calculate total price including addons
  const itemTotalPrice = useMemo(() => {
    let total = menuItem.price * quantity;
    if (selectedAddons) {
      selectedAddons.forEach(({ addon, quantity: addonQty }) => {
        total += addon.price * addonQty * quantity;
      });
    }
    return total;
  }, [menuItem.price, quantity, selectedAddons]);

  const handleQuantityChange = (newQuantity: number) => {
    updateQuantity(cartItemId, newQuantity);
  };

  const handleRemove = () => {
    removeFromCart(cartItemId);
  };

  return (
    <div className="border-b pb-3 pt-2">
      <div className="flex gap-2">
        {/* Compact image */}
        <div className="w-16 h-16 flex-shrink-0">
          <img
            src={menuItem.image}
            alt={menuItem.name}
            className="w-full h-full object-cover rounded-md"
          />
        </div>

        {/* Main item info */}
        <div className="flex-grow min-w-0">
          <div className="flex justify-between items-start">
            <h3 className="font-medium text-sm truncate">{menuItem.name}</h3>
            <p className="font-semibold text-food-primary text-sm whitespace-nowrap">
              ${itemTotalPrice.toFixed(2)}
            </p>
          </div>

          {/* Quantity controls in compact row */}
          <div className="flex items-center justify-between mt-1">
            <div className="flex items-center border rounded-md">
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 rounded-none"
                onClick={() => handleQuantityChange(quantity - 1)}
                disabled={quantity <= 1}
              >
                <Minus size={12} />
              </Button>
              <Input
                type="number"
                value={quantity}
                onChange={(e) => {
                  const val = parseInt(e.target.value);
                  if (!isNaN(val) && val > 0) {
                    handleQuantityChange(val);
                  }
                }}
                className="w-8 h-7 text-center border-0 focus:ring-0 p-0 text-sm"
                min={1}
              />
              <Button
                variant="ghost"
                size="icon"
                className="h-7 w-7 rounded-none"
                onClick={() => handleQuantityChange(quantity + 1)}
              >
                <Plus size={12} />
              </Button>
            </div>

            <Button
              variant="ghost"
              size="sm"
              onClick={handleRemove}
              className="text-red-500 hover:text-red-700 hover:bg-red-50 h-7 px-2 text-xs"
            >
              <Trash2 size={14} className="mr-1" />
              Remove
            </Button>
          </div>

          {/* Selected Addons Display */}
          {selectedAddons && selectedAddons.length > 0 && (
            <div className="mt-1 text-xs text-gray-600">
              {selectedAddons.map(({ addon, quantity: addonQty }) => (
                <div key={addon.id} className="flex justify-between">
                  <span>+ {addon.name} x{addonQty}</span>
                  <span className="text-food-primary">+${(addon.price * addonQty).toFixed(2)}</span>
                </div>
              ))}
            </div>
          )}

          {/* Special Instructions Display */}
          {item.specialInstructions && (
            <div className="mt-1 text-xs text-gray-500 italic border-t pt-1">
              Note: {item.specialInstructions}
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default CartItemComponent;

