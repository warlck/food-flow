import React, { useState, useEffect } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { Card, CardContent } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { useCart } from '@/context/CartContext';
import { useRestaurantDetails } from '@/hooks/useRestaurantDetails';
import { transformApiRestaurant, transformApiMenuItems } from '@/lib/transformers';
import { MenuItem as MenuItemType } from '@/types';
import { Plus, Minus, ShoppingCart, Clock, Star, ArrowUp, Loader2, AlertCircle } from 'lucide-react';
import { toast } from '@/components/ui/use-toast';
import MenuItemDialog from '@/components/MenuItemDialog';

const MobileMenu: React.FC = () => {
  const { restaurantId: urlRestaurantId } = useParams<{ restaurantId: string }>();
  const [searchParams] = useSearchParams();
  const queryRestaurantId = searchParams.get('restaurant_id');
  
  // Use restaurant_id from URL path or query param - no default
  const restaurantId = urlRestaurantId || queryRestaurantId;
  
  // Fetch restaurant details from API only if restaurantId is provided
  const { data: apiData, isLoading, error } = useRestaurantDetails(restaurantId || '');
  
  const [selectedCategory, setSelectedCategory] = useState('All');
  const [selectedItem, setSelectedItem] = useState<MenuItemType | null>(null);
  const [isDialogOpen, setIsDialogOpen] = useState(false);
  const { addToCart, updateQuantity, removeFromCart, items, getTotalItems, getTotalPrice, setRestaurantId } = useCart();

  // Set restaurant ID in cart context when component mounts or restaurantId changes
  useEffect(() => {
    if (restaurantId) {
      setRestaurantId(restaurantId);
    }
  }, [restaurantId, setRestaurantId]);

  // Check if restaurant ID is missing
  if (!restaurantId) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 flex items-center justify-center p-4">
        <Alert variant="destructive" className="max-w-md">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Restaurant ID Required</AlertTitle>
          <AlertDescription>
            Please provide a valid restaurant ID in the URL. 
            The URL should be in the format: /mobile-restaurant/[restaurant-id]
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  // Loading state
  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 flex items-center justify-center">
        <div className="text-center">
          <Loader2 className="h-12 w-12 animate-spin mx-auto mb-4 text-food-primary" />
          <h2 className="text-xl font-semibold">Loading restaurant...</h2>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 flex items-center justify-center p-4">
        <Alert variant="destructive" className="max-w-md">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>Error loading restaurant</AlertTitle>
          <AlertDescription>
            {error.message || 'Failed to load restaurant details. Please try again later.'}
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  // Require API data - no fallback to mock
  if (!apiData) {
    return (
      <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100 flex items-center justify-center p-4">
        <Alert variant="destructive" className="max-w-md">
          <AlertCircle className="h-4 w-4" />
          <AlertTitle>No data available</AlertTitle>
          <AlertDescription>
            Unable to load restaurant data. Please try again later.
          </AlertDescription>
        </Alert>
      </div>
    );
  }

  // Transform API data
  const restaurant = transformApiRestaurant(apiData);
  const { items: menuItems, categories } = transformApiMenuItems(apiData);

  const getItemQuantity = (itemId: string) => {
    const cartItem = items.find(item => item.menuItem.id === itemId);
    return cartItem ? cartItem.quantity : 0;
  };

  const handleAddItem = (item: MenuItemType) => {
    addToCart(item, 1);
    toast({
      title: "Added to cart",
      description: `${item.name} has been added to your cart.`,
      duration: 1500,
    });
  };

  const handleUpdateQuantity = (item: MenuItemType, newQuantity: number) => {
    if (newQuantity === 0) {
      removeFromCart(item.id);
    } else {
      updateQuantity(item.id, newQuantity);
    }
  };

  const filteredItems = menuItems.filter(item => 
    selectedCategory === 'All' || item.category === selectedCategory
  );

  const categoriesWithAll = ['All', ...categories];

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100">
      {/* Main Container - Centralized Design */}
      <div className="max-w-md mx-auto bg-white shadow-2xl min-h-screen relative">
        
        {/* Elegant Header */}
        <div className="bg-gradient-to-r from-food-primary to-food-accent text-white sticky top-0 z-30 shadow-lg">
          <div className="px-6 py-6">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center space-x-3">
                <div className="w-12 h-12 bg-white/20 rounded-full flex items-center justify-center backdrop-blur-sm">
                  <Star className="w-6 h-6 text-white" />
                </div>
                <div>
                  <h1 className="text-xl font-bold">{restaurant.name} (v0.0.4)</h1>
                  <div className="flex items-center space-x-2 text-white/90 text-sm">
                    <Clock className="w-4 h-4" />
                    <span>25-35 min</span>
                  </div>
                </div>
              </div>
              <div className="text-right">
                <div className="flex items-center space-x-1 text-sm">
                  <Star className="w-4 h-4 fill-current" />
                  <span className="font-semibold">{restaurant.rating}</span>
                </div>
                <p className="text-white/80 text-xs">Pick Up</p>
              </div>
            </div>
            <p className="text-white/90 text-sm">{restaurant.address}</p>
          </div>
        </div>

        {/* Category Navigation */}
        <div className="bg-white border-b sticky top-[132px] z-20 shadow-sm">
          <div className="overflow-x-auto">
            <div className="flex px-4 py-4 space-x-3 min-w-max">
              {categoriesWithAll.map((category) => (
                <button
                  key={category}
                  onClick={() => setSelectedCategory(category)}
                  className={`px-6 py-3 rounded-full text-sm font-medium whitespace-nowrap transition-all duration-300 border-2 ${
                    selectedCategory === category
                      ? 'bg-food-primary text-white border-food-primary shadow-lg scale-105'
                      : 'bg-white text-gray-700 border-gray-200 hover:border-food-primary hover:text-food-primary hover:shadow-md'
                  }`}
                >
                  {category}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Menu Items */}
        <div className="px-4 py-6 pb-32">
          <div className="space-y-4">
            {filteredItems.map((item) => {
              const quantity = getItemQuantity(item.id);
              return (
                <Card 
                  key={item.id} 
                  className="overflow-hidden hover:shadow-xl transition-all duration-300 border-0 shadow-md cursor-pointer"
                  onClick={() => {
                    console.log('Card clicked:', item.name);
                    setSelectedItem(item);
                    setIsDialogOpen(true);
                  }}
                >
                  <CardContent className="p-0">
                    <div className="relative">
                      {/* Item Image */}
                      <div className="relative w-full h-48 overflow-hidden">
                        <img
                          src={item.image}
                          alt={item.name}
                          className="w-full h-full object-cover transition-transform duration-300 hover:scale-105"
                          onError={(e) => {
                            const target = e.target as HTMLImageElement;
                            target.src = 'https://images.unsplash.com/photo-1504674900247-0877df9cc836?auto=format&fit=crop&q=80';
                          }}
                        />
                        <div className="absolute top-3 right-3">
                          <Badge className="bg-white/90 text-food-primary border-0 font-semibold">
                            {item.category}
                          </Badge>
                        </div>
                        {!item.available && (
                          <div className="absolute inset-0 bg-black/50 flex items-center justify-center">
                            <span className="text-white font-semibold">Not Available</span>
                          </div>
                        )}
                      </div>

                      {/* Item Details */}
                      <div className="p-4">
                        <div className="mb-3">
                          <h3 className="font-bold text-gray-900 text-lg leading-tight mb-2">{item.name}</h3>
                          <p className="text-gray-600 text-sm leading-relaxed">{item.description}</p>
                        </div>

                        {/* Price and Controls */}
                        <div className="flex items-center justify-between">
                          <div className="text-food-primary font-bold text-xl">
                            ${item.price.toFixed(2)}
                          </div>

                          {/* Quantity Controls */}
                          <div className="flex items-center">
                            {quantity === 0 ? (
                              <Button
                                onClick={(e) => {
                                  e.stopPropagation();
                                  console.log('Add button clicked:', item.name);
                                  setSelectedItem(item);
                                  setIsDialogOpen(true);
                                }}
                                disabled={!item.available}
                                className="bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white px-6 py-2 rounded-full font-semibold shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-105"
                              >
                                <Plus className="w-4 h-4 mr-2" />
                                Add
                              </Button>
                            ) : (
                              <div className="flex items-center space-x-3 bg-gray-50 rounded-full p-2 border-2 border-food-primary/20">
                                <Button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleUpdateQuantity(item, quantity - 1);
                                  }}
                                  size="sm"
                                  variant="ghost"
                                  className="w-8 h-8 rounded-full p-0 hover:bg-food-primary hover:text-white transition-colors"
                                >
                                  <Minus className="w-4 h-4" />
                                </Button>
                                <span className="w-8 text-center font-bold text-food-primary text-lg">{quantity}</span>
                                <Button
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    handleUpdateQuantity(item, quantity + 1);
                                  }}
                                  size="sm"
                                  variant="ghost"
                                  className="w-8 h-8 rounded-full p-0 hover:bg-food-primary hover:text-white transition-colors"
                                >
                                  <Plus className="w-4 h-4" />
                                </Button>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        </div>

        {/* Floating Cart Button */}
        {getTotalItems() > 0 && (
          <div className="fixed bottom-4 left-1/2 transform -translate-x-1/2 z-40 w-full max-w-md px-4">
            <Button
              className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-2xl font-bold text-lg shadow-2xl hover:shadow-3xl transition-all duration-300 transform hover:scale-105 border-0"
              onClick={() => window.location.href = '/mobile-cart'}
            >
              <div className="flex items-center justify-between w-full">
                <div className="flex items-center space-x-3">
                  <div className="bg-white/20 rounded-full p-2">
                    <ShoppingCart className="w-5 h-5" />
                  </div>
                  <span>{getTotalItems()} items</span>
                </div>
                <div className="flex items-center space-x-2">
                  <span>${getTotalPrice().toFixed(2)}</span>
                  <ArrowUp className="w-5 h-5 rotate-45" />
                </div>
              </div>
            </Button>
          </div>
        )}
      </div>
      {/* Menu Item Dialog */}
      <MenuItemDialog
        item={selectedItem}
        isOpen={isDialogOpen}
        onClose={() => setIsDialogOpen(false)}
      />
    </div>
  );
};

export default MobileMenu;
