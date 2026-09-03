export interface ModifierOption {
  id: string;
  modifierGroupId?: string;
  name: string;
  description?: string;
  priceDelta: number;
  available: boolean;
  rank?: number | null;
}

export interface ModifierGroup {
  id: string;
  menuItemId?: string;
  name: string;
  description?: string;
  minSelections: number;
  maxSelections: number;
  available: boolean;
  rank?: number | null;
  options: ModifierOption[];
}

export interface Addon {
  id: string;
  addonId?: string;
  name: string;
  description: string;
  price: number;
  available: boolean;
  maxQuantity: number;
  rank?: number | null;
}

export interface SelectedModifier {
  modifierGroupId: string;
  modifierGroupName: string;
  modifierOptionId: string;
  modifierOptionName: string;
  priceDelta: number;
}

export interface SelectedAddon {
  addon: Addon;
  quantity: number;
}

export interface MenuItem {
  id: string;
  name: string;
  description: string;
  price: number;
  image: string;
  category: string;
  available: boolean;
  orderable: boolean;
  preparationTime: number; // in minutes
  restaurantId: string;
  rank?: number | null;
  tags?: string[];
  modifierGroups?: ModifierGroup[];
  addons?: Addon[];
}

export interface DaySchedule {
  open: string;
  close: string;
  isClosed?: boolean;
}

export interface Restaurant {
  id: string;
  name: string;
  description: string;
  logo: string;
  coverImage: string;
  address: string;
  phone: string;
  email: string;
  enabled?: boolean;
  openingHours: {
    [key: string]: DaySchedule;
  };
  latitude?: number;
  longitude?: number;
  maxDeliveryDistanceKm?: number;
  taxRate?: number;
  deliveryFee: number;
  minimumOrder: number;
  estimatedDeliveryTime: {
    min: number;
    max: number;
  };
  estimatedPickupTime?: {
    min: number;
    max: number;
  };
  rating: number;
}

export interface CartItem {
  cartItemId: string; // Unique identifier for this cart entry
  menuItem: MenuItem;
  quantity: number;
  specialInstructions?: string;
  selectedModifiers?: SelectedModifier[];
  selectedAddons?: SelectedAddon[];
  unitPrice?: number;
}

export type OrderType = 'delivery' | 'pickup';

export interface Order {
  id: string;
  restaurantId: string;
  customerId: string;
  items: CartItem[];
  status: 'pending' | 'confirmed' | 'preparing' | 'ready' | 'out-for-delivery' | 'delivered' | 'cancelled';
  totalAmount: number;
  deliveryAddress?: Address;
  orderType: OrderType;
  createdAt: Date;
  estimatedDeliveryTime?: Date;
  estimatedPickupTime?: Date;
}

export interface Address {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  country: string;
  additionalInstructions?: string;
}

export interface Customer {
  id: string;
  name: string;
  email: string;
  phone: string;
  addresses: Address[];
  // Extended customer data for analytics
  orderCount?: number;
  totalSpent?: number;
  lastVisit?: Date;
  preferredPayment?: string;
  dietaryPreferences?: string[];
  favoriteItems?: string[];
}
