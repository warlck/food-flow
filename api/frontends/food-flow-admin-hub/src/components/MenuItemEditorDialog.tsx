import React, { FormEvent, useEffect, useState } from 'react';
import {
  ChevronDown,
  ChevronUp,
  GripVertical,
  Loader2,
  Plus,
  Puzzle,
  Trash2,
  UtensilsCrossed,
  Layers,
  Check,
  AlertCircle,
} from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import { ImageField } from '@/components/ImageField';
import {
  AdminCategory,
  AdminMenuItem,
  AdminMenuItemAddonInfo,
  AdminModifierGroup,
  AdminModifierOption,
  AdminWorkspace,
  MenuItemInput,
  ModifierGroupInput,
  ModifierOptionInput,
  adminApi,
} from '@/lib/admin-api';

interface MenuItemEditorDialogProps {
  item?: AdminMenuItem;
  defaultCategoryId?: string;
  workspace: AdminWorkspace;
  open: boolean;
  onClose: () => void;
  onSaveItem: (input: MenuItemInput, existingId?: string) => Promise<AdminMenuItem>;
  onRefreshWorkspace: () => Promise<void>;
}

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('en-US', { style: 'currency', currency: 'USD' }).format(amount);
}

export function MenuItemEditorDialog({
  item,
  defaultCategoryId,
  workspace,
  open,
  onClose,
  onSaveItem,
  onRefreshWorkspace,
}: MenuItemEditorDialogProps) {
  const [activeTab, setActiveTab] = useState<'details' | 'modifiers' | 'addons'>('details');
  const [savingDetails, setSavingDetails] = useState(false);
  const [detailsError, setDetailsError] = useState('');

  // Current item ID (if saved or already existing)
  const [currentItem, setCurrentItem] = useState<AdminMenuItem | undefined>(item);

  // Modifiers state
  const [modifierGroups, setModifierGroups] = useState<AdminModifierGroup[]>([]);
  const [loadingGroups, setLoadingGroups] = useState(false);
  const [editingGroup, setEditingGroup] = useState<AdminModifierGroup | null>(null);
  const [isAddingGroup, setIsAddingGroup] = useState(false);
  const [editingOption, setEditingOption] = useState<{ group: AdminModifierGroup; option?: AdminModifierOption } | null>(null);

  // Addons state
  const [assignedAddons, setAssignedAddons] = useState<AdminMenuItemAddonInfo[]>([]);
  const [loadingAddons, setLoadingAddons] = useState(false);
  const [savingAddons, setSavingAddons] = useState(false);

  useEffect(() => {
    setCurrentItem(item);
    setActiveTab('details');
    setDetailsError('');
  }, [item, open]);

  // Load modifier groups and options when switching to modifiers tab
  useEffect(() => {
    if (activeTab === 'modifiers' && currentItem?.id) {
      loadModifierGroups(currentItem.id);
    } else if (activeTab === 'addons' && currentItem?.id) {
      loadAssignedAddons(currentItem.id);
    }
  }, [activeTab, currentItem?.id]);

  const loadModifierGroups = async (itemId: string) => {
    setLoadingGroups(true);
    try {
      const groupPage = await adminApi.listModifierGroups(itemId);
      const groupsWithOptions = await Promise.all(
        groupPage.items.map(async (group) => {
          const optPage = await adminApi.listModifierOptions(group.id);
          return { ...group, options: optPage.items };
        })
      );
      setModifierGroups(groupsWithOptions);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to load modifier groups');
    } finally {
      setLoadingGroups(false);
    }
  };

  const loadAssignedAddons = async (itemId: string) => {
    setLoadingAddons(true);
    try {
      const addons = await adminApi.getMenuItemAddons(itemId);
      setAssignedAddons(addons);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to load assigned add-ons');
    } finally {
      setLoadingAddons(false);
    }
  };

  const handleSaveDetails = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSavingDetails(true);
    setDetailsError('');

    const data = new FormData(e.currentTarget);
    const rankStr = String(data.get('rank') ?? '').trim();
    const rankVal = rankStr !== '' ? Number(rankStr) : undefined;

    const input: MenuItemInput = {
      name: String(data.get('name')),
      description: String(data.get('description')),
      price: Number(data.get('price')),
      categoryId: String(data.get('categoryId')),
      restaurantId: workspace.restaurant.id,
      imageUrl: String(data.get('imageUrl') ?? ''),
      rank: rankVal,
    };

    try {
      const saved = await onSaveItem(input, currentItem?.id);
      setCurrentItem(saved);
      toast.success(currentItem?.id ? 'Menu item updated' : 'Menu item created');
      await onRefreshWorkspace();
    } catch (err) {
      setDetailsError(err instanceof Error ? err.message : 'Failed to save menu item details');
    } finally {
      setSavingDetails(false);
    }
  };

  // Modifier Group Operations
  const handleSaveGroup = async (groupInput: { name: string; description: string; minSelections: number; maxSelections: number; available: boolean }, existingGroupId?: string) => {
    if (!currentItem?.id) return;

    if (groupInput.minSelections > groupInput.maxSelections) {
      toast.error('Minimum selections cannot exceed maximum selections');
      return;
    }

    try {
      if (existingGroupId) {
        await adminApi.updateModifierGroup(existingGroupId, groupInput);
        toast.success('Modifier group updated');
      } else {
        await adminApi.createModifierGroup({
          ...groupInput,
          menuItemId: currentItem.id,
          restaurantId: workspace.restaurant.id,
        });
        toast.success('Modifier group created');
      }
      setIsAddingGroup(false);
      setEditingGroup(null);
      await loadModifierGroups(currentItem.id);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to save modifier group');
    }
  };

  const handleDeleteGroup = async (group: AdminModifierGroup) => {
    if (!currentItem?.id) return;
    if (!window.confirm(`Delete modifier group "${group.name}" and all its options? This cannot be undone.`)) return;

    try {
      await adminApi.deleteModifierGroup(group.id);
      toast.success('Modifier group deleted');
      await loadModifierGroups(currentItem.id);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete modifier group');
    }
  };

  const handleToggleGroupAvailable = async (group: AdminModifierGroup, available: boolean) => {
    if (!currentItem?.id) return;
    if (available && group.minSelections > 0) {
      const hasAvailableOption = (group.options ?? []).some((o) => o.available);
      if (!hasAvailableOption) {
        toast.error('Cannot enable a required group without at least one available option');
        return;
      }
    }

    try {
      await adminApi.updateModifierGroup(group.id, { available });
      toast.success(available ? 'Group enabled' : 'Group disabled');
      await loadModifierGroups(currentItem.id);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update group availability');
    }
  };

  const handleMoveGroup = async (index: number, direction: 'up' | 'down') => {
    if (!currentItem?.id) return;
    const targetIndex = direction === 'up' ? index - 1 : index + 1;
    if (targetIndex < 0 || targetIndex >= modifierGroups.length) return;

    const reordered = [...modifierGroups];
    const [moved] = reordered.splice(index, 1);
    reordered.splice(targetIndex, 0, moved);

    const orderedIds = reordered.map((g) => g.id);
    try {
      const updated = await adminApi.reorderModifierGroups({ menuItemId: currentItem.id, orderedIds });
      const optMap = new Map(modifierGroups.map((g) => [g.id, g.options]));
      setModifierGroups(updated.map((g) => ({ ...g, options: optMap.get(g.id) ?? [] })));
      toast.success('Modifier groups reordered');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to reorder modifier groups');
      await loadModifierGroups(currentItem.id);
    }
  };

  // Modifier Option Operations
  const handleSaveOption = async (optionInput: { name: string; description: string; priceDelta: number; available: boolean }, group: AdminModifierGroup, existingOptionId?: string) => {
    if (!currentItem?.id) return;

    try {
      if (existingOptionId) {
        await adminApi.updateModifierOption(existingOptionId, optionInput);
        toast.success('Option updated');
      } else {
        await adminApi.createModifierOption({
          ...optionInput,
          modifierGroupId: group.id,
          restaurantId: workspace.restaurant.id,
        });
        toast.success('Option created');
      }
      setEditingOption(null);
      await loadModifierGroups(currentItem.id);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to save modifier option');
    }
  };

  const handleDeleteOption = async (group: AdminModifierGroup, option: AdminModifierOption) => {
    if (!currentItem?.id) return;
    if (group.available && group.minSelections > 0) {
      const remainingAvailable = (group.options ?? []).filter((o) => o.available && o.id !== option.id);
      if (remainingAvailable.length === 0) {
        toast.error('Cannot delete the last available option of an active required group. Disable the group or menu item first.');
        return;
      }
    }

    if (!window.confirm(`Delete option "${option.name}"? This cannot be undone.`)) return;

    try {
      await adminApi.deleteModifierOption(option.id);
      toast.success('Option deleted');
      await loadModifierGroups(currentItem.id);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to delete modifier option');
    }
  };

  const handleToggleOptionAvailable = async (group: AdminModifierGroup, option: AdminModifierOption, available: boolean) => {
    if (!currentItem?.id) return;
    if (!available && group.available && group.minSelections > 0) {
      const remainingAvailable = (group.options ?? []).filter((o) => o.available && o.id !== option.id);
      if (remainingAvailable.length === 0) {
        toast.error('Cannot disable the last available option of an active required group. Disable the group or menu item first.');
        return;
      }
    }

    try {
      await adminApi.updateModifierOption(option.id, { available });
      toast.success(available ? 'Option enabled' : 'Option disabled');
      await loadModifierGroups(currentItem.id);
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update option availability');
    }
  };

  const handleMoveOption = async (group: AdminModifierGroup, optIndex: number, direction: 'up' | 'down') => {
    if (!currentItem?.id || !group.options) return;
    const targetIndex = direction === 'up' ? optIndex - 1 : optIndex + 1;
    if (targetIndex < 0 || targetIndex >= group.options.length) return;

    const reordered = [...group.options];
    const [moved] = reordered.splice(optIndex, 1);
    reordered.splice(targetIndex, 0, moved);

    const orderedIds = reordered.map((o) => o.id);
    try {
      const updatedOpts = await adminApi.reorderModifierOptions({ modifierGroupId: group.id, orderedIds });
      setModifierGroups((prev) =>
        prev.map((g) => (g.id === group.id ? { ...g, options: updatedOpts } : g))
      );
      toast.success('Options reordered');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to reorder options');
      await loadModifierGroups(currentItem.id);
    }
  };

  // Addon Assignment Operations
  const handleToggleAddonAssignment = async (addonId: string, assigned: boolean) => {
    if (!currentItem?.id) return;

    let updatedList: AdminMenuItemAddonInfo[];
    if (assigned) {
      const fullAddon = workspace.addons.find((a) => a.id === addonId);
      if (!fullAddon) return;
      const newRank = (assignedAddons.length + 1) * 10;
      updatedList = [
        ...assignedAddons,
        {
          id: addonId,
          addonId: fullAddon.id,
          name: fullAddon.name,
          description: fullAddon.description,
          price: fullAddon.price,
          available: fullAddon.available,
          maxQuantity: fullAddon.maxQuantity,
          rank: newRank,
        },
      ];
    } else {
      updatedList = assignedAddons.filter((a) => a.addonId !== addonId);
    }

    setAssignedAddons(updatedList);
    setSavingAddons(true);
    try {
      const saved = await adminApi.replaceMenuItemAddons(currentItem.id, {
        addons: updatedList.map((a, idx) => ({ addonId: a.addonId, rank: (idx + 1) * 10 })),
      });
      setAssignedAddons(saved);
      toast.success(assigned ? 'Add-on assigned' : 'Add-on removed');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to update add-on assignments');
      await loadAssignedAddons(currentItem.id);
    } finally {
      setSavingAddons(false);
    }
  };

  const handleMoveAssignedAddon = async (index: number, direction: 'up' | 'down') => {
    if (!currentItem?.id) return;
    const targetIndex = direction === 'up' ? index - 1 : index + 1;
    if (targetIndex < 0 || targetIndex >= assignedAddons.length) return;

    const reordered = [...assignedAddons];
    const [moved] = reordered.splice(index, 1);
    reordered.splice(targetIndex, 0, moved);

    const orderedIds = reordered.map((a) => a.addonId);
    setAssignedAddons(reordered);
    setSavingAddons(true);
    try {
      const saved = await adminApi.reorderMenuItemAddons(currentItem.id, { orderedIds });
      setAssignedAddons(saved);
      toast.success('Add-ons reordered');
    } catch (err) {
      toast.error(err instanceof Error ? err.message : 'Failed to reorder add-ons');
      await loadAssignedAddons(currentItem.id);
    } finally {
      setSavingAddons(false);
    }
  };

  const isNewItem = !currentItem?.id;

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="flex h-[90vh] max-h-[92vh] w-[95vw] max-w-4xl flex-col overflow-hidden rounded-2xl border-[#E5E7EB] bg-[#F9FAFB] p-0 shadow-2xl">
        {/* Header */}
        <div className="flex shrink-0 items-center justify-between border-b border-[#E5E7EB] bg-white px-6 py-4">
          <div className="flex items-center gap-3">
            <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#FFF1EB] text-[#FF4500]">
              <UtensilsCrossed size={20} />
            </div>
            <div>
              <h2 className="text-lg font-bold tracking-[-.025em] text-[#111827]">
                {currentItem?.name ? `Edit: ${currentItem.name}` : 'New Menu Item'}
              </h2>
              <p className="text-xs text-[#6B7280]">
                Configure item details, required/optional modifier choices, and add-on assignments.
              </p>
            </div>
          </div>

          {/* Navigation Tabs */}
          <div className="flex items-center rounded-xl bg-[#F3F4F6] p-1">
            <button
              type="button"
              onClick={() => setActiveTab('details')}
              className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-semibold transition-all ${
                activeTab === 'details'
                  ? 'bg-white text-[#111827] shadow-sm'
                  : 'text-[#6B7280] hover:text-[#111827]'
              }`}
            >
              <UtensilsCrossed size={14} />
              Details
            </button>
            <button
              type="button"
              onClick={() => {
                if (isNewItem) {
                  toast.error('Please save the menu item details before adding modifiers');
                  return;
                }
                setActiveTab('modifiers');
              }}
              disabled={isNewItem}
              className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-semibold transition-all ${
                isNewItem ? 'cursor-not-allowed opacity-50 text-[#9CA3AF]' : ''
              } ${
                activeTab === 'modifiers'
                  ? 'bg-white text-[#111827] shadow-sm'
                  : 'text-[#6B7280] hover:text-[#111827]'
              }`}
            >
              <Layers size={14} />
              Modifiers ({modifierGroups.length})
            </button>
            <button
              type="button"
              onClick={() => {
                if (isNewItem) {
                  toast.error('Please save the menu item details before assigning add-ons');
                  return;
                }
                setActiveTab('addons');
              }}
              disabled={isNewItem}
              className={`flex items-center gap-1.5 rounded-lg px-3 py-1.5 text-xs font-semibold transition-all ${
                isNewItem ? 'cursor-not-allowed opacity-50 text-[#9CA3AF]' : ''
              } ${
                activeTab === 'addons'
                  ? 'bg-white text-[#111827] shadow-sm'
                  : 'text-[#6B7280] hover:text-[#111827]'
              }`}
            >
              <Puzzle size={14} />
              Add-ons ({assignedAddons.length})
            </button>
          </div>
        </div>

        {/* Tab Content */}
        <div className="flex-1 overflow-y-auto p-6">
          {activeTab === 'details' && (
            <form onSubmit={handleSaveDetails} id="item-details-form" className="space-y-5">
              <div className="rounded-2xl border border-[#E5E7EB] bg-white p-5 shadow-sm space-y-4">
                <div className="flex items-center gap-2 border-b border-[#F3F4F6] pb-3 text-xs font-bold uppercase tracking-wider text-[#FF4500]">
                  <UtensilsCrossed size={14} />
                  <span>Item Information</span>
                </div>

                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-1.5">
                    <Label htmlFor="name" className="text-xs font-semibold text-[#374151]">
                      Item Name <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id="name"
                      name="name"
                      defaultValue={currentItem?.name ?? ''}
                      required
                      minLength={3}
                      maxLength={100}
                      placeholder="e.g. Truffle Cheeseburger"
                      className="admin-input"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="categoryId" className="text-xs font-semibold text-[#374151]">
                      Category <span className="text-red-500">*</span>
                    </Label>
                    <Select
                      name="categoryId"
                      defaultValue={currentItem?.categoryId ?? defaultCategoryId ?? workspace.categories[0]?.id}
                    >
                      <SelectTrigger className="admin-input">
                        <SelectValue placeholder="Select category" />
                      </SelectTrigger>
                      <SelectContent>
                        {workspace.categories.map((c) => (
                          <SelectItem key={c.id} value={c.id}>
                            {c.name}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <div className="space-y-1.5">
                  <Label htmlFor="description" className="text-xs font-semibold text-[#374151]">
                    Description
                  </Label>
                  <Textarea
                    id="description"
                    name="description"
                    defaultValue={currentItem?.description ?? ''}
                    rows={2}
                    placeholder="Short appetising description of the dish"
                    className="admin-input resize-none text-xs"
                  />
                </div>

                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-1.5">
                    <Label htmlFor="price" className="text-xs font-semibold text-[#374151]">
                      Base Price ($) <span className="text-red-500">*</span>
                    </Label>
                    <div className="relative">
                      <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-[#6B7280]">$</span>
                      <Input
                        id="price"
                        name="price"
                        type="number"
                        min="0.01"
                        step="0.01"
                        defaultValue={currentItem?.price ?? ''}
                        required
                        placeholder="0.00"
                        className="admin-input pl-7"
                      />
                    </div>
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="rank" className="text-xs font-semibold text-[#374151]">
                      Display Rank (Order)
                    </Label>
                    <Input
                      id="rank"
                      name="rank"
                      type="number"
                      min="1"
                      step="1"
                      defaultValue={currentItem?.rank ?? ''}
                      placeholder="e.g. 10 (lower appears first)"
                      className="admin-input"
                    />
                  </div>
                </div>

                <div className="space-y-1.5 pt-2">
                  <Label htmlFor="imageUrl" className="text-xs font-semibold text-[#374151]">
                    Dish Photo
                  </Label>
                  <ImageField
                    name="imageUrl"
                    entityType="menu_item"
                    restaurantId={workspace.restaurant.id}
                    defaultValue={currentItem?.imageUrl ?? ''}
                  />
                </div>
              </div>

              {detailsError && (
                <div className="rounded-xl border border-red-200 bg-red-50 p-3 text-xs font-medium text-red-700">
                  {detailsError}
                </div>
              )}
            </form>
          )}

          {activeTab === 'modifiers' && (
            <div className="space-y-5">
              <div className="flex items-center justify-between">
                <div>
                  <h3 className="text-sm font-bold text-[#111827]">Modifier Groups</h3>
                  <p className="text-xs text-[#6B7280]">
                    Create required options (e.g. Choose Size) or optional choices (e.g. Extra Sauce).
                  </p>
                </div>
                <Button
                  type="button"
                  size="sm"
                  onClick={() => {
                    setIsAddingGroup(true);
                    setEditingGroup(null);
                  }}
                  className="admin-primary h-8 gap-1.5 text-xs font-semibold"
                >
                  <Plus size={14} /> Add Modifier Group
                </Button>
              </div>

              {loadingGroups ? (
                <div className="flex h-40 items-center justify-center">
                  <Loader2 className="h-6 w-6 animate-spin text-[#FF4500]" />
                </div>
              ) : modifierGroups.length === 0 ? (
                <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-[#D1D5DB] bg-white p-8 text-center">
                  <Layers className="h-8 w-8 text-[#9CA3AF] mb-2" />
                  <p className="text-xs font-bold text-[#374151]">No modifier groups configured</p>
                  <p className="text-[11px] text-[#6B7280] max-w-xs mt-1">
                    Add choices like proteins, cook temperature, or sizes for customers to select when ordering.
                  </p>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setIsAddingGroup(true)}
                    className="mt-4 h-8 gap-1.5 text-xs text-[#FF4500] border-[#FF4500]/30 hover:bg-[#FFF1EB]"
                  >
                    <Plus size={13} /> Add first group
                  </Button>
                </div>
              ) : (
                <div className="space-y-4">
                  {modifierGroups.map((group, groupIdx) => {
                    const isRequired = group.minSelections > 0;
                    return (
                      <div
                        key={group.id}
                        className="rounded-2xl border border-[#E5E7EB] bg-white p-4 shadow-sm space-y-3"
                      >
                        {/* Group Header */}
                        <div className="flex items-center justify-between border-b border-[#F3F4F6] pb-3">
                          <div className="flex items-center gap-2">
                            <div className="flex flex-col">
                              <button
                                type="button"
                                disabled={groupIdx === 0}
                                onClick={() => handleMoveGroup(groupIdx, 'up')}
                                className="p-0.5 text-[#9CA3AF] hover:text-[#111827] disabled:opacity-30"
                              >
                                <ChevronUp size={14} />
                              </button>
                              <button
                                type="button"
                                disabled={groupIdx === modifierGroups.length - 1}
                                onClick={() => handleMoveGroup(groupIdx, 'down')}
                                className="p-0.5 text-[#9CA3AF] hover:text-[#111827] disabled:opacity-30"
                              >
                                <ChevronDown size={14} />
                              </button>
                            </div>

                            <div>
                              <div className="flex items-center gap-2">
                                <h4 className="text-sm font-bold text-[#111827]">{group.name}</h4>
                                <span
                                  className={`rounded-full px-2 py-0.5 text-[10px] font-bold uppercase tracking-wider ${
                                    isRequired
                                      ? 'bg-[#FFF1EB] text-[#FF4500]'
                                      : 'bg-[#F3F4F6] text-[#6B7280]'
                                  }`}
                                >
                                  {isRequired ? `Required (Min ${group.minSelections}, Max ${group.maxSelections})` : `Optional (Max ${group.maxSelections})`}
                                </span>
                                {!group.available && (
                                  <span className="rounded-full bg-[#FFEBEE] px-2 py-0.5 text-[10px] font-bold uppercase text-[#C62828]">
                                    Disabled
                                  </span>
                                )}
                              </div>
                              {group.description && (
                                <p className="text-xs text-[#6B7280]">{group.description}</p>
                              )}
                            </div>
                          </div>

                          <div className="flex items-center gap-3">
                            <div className="flex items-center gap-1.5">
                              <Label className="text-[11px] text-[#6B7280]">Available</Label>
                              <Switch
                                checked={group.available}
                                onCheckedChange={(val) => handleToggleGroupAvailable(group, val)}
                              />
                            </div>
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={() => setEditingGroup(group)}
                              className="h-7 px-2 text-xs text-[#6B7280] hover:text-[#111827]"
                            >
                              Edit
                            </Button>
                            <Button
                              type="button"
                              variant="ghost"
                              size="sm"
                              onClick={() => handleDeleteGroup(group)}
                              className="h-7 px-2 text-xs text-[#EF4444] hover:bg-red-50"
                            >
                              <Trash2 size={13} />
                            </Button>
                          </div>
                        </div>

                        {/* Options List */}
                        <div className="pl-6 space-y-2">
                          <div className="flex items-center justify-between">
                            <span className="text-[11px] font-bold uppercase tracking-wider text-[#9CA3AF]">
                              Options ({group.options?.length ?? 0})
                            </span>
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              onClick={() => setEditingOption({ group })}
                              className="h-6 px-2 text-[11px] font-semibold text-[#FF4500] border-[#FF4500]/30 hover:bg-[#FFF1EB]"
                            >
                              <Plus size={12} /> Add Option
                            </Button>
                          </div>

                          {!group.options || group.options.length === 0 ? (
                            <p className="text-xs italic text-[#9CA3AF] py-1">
                              No options in this group yet. Add at least one option.
                            </p>
                          ) : (
                            <div className="space-y-1.5">
                              {group.options.map((opt, optIdx) => (
                                <div
                                  key={opt.id}
                                  className="flex items-center justify-between rounded-xl bg-[#F9FAFB] px-3 py-2 border border-[#F3F4F6]"
                                >
                                  <div className="flex items-center gap-2">
                                    <div className="flex flex-col">
                                      <button
                                        type="button"
                                        disabled={optIdx === 0}
                                        onClick={() => handleMoveOption(group, optIdx, 'up')}
                                        className="p-0.5 text-[#9CA3AF] hover:text-[#111827] disabled:opacity-30"
                                      >
                                        <ChevronUp size={12} />
                                      </button>
                                      <button
                                        type="button"
                                        disabled={optIdx === (group.options?.length ?? 1) - 1}
                                        onClick={() => handleMoveOption(group, optIdx, 'down')}
                                        className="p-0.5 text-[#9CA3AF] hover:text-[#111827] disabled:opacity-30"
                                      >
                                        <ChevronDown size={12} />
                                      </button>
                                    </div>
                                    <div>
                                      <span className="text-xs font-semibold text-[#111827]">{opt.name}</span>
                                      {opt.description && (
                                        <span className="ml-2 text-[11px] text-[#6B7280]">
                                          ({opt.description})
                                        </span>
                                      )}
                                    </div>
                                  </div>

                                  <div className="flex items-center gap-3">
                                    <span className="font-mono text-xs font-semibold text-[#374151]">
                                      {opt.priceDelta === 0 ? '+$0.00' : `+${formatCurrency(opt.priceDelta)}`}
                                    </span>
                                    <div className="flex items-center gap-1.5">
                                      <Switch
                                        checked={opt.available}
                                        onCheckedChange={(val) => handleToggleOptionAvailable(group, opt, val)}
                                      />
                                    </div>
                                    <Button
                                      type="button"
                                      variant="ghost"
                                      size="sm"
                                      onClick={() => setEditingOption({ group, option: opt })}
                                      className="h-6 px-1.5 text-[11px] text-[#6B7280] hover:text-[#111827]"
                                    >
                                      Edit
                                    </Button>
                                    <Button
                                      type="button"
                                      variant="ghost"
                                      size="sm"
                                      onClick={() => handleDeleteOption(group, opt)}
                                      className="h-6 px-1.5 text-[11px] text-[#EF4444] hover:bg-red-50"
                                    >
                                      <Trash2 size={12} />
                                    </Button>
                                  </div>
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                      </div>
                    );
                  })}
                </div>
              )}
            </div>
          )}

          {activeTab === 'addons' && (
            <div className="space-y-5">
              <div>
                <h3 className="text-sm font-bold text-[#111827]">Assigned Add-ons</h3>
                <p className="text-xs text-[#6B7280]">
                  Select which restaurant add-ons can be selected as extras with this dish, and set their ordering.
                </p>
              </div>

              {loadingAddons ? (
                <div className="flex h-40 items-center justify-center">
                  <Loader2 className="h-6 w-6 animate-spin text-[#FF4500]" />
                </div>
              ) : workspace.addons.length === 0 ? (
                <div className="flex flex-col items-center justify-center rounded-2xl border border-dashed border-[#D1D5DB] bg-white p-8 text-center">
                  <Puzzle className="h-8 w-8 text-[#9CA3AF] mb-2" />
                  <p className="text-xs font-bold text-[#374151]">No restaurant add-on definitions created yet</p>
                  <p className="text-[11px] text-[#6B7280] max-w-xs mt-1">
                    Create global add-on definitions under the Menu &gt; Add-ons tab first, then assign them here.
                  </p>
                </div>
              ) : (
                <div className="space-y-3">
                  <div className="rounded-2xl border border-[#E5E7EB] bg-white p-4 shadow-sm space-y-2">
                    <span className="text-[11px] font-bold uppercase tracking-wider text-[#9CA3AF]">
                      Available Restaurant Add-ons
                    </span>

                    <div className="space-y-2">
                      {workspace.addons.map((addon) => {
                        const assignedIdx = assignedAddons.findIndex((a) => a.addonId === addon.id);
                        const isAssigned = assignedIdx !== -1;

                        return (
                          <div
                            key={addon.id}
                            className={`flex items-center justify-between rounded-xl px-3 py-2.5 border transition-all ${
                              isAssigned
                                ? 'border-[#FF8C42]/50 bg-[#FFF7F3]'
                                : 'border-[#F3F4F6] bg-[#F9FAFB]'
                            }`}
                          >
                            <div className="flex items-center gap-3">
                              <input
                                type="checkbox"
                                id={`assign-addon-${addon.id}`}
                                checked={isAssigned}
                                disabled={savingAddons}
                                onChange={(e) => handleToggleAddonAssignment(addon.id, e.target.checked)}
                                className="h-4 w-4 rounded border-gray-300 text-[#FF4500] focus:ring-[#FF4500]"
                              />
                              <div>
                                <label
                                  htmlFor={`assign-addon-${addon.id}`}
                                  className="text-xs font-bold text-[#111827] cursor-pointer"
                                >
                                  {addon.name}
                                </label>
                                {addon.description && (
                                  <p className="text-[11px] text-[#6B7280]">{addon.description}</p>
                                )}
                              </div>
                            </div>

                            <div className="flex items-center gap-3">
                              <span className="font-mono text-xs font-semibold text-[#374151]">
                                +{formatCurrency(addon.price)} (Max {addon.maxQuantity})
                              </span>

                              {isAssigned && (
                                <div className="flex items-center gap-1 border-l border-[#E5E7EB] pl-3">
                                  <button
                                    type="button"
                                    disabled={assignedIdx === 0 || savingAddons}
                                    onClick={() => handleMoveAssignedAddon(assignedIdx, 'up')}
                                    className="p-1 text-[#9CA3AF] hover:text-[#111827] disabled:opacity-30"
                                    title="Move Up"
                                  >
                                    <ChevronUp size={14} />
                                  </button>
                                  <button
                                    type="button"
                                    disabled={assignedIdx === assignedAddons.length - 1 || savingAddons}
                                    onClick={() => handleMoveAssignedAddon(assignedIdx, 'down')}
                                    className="p-1 text-[#9CA3AF] hover:text-[#111827] disabled:opacity-30"
                                    title="Move Down"
                                  >
                                    <ChevronDown size={14} />
                                  </button>
                                </div>
                              )}
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Footer */}
        <div className="flex shrink-0 items-center justify-end gap-2 border-t border-[#E5E7EB] bg-white px-6 py-4">
          <Button type="button" variant="outline" onClick={onClose} className="h-9 px-4 text-xs font-semibold">
            Close
          </Button>
          {activeTab === 'details' && (
            <Button
              type="submit"
              form="item-details-form"
              disabled={savingDetails}
              className="admin-primary h-9 min-w-[120px] gap-2 text-xs font-semibold"
            >
              {savingDetails && <Loader2 size={14} className="animate-spin" />}
              {currentItem?.id ? 'Save Changes' : 'Create & Continue'}
            </Button>
          )}
        </div>

        {/* Group Add/Edit Modal */}
        {(isAddingGroup || editingGroup) && (
          <Dialog open onOpenChange={() => { setIsAddingGroup(false); setEditingGroup(null); }}>
            <DialogContent className="sm:max-w-[480px]">
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  const form = new FormData(e.currentTarget);
                  handleSaveGroup(
                    {
                      name: String(form.get('name')),
                      description: String(form.get('description')),
                      minSelections: Number(form.get('minSelections')),
                      maxSelections: Number(form.get('maxSelections')),
                      available: form.get('available') === 'true',
                    },
                    editingGroup?.id
                  );
                }}
              >
                <DialogHeader>
                  <DialogTitle>{editingGroup ? 'Edit Modifier Group' : 'Add Modifier Group'}</DialogTitle>
                  <DialogDescription>
                    Configure group name, min/max selection bounds, and availability.
                  </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4">
                  <div className="space-y-1.5">
                    <Label htmlFor="grp-name" className="text-xs font-semibold">Group Name *</Label>
                    <Input
                      id="grp-name"
                      name="name"
                      defaultValue={editingGroup?.name ?? ''}
                      required
                      placeholder="e.g. Choose Size or Select Protein"
                      className="admin-input"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="grp-desc" className="text-xs font-semibold">Description</Label>
                    <Input
                      id="grp-desc"
                      name="description"
                      defaultValue={editingGroup?.description ?? ''}
                      placeholder="Optional guidance for guests"
                      className="admin-input"
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div className="space-y-1.5">
                      <Label htmlFor="grp-min" className="text-xs font-semibold">Min Selections</Label>
                      <Input
                        id="grp-min"
                        name="minSelections"
                        type="number"
                        min="0"
                        max="10"
                        defaultValue={editingGroup?.minSelections ?? 1}
                        required
                        className="admin-input"
                      />
                      <span className="text-[10px] text-[#6B7280]">0 = optional, 1+ = required</span>
                    </div>

                    <div className="space-y-1.5">
                      <Label htmlFor="grp-max" className="text-xs font-semibold">Max Selections</Label>
                      <Input
                        id="grp-max"
                        name="maxSelections"
                        type="number"
                        min="1"
                        max="10"
                        defaultValue={editingGroup?.maxSelections ?? 1}
                        required
                        className="admin-input"
                      />
                    </div>
                  </div>

                  <div className="flex items-center justify-between rounded-xl bg-[#F9FAFB] p-3 border border-[#F3F4F6]">
                    <div>
                      <Label className="text-xs font-semibold">Group Available</Label>
                      <p className="text-[11px] text-[#6B7280]">Allow selection on customer storefront</p>
                    </div>
                    <select
                      name="available"
                      defaultValue={String(editingGroup?.available ?? true)}
                      className="h-8 rounded-lg border border-[#E5E7EB] bg-white px-2 text-xs font-semibold"
                    >
                      <option value="true">Yes (Enabled)</option>
                      <option value="false">No (Disabled)</option>
                    </select>
                  </div>
                </div>

                <div className="flex justify-end gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => { setIsAddingGroup(false); setEditingGroup(null); }}
                  >
                    Cancel
                  </Button>
                  <Button type="submit" className="admin-primary">
                    Save Group
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
        )}

        {/* Option Add/Edit Modal */}
        {editingOption && (
          <Dialog open onOpenChange={() => setEditingOption(null)}>
            <DialogContent className="sm:max-w-[440px]">
              <form
                onSubmit={(e) => {
                  e.preventDefault();
                  const form = new FormData(e.currentTarget);
                  handleSaveOption(
                    {
                      name: String(form.get('name')),
                      description: String(form.get('description')),
                      priceDelta: Number(form.get('priceDelta')),
                      available: form.get('available') === 'true',
                    },
                    editingOption.group,
                    editingOption.option?.id
                  );
                }}
              >
                <DialogHeader>
                  <DialogTitle>{editingOption.option ? 'Edit Option' : `Add Option to ${editingOption.group.name}`}</DialogTitle>
                  <DialogDescription>
                    Configure option name, price adjustment, and availability.
                  </DialogDescription>
                </DialogHeader>

                <div className="space-y-4 py-4">
                  <div className="space-y-1.5">
                    <Label htmlFor="opt-name" className="text-xs font-semibold">Option Name *</Label>
                    <Input
                      id="opt-name"
                      name="name"
                      defaultValue={editingOption.option?.name ?? ''}
                      required
                      placeholder="e.g. Regular, Large, or Spicy"
                      className="admin-input"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="opt-desc" className="text-xs font-semibold">Description</Label>
                    <Input
                      id="opt-desc"
                      name="description"
                      defaultValue={editingOption.option?.description ?? ''}
                      placeholder="Optional details"
                      className="admin-input"
                    />
                  </div>

                  <div className="space-y-1.5">
                    <Label htmlFor="opt-price" className="text-xs font-semibold">Price Delta ($)</Label>
                    <div className="relative">
                      <span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-[#6B7280]">$</span>
                      <Input
                        id="opt-price"
                        name="priceDelta"
                        type="number"
                        min="0.00"
                        step="0.01"
                        defaultValue={editingOption.option?.priceDelta ?? 0}
                        required
                        placeholder="0.00"
                        className="admin-input pl-7"
                      />
                    </div>
                    <span className="text-[10px] text-[#6B7280]">Enter 0.00 for no price change</span>
                  </div>

                  <div className="flex items-center justify-between rounded-xl bg-[#F9FAFB] p-3 border border-[#F3F4F6]">
                    <div>
                      <Label className="text-xs font-semibold">Option Available</Label>
                      <p className="text-[11px] text-[#6B7280]">In-stock for customer orders</p>
                    </div>
                    <select
                      name="available"
                      defaultValue={String(editingOption.option?.available ?? true)}
                      className="h-8 rounded-lg border border-[#E5E7EB] bg-white px-2 text-xs font-semibold"
                    >
                      <option value="true">Yes (Available)</option>
                      <option value="false">No (Sold Out)</option>
                    </select>
                  </div>
                </div>

                <div className="flex justify-end gap-2">
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setEditingOption(null)}
                  >
                    Cancel
                  </Button>
                  <Button type="submit" className="admin-primary">
                    Save Option
                  </Button>
                </div>
              </form>
            </DialogContent>
          </Dialog>
        )}
      </DialogContent>
    </Dialog>
  );
}
