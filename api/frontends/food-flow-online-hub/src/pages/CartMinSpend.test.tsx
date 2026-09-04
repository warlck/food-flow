import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent } from '@testing-library/react';
import React from 'react';
import CartDesktop from './CartDesktop';
import CartMobile from './CartMobile';
import * as CartContextModule from '@/context/CartContext';
import * as useRestaurantDetailsModule from '@/hooks/useRestaurantDetails';
import { toast } from '@/components/ui/use-toast';

const mockNavigate = vi.fn();

vi.mock('react-router-dom', () => ({
  useNavigate: () => mockNavigate,
  Link: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

vi.mock('@/components/Layout', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div>{children}</div>,
}));

vi.mock('@/components/ui/use-toast', () => ({
  toast: vi.fn(),
}));

vi.mock('@/context/CartContext', async (importOriginal) => ({
  ...(await importOriginal<typeof import('@/context/CartContext')>()),
  useCart: vi.fn(),
}));

vi.mock('@/hooks/useRestaurantDetails', () => ({
  useRestaurantDetails: vi.fn(),
}));

const sampleCartItem = {
  menuItem: {
    id: 'item-1',
    name: 'Falafel Roll',
    price: 10.0,
    description: 'Fresh roll',
    category: 'Rolls',
    available: true,
  },
  quantity: 1,
  selectedModifiers: [],
  selectedAddons: [],
  specialInstructions: '',
};

const restaurantWithMinSpend = {
  id: 'rest-1',
  name: 'Test Kitchen',
  enabled: true,
  taxRate: 0.10,
  minSpend: 25.0, // Minimum spend of $25.00
};

describe('Cart Minimum Spend - Pickup Exemption', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('CartDesktop', () => {
    it('enforces minimum spend when orderType is delivery', () => {
      vi.mocked(useRestaurantDetailsModule.useRestaurantDetails).mockReturnValue({
        data: restaurantWithMinSpend as any,
        isLoading: false,
        isError: false,
      } as any);

      vi.mocked(CartContextModule.useCart).mockReturnValue({
        items: [sampleCartItem],
        getTotalPrice: () => 10.0, // $10 < $25 minSpend
        hasItems: () => true,
        clearCart: vi.fn(),
        orderType: 'delivery',
        setOrderType: vi.fn(),
        restaurantId: 'rest-1',
        appliedPromo: null,
        applyPromoCode: vi.fn(),
        removePromoCode: vi.fn(),
      } as any);

      render(<CartDesktop />);

      // Warning banner is displayed
      expect(screen.getByText('Minimum Order Required')).toBeDefined();
      expect(screen.getByText(/Add \$15.00 more to meet the minimum of \$25.00/i)).toBeDefined();

      // Proceed button is disabled
      const checkoutBtn = screen.getByRole('button', { name: /Proceed to Checkout/i }) as HTMLButtonElement;
      expect(checkoutBtn.disabled).toBe(true);

      // Clicking does not navigate and triggers toast error
      fireEvent.click(checkoutBtn);
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('exempts minimum spend when orderType is pickup', () => {
      vi.mocked(useRestaurantDetailsModule.useRestaurantDetails).mockReturnValue({
        data: restaurantWithMinSpend as any,
        isLoading: false,
        isError: false,
      } as any);

      vi.mocked(CartContextModule.useCart).mockReturnValue({
        items: [sampleCartItem],
        getTotalPrice: () => 10.0, // $10 < $25 minSpend
        hasItems: () => true,
        clearCart: vi.fn(),
        orderType: 'pickup',
        setOrderType: vi.fn(),
        restaurantId: 'rest-1',
        appliedPromo: null,
        applyPromoCode: vi.fn(),
        removePromoCode: vi.fn(),
      } as any);

      render(<CartDesktop />);

      // Warning banner is NOT displayed
      expect(screen.queryByText('Minimum Order Required')).toBeNull();

      // Proceed button is ENABLED
      const checkoutBtn = screen.getByRole('button', { name: /Proceed to Checkout/i }) as HTMLButtonElement;
      expect(checkoutBtn.disabled).toBe(false);

      // Clicking proceeds to checkout without toast error
      fireEvent.click(checkoutBtn);
      expect(mockNavigate).toHaveBeenCalledWith('/checkout');
      expect(toast).not.toHaveBeenCalled();
    });
  });

  describe('CartMobile', () => {
    it('enforces minimum spend on mobile when orderType is delivery', () => {
      vi.mocked(useRestaurantDetailsModule.useRestaurantDetails).mockReturnValue({
        data: restaurantWithMinSpend as any,
        isLoading: false,
        isError: false,
      } as any);

      vi.mocked(CartContextModule.useCart).mockReturnValue({
        items: [sampleCartItem],
        getTotalPrice: () => 10.0, // $10 < $25 minSpend
        hasItems: () => true,
        clearCart: vi.fn(),
        orderType: 'delivery',
        setOrderType: vi.fn(),
        restaurantId: 'rest-1',
        appliedPromo: null,
        applyPromoCode: vi.fn(),
        removePromoCode: vi.fn(),
      } as any);

      render(<CartMobile />);

      expect(screen.getByText('Minimum Order Required')).toBeDefined();

      const checkoutBtn = screen.getByRole('button', { name: /Proceed to Checkout/i }) as HTMLButtonElement;
      expect(checkoutBtn.disabled).toBe(true);

      fireEvent.click(checkoutBtn);
      expect(mockNavigate).not.toHaveBeenCalled();
    });

    it('exempts minimum spend on mobile when orderType is pickup', () => {
      vi.mocked(useRestaurantDetailsModule.useRestaurantDetails).mockReturnValue({
        data: restaurantWithMinSpend as any,
        isLoading: false,
        isError: false,
      } as any);

      vi.mocked(CartContextModule.useCart).mockReturnValue({
        items: [sampleCartItem],
        getTotalPrice: () => 10.0, // $10 < $25 minSpend
        hasItems: () => true,
        clearCart: vi.fn(),
        orderType: 'pickup',
        setOrderType: vi.fn(),
        restaurantId: 'rest-1',
        appliedPromo: null,
        applyPromoCode: vi.fn(),
        removePromoCode: vi.fn(),
      } as any);

      render(<CartMobile />);

      expect(screen.queryByText('Minimum Order Required')).toBeNull();

      const checkoutBtn = screen.getByRole('button', { name: /Proceed to Checkout/i }) as HTMLButtonElement;
      expect(checkoutBtn.disabled).toBe(false);

      fireEvent.click(checkoutBtn);
      expect(mockNavigate).toHaveBeenCalledWith('/mobile-checkout');
      expect(toast).not.toHaveBeenCalled();
    });
  });
});
