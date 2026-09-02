// API Configuration
// Use empty string for relative URLs (same origin) so nginx can proxy the requests
// In development, use the full URL to localhost:3000
const API_BASE_URL = import.meta.env.VITE_API_URL || '';

// API Types based on backend response
export interface ApiModifierOption {
  id: string;
  name: string;
  description: string;
  priceDelta: number;
  available: boolean;
  rank?: number | null;
}

export interface ApiModifierGroup {
  id: string;
  name: string;
  description: string;
  minSelections: number;
  maxSelections: number;
  available: boolean;
  rank?: number | null;
  options: ApiModifierOption[];
}

export interface ApiAddon {
  id: string;
  addonId?: string;
  name: string;
  description: string;
  price: number;
  available: boolean;
  maxQuantity: number;
  rank?: number | null;
}

export interface ApiMenuItem {
  id: string;
  name: string;
  description: string;
  price: number;
  imageUrl: string;
  available: boolean;
  orderable?: boolean;
  rank?: number | null;
  modifierGroups?: ApiModifierGroup[];
  addons?: ApiAddon[];
}

export interface ApiCategory {
  id: string;
  name: string;
  description: string;
  enabled: boolean;
  rank?: number | null;
  menuItems?: ApiMenuItem[];
  mentuItems?: ApiMenuItem[];
}

export interface ApiDaySchedule {
  open: string;
  close: string;
  isClosed: boolean;
}

export type ApiOperatingHours = Record<string, ApiDaySchedule>;

export interface ApiRestaurantDetails {
  id: string;
  name: string;
  description: string;
  address: string;
  phone: string;
  email: string;
  imageUrl: string;
  logoUrl?: string;
  operatingHours?: ApiOperatingHours;
  enabled: boolean;
  latitude?: number;
  longitude?: number;
  maxDeliveryDistanceKm?: number;
  minSpend?: number;
  taxRate?: number;
  categories: ApiCategory[];
  dateCreated: string;
  dateUpdated: string;
}

// API Service
export class RestaurantApiService {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  /**
   * Fetch restaurant details including categories and menu items
   */
  async getRestaurantDetails(restaurantId: string): Promise<ApiRestaurantDetails> {
    const url = `${this.baseUrl}/v1/restaurants/${restaurantId}/details`;
    
    try {
      const response = await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
      });

      if (!response.ok) {
        throw new Error(`Failed to fetch restaurant details: ${response.status} ${response.statusText}`);
      }

      const data = await response.json();
      return data;
    } catch (error) {
      console.error('Error fetching restaurant details:', error);
      throw error;
    }
  }
}

// Export singleton instance
export const restaurantApi = new RestaurantApiService();
