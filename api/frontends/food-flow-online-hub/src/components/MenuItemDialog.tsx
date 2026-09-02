import React, { useMemo, useState, useEffect } from 'react';
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
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Checkbox } from '@/components/ui/checkbox';
import { MenuItem as MenuItemType, Addon, SelectedAddon, ModifierGroup, ModifierOption, SelectedModifier } from '@/types';
import { useCart } from '@/context/CartContext';
import { Plus, Minus, ShoppingCart, AlertCircle, Check } from 'lucide-react';
import { Separator } from '@/components/ui/separator';
import { Textarea } from '@/components/ui/textarea';

interface MenuItemDialogProps {
  item: MenuItemType | null;
  categoryItems?: MenuItemType[];
  isOpen: boolean;
  onClose: () => void;
}

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount);
}

const MenuItemDialog: React.FC<MenuItemDialogProps> = ({
  item,
  categoryItems,
  isOpen,
  onClose,
}) => {
  const { addToCart } = useCart();
  const [quantity, setQuantity] = useState(1);
  const [selectedModifiers, setSelectedModifiers] = useState<Record<string, string[]>>({}); // groupId -> array of optionIds
  const [addonQuantities, setAddonQuantities] = useState<Record<string, number>>({});
  const [specialInstructions, setSpecialInstructions] = useState('');
  const [imageError, setImageError] = useState(false);

  const activeItem = item;

  // Reset state when dialog opens
  useEffect(() => {
    if (isOpen && activeItem) {
      setQuantity(1);
      setAddonQuantities({});
      setSpecialInstructions('');
      setImageError(false);

      // Pre-select default modifier options if a single required choice exists
      const initialMods: Record<string, string[]> = {};
      if (activeItem.modifierGroups) {
        activeItem.modifierGroups.forEach((group) => {
          if (group.available && group.minSelections === 1 && group.maxSelections === 1) {
            const firstAvailable = group.options.find((o) => o.available);
            if (firstAvailable) {
              initialMods[group.id] = [firstAvailable.id];
            }
          }
        });
      }
      setSelectedModifiers(initialMods);
    }
  }, [isOpen, activeItem]);

  // Handle single-select modifier option (Radio)
  const handleSelectRadioOption = (groupId: string, optionId: string) => {
    setSelectedModifiers((prev) => ({
      ...prev,
      [groupId]: [optionId],
    }));
  };

  // Handle multi-select modifier option (Checkbox)
  const handleToggleCheckboxOption = (group: ModifierGroup, option: ModifierOption, checked: boolean) => {
    setSelectedModifiers((prev) => {
      const current = prev[group.id] || [];
      if (checked) {
        if (current.length < group.maxSelections && !current.includes(option.id)) {
          return { ...prev, [group.id]: [...current, option.id] };
        }
      } else {
        return { ...prev, [group.id]: current.filter((id) => id !== option.id) };
      }
      return prev;
    });
  };

  // Handle addon increment / decrement
  const handleAddonIncrement = (addon: Addon) => {
    setAddonQuantities((prev) => {
      const current = prev[addon.id] || 0;
      const max = addon.maxQuantity && addon.maxQuantity > 0 ? addon.maxQuantity : 10;
      if (current < max) {
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

  // Compute structured selected modifiers
  const selectedModifiersList: SelectedModifier[] = useMemo(() => {
    if (!activeItem?.modifierGroups) return [];
    const list: SelectedModifier[] = [];

    activeItem.modifierGroups.forEach((group) => {
      const chosenOptionIds = selectedModifiers[group.id] || [];
      chosenOptionIds.forEach((optId) => {
        const opt = group.options.find((o) => o.id === optId);
        if (opt) {
          list.push({
            modifierGroupId: group.id,
            modifierGroupName: group.name,
            modifierOptionId: opt.id,
            modifierOptionName: opt.name,
            priceDelta: opt.priceDelta,
          });
        }
      });
    });

    return list;
  }, [activeItem?.modifierGroups, selectedModifiers]);

  // Compute structured selected addons
  const selectedAddonsList: SelectedAddon[] = useMemo(() => {
    if (!activeItem?.addons) return [];
    return activeItem.addons
      .filter((addon) => (addonQuantities[addon.id] || 0) > 0)
      .map((addon) => ({
        addon,
        quantity: addonQuantities[addon.id],
      }));
  }, [activeItem?.addons, addonQuantities]);

  // Validate modifier requirements
  const modifierValidation = useMemo(() => {
    if (!activeItem?.modifierGroups) return { isValid: true, errors: [] as string[] };
    const errors: string[] = [];

    activeItem.modifierGroups.forEach((group) => {
      if (!group.available) return;
      const count = (selectedModifiers[group.id] || []).length;
      if (group.minSelections > 0 && count < group.minSelections) {
        errors.push(`Please select at least ${group.minSelections} option${group.minSelections > 1 ? 's' : ''} for "${group.name}".`);
      } else if (count > group.maxSelections) {
        errors.push(`You can select at most ${group.maxSelections} option${group.maxSelections > 1 ? 's' : ''} for "${group.name}".`);
      }
    });

    return {
      isValid: errors.length === 0,
      errors,
    };
  }, [activeItem?.modifierGroups, selectedModifiers]);

  // Calculate unit price and total
  const unitPrice = useMemo(() => {
    if (!activeItem) return 0;
    let price = activeItem.price;
    selectedModifiersList.forEach((mod) => {
      price += mod.priceDelta;
    });
    selectedAddonsList.forEach(({ addon, quantity: addonQty }) => {
      price += addon.price * addonQty;
    });
    return price;
  }, [activeItem, selectedModifiersList, selectedAddonsList]);

  const totalPrice = unitPrice * quantity;

  const isOrderable = activeItem?.orderable ?? activeItem?.available ?? false;
  const canAddToCart = isOrderable && modifierValidation.isValid;

  const handleAddToCart = () => {
    if (!activeItem || !canAddToCart) return;
    addToCart(
      activeItem,
      quantity,
      selectedModifiersList.length > 0 ? selectedModifiersList : undefined,
      selectedAddonsList.length > 0 ? selectedAddonsList : undefined,
      specialInstructions
    );
    onClose();
  };

  const getCategoryImage = () => {
    if (!item) return '';
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
      categoryImageMap[item.category] ||
      'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&q=80'
    );
  };

  if (!item || !activeItem) return null;

  return (
    <Dialog open={isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="sm:max-w-[540px] max-h-[90vh] overflow-y-auto bg-white text-gray-900 p-0 rounded-2xl">
        <div className="p-6">
          <DialogHeader className="text-left">
            <div className="relative w-full aspect-video rounded-xl overflow-hidden mb-4 bg-muted">
              <img
                src={
                  imageError
                    ? 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&q=80'
                    : activeItem?.image || getCategoryImage()
                }
                alt={activeItem.name}
                className="w-full h-full object-cover"
                onError={() => setImageError(true)}
              />
              <div className="absolute top-2 right-2">
                <Badge className="bg-food-primary text-white">{activeItem.category}</Badge>
              </div>
            </div>

            <div className="flex items-start justify-between gap-3">
              <div>
                <DialogTitle className="text-xl font-bold text-gray-900">{activeItem.name}</DialogTitle>
                {activeItem.description && (
                  <DialogDescription className="text-sm text-gray-600 mt-1">
                    {activeItem.description}
                  </DialogDescription>
                )}
              </div>
              <div className="text-lg font-bold text-food-primary shrink-0">
                {formatCurrency(activeItem.price)}
              </div>
            </div>

            {!isOrderable && (
              <div className="mt-3 flex items-center gap-2 rounded-xl bg-red-50 p-3 text-xs font-semibold text-red-700 border border-red-200">
                <AlertCircle size={16} className="shrink-0" />
                <span>This dish is currently sold out or has unavailable required options.</span>
              </div>
            )}
          </DialogHeader>

          {/* Modifier Groups Section - unavailable groups are suspended:
              not presented, no requirement, no selections */}
          {activeItem.modifierGroups && activeItem.modifierGroups.some((group) => group.available) && (
            <div className="mt-6 space-y-6">
              {activeItem.modifierGroups.filter((group) => group.available).map((group) => {
                const isSingleSelect = group.maxSelections === 1;
                const isRequired = group.minSelections > 0;
                const currentSelections = selectedModifiers[group.id] || [];

                return (
                  <div key={group.id} className="rounded-xl border border-gray-200 bg-gray-50/50 p-4 space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <div className="flex items-center gap-2">
                          <h4 className="font-bold text-sm text-gray-900">{group.name}</h4>
                          <span
                            className={`rounded-full px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider ${
                              isRequired
                                ? 'bg-food-primary/10 text-food-primary'
                                : 'bg-gray-200 text-gray-700'
                            }`}
                          >
                            {isRequired
                              ? isSingleSelect
                                ? 'Required - Choose 1'
                                : `Required - Choose ${group.minSelections}${group.maxSelections > group.minSelections ? ` to ${group.maxSelections}` : ''}`
                              : `Optional (up to ${group.maxSelections})`}
                          </span>
                        </div>
                        {group.description && (
                          <p className="text-xs text-gray-500 mt-0.5">{group.description}</p>
                        )}
                      </div>
                    </div>

                    {isSingleSelect ? (
                      <RadioGroup
                        value={currentSelections[0] || ''}
                        onValueChange={(val) => handleSelectRadioOption(group.id, val)}
                        className="space-y-2"
                      >
                        {group.options.map((opt) => {
                          const id = `opt-${opt.id}`;
                          const isSelected = currentSelections.includes(opt.id);
                          return (
                            <div
                              key={opt.id}
                              className={`flex items-center justify-between rounded-lg p-3 transition-colors ${
                                isSelected
                                  ? 'border border-food-primary bg-food-primary/5'
                                  : 'bg-white border border-gray-200'
                              } ${!opt.available ? 'opacity-50 cursor-not-allowed' : 'cursor-pointer'}`}
                              onClick={() => {
                                if (opt.available) handleSelectRadioOption(group.id, opt.id);
                              }}
                            >
                              <div className="flex items-center gap-3">
                                <RadioGroupItem
                                  id={id}
                                  value={opt.id}
                                  disabled={!opt.available}
                                />
                                <div>
                                  <Label
                                    htmlFor={id}
                                    className={`font-medium text-sm text-gray-900 ${
                                      !opt.available ? 'cursor-not-allowed' : 'cursor-pointer'
                                    }`}
                                  >
                                    {opt.name}
                                  </Label>
                                  {opt.description && (
                                    <p className="text-xs text-gray-500">{opt.description}</p>
                                  )}
                                </div>
                              </div>

                              <div className="text-xs font-semibold text-gray-700">
                                {!opt.available ? (
                                  <span className="text-red-500 font-normal">Sold out</span>
                                ) : opt.priceDelta === 0 ? (
                                  <span className="text-gray-400">Free</span>
                                ) : (
                                  <span className="text-food-primary">+{formatCurrency(opt.priceDelta)}</span>
                                )}
                              </div>
                            </div>
                          );
                        })}
                      </RadioGroup>
                    ) : (
                      <div className="space-y-2">
                        {group.options.map((opt) => {
                          const id = `opt-${opt.id}`;
                          const isSelected = currentSelections.includes(opt.id);
                          const reachedMax = currentSelections.length >= group.maxSelections && !isSelected;

                          return (
                            <div
                              key={opt.id}
                              className={`flex items-center justify-between rounded-lg p-3 transition-colors ${
                                isSelected
                                  ? 'border border-food-primary bg-food-primary/5'
                                  : 'bg-white border border-gray-200'
                              } ${!opt.available || (reachedMax && !isSelected) ? 'opacity-50' : 'cursor-pointer'}`}
                              onClick={() => {
                                if (opt.available && (!reachedMax || isSelected)) {
                                  handleToggleCheckboxOption(group, opt, !isSelected);
                                }
                              }}
                            >
                              <div className="flex items-center gap-3">
                                <Checkbox
                                  id={id}
                                  checked={isSelected}
                                  disabled={!opt.available || (reachedMax && !isSelected)}
                                  onCheckedChange={(val) => handleToggleCheckboxOption(group, opt, Boolean(val))}
                                />
                                <div>
                                  <Label
                                    htmlFor={id}
                                    className={`font-medium text-sm text-gray-900 ${
                                      !opt.available ? 'cursor-not-allowed' : 'cursor-pointer'
                                    }`}
                                  >
                                    {opt.name}
                                  </Label>
                                  {opt.description && (
                                    <p className="text-xs text-gray-500">{opt.description}</p>
                                  )}
                                </div>
                              </div>

                              <div className="text-xs font-semibold text-gray-700">
                                {!opt.available ? (
                                  <span className="text-red-500 font-normal">Sold out</span>
                                ) : opt.priceDelta === 0 ? (
                                  <span className="text-gray-400">Free</span>
                                ) : (
                                  <span className="text-food-primary">+{formatCurrency(opt.priceDelta)}</span>
                                )}
                              </div>
                            </div>
                          );
                        })}
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}

          {/* Add-ons Section */}
          {activeItem.addons && activeItem.addons.length > 0 && (
            <div className="mt-6 space-y-3">
              <Separator />
              <h3 className="font-bold text-sm text-gray-900 pt-2">Add-ons (Optional)</h3>
              <div className="space-y-2">
                {activeItem.addons.map((addon) => {
                  const max = addon.maxQuantity && addon.maxQuantity > 0 ? addon.maxQuantity : 10;
                  const currentQty = addonQuantities[addon.id] || 0;

                  return (
                    <div
                      key={addon.id}
                      className="flex items-center justify-between p-3 bg-gray-50 rounded-xl border border-gray-200"
                    >
                      <div className="flex-1">
                        <div className="font-medium text-sm text-gray-900">{addon.name}</div>
                        {addon.description && (
                          <div className="text-xs text-gray-500">{addon.description}</div>
                        )}
                        <div className="text-xs font-semibold text-food-primary mt-0.5">
                          {!addon.available ? (
                            <span className="text-red-500 font-normal">Sold out</span>
                          ) : (
                            `+${formatCurrency(addon.price)} (Max ${max})`
                          )}
                        </div>
                      </div>

                      <div className="flex items-center gap-2">
                        <Button
                          type="button"
                          variant="outline"
                          size="icon"
                          className="h-8 w-8 rounded-full text-gray-800"
                          onClick={() => handleAddonDecrement(addon.id)}
                          disabled={!addon.available || currentQty <= 0}
                        >
                          <Minus className="h-4 w-4" />
                        </Button>
                        <span className="w-6 text-center font-semibold text-sm text-gray-900">
                          {currentQty}
                        </span>
                        <Button
                          type="button"
                          variant="outline"
                          size="icon"
                          className="h-8 w-8 rounded-full text-gray-800"
                          onClick={() => handleAddonIncrement(addon)}
                          disabled={!addon.available || currentQty >= max}
                        >
                          <Plus className="h-4 w-4" />
                        </Button>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>
          )}

          {/* Special Instructions */}
          <div className="mt-6 space-y-2">
            <Separator />
            <h3 className="font-bold text-sm text-gray-900 pt-2">Special Instructions</h3>
            <Textarea
              placeholder="Add a note for the kitchen (e.g. dressing on the side)..."
              value={specialInstructions}
              onChange={(e) => setSpecialInstructions(e.target.value)}
              className="resize-none text-xs"
              rows={2}
            />
          </div>

          {/* Quantity Section */}
          <div className="mt-6 space-y-2">
            <Separator />
            <div className="flex items-center justify-between pt-2">
              <span className="font-bold text-sm text-gray-900">Quantity</span>
              <div className="flex items-center gap-3">
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  className="h-8 w-8 rounded-full text-gray-800"
                  onClick={() => setQuantity((q) => Math.max(1, q - 1))}
                  disabled={quantity <= 1}
                >
                  <Minus className="h-4 w-4" />
                </Button>
                <span className="w-8 text-center font-bold text-base text-gray-900">{quantity}</span>
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  className="h-8 w-8 rounded-full text-gray-800"
                  onClick={() => setQuantity((q) => q + 1)}
                >
                  <Plus className="h-4 w-4" />
                </Button>
              </div>
            </div>
          </div>

          {/* Validation Warnings */}
          {modifierValidation.errors.length > 0 && (
            <div className="mt-4 space-y-1">
              {modifierValidation.errors.map((err, i) => (
                <p key={i} className="text-xs font-semibold text-red-600 flex items-center gap-1.5">
                  <AlertCircle size={13} /> {err}
                </p>
              ))}
            </div>
          )}

          <DialogFooter className="mt-6">
            <Button
              onClick={handleAddToCart}
              disabled={!canAddToCart}
              className="w-full bg-food-primary hover:bg-food-accent text-white py-6 text-base font-bold rounded-xl shadow-lg shadow-food-primary/20"
            >
              <ShoppingCart className="mr-2 h-5 w-5" />
              Add to Cart - {formatCurrency(totalPrice)}
            </Button>
          </DialogFooter>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default MenuItemDialog;
