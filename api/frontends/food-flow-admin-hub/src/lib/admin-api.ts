import { NotAuthenticatedError, getValidToken, logout as clearStoredToken } from '@/lib/auth';

const SALES_API_BASE_URL = import.meta.env.VITE_SALES_API_URL || '';

export interface DaySchedule {
  open: string;
  close: string;
  isClosed: boolean;
}

export type OperatingHours = Record<string, DaySchedule>;

export const DEFAULT_OPERATING_HOURS: OperatingHours = {
  monday:    { open: '10:00', close: '22:00', isClosed: false },
  tuesday:   { open: '10:00', close: '22:00', isClosed: false },
  wednesday: { open: '10:00', close: '22:00', isClosed: false },
  thursday:  { open: '10:00', close: '22:00', isClosed: false },
  friday:    { open: '10:00', close: '23:00', isClosed: false },
  saturday:  { open: '11:00', close: '23:00', isClosed: false },
  sunday:    { open: '11:00', close: '22:00', isClosed: false },
};

export interface AdminRestaurant {
  id: string;
  name: string;
  description: string;
  address: string;
  phone: string;
  email: string;
  imageUrl: string;
  logoUrl?: string;
  operatingHours?: OperatingHours;
  enabled: boolean;
  latitude?: number | null;
  longitude?: number | null;
  maxDeliveryDistanceKm: number;
  minSpend?: number;
  taxRate: number;
  dateCreated: string;
  dateUpdated: string;
}

export interface AdminCategory {
  id: string;
  name: string;
  description: string;
  restaurantId: string;
  enabled: boolean;
  rank?: number | null;
  dateCreated?: string;
  dateUpdated?: string;
}

export interface AdminMenuItem {
  id: string;
  name: string;
  description: string;
  price: number;
  categoryId: string;
  restaurantId: string;
  imageUrl: string;
  available: boolean;
  rank?: number | null;
  modifierGroups?: AdminModifierGroup[];
  addons?: AdminAddon[];
  dateCreated?: string;
  dateUpdated?: string;
}

export interface AdminModifierOption {
  id: string;
  modifierGroupId: string;
  restaurantId: string;
  name: string;
  description: string;
  priceDelta: number;
  available: boolean;
  rank?: number | null;
  dateCreated?: string;
  dateUpdated?: string;
}

export interface AdminModifierGroup {
  id: string;
  menuItemId: string;
  restaurantId: string;
  name: string;
  description: string;
  minSelections: number;
  maxSelections: number;
  available: boolean;
  rank?: number | null;
  options?: AdminModifierOption[];
  dateCreated?: string;
  dateUpdated?: string;
}

export interface AdminAddon {
  id: string;
  menuItemId: string;
  restaurantId: string;
  name: string;
  description: string;
  price: number;
  available: boolean;
  maxQuantity: number;
  rank?: number | null;
  dateCreated?: string;
  dateUpdated?: string;
}

export interface AddonInput {
  menuItemId: string;
  restaurantId: string;
  name: string;
  description?: string;
  price: number;
  available?: boolean;
  maxQuantity?: number;
  rank?: number | null;
}

export interface AdminWorkspace {
  restaurant: AdminRestaurant;
  categories: AdminCategory[];
  menuItems: AdminMenuItem[];
  addons: AdminAddon[];
}

export type OrderStatus = 'pending' | 'confirmed' | 'preparing' | 'ready' | 'out_for_delivery' | 'completed' | 'cancelled';
export type PaymentStatus = 'pending' | 'processing' | 'paid' | 'failed' | 'refunded';
export type OrderType = 'pickup' | 'delivery';

export interface AdminOrderItemModifier {
  id: string;
  modifierGroupId: string;
  modifierGroupName: string;
  modifierOptionId: string;
  modifierOptionName: string;
  priceDelta: number;
}

export interface AdminOrderItemAddon {
  id: string;
  addonId: string;
  addonName: string;
  addonPrice: number;
  quantity: number;
}

export interface AdminOrderItem {
  id: string;
  menuItemId: string;
  menuItemName: string;
  menuItemPrice: number;
  quantity: number;
  specialInstructions?: string;
  modifiers?: AdminOrderItemModifier[];
  addons?: AdminOrderItemAddon[];
  dateCreated: string;
}

