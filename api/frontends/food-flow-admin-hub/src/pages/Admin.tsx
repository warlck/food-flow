import { FormEvent, useCallback, useEffect, useMemo, useState } from 'react';
import {
  ArrowRight, Banknote, BarChart3, Bike, BookOpen, Boxes, Building2, Check, ChevronDown, ChevronRight,
  ChefHat, CircleAlert, Clock3, CreditCard, Grid2X2, HelpCircle, ImageOff, LayoutDashboard, List,
  Loader2, Mail, MapPin, Menu, MoreHorizontal, PackageCheck, Pencil, Phone, Plus, ReceiptText,
  Puzzle, RefreshCw, Search, Settings, ShoppingBag, Sparkles, Store, Tag, Trash2, UtensilsCrossed, XCircle,
} from 'lucide-react';
import { toast } from 'sonner';
import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from '@/components/ui/dropdown-menu';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Switch } from '@/components/ui/switch';
import { Textarea } from '@/components/ui/textarea';
import {
  AddonInput, AdminAddon, AdminCategory, AdminMenuItem, AdminOrder, AdminPromotion, AdminRestaurant, AdminWorkspace,
  CategoryInput, MenuItemInput, OrderStatus, OrderType, PaymentStatus, PromotionInput, RestaurantInput, adminApi,
} from '@/lib/admin-api';
import './Admin.css';

type EditorState =
  | { kind: 'restaurant'; value?: AdminRestaurant }
  | { kind: 'category'; value?: AdminCategory }
  | { kind: 'item'; value?: AdminMenuItem; categoryId?: string }
  | { kind: 'addon'; value?: AdminAddon; categoryId?: string }
  | { kind: 'promotion'; value?: AdminPromotion }
  | null;

