// Prefer same-origin in production (nginx proxies /v1 -> sales-service in k8s).
// In local dev (vite), default to talking to the backend directly.
const API_BASE_URL =
  import.meta.env.VITE_API_URL || (import.meta.env.DEV ? 'http://localhost:3000' : '');

// Storefront endpoints are public; no credentials are ever attached.
const jsonHeaders = (): HeadersInit => ({
  'Content-Type': 'application/json',
});

// =============================================================================
// Types

export interface CreateOrderRequest {
  restaurantId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  orderType: 'delivery' | 'pickup';
  paymentMethod: 'creditCard' | 'cash';
  promoCode?: string;
  specialInstructions?: string;
  items: OrderItemRequest[];
  deliveryAddress?: DeliveryAddressRequest;
}

export interface OrderItemModifierRequest {
  modifierGroupId: string;
  modifierOptionId: string;
}

export interface OrderItemAddonRequest {
  addonId: string;
  quantity: number;
}

export interface OrderItemRequest {
  menuItemId: string;
  quantity: number;
  specialInstructions?: string;
  modifiers?: OrderItemModifierRequest[];
  addons?: OrderItemAddonRequest[];
}

export interface DeliveryAddressRequest {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  deliveryInstructions?: string;
  latitude: number;
  longitude: number;
}

export interface OrderItemModifier {
  id: string;
  modifierGroupId: string;
  modifierGroupName: string;
  modifierOptionId: string;
  modifierOptionName: string;
  priceDelta: number;
}

export interface OrderItemAddon {
  id: string;
  addonId: string;
  addonName: string;
  addonPrice: number;
  quantity: number;
}

export interface OrderItem {
  id: string;
  menuItemId: string;
  menuItemName: string;
  menuItemPrice: number;
  quantity: number;
  specialInstructions?: string;
  modifiers?: OrderItemModifier[];
  addons?: OrderItemAddon[];
}

export interface DeliveryAddress {
  id: string;
  street: string;
  city: string;
  state: string;
  postalCode: string;
  deliveryInstructions?: string;
  latitude?: number;
  longitude?: number;
}

export interface Order {
  id: string;
  restaurantId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  orderType: 'delivery' | 'pickup';
  orderStatus: string;
  paymentStatus: string;
  paymentMethod: string;
  promoCode?: string;
  subtotal: number;
  discount?: number;
  deliveryFee: number;
  tax: number;
  total: number;
  specialInstructions?: string;
  stripePaymentIntentId?: string;
  items: OrderItem[];
  deliveryAddress?: DeliveryAddress;
  dateCreated: string;
  dateUpdated: string;
}

export interface DeliveryQuote {
  distanceKm: number;
  deliveryFee: number;
  maxDeliveryDistanceKm: number;
  withinLimit: boolean;
}

export interface PaymentIntentResponse {
  clientSecret: string;
  orderId: string;
  amount: number;
  currency: string;
}

export interface ValidatePromoRequest {
  promoCode: string;
  restaurantId?: string;
  subtotal?: number;
}

export interface ValidatePromoResponse {
  valid: boolean;
  reason?: string;
  code?: string;
  discountType?: string;
  discountValue?: number;
  discountAmount: number;
  finalSubtotal: number;
}

// =============================================================================
// API Functions

export const orderService = {
  // Create a new order
  createOrder: async (request: CreateOrderRequest): Promise<Order> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders`, {
      method: 'POST',
      headers: jsonHeaders(),
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: 'Failed to create order' }));

      // Backend uses {code, message}; older code paths may use {error}.
      throw new Error(error.message || error.error || 'Failed to create order');
    }

    return response.json();
  },

  // Get order by ID
  getOrder: async (orderId: string): Promise<Order> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders/${orderId}`, {
      method: 'GET',
      headers: jsonHeaders(),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: 'Failed to fetch order' }));
      throw new Error(error.message || error.error || 'Failed to fetch order');
    }

    return response.json();
  },

  // Create payment intent for an order
  createPaymentIntent: async (orderId: string): Promise<PaymentIntentResponse> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders/${orderId}/payment/intent`, {
      method: 'POST',
      headers: jsonHeaders(),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: 'Failed to create payment intent' }));
      throw new Error(error.message || error.error || 'Failed to create payment intent');
    }

    return response.json();
  },

  // Get a delivery fee quote for a destination
  getDeliveryQuote: async (
    restLat: number,
    restLng: number,
    lat: number,
    lng: number,
    maxDeliveryDistanceKm?: number,
  ): Promise<DeliveryQuote> => {
    const params = new URLSearchParams({
      restLat: restLat.toString(),
      restLng: restLng.toString(),
      lat: lat.toString(),
      lng: lng.toString(),
    });

    if (maxDeliveryDistanceKm !== undefined) {
      params.append('maxDeliveryDistanceKm', maxDeliveryDistanceKm.toString());
    }

    const response = await fetch(`${API_BASE_URL}/v1/orders/delivery-quote?${params}`, {
      method: 'GET',
      headers: jsonHeaders(),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: 'Failed to fetch delivery quote' }));
      throw new Error(error.message || error.error || 'Failed to fetch delivery quote');
    }

    return response.json();
  },

  // Confirm payment for an order
  confirmPayment: async (orderId: string): Promise<Order> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders/${orderId}/payment/confirm`, {
      method: 'POST',
      headers: jsonHeaders(),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: 'Failed to confirm payment' }));
      throw new Error(error.message || error.error || 'Failed to confirm payment');
    }

    return response.json();
  },

  // Validate a promo code
  validatePromoCode: async (request: ValidatePromoRequest): Promise<ValidatePromoResponse> => {
    const response = await fetch(`${API_BASE_URL}/v1/promotions/validate`, {
      method: 'POST',
      headers: jsonHeaders(),
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response
        .json()
        .catch(() => ({ message: 'Failed to validate promo code' }));
      throw new Error(error.message || error.error || 'Failed to validate promo code');
    }

    return response.json();
  },
};