export interface AdminDeliveryAddress {
  id: string;
  street: string;
  city: string;
  state: string;
  postalCode: string;
  deliveryInstructions?: string;
  latitude?: number;
  longitude?: number;
  dateCreated: string;
}

export interface AdminOrder {
  id: string;
  restaurantId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  orderType: OrderType;
  orderStatus: OrderStatus;
  paymentStatus: PaymentStatus;
  paymentMethod: string;
  subtotal: number;
  deliveryFee: number;
  tax: number;
  total: number;
  specialInstructions?: string;
  stripePaymentIntentId?: string;
  items: AdminOrderItem[];
  deliveryAddress?: AdminDeliveryAddress;
  dateCreated: string;
  dateUpdated: string;
}

export interface OrderFilters {
  orderStatus?: OrderStatus;
  paymentStatus?: PaymentStatus;
  orderType?: OrderType;
  customerEmail?: string;
  startDate?: string;
  endDate?: string;
}

export interface OrderStatusInput {
  orderStatus?: OrderStatus;
  paymentStatus?: PaymentStatus;
}

export type RestaurantInput = Pick<AdminRestaurant, 'name' | 'description' | 'address' | 'phone' | 'email' | 'imageUrl' | 'logoUrl' | 'operatingHours' | 'enabled' | 'latitude' | 'longitude' | 'maxDeliveryDistanceKm' | 'minSpend' | 'taxRate'> & {
  organizationId?: string;
};
export type CategoryInput = Pick<AdminCategory, 'name' | 'description' | 'restaurantId'> & { rank?: number | null };
export type MenuItemInput = Pick<AdminMenuItem, 'name' | 'description' | 'price' | 'categoryId' | 'restaurantId' | 'imageUrl'> & { rank?: number | null };
export type ModifierGroupInput = Pick<AdminModifierGroup, 'name' | 'description' | 'minSelections' | 'maxSelections' | 'menuItemId' | 'restaurantId'> & { available?: boolean; rank?: number | null };
export type ModifierOptionInput = Pick<AdminModifierOption, 'name' | 'description' | 'priceDelta' | 'modifierGroupId' | 'restaurantId'> & { available?: boolean; rank?: number | null };

interface ApiPage<T> {
  items: T[];
  total: number;
  page: number;
  rowsPerPage: number;
}

function readErrorMessage(payload: unknown, fallback: string) {
  if (!payload || typeof payload !== 'object') return fallback;
  const record = payload as Record<string, unknown>;
  if (typeof record.message === 'string') return record.message;
  if (record.error && typeof record.error === 'object') {
    const message = (record.error as Record<string, unknown>).message;
    if (typeof message === 'string') return message;
  }
  return fallback;
}

export interface AdminOrganization {
  id: string;
  name: string;
  dateCreated: string;
  dateUpdated: string;
}

class AdminApi {
  private unauthorizedHandler: (() => void) | null = null;

  setUnauthorizedHandler(cb: () => void) {
    this.unauthorizedHandler = cb;
  }

  private getToken(): string {
    const token = getValidToken();
    if (!token) {
      this.handleUnauthorized();
      throw new NotAuthenticatedError();
    }
    return token;
  }

  private handleUnauthorized() {
    clearStoredToken();
    this.unauthorizedHandler?.();
  }

  private async request<T>(path: string, init: RequestInit = {}, authenticated = true): Promise<T> {
    const headers = new Headers(init.headers);
    headers.set('Content-Type', 'application/json');
    if (authenticated) headers.set('Authorization', `Bearer ${this.getToken()}`);

    const response = await fetch(`${SALES_API_BASE_URL}${path}`, { ...init, headers });
    if (response.status === 401 && authenticated) {
      this.handleUnauthorized();
    }
    if (!response.ok) {
      let payload: unknown;
      try {
        payload = await response.json();
      } catch {
        payload = null;
      }
      throw new Error(readErrorMessage(payload, `${response.status} ${response.statusText}`));
    }
    if (response.status === 204) return undefined as T;
    return response.json() as Promise<T>;
  }

  listMyOrganizations() {
    return this.request<AdminOrganization[]>('/v1/organizations/me');
  }

  listRestaurants() {
    return this.request<ApiPage<AdminRestaurant>>('/v1/restaurants?page=1&rows=100&orderBy=name,ASC');
  }

