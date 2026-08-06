import React, { useState } from 'react';
import { useCart } from '@/context/CartContext';
import { useRestaurantDetails } from '@/hooks/useRestaurantDetails';
import CartItemComponent from '@/components/CartItemComponent';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { AlertCircle, ShoppingCart, ArrowRight, Package, MapPin, ArrowLeft, Truck, Clock, Tag } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { toast } from '@/components/ui/use-toast';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

// Default values for restaurant settings
const DEFAULT_DELIVERY_FEE = 3.99;
const DEFAULT_MINIMUM_ORDER = 15.00;
const DEFAULT_DELIVERY_TIME = { min: 30, max: 45 };
const DEFAULT_PICKUP_TIME = { min: 15, max: 25 };

const CartMobile: React.FC = () => {
  const { items, getTotalPrice, hasItems, clearCart, orderType, setOrderType, restaurantId } = useCart();
  const [promoCode, setPromoCode] = useState('');
  const [isApplying, setIsApplying] = useState(false);
  const navigate = useNavigate();

  const { data: restaurant } = useRestaurantDetails(restaurantId || "");
  const subtotal = getTotalPrice();
  const taxRate = restaurant?.taxRate ?? 0.10;
  const tax = taxRate > 0 ? subtotal * taxRate : 0;
  const total = subtotal + tax;

  const handleApplyPromo = () => {
    setIsApplying(true);
    setTimeout(() => {
      setIsApplying(false);
      if (promoCode.toLowerCase() === 'discount20') {
        toast({
          title: "Promo code applied!",
          description: "You received 20% discount on your order.",
          variant: "default",
        });
      } else {
        toast({
          title: "Invalid promo code",
          description: "Please try a different code.",
          variant: "destructive",
        });
      }
    }, 1000);
  };

  const handleCheckout = () => {
    if (subtotal < DEFAULT_MINIMUM_ORDER) {
      toast({
        title: "Minimum order not met",
        description: `Please add more items to meet the minimum order amount of $${DEFAULT_MINIMUM_ORDER.toFixed(2)}`,
        variant: "destructive",
      });
      return;
    }
    
    navigate('/mobile-checkout');
  };

  const getEstimatedTime = () => {
    if (orderType === 'delivery') {
      return 'Calculated at checkout';
    } else {
      return `${DEFAULT_PICKUP_TIME.min}-${DEFAULT_PICKUP_TIME.max} minutes`;
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100">
      {/* Main Container - Mobile Design */}
      <div className="max-w-md mx-auto bg-white shadow-2xl min-h-screen">
        
        {/* Header */}
        <div className="bg-gradient-to-r from-food-primary to-food-accent text-white sticky top-0 z-30">
          <div className="px-6 py-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                {restaurantId && (
                  <Link to={`/mobile-restaurant/${restaurantId}`}>
                    <Button
                      variant="ghost"
                      size="sm"
                      className="text-white hover:bg-white/20 p-2 rounded-full"
                    >
                      <ArrowLeft className="w-5 h-5" />
                    </Button>
                  </Link>
                )}
                <div>
                  <h1 className="text-xl font-bold">Your Cart</h1>
                  <p className="text-white/90 text-sm">{items.length} items</p>
                </div>
              </div>
              <div className="text-right">
                <p className="text-white/90 text-xs">Restaurant</p>
                <div className="flex items-center text-white/80 text-xs">
                  <Clock className="w-3 h-3 mr-1" />
                  <span>{getEstimatedTime()}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {hasItems() ? (
          <div className="px-6 pb-6">
            
            {/* Cart Items */}
            <Card className="mt-6 overflow-hidden border-0 shadow-lg">
              <CardHeader className="bg-gradient-to-r from-gray-50 to-gray-100 pb-3">
                <CardTitle className="text-lg font-semibold flex items-center justify-between">
                  <div className="flex items-center">
                    <ShoppingCart className="w-5 h-5 mr-2 text-food-primary" />
                    Order Items
                  </div>
                  <Button 
                    variant="ghost"
                    size="sm"
                    onClick={clearCart}
                    className="text-red-600 hover:text-red-700 hover:bg-red-50 px-3 py-1 h-8"
                  >
                    Clear All
                  </Button>
                </CardTitle>
              </CardHeader>
              <CardContent className="p-0">
                <div className="divide-y divide-gray-100">
                  {items.map((item) => (
                    <div key={item.cartItemId} className="p-4">
                      <CartItemComponent item={item} />
                    </div>
                  ))}
                </div>
                <div className="p-4 bg-gray-50 border-t">
                  {restaurantId && (
                    <Link to={`/mobile-restaurant/${restaurantId}`}>
                      <Button variant="outline" className="w-full border-dashed border-2 border-food-primary text-food-primary hover:bg-food-primary/10">
                        + Add More Items
                      </Button>
                    </Link>
                  )}
                </div>
              </CardContent>
            </Card>

            {/* Order Type Selection */}
            <Card className="mt-6 overflow-hidden border-0 shadow-lg">
              <CardHeader className="bg-gradient-to-r from-gray-50 to-gray-100 pb-3">
                <CardTitle className="text-lg font-semibold flex items-center">
                  <Truck className="w-5 h-5 mr-2 text-food-primary" />
                  Order Type
                </CardTitle>
              </CardHeader>
              <CardContent className="p-4">
                <RadioGroup 
                  value={orderType} 
                  onValueChange={(value) => setOrderType(value as 'delivery' | 'pickup')}
                  className="space-y-3"
                >
                  <div className={`flex items-center space-x-3 border-2 rounded-xl p-4 transition-all cursor-pointer ${
                    orderType === 'delivery' 
                      ? 'border-food-primary bg-food-primary/5' 
                      : 'border-gray-200 hover:border-gray-300'
                  }`}>
                    <RadioGroupItem value="delivery" id="delivery" className="text-food-primary" />
                    <label htmlFor="delivery" className="flex items-center justify-between w-full cursor-pointer">
                      <div className="flex items-center">
                        <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center mr-3">
                          <MapPin className="h-5 w-5 text-blue-600" />
                        </div>
                        <div>
                          <p className="font-semibold text-gray-900">Delivery</p>
                          <p className="text-sm text-gray-500">Delivered to your address</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="text-xs font-medium text-gray-500">Calculated at checkout</p>
                      </div>
                    </label>
                  </div>
                  
                  <div className={`flex items-center space-x-3 border-2 rounded-xl p-4 transition-all cursor-pointer ${
                    orderType === 'pickup' 
                      ? 'border-food-primary bg-food-primary/5' 
                      : 'border-gray-200 hover:border-gray-300'
                  }`}>
                    <RadioGroupItem value="pickup" id="pickup" className="text-food-primary" />
                    <label htmlFor="pickup" className="flex items-center justify-between w-full cursor-pointer">
                      <div className="flex items-center">
                        <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center mr-3">
                          <Package className="h-5 w-5 text-green-600" />
                        </div>
                        <div>
                          <p className="font-semibold text-gray-900">Pickup</p>
                          <p className="text-sm text-gray-500">Collect from restaurant</p>
                        </div>
                      </div>
                      <div className="text-right">
                        <p className="font-semibold text-green-600">Free</p>
                        <p className="text-xs text-gray-500">no fee</p>
                      </div>
                    </label>
                  </div>
                </RadioGroup>
              </CardContent>
            </Card>

            {/* Promo Code */}
            <Card className="mt-6 overflow-hidden border-0 shadow-lg">
              <CardContent className="p-4">
                <div className="flex items-center mb-3">
                  <Tag className="w-4 h-4 mr-2 text-food-primary" />
                  <h3 className="font-semibold text-gray-900">Promo Code</h3>
                </div>
                <div className="flex gap-3">
                  <Input
                    placeholder="Enter promo code"
                    value={promoCode}
                    onChange={(e) => setPromoCode(e.target.value)}
                    className="flex-1 border-gray-300 focus:border-food-primary"
                  />
                  <Button 
                    variant="outline"
                    onClick={handleApplyPromo}
                    disabled={!promoCode || isApplying}
                    className="border-food-primary text-food-primary hover:bg-food-primary hover:text-white px-6"
                  >
                    {isApplying ? 'Applying...' : 'Apply'}
                  </Button>
                </div>
              </CardContent>
            </Card>

            {/* Order Summary */}
            <Card className="mt-6 overflow-hidden border-0 shadow-lg">
              <CardHeader className="bg-gradient-to-r from-gray-50 to-gray-100 pb-3">
                <CardTitle className="text-lg font-semibold text-gray-900">Order Summary</CardTitle>
              </CardHeader>
              <CardContent className="p-4 space-y-4">
                <div className="space-y-3">
                  <div className="flex justify-between items-center">
                    <span className="text-gray-600">Subtotal</span>
                    <span className="font-semibold">${subtotal.toFixed(2)}</span>
                  </div>
                  {orderType === 'delivery' && (
                    <div className="flex justify-between items-center">
                      <span className="text-gray-600">Delivery Fee</span>
                      <span className="font-semibold text-sm text-gray-500">Calculated at checkout</span>
                    </div>
                  )}
                  {tax > 0 && (
                    <div className="flex justify-between items-center">
                      <span className="text-gray-600">Tax ({(taxRate * 100).toFixed(0)}%)</span>
                      <span className="font-semibold">${tax.toFixed(2)}</span>
                    </div>
                  )}
                  
                  <div className="border-t border-gray-200 pt-3">
                    <div className="flex justify-between items-center">
                      <span className="text-lg font-bold text-gray-900">Total</span>
                      <span className="text-xl font-bold text-food-primary">${total.toFixed(2)}</span>
                    </div>
                  </div>
                </div>
                
                {subtotal < DEFAULT_MINIMUM_ORDER && (
                  <div className="bg-amber-50 border border-amber-200 rounded-xl p-4 mt-4">
                    <div className="flex items-start gap-3">
                      <AlertCircle className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
                      <div>
                        <p className="font-semibold text-amber-800 text-sm">Minimum Order Required</p>
                        <p className="text-amber-700 text-sm mt-1">
                          Add ${(DEFAULT_MINIMUM_ORDER - subtotal).toFixed(2)} more to meet the minimum order of ${DEFAULT_MINIMUM_ORDER.toFixed(2)}
                        </p>
                      </div>
                    </div>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-16 px-6 text-center">
            <div className="w-24 h-24 rounded-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center mb-6">
              <ShoppingCart className="w-12 h-12 text-gray-400" />
            </div>
            <h2 className="text-2xl font-bold text-gray-900 mb-2">Your cart is empty</h2>
            <p className="text-gray-600 mb-8 max-w-sm">Discover delicious meals from our menu and start building your perfect order!</p>
            {restaurantId && (
              <Link to={`/mobile-restaurant/${restaurantId}`}>
                <Button 
                  size="lg" 
                  className="bg-gradient-to-r from-food-primary to-food-accent text-white px-8 py-3 rounded-xl shadow-lg hover:shadow-xl transition-all"
                >
                  Browse Menu
                  <ArrowRight className="w-5 h-5 ml-2" />
                </Button>
              </Link>
            )}
          </div>
        )}

        {/* Sticky Bottom Checkout Button */}
        {hasItems() && (
          <div className="sticky bottom-0 bg-white border-t border-gray-200 p-6 shadow-lg">
            <Button 
              className="w-full bg-gradient-to-r from-food-primary to-food-accent text-white py-4 text-lg font-semibold rounded-xl shadow-lg hover:shadow-xl transition-all disabled:from-gray-300 disabled:to-gray-400"
              size="lg"
              onClick={handleCheckout}
              disabled={subtotal < DEFAULT_MINIMUM_ORDER}
            >
              <div className="flex items-center justify-between w-full">
                <span>Proceed to Checkout</span>
                <div className="flex items-center">
                  <span className="mr-2">${total.toFixed(2)}</span>
                  <ArrowRight className="w-5 h-5" />
                </div>
              </div>
            </Button>
            {orderType === 'pickup' && (
              <p className="text-center text-sm text-gray-500 mt-3">
                Estimated pickup time: {getEstimatedTime()}
              </p>
            )}
          </div>
        )}
      </div>
    </div>
  );
};

export default CartMobile;
