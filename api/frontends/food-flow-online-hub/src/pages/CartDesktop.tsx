import React, { useState } from 'react';
import { useCart } from '@/context/CartContext';
import { useRestaurantDetails } from '@/hooks/useRestaurantDetails';
import CartItemComponent from '@/components/CartItemComponent';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import Layout from '@/components/Layout';
import { AlertCircle, ShoppingCart, ArrowRight, Package, MapPin, Truck, Tag } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';
import { toast } from '@/components/ui/use-toast';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';

// Default values for restaurant settings
const DEFAULT_DELIVERY_FEE = 3.99;
const DEFAULT_DELIVERY_TIME = { min: 30, max: 45 };
const DEFAULT_PICKUP_TIME = { min: 15, max: 25 };

const CartDesktop: React.FC = () => {
  const { items, getTotalPrice, hasItems, clearCart, orderType, setOrderType, restaurantId, appliedPromo, applyPromoCode, removePromoCode } = useCart();
  const [promoInput, setPromoInput] = useState('');
  const [promoError, setPromoError] = useState('');
  const [isApplying, setIsApplying] = useState(false);
  const navigate = useNavigate();

  const { data: restaurant } = useRestaurantDetails(restaurantId || "");
  const minSpend = restaurant?.minSpend ?? 0;
  const subtotal = getTotalPrice();
  const discount = appliedPromo ? appliedPromo.discountAmount : 0;
  const taxableSubtotal = Math.max(0, subtotal - discount);
  const taxRate = restaurant?.taxRate ?? 0.10;
  const tax = taxRate > 0 ? taxableSubtotal * taxRate : 0;
  const total = taxableSubtotal + tax;

  const handleApplyPromo = async () => {
    setIsApplying(true);
    setPromoError('');
    const res = await applyPromoCode(promoInput);
    setIsApplying(false);
    if (!res.success) {
      setPromoError(res.message);
    } else {
      setPromoInput('');
    }
  };

  const handleCheckout = () => {
    if (restaurant?.enabled === false) {
      toast({
        title: "Restaurant Paused",
        description: `${restaurant.name || "This restaurant"} is currently paused and not accepting orders.`,
        variant: "destructive",
      });
      return;
    }
    if (minSpend > 0 && subtotal < minSpend) {
      toast({
        title: "Minimum order not met",
        description: `Please add more items to meet the minimum order amount of $${minSpend.toFixed(2)}`,
        variant: "destructive",
      });
      return;
    }
    
    navigate('/checkout');
  };

  const getEstimatedTime = () => {
    if (orderType === 'delivery') {
      return 'Calculated at checkout';
    } else {
      return `${DEFAULT_PICKUP_TIME.min}-${DEFAULT_PICKUP_TIME.max} minutes`;
    }
  };

  return (
    <Layout>
      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Your Cart</h1>
          <p className="text-gray-600 mt-2">{items.length} items in your cart</p>
        </div>

        {restaurant?.enabled === false && (
          <div className="mb-6">
            <div className="bg-red-50 border border-red-200 rounded-xl p-4 text-red-900 flex items-center justify-between gap-3 shadow-sm">
              <div className="flex items-center gap-3">
                <AlertCircle className="w-5 h-5 text-red-600 shrink-0" />
                <div>
                  <h3 className="font-bold text-sm">Restaurant Paused</h3>
                  <p className="text-xs text-red-700 mt-0.5">{restaurant.name} is currently paused and not taking new orders. Checkout is temporarily disabled.</p>
                </div>
              </div>
              <span className="text-[10px] font-bold uppercase bg-red-200 text-red-900 px-2.5 py-1 rounded-full shrink-0">Paused</span>
            </div>
          </div>
        )}

        {hasItems() ? (
          <div className="grid grid-cols-3 gap-8">
            
            {/* Left Column: Cart Items */}
            <div className="col-span-2 space-y-6">
              
              {/* Cart Items */}
              <Card className="overflow-hidden border shadow-md">
                <CardHeader className="bg-gray-50 border-b">
                  <CardTitle className="text-xl font-semibold flex items-center justify-between">
                    <div className="flex items-center">
                      <ShoppingCart className="w-6 h-6 mr-3 text-food-primary" />
                      Order Items
                    </div>
                    <Button 
                      variant="outline"
                      size="sm"
                      onClick={clearCart}
                      className="text-red-600 hover:text-red-700 hover:bg-red-50"
                    >
                      Clear All
                    </Button>
                  </CardTitle>
                </CardHeader>
                <CardContent className="p-0">
                  <div className="divide-y divide-gray-200">
                    {items.map((item) => (
                      <div key={item.cartItemId} className="p-6">
                        <CartItemComponent item={item} />
                      </div>
                    ))}
                  </div>
                  <div className="p-6 bg-gray-50 border-t">
                    {restaurantId && (
                      <Link to={`/restaurant/${restaurantId}`}>
                        <Button variant="outline" className="w-full border-dashed border-2 border-food-primary text-food-primary hover:bg-food-primary/10">
                          + Add More Items
                        </Button>
                      </Link>
                    )}
                  </div>
                </CardContent>
              </Card>

              {/* Order Type Selection */}
              <Card className="overflow-hidden border shadow-md">
                <CardHeader className="bg-gray-50 border-b">
                  <CardTitle className="text-xl font-semibold flex items-center">
                    <Truck className="w-6 h-6 mr-3 text-food-primary" />
                    Order Type
                  </CardTitle>
                </CardHeader>
                <CardContent className="p-6">
                  <RadioGroup 
                    value={orderType} 
                    onValueChange={(value) => setOrderType(value as 'delivery' | 'pickup')}
                    className="grid grid-cols-2 gap-4"
                  >
                    <div className={`flex items-center space-x-3 border-2 rounded-lg p-4 transition-all cursor-pointer ${
                      orderType === 'delivery' 
                        ? 'border-food-primary bg-food-primary/5' 
                        : 'border-gray-200 hover:border-gray-300'
                    }`}>
                      <RadioGroupItem value="delivery" id="delivery-desktop" className="text-food-primary" />
                      <label htmlFor="delivery-desktop" className="flex-1 cursor-pointer">
                        <div className="flex items-center mb-2">
                          <div className="w-10 h-10 rounded-full bg-blue-100 flex items-center justify-center mr-3">
                            <MapPin className="h-5 w-5 text-blue-600" />
                          </div>
                          <p className="font-semibold text-gray-900">Delivery</p>
                        </div>
                        <p className="text-sm text-gray-500 ml-13">Delivered to your address</p>
                        <p className="text-sm font-medium text-gray-500 mt-2">Calculated at checkout</p>
                      </label>
                    </div>
                    
                    <div className={`flex items-center space-x-3 border-2 rounded-lg p-4 transition-all cursor-pointer ${
                      orderType === 'pickup' 
                        ? 'border-food-primary bg-food-primary/5' 
                        : 'border-gray-200 hover:border-gray-300'
                    }`}>
                      <RadioGroupItem value="pickup" id="pickup-desktop" className="text-food-primary" />
                      <label htmlFor="pickup-desktop" className="flex-1 cursor-pointer">
                        <div className="flex items-center mb-2">
                          <div className="w-10 h-10 rounded-full bg-green-100 flex items-center justify-center mr-3">
                            <Package className="h-5 w-5 text-green-600" />
                          </div>
                          <p className="font-semibold text-gray-900">Pickup</p>
                        </div>
                        <p className="text-sm text-gray-500 ml-13">Collect from restaurant</p>
                      </label>
                    </div>
                  </RadioGroup>
                </CardContent>
              </Card>

              {/* Items List */}
              <div className="space-y-4">
                {items.map((item) => (
                  <CartItemComponent key={item.cartItemId} item={item} />
                ))}
              </div>
            </div>

            {/* Sidebar Summary (Right 1 col) */}
            <div className="col-span-1">
              <div className="sticky top-24 space-y-6">
                
                {/* Promo Code */}
                <Card className="overflow-hidden border shadow-md">
                  <CardContent className="p-6">
                    <div className="flex items-center mb-4">
                      <Tag className="w-5 h-5 mr-2 text-food-primary" />
                      <h3 className="font-semibold text-gray-900">Promo Code</h3>
                    </div>
                    {appliedPromo ? (
                      <div className="flex items-center justify-between bg-green-50 border border-green-200 rounded-lg p-3">
                        <div>
                          <p className="font-semibold text-green-800 text-sm">Code: {appliedPromo.code}</p>
                          <p className="text-green-700 text-xs">Saved ${appliedPromo.discountAmount.toFixed(2)}</p>
                        </div>
                        <Button variant="ghost" size="sm" onClick={removePromoCode} className="text-red-600 hover:text-red-800 text-xs">
                          Remove
                        </Button>
                      </div>
                    ) : (
                      <div className="space-y-3">
                        <Input
                          placeholder="Enter promo code"
                          value={promoInput}
                          onChange={(e) => setPromoInput(e.target.value)}
                          className="border-gray-300 focus:border-food-primary"
                        />
                        {promoError && (
                          <p className="text-red-500 text-xs">{promoError}</p>
                        )}
                        <Button 
                          variant="outline"
                          onClick={handleApplyPromo}
                          disabled={!promoInput.trim() || isApplying}
                          className="w-full border-food-primary text-food-primary hover:bg-food-primary hover:text-white"
                        >
                          {isApplying ? 'Applying...' : 'Apply Code'}
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>

                {/* Order Summary */}
                <Card className="overflow-hidden border shadow-md">
                  <CardHeader className="bg-gray-50 border-b">
                    <CardTitle className="text-xl font-semibold text-gray-900">Order Summary</CardTitle>
                  </CardHeader>
                  <CardContent className="p-6 space-y-4">
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-gray-600">Subtotal</span>
                        <span className="font-semibold text-lg">${subtotal.toFixed(2)}</span>
                      </div>
                      {discount > 0 && (
                        <div className="flex justify-between items-center text-green-600 font-medium">
                          <span>Discount ({appliedPromo?.code})</span>
                          <span>-${discount.toFixed(2)}</span>
                        </div>
                      )}
                      {orderType === 'delivery' && (
                        <div className="flex justify-between items-center">
                          <span className="text-gray-600">Delivery Fee</span>
                          <span className="font-semibold text-sm text-gray-500">Calculated at checkout</span>
                        </div>
                      )}
                      {tax > 0 && (
                        <div className="flex justify-between items-center">
                          <span className="text-gray-600">Tax ({(taxRate * 100).toFixed(0)}%)</span>
                          <span className="font-semibold text-lg">${tax.toFixed(2)}</span>
                        </div>
                      )}
                      
                      <div className="border-t border-gray-200 pt-3 mt-4">
                        <div className="flex justify-between items-center">
                          <span className="text-xl font-bold text-gray-900">Total</span>
                          <span className="text-2xl font-bold text-food-primary">${total.toFixed(2)}</span>
                        </div>
                      </div>
                    </div>
                    
                    {minSpend > 0 && subtotal < minSpend && (
                      <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 mt-4">
                        <div className="flex items-start gap-3">
                          <AlertCircle className="w-5 h-5 text-amber-600 flex-shrink-0 mt-0.5" />
                          <div>
                            <p className="font-semibold text-amber-800 text-sm">Minimum Order Required</p>
                            <p className="text-amber-700 text-sm mt-1">
                              Add ${(minSpend - subtotal).toFixed(2)} more to meet the minimum of ${minSpend.toFixed(2)}
                            </p>
                          </div>
                        </div>
                      </div>
                    )}

                    <Button 
                      className="w-full bg-food-primary hover:bg-food-accent text-white py-6 text-lg font-semibold rounded-lg mt-6 disabled:bg-gray-400"
                      size="lg"
                      onClick={handleCheckout}
                      disabled={(minSpend > 0 && subtotal < minSpend) || restaurant?.enabled === false}
                    >
                      {restaurant?.enabled === false ? 'Restaurant Paused' : 'Proceed to Checkout'}
                      {restaurant?.enabled !== false && <ArrowRight className="w-5 h-5 ml-2" />}
                    </Button>

                    {orderType === 'pickup' && (
                      <p className="text-center text-sm text-gray-500 mt-3">
                        Estimated pickup: {getEstimatedTime()}
                      </p>
                    )}
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <div className="w-32 h-32 rounded-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center mb-6">
              <ShoppingCart className="w-16 h-16 text-gray-400" />
            </div>
            <h2 className="text-3xl font-bold text-gray-900 mb-3">Your cart is empty</h2>
            <p className="text-gray-600 mb-8 max-w-md text-lg">Discover delicious meals from our menu and start building your perfect order!</p>
            {restaurantId && (
              <Link to={`/restaurant/${restaurantId}`}>
                <Button 
                  size="lg" 
                  className="bg-food-primary hover:bg-food-accent text-white px-10 py-6 text-lg rounded-lg shadow-lg hover:shadow-xl transition-all"
                >
                  Browse Menu
                  <ArrowRight className="w-6 h-6 ml-2" />
                </Button>
              </Link>
            )}
          </div>
        )}
      </div>
    </Layout>
  );
};

export default CartDesktop;
