const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

// Helper to get auth token
const getAuthHeaders = (): HeadersInit => {
  const token = localStorage.getItem('authToken');
  return {
    'Content-Type': 'application/json',
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
  };
};

// =============================================================================
// Types

export interface CreateOrderRequest {
  restaurantId: string;
  customerName: string;
  customerEmail: string;
  customerPhone: string;
  orderType: 'delivery' | 'pickup';
  paymentMethod: 'creditCard' | 'cash';
  specialInstructions?: string;
  items: OrderItemRequest[];
  deliveryAddress?: DeliveryAddressRequest;
}

export interface OrderItemRequest {
  menuItemId: string;
  quantity: number;
  specialInstructions?: string;
}

export interface DeliveryAddressRequest {
  street: string;
  city: string;
  state: string;
  postalCode: string;
  deliveryInstructions?: string;
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
  subtotal: number;
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

export interface OrderItem {
  id: string;
  menuItemId: string;
  menuItemName: string;
  menuItemPrice: number;
  quantity: number;
  specialInstructions?: string;
}

export interface DeliveryAddress {
  id: string;
  street: string;
  city: string;
  state: string;
  postalCode: string;
  deliveryInstructions?: string;
}

export interface PaymentIntentResponse {
  clientSecret: string;
  orderId: string;
  amount: number;
  currency: string;
}

// =============================================================================
// API Functions

export const orderService = {
  // Create a new order
  createOrder: async (request: CreateOrderRequest): Promise<Order> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders`, {
      method: 'POST',
      headers: getAuthHeaders(),
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Failed to create order' }));
      throw new Error(error.error || 'Failed to create order');
    }

    return response.json();
  },

  // Get order by ID
  getOrder: async (orderId: string): Promise<Order> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders/${orderId}`, {
      method: 'GET',
      headers: getAuthHeaders(),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Failed to fetch order' }));
      throw new Error(error.error || 'Failed to fetch order');
    }

    return response.json();
  },

  // Create payment intent for an order
  createPaymentIntent: async (orderId: string): Promise<PaymentIntentResponse> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders/${orderId}/payment/intent`, {
      method: 'POST',
      headers: getAuthHeaders(),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Failed to create payment intent' }));
      throw new Error(error.error || 'Failed to create payment intent');
    }

    return response.json();
  },

  // Confirm payment for an order
  confirmPayment: async (orderId: string): Promise<Order> => {
    const response = await fetch(`${API_BASE_URL}/v1/orders/${orderId}/payment/confirm`, {
      method: 'POST',
      headers: getAuthHeaders(),
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({ error: 'Failed to confirm payment' }));
      throw new Error(error.error || 'Failed to confirm payment');
    }

    return response.json();
  },
};
