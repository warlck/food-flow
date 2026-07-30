import { FormEvent, useCallback, useEffect, useMemo, useState } from 'react';
import {
  BarChart3, Bell, BookOpen, Boxes, Building2, Check, ChevronDown, ChevronRight,
  ChefHat, CircleAlert, Clock3, Grid2X2, HelpCircle, ImageOff, LayoutDashboard, List,
  Loader2, MapPin, Menu, MoreHorizontal, PackageCheck, Pencil, Plus, ReceiptText,
  Puzzle, RefreshCw, Search, Settings, ShoppingBag, Sparkles, Store, Trash2, UtensilsCrossed,
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
  AddonInput, AdminAddon, AdminCategory, AdminMenuItem, AdminRestaurant, AdminWorkspace,
  CategoryInput, MenuItemInput, RestaurantInput, adminApi,
} from '@/lib/admin-api';
import './Admin.css';

type EditorState =
  | { kind: 'restaurant'; value?: AdminRestaurant }
  | { kind: 'category'; value?: AdminCategory }
  | { kind: 'item'; value?: AdminMenuItem }
  | { kind: 'addon'; value?: AdminAddon; categoryId?: string }
  | null;

const NAME_PATTERN = /^[\p{L}\p{N}' -]{3,100}$/u;

const navItems = [
  { icon: LayoutDashboard, label: 'Overview' },
  { icon: BookOpen, label: 'Menu & inventory', active: true },
  { icon: ReceiptText, label: 'Orders', soon: true },
  { icon: BarChart3, label: 'Sales & insights', soon: true },
  { icon: Settings, label: 'Settings' },
];

const demoRestaurant: AdminRestaurant = {
  id: 'demo-restaurant', name: 'Juniper & Grain', description: 'Seasonal plates and thoughtful pantry staples.',
  address: '28 Greenwood Avenue, Singapore', phone: '+65 6123 7788', email: 'hello@junipergrain.co',
  imageUrl: 'https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?auto=format&fit=crop&w=1200&q=80',
  enabled: true, dateCreated: '2026-07-21T08:00:00Z', dateUpdated: '2026-07-28T08:00:00Z',
};

const demoCategories: AdminCategory[] = [
  { id: 'cat-breakfast', name: 'All day breakfast', description: 'Bright starts and slow mornings', restaurantId: demoRestaurant.id, enabled: true },
  { id: 'cat-bowls', name: 'Harvest bowls', description: 'Wholesome seasonal bowls', restaurantId: demoRestaurant.id, enabled: true },
  { id: 'cat-drinks', name: 'Drinks', description: 'Coffee, tonics and juices', restaurantId: demoRestaurant.id, enabled: true },
  { id: 'cat-bakes', name: 'Bakes', description: 'Made fresh each morning', restaurantId: demoRestaurant.id, enabled: true },
];

const demoItems: AdminMenuItem[] = [
  { id: 'item-1', name: 'Sourdough garden toast', description: 'Whipped ricotta, heirloom tomatoes, basil oil and toasted seeds.', price: 18, categoryId: 'cat-breakfast', restaurantId: demoRestaurant.id, imageUrl: 'https://images.unsplash.com/photo-1541519227354-08fa5d50c44d?auto=format&fit=crop&w=900&q=80', available: true },
  { id: 'item-2', name: 'Miso mushroom eggs', description: 'Soft eggs, roasted mushrooms, miso butter and sourdough.', price: 21, categoryId: 'cat-breakfast', restaurantId: demoRestaurant.id, imageUrl: 'https://images.unsplash.com/photo-1525351484163-7529414344d8?auto=format&fit=crop&w=900&q=80', available: true },
  { id: 'item-3', name: 'Green goddess bowl', description: 'Charred greens, avocado, herbed grains, edamame and tahini.', price: 19.5, categoryId: 'cat-bowls', restaurantId: demoRestaurant.id, imageUrl: 'https://images.unsplash.com/photo-1512621776951-a57141f2eefd?auto=format&fit=crop&w=900&q=80', available: true },
  { id: 'item-4', name: 'Roasted pumpkin bowl', description: 'Maple pumpkin, lentils, feta, leaves and pepita crunch.', price: 18.5, categoryId: 'cat-bowls', restaurantId: demoRestaurant.id, imageUrl: 'https://images.unsplash.com/photo-1547592180-85f173990554?auto=format&fit=crop&w=900&q=80', available: false },
  { id: 'item-5', name: 'Cloud cold brew', description: 'Slow steeped coffee with a silky oat cloud.', price: 7.5, categoryId: 'cat-drinks', restaurantId: demoRestaurant.id, imageUrl: 'https://images.unsplash.com/photo-1461023058943-07fcbe16d735?auto=format&fit=crop&w=900&q=80', available: true },
  { id: 'item-6', name: 'Pistachio morning bun', description: 'Laminated pastry with pistachio praline and citrus sugar.', price: 8, categoryId: 'cat-bakes', restaurantId: demoRestaurant.id, imageUrl: 'https://images.unsplash.com/photo-1555507036-ab1f4038808a?auto=format&fit=crop&w=900&q=80', available: false },
];

const demoAddons: AdminAddon[] = [
  { id: 'addon-1', name: 'Extra avocado', description: 'Fresh sliced avocado', price: 4, maxQuantity: 2, categoryId: 'cat-breakfast', restaurantId: demoRestaurant.id, available: true },
  { id: 'addon-2', name: 'Smoked salmon', description: 'House cured salmon', price: 6, maxQuantity: 2, categoryId: 'cat-breakfast', restaurantId: demoRestaurant.id, available: true },
  { id: 'addon-3', name: 'Soft egg', description: 'Jammy free range egg', price: 3, maxQuantity: 3, categoryId: 'cat-bowls', restaurantId: demoRestaurant.id, available: true },
  { id: 'addon-4', name: 'Extra greens', description: 'Seasonal dressed leaves', price: 2.5, maxQuantity: 2, categoryId: 'cat-bowls', restaurantId: demoRestaurant.id, available: false },
];

function demoWorkspace(): AdminWorkspace {
  return { restaurant: demoRestaurant, categories: demoCategories, menuItems: demoItems, addons: demoAddons };
}

function formatCurrency(value: number) {
  return new Intl.NumberFormat('en-SG', { style: 'currency', currency: 'SGD' }).format(value);
}

function initials(name: string) {
  return name.split(/\s+/).slice(0, 2).map((part) => part[0]).join('').toUpperCase();
}

function validateName(name: string) {
  if (!NAME_PATTERN.test(name)) throw new Error('Use 3–100 letters, numbers, spaces, apostrophes or hyphens for names.');
}

const Field = ({ label, htmlFor, required, hint, children }: { label: string; htmlFor: string; required?: boolean; hint?: string; children: React.ReactNode }) => (
  <div className="space-y-2">
    <div className="flex items-center justify-between gap-3">
      <Label htmlFor={htmlFor} className="text-[13px] font-semibold text-[#374151]">{label}{required && <span className="ml-1 text-[#F44336]">*</span>}</Label>
      {hint && <span className="text-[11px] text-[#9CA3AF]">{hint}</span>}
    </div>
    {children}
  </div>
);

export default function Admin() {
  const [restaurants, setRestaurants] = useState<AdminRestaurant[]>([]);
  const [selectedId, setSelectedId] = useState('');
  const [workspace, setWorkspace] = useState<AdminWorkspace | null>(null);
  const [loading, setLoading] = useState(true);
  const [refreshing, setRefreshing] = useState(false);
  const [isDemo, setIsDemo] = useState(false);
  const [editor, setEditor] = useState<EditorState>(null);
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [addonCategory, setAddonCategory] = useState('');
  const [availability, setAvailability] = useState('all');
  const [query, setQuery] = useState('');
  const [view, setView] = useState<'grid' | 'list'>('grid');

  const loadWorkspace = useCallback(async (restaurantId: string, quiet = false) => {
    if (!restaurantId) return;
    if (!quiet) setLoading(true); else setRefreshing(true);
    try {
      const data = await adminApi.getWorkspace(restaurantId);
      setWorkspace(data);
      setAddonCategory((current) => data.categories.some((category) => category.id === current) ? current : data.categories[0]?.id ?? '');
      setIsDemo(false);
    } catch (error) {
      if (!workspace) {
        const preview = demoWorkspace();
        setWorkspace(preview);
        setAddonCategory(preview.categories[0]?.id ?? '');
        setRestaurants([preview.restaurant]);
        setSelectedId(preview.restaurant.id);
        setIsDemo(true);
      } else {
        toast.error(error instanceof Error ? error.message : 'Could not refresh workspace');
      }
    } finally {
      setLoading(false);
      setRefreshing(false);
    }
  }, [workspace]);

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
      } catch {
        if (!active) return;
        const preview = demoWorkspace();
        setRestaurants([preview.restaurant]);
        setSelectedId(preview.restaurant.id);
        setWorkspace(preview);
        setAddonCategory(preview.categories[0]?.id ?? '');
        setIsDemo(true);
        setLoading(false);
      }
    })();
    return () => { active = false; };
    // Initial bootstrap only; restaurant changes are handled explicitly.
    // eslint-disable-next-line react-hooks/exhaustive-deps
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
      if (selectedId && !isDemo) await loadWorkspace(selectedId, true);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Something went wrong');
      throw error;
    }
  };

  const toggleAvailability = async (item: AdminMenuItem, available: boolean) => {
    setWorkspace((current) => current ? { ...current, menuItems: current.menuItems.map((entry) => entry.id === item.id ? { ...entry, available } : entry) } : current);
    if (isDemo) return toast.success(available ? 'Item is now available' : 'Item marked unavailable');
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
    if (id === demoRestaurant.id) {
      setWorkspace(demoWorkspace());
      setAddonCategory(demoCategories[0]?.id ?? '');
      setIsDemo(true);
    } else {
      await loadWorkspace(id);
    }
  };

  const deleteItem = async (item: AdminMenuItem) => {
    if (!window.confirm(`Delete ${item.name}? This cannot be undone.`)) return;
    if (isDemo) {
      setWorkspace((current) => current ? { ...current, menuItems: current.menuItems.filter((entry) => entry.id !== item.id) } : current);
      return toast.success('Menu item deleted');
    }
    await mutateWorkspace(() => adminApi.deleteMenuItem(item.id), 'Menu item deleted');
  };

  const toggleAddonAvailability = async (addon: AdminAddon, available: boolean) => {
    setWorkspace((current) => current ? { ...current, addons: current.addons.map((entry) => entry.id === addon.id ? { ...entry, available } : entry) } : current);
    if (isDemo) return toast.success(available ? 'Add-on is now available' : 'Add-on marked unavailable');
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
    if (isDemo) {
      setWorkspace((current) => current ? { ...current, addons: current.addons.filter((entry) => entry.id !== addon.id) } : current);
      return toast.success('Add-on deleted');
    }
    await mutateWorkspace(() => adminApi.deleteAddon(addon.id), 'Add-on deleted');
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
          <button className="flex w-full items-center gap-2.5 text-left" onClick={() => setEditor({ kind: 'restaurant', value: workspace?.restaurant })}>
            <div className="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg bg-[#FFF1EB] text-xs font-bold text-[#FF4500]">{workspace ? initials(workspace.restaurant.name) : 'FF'}</div>
            <div className="min-w-0 flex-1">
              <div className="truncate text-[12px] font-semibold text-white">{workspace?.restaurant.name ?? 'Choose restaurant'}</div>
              <div className="mt-0.5 flex items-center gap-1 text-[10px] text-[#D1D5DB]"><span className="h-1.5 w-1.5 rounded-full bg-[#FFB72B]" /> {workspace?.restaurant.enabled ? 'Open for business' : 'Paused'}</div>
            </div>
            <ChevronDown size={15} className="text-[#9CA3AF]" />
          </button>
        </div>

        <nav className="space-y-1">
          <div className="admin-nav-label px-3 pb-2 text-[10px] font-bold uppercase tracking-[.16em] text-[#9CA3AF]">Workspace</div>
          {navItems.map(({ icon: Icon, label, active, soon }) => (
            <button key={label} className={`admin-nav-item flex w-full items-center gap-3 rounded-xl px-3 py-2.5 text-left text-[13px] font-semibold ${active ? 'active' : ''}`}>
              <Icon size={18} strokeWidth={active ? 2.3 : 1.8} />
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
                {workspace && <span className={`rounded-full px-2 py-0.5 text-[9px] font-bold uppercase tracking-[.1em] ${workspace.restaurant.enabled ? 'bg-[#E8F5E9] text-[#2E7D32]' : 'bg-[#FFEBEE] text-[#C62828]'}`}>{workspace.restaurant.enabled ? 'Live' : 'Paused'}</span>}
              </div>
              <div className="mt-0.5 flex items-center gap-1 text-[11px] text-[#6B7280]"><MapPin size={11} /> <span className="truncate">{workspace?.restaurant.address ?? 'Restaurant workspace'}</span></div>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <button className="relative flex h-9 w-9 items-center justify-center rounded-full border border-[#E5E7EB] bg-white text-[#4B5563]"><Bell size={17} /><span className="absolute right-1.5 top-1.5 h-1.5 w-1.5 rounded-full bg-[#FF4500]" /></button>
            <Button className="admin-primary h-9 gap-2 rounded-lg px-3.5 text-xs font-semibold" onClick={() => setEditor({ kind: 'item' })} disabled={!workspace?.categories.length}>
              <Plus size={16} /> <span className="hidden sm:inline">New menu item</span><span className="sm:hidden">New</span>
            </Button>
          </div>
        </header>

        <div className="admin-content">
          {isDemo && (
            <div className="mb-5 flex flex-wrap items-center justify-between gap-3 rounded-xl border border-[#FFD78A] bg-[#FFF8E7] px-4 py-3 text-xs text-[#6B4E16]">
              <div className="flex items-center gap-2"><Sparkles size={15} className="text-[#D97706]" /><span><strong>Preview workspace.</strong> Start the local API services to load and save live restaurant data.</span></div>
              <Button variant="ghost" size="sm" className="h-7 gap-1.5 text-[#92400E] hover:bg-[#FFF0CC]" onClick={() => window.location.reload()}><RefreshCw size={13} /> Reconnect</Button>
            </div>
          )}

          <section className="mb-7 flex flex-col justify-between gap-4 md:flex-row md:items-end">
            <div>
              <div className="mb-1.5 flex items-center gap-2 text-[11px] font-bold uppercase tracking-[.14em] text-[#FF4500]"><BookOpen size={13} /> Menu workspace</div>
              <h1 className="text-[28px] font-bold tracking-[-.035em] text-[#333333] sm:text-[32px]">Build a menu people remember.</h1>
              <p className="mt-1 max-w-2xl text-[13px] text-[#6B7280]">Organise categories, curate every dish, and keep availability accurate from one calm workspace.</p>
            </div>
            <div className="flex gap-2">
              <Button variant="outline" className="h-9 gap-2 rounded-lg border-[#E5E7EB] bg-white text-xs hover:border-[#FF8C42] hover:bg-[#FFF7F3] hover:text-[#FF4500]" onClick={() => setEditor({ kind: 'category' })} disabled={!workspace}><Plus size={15} /> Add category</Button>
              <Button variant="outline" className="h-9 gap-2 rounded-lg border-[#E5E7EB] bg-white text-xs hover:border-[#FF8C42] hover:bg-[#FFF7F3] hover:text-[#FF4500]" onClick={() => setEditor({ kind: 'restaurant' })}><Building2 size={15} /> Add restaurant</Button>
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

              <section className="admin-panel mt-6 overflow-hidden rounded-2xl">
                <div className="flex flex-col border-b border-[#E5E7EB] lg:flex-row">
                  <CategoryRail
                    categories={workspace.categories}
                    counts={categoryCounts}
                    total={workspace.menuItems.length}
                    selected={selectedCategory}
                    onSelect={setSelectedCategory}
                    onAdd={() => setEditor({ kind: 'category' })}
                    onEdit={(category) => setEditor({ kind: 'category', value: category })}
                  />
                  <div className="min-w-0 flex-1">
                    <div className="flex flex-col gap-3 border-b border-[#E5E7EB] p-4 sm:flex-row sm:items-center sm:justify-between sm:px-5">
                      <div>
                        <h2 className="text-[17px] font-bold tracking-[-.02em]">{selectedCategoryName}</h2>
                        <p className="mt-0.5 text-[11px] text-[#9CA3AF]">{filteredItems.length} {filteredItems.length === 1 ? 'item' : 'items'} shown</p>
                      </div>
                      <div className="flex flex-wrap items-center gap-2">
                        <div className="relative min-w-[180px] flex-1 sm:w-[220px] sm:flex-none">
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
                        <Button variant="ghost" size="icon" className="h-9 w-9" onClick={() => loadWorkspace(selectedId, true)} disabled={refreshing || isDemo} aria-label="Refresh workspace"><RefreshCw size={15} className={refreshing ? 'animate-spin' : ''} /></Button>
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
                  <span className="flex items-center gap-1.5"><Clock3 size={12} /> Last synced {isDemo ? 'in preview mode' : 'just now'}</span>
                </div>
              </section>
            </>
          )}
        </div>
      </main>

      <EditorDialog
        editor={editor}
        workspace={workspace}
        isDemo={isDemo}
        onClose={() => setEditor(null)}
        onSave={async (kind, input, existingId) => {
          if (kind === 'restaurant') {
            const restaurantInput = input as RestaurantInput;
            validateName(restaurantInput.name);
            if (isDemo) {
              const created = { ...demoRestaurant, ...restaurantInput, id: existingId ?? `demo-${Date.now()}` };
              setRestaurants((current) => existingId ? current.map((entry) => entry.id === existingId ? created : entry) : [...current, created]);
              setWorkspace((current) => existingId && current ? { ...current, restaurant: created } : { restaurant: created, categories: [], menuItems: [], addons: [] });
              setSelectedId(created.id); setEditor(null); toast.success(existingId ? 'Restaurant updated' : 'Restaurant created'); return;
            }
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
            if (isDemo) {
              setWorkspace((current) => current ? { ...current, categories: existingId ? current.categories.map((entry) => entry.id === existingId ? { ...entry, ...categoryInput } : entry) : [...current.categories, { ...categoryInput, id: `demo-category-${Date.now()}`, enabled: true }] } : current);
              setEditor(null); toast.success(existingId ? 'Category updated' : 'Category created'); return;
            }
            await mutateWorkspace(() => existingId ? adminApi.updateCategory(existingId, categoryInput) : adminApi.createCategory(categoryInput), existingId ? 'Category updated' : 'Category created');
          }
          if (kind === 'item' && workspace) {
            const itemInput = input as MenuItemInput;
            validateName(itemInput.name);
            if (isDemo) {
              setWorkspace((current) => current ? { ...current, menuItems: existingId ? current.menuItems.map((entry) => entry.id === existingId ? { ...entry, ...itemInput } : entry) : [...current.menuItems, { ...itemInput, id: `demo-item-${Date.now()}`, available: true }] } : current);
              setEditor(null); toast.success(existingId ? 'Menu item updated' : 'Menu item created'); return;
            }
            await mutateWorkspace(() => existingId ? adminApi.updateMenuItem(existingId, itemInput) : adminApi.createMenuItem(itemInput), existingId ? 'Menu item updated' : 'Menu item created');
          }
          if (kind === 'addon' && workspace) {
            const addonInput = input as AddonInput;
            validateName(addonInput.name);
            if (isDemo) {
              setWorkspace((current) => current ? { ...current, addons: existingId ? current.addons.map((entry) => entry.id === existingId ? { ...entry, ...addonInput } : entry) : [...current.addons, { ...addonInput, id: `demo-addon-${Date.now()}`, available: true }] } : current);
              setAddonCategory(addonInput.categoryId);
              setEditor(null); toast.success(existingId ? 'Add-on updated' : 'Add-on created'); return;
            }
            await mutateWorkspace(() => existingId ? adminApi.updateAddon(existingId, { name: addonInput.name, description: addonInput.description, price: addonInput.price, maxQuantity: addonInput.maxQuantity }) : adminApi.createAddon(addonInput), existingId ? 'Add-on updated' : 'Add-on created');
            setAddonCategory(addonInput.categoryId);
          }
        }}
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

function CategoryRail({ categories, counts, total, selected, onSelect, onAdd, onEdit }: { categories: AdminCategory[]; counts: Map<string, number>; total: number; selected: string; onSelect: (id: string) => void; onAdd: () => void; onEdit: (category: AdminCategory) => void }) {
  return (
    <aside className="w-full border-b border-[#E5E7EB] bg-[#FAFAFA] lg:w-[230px] lg:shrink-0 lg:border-b-0 lg:border-r">
      <div className="flex items-center justify-between px-4 pb-2 pt-4">
        <span className="text-[10px] font-bold uppercase tracking-[.13em] text-[#9CA3AF]">Categories</span>
        <button onClick={onAdd} className="flex h-6 w-6 items-center justify-center rounded-md text-[#FF4500] hover:bg-[#FFF1EB]" aria-label="Add category"><Plus size={14} /></button>
      </div>
      <div className="flex gap-1 overflow-x-auto px-2 pb-3 lg:block lg:space-y-0.5 lg:overflow-visible lg:pb-5">
        <button onClick={() => onSelect('all')} className={`admin-category-button flex min-w-fit items-center gap-2 rounded-lg px-3 py-2 text-left text-[12px] font-semibold lg:w-full ${selected === 'all' ? 'active' : ''}`}>
          <Grid2X2 size={14} /><span className="flex-1">All items</span><span className="rounded bg-black/[.04] px-1.5 py-0.5 text-[9px]">{total}</span>
        </button>
        {categories.map((category) => (
          <div key={category.id} className={`group flex min-w-fit items-center rounded-lg lg:w-full ${selected === category.id ? 'admin-category-button active' : ''}`}>
            <button onClick={() => onSelect(category.id)} className={`admin-category-button flex min-w-0 flex-1 items-center gap-2 rounded-lg px-3 py-2 text-left text-[12px] font-medium ${selected === category.id ? 'active !bg-transparent !shadow-none' : ''}`}>
              <span className={`h-2 w-2 shrink-0 rounded-full ${category.enabled ? 'bg-[#FFB72B]' : 'bg-[#9CA3AF]'}`} />
              <span className="max-w-[120px] truncate">{category.name}</span>
              <span className="ml-auto rounded bg-black/[.04] px-1.5 py-0.5 text-[9px]">{counts.get(category.id) ?? 0}</span>
            </button>
            <button onClick={() => onEdit(category)} className="mr-1 hidden p-1 text-[#9CA3AF] hover:text-[#FF4500] lg:group-hover:block" aria-label={`Edit ${category.name}`}><Pencil size={12} /></button>
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

function EditorDialog({ editor, workspace, isDemo, onClose, onSave }: { editor: EditorState; workspace: AdminWorkspace | null; isDemo: boolean; onClose: () => void; onSave: (kind: 'restaurant' | 'category' | 'item' | 'addon', input: RestaurantInput | CategoryInput | MenuItemInput | AddonInput, existingId?: string) => Promise<void> }) {
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState('');
  const kind = editor?.kind;
  const existing = editor?.value;
  const addonCategoryId = editor?.kind === 'addon' ? (editor.value as AdminAddon | undefined)?.categoryId ?? editor.categoryId ?? workspace?.categories[0]?.id ?? '' : '';

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault(); setSaving(true); setError('');
    const data = new FormData(event.currentTarget);
    try {
      if (kind === 'restaurant') await onSave(kind, { name: String(data.get('name')), description: String(data.get('description')), address: String(data.get('address')), phone: String(data.get('phone')), email: String(data.get('email')), imageUrl: String(data.get('imageUrl')) }, existing?.id);
      if (kind === 'category' && workspace) await onSave(kind, { name: String(data.get('name')), description: String(data.get('description')), restaurantId: workspace.restaurant.id }, existing?.id);
      if (kind === 'item' && workspace) await onSave(kind, { name: String(data.get('name')), description: String(data.get('description')), price: Number(data.get('price')), categoryId: String(data.get('categoryId')), restaurantId: workspace.restaurant.id, imageUrl: String(data.get('imageUrl')) }, existing?.id);
      if (kind === 'addon' && workspace) await onSave(kind, { name: String(data.get('name')), description: String(data.get('description')), price: Number(data.get('price')), maxQuantity: Number(data.get('maxQuantity')), categoryId: String(data.get('categoryId')), restaurantId: workspace.restaurant.id }, existing?.id);
    } catch (submitError) { setError(submitError instanceof Error ? submitError.message : 'Could not save changes'); }
    finally { setSaving(false); }
  };

  const titles = { restaurant: existing ? 'Edit restaurant' : 'Create a restaurant', category: existing ? 'Edit category' : 'Create a category', item: existing ? 'Edit menu item' : 'Create a menu item', addon: existing ? 'Edit add-on' : 'Create an add-on' };
  const descriptions = { restaurant: 'The profile guests see across your storefront.', category: 'Group related items so your menu is easy to browse.', item: 'Add the details your team and guests need.', addon: 'Offer an optional extra across every item in a category.' };

  return (
    <Dialog open={Boolean(editor)} onOpenChange={(open) => !open && onClose()}>
      <DialogContent className="max-h-[92vh] overflow-y-auto border-[#E5E7EB] p-0 sm:max-w-[560px]">
        {kind && <form onSubmit={submit}>
          <DialogHeader className="border-b border-[#E5E7EB] px-6 py-5 text-left"><div className="mb-2 flex h-9 w-9 items-center justify-center rounded-xl bg-[#FFF1EB] text-[#FF4500]">{kind === 'restaurant' ? <Building2 size={18} /> : kind === 'category' ? <Boxes size={18} /> : kind === 'addon' ? <Puzzle size={18} /> : <UtensilsCrossed size={18} />}</div><DialogTitle className="text-xl tracking-[-.025em]">{titles[kind]}</DialogTitle><DialogDescription className="text-xs">{descriptions[kind]}</DialogDescription></DialogHeader>
          <div className="space-y-4 px-6 py-5">
            <Field label={kind === 'restaurant' ? 'Restaurant name' : kind === 'category' ? 'Category name' : kind === 'addon' ? 'Add-on name' : 'Item name'} htmlFor="name" required hint="3–100 characters"><Input id="name" name="name" defaultValue={existing?.name ?? ''} required minLength={3} maxLength={100} placeholder={kind === 'restaurant' ? 'e.g. Juniper Kitchen' : kind === 'category' ? 'e.g. Seasonal plates' : kind === 'addon' ? 'e.g. Extra avocado' : 'e.g. Garden harvest bowl'} className="admin-input" /></Field>
            <Field label="Description" htmlFor="description" hint="Recommended"><Textarea id="description" name="description" defaultValue={existing?.description ?? ''} rows={3} placeholder="Add a concise, useful description" className="admin-input resize-none" /></Field>
            {kind === 'restaurant' && <>
              <Field label="Address" htmlFor="address" required><Input id="address" name="address" defaultValue={(existing as AdminRestaurant | undefined)?.address ?? ''} required placeholder="Street, city and postal code" className="admin-input" /></Field>
              <div className="grid gap-4 sm:grid-cols-2"><Field label="Phone" htmlFor="phone" required><Input id="phone" name="phone" defaultValue={(existing as AdminRestaurant | undefined)?.phone ?? ''} required placeholder="+65 6123 4567" className="admin-input" /></Field><Field label="Email" htmlFor="email" required><Input id="email" name="email" type="email" defaultValue={(existing as AdminRestaurant | undefined)?.email ?? ''} required placeholder="hello@restaurant.com" className="admin-input" /></Field></div>
              <Field label="Cover image URL" htmlFor="imageUrl" hint="Optional"><Input id="imageUrl" name="imageUrl" type="url" defaultValue={(existing as AdminRestaurant | undefined)?.imageUrl ?? ''} placeholder="https://images.example.com/restaurant.jpg" className="admin-input" /></Field>
            </>}
            {kind === 'item' && workspace && <>
              <div className="grid gap-4 sm:grid-cols-2"><Field label="Price" htmlFor="price" required><div className="relative"><span className="absolute left-3 top-1/2 -translate-y-1/2 text-sm text-[#6B7280]">$</span><Input id="price" name="price" type="number" min="0.01" step="0.01" defaultValue={(existing as AdminMenuItem | undefined)?.price ?? ''} required placeholder="0.00" className="admin-input pl-7" /></div></Field><Field label="Category" htmlFor="categoryId" required><Select name="categoryId" defaultValue={(existing as AdminMenuItem | undefined)?.categoryId ?? workspace.categories[0]?.id}><SelectTrigger className="admin-input"><SelectValue placeholder="Choose category" /></SelectTrigger><SelectContent>{workspace.categories.map((category) => <SelectItem key={category.id} value={category.id}>{category.name}</SelectItem>)}</SelectContent></Select></Field></div>
              <Field label="Dish image URL" htmlFor="imageUrl" hint="Optional"><Input id="imageUrl" name="imageUrl" type="url" defaultValue={(existing as AdminMenuItem | undefined)?.imageUrl ?? ''} placeholder="https://images.example.com/dish.jpg" className="admin-input" /></Field>
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
            {isDemo && <div className="rounded-lg bg-[#FFF8E7] px-3 py-2.5 text-[11px] leading-relaxed text-[#6B4E16]">Preview mode changes stay in this browser session. Connect the API to save them permanently.</div>}
            {error && <div className="rounded-lg bg-red-50 px-3 py-2.5 text-xs text-red-700">{error}</div>}
          </div>
          <div className="flex items-center justify-end gap-2 border-t border-[#E5E7EB] bg-[#FAFAFA] px-6 py-4"><Button type="button" variant="outline" onClick={onClose} className="h-9 text-xs">Cancel</Button><Button type="submit" disabled={saving} className="admin-primary h-9 min-w-[112px] gap-2 text-xs">{saving && <Loader2 size={14} className="animate-spin" />}{existing ? 'Save changes' : kind === 'addon' ? 'Create add-on' : `Create ${kind}`}</Button></div>
        </form>}
      </DialogContent>
    </Dialog>
  );
}
