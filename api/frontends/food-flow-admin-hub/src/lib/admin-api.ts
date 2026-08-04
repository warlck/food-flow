const SALES_API_BASE_URL = import.meta.env.VITE_SALES_API_URL || '';
const AUTH_API_BASE_URL = import.meta.env.VITE_AUTH_API_URL || '';
const AUTH_KID = import.meta.env.VITE_AUTH_KID || '54bb2165-71e1-41a6-af3e-7da4a0e1e2c1';
const TOKEN_STORAGE_KEY = 'foodflow.admin-token';

export interface AdminRestaurant {
  id: string;
  name: string;
  description: string;
  address: string;
  phone: string;
  email: string;
  imageUrl: string;
  enabled: boolean;
  latitude?: number | null;
  longitude?: number | null;
  maxDeliveryDistanceKm: number;
  dateCreated: string;
  dateUpdated: string;
}

export interface AdminCategory {
  id: string;
  name: string;
  description: string;
  restaurantId: string;
  enabled: boolean;
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
  dateCreated?: string;
  dateUpdated?: string;
}

export interface AdminAddon {
  id: string;
  categoryId: string;
  restaurantId: string;
  name: string;
  description: string;
  price: number;
  available: boolean;
  maxQuantity: number;
  dateCreated?: string;
  dateUpdated?: string;
}

export interface AdminWorkspace {
  restaurant: AdminRestaurant;
  categories: AdminCategory[];
  menuItems: AdminMenuItem[];
  addons: AdminAddon[];
}

export type RestaurantInput = Pick<AdminRestaurant, 'name' | 'description' | 'address' | 'phone' | 'email' | 'imageUrl' | 'latitude' | 'longitude' | 'maxDeliveryDistanceKm'>;
export type CategoryInput = Pick<AdminCategory, 'name' | 'description' | 'restaurantId'>;
export type MenuItemInput = Pick<AdminMenuItem, 'name' | 'description' | 'price' | 'categoryId' | 'restaurantId' | 'imageUrl'>;
export type AddonInput = Pick<AdminAddon, 'name' | 'description' | 'price' | 'categoryId' | 'restaurantId' | 'maxQuantity'>;

interface ApiPage<T> {
  items: T[];
  total: number;
  page: number;
  rowsPerPage: number;
}

interface ApiDetailsCategory {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  mentuItems?: Array<Omit<AdminMenuItem, 'categoryId' | 'restaurantId'>>;
  menuItems?: Array<Omit<AdminMenuItem, 'categoryId' | 'restaurantId'>>;
}

interface ApiRestaurantDetails extends AdminRestaurant {
  categories: ApiDetailsCategory[];
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

class AdminApi {
  private token: string | null = import.meta.env.VITE_ADMIN_TOKEN || null;

  private async getToken() {
    if (this.token) return this.token;

    const storedToken = window.localStorage.getItem(TOKEN_STORAGE_KEY);
    if (storedToken) {
      this.token = storedToken;
      return storedToken;
    }

    // The current auth service exposes a development token endpoint keyed by its
    // public KID. This keeps local admin development aligned with the Go service.
    const response = await fetch(`${AUTH_API_BASE_URL}/v1/auth/token/${AUTH_KID}`, {
      headers: { Authorization: `Basic ${window.btoa('foodflow-admin:local')}` },
    });
    if (!response.ok) throw new Error('Admin authentication is unavailable. Start the auth service and try again.');

    const payload = (await response.json()) as { token: string };
    this.token = payload.token;
    window.localStorage.setItem(TOKEN_STORAGE_KEY, payload.token);
    return payload.token;
  }

  private async request<T>(path: string, init: RequestInit = {}, authenticated = true): Promise<T> {
    const headers = new Headers(init.headers);
    headers.set('Content-Type', 'application/json');
    if (authenticated) headers.set('Authorization', `Bearer ${await this.getToken()}`);

    const response = await fetch(`${SALES_API_BASE_URL}${path}`, { ...init, headers });
    if (response.status === 401 && authenticated) {
      this.token = null;
      window.localStorage.removeItem(TOKEN_STORAGE_KEY);
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

  listRestaurants() {
    return this.request<ApiPage<AdminRestaurant>>('/v1/restaurants?page=1&rows=100&orderBy=name,ASC');
  }

  async getWorkspace(restaurantId: string): Promise<AdminWorkspace> {
    const [details, addonPage] = await Promise.all([
      this.request<ApiRestaurantDetails>(`/v1/restaurants/${restaurantId}/details`, {}, false),
      this.listAddons(restaurantId),
    ]);
    const categories = details.categories.map((category) => ({
      id: category.id,
      name: category.name,
      description: category.description,
      restaurantId,
      enabled: category.enabled,
    }));
    const menuItems = details.categories.flatMap((category) =>
      (category.mentuItems ?? category.menuItems ?? []).map((item) => ({
        ...item,
        categoryId: category.id,
        restaurantId,
      })),
    );
    const { categories: _categories, ...restaurant } = details;
    return { restaurant, categories, menuItems, addons: addonPage.items };
  }

  createRestaurant(input: RestaurantInput) {
    return this.request<AdminRestaurant>('/v1/restaurants', { method: 'POST', body: JSON.stringify(input) });
  }

  updateRestaurant(id: string, input: Partial<RestaurantInput> & { enabled?: boolean }) {
    return this.request<AdminRestaurant>(`/v1/restaurants/${id}`, { method: 'PUT', body: JSON.stringify(input) });
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

  createMenuItem(input: MenuItemInput) {
    return this.request<AdminMenuItem>('/v1/menuitems', { method: 'POST', body: JSON.stringify(input) });
  }

  updateMenuItem(id: string, input: Partial<Omit<MenuItemInput, 'restaurantId'>> & { available?: boolean }) {
    return this.request<AdminMenuItem>(`/v1/menuitems/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deleteMenuItem(id: string) {
    return this.request<void>(`/v1/menuitems/${id}`, { method: 'DELETE' });
  }

  listAddons(restaurantId: string) {
    const params = new URLSearchParams({
      page: '1',
      rows: '100',
      orderBy: 'name,ASC',
      restaurant_id: restaurantId,
    });
    return this.request<ApiPage<AdminAddon>>(`/v1/addons?${params}`);
  }

  createAddon(input: AddonInput) {
    return this.request<AdminAddon>('/v1/addons', { method: 'POST', body: JSON.stringify(input) });
  }

  updateAddon(id: string, input: Partial<Omit<AddonInput, 'categoryId' | 'restaurantId'>> & { available?: boolean }) {
    return this.request<AdminAddon>(`/v1/addons/${id}`, { method: 'PUT', body: JSON.stringify(input) });
  }

  deleteAddon(id: string) {
    return this.request<void>(`/v1/addons/${id}`, { method: 'DELETE' });
  }
}

export const adminApi = new AdminApi();
