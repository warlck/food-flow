import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { beforeEach, describe, expect, it, vi } from 'vitest';
import { MenuItemEditorDialog } from './MenuItemEditorDialog';
import type {
  AdminAddon,
  AdminCategory,
  AdminMenuItem,
  AdminModifierGroup,
  AdminModifierOption,
  AdminWorkspace,
  MenuItemInput,
} from '@/lib/admin-api';

const apiMocks = vi.hoisted(() => ({
  listModifierGroups: vi.fn(),
  listModifierOptions: vi.fn(),
  createModifierGroup: vi.fn(),
  updateModifierGroup: vi.fn(),
  deleteModifierGroup: vi.fn(),
  reorderModifierGroups: vi.fn(),
  createModifierOption: vi.fn(),
  updateModifierOption: vi.fn(),
  deleteModifierOption: vi.fn(),
  reorderModifierOptions: vi.fn(),
  listAddons: vi.fn(),
  createAddon: vi.fn(),
  updateAddon: vi.fn(),
  deleteAddon: vi.fn(),
  reorderAddons: vi.fn(),
}));

const toastMocks = vi.hoisted(() => ({
  error: vi.fn(),
  success: vi.fn(),
}));

vi.mock('@/lib/admin-api', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@/lib/admin-api')>()),
  adminApi: apiMocks,
}));

vi.mock('sonner', () => ({ toast: toastMocks }));

vi.mock('@/components/ImageField', () => ({
  ImageField: ({ name, defaultValue }: { name: string; defaultValue?: string }) => (
    <input type="hidden" name={name} value={defaultValue ?? ''} readOnly />
  ),
}));

vi.mock('@/components/ui/dialog', () => ({
  Dialog: ({ open, children }: React.PropsWithChildren<{ open?: boolean }>) => open ? <>{children}</> : null,
  DialogContent: ({ children }: React.PropsWithChildren) => <div>{children}</div>,
  DialogDescription: ({ children }: React.PropsWithChildren) => <p>{children}</p>,
  DialogHeader: ({ children }: React.PropsWithChildren) => <header>{children}</header>,
  DialogTitle: ({ children }: React.PropsWithChildren) => <h2>{children}</h2>,
}));

vi.mock('@/components/ui/select', () => ({
  Select: ({
    children,
    name,
    defaultValue,
  }: React.PropsWithChildren<{ name?: string; defaultValue?: string }>) => (
    <select name={name} defaultValue={defaultValue}>{children}</select>
  ),
  SelectContent: ({ children }: React.PropsWithChildren) => <>{children}</>,
  SelectItem: ({ children, value }: React.PropsWithChildren<{ value: string }>) => (
    <option value={value}>{children}</option>
  ),
  SelectTrigger: ({ children }: React.PropsWithChildren) => <>{children}</>,
  SelectValue: () => null,
}));

vi.mock('@/components/ui/switch', () => ({
  Switch: ({
    checked,
    onCheckedChange,
    ...props
  }: React.ButtonHTMLAttributes<HTMLButtonElement> & {
    checked: boolean;
    onCheckedChange: (checked: boolean) => void;
  }) => (
    <button
      type="button"
      role="switch"
      aria-checked={checked}
      onClick={() => onCheckedChange(!checked)}
      {...props}
    />
  ),
}));

const restaurant = {
  id: 'restaurant-1',
  name: 'Test Restaurant',
  description: '',
  address: '',
  phone: '',
  email: 'test@example.com',
  imageUrl: '',
  enabled: true,
  maxDeliveryDistanceKm: 10,
  taxRate: 0.1,
  dateCreated: '',
  dateUpdated: '',
};

const categories: AdminCategory[] = [
  {
    id: 'category-enabled',
    name: 'Enabled Category',
    description: '',
    restaurantId: restaurant.id,
    enabled: true,
  },
  {
    id: 'category-disabled',
    name: 'Disabled Category',
    description: '',
    restaurantId: restaurant.id,
    enabled: false,
  },
];

const item: AdminMenuItem = {
  id: 'item-1',
  name: 'Test Item',
  description: 'Test description',
  price: 12,
  categoryId: categories[0].id,
  restaurantId: restaurant.id,
  imageUrl: '',
  available: true,
};

const option: AdminModifierOption = {
  id: 'option-1',
  modifierGroupId: 'group-1',
  restaurantId: restaurant.id,
  name: 'Mild',
  description: '',
  priceDelta: 0.5,
  available: true,
};

const group: AdminModifierGroup = {
  id: 'group-1',
  menuItemId: item.id,
  restaurantId: restaurant.id,
  name: 'Spice Level',
  description: '',
  minSelections: 0,
  maxSelections: 1,
  available: true,
  options: [option],
};

