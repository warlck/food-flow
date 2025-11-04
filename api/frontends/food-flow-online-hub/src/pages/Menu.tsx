import React, { useState, useEffect, useRef } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import Layout from '@/components/Layout';
import MenuGrid from '@/components/MenuGrid';
import RestaurantInfo from '@/components/RestaurantInfo';
import { mockMenuItems, mockRestaurant, mockCategories } from '@/data/mockData';
import { useRestaurantDetails } from '@/hooks/useRestaurantDetails';
import { transformApiRestaurant, transformApiMenuItems } from '@/lib/transformers';
import { Button } from '@/components/ui/button';
import { ScrollArea } from '@/components/ui/scroll-area';
import { useCart } from '@/context/CartContext';
import CartItemComponent from '@/components/CartItemComponent';
import { ShoppingBag, Loader2, AlertCircle } from 'lucide-react';
import { useIsMobile } from '@/hooks/use-mobile';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';

const Menu: React.FC = () => {
  const { restaurantId: urlRestaurantId } = useParams<{ restaurantId: string }>();
  const [searchParams] = useSearchParams();
  const queryRestaurantId = searchParams.get('restaurant_id');
  
  // Use restaurant_id from URL path, query param, or default
  const restaurantId = urlRestaurantId || queryRestaurantId || 'a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d';
  
  // Fetch restaurant details from API
  const { data: apiData, isLoading, error } = useRestaurantDetails(restaurantId);
  
  const { items, getTotalItems, getTotalPrice, hasItems } = useCart();
  const isMobile = useIsMobile();
  const menuSectionRef = useRef<HTMLElement>(null);
  
  // Transform API data or use mock data as fallback
  const restaurant = apiData ? transformApiRestaurant(apiData) : mockRestaurant;
  const { items: menuItems, categories } = apiData 
    ? transformApiMenuItems(apiData) 
    : { items: mockMenuItems, categories: mockCategories };
  
  useEffect(() => {
    if (menuSectionRef.current) {
      menuSectionRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  }, []);

  const CartComponent = () => (
    <div className="flex flex-col h-full">
      <div className="flex items-center justify-between mb-3">
        <div>
          <h2 className="text-lg font-bold">Your Order</h2>
          <p className="text-xs text-gray-500">{getTotalItems()} items</p>
        </div>
      </div>
      
      {hasItems() ? (
        <>
          <ScrollArea className="flex-1 pr-3" style={{ height: 'calc(100vh - 240px)' }}>
            <div className="space-y-0">
              {items.map((item) => (
                <CartItemComponent key={item.menuItem.id} item={item} />
              ))}
            </div>
          </ScrollArea>
          
          <div className="border-t pt-3 mt-auto">
            <div className="flex justify-between text-base font-semibold mb-3">
              <span>Total:</span>
              <span>${getTotalPrice().toFixed(2)}</span>
            </div>
            <Button className="w-full" size="lg" asChild>
              <a href="/checkout">Proceed to Checkout</a>
            </Button>
          </div>
        </>
      ) : (
        <div className="flex flex-col items-center justify-center h-40 text-center">
          <ShoppingBag className="h-12 w-12 text-gray-300 mb-2" />
          <h3 className="text-lg font-medium">Your cart is empty</h3>
          <p className="text-sm text-gray-500 mt-1">Add items from the menu to get started</p>
        </div>
      )}
    </div>
  );

  // Loading state
  if (isLoading) {
    return (
      <Layout>
        <div className="container mx-auto px-4 py-8 flex items-center justify-center min-h-[50vh]">
          <div className="text-center">
            <Loader2 className="h-12 w-12 animate-spin mx-auto mb-4 text-food-primary" />
            <h2 className="text-xl font-semibold">Loading restaurant details...</h2>
          </div>
        </div>
      </Layout>
    );
  }

  // Error state
  if (error) {
    return (
      <Layout>
        <div className="container mx-auto px-4 py-8">
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertTitle>Error loading restaurant</AlertTitle>
            <AlertDescription>
              {error.message || 'Failed to load restaurant details. Please try again later.'}
            </AlertDescription>
          </Alert>
          <div className="mt-6">
            <Button onClick={() => window.location.reload()}>
              Try Again
            </Button>
          </div>
        </div>
      </Layout>
    );
  }

  return (
    <Layout>
      {/* Restaurant Info */}
      <RestaurantInfo restaurant={restaurant} />
      
      {/* Menu Grid with persistent side cart */}
      <section ref={menuSectionRef} className="container mx-auto px-4 py-8 relative">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-3xl font-bold">Menu</h2>
        </div>
        
        <div className="flex gap-4">
          {/* Main menu content */}
          <div className="flex-1 pr-0 lg:pr-4">
            <MenuGrid items={menuItems} categories={categories} />
          </div>
          
          {/* Persistent side cart for desktop */}
          {!isMobile && (
            <div className="hidden lg:flex lg:flex-col w-[280px] xl:w-[320px] h-[calc(100vh-180px)] sticky top-24 bg-white shadow-md border border-gray-200 rounded-lg p-3 self-start">
              <CartComponent />
            </div>
          )}
        </div>
        
        {/* Mobile cart button */}
        {isMobile && hasItems() && (
          <div className="fixed bottom-4 right-4 z-50">
            <Button className="shadow-lg flex items-center gap-2 bg-food-primary hover:bg-food-accent">
              <ShoppingBag className="h-5 w-5" />
              <span>{getTotalItems()}</span>
              <span className="bg-white text-food-primary rounded-full w-6 h-6 flex items-center justify-center text-xs font-bold">
                {getTotalItems()}
              </span>
              <span className="font-bold">${getTotalPrice().toFixed(2)}</span>
            </Button>
          </div>
        )}
      </section>
    </Layout>
  );
};

export default Menu;