  async getWorkspace(restaurantId: string): Promise<AdminWorkspace> {
    const [restaurant, categoryPage, menuItemPage, addonPage] = await Promise.all([
      this.request<AdminRestaurant>(`/v1/restaurants/${restaurantId}`),
      this.listCategories(restaurantId),
      this.listMenuItems(restaurantId),
      this.listAddons(restaurantId),
    ]);

    return {
      restaurant,
      categories: categoryPage.items,
      menuItems: menuItemPage.items,
      addons: addonPage.items,
    };
  }

  createRestaurant(input: RestaurantInput) {
    return this.request<AdminRestaurant>('/v1/restaurants', { method: 'POST', body: JSON.stringify(input) });
  }

  updateRestaurant(id: string, input: Partial<RestaurantInput> & { enabled?: boolean }) {
    return this.request<AdminRestaurant>(`/v1/restaurants/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  listCategories(restaurantId: string) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'rank,ASC',
      restaurant_id: restaurantId,
    });
    return this.request<ApiPage<AdminCategory>>(`/v1/categories?${params}`);
  }

  createCategory(input: CategoryInput) {
    return this.request<AdminCategory>('/v1/categories', { method: 'POST', body: JSON.stringify(input) });
  }

  updateCategory(id: string, input: Partial<Omit<CategoryInput, 'restaurantId'>> & { enabled?: boolean }) {
    return this.request<AdminCategory>(`/v1/categories/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deleteCategory(id: string) {
    return this.request<void>(`/v1/categories/${id}`, { method: 'DELETE' });
  }

  reorderCategories(input: { restaurantId: string; orderedIds: string[] }) {
    return this.request<AdminCategory[]>('/v1/categories/order', { method: 'PUT', body: JSON.stringify(input) });
  }

  listMenuItems(restaurantId: string, categoryId?: string) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'rank,ASC',
      restaurant_id: restaurantId,
    });
    if (categoryId) params.set('category_id', categoryId);
    return this.request<ApiPage<AdminMenuItem>>(`/v1/menuitems?${params}`);
  }

  createMenuItem(input: MenuItemInput) {
    return this.request<AdminMenuItem>('/v1/menuitems', { method: 'POST', body: JSON.stringify(input) });
  }

  updateMenuItem(id: string, input: Partial<Omit<MenuItemInput, 'restaurantId'>> & { available?: boolean }) {
    return this.request<AdminMenuItem>(`/v1/menuitems/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deleteMenuItem(id: string) {
    return this.request<void>(`/v1/menuitems/${id}`, { method: 'DELETE' });
  }

  reorderMenuItems(input: { categoryId: string; orderedIds: string[] }) {
    return this.request<AdminMenuItem[]>('/v1/menuitems/order', { method: 'PUT', body: JSON.stringify(input) });
  }

  listModifierGroups(menuItemId: string) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'rank,ASC',
      menu_item_id: menuItemId,
    });
    return this.request<ApiPage<AdminModifierGroup>>(`/v1/modifiergroups?${params}`);
  }

  createModifierGroup(input: ModifierGroupInput) {
    return this.request<AdminModifierGroup>('/v1/modifiergroups', { method: 'POST', body: JSON.stringify(input) });
  }

  updateModifierGroup(id: string, input: Partial<Omit<ModifierGroupInput, 'menuItemId' | 'restaurantId'>>) {
    return this.request<AdminModifierGroup>(`/v1/modifiergroups/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deleteModifierGroup(id: string) {
    return this.request<void>(`/v1/modifiergroups/${id}`, { method: 'DELETE' });
  }

  reorderModifierGroups(input: { menuItemId: string; orderedIds: string[] }) {
    return this.request<AdminModifierGroup[]>('/v1/modifiergroups/order', { method: 'PUT', body: JSON.stringify(input) });
  }

  listModifierOptions(modifierGroupId: string) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'rank,ASC',
      modifier_group_id: modifierGroupId,
    });
    return this.request<ApiPage<AdminModifierOption>>(`/v1/modifieroptions?${params}`);
  }

  createModifierOption(input: ModifierOptionInput) {
    return this.request<AdminModifierOption>('/v1/modifieroptions', { method: 'POST', body: JSON.stringify(input) });
  }

  updateModifierOption(id: string, input: Partial<Omit<ModifierOptionInput, 'modifierGroupId' | 'restaurantId'>>) {
    return this.request<AdminModifierOption>(`/v1/modifieroptions/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deleteModifierOption(id: string) {
    return this.request<void>(`/v1/modifieroptions/${id}`, { method: 'DELETE' });
  }

  reorderModifierOptions(input: { modifierGroupId: string; orderedIds: string[] }) {
    return this.request<AdminModifierOption[]>('/v1/modifieroptions/order', { method: 'PUT', body: JSON.stringify(input) });
  }

  listAddons(restaurantId: string, menuItemId?: string) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'rank,ASC',
      restaurant_id: restaurantId,
    });
    if (menuItemId) params.set('menu_item_id', menuItemId);
    return this.request<ApiPage<AdminAddon>>(`/v1/addons?${params}`);
  }

  createAddon(input: AddonInput) {
    return this.request<AdminAddon>('/v1/addons', { method: 'POST', body: JSON.stringify(input) });
  }

  updateAddon(id: string, input: Partial<Omit<AddonInput, 'restaurantId' | 'menuItemId'>> & { available?: boolean }) {
    return this.request<AdminAddon>(`/v1/addons/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deleteAddon(id: string) {
    return this.request<void>(`/v1/addons/${id}`, { method: 'DELETE' });
  }

  reorderAddons(input: { menuItemId: string; orderedIds: string[] }) {
    return this.request<AdminAddon[]>('/v1/addons/reorder', { method: 'POST', body: JSON.stringify(input) });
  }

  listOrders(restaurantId: string, filters: OrderFilters = {}) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'date,DESC',
      restaurant_id: restaurantId,
    });
    if (filters.orderStatus) params.set('order_status', filters.orderStatus);
    if (filters.paymentStatus) params.set('payment_status', filters.paymentStatus);
    if (filters.orderType) params.set('order_type', filters.orderType);
    if (filters.customerEmail) params.set('customer_email', filters.customerEmail);
    if (filters.startDate) params.set('start_date', filters.startDate);
    if (filters.endDate) params.set('end_date', filters.endDate);
    return this.request<ApiPage<AdminOrder>>(`/v1/orders?${params}`);
  }

  updateOrderStatus(id: string, input: OrderStatusInput) {
    return this.request<AdminOrder>(`/v1/orders/${id}/status`, { method: 'PATCH', body: JSON.stringify(input) });
  }

  cancelOrder(id: string) {
    return this.request<void>(`/v1/orders/${id}/cancel`, { method: 'POST' });
  }

  getInsights(filters: InsightsFilters = {}) {
    const params = new URLSearchParams();
    if (filters.restaurantId) params.set('restaurant_id', filters.restaurantId);
    if (filters.startDate) params.set('start_date', filters.startDate);
    if (filters.endDate) params.set('end_date', filters.endDate);
    const qs = params.toString();
    return this.request<AdminInsights>(`/v1/insights${qs ? `?${qs}` : ''}`);
  }

  listPromotions(restaurantId?: string) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'code,ASC',
    });
    if (restaurantId) params.set('restaurant_id', restaurantId);
    return this.request<ApiPage<AdminPromotion>>(`/v1/promotions?${params}`);
  }

  createPromotion(input: PromotionInput) {
    return this.request<AdminPromotion>('/v1/promotions', { method: 'POST', body: JSON.stringify(input) });
  }

  updatePromotion(id: string, input: Partial<PromotionInput>) {
    return this.request<AdminPromotion>(`/v1/promotions/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deletePromotion(id: string) {
    return this.request<void>(`/v1/promotions/${id}`, { method: 'DELETE' });
  }

  requestImageUploadUrl(input: { restaurantId: string; entityType: ImageEntityType; contentType: string; sizeBytes: number }) {
    return this.request<ImageUploadGrant>('/v1/images/upload-url', { method: 'POST', body: JSON.stringify(input) });
  }

  completeImageUpload(imageId: string) {
    return this.request<AdminImage>(`/v1/images/${imageId}/complete`, { method: 'POST' });
  }

  deleteImage(imageId: string) {
    return this.request<void>(`/v1/images/${imageId}`, { method: 'DELETE' });
  }

  private async uploadToSignedUrl(grant: ImageUploadGrant, file: File) {
    const headers = new Headers(grant.headers);
    if (grant.uploadUrl.startsWith('/')) {
      headers.set('Authorization', `Bearer ${this.getToken()}`);
    }
    const response = await fetch(grant.uploadUrl, {
      method: grant.method || 'PUT',
      headers,
      body: file,
    });
    if (!response.ok) {
      throw new Error(`Upload to storage failed with ${response.status}. Please try again.`);
    }
  }

  async uploadEntityImage(input: { restaurantId: string; entityType: ImageEntityType; file: File }): Promise<AdminImage> {
    const grant = await this.requestImageUploadUrl({
      restaurantId: input.restaurantId,
      entityType: input.entityType,
      contentType: input.file.type,
      sizeBytes: input.file.size,
    });
    try {
      await this.uploadToSignedUrl(grant, input.file);
      return await this.completeImageUpload(grant.image.imageId);
    } catch (err) {
      await this.deleteImage(grant.image.imageId).catch(() => undefined);
      throw err;
    }
  }
}

export interface AdminPromotion {
  id: string;
  restaurantId?: string | null;
  code: string;
  name: string;
  description?: string;
  discountType: 'percentage' | 'fixed_amount';
  discountValue: number;
  minOrderAmount: number;
  maxDiscountAmount?: number | null;
  usageLimit?: number | null;
  usageCount: number;
  startDate?: string | null;
  endDate?: string | null;
  enabled: boolean;
  dateCreated: string;
  dateUpdated: string;
}

export type PromotionInput = {
  restaurantId?: string | null;
  code: string;
  name: string;
  description?: string;
  discountType: 'percentage' | 'fixed_amount';
  discountValue: number;
  minOrderAmount: number;
  maxDiscountAmount?: number | null;
  usageLimit?: number | null;
  enabled: boolean;
};

export interface AdminImage {
  imageId: string;
  restaurantId: string;
  entityType: ImageEntityType;
  objectPath: string;
  publicUrl: string;
  contentType: string;
  sizeBytes: number;
  status: 'pending' | 'confirmed';
  dateCreated: string;
  dateUpdated: string;
}

export type ImageEntityType = 'restaurant' | 'menu_item';

export interface ImageUploadGrant {
  image: AdminImage;
  uploadUrl: string;
  method: string;
  headers: Record<string, string>;
  expiresAt: string;
}

export const IMAGE_UPLOAD_ACCEPT = 'image/jpeg,image/png,image/webp';
export const IMAGE_UPLOAD_MAX_BYTES = 5 * 1024 * 1024;

export interface InsightsFilters {
  restaurantId?: string;
  startDate?: string;
  endDate?: string;
}

export interface AdminSalesSummary {
  grossSales: number;
  netSales: number;
  totalOrders: number;
  completedOrders: number;
  cancelledOrders: number;
  averageOrderValue: number;
  totalDiscounts: number;
  totalDeliveryFees: number;
  totalTax: number;
  totalCollected: number;
}

export interface AdminTimeSeriesPoint {
  date: string;
  grossSales: number;
  netSales: number;
  totalCollected: number;
  orderCount: number;
  averageOrder: number;
}

export interface AdminTopItemMetric {
  menuItemId: string;
  menuItemName: string;
  categoryName: string;
  quantitySold: number;
  totalRevenue: number;
}

export interface AdminTopCategoryMetric {
  categoryId: string;
  categoryName: string;
  quantitySold: number;
  totalRevenue: number;
  percentage: number;
}

export interface AdminTopAddonMetric {
  addonId: string;
  addonName: string;
  quantitySold: number;
  totalRevenue: number;
}

export interface AdminOrderTypeMetric {
  orderType: string;
  count: number;
  totalRevenue: number;
  percentage: number;
}

export interface AdminHourlyMetric {
  hour: number;
  count: number;
  totalRevenue: number;
}

export interface AdminInsights {
  summary: AdminSalesSummary;
  salesOverTime: AdminTimeSeriesPoint[];
  topItems: AdminTopItemMetric[];
  topCategories: AdminTopCategoryMetric[];
  topAddons: AdminTopAddonMetric[];
  orderTypes: AdminOrderTypeMetric[];
  peakHours: AdminHourlyMetric[];
}

export const adminApi = new AdminApi();