const addons: AdminAddon[] = [
  {
    id: 'addon-1',
    menuItemId: item.id,
    restaurantId: restaurant.id,
    name: 'First Addon',
    description: '',
    price: 1,
    available: true,
    maxQuantity: 2,
    rank: 10,
  },
  {
    id: 'addon-2',
    menuItemId: item.id,
    restaurantId: restaurant.id,
    name: 'Second Addon',
    description: '',
    price: 2,
    available: true,
    maxQuantity: 3,
    rank: 20,
  },
];

const page = <T,>(items: T[]) => ({
  items,
  total: items.length,
  page: 1,
  rowsPerPage: 100,
});

function workspaceFor(menuItem: AdminMenuItem = item): AdminWorkspace {
  return {
    restaurant,
    categories,
    menuItems: [menuItem],
    addons: [],
  };
}

function renderEditor({
  currentItem = item,
  workspace = workspaceFor(currentItem ?? item),
  onSaveItem = vi.fn().mockResolvedValue(currentItem ?? item),
  onRefreshWorkspace = vi.fn().mockResolvedValue(undefined),
}: {
  currentItem?: AdminMenuItem | null;
  workspace?: AdminWorkspace;
  onSaveItem?: (input: MenuItemInput, existingId?: string) => Promise<AdminMenuItem>;
  onRefreshWorkspace?: () => Promise<void>;
} = {}) {
  const view = render(
    <MenuItemEditorDialog
      item={currentItem ?? undefined}
      workspace={workspace}
      open
      onClose={vi.fn()}
      onSaveItem={onSaveItem}
      onRefreshWorkspace={onRefreshWorkspace}
    />,
  );
  return { ...view, onSaveItem, onRefreshWorkspace };
}