const NAME_PATTERN = /^[\p{L}\p{N}' -]{3,100}$/u;

type Section = 'menu' | 'orders' | 'promotions';

const navItems: { icon: typeof LayoutDashboard; label: string; section?: Section; soon?: boolean }[] = [
  { icon: LayoutDashboard, label: 'Overview' },
  { icon: BookOpen, label: 'Menu & inventory', section: 'menu' },
  { icon: ReceiptText, label: 'Orders', section: 'orders' },
  { icon: Tag, label: 'Promotions', section: 'promotions' },
  { icon: BarChart3, label: 'Sales & insights', soon: true },
  { icon: Settings, label: 'Settings' },
];

function nextOrderStatus(order: Pick<AdminOrder, 'orderStatus' | 'orderType'>): OrderStatus | null {
  switch (order.orderStatus) {
    case 'pending': return 'confirmed';
    case 'confirmed': return 'preparing';
    case 'preparing': return 'ready';
    case 'ready': return order.orderType === 'delivery' ? 'out_for_delivery' : 'completed';
    case 'out_for_delivery': return 'completed';
    default: return null;
  }
}

const NEXT_STATUS_LABELS: Partial<Record<OrderStatus, string>> = {
  confirmed: 'Confirm order',
  preparing: 'Start preparing',
  ready: 'Mark ready',
  out_for_delivery: 'Mark out for delivery',
  completed: 'Complete order',
};

const ORDER_STATUS_LABELS: Record<OrderStatus, string> = {
  pending: 'Pending',
  confirmed: 'Confirmed',
  preparing: 'Preparing',
  ready: 'Ready',
  out_for_delivery: 'Out for delivery',
  completed: 'Completed',
  cancelled: 'Cancelled',
};

const ORDER_STATUS_STYLES: Record<OrderStatus, string> = {
  pending: 'bg-[#FFF8E7] text-[#B45309]',
  confirmed: 'bg-[#EFF6FF] text-[#1D4ED8]',
  preparing: 'bg-[#FFF1EB] text-[#FF4500]',
  ready: 'bg-[#F0FDF4] text-[#15803D]',
  out_for_delivery: 'bg-[#EEF2FF] text-[#4F46E5]',
  completed: 'bg-[#E8F5E9] text-[#2E7D32]',
  cancelled: 'bg-[#FFEBEE] text-[#C62828]',
};

const PAYMENT_STATUS_STYLES: Record<PaymentStatus, string> = {
  pending: 'bg-[#F3F4F6] text-[#4B5563]',
  processing: 'bg-[#EFF6FF] text-[#1D4ED8]',
  paid: 'bg-[#E8F5E9] text-[#2E7D32]',
  failed: 'bg-[#FFEBEE] text-[#C62828]',
  refunded: 'bg-[#F3E8FF] text-[#7E22CE]',
};

const ORDER_STATUS_OPTIONS: { value: 'all' | OrderStatus; label: string }[] = [
  { value: 'all', label: 'All statuses' },
  { value: 'pending', label: 'Pending' },
  { value: 'confirmed', label: 'Confirmed' },
  { value: 'preparing', label: 'Preparing' },
  { value: 'ready', label: 'Ready' },
  { value: 'out_for_delivery', label: 'Out for delivery' },
  { value: 'completed', label: 'Completed' },
  { value: 'cancelled', label: 'Cancelled' },
];

const ORDER_TYPE_OPTIONS: { value: 'all' | OrderType; label: string }[] = [
  { value: 'all', label: 'All types' },
  { value: 'pickup', label: 'Pickup' },
  { value: 'delivery', label: 'Delivery' },
];



function formatCurrency(value: number) {
  return new Intl.NumberFormat('en-SG', { style: 'currency', currency: 'SGD' }).format(value);
}

function formatOrderTime(value: string) {
  return new Intl.DateTimeFormat('en-SG', { dateStyle: 'medium', timeStyle: 'short' }).format(new Date(value));
}

function initials(name: string) {
  return name.split(/\s+/).slice(0, 2).map((part) => part[0]).join('').toUpperCase();
}

function validateName(name: string) {
  if (!NAME_PATTERN.test(name)) throw new Error('Use 3–100 letters, numbers, spaces, apostrophes or hyphens for names.');
}

const Field = ({ label, htmlFor, required, hint, children }: { label: string; htmlFor: string; required?: boolean; hint?: string; children: React.ReactNode }) => (
  <div className="flex flex-col justify-between h-full space-y-1.5">
    <div>
      <Label htmlFor={htmlFor} className="text-[13px] font-semibold text-[#374151]">{label}{required && <span className="ml-1 text-[#F44336]">*</span>}</Label>
      {hint && <p className="mt-0.5 text-[11px] text-[#9CA3AF]">{hint}</p>}
    </div>
    <div className="mt-auto">{children}</div>
  </div>
);

export default function Admin() {
  const [restaurants, setRestaurants] = useState<AdminRestaurant[]>([]);
  const [selectedId, setSelectedId] = useState('');
  const [workspace, setWorkspace] = useState<AdminWorkspace | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [editor, setEditor] = useState<EditorState>(null);
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [addonCategory, setAddonCategory] = useState('');
  const [availability, setAvailability] = useState('all');
  const [query, setQuery] = useState('');
  const [view, setView] = useState<'grid' | 'list'>('grid');
  const [section, setSection] = useState<Section>('menu');
  const [orders, setOrders] = useState<AdminOrder[]>([]);
  const [ordersLoading, setOrdersLoading] = useState(false);
  const [ordersLoaded, setOrdersLoaded] = useState(false);
  const [ordersRefreshKey, setOrdersRefreshKey] = useState(0);
  const [orderStatusFilter, setOrderStatusFilter] = useState<'all' | OrderStatus>('all');
  const [orderTypeFilter, setOrderTypeFilter] = useState<'all' | OrderType>('all');
  const [orderQuery, setOrderQuery] = useState('');
  const [selectedOrder, setSelectedOrder] = useState<AdminOrder | null>(null);
  const [orderActionBusy, setOrderActionBusy] = useState(false);

  const [promotions, setPromotions] = useState<AdminPromotion[]>([]);
  const [promotionsLoading, setPromotionsLoading] = useState(false);
  const [promotionsQuery, setPromotionsQuery] = useState('');

  const loadPromotions = useCallback(async (quiet = false) => {
    if (!quiet) setPromotionsLoading(true);
    try {
      const page = await adminApi.listPromotions(selectedId || undefined);
      setPromotions(page.items);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Could not load promotions');
    } finally {
      setPromotionsLoading(false);
    }
  }, [selectedId]);

  useEffect(() => {
    if (section === 'promotions') {
      loadPromotions();
    }
  }, [section, selectedId, loadPromotions]);

  const togglePromotionEnabled = async (promo: AdminPromotion, enabled: boolean) => {
    try {
      await adminApi.updatePromotion(promo.id, { enabled });
      setPromotions((current) => current.map((p) => (p.id === promo.id ? { ...p, enabled } : p)));
      toast.success(enabled ? 'Promotion campaign enabled' : 'Promotion campaign disabled');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Could not update promotion');
    }
  };

  const deletePromotion = async (promo: AdminPromotion) => {
    if (!window.confirm(`Delete promo code ${promo.code}? This cannot be undone.`)) return;
    try {
      await adminApi.deletePromotion(promo.id);
      setPromotions((current) => current.filter((p) => p.id !== promo.id));
      toast.success('Promotion campaign deleted');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Could not delete promotion');
    }
  };

  const toggleRestaurantStatus = async () => {
    if (!workspace) return;
    const newStatus = !workspace.restaurant.enabled;
    try {
      await mutateWorkspace(
        () => adminApi.updateRestaurant(workspace.restaurant.id, { enabled: newStatus }),
        newStatus ? 'Restaurant is now Live' : 'Restaurant is now Paused'
      );
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Could not update restaurant status');
    }
  };

  const loadWorkspace = useCallback(async (restaurantId: string, quiet = false) => {
    if (!restaurantId) return;
    if (!quiet) setLoading(true); else setRefreshing(true);
    try {
      const data = await adminApi.getWorkspace(restaurantId);
      setWorkspace(data);
      setAddonCategory((current) => data.categories.some((category) => category.id === current) ? current : data.categories[0]?.id ?? '');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Could not load workspace');
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, []);

  useEffect(() => {
    let active = true;
    (async () => {
      try {
        const page = await adminApi.listRestaurants();
        if (!active) return;
        setRestaurants(page.items);
        if (page.items.length) {
          setSelectedId(page.items[0].id);
          await loadWorkspace(page.items[0].id);
        } else {
          setLoading(false);
        }
      } catch (error) {
        if (!active) return;
        toast.error(error instanceof Error ? error.message : 'Could not load restaurants');
        setLoading(false);
      }
    })();
    return () => { active = false; };
  }, []);

  const categoryCounts = useMemo(() => {
    const counts = new Map<string, number>();
    workspace?.menuItems.forEach((item) => counts.set(item.categoryId, (counts.get(item.categoryId) ?? 0) + 1));
    return counts;
  }, [workspace]);

  const filteredItems = useMemo(() => {
    if (!workspace) return [];
    const normalizedQuery = query.trim().toLowerCase();
    return workspace.menuItems.filter((item) => {
      const matchesCategory = selectedCategory === 'all' || item.categoryId === selectedCategory;
      const matchesAvailability = availability === 'all' || (availability === 'available' ? item.available : !item.available);
      const matchesQuery = !normalizedQuery || `${item.name} ${item.description}`.toLowerCase().includes(normalizedQuery);
      return matchesCategory && matchesAvailability && matchesQuery;
    });
  }, [availability, query, selectedCategory, workspace]);

  const availableItems = workspace?.menuItems.filter((item) => item.available).length ?? 0;
  const unavailableItems = (workspace?.menuItems.length ?? 0) - availableItems;
  const availabilityPercent = workspace?.menuItems.length ? Math.round((availableItems / workspace.menuItems.length) * 100) : 0;

  const mutateWorkspace = async (operation: () => Promise<unknown>, successMessage: string) => {
    try {
      await operation();
      toast.success(successMessage);
      setEditor(null);
      if (selectedId) await loadWorkspace(selectedId, true);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Something went wrong');
      throw error;
    }
  };

  const toggleAvailability = async (item: AdminMenuItem, available: boolean) => {
    setWorkspace((current) => current ? { ...current, menuItems: current.menuItems.map((entry) => entry.id === item.id ? { ...entry, available } : entry) } : current);
    try {
      await adminApi.updateMenuItem(item.id, { available });
      toast.success(available ? 'Item is now available' : 'Item marked unavailable');
    } catch (error) {
      setWorkspace((current) => current ? { ...current, menuItems: current.menuItems.map((entry) => entry.id === item.id ? { ...entry, available: item.available } : entry) } : current);
      toast.error(error instanceof Error ? error.message : 'Availability could not be changed');
    }
  };

  const handleRestaurantChange = async (id: string) => {
    setSelectedId(id);
    setSelectedCategory('all');
    setOrders([]);
    setOrdersLoaded(false);
    setSelectedOrder(null);
    await loadWorkspace(id);
  };

  const deleteItem = async (item: AdminMenuItem) => {
    if (!window.confirm(`Delete ${item.name}? This cannot be undone.`)) return;
    await mutateWorkspace(() => adminApi.deleteMenuItem(item.id), 'Menu item deleted');
  };

  const toggleAddonAvailability = async (addon: AdminAddon, available: boolean) => {
    setWorkspace((current) => current ? { ...current, addons: current.addons.map((entry) => entry.id === addon.id ? { ...entry, available } : entry) } : current);
    try {
      await adminApi.updateAddon(addon.id, { available });
      toast.success(available ? 'Add-on is now available' : 'Add-on marked unavailable');
    } catch (error) {
      setWorkspace((current) => current ? { ...current, addons: current.addons.map((entry) => entry.id === addon.id ? { ...entry, available: addon.available } : entry) } : current);
      toast.error(error instanceof Error ? error.message : 'Add-on availability could not be changed');
    }
  };

  const deleteAddon = async (addon: AdminAddon) => {
    if (!window.confirm(`Delete ${addon.name}? This removes it from every item in this category.`)) return;
    await mutateWorkspace(() => adminApi.deleteAddon(addon.id), 'Add-on deleted');
  };

  useEffect(() => {
    if (section !== 'orders' || !selectedId) return;
    let active = true;
    setOrdersLoading(true);
    adminApi.listOrders(selectedId, {
      orderStatus: orderStatusFilter === 'all' ? undefined : orderStatusFilter,
      orderType: orderTypeFilter === 'all' ? undefined : orderTypeFilter,
    })
      .then((page) => {
        if (!active) return;
        setOrders(page.items);
        setOrdersLoaded(true);
      })
      .catch((error) => {
        if (!active) return;
        toast.error(error instanceof Error ? error.message : 'Could not load orders');
      })
      .finally(() => {
        if (active) setOrdersLoading(false);
      });
    return () => { active = false; };
  }, [section, selectedId, orderStatusFilter, orderTypeFilter, ordersRefreshKey]);

  const visibleOrders = useMemo(() => {
    const normalizedQuery = orderQuery.trim().toLowerCase();
    if (!normalizedQuery) return orders;
    return orders.filter((order) =>
      `${order.customerName} ${order.customerEmail} ${order.customerPhone} ${order.id}`.toLowerCase().includes(normalizedQuery));
  }, [orders, orderQuery]);

  const applyOrderUpdate = (updated: AdminOrder) => {
    setOrders((current) => current.map((entry) => (entry.id === updated.id ? updated : entry)));
    setSelectedOrder((current) => (current?.id === updated.id ? updated : current));
  };

  const advanceOrder = async (order: AdminOrder) => {
    const next = nextOrderStatus(order);
    if (!next) return;
    setOrderActionBusy(true);
    try {
      const updated = await adminApi.updateOrderStatus(order.id, { orderStatus: next });
      applyOrderUpdate(updated);
      toast.success(`Order #${order.id.slice(0, 8)} marked ${ORDER_STATUS_LABELS[next].toLowerCase()}`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Order status could not be updated');
    } finally {
      setOrderActionBusy(false);
    }
  };

  const markOrderPaid = async (order: AdminOrder) => {
    setOrderActionBusy(true);
    try {
      const updated = await adminApi.updateOrderStatus(order.id, { paymentStatus: 'paid' });
      applyOrderUpdate(updated);
      toast.success(`Order #${order.id.slice(0, 8)} marked as paid`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Payment status could not be updated');
    } finally {
      setOrderActionBusy(false);
    }
  };

  const cancelOrder = async (order: AdminOrder) => {
    if (!window.confirm(`Cancel order #${order.id.slice(0, 8)} for ${order.customerName}?`)) return;
    setOrderActionBusy(true);
    try {
      await adminApi.cancelOrder(order.id);
      applyOrderUpdate({ ...order, orderStatus: 'cancelled' });
      toast.success(`Order #${order.id.slice(0, 8)} cancelled`);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Order could not be cancelled');
    } finally {
      setOrderActionBusy(false);
    }
  };

  const selectedCategoryName = selectedCategory === 'all' ? 'All menu items' : workspace?.categories.find((category) => category.id === selectedCategory)?.name ?? 'Menu items';

  return (
    <div className="admin-shell">
      <aside className="admin-sidebar px-3 py-5">
        <div className="flex items-center gap-3 px-2.5 pb-7">
          <div className="admin-brand-mark flex h-10 w-10 shrink-0 items-center justify-center rounded-lg bg-[#FF4500] text-white">
            <ChefHat size={22} strokeWidth={2.2} />
          </div>
          <div className="admin-sidebar-copy min-w-0">
            <div className="text-[17px] font-bold tracking-[-.02em] text-[#FF4500]">FoodFlow</div>
            <div className="text-[10px] font-semibold uppercase tracking-[.18em] text-[#9CA3AF]">Restaurant studio</div>
          </div>
        </div>

        <div className="admin-sidebar-copy mx-1 mb-6 rounded-xl border border-white/10 bg-white/[.055] p-2.5">
          <div className="flex w-full items-center gap-2.5 text-left">
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-[#FFF1EB] text-xs font-bold text-[#FF4500]">{workspace ? initials(workspace.restaurant.name) : 'FF'}</div>
            <div className="min-w-0 flex-1">
              <div className="truncate text-[12px] font-semibold text-white">{workspace?.restaurant.name ?? 'Choose restaurant'}</div>
              <div className="mt-0.5 flex items-center gap-1 text-[10px] text-[#D1D5DB]"><span className="h-1.5 w-1.5 rounded-full bg-[#FFB72B]" /> {workspace?.restaurant.enabled ? 'Open for business' : 'Paused'}</div>
            </div>
          </div>
        </div>

        <nav className="space-y-1">
          <div className="admin-nav-label px-3 pb-2 text-[10px] font-bold uppercase tracking-[.16em] text-[#9CA3AF]">Workspace</div>
          {navItems.map(({ icon: Icon, label, section: navSection, soon }) => (
            <button key={label} onClick={() => navSection && setSection(navSection)} className={`admin-nav-item flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left text-[13px] font-semibold ${navSection === section ? 'active' : ''}`}>
              <Icon size={18} strokeWidth={navSection === section ? 2.3 : 1.8} />
              <span className="admin-nav-label flex-1">{label}</span>
              {soon && <span className="admin-nav-label rounded-full border border-current/20 px-1.5 py-0.5 text-[8px] font-bold uppercase tracking-wider opacity-70">Soon</span>}
            </button>
          ))}
        </nav>

        <div className="admin-sidebar-footer mt-auto space-y-2 px-2">
          <button className="flex w-full items-center gap-2 rounded-lg px-2 py-2 text-xs font-medium text-[#D1D5DB]"><HelpCircle size={16} /> Help centre</button>
          <div className="border-t border-white/10 pt-4 text-[10px] leading-relaxed text-[#9CA3AF]">Built around your live FoodFlow menu data.</div>
        </div>
      </aside>

      <main className="admin-main">
        <header className="admin-topbar sticky top-0 z-20 flex min-h-[72px] items-center justify-between gap-4 px-4 sm:px-8">
          <div className="flex min-w-0 items-center gap-3">
            <button className="md:hidden"><Menu size={21} /></button>
            <div className="hidden h-8 w-px bg-[#E5E7EB] sm:block" />
            <div className="min-w-0">
              <div className="flex items-center gap-2">
                <Select value={selectedId} onValueChange={handleRestaurantChange}>
                  <SelectTrigger className="h-auto w-auto max-w-[260px] gap-2 border-0 bg-transparent p-0 text-[14px] font-bold shadow-none focus:ring-0">
                    <SelectValue placeholder="Select a restaurant" />
                  </SelectTrigger>
                  <SelectContent>
                    {restaurants.map((restaurant) => <SelectItem value={restaurant.id} key={restaurant.id}>{restaurant.name}</SelectItem>)}
                  </SelectContent>
                </Select>
                {workspace && (
                  <span
                    className={`rounded-full px-2.5 py-0.5 text-[9px] font-bold uppercase tracking-[.1em] ${
                      workspace.restaurant.enabled
                        ? 'bg-[#E8F5E9] text-[#2E7D32]'
                        : 'bg-[#FFEBEE] text-[#C62828]'
                    }`}
                  >
                    {workspace.restaurant.enabled ? 'Live' : 'Paused'}
                  </span>
                )}
                <Button
                  variant="ghost"
                  size="sm"
                  className="h-6 gap-1 rounded-md px-2 text-[11px] font-semibold text-[#6B7280] hover:bg-[#FFF7F3] hover:text-[#FF4500]"
                  onClick={() => workspace && setEditor({ kind: 'restaurant', value: workspace.restaurant })}
                  disabled={!workspace}
                  title="Edit restaurant profile and settings"
                >
                  <Pencil size={12} className="text-[#FF4500]" />
                  <span>Edit</span>
                </Button>
              </div>
              <div className="mt-0.5 flex items-center gap-1 text-[11px] text-[#6B7280]"><MapPin size={11} /> <span className="truncate">{workspace?.restaurant.address ?? 'Restaurant workspace'}</span></div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="outline"
              className="h-9 gap-1.5 rounded-xl border-[#E5E7EB] bg-white px-3 text-xs font-semibold text-[#374151] shadow-sm hover:border-[#FF8C42] hover:bg-[#FFF7F3] hover:text-[#FF4500] transition-colors"
              onClick={() => setEditor({ kind: 'restaurant' })}
              title="Add a new restaurant"
            >
              <Plus size={14} className="text-[#FF4500]" />
              <span className="hidden sm:inline">Add Restaurant</span>
              <span className="sm:hidden">Add</span>
            </Button>
          </div>
        </header>

        <div className="admin-content">
          {section === 'menu' ? (
            <>
          <section className="mb-7 flex flex-col justify-between gap-4 md:flex-row md:items-end">
            <div>
              <div className="mb-1.5 flex items-center gap-2 text-[11px] font-bold uppercase tracking-[.14em] text-[#FF4500]"><BookOpen size={13} /> Menu workspace</div>
              <h1 className="text-[28px] font-bold tracking-[-.035em] text-[#333333] sm:text-[32px]">Build a menu people remember.</h1>
              <p className="mt-1 max-w-2xl text-[13px] text-[#6B7280]">Organise categories, curate every dish, and keep availability accurate from one calm workspace.</p>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" className="h-9 gap-2 rounded-lg border-[#E5E7EB] bg-white text-xs hover:border-[#FF8C42] hover:bg-[#FFF7F3] hover:text-[#FF4500]" onClick={() => setEditor({ kind: 'category' })} disabled={!workspace}><Plus size={15} /> Add category</Button>
            </div>
          </section>

          {loading ? (
            <div className="flex min-h-[420px] items-center justify-center"><Loader2 className="animate-spin text-[#FF4500]" size={28} /></div>
          ) : !workspace ? (
            <EmptyRestaurant onCreate={() => setEditor({ kind: 'restaurant' })} />
          ) : (
            <>
              <section className="mb-6 grid grid-cols-2 gap-3 xl:grid-cols-4">
                <Stat icon={Boxes} label="Categories" value={workspace.categories.length} note={`${workspace.addons.length} add-ons configured`} />
                <Stat icon={ShoppingBag} label="Menu items" value={workspace.menuItems.length} note={`${formatCurrency(workspace.menuItems.reduce((sum, item) => sum + item.price, 0) / Math.max(workspace.menuItems.length, 1))} avg. price`} />
                <Stat icon={PackageCheck} label="Available now" value={`${availabilityPercent}%`} note={`${availableItems} items can be ordered`} progress={availabilityPercent} />
                <Stat icon={CircleAlert} label="Needs attention" value={unavailableItems} note={unavailableItems ? 'Unavailable items' : 'Everything looks good'} attention={unavailableItems > 0} />
              </section>

              <SetupGuide restaurant={workspace.restaurant} categoryCount={workspace.categories.length} itemCount={workspace.menuItems.length} />

              <section className="admin-panel mt-6 overflow-hidden rounded-2xl">
                <div className="flex flex-col border-b border-[#E5E7EB] lg:flex-row">
                  <CategoryRail
                    categories={workspace.categories}
                    counts={categoryCounts}
                    total={workspace.menuItems.length}
                    selected={selectedCategory}
                    onSelect={setSelectedCategory}
                    onAdd={() => setEditor({ kind: 'category' })}
                    onAddItem={(catId) => setEditor({ kind: 'item', categoryId: catId })}
                    onEdit={(category) => setEditor({ kind: 'category', value: category })}
                  />
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-col gap-3 border-b border-[#E5E7EB] p-4 sm:flex-row sm:items-center sm:justify-between sm:px-5">
                      <div>
                        <h2 className="text-[17px] font-bold tracking-[-.02em]">{selectedCategoryName}</h2>
                        <p className="mt-0.5 text-[11px] text-[#9CA3AF]">{filteredItems.length} {filteredItems.length === 1 ? 'item' : 'items'} shown</p>
                      </div>
                      <div className="flex flex-wrap items-center gap-2">
                        <Button
                          className="admin-primary h-9 gap-1.5 rounded-lg px-3 text-xs"
                          onClick={() => setEditor({ kind: 'item', categoryId: selectedCategory !== 'all' ? selectedCategory : undefined })}
                          disabled={!workspace.categories.length}
                        >
                          <Plus size={15} /> New Menu Item
                        </Button>
                        <div className="relative min-w-[180px] flex-1 sm:w-[200px] sm:flex-none">
                          <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-[#9CA3AF]" size={15} />
                          <Input value={query} onChange={(event) => setQuery(event.target.value)} placeholder="Search menu" className="admin-input h-9 rounded-lg pl-9 text-xs" />
                        </div>
                        <Select value={availability} onValueChange={setAvailability}>
                          <SelectTrigger className="admin-input h-9 w-[126px] rounded-lg text-xs"><SelectValue /></SelectTrigger>
                          <SelectContent><SelectItem value="all">All stock</SelectItem><SelectItem value="available">Available</SelectItem><SelectItem value="unavailable">Unavailable</SelectItem></SelectContent>
                        </Select>
                        <div className="flex rounded-lg border border-[#E5E7EB] bg-[#F3F4F6] p-0.5">
                          <button onClick={() => setView('grid')} className={`rounded-md p-1.5 ${view === 'grid' ? 'bg-white text-[#FF4500] shadow-sm' : 'text-[#9CA3AF]'}`} aria-label="Grid view"><Grid2X2 size={15} /></button>
                          <button onClick={() => setView('list')} className={`rounded-md p-1.5 ${view === 'list' ? 'bg-white text-[#FF4500] shadow-sm' : 'text-[#9CA3AF]'}`} aria-label="List view"><List size={15} /></button>
                        </div>
                        <Button variant="ghost" size="icon" className="h-9 w-9" onClick={() => loadWorkspace(selectedId, true)} disabled={refreshing} aria-label="Refresh workspace"><RefreshCw size={15} className={refreshing ? 'animate-spin' : ''} /></Button>
                      </div>
                    </div>

                    <div className="admin-empty-pattern min-h-[440px] p-4 sm:p-5">
                      {filteredItems.length ? (
                        <div className={view === 'grid' ? 'grid grid-cols-1 gap-4 sm:grid-cols-2 2xl:grid-cols-3' : 'space-y-3'}>
                          {filteredItems.map((item) => (
                            <MenuCard key={item.id} item={item} category={workspace.categories.find((category) => category.id === item.categoryId)} view={view} onEdit={() => setEditor({ kind: 'item', value: item })} onDelete={() => deleteItem(item)} onAvailability={(value) => toggleAvailability(item, value)} />
                          ))}
                        </div>
                      ) : (
                        <div className="flex min-h-[390px] flex-col items-center justify-center px-5 text-center">
                          <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl border border-[#FED7C7] bg-white text-[#FF4500] shadow-sm"><UtensilsCrossed size={23} /></div>
                          <h3 className="text-base font-bold">No menu items here yet</h3>
                          <p className="mt-1 max-w-sm text-xs leading-relaxed text-[#6B7280]">Create a dish in this category or clear your filters to see the rest of the menu.</p>
                          <Button className="admin-primary mt-4 h-9 gap-2 text-xs" onClick={() => setEditor({ kind: 'item' })} disabled={!workspace.categories.length}><Plus size={15} /> Create menu item</Button>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
                <div className="flex flex-wrap items-center justify-between gap-3 bg-[#FAFAFA] px-5 py-3 text-[10px] text-[#6B7280]">
                  <span>Changes update the customer menu immediately.</span>
                  <span className="flex items-center gap-1.5"><Clock3 size={12} /> Last synced just now</span>
                </div>
              </section>

              <AddonManager
                categories={workspace.categories}
                addons={workspace.addons}
                selectedCategory={addonCategory}
                onCategoryChange={setAddonCategory}
                onCreate={() => setEditor({ kind: 'addon', categoryId: addonCategory || workspace.categories[0]?.id })}
                onEdit={(addon) => setEditor({ kind: 'addon', value: addon })}
                onDelete={deleteAddon}
                onAvailability={toggleAddonAvailability}
              />
            </>
          )}
            </>
          ) : section === 'promotions' ? (
            <PromotionsSection
              promotions={promotions}
              loading={promotionsLoading}
              query={promotionsQuery}
              onQueryChange={setPromotionsQuery}
              onCreate={() => setEditor({ kind: 'promotion' })}
              onEdit={(promo) => setEditor({ kind: 'promotion', value: promo })}
              onToggleEnabled={togglePromotionEnabled}
              onDelete={deletePromotion}
              onRefresh={() => loadPromotions()}
            />
          ) : (
            <OrdersSection
              orders={visibleOrders}
              loading={ordersLoading}
              loaded={ordersLoaded}
              statusFilter={orderStatusFilter}
              onStatusFilter={setOrderStatusFilter}
              typeFilter={orderTypeFilter}
              onTypeFilter={setOrderTypeFilter}
              query={orderQuery}
              onQuery={setOrderQuery}
              onRefresh={() => setOrdersRefreshKey((key) => key + 1)}
              onSelect={setSelectedOrder}
            />
          )}
        </div>
      </main>

      <EditorDialog
        editor={editor}
        workspace={workspace}
        onClose={() => setEditor(null)}
        onSave={async (kind, input, existingId) => {
          if (kind === 'promotion') {
            const promoInput = input as PromotionInput;
            if (existingId) {
              await adminApi.updatePromotion(existingId, promoInput);
              toast.success('Promotion campaign updated');
            } else {
              await adminApi.createPromotion(promoInput);
              toast.success('Promotion campaign created');
            }
            setEditor(null);
            await loadPromotions();
            return;
          }
          if (kind === 'restaurant') {
            const restaurantInput = input as RestaurantInput;
            validateName(restaurantInput.name);
            if (existingId) await mutateWorkspace(() => adminApi.updateRestaurant(existingId, restaurantInput), 'Restaurant updated');
            else {
              await mutateWorkspace(async () => {
                const created = await adminApi.createRestaurant(restaurantInput);
                const page = await adminApi.listRestaurants();
                setRestaurants(page.items); setSelectedId(created.id); await loadWorkspace(created.id);
              }, 'Restaurant created');
            }
          }
          if (kind === 'category' && workspace) {
            const categoryInput = input as CategoryInput;
            validateName(categoryInput.name);
            await mutateWorkspace(() => existingId ? adminApi.updateCategory(existingId, categoryInput) : adminApi.createCategory(categoryInput), existingId ? 'Category updated' : 'Category created');
          }
          if (kind === 'item' && workspace) {
            const itemInput = input as MenuItemInput;
            validateName(itemInput.name);
            await mutateWorkspace(() => existingId ? adminApi.updateMenuItem(existingId, itemInput) : adminApi.createMenuItem(itemInput), existingId ? 'Menu item updated' : 'Menu item created');
          }
          if (kind === 'addon' && workspace) {
            const addonInput = input as AddonInput;
            validateName(addonInput.name);
            await mutateWorkspace(() => existingId ? adminApi.updateAddon(existingId, { name: addonInput.name, description: addonInput.description, price: addonInput.price, maxQuantity: addonInput.maxQuantity }) : adminApi.createAddon(addonInput), existingId ? 'Add-on updated' : 'Add-on created');
            setAddonCategory(addonInput.categoryId);
          }
        }}
      />

      <OrderDetailDialog
        order={selectedOrder}
        busy={orderActionBusy}
        onClose={() => setSelectedOrder(null)}
        onAdvance={advanceOrder}
        onMarkPaid={markOrderPaid}
        onCancel={cancelOrder}
      />
    </div>
  );
}

function Stat({ icon: Icon, label, value, note, progress, attention }: { icon: typeof Boxes; label: string; value: string | number; note: string; progress?: number; attention?: boolean }) {
  return (
    <div className="admin-stat-card rounded-xl p-3.5 sm:p-4">
      <div className="flex items-start justify-between gap-3">
        <div><p className="text-[10px] font-bold uppercase tracking-[.11em] text-[#9CA3AF]">{label}</p><p className="mt-1 text-[22px] font-bold tracking-[-.04em] text-[#333333]">{value}</p></div>
        <div className={`admin-stat-icon flex h-8 w-8 items-center justify-center rounded-lg ${attention ? '!bg-[#FFEBEE] !text-[#F44336]' : ''}`}><Icon size={16} /></div>
      </div>
      {progress !== undefined && <div className="mt-2.5 h-1 overflow-hidden rounded-full bg-[#F3F4F6]"><div className="h-full rounded-full bg-[#FF8C42]" style={{ width: `${progress}%` }} /></div>}
      <p className="mt-2 text-[10px] text-[#6B7280]">{note}</p>
    </div>
  );
}

function SetupGuide({ restaurant, categoryCount, itemCount }: { restaurant: AdminRestaurant; categoryCount: number; itemCount: number }) {
  const steps = [
    { label: 'Restaurant details', description: restaurant.address ? 'Profile is ready' : 'Add your location', done: Boolean(restaurant.address) },
    { label: 'Build categories', description: categoryCount ? `${categoryCount} organised` : 'Create the first section', done: categoryCount > 0 },
    { label: 'Add menu items', description: itemCount ? `${itemCount} dishes added` : 'Bring your menu to life', done: itemCount > 0 },
  ];
  return (
    <section className="admin-setup-card overflow-hidden rounded-2xl px-5 py-5 text-white sm:px-6">
      <div className="flex flex-col gap-5 lg:flex-row lg:items-center lg:justify-between">
        <div className="max-w-md">
          <div className="flex items-center gap-2 text-[10px] font-bold uppercase tracking-[.15em] text-white/80"><Sparkles size={13} /> Launch checklist</div>
          <h2 className="mt-1.5 text-lg font-bold tracking-[-.02em]">Your menu foundation is taking shape.</h2>
          <p className="mt-1 text-[11px] leading-relaxed text-white/80">Complete these essentials, then you’ll be ready for ordering and sales tools.</p>
        </div>
        <div className="grid flex-1 grid-cols-1 gap-3 sm:grid-cols-3 lg:max-w-[700px]">
          {steps.map((step, index) => (
            <div key={step.label} className="admin-setup-step relative flex items-center gap-3 rounded-xl border border-white/10 bg-white/[.06] p-3">
              <div className={`relative z-10 flex h-8 w-8 shrink-0 items-center justify-center rounded-full text-xs font-bold ${step.done ? 'bg-[#FFB72B] text-[#4A2F00]' : 'border border-white/30 bg-white/10 text-white'}`}>{step.done ? <Check size={15} strokeWidth={3} /> : index + 1}</div>
              <div className="min-w-0"><div className="truncate text-[11px] font-semibold">{step.label}</div><div className="mt-0.5 truncate text-[9px] text-white/70">{step.description}</div></div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}

function AddonManager({ categories, addons, selectedCategory, onCategoryChange, onCreate, onEdit, onDelete, onAvailability }: {
  categories: AdminCategory[];
  addons: AdminAddon[];
  selectedCategory: string;
  onCategoryChange: (id: string) => void;
  onCreate: () => void;
  onEdit: (addon: AdminAddon) => void;
  onDelete: (addon: AdminAddon) => void;
  onAvailability: (addon: AdminAddon, value: boolean) => void;
}) {
  const category = categories.find((entry) => entry.id === selectedCategory) ?? categories[0];
  const visibleAddons = category ? addons.filter((addon) => addon.categoryId === category.id) : [];

  return (
    <section className="admin-panel mt-6 overflow-hidden rounded-2xl">
      <div className="flex flex-col gap-4 border-b border-[#E5E7EB] px-5 py-4 sm:flex-row sm:items-center sm:justify-between">
        <div className="flex min-w-0 items-center gap-3">
          <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-xl bg-[#FFF1EB] text-[#FF4500]"><Puzzle size={19} /></div>
          <div className="min-w-0">
            <h2 className="text-[17px] font-bold tracking-[-.02em]">Category add-ons</h2>
            <p className="mt-0.5 text-[11px] text-[#6B7280]">Options are shared by every menu item in the selected category.</p>
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-2">
          {categories.length > 0 && (
            <Select value={category?.id ?? ''} onValueChange={onCategoryChange}>
              <SelectTrigger className="admin-input h-9 min-w-[180px] flex-1 rounded-lg text-xs sm:flex-none"><SelectValue placeholder="Choose category" /></SelectTrigger>
              <SelectContent>{categories.map((entry) => <SelectItem key={entry.id} value={entry.id}>{entry.name}</SelectItem>)}</SelectContent>
            </Select>
          )}
          <Button className="admin-primary h-9 gap-2 rounded-lg px-3.5 text-xs" onClick={onCreate} disabled={!category}>
            <Plus size={15} /> Add add-on
          </Button>
        </div>
      </div>

      {!category ? (
        <div className="flex min-h-[150px] flex-col items-center justify-center px-5 py-8 text-center">
          <Puzzle size={22} className="text-[#FF8C42]" />
          <h3 className="mt-2 text-sm font-bold">Create a category first</h3>
          <p className="mt-1 text-[11px] text-[#6B7280]">Add-ons belong to a menu category and will appear on all of its items.</p>
        </div>
      ) : visibleAddons.length ? (
        <div className="grid gap-3 bg-[#FAFAFA] p-4 sm:grid-cols-2 xl:grid-cols-3">
          {visibleAddons.map((addon) => (
            <article key={addon.id} className="admin-addon-card rounded-xl border border-[#E5E7EB] bg-white p-4">
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0">
                  <h3 className="truncate text-[13px] font-bold text-[#333333]">{addon.name}</h3>
                  <p className="mt-1 line-clamp-1 text-[10px] text-[#6B7280]">{addon.description || 'No description added yet.'}</p>
                </div>
                <div className="flex shrink-0 items-center gap-1.5">
                  <span className="text-[12px] font-bold text-[#FF4500]">+{formatCurrency(addon.price)}</span>
                  <ItemMenu onEdit={() => onEdit(addon)} onDelete={() => onDelete(addon)} noun="add-on" />
                </div>
              </div>
              <div className="mt-3 flex items-center justify-between border-t border-[#F3F4F6] pt-3">
                <div>
                  <div className="flex items-center gap-2"><span className={`h-2 w-2 rounded-full ${addon.available ? 'bg-[#4CAF50]' : 'bg-[#F44336]'}`} /><span className={`text-[10px] font-semibold ${addon.available ? 'text-[#2E7D32]' : 'text-[#C62828]'}`}>{addon.available ? 'Available' : 'Unavailable'}</span></div>
                  <p className="mt-1 text-[9px] text-[#9CA3AF]">Up to {addon.maxQuantity} per item</p>
                </div>
                <Switch checked={addon.available} onCheckedChange={(value) => onAvailability(addon, value)} aria-label={`Toggle ${addon.name} availability`} />
              </div>
            </article>
          ))}
        </div>
      ) : (
        <div className="flex min-h-[150px] flex-col items-center justify-center px-5 py-8 text-center">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#FFF1EB] text-[#FF4500]"><Puzzle size={18} /></div>
          <h3 className="mt-2 text-sm font-bold">No add-ons for {category.name}</h3>
          <p className="mt-1 text-[11px] text-[#6B7280]">Offer extras such as toppings, sides, sauces or size upgrades.</p>
          <Button variant="outline" className="mt-3 h-8 gap-1.5 text-[11px] hover:border-[#FF8C42] hover:text-[#FF4500]" onClick={onCreate}><Plus size={13} /> Create first add-on</Button>
        </div>
      )}
    </section>
  );
}

function CategoryRail({ categories, counts, total, selected, onSelect, onAdd, onAddItem, onEdit }: { categories: AdminCategory[]; counts: Map<string, number>; total: number; selected: string; onSelect: (id: string) => void; onAdd: () => void; onAddItem: (categoryId?: string) => void; onEdit: (category: AdminCategory) => void }) {
  return (
    <aside className="w-full border-b border-[#E5E7EB] bg-[#FAFAFA] lg:w-[230px] lg:shrink-0 lg:border-b-0 lg:border-r">
      <div className="flex items-center justify-between px-4 pb-2 pt-4">
        <span className="text-[10px] font-bold uppercase tracking-[.13em] text-[#9CA3AF]">Categories</span>
        <div className="flex items-center gap-1">
          <button onClick={() => onAddItem(selected !== 'all' ? selected : undefined)} className="flex h-6 items-center gap-1 rounded-md px-1.5 text-[11px] font-semibold text-[#FF4500] hover:bg-[#FFF1EB]" title="Add menu item" aria-label="Add menu item"><Plus size={13} /> Item</button>
          <button onClick={onAdd} className="flex h-6 w-6 items-center justify-center rounded-md text-[#6B7280] hover:bg-[#F3F4F6] hover:text-[#111827]" title="Add category" aria-label="Add category"><Plus size={14} /></button>
        </div>
      </div>
      <div className="flex gap-1 overflow-x-auto px-2 pb-3 lg:block lg:space-y-0.5 lg:overflow-visible lg:pb-5">
        <button onClick={() => onSelect('all')} className={`admin-category-button flex min-w-fit items-center gap-2 rounded-lg px-3 py-2 text-left text-[12px] font-semibold lg:w-full ${selected === 'all' ? 'active' : ''}`}>
          <Grid2X2 size={14} /><span className="flex-1">All items</span><span className="rounded bg-black/[.04] px-1.5 py-0.5 text-[9px]">{total}</span>
        </button>
        {categories.map((category) => (
          <div key={category.id} className={`group flex min-w-fit items-center rounded-lg lg:w-full ${selected === category.id ? 'admin-category-button active' : ''}`}>
            <button onClick={() => onSelect(category.id)} className={`admin-category-button flex min-w-0 flex-1 items-center gap-2 rounded-lg px-3 py-2 text-left text-[12px] font-medium ${selected === category.id ? 'active !bg-transparent !shadow-none' : ''}`}>
              <span className={`h-2 w-2 shrink-0 rounded-full ${category.enabled ? 'bg-[#FFB72B]' : 'bg-[#9CA3AF]'}`} />
              <span className="max-w-[110px] truncate">{category.name}</span>
              <span className="ml-auto rounded bg-black/[.04] px-1.5 py-0.5 text-[9px]">{counts.get(category.id) ?? 0}</span>
            </button>
            <div className="mr-1 hidden items-center gap-0.5 lg:group-hover:flex">
              <button onClick={() => onAddItem(category.id)} className="p-1 text-[#9CA3AF] hover:text-[#FF4500]" title={`Add menu item to ${category.name}`} aria-label={`Add menu item to ${category.name}`}><Plus size={13} /></button>
              <button onClick={() => onEdit(category)} className="p-1 text-[#9CA3AF] hover:text-[#FF4500]" title={`Edit ${category.name}`} aria-label={`Edit ${category.name}`}><Pencil size={12} /></button>
            </div>
          </div>
        ))}
      </div>
    </aside>
  );
}

function MenuCard({ item, category, view, onEdit, onDelete, onAvailability }: { item: AdminMenuItem; category?: AdminCategory; view: 'grid' | 'list'; onEdit: () => void; onDelete: () => void; onAvailability: (value: boolean) => void }) {
  if (view === 'list') {
    return (
      <article className="admin-menu-card flex items-center gap-3 rounded-xl p-3">
        <div className="admin-menu-image relative h-16 w-20 shrink-0 overflow-hidden rounded-lg">{item.imageUrl ? <img src={item.imageUrl} alt="" className="h-full w-full object-cover" /> : <ImageOff className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 text-[#9CA3AF]" size={18} />}</div>
        <div className="min-w-0 flex-1"><div className="flex items-center gap-2"><h3 className="truncate text-sm font-bold">{item.name}</h3>{!item.available && <span className="rounded-full bg-[#FFEBEE] px-2 py-0.5 text-[9px] font-bold text-[#C62828]">Unavailable</span>}</div><p className="mt-1 line-clamp-1 text-[11px] text-[#6B7280]">{item.description}</p><div className="mt-1.5 flex items-center gap-2 text-[10px] text-[#6B7280]"><span>{category?.name}</span><span>•</span><strong className="text-[#FF4500]">{formatCurrency(item.price)}</strong></div></div>
        <Switch checked={item.available} onCheckedChange={onAvailability} aria-label={`Toggle ${item.name} availability`} />
        <ItemMenu onEdit={onEdit} onDelete={onDelete} />
      </article>
    );
  }
  return (
    <article className="admin-menu-card rounded-xl">
      <div className="admin-menu-image relative aspect-[16/9] overflow-hidden">
        {item.imageUrl ? <img src={item.imageUrl} alt={item.name} className={`h-full w-full object-cover transition duration-500 hover:scale-[1.03] ${item.available ? '' : 'grayscale-[45%] opacity-80'}`} /> : <ImageOff className="absolute left-1/2 top-1/2 -translate-x-1/2 -translate-y-1/2 text-[#9CA3AF]" size={25} />}
        <div className="absolute left-3 top-3 z-10 rounded-full bg-[#FF4500] px-2.5 py-1 text-[9px] font-bold text-white shadow-sm">{category?.name ?? 'Uncategorised'}</div>
        <div className="absolute right-2.5 top-2.5 z-10"><ItemMenu onEdit={onEdit} onDelete={onDelete} contrast /></div>
      </div>
      <div className="p-3.5">
        <div className="flex items-start justify-between gap-3"><h3 className="line-clamp-1 text-[14px] font-bold tracking-[-.01em]">{item.name}</h3><span className="shrink-0 text-[13px] font-bold text-[#FF4500]">{formatCurrency(item.price)}</span></div>
        <p className="mt-1.5 line-clamp-2 min-h-8 text-[10px] leading-4 text-[#6B7280]">{item.description || 'No description added yet.'}</p>
        <div className="mt-3 flex items-center justify-between border-t border-[#F3F4F6] pt-3">
          <div className="flex items-center gap-2"><span className={`h-2 w-2 rounded-full ${item.available ? 'bg-[#4CAF50]' : 'bg-[#F44336]'}`} /><span className={`text-[10px] font-semibold ${item.available ? 'text-[#2E7D32]' : 'text-[#C62828]'}`}>{item.available ? 'Available' : 'Unavailable'}</span></div>
          <Switch checked={item.available} onCheckedChange={onAvailability} aria-label={`Toggle ${item.name} availability`} />
        </div>
      </div>
    </article>
  );
}

function ItemMenu({ onEdit, onDelete, contrast, noun = 'item' }: { onEdit: () => void; onDelete: () => void; contrast?: boolean; noun?: string }) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild><button className={`flex h-7 w-7 items-center justify-center rounded-lg ${contrast ? 'bg-white/90 text-[#374151] shadow-sm' : 'text-[#6B7280] hover:bg-[#F3F4F6]'}`} aria-label={`${noun === 'item' ? 'Item' : 'Add-on'} actions`}><MoreHorizontal size={16} /></button></DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-36"><DropdownMenuItem onClick={onEdit} className="gap-2 text-xs"><Pencil size={13} /> Edit {noun}</DropdownMenuItem><DropdownMenuItem onClick={onDelete} className="gap-2 text-xs text-red-600"><Trash2 size={13} /> Delete</DropdownMenuItem></DropdownMenuContent>
    </DropdownMenu>
  );
}

function EmptyRestaurant({ onCreate }: { onCreate: () => void }) {
  return <div className="admin-panel admin-empty-pattern flex min-h-[480px] flex-col items-center justify-center rounded-2xl px-5 text-center"><div className="mb-5 flex h-16 w-16 items-center justify-center rounded-2xl bg-[#FFF1EB] text-[#FF4500]"><Store size={28} /></div><h2 className="text-xl font-bold">Create your first restaurant</h2><p className="mt-2 max-w-md text-sm leading-relaxed text-[#6B7280]">Add the essentials, then build categories and menu items in a workspace designed to keep everything connected.</p><Button className="admin-primary mt-5 gap-2" onClick={onCreate}><Plus size={16} /> Create restaurant</Button></div>;
}

function OrderStatusBadge({ status }: { status: OrderStatus }) {
  return <span className={`rounded-full px-2 py-0.5 text-[9px] font-bold uppercase tracking-[.08em] ${ORDER_STATUS_STYLES[status]}`}>{ORDER_STATUS_LABELS[status]}</span>;
}

function PaymentStatusBadge({ status }: { status: PaymentStatus }) {
  return <span className={`rounded-full px-2 py-0.5 text-[9px] font-bold uppercase tracking-[.08em] ${PAYMENT_STATUS_STYLES[status]}`}>{status}</span>;
}

function OrdersSection({ orders, loading, loaded, statusFilter, onStatusFilter, typeFilter, onTypeFilter, query, onQuery, onRefresh, onSelect }: {
  orders: AdminOrder[];
  loading: boolean;
  loaded: boolean;
  statusFilter: 'all' | OrderStatus;
  onStatusFilter: (value: 'all' | OrderStatus) => void;
  typeFilter: 'all' | OrderType;
  onTypeFilter: (value: 'all' | OrderType) => void;
  query: string;
  onQuery: (value: string) => void;
  onRefresh: () => void;
  onSelect: (order: AdminOrder) => void;
}) {
  const activeCount = orders.filter((order) => order.orderStatus !== 'completed' && order.orderStatus !== 'cancelled').length;
  const pendingCount = orders.filter((order) => order.orderStatus === 'pending').length;
  const cancelledCount = orders.filter((order) => order.orderStatus === 'cancelled').length;
  const revenue = orders.filter((order) => order.orderStatus !== 'cancelled').reduce((sum, order) => sum + order.total, 0);

  return (
    <>
      <section className="mb-7 flex flex-col justify-between gap-4 md:flex-row md:items-end">
        <div>
          <div className="mb-1.5 flex items-center gap-2 text-[11px] font-bold uppercase tracking-[.14em] text-[#FF4500]"><ReceiptText size={13} /> Order management</div>
          <h1 className="text-[28px] font-bold tracking-[-.035em] text-[#333333] sm:text-[32px]">Stay ahead of every order.</h1>
          <p className="mt-1 max-w-2xl text-[13px] text-[#6B7280]">Track incoming orders, move them through preparation, and keep the kitchen flowing.</p>
        </div>
      </section>

      <section className="mb-6 grid grid-cols-2 gap-3 xl:grid-cols-4">
        <Stat icon={ReceiptText} label="Active orders" value={activeCount} note="Being prepared right now" />
        <Stat icon={Clock3} label="Pending" value={pendingCount} note="Awaiting confirmation" attention={pendingCount > 0} />
        <Stat icon={CreditCard} label="Revenue" value={formatCurrency(revenue)} note="Across listed orders" />
        <Stat icon={XCircle} label="Cancelled" value={cancelledCount} note="Across listed orders" />
      </section>

      <section className="admin-panel overflow-hidden rounded-2xl">
        <div className="flex flex-col gap-3 border-b border-[#E5E7EB] p-4 sm:flex-row sm:items-center sm:justify-between sm:px-5">
          <div>
            <h2 className="text-[17px] font-bold tracking-[-.02em]">Order feed</h2>
            <p className="mt-0.5 text-[11px] text-[#9CA3AF]">{orders.length} {orders.length === 1 ? 'order' : 'orders'} shown</p>
          </div>
          <div className="flex flex-wrap items-center gap-2">
            <div className="relative min-w-[180px] flex-1 sm:w-[220px] sm:flex-none">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-[#9CA3AF]" size={15} />
              <Input value={query} onChange={(event) => onQuery(event.target.value)} placeholder="Search customer or id" className="admin-input h-9 rounded-lg pl-9 text-xs" />
            </div>
            <Select value={statusFilter} onValueChange={(value) => onStatusFilter(value as 'all' | OrderStatus)}>
              <SelectTrigger className="admin-input h-9 w-[136px] rounded-lg text-xs"><SelectValue /></SelectTrigger>
              <SelectContent>{ORDER_STATUS_OPTIONS.map((option) => <SelectItem key={option.value} value={option.value}>{option.label}</SelectItem>)}</SelectContent>
            </Select>
            <Select value={typeFilter} onValueChange={(value) => onTypeFilter(value as 'all' | OrderType)}>
              <SelectTrigger className="admin-input h-9 w-[120px] rounded-lg text-xs"><SelectValue /></SelectTrigger>
              <SelectContent>{ORDER_TYPE_OPTIONS.map((option) => <SelectItem key={option.value} value={option.value}>{option.label}</SelectItem>)}</SelectContent>
            </Select>
            <Button variant="ghost" size="icon" className="h-9 w-9" onClick={onRefresh} disabled={loading} aria-label="Refresh orders"><RefreshCw size={15} className={loading ? 'animate-spin' : ''} /></Button>
          </div>
        </div>

        <div className="admin-empty-pattern min-h-[440px] p-4 sm:p-5">
          {loading ? (
            <div className="flex min-h-[390px] items-center justify-center"><Loader2 className="animate-spin text-[#FF4500]" size={28} /></div>
          ) : orders.length ? (
            <div className="space-y-3">
              {orders.map((order) => <OrderRow key={order.id} order={order} onSelect={() => onSelect(order)} />)}
            </div>
          ) : (
            <div className="flex min-h-[390px] flex-col items-center justify-center px-5 text-center">
              <div className="mb-4 flex h-14 w-14 items-center justify-center rounded-2xl border border-[#FED7C7] bg-white text-[#FF4500] shadow-sm"><ReceiptText size={23} /></div>
              <h3 className="text-base font-bold">{!loaded ? 'Loading orders...' : 'No orders match these filters'}</h3>
              <p className="mt-1 max-w-sm text-xs leading-relaxed text-[#6B7280]">{!loaded ? 'Connecting to backend service...' : 'Try a different status or clear the search to see more orders.'}</p>
            </div>
          )}
        </div>
        <div className="flex flex-wrap items-center justify-between gap-3 bg-[#FAFAFA] px-5 py-3 text-[10px] text-[#6B7280]">
          <span>Status updates apply to the customer order immediately.</span>
          <span className="flex items-center gap-1.5"><Clock3 size={12} /> Last synced just now</span>
        </div>
      </section>
    </>
  );
}

function OrderRow({ order, onSelect }: { order: AdminOrder; onSelect: () => void }) {
  const itemCount = order.items.reduce((sum, item) => sum + item.quantity, 0);
  return (
    <button onClick={onSelect} className="admin-order-row flex w-full flex-col gap-3 rounded-xl border border-[#E5E7EB] bg-white p-4 text-left sm:flex-row sm:items-center">
      <div className="flex min-w-0 flex-1 items-center gap-3">
        <div className={`flex h-10 w-10 shrink-0 items-center justify-center rounded-xl ${order.orderType === 'delivery' ? 'bg-[#EFF6FF] text-[#1D4ED8]' : 'bg-[#FFF1EB] text-[#FF4500]'}`}>
          {order.orderType === 'delivery' ? <Bike size={18} /> : <ShoppingBag size={18} />}
        </div>
        <div className="min-w-0">
          <div className="flex flex-wrap items-center gap-2">
            <span className="text-[13px] font-bold text-[#333333]">#{order.id.slice(0, 8)}</span>
            <OrderStatusBadge status={order.orderStatus} />
          </div>
          <div className="mt-0.5 truncate text-[11px] text-[#6B7280]">{order.customerName} • {itemCount} {itemCount === 1 ? 'item' : 'items'} • {order.orderType}</div>
        </div>
      </div>
      <div className="flex items-center justify-between gap-3 sm:justify-end">
        <PaymentStatusBadge status={order.paymentStatus} />
        <div className="text-left sm:text-right">
          <div className="text-[13px] font-bold text-[#FF4500]">{formatCurrency(order.total)}</div>
          <div className="mt-0.5 text-[10px] text-[#9CA3AF]">{formatOrderTime(order.dateCreated)}</div>
        </div>
        <ChevronRight size={16} className="shrink-0 text-[#9CA3AF]" />
      </div>
    </button>
  );
}

function OrderDetailDialog({ order, busy, onClose, onAdvance, onMarkPaid, onCancel }: {
  order: AdminOrder | null;
  busy: boolean;
  onClose: () => void;
  onAdvance: (order: AdminOrder) => void;
  onMarkPaid: (order: AdminOrder) => void;
  onCancel: (order: AdminOrder) => void;
}) {
  const nextStatus = order ? nextOrderStatus(order) : null;
  const canCancel = order ? order.orderStatus === 'pending' || order.orderStatus === 'confirmed' : false;
  const canMarkPaid = order ? order.paymentMethod === 'cash' && order.paymentStatus !== 'paid' && order.paymentStatus !== 'refunded' && order.orderStatus !== 'cancelled' : false;

  return (
    <Dialog open={Boolean(order)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-h-[92vh] overflow-y-auto border-[#E5E7EB] p-0 sm:max-w-[640px]">
        {order && (
          <>
            <DialogHeader className="border-b border-[#E5E7EB] px-6 py-5 text-left">
              <div className="flex flex-wrap items-center justify-between gap-3 pr-6">
                <div className="flex items-center gap-3">
                  <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-xl bg-[#FFF1EB] text-[#FF4500]"><ReceiptText size={18} /></div>
                  <div>
                    <DialogTitle className="text-xl tracking-[-.025em]">Order #{order.id.slice(0, 8)}</DialogTitle>
                    <DialogDescription className="text-xs">{formatOrderTime(order.dateCreated)} • {order.orderType === 'delivery' ? 'Delivery' : 'Pickup'}</DialogDescription>
                  </div>
                </div>
                <div className="flex items-center gap-1.5">
                  <OrderStatusBadge status={order.orderStatus} />
                  <PaymentStatusBadge status={order.paymentStatus} />
                </div>
              </div>
            </DialogHeader>

            <div className="space-y-5 px-6 py-5">
              <section>
                <h3 className="mb-2 text-[10px] font-bold uppercase tracking-[.13em] text-[#9CA3AF]">Customer</h3>
                <div className="rounded-xl border border-[#E5E7EB] bg-white p-3.5 text-xs">
                  <div className="text-[13px] font-bold text-[#333333]">{order.customerName}</div>
                  <div className="mt-1.5 flex flex-wrap gap-x-4 gap-y-1 text-[#6B7280]">
                    <span className="flex items-center gap-1.5"><Mail size={12} /> {order.customerEmail}</span>
                    <span className="flex items-center gap-1.5"><Phone size={12} /> {order.customerPhone}</span>
                  </div>
                </div>
              </section>

              <section>
                <h3 className="mb-2 text-[10px] font-bold uppercase tracking-[.13em] text-[#9CA3AF]">Items</h3>
                <div className="divide-y divide-[#F3F4F6] rounded-xl border border-[#E5E7EB] bg-white">
                  {order.items.map((item) => (
                    <div key={item.id} className="flex items-start justify-between gap-3 p-3.5 text-xs">
                      <div className="min-w-0">
                        <div className="font-semibold text-[#333333]">{item.quantity} × {item.menuItemName}</div>
                        {item.addons && item.addons.length > 0 && (
                          <div className="mt-1 space-y-0.5">
                            {item.addons.map((addon) => (
                              <div key={addon.id} className="flex items-center justify-between gap-3 text-[11px] text-[#6B7280]">
                                <span>+ {addon.addonName} ×{addon.quantity}</span>
                                <span className="text-[#FF4500]">+{formatCurrency(addon.addonPrice * addon.quantity * item.quantity)}</span>
                              </div>
                            ))}
                          </div>
                        )}
                        {item.specialInstructions && <div className="mt-1 text-[11px] italic text-[#9CA3AF]">&ldquo;{item.specialInstructions}&rdquo;</div>}
                      </div>
                      <div className="shrink-0 font-semibold text-[#333333]">{formatCurrency(item.menuItemPrice * item.quantity)}</div>
                    </div>
                  ))}
                </div>
              </section>

              {order.deliveryAddress && (
                <section>
                  <h3 className="mb-2 text-[10px] font-bold uppercase tracking-[.13em] text-[#9CA3AF]">Delivery address</h3>
                  <div className="rounded-xl border border-[#E5E7EB] bg-white p-3.5 text-xs text-[#374151]">
                    <div className="flex items-start gap-2">
                      <MapPin size={13} className="mt-0.5 shrink-0 text-[#FF4500]" />
                      <span>{order.deliveryAddress.street}, {order.deliveryAddress.city}, {order.deliveryAddress.state} {order.deliveryAddress.postalCode}</span>
                    </div>
                    {order.deliveryAddress.deliveryInstructions && <div className="mt-1.5 pl-5 text-[11px] text-[#6B7280]">Note: {order.deliveryAddress.deliveryInstructions}</div>}
                  </div>
                </section>
              )}

              {order.specialInstructions && (
                <section>
                  <h3 className="mb-2 text-[10px] font-bold uppercase tracking-[.13em] text-[#9CA3AF]">Order notes</h3>
                  <div className="rounded-xl border border-[#FDE68A] bg-[#FFFBEB] p-3.5 text-xs text-[#92400E]">{order.specialInstructions}</div>
                </section>
              )}

              <section className="rounded-xl border border-[#E5E7EB] bg-[#FAFAFA] p-3.5 text-xs">
                <div className="flex justify-between py-0.5 text-[#6B7280]"><span>Subtotal</span><span>{formatCurrency(order.subtotal)}</span></div>
                <div className="flex justify-between py-0.5 text-[#6B7280]"><span>Delivery fee</span><span>{formatCurrency(order.deliveryFee)}</span></div>
                <div className="flex justify-between py-0.5 text-[#6B7280]"><span>Tax</span><span>{formatCurrency(order.tax)}</span></div>
                <div className="mt-1.5 flex justify-between border-t border-[#E5E7EB] pt-2 text-[13px] font-bold text-[#333333]"><span>Total</span><span>{formatCurrency(order.total)}</span></div>
              </section>

              <div className="flex flex-wrap items-center gap-2 text-[11px] text-[#6B7280]">
                {order.paymentMethod === 'cash' ? <Banknote size={13} /> : <CreditCard size={13} />}
                <span>{order.paymentMethod === 'cash' ? 'Cash on handover' : 'Card payment'}</span>
                {order.stripePaymentIntentId && <span className="truncate text-[#9CA3AF]">• {order.stripePaymentIntentId}</span>}
              </div>
            </div>

            {(canCancel || canMarkPaid || nextStatus) && (
              <div className="flex flex-wrap items-center justify-between gap-2 border-t border-[#E5E7EB] bg-[#FAFAFA] px-6 py-4">
                <div>
                  {canCancel && (
                    <Button variant="outline" className="h-9 gap-2 border-[#FECACA] text-xs text-red-600 hover:bg-[#FEF2F2] hover:text-red-700" onClick={() => onCancel(order)} disabled={busy}>
                      <XCircle size={14} /> Cancel order
                    </Button>
                  )}
                </div>
                <div className="flex items-center gap-2">
                  {canMarkPaid && (
                    <Button variant="outline" className="h-9 gap-2 text-xs" onClick={() => onMarkPaid(order)} disabled={busy}>
                      <Banknote size={14} /> Mark as paid
                    </Button>
                  )}
                  {nextStatus && (
                    <Button className="admin-primary h-9 gap-2 text-xs" onClick={() => onAdvance(order)} disabled={busy}>
                      {busy && <Loader2 size={14} className="animate-spin" />}
                      {NEXT_STATUS_LABELS[nextStatus]} <ArrowRight size={14} />
                    </Button>
                  )}
                </div>
              </div>
            )}
          </>
        )}
      </DialogContent>
    </Dialog>
  );
}

function EditorDialog({ editor, workspace, onClose, onSave }: { editor: EditorState; workspace: AdminWorkspace | null; onClose: () => void; onSave: (kind: 'restaurant' | 'category' | 'item' | 'addon' | 'promotion', input: RestaurantInput | CategoryInput | MenuItemInput | AddonInput | PromotionInput, existingId?: string) => Promise<void> }) {
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const kind = editor?.kind;
  const existing = editor?.value;
  const addonCategoryId = editor?.kind === 'addon' ? (editor.value as AdminAddon | undefined)?.categoryId ?? editor.categoryId ?? workspace?.categories[0]?.id ?? '' : '';
  const itemCategoryId = editor?.kind === 'item' ? (editor.value as AdminMenuItem | undefined)?.categoryId ?? editor.categoryId ?? workspace?.categories[0]?.id ?? '' : '';

  // Address search & resolution state for restaurant profile
  const [addressInput, setAddressInput] = useState<string>((existing as AdminRestaurant | undefined)?.address ?? '');
  const [lat, setLat] = useState<number | null>((existing as AdminRestaurant | undefined)?.latitude ?? null);
  const [lon, setLon] = useState<number | null>((existing as AdminRestaurant | undefined)?.longitude ?? null);
  const [isGeocoding, setIsGeocoding] = useState<boolean>(false);
  const [geoResults, setGeoResults] = useState<{ display_name: string; lat: string; lon: string }[]>([]);
  const [geoError, setGeoError] = useState<string>('');
  // Promotion discount type state for max attribute & validation
  const [promoDiscountType, setPromoDiscountType] = useState<'percentage' | 'fixed_amount'>(
    (existing as AdminPromotion | undefined)?.discountType ?? 'percentage'
  );

  useEffect(() => {
    if (editor?.kind === 'restaurant') {
      const rest = editor.value as AdminRestaurant | undefined;
      setAddressInput(rest?.address ?? '');
      setLat(rest?.latitude ?? null);
      setLon(rest?.longitude ?? null);
      setGeoResults([]);
      setGeoError('');
    }
    if (editor?.kind === 'promotion') {
      const promo = editor.value as AdminPromotion | undefined;
      setPromoDiscountType(promo?.discountType ?? 'percentage');
    }
  }, [editor]);

  const handleResolveAddress = async () => {
    const query = addressInput.trim();
    if (!query) return;
    setIsGeocoding(true);
    setGeoError('');
    setGeoResults([]);

    try {
      const isPostalCode = /^\d{6}$/.test(query);
      const geoUrl = new URL('https://nominatim.openstreetmap.org/search');
      geoUrl.searchParams.set('q', isPostalCode ? `${query}, Singapore` : query);
      geoUrl.searchParams.set('format', 'jsonv2');
      geoUrl.searchParams.set('addressdetails', '1');
      geoUrl.searchParams.set('countrycodes', 'sg');
      geoUrl.searchParams.set('accept-language', 'en');
      geoUrl.searchParams.set('limit', '5');

      const res = await fetch(geoUrl.toString(), {
        headers: { Accept: 'application/json', 'Accept-Language': 'en' },
      });

      if (res.ok) {
        const data = await res.json();
        if (Array.isArray(data) && data.length > 0) {
          if (data.length === 1) {
            const chosenLat = parseFloat(data[0].lat);
            const chosenLon = parseFloat(data[0].lon);
            setLat(chosenLat);
            setLon(chosenLon);
            setAddressInput(data[0].display_name);
            toast.success('Address verified successfully');
          } else {
            setGeoResults(data);
          }
        } else {
          setGeoError('No matching address found. Please refine the address and try again.');
        }
      } else {
        setGeoError('Failed to contact address lookup service.');
      }
    } catch (e) {
      console.error('Error resolving address:', e);
      setGeoError('Failed to resolve address. Please try again.');
    } finally {
      setIsGeocoding(false);
    }
  };

  const handleSelectResult = (item: { display_name: string; lat: string; lon: string }) => {
    const chosenLat = parseFloat(item.lat);
    const chosenLon = parseFloat(item.lon);
    setLat(chosenLat);
    setLon(chosenLon);
    setAddressInput(item.display_name);
    setGeoResults([]);
    setGeoError('');
    toast.success('Address selected');
  };

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault(); setSaving(true); setError('');
    const data = new FormData(event.currentTarget);
    try {
      if (kind === 'restaurant') {
        const addressStr = addressInput.trim();
        let latitudeVal: number | null = lat;
        let longitudeVal: number | null = lon;

        if (addressStr && (latitudeVal === null || longitudeVal === null)) {
          try {
            const isPostalCode = /^\d{6}$/.test(addressStr);
            const geoUrl = new URL('https://nominatim.openstreetmap.org/search');
            geoUrl.searchParams.set('q', isPostalCode ? `${addressStr}, Singapore` : addressStr);
            geoUrl.searchParams.set('format', 'jsonv2');
            geoUrl.searchParams.set('addressdetails', '1');
            geoUrl.searchParams.set('countrycodes', 'sg');
            geoUrl.searchParams.set('accept-language', 'en');
            geoUrl.searchParams.set('limit', '1');

            const res = await fetch(geoUrl.toString(), {
              headers: { Accept: 'application/json', 'Accept-Language': 'en' },
            });
            if (res.ok) {
              const geoData = await res.json();
              if (Array.isArray(geoData) && geoData.length > 0) {
                latitudeVal = parseFloat(geoData[0].lat);
                longitudeVal = parseFloat(geoData[0].lon);
              }
            }
          } catch (e) {
            console.error('Failed to resolve address coordinates:', e);
          }
        }

        const maxDeliveryDistanceKm = String(data.get('maxDeliveryDistanceKm') ?? '').trim();
        const minSpendStr = String(data.get('minSpend') ?? '').trim();
        const taxRatePct = String(data.get('taxRatePct') ?? '').trim();
        const taxRateVal = taxRatePct === '' ? 0.10 : Number(taxRatePct) / 100;
        await onSave(kind, {
          name: String(data.get('name')),
          description: String(data.get('description')),
          address: addressStr,
          phone: String(data.get('phone')),
          email: String(data.get('email')),
          imageUrl: String(data.get('imageUrl')),
          enabled: data.get('enabled') === 'true',
          latitude: latitudeVal,
          longitude: longitudeVal,
          maxDeliveryDistanceKm: maxDeliveryDistanceKm === '' ? 0 : Number(maxDeliveryDistanceKm),
          minSpend: minSpendStr === '' ? 0 : Number(minSpendStr),
          taxRate: taxRateVal,
        }, existing?.id);
      }
      if (kind === 'category' && workspace) await onSave(kind, { name: String(data.get('name')), description: String(data.get('description')), restaurantId: workspace.restaurant.id }, existing?.id);
      if (kind === 'item' && workspace) await onSave(kind, { name: String(data.get('name')), description: String(data.get('description')), price: Number(data.get('price')), categoryId: String(data.get('categoryId')), restaurantId: workspace.restaurant.id, imageUrl: String(data.get('imageUrl')) }, existing?.id);
      if (kind === 'addon' && workspace) await onSave(kind, { name: String(data.get('name')), description: String(data.get('description')), price: Number(data.get('price')), maxQuantity: Number(data.get('maxQuantity')), categoryId: String(data.get('categoryId')), restaurantId: workspace.restaurant.id }, existing?.id);
      if (kind === 'promotion') {
        const code = String(data.get('code')).trim();
        const discountType = String(data.get('discountType')) as 'percentage' | 'fixed_amount';
        const discountVal = Number(data.get('discountValue'));
        const minOrder = Number(data.get('minOrderAmount') ?? 0);
        const maxDiscount = data.get('maxDiscountAmount') ? Number(data.get('maxDiscountAmount')) : null;
        const usageLim = data.get('usageLimit') ? Number(data.get('usageLimit')) : null;

        if (discountType === 'percentage' && discountVal > 100) {
          setError('Percentage discount cannot exceed 100%');
          setSaving(false);
          return;
        }

        await onSave(kind, {
          code,
          name: String(data.get('name')),
          description: String(data.get('description')),
          discountType,
          discountValue: discountVal,
          minOrderAmount: minOrder,
          maxDiscountAmount: maxDiscount,
          usageLimit: usageLim,
          enabled: (existing as AdminPromotion | undefined)?.enabled ?? true,
        }, existing?.id);
      }
    } catch (submitError) { setError(submitError instanceof Error ? submitError.message : 'Could not save changes'); }
    finally { setSaving(false); }
  };

  const titles = { restaurant: existing ? 'Edit restaurant' : 'Create a restaurant', category: existing ? 'Edit category' : 'Create a category', item: existing ? 'Edit menu item' : 'Create a menu item', addon: existing ? 'Edit add-on' : 'Create an add-on', promotion: existing ? 'Edit promotion' : 'Create promotion' };
  const descriptions = { restaurant: 'The profile guests see across your storefront.', category: 'Group related items so your menu is easy to browse.', item: 'Add the details your team and guests need.', addon: 'Offer an optional extra across every item in a category.', promotion: 'Set up discount codes and redemption rules.' };

  return (
    <Dialog open={Boolean(editor)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-h-[92vh] overflow-y-auto border-[#E5E7EB] p-0 sm:max-w-[560px]">
        {kind && <form onSubmit={submit}>
          <DialogHeader className="border-b border-[#E5E7EB] px-6 py-5 text-left"><div className="mb-2 flex h-9 w-9 items-center justify-center rounded-xl bg-[#FFF1EB] text-[#FF4500]">{kind === 'restaurant' ? <Building2 size={18} /> : kind === 'category' ? <Boxes size={18} /> : kind === 'addon' ? <Puzzle size={18} /> : kind === 'promotion' ? <Tag size={18} /> : <UtensilsCrossed size={18} />}</div><DialogTitle className="text-xl tracking-[-.025em]">{titles[kind]}</DialogTitle><DialogDescription className="text-xs">{descriptions[kind]}</DialogDescription></DialogHeader>
          <div className="space-y-4 px-6 py-5">
            {kind !== 'promotion' && <Field label={kind === 'restaurant' ? 'Restaurant name' : kind === 'category' ? 'Category name' : kind === 'addon' ? 'Add-on name' : 'Item name'} htmlFor="name" required hint="3–100 characters"><Input id="name" name="name" defaultValue={existing?.name ?? ''} required minLength={3} maxLength={100} placeholder={kind === 'restaurant' ? 'e.g. Juniper Kitchen' : kind === 'category' ? 'e.g. Seasonal plates' : kind === 'addon' ? 'e.g. Extra avocado' : 'e.g. Garden harvest bowl'} className="admin-input" /></Field>}
            {kind !== 'promotion' && <Field label="Description" htmlFor="description" hint="Recommended"><Textarea id="description" name="description" defaultValue={existing?.description ?? ''} rows={3} placeholder="Add a concise, useful description" className="admin-input resize-none" /></Field>}
            {kind === 'restaurant' && <>
              <Field label="Address" htmlFor="address" required hint="Enter street or 6-digit postal code">
                <div className="flex gap-2">
                  <Input
                    id="address"
                    name="address"
                    value={addressInput}
                    onChange={(e) => {
                      setAddressInput(e.target.value);
                    }}
                    onKeyDown={(e) => {
                      if (e.key === 'Enter') {
                        e.preventDefault();
                        handleResolveAddress();
                      }
                    }}
                    required
                    placeholder="Street, city and postal code"
                    className="admin-input flex-1"
                  />
                  <Button
                    type="button"
                    variant="outline"
                    onClick={handleResolveAddress}
                    disabled={isGeocoding || !addressInput.trim()}
                    className="border-[#E5E7EB] bg-white text-xs font-medium text-[#374151] hover:bg-[#F9FAFB] hover:text-[#111827] shrink-0"
                  >
                    {isGeocoding ? <Loader2 size={14} className="animate-spin mr-1" /> : <Search size={14} className="mr-1" />}
                    Search
                  </Button>
                </div>
              </Field>

              {/* Geo Candidate Selection */}
              {geoResults.length > 0 && (
                <div className="rounded-xl border border-[#E5E7EB] bg-white p-2 shadow-sm space-y-1">
                  <div className="px-2 py-1 text-[11px] font-semibold text-[#6B7280]">Select matching address:</div>
                  <div className="max-h-40 overflow-y-auto divide-y divide-[#F3F4F6]">
                    {geoResults.map((item, idx) => (
                      <button
                        key={idx}
                        type="button"
                        onClick={() => handleSelectResult(item)}
                        className="w-full text-left px-3 py-2 text-xs text-[#374151] hover:bg-[#FFF5F0] hover:text-[#FF4500] transition-colors rounded-lg flex items-center justify-between"
                      >
                        <span className="truncate mr-2">{item.display_name}</span>
                        <span className="shrink-0 text-[10px] text-[#9CA3AF]">Lat: {parseFloat(item.lat).toFixed(4)}, Lon: {parseFloat(item.lon).toFixed(4)}</span>
                      </button>
                    ))}
                  </div>
                </div>
              )}

              {geoError && (
                <div className="flex items-center gap-1.5 text-xs text-[#DC2626] bg-[#FEF2F2] border border-[#FCA5A5] rounded-lg p-2.5">
                  <CircleAlert size={14} className="shrink-0" />
                  <span>{geoError}</span>
                </div>
              )}

              {/* Address Resolution Status Indicator */}
              {lat !== null && lon !== null ? (
                <div className="flex items-center justify-between rounded-xl border border-[#D1FAE5] bg-[#ECFDF5] px-3.5 py-2 text-xs text-[#065F46]">
                  <div className="flex items-center gap-2">
                    <div className="flex h-5 w-5 items-center justify-center rounded-full bg-[#10B981] text-white">
                      <Check size={12} />
                    </div>
                    <div>
                      <span className="font-semibold text-xs">Address verified for delivery</span>
                    </div>
                  </div>
                  <span className="text-[10px] uppercase font-bold tracking-wider text-[#059669] bg-[#A7F3D0] px-2 py-0.5 rounded-full">Verified</span>
                </div>
              ) : (
                <div className="flex items-center gap-2 rounded-xl border border-[#FDE68A] bg-[#FFFBEB] px-3.5 py-2 text-xs text-[#92400E]">
                  <CircleAlert size={14} className="shrink-0 text-[#D97706]" />
                  <span>Click <strong>Search</strong> to verify this address for delivery.</span>
                </div>
              )}

              <div className="grid gap-4 sm:grid-cols-2"><Field label="Phone" htmlFor="phone" required><Input id="phone" name="phone" defaultValue={(existing as AdminRestaurant | undefined)?.phone ?? ''} required placeholder="+65 6123 4567" className="admin-input" /></Field><Field label="Email" htmlFor="email" required><Input id="email" name="email" type="email" defaultValue={(existing as AdminRestaurant | undefined)?.email ?? ''} required placeholder="hello@restaurant.com" className="admin-input" /></Field></div>
              <Field label="Cover image URL" htmlFor="imageUrl" hint="Optional"><Input id="imageUrl" name="imageUrl" type="url" defaultValue={(existing as AdminRestaurant | undefined)?.imageUrl ?? ''} placeholder="https://images.unsplash.com/restaurant.jpg" className="admin-input" /></Field>
              <input type="hidden" name="latitude" value={lat ?? ''} />
              <input type="hidden" name="longitude" value={lon ?? ''} />
              <div className="grid gap-4 sm:grid-cols-3">
                <Field label="Max delivery distance (km)" htmlFor="maxDeliveryDistanceKm" hint="0 = unlimited"><Input id="maxDeliveryDistanceKm" name="maxDeliveryDistanceKm" type="number" min="0" step="0.1" defaultValue={(existing as AdminRestaurant | undefined)?.maxDeliveryDistanceKm ?? 0} placeholder="0" className="admin-input" /></Field>
                <Field label="Minimum spend ($)" htmlFor="minSpend" hint="0 = no minimum"><Input id="minSpend" name="minSpend" type="number" min="0" step="0.01" defaultValue={(existing as AdminRestaurant | undefined)?.minSpend ?? 0} placeholder="0.00" className="admin-input" /></Field>
                <Field label="Tax rate (%)" htmlFor="taxRatePct" hint="e.g. 10 = 10%"><Input id="taxRatePct" name="taxRatePct" type="number" min="0" max="100" step="0.1" defaultValue={((existing as AdminRestaurant | undefined)?.taxRate ?? 0.10) * 100} placeholder="10" className="admin-input" /></Field>
              </div>
              <Field label="Storefront Status" htmlFor="enabled" hint="Live allows guests to place online orders">
                <div className="flex items-center gap-3 pt-1">
                  <Switch
                    id="enabled"
                    name="enabled"
                    defaultChecked={(existing as AdminRestaurant | undefined)?.enabled ?? true}
                    value="true"
                  />
                  <span className="text-xs font-semibold text-[#374151]">
                    {(existing as AdminRestaurant | undefined)?.enabled ?? true ? 'Live (Storefront enabled)' : 'Paused (Storefront paused)'}
                  </span>
                </div>
              </Field>
            </>}
            {kind === 'item' && workspace && <>
              <div className="grid gap-4 sm:grid-cols-2"><Field label="Price" htmlFor="price" required><div className="relative"><span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-[#6B7280]">$</span><Input id="price" name="price" type="number" min="0.01" step="0.01" defaultValue={(existing as AdminMenuItem | undefined)?.price ?? ''} required placeholder="0.00" className="admin-input pl-7" /></div></Field><Field label="Category" htmlFor="categoryId" required><Select name="categoryId" defaultValue={itemCategoryId}><SelectTrigger className="admin-input"><SelectValue placeholder="Choose category" /></SelectTrigger><SelectContent>{workspace.categories.map((category) => <SelectItem key={category.id} value={category.id}>{category.name}</SelectItem>)}</SelectContent></Select></Field></div>
              <Field label="Dish image URL" htmlFor="imageUrl" hint="Optional"><Input id="imageUrl" name="imageUrl" type="url" defaultValue={(existing as AdminMenuItem | undefined)?.imageUrl ?? ''} placeholder="https://images.unsplash.com/dish.jpg" className="admin-input" /></Field>
            </>}
            {kind === 'addon' && workspace && <>
              {existing ? (
                <div className="rounded-lg border border-[#F3E1D9] bg-[#FFF7F3] px-3 py-2.5 text-[11px] text-[#6B4E3D]">
                  Shared across <strong>{workspace.categories.find((category) => category.id === addonCategoryId)?.name ?? 'this category'}</strong>
                  <input type="hidden" name="categoryId" value={addonCategoryId} />
                </div>
              ) : (
                <Field label="Category" htmlFor="categoryId" required><Select name="categoryId" defaultValue={addonCategoryId}><SelectTrigger className="admin-input"><SelectValue placeholder="Choose category" /></SelectTrigger><SelectContent>{workspace.categories.map((category) => <SelectItem key={category.id} value={category.id}>{category.name}</SelectItem>)}</SelectContent></Select></Field>
              )}
              <div className="grid gap-4 sm:grid-cols-2">
                <Field label="Additional price" htmlFor="price" required><div className="relative"><span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-[#6B7280]">$</span><Input id="price" name="price" type="number" min="0.01" step="0.01" defaultValue={(existing as AdminAddon | undefined)?.price ?? ''} required placeholder="0.00" className="admin-input pl-7" /></div></Field>
                <Field label="Maximum quantity" htmlFor="maxQuantity" required hint="Per order item"><Input id="maxQuantity" name="maxQuantity" type="number" min="1" max="20" step="1" defaultValue={(existing as AdminAddon | undefined)?.maxQuantity ?? 1} required className="admin-input" /></Field>
              </div>
            </>}
            {kind === 'promotion' && (
              <>
                <div className="grid gap-4 sm:grid-cols-2">
                  <Field label="Promo Code" htmlFor="code" required hint="e.g. WELCOME10">
                    <Input
                      id="code"
                      name="code"
                      defaultValue={(existing as AdminPromotion | undefined)?.code ?? ''}
                      required
                      placeholder="e.g. SAVE20"
                      className="admin-input uppercase"
                    />
                  </Field>
                  <Field label="Campaign Name" htmlFor="name" required hint="3–100 characters">
                    <Input
                      id="name"
                      name="name"
                      defaultValue={existing?.name ?? ''}
                      required
                      minLength={3}
                      maxLength={100}
                      placeholder="e.g. Save 20% on First Order"
                      className="admin-input"
                    />
                  </Field>
                </div>
                <Field label="Description" htmlFor="description" hint="Optional">
                  <Textarea
                    id="description"
                    name="description"
                    defaultValue={existing?.description ?? ''}
                    rows={2}
                    placeholder="Short description for campaign details"
                    className="admin-input resize-none"
                  />
                </Field>
                <div className="grid gap-4 sm:grid-cols-2">
                  <Field label="Discount Type" htmlFor="discountType" required>
                    <Select
                      name="discountType"
                      value={promoDiscountType}
                      onValueChange={(val) => setPromoDiscountType(val as 'percentage' | 'fixed_amount')}
                    >
                      <SelectTrigger className="admin-input">
                        <SelectValue placeholder="Select type" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="percentage">Percentage (%)</SelectItem>
                        <SelectItem value="fixed_amount">Fixed Amount ($)</SelectItem>
                      </SelectContent>
                    </Select>
                  </Field>
                  <Field label="Discount Value" htmlFor="discountValue" required hint={promoDiscountType === 'percentage' ? "Percentage (1–100%)" : "$ amount"}>
                    <Input
                      id="discountValue"
                      name="discountValue"
                      type="number"
                      min="0.01"
                      max={promoDiscountType === 'percentage' ? '100' : undefined}
                      step="0.01"
                      defaultValue={(existing as AdminPromotion | undefined)?.discountValue ?? ''}
                      required
                      placeholder={promoDiscountType === 'percentage' ? 'e.g. 15' : 'e.g. 5.00'}
                      className="admin-input"
                    />
                  </Field>
                </div>
                <div className="grid gap-4 sm:grid-cols-3">
                  <Field label="Min Order ($)" htmlFor="minOrderAmount" hint="0 for none">
                    <Input
                      id="minOrderAmount"
                      name="minOrderAmount"
                      type="number"
                      min="0"
                      step="0.01"
                      defaultValue={(existing as AdminPromotion | undefined)?.minOrderAmount ?? 0}
                      placeholder="0.00"
                      className="admin-input"
                    />
                  </Field>
                  <Field label="Max Discount ($)" htmlFor="maxDiscountAmount" hint="Cap for % off">
                    <Input
                      id="maxDiscountAmount"
                      name="maxDiscountAmount"
                      type="number"
                      min="0.01"
                      step="0.01"
                      defaultValue={(existing as AdminPromotion | undefined)?.maxDiscountAmount ?? ''}
                      placeholder="Optional cap"
                      className="admin-input"
                    />
                  </Field>
                  <Field label="Usage Limit" htmlFor="usageLimit" hint="Max redemptions">
                    <Input
                      id="usageLimit"
                      name="usageLimit"
                      type="number"
                      min="1"
                      step="1"
                      defaultValue={(existing as AdminPromotion | undefined)?.usageLimit ?? ''}
                      placeholder="Unlimited"
                      className="admin-input"
                    />
                  </Field>
                </div>
              </>
            )}
            {error && <div className="rounded-lg bg-red-50 px-3 py-2.5 text-xs text-red-700">{error}</div>}
          </div>
          <div className="flex items-center justify-end gap-2 border-t border-[#E5E7EB] bg-[#FAFAFA] px-6 py-4"><Button type="button" variant="outline" onClick={onClose} className="h-9 text-xs">Cancel</Button><Button type="submit" disabled={saving} className="admin-primary h-9 min-w-[112px] gap-2 text-xs">{saving && <Loader2 size={14} className="animate-spin" />}{existing ? 'Save changes' : kind === 'addon' ? 'Create add-on' : `Create ${kind}`}</Button></div>
        </form>}
      </DialogContent>
    </Dialog>
  );
}

function PromotionsSection({
  promotions,
  loading,
  query,
  onQueryChange,
  onCreate,
  onEdit,
  onToggleEnabled,
  onDelete,
  onRefresh,
}: {
  promotions: AdminPromotion[];
  loading: boolean;
  query: string;
  onQueryChange: (query: string) => void;
  onCreate: () => void;
  onEdit: (promo: AdminPromotion) => void;
  onToggleEnabled: (promo: AdminPromotion, enabled: boolean) => void;
  onDelete: (promo: AdminPromotion) => void;
  onRefresh: () => void;
}) {
  const filtered = useMemo(() => {
    const q = query.trim().toLowerCase();
    if (!q) return promotions;
    return promotions.filter(
      (p) => p.code.toLowerCase().includes(q) || p.name.toLowerCase().includes(q) || (p.description && p.description.toLowerCase().includes(q))
    );
  }, [promotions, query]);

  return (
    <div className="space-y-6">
      <section className="flex flex-col justify-between gap-4 md:flex-row md:items-end">
        <div>
          <div className="mb-1.5 flex items-center gap-2 text-[11px] font-bold uppercase tracking-[.14em] text-[#FF4500]">
            <Tag size={13} /> Promotion Campaigns
          </div>
          <h1 className="text-[28px] font-bold tracking-[-.035em] text-[#333333] sm:text-[32px]">Drive orders with discounts.</h1>
          <p className="mt-1 max-w-2xl text-[13px] text-[#6B7280]">
            Create promo codes, set minimum order rules, and track campaign redemptions.
          </p>
        </div>
        <div className="flex gap-2">
          <Button variant="outline" className="h-9 gap-2 rounded-lg border-[#E5E7EB] bg-white text-xs" onClick={onRefresh}>
            <RefreshCw size={14} className={loading ? 'animate-spin' : ''} /> Refresh
          </Button>
          <Button className="admin-primary h-9 gap-2 rounded-lg text-xs" onClick={onCreate}>
            <Plus size={15} /> Create promo code
          </Button>
        </div>
      </section>

      <section className="admin-panel overflow-hidden rounded-2xl">
        <div className="flex items-center justify-between border-b border-[#E5E7EB] p-4 sm:px-5">
          <div className="relative min-w-[200px] flex-1 sm:w-[280px] sm:flex-none">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-[#9CA3AF]" size={15} />
            <Input
              value={query}
              onChange={(e) => onQueryChange(e.target.value)}
              placeholder="Search by code or name"
              className="admin-input h-9 rounded-lg pl-9 text-xs"
            />
          </div>
          <div className="text-xs text-[#6B7280]">
            {filtered.length} {filtered.length === 1 ? 'campaign' : 'campaigns'}
          </div>
        </div>

        {loading ? (
          <div className="flex min-h-[300px] items-center justify-center">
            <Loader2 className="animate-spin text-[#FF4500]" size={28} />
          </div>
        ) : filtered.length === 0 ? (
          <div className="flex min-h-[300px] flex-col items-center justify-center p-8 text-center">
            <div className="mb-3 flex h-12 w-12 items-center justify-center rounded-2xl bg-[#FFF1EB] text-[#FF4500]">
              <Tag size={20} />
            </div>
            <h3 className="text-base font-bold text-[#111827]">No promo campaigns found</h3>
            <p className="mt-1 text-xs text-[#6B7280]">Create your first promotion code to offer customer discounts.</p>
            <Button className="admin-primary mt-4 h-9 gap-2 text-xs" onClick={onCreate}>
              <Plus size={15} /> Create promo code
            </Button>
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full text-left text-xs text-[#374151]">
              <thead className="bg-[#FAFAFA] text-[11px] uppercase tracking-wider text-[#6B7280] border-b border-[#E5E7EB]">
                <tr>
                  <th className="px-5 py-3.5 font-semibold">Code & Name</th>
                  <th className="px-5 py-3.5 font-semibold">Discount</th>
                  <th className="px-5 py-3.5 font-semibold">Min Order</th>
                  <th className="px-5 py-3.5 font-semibold">Usage</th>
                  <th className="px-5 py-3.5 font-semibold">Status</th>
                  <th className="px-5 py-3.5 font-semibold text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-[#E5E7EB] bg-white">
                {filtered.map((promo) => (
                  <tr key={promo.id} className="hover:bg-[#F9FAFB] transition-colors">
                    <td className="px-5 py-4">
                      <div className="font-bold text-sm text-[#111827] flex items-center gap-2">
                        <span className="font-mono bg-[#FFF1EB] text-[#FF4500] px-2 py-0.5 rounded border border-[#FFD8C8]">
                          {promo.code}
                        </span>
                        <span>{promo.name}</span>
                      </div>
                      {promo.description && (
                        <div className="text-[11px] text-[#6B7280] mt-1">{promo.description}</div>
                      )}
                    </td>
                    <td className="px-5 py-4 font-semibold text-[#111827]">
                      {promo.discountType === 'percentage'
                        ? `${promo.discountValue}% off${promo.maxDiscountAmount ? ` (cap $${promo.maxDiscountAmount.toFixed(2)})` : ''}`
                        : `$${promo.discountValue.toFixed(2)} off`}
                    </td>
                    <td className="px-5 py-4 text-[#6B7280]">
                      {promo.minOrderAmount > 0 ? `$${promo.minOrderAmount.toFixed(2)}` : 'None'}
                    </td>
                    <td className="px-5 py-4 text-[#6B7280]">
                      <div className="font-medium text-[#111827]">{promo.usageCount} used</div>
                      {promo.usageLimit ? (
                        <div className="text-[10px]">Limit: {promo.usageLimit}</div>
                      ) : (
                        <div className="text-[10px] text-[#9CA3AF]">Unlimited</div>
                      )}
                    </td>
                    <td className="px-5 py-4">
                      <div className="flex items-center gap-2">
                        <Switch
                          checked={promo.enabled}
                          onCheckedChange={(val) => onToggleEnabled(promo, val)}
                        />
                        <span className={`text-[11px] font-semibold ${promo.enabled ? 'text-green-700' : 'text-gray-500'}`}>
                          {promo.enabled ? 'Active' : 'Inactive'}
                        </span>
                      </div>
                    </td>
                    <td className="px-5 py-4 text-right">
                      <div className="flex items-center justify-end gap-1">
                        <Button variant="ghost" size="icon" className="h-8 w-8 text-[#6B7280] hover:text-[#FF4500]" onClick={() => onEdit(promo)}>
                          <Pencil size={14} />
                        </Button>
                        <Button variant="ghost" size="icon" className="h-8 w-8 text-[#6B7280] hover:text-red-600" onClick={() => onDelete(promo)}>
                          <Trash2 size={14} />
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>
    </div>
  );
}
