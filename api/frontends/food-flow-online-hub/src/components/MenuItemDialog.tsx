import React, { useMemo, useState } from 'react';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { MenuItem as MenuItemType, Addon, SelectedAddon } from '@/types';
import { useCart } from '@/context/CartContext';
import { Plus, Minus, Tag, ShoppingCart } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Textarea } from '@/components/ui/textarea';

interface MenuItemDialogProps {
  item: MenuItemType | null;
  categoryItems?: MenuItemType[];
  isOpen: boolean;
  onClose: () => void;
}

const MenuItemDialog: React.FC<MenuItemDialogProps> = ({
  item,
  categoryItems,
  isOpen,
  onClose,
}) => {
  const { addToCart } = useCart();
  const [quantity, setQuantity] = useState(1);
  const [addonQuantities, setAddonQuantities] = useState<Record<string, number>>({});
  const [specialInstructions, setSpecialInstructions] = useState('');
  const [imageError, setImageError] = useState(false);
  const [selectedItemId, setSelectedItemId] = useState<string | null>(null);

  const itemsInCategory = useMemo(() => {
    if (!item) return [];

    const list = categoryItems && categoryItems.length > 0 ? categoryItems : [item];

    // Ensure the triggering item is always present.
    if (!list.some((mi) => mi.id === item.id)) {
      return [item, ...list];
    }

    return list;
  }, [categoryItems, item]);

  // Backend sends category items sorted by price (cheapest first).
  const cheapestItem = useMemo(() => {
    if (itemsInCategory.length === 0) return null;
    return itemsInCategory[0];
  }, [itemsInCategory]);

  const activeItem = useMemo(() => {
    if (!item) return null;

    if (!selectedItemId) {
      return item;
    }

    return itemsInCategory.find((mi) => mi.id === selectedItemId) ?? item;
  }, [item, itemsInCategory, selectedItemId]);

  // Reset state when dialog opens with new item
  React.useEffect(() => {
    if (isOpen && item) {
      setQuantity(1);
      setAddonQuantities({});
      setSpecialInstructions('');
      setImageError(false);
      setSelectedItemId(item.id);
    }
  }, [isOpen, item]);

  const handleSelectMenuItem = (menuItemId: string) => {
    setSelectedItemId(menuItemId);
    // Addons are menu-item specific, so reset when switching.
    setAddonQuantities({});
    setImageError(false);
  };

  const handleAddonIncrement = (addon: Addon) => {
    setAddonQuantities((prev) => {
      const current = prev[addon.id] || 0;
      if (current < addon.maxQuantity) {
        return { ...prev, [addon.id]: current + 1 };
      }
      return prev;
    });
  };

  const handleAddonDecrement = (addonId: string) => {
    setAddonQuantities((prev) => {
      const current = prev[addonId] || 0;
      if (current > 0) {
        return { ...prev, [addonId]: current - 1 };
      }
      return prev;
    });
  };

  const selectedAddons: SelectedAddon[] = useMemo(() => {
    if (!activeItem?.addons) return [];
    return activeItem.addons
      .filter((addon) => (addonQuantities[addon.id] || 0) > 0)
      .map((addon) => ({
        addon,
        quantity: addonQuantities[addon.id],
      }));
  }, [activeItem?.addons, addonQuantities]);

  const totalPrice = useMemo(() => {
    if (!activeItem) return 0;
    let total = activeItem.price * quantity;
    selectedAddons.forEach(({ addon, quantity: addonQty }) => {
      total += addon.price * addonQty * quantity;
    });
    return total;
  }, [activeItem, quantity, selectedAddons]);

  const handleAddToCart = () => {
    if (!activeItem) return;
    addToCart(
      activeItem,
      quantity,
      selectedAddons.length > 0 ? selectedAddons : undefined,
      specialInstructions
    );
    onClose();
  };

  const getFallbackImage = () => {
    if (!activeItem) return '';
    const categoryImageMap: Record<string, string> = {
      Appetizers: 'https://images.unsplash.com/photo-1546241072-48010ad2862c?auto=format&fit=crop&q=80',
      'Main Course': 'https://images.unsplash.com/photo-1574484284002-952d92456975?auto=format&fit=crop&q=80',
      Desserts: 'https://images.unsplash.com/photo-1563729784474-d77dbb933a9e?auto=format&fit=crop&q=80',
      Beverages: 'https://images.unsplash.com/photo-1544145945-f90425340c7e?auto=format&fit=crop&q=80',
      Sides: 'https://images.unsplash.com/photo-1573080496219-bb080dd4f877?auto=format&fit=crop&q=80',
      Pizza: 'https://images.unsplash.com/photo-1565299624946-b28f40a0ae38?auto=format&fit=crop&q=80',
      Burgers: 'https://images.unsplash.com/photo-1568901346375-23c9450c58cd?auto=format&fit=crop&q=80',
      Pasta: 'https://images.unsplash.com/photo-1473093226795-af9932fe5856?auto=format&fit=crop&q=80',
    };
    return (
      categoryImageMap[activeItem.category] ||
      'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&q=80'
    );
  };

  if (!item || !activeItem) return null;

  const priceDeltaVsCheapest =
    cheapestItem && activeItem.price > cheapestItem.price
      ? activeItem.price - cheapestItem.price
      : 0;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[500px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <div className="relative w-full aspect-video rounded-lg overflow-hidden mb-4 bg-muted">
            <img
              src={imageError ? getFallbackImage() : activeItem.image}
              alt={activeItem.name}
              className="w-full h-full object-cover"
              onError={() => setImageError(true)}
            />
            {!activeItem.available && (
              <div className="absolute inset-0 bg-black/60 flex items-center justify-center">
                <span className="text-white font-semibold text-lg">Out of Stock</span>
              </div>
            )}
            <div className="absolute top-2 right-2">
              <Badge className="bg-food-primary text-white">{activeItem.category}</Badge>
            </div>
          </div>
          <DialogTitle className="text-xl font-bold">{activeItem.name}</DialogTitle>
          <DialogDescription className="text-base text-gray-600 mt-2">
            {activeItem.description}
          </DialogDescription>
          {activeItem.tags && activeItem.tags.length > 0 && (
            <div className="flex flex-wrap gap-1 mt-2">
              {activeItem.tags.map((tag) => (
                <Badge key={tag} variant="secondary" className="text-xs flex items-center">
                  <Tag size={10} className="mr-1" />
                  {tag}
                </Badge>
              ))}
            </div>
          )}

          {itemsInCategory.length > 1 && cheapestItem && (
            <div className="mt-4">
              <Label htmlFor="menu-item-select" className="mb-2 block">
                Choose an item
              </Label>
              <Select value={activeItem.id} onValueChange={handleSelectMenuItem}>
                <SelectTrigger id="menu-item-select">
                  <SelectValue placeholder="Select an item" />
                </SelectTrigger>
                <SelectContent>
                  {itemsInCategory.map((mi) => {
                    const delta = mi.price - cheapestItem.price;
                    const deltaLabel = delta > 0 ? ` (+$${delta.toFixed(2)})` : '';
                    const stockLabel = !mi.available ? ' (out of stock)' : '';
                    return (
                      <SelectItem key={mi.id} value={mi.id} disabled={!mi.available}>
                        {mi.name} — ${mi.price.toFixed(2)}{deltaLabel}{stockLabel}
                      </SelectItem>
                    );
                  })}
                </SelectContent>
              </Select>
            </div>
          )}

          <div className="text-xl font-bold text-food-primary mt-4">
            ${activeItem.price.toFixed(2)}
          </div>
          {priceDeltaVsCheapest > 0 && (
            <div className="text-sm text-gray-500">
              +${priceDeltaVsCheapest.toFixed(2)} vs. cheapest option
            </div>
          )}
        </DialogHeader>

        {/* Addons Section */}
        {activeItem.addons && activeItem.addons.length > 0 && (
          <div className="mt-4">
            <Separator className="mb-4" />
            <h3 className="font-semibold text-lg mb-3">Add-ons (Optional)</h3>
            <div className="space-y-3">
              {activeItem.addons.map((addon) => (
                <div
                  key={addon.id}
                  className="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
                >
                  <div className="flex-1">
                    <div className="font-medium">{addon.name}</div>
                    {addon.description && (
                      <div className="text-sm text-gray-500">{addon.description}</div>
                    )}
                    <div className="text-sm font-semibold text-food-primary">
                      +${addon.price.toFixed(2)}
                    </div>
                  </div>
                  <div className="flex items-center gap-2">
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8 rounded-full"
                      onClick={() => handleAddonDecrement(addon.id)}
                      disabled={!addonQuantities[addon.id]}
                    >
                      <Minus className="h-4 w-4" />
                    </Button>
                    <span className="w-8 text-center font-medium">
                      {addonQuantities[addon.id] || 0}
                    </span>
                    <Button
                      variant="outline"
                      size="icon"
                      className="h-8 w-8 rounded-full"
                      onClick={() => handleAddonIncrement(addon)}
                      disabled={(addonQuantities[addon.id] || 0) >= addon.maxQuantity}
                    >
                      <Plus className="h-4 w-4" />
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* Special Instructions Section */}
        <div className="mt-4">
          <Separator className="mb-4" />
          <h3 className="font-semibold text-lg mb-3">Special Instructions</h3>
          <Textarea
            placeholder="Add a note for the kitchen (e.g. no onions, extra spicy)..."
            value={specialInstructions}
            onChange={(e) => setSpecialInstructions(e.target.value)}
            className="resize-none"
          />
        </div>

        {/* Quantity Section */}
        <div className="mt-4">
          <Separator className="mb-4" />
          <div className="flex items-center justify-between">
            <span className="font-semibold">Quantity</span>
            <div className="flex items-center gap-3">
              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8 rounded-full"
                onClick={() => setQuantity(q => Math.max(1, q - 1))}
                disabled={quantity <= 1}
              >
                <Minus className="h-4 w-4" />
              </Button>
              <span className="w-8 text-center font-medium text-lg">{quantity}</span>
              <Button
                variant="outline"
                size="icon"
                className="h-8 w-8 rounded-full"
                onClick={() => setQuantity(q => q + 1)}
              >
                <Plus className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>

        <DialogFooter className="mt-6">
          <Button
            onClick={handleAddToCart}
            disabled={!activeItem.available}
            className="w-full bg-food-primary hover:bg-food-accent text-white py-6 text-lg"
          >
            <ShoppingCart className="mr-2 h-5 w-5" />
            Add to Cart - ${totalPrice.toFixed(2)}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

export default MenuItemDialog;