describe('MenuItemEditorDialog', () => {
  beforeEach(() => {
    Object.values(apiMocks).forEach((mock) => mock.mockReset());
    Object.values(toastMocks).forEach((mock) => mock.mockReset());
    apiMocks.listModifierGroups.mockResolvedValue(page([]));
    apiMocks.listModifierOptions.mockResolvedValue(page([]));
    apiMocks.listAddons.mockResolvedValue(page([]));
    apiMocks.createModifierGroup.mockResolvedValue(group);
    apiMocks.updateModifierGroup.mockResolvedValue(group);
    apiMocks.deleteModifierGroup.mockResolvedValue(undefined);
    apiMocks.reorderModifierGroups.mockResolvedValue([]);
    apiMocks.createModifierOption.mockResolvedValue(option);
    apiMocks.updateModifierOption.mockResolvedValue(option);
    apiMocks.deleteModifierOption.mockResolvedValue(undefined);
    apiMocks.reorderModifierOptions.mockResolvedValue([]);
    apiMocks.createAddon.mockResolvedValue(addons[0]);
    apiMocks.updateAddon.mockResolvedValue(addons[0]);
    apiMocks.deleteAddon.mockResolvedValue(undefined);
    apiMocks.reorderAddons.mockResolvedValue(addons);
    vi.spyOn(window, 'confirm').mockReturnValue(true);
  });

  it('saves a new item before enabling nested editors', async () => {
    const saved = { ...item, id: 'saved-item' };
    const onSaveItem = vi.fn().mockResolvedValue(saved);
    renderEditor({
      currentItem: null,
      workspace: workspaceFor(),
      onSaveItem,
    });

    const modifiersTab = screen.getByRole('button', { name: 'Modifiers (0)' });
    const addonsTab = screen.getByRole('button', { name: 'Add-ons (0)' });
    expect(modifiersTab).toBeDisabled();
    expect(addonsTab).toBeDisabled();

    fireEvent.change(screen.getByLabelText(/Item Name/), { target: { value: 'Saved Item' } });
    fireEvent.change(screen.getByLabelText(/Base Price/), { target: { value: '12.00' } });
    fireEvent.click(screen.getByRole('button', { name: 'Create & Continue' }));

    await waitFor(() => expect(onSaveItem).toHaveBeenCalledTimes(1));
    expect(modifiersTab).not.toBeDisabled();
    expect(addonsTab).not.toBeDisabled();

    fireEvent.click(modifiersTab);
    await waitFor(() => expect(apiMocks.listModifierGroups).toHaveBeenCalledWith(saved.id));
    fireEvent.click(addonsTab);
    await waitFor(() => {
      expect(apiMocks.listAddons).toHaveBeenCalledWith(restaurant.id, saved.id);
    });
  });

  it.each([
    ['enabled', categories[0]],
    ['disabled', categories[1]],
  ])('loads item-scoped add-ons for an item in an %s category', async (_state, category) => {
    const categoryItem = { ...item, categoryId: category.id };
    apiMocks.listAddons.mockResolvedValue(page(addons));
    renderEditor({
      currentItem: categoryItem,
      workspace: workspaceFor(categoryItem),
    });

    fireEvent.click(screen.getByRole('button', { name: 'Add-ons (0)' }));

    await waitFor(() => expect(screen.getByText('First Addon')).toBeInTheDocument());
    expect(apiMocks.listAddons).toHaveBeenCalledWith(restaurant.id, categoryItem.id);
  });

  it('covers modifier CRUD and preserves a zero-price option', async () => {
    apiMocks.listModifierGroups.mockResolvedValue(page([group]));
    apiMocks.listModifierOptions.mockResolvedValue(page([option]));
    renderEditor();

    fireEvent.click(screen.getByRole('button', { name: 'Modifiers (0)' }));
    await waitFor(() => expect(screen.getByText('Spice Level')).toBeInTheDocument());

    fireEvent.click(screen.getByRole('button', { name: 'Add Modifier Group' }));
    expect(screen.queryByLabelText('Group Available')).not.toBeInTheDocument();
    expect(screen.getByText(/New modifier groups are disabled by default/)).toBeInTheDocument();

    fireEvent.change(screen.getByLabelText('Group Name *'), { target: { value: 'Sauce Choice' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save Group' }));
    await waitFor(() => {
      expect(apiMocks.createModifierGroup).toHaveBeenCalledWith({
        name: 'Sauce Choice',
        description: '',
        minSelections: 1,
        maxSelections: 1,
        available: false,
        menuItemId: item.id,
        restaurantId: restaurant.id,
      });
    });

    fireEvent.click(screen.getByRole('button', { name: 'Add Option' }));
    fireEvent.change(screen.getByLabelText('Option Name *'), { target: { value: 'No Charge' } });
    fireEvent.change(screen.getByLabelText('Price Delta ($)'), { target: { value: '0' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save Option' }));
    await waitFor(() => {
      expect(apiMocks.createModifierOption).toHaveBeenCalledWith(expect.objectContaining({
        name: 'No Charge',
        priceDelta: 0,
        modifierGroupId: group.id,
        restaurantId: restaurant.id,
      }));
    });

    fireEvent.click(screen.getByRole('switch', { name: 'Toggle Mild availability' }));
    await waitFor(() => {
      expect(apiMocks.updateModifierOption).toHaveBeenCalledWith(option.id, { available: false });
    });

    fireEvent.click(screen.getByRole('button', { name: 'Delete Mild' }));
    await waitFor(() => expect(apiMocks.deleteModifierOption).toHaveBeenCalledWith(option.id));
    fireEvent.click(screen.getByRole('button', { name: 'Delete Spice Level' }));
    await waitFor(() => expect(apiMocks.deleteModifierGroup).toHaveBeenCalledWith(group.id));
  });

  it('disables availability toggle and displays No Options badge for groups with no options', async () => {
    const emptyGroup: AdminModifierGroup = {
      ...group,
      id: 'empty-group',
      name: 'Empty Group',
      available: false,
      options: [],
    };
    apiMocks.listModifierGroups.mockResolvedValue(page([emptyGroup]));
    apiMocks.listModifierOptions.mockResolvedValue(page([]));
    renderEditor();

    fireEvent.click(screen.getByRole('button', { name: 'Modifiers (0)' }));
    await waitFor(() => expect(screen.getByText('Empty Group')).toBeInTheDocument());

    expect(screen.getByText('No Options')).toBeInTheDocument();
    const switchBtn = screen.getByRole('switch', { name: 'Toggle Empty Group availability' });
    expect(switchBtn).toBeDisabled();

    // Clicking edit on an empty group does not provide an enable selector
    fireEvent.click(screen.getByRole('button', { name: 'Edit Empty Group' }));
    expect(screen.queryByLabelText('Group Available')).not.toBeInTheDocument();
    expect(screen.getByText(/This group has no options and must remain disabled/)).toBeInTheDocument();
  });

  it('covers add-on create, update, delete, and reorder rollback controls', async () => {
    apiMocks.listAddons.mockResolvedValue(page(addons));
    apiMocks.reorderAddons.mockRejectedValue(new Error('stale reorder'));
    const { onRefreshWorkspace } = renderEditor();

    fireEvent.click(screen.getByRole('button', { name: 'Add-ons (0)' }));
    await waitFor(() => expect(screen.getByText('Second Addon')).toBeInTheDocument());

    const moveUp = screen.getByRole('button', { name: 'Move Second Addon up' });
    expect(moveUp.tagName).toBe('BUTTON');
    fireEvent.click(moveUp);
    await waitFor(() => {
      expect(apiMocks.reorderAddons).toHaveBeenCalledWith({
        menuItemId: item.id,
        orderedIds: [addons[1].id, addons[0].id],
      });
      expect(apiMocks.listAddons).toHaveBeenCalledTimes(2);
      expect(toastMocks.error).toHaveBeenCalledWith('stale reorder');
    });

    fireEvent.click(screen.getByRole('button', { name: 'Add Add-on' }));
    fireEvent.change(screen.getByLabelText('Add-on Name *'), { target: { value: 'Garlic Dip' } });
    fireEvent.change(screen.getByLabelText('Price ($) *'), { target: { value: '0' } });
    fireEvent.click(screen.getByRole('button', { name: 'Save Add-on' }));
    await waitFor(() => {
      expect(apiMocks.createAddon).toHaveBeenCalledWith(expect.objectContaining({
        name: 'Garlic Dip',
        price: 0,
        menuItemId: item.id,
        restaurantId: restaurant.id,
      }));
      expect(onRefreshWorkspace).toHaveBeenCalled();
    });

    fireEvent.click(screen.getByRole('switch', { name: 'Toggle First Addon availability' }));
    await waitFor(() => {
      expect(apiMocks.updateAddon).toHaveBeenCalledWith(addons[0].id, { available: false });
    });
    fireEvent.click(screen.getByRole('button', { name: 'Delete First Addon' }));
    await waitFor(() => expect(apiMocks.deleteAddon).toHaveBeenCalledWith(addons[0].id));
  });

  it('clears nested editor state and ignores stale loads when the item changes', async () => {
    const nextItem = {
      ...item,
      id: 'item-2',
      name: 'Second Item',
      categoryId: categories[1].id,
    };
    let resolveFirstLoad: ((value: ReturnType<typeof page<AdminAddon>>) => void) | undefined;
    apiMocks.listAddons
      .mockImplementationOnce(() => new Promise((resolve) => {
        resolveFirstLoad = resolve;
      }))
      .mockResolvedValueOnce(page([{ ...addons[1], menuItemId: nextItem.id }]));

    const onSaveItem = vi.fn(async () => item);
    const onRefreshWorkspace = vi.fn(async () => undefined);
    const editor = (currentItem: AdminMenuItem) => (
      <MenuItemEditorDialog
        item={currentItem}
        workspace={workspaceFor(currentItem)}
        open
        onClose={vi.fn()}
        onSaveItem={onSaveItem}
        onRefreshWorkspace={onRefreshWorkspace}
      />
    );

    const { rerender } = render(editor(item));
    fireEvent.click(screen.getByRole('button', { name: 'Add-ons (0)' }));
    await waitFor(() => expect(apiMocks.listAddons).toHaveBeenCalledWith(restaurant.id, item.id));
    fireEvent.click(screen.getByRole('button', { name: 'Add Add-on' }));
    expect(screen.getByRole('heading', { name: 'Add Menu Item Add-on' })).toBeInTheDocument();

    rerender(editor(nextItem));
    expect(screen.queryByRole('heading', { name: 'Add Menu Item Add-on' })).not.toBeInTheDocument();
    expect(screen.queryByText('First Addon')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Add-ons (0)' })).toBeInTheDocument();

    fireEvent.click(screen.getByRole('button', { name: 'Add-ons (0)' }));
    await waitFor(() => expect(screen.getByText('Second Addon')).toBeInTheDocument());

    resolveFirstLoad?.(page(addons));
    await waitFor(() => expect(screen.queryByText('First Addon')).not.toBeInTheDocument());
    expect(screen.getByText('Second Addon')).toBeInTheDocument();
  });

  it('closes modifier editors and clears loaded groups when the item changes', async () => {
    const nextItem = { ...item, id: 'item-2', name: 'Second Item' };
    apiMocks.listModifierGroups.mockResolvedValue(page([group]));
    apiMocks.listModifierOptions.mockResolvedValue(page([option]));
    const onSaveItem = vi.fn(async () => item);
    const onRefreshWorkspace = vi.fn(async () => undefined);
    const editor = (currentItem: AdminMenuItem) => (
      <MenuItemEditorDialog
        item={currentItem}
        workspace={workspaceFor(currentItem)}
        open
        onClose={vi.fn()}
        onSaveItem={onSaveItem}
        onRefreshWorkspace={onRefreshWorkspace}
      />
    );

    const { rerender } = render(editor(item));
    fireEvent.click(screen.getByRole('button', { name: 'Modifiers (0)' }));
    await waitFor(() => expect(screen.getByText('Spice Level')).toBeInTheDocument());
    fireEvent.click(screen.getByRole('button', { name: 'Add Option' }));
    expect(screen.getByRole('heading', { name: 'Add Option to Spice Level' })).toBeInTheDocument();

    rerender(editor(nextItem));

    expect(screen.queryByRole('heading', { name: 'Add Option to Spice Level' })).not.toBeInTheDocument();
    expect(screen.queryByText('Spice Level')).not.toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Modifiers (0)' })).toBeInTheDocument();
  });
});
