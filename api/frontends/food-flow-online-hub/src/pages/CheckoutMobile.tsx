import React, { useState } from "react";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { MapPin, Clock, CreditCard, ArrowLeft, CheckCircle, User, Phone, Mail, Home, MapPinCheck, ShoppingBag, Loader2, Search } from "lucide-react";
import { loadStripe } from "@stripe/stripe-js";
import { Elements } from "@stripe/react-stripe-js";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useCart } from "@/context/CartContext";
import { useRestaurantDetails } from "@/hooks/useRestaurantDetails";
import { Label } from "@/components/ui/label";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { orderService, Order, DeliveryQuote, CreateOrderRequest } from "@/services/orderService";
import { searchAddress, GeocodingResult } from "@/lib/geocoding";
import StripePaymentForm from "@/components/StripePaymentForm";

// Initialize Stripe (may be null when key isn't configured)
const stripePromise = import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY
  ? loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY)
  : Promise.resolve(null);

// Default values for restaurant settings
const DEFAULT_DELIVERY_TIME = { min: 30, max: 45 };
const DEFAULT_PICKUP_TIME = { min: 15, max: 25 };

// Form schemas
const deliveryFormSchema = z.object({
  name: z.string().min(2, { message: "Name must be at least 2 characters" }),
  email: z.string().email({ message: "Invalid email address" }),
  phone: z.string().min(10, { message: "Phone number must be at least 10 characters" }),
  street: z.string().min(5, { message: "Street address must be at least 5 characters" }),
  city: z.string().min(2, { message: "City must be at least 2 characters" }),
  state: z.string().optional(),
  postalCode: z.string().min(5, { message: "Postal code must be at least 5 characters" }),
  deliveryInstructions: z.string().optional(),
});

const pickupFormSchema = z.object({
  name: z.string().min(2, { message: "Name must be at least 2 characters" }),
  email: z.string().email({ message: "Invalid email address" }),
  phone: z.string().min(10, { message: "Phone number must be at least 10 characters" }),
});

const CheckoutMobile: React.FC = () => {
  const navigate = useNavigate();
  const { items, getTotalPrice, clearCart, orderType, restaurantId, appliedPromo } = useCart();
  const [step, setStep] = useState(1);
  const [paymentMethod, setPaymentMethod] = useState<"creditCard" | "payAtLocation">("creditCard");
  const [orderDetails, setOrderDetails] = useState<Record<string, unknown>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Stripe-related state
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null);
  const [clientSecret, setClientSecret] = useState<string | null>(null);

  // Delivery destination state
  const [addressQuery, setAddressQuery] = useState("");
  const [searchResults, setSearchResults] = useState<GeocodingResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState<GeocodingResult | null>(null);
  const [deliveryQuote, setDeliveryQuote] = useState<DeliveryQuote | null>(null);

  const { data: restaurant } = useRestaurantDetails(restaurantId || "");
  const subtotal = getTotalPrice();
  const discount = appliedPromo ? appliedPromo.discountAmount : 0;
  const taxableSubtotal = Math.max(0, subtotal - discount);
  const deliveryFee = orderType === "delivery" && deliveryQuote?.withinLimit ? deliveryQuote.deliveryFee : 0;
  const taxRate = restaurant?.taxRate ?? 0.10;
  const tax = taxRate > 0 ? taxableSubtotal * taxRate : 0;
  const total = taxableSubtotal + deliveryFee + tax;

  // Setup forms based on order type
  const deliveryForm = useForm<z.infer<typeof deliveryFormSchema>>({
    resolver: zodResolver(deliveryFormSchema),
    defaultValues: {
      name: "",
      email: "",
      phone: "",
      street: "",
      city: "Singapore",
      state: "SG",
      postalCode: "",
      deliveryInstructions: "",
    },
  });

  const pickupForm = useForm<z.infer<typeof pickupFormSchema>>({
    resolver: zodResolver(pickupFormSchema),
    defaultValues: {
      name: "",
      email: "",
      phone: "",
    },
  });

  // Address search handlers
  const handleAddressSearch = async () => {
    const query = addressQuery.trim();
    if (!query) {
      return;
    }

    setIsSearching(true);
    setSearchResults([]);

    try {
      const results = await searchAddress(query);
      setSearchResults(results);
      if (results.length === 0) {
        toast.error("No addresses found in Singapore. Try a 6-digit postal code or street name.");
      }
    } catch (error) {
      console.error("Address search failed:", error);
      toast.error("Address search failed. Please try again.");
    } finally {
      setIsSearching(false);
    }
  };

  const handleSelectAddress = async (result: GeocodingResult) => {
    setSelectedLocation(result);
    setSearchResults([]);
    setAddressQuery(result.displayName);

    // Pre-fill address fields from the geocoded result; the customer can
    // still adjust them manually.
    if (result.address.street) deliveryForm.setValue("street", result.address.street);
    deliveryForm.setValue("city", result.address.city || "Singapore");
    deliveryForm.setValue("state", result.address.state || "SG");
    if (result.address.postalCode) deliveryForm.setValue("postalCode", result.address.postalCode);

    // Fetch the delivery fee quote for the selected destination
    if (restaurantId && restaurant?.latitude !== undefined && restaurant?.longitude !== undefined) {
      try {
        const quote = await orderService.getDeliveryQuote(
          restaurant.latitude,
          restaurant.longitude,
          result.latitude,
          result.longitude,
          restaurant.maxDeliveryDistanceKm
        );
        setDeliveryQuote(quote);
        if (!quote.withinLimit) {
          toast.error(`This address is ${quote.distanceKm.toFixed(1)} km away, outside the ${quote.maxDeliveryDistanceKm} km delivery area.`);
        }
      } catch (error) {
        console.error("Failed to fetch delivery quote:", error);
        setDeliveryQuote(null);
        toast.error("Could not calculate the delivery fee for this address.");
      }
    } else if (restaurantId) {
      toast.error("Restaurant location data is missing, cannot calculate delivery fee.");
    }
  };

  const handleClearAddress = () => {
    setSelectedLocation(null);
    setDeliveryQuote(null);
    setAddressQuery("");
    setSearchResults([]);
    deliveryForm.setValue("street", "");
    deliveryForm.setValue("postalCode", "");
  };

  const onInvalidCustomerInfo = () => {
    if (orderType === "delivery" && !selectedLocation) {
      toast.error("Please search and select your delivery address.");
    }
  };

  // Submit handlers
  const onSubmitCustomerInfo = (data: z.infer<typeof deliveryFormSchema> | z.infer<typeof pickupFormSchema>) => {
    if (restaurant?.enabled === false) {
      toast.error(`${restaurant.name || "This restaurant"} is currently paused and not accepting new orders.`);
      return;
    }

    const minSpend = restaurant?.minSpend ?? 0;
    if (minSpend > 0 && subtotal < minSpend) {
      toast.error(`Minimum spend not met. Subtotal must be at least $${minSpend.toFixed(2)}.`);
      return;
    }

    if (orderType === "delivery") {
      if (!selectedLocation) {
        toast.error("Please search and select your delivery address.");
        return;
      }
      if (deliveryQuote && !deliveryQuote.withinLimit) {
        toast.error("The selected address is outside the delivery area.");
        return;
      }
    }

    setOrderDetails({ ...orderDetails, ...data });
    setStep(2);
  };

  const buildOrderRequest = (method: "creditCard" | "cash"): CreateOrderRequest => {
    const orderData = orderDetails as z.infer<typeof deliveryFormSchema>;

    return {
      restaurantId: restaurantId || "",
      customerName: orderData.name,
      customerEmail: orderData.email,
      customerPhone: orderData.phone,
      orderType: orderType as "delivery" | "pickup",
      paymentMethod: method,
      promoCode: appliedPromo?.code,
      items: items.map((item) => ({
        menuItemId: item.menuItem.id,
        quantity: item.quantity,
        specialInstructions: item.specialInstructions,
        modifiers: item.selectedModifiers?.map((m) => ({
          modifierGroupId: m.modifierGroupId,
          modifierOptionId: m.modifierOptionId,
        })),
        addons: item.selectedAddons?.map(({ addon, quantity: addonQty }) => ({
          addonId: addon.addonId || addon.id,
          quantity: addonQty,
        })),
      })),
      deliveryAddress: orderType === "delivery" && selectedLocation ? {
        street: orderData.street,
        city: orderData.city,
        state: orderData.state && orderData.state.trim() !== "" ? orderData.state.trim() : "SG",
        postalCode: orderData.postalCode,
        deliveryInstructions: orderData.deliveryInstructions,
        latitude: selectedLocation.latitude,
        longitude: selectedLocation.longitude,
      } : undefined,
    };
  };

  const onSubmitPayment = async () => {
    setIsSubmitting(true);

    try {
      if (paymentMethod === "creditCard") {
        // Create the order, then initialize the Stripe PaymentIntent. The
        // Stripe form appears once we have a client secret.
        const order = await orderService.createOrder(buildOrderRequest("creditCard"));
        setCurrentOrder(order);

        try {
          const paymentIntent = await orderService.createPaymentIntent(order.id);
          setClientSecret(paymentIntent.clientSecret);
        } catch (error) {
          console.error("Failed to create payment intent:", error);
          toast.error("Failed to initialize Stripe payment. You can retry below.");
        }
      } else {
        // Pay at location: create the order and confirm it directly.
        const order = await orderService.createOrder(buildOrderRequest("cash"));
        await orderService.confirmPayment(order.id);

        clearCart();
        toast.success("Order placed successfully!");
        navigate(`/order-confirmation/${order.id}`);
      }
    } catch (error) {
      console.error("Failed to place order:", error);
      const message = error instanceof Error ? error.message : "Failed to place order. Please try again.";
      toast.error(message);
    } finally {
      setIsSubmitting(false);
    }
  };

  const retryCreatePaymentIntent = async () => {
    if (!currentOrder) {
      toast.error("No order found. Please go back and try again.");
      return;
    }

    try {
      setIsSubmitting(true);
      const paymentIntent = await orderService.createPaymentIntent(currentOrder.id);
      setClientSecret(paymentIntent.clientSecret);
    } catch (error) {
      console.error("Failed to create payment intent:", error);
      toast.error("Failed to initialize Stripe payment. Please try again.");
    } finally {
      setIsSubmitting(false);
    }
  };

  const handlePaymentSuccess = async () => {
    // Ensure backend order/payment status is updated (useful even without webhooks in dev).
    if (currentOrder) {
      try {
        await orderService.confirmPayment(currentOrder.id);
      } catch (err) {
        // Payment may have succeeded on Stripe but backend confirmation failed.
        console.error("Failed to confirm payment with backend:", err);
      }
    }

    clearCart();
    toast.success("Payment successful!");
    navigate(`/order-confirmation/${currentOrder?.id}`);
  };

  const handlePaymentError = (error: Error) => {
    console.error("Payment error:", error);
    // Keep user on payment page to retry
  };

  const handleBackStep = () => {
    setStep(step - 1);
    // Reset order and payment state if going back
    if (step === 2) {
      setCurrentOrder(null);
      setClientSecret(null);
    }
  };

  const getEstimatedTime = () => {
    if (orderType === "delivery") {
      return `${DEFAULT_DELIVERY_TIME.min}-${DEFAULT_DELIVERY_TIME.max} minutes`;
    } else {
      return `${DEFAULT_PICKUP_TIME.min}-${DEFAULT_PICKUP_TIME.max} minutes`;
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-b from-gray-50 to-gray-100">
      {/* Main Container - Centralized Design */}
      <div className="max-w-md mx-auto bg-white shadow-2xl min-h-screen">

        {/* Header */}
        <div className="bg-gradient-to-r from-food-primary to-food-accent text-white sticky top-0 z-30">
          <div className="px-6 py-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center space-x-3">
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => navigate("/mobile-cart")}
                  className="text-white hover:bg-white/20 p-2 rounded-full"
                >
                  <ArrowLeft className="w-5 h-5" />
                </Button>
                <div>
                  <h1 className="text-xl font-bold">Checkout</h1>
                  <p className="text-white/90 text-sm">{orderType === "delivery" ? "Delivery" : "Pickup"} Order</p>
                </div>
              </div>
              <div className="text-right">
                <div className="flex items-center space-x-1 text-sm">
                  <Clock className="w-4 h-4" />
                  <span className="font-semibold">{getEstimatedTime()}</span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Progress Steps */}
        <div className="px-6 py-4 bg-white border-b">
          <div className="flex items-center justify-center space-x-4">
            <div className="flex items-center">
              <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold ${
                step >= 1 ? 'bg-food-primary text-white' : 'bg-gray-200 text-gray-600'
              }`}>
                {step > 1 ? <CheckCircle className="w-5 h-5" /> : '1'}
              </div>
              <span className={`ml-2 text-sm font-medium ${step === 1 ? 'text-food-primary' : 'text-gray-600'}`}>
                Info
              </span>
            </div>
            <div className="flex-1 h-1 bg-gray-200 rounded-full mx-4">
              <div className={`h-full rounded-full transition-all duration-300 ${
                step >= 2 ? 'bg-food-primary' : 'bg-gray-200'
              }`} style={{ width: step >= 2 ? '100%' : '0%' }} />
            </div>
            <div className="flex items-center">
              <div className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-semibold ${
                step >= 2 ? 'bg-food-primary text-white' : 'bg-gray-200 text-gray-600'
              }`}>
                {step > 2 ? <CheckCircle className="w-5 h-5" /> : '2'}
              </div>
              <span className={`ml-2 text-sm font-medium ${step === 2 ? 'text-food-primary' : 'text-gray-600'}`}>
                Payment
              </span>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="px-6 py-6 pb-72">
          {/* Order Items Summary */}
          <Card className="border-0 shadow-lg mb-6">
            <CardHeader className="pb-2">
              <CardTitle className="flex items-center space-x-2 text-lg">
                <ShoppingBag className="w-5 h-5 text-food-primary" />
                <span>Order Items ({items.length})</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="pt-0">
              <div className="divide-y divide-gray-100">
                {items.map((item) => {
                  // Calculate item total including addons
                  let itemTotal = item.menuItem.price * item.quantity;
                  if (item.selectedAddons) {
                    item.selectedAddons.forEach(({ addon, quantity: addonQty }) => {
                      itemTotal += addon.price * addonQty * item.quantity;
                    });
                  }

                  return (
                    <div key={item.cartItemId} className="py-3">
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <div className="flex items-center gap-2">
                            <span className="font-medium text-sm">{item.quantity}x</span>
                            <span className="font-medium text-sm">{item.menuItem.name}</span>
                          </div>

                          {/* Addons */}
                          {item.selectedAddons && item.selectedAddons.length > 0 && (
                            <div className="ml-6 mt-1 space-y-0.5">
                              {item.selectedAddons.map(({ addon, quantity: addonQty }) => (
                                <div key={addon.id} className="flex justify-between text-xs text-gray-600">
                                  <span>+ {addon.name} x{addonQty}</span>
                                  <span className="text-food-primary">+${(addon.price * addonQty).toFixed(2)}</span>
                                </div>
                              ))}
                            </div>
                          )}

                          {/* Special Instructions */}
                          {item.specialInstructions && (
                            <div className="ml-6 mt-1 text-xs text-gray-500 italic">
                              Note: {item.specialInstructions}
                            </div>
                          )}
                        </div>
                        <span className="font-semibold text-sm text-food-primary ml-2">
                          ${itemTotal.toFixed(2)}
                        </span>
                      </div>
                    </div>
                  );
                })}
              </div>
            </CardContent>
          </Card>

          {step === 1 && (
            <Card className="border-0 shadow-lg">
              <CardHeader className="pb-4">
                <CardTitle className="flex items-center space-x-2">
                  {orderType === "delivery" ? <MapPin className="w-5 h-5 text-food-primary" /> : <MapPinCheck className="w-5 h-5 text-food-primary" />}
                  <span>{orderType === "delivery" ? "Delivery" : "Pickup"} Information</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                {orderType === "delivery" ? (
                  <Form {...deliveryForm}>
                    <form onSubmit={deliveryForm.handleSubmit(onSubmitCustomerInfo, onInvalidCustomerInfo)} className="space-y-4">
                      <FormField
                        control={deliveryForm.control}
                        name="name"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel className="flex items-center space-x-2">
                              <User className="w-4 h-4" />
                              <span>Full Name</span>
                            </FormLabel>
                            <FormControl>
                              <Input placeholder="e.g. Alex Tan" className="border-2 focus:border-food-primary" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <div className="grid grid-cols-1 gap-4">
                        <FormField
                          control={deliveryForm.control}
                          name="email"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel className="flex items-center space-x-2">
                                <Mail className="w-4 h-4" />
                                <span>Email</span>
                              </FormLabel>
                              <FormControl>
                                <Input placeholder="alex@example.com" type="email" className="border-2 focus:border-food-primary" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <FormField
                          control={deliveryForm.control}
                          name="phone"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel className="flex items-center space-x-2">
                                <Phone className="w-4 h-4" />
                                <span>Phone Number</span>
                              </FormLabel>
                              <FormControl>
                                <Input placeholder="+65 9123 4567" type="tel" className="border-2 focus:border-food-primary" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

                      {/* Address Search */}
                      <div className="space-y-3">
                        <Label className="flex items-center space-x-2">
                          <Search className="w-4 h-4" />
                          <span>Find Your Address</span>
                        </Label>
                        <div className="flex space-x-2">
                          <Input
                            placeholder="Search for your delivery address..."
                            value={addressQuery}
                            onChange={(e) => setAddressQuery(e.target.value)}
                            onKeyDown={(e) => {
                              if (e.key === "Enter") {
                                e.preventDefault();
                                handleAddressSearch();
                              }
                            }}
                            className="border-2 focus:border-food-primary"
                          />
                          <Button
                            type="button"
                            onClick={handleAddressSearch}
                            disabled={isSearching}
                            className="px-5 bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white font-medium shadow-md transition-all flex items-center space-x-1.5 shrink-0"
                          >
                            {isSearching ? (
                              <Loader2 className="w-4 h-4 animate-spin" />
                            ) : (
                              <>
                                <Search className="w-4 h-4" />
                                <span className="text-sm">Search</span>
                              </>
                            )}
                          </Button>
                        </div>

                        {searchResults.length > 0 && (
                          <div className="space-y-2 mt-3">
                            <div className="flex items-center space-x-1.5 text-xs font-semibold text-food-primary uppercase tracking-wider">
                              <MapPin className="w-3.5 h-3.5 animate-bounce" />
                              <span>Select your address below:</span>
                            </div>
                            <div className="border-2 border-food-primary/30 shadow-lg rounded-xl divide-y divide-gray-100 max-h-60 overflow-y-auto bg-white">
                              {searchResults.map((result, index) => (
                                <button
                                  key={index}
                                  type="button"
                                  onClick={() => handleSelectAddress(result)}
                                  className="w-full text-left px-3.5 py-3 hover:bg-food-primary/10 text-sm transition-all group flex items-center justify-between space-x-2 cursor-pointer"
                                >
                                  <div className="flex items-start space-x-2 min-w-0">
                                    <MapPin className="w-4 h-4 text-food-primary mt-0.5 shrink-0 group-hover:scale-110 transition-transform" />
                                    <span className="text-gray-800 font-medium group-hover:text-food-primary transition-colors text-xs leading-snug">
                                      {result.displayName}
                                    </span>
                                  </div>
                                  <span className="shrink-0 text-xs font-semibold px-2.5 py-1 rounded-lg bg-food-primary/10 text-food-primary group-hover:bg-food-primary group-hover:text-white transition-colors shadow-sm">
                                    Select
                                  </span>
                                </button>
                              ))}
                            </div>
                          </div>
                        )}

                        {selectedLocation && (
                          <div className={`rounded-xl border p-3 text-sm flex items-start justify-between ${
                            deliveryQuote && !deliveryQuote.withinLimit
                              ? "border-red-200 bg-red-50 text-red-800"
                              : "border-green-200 bg-green-50 text-green-800"
                          }`}>
                            <div className="flex items-start space-x-2">
                              <MapPinCheck className="w-4 h-4 mt-0.5 shrink-0" />
                              <div>
                                <div className="font-medium">Delivery destination</div>
                                <div>{selectedLocation.displayName}</div>
                                {deliveryQuote && (
                                  <div className="mt-1 text-xs">
                                    {deliveryQuote.distanceKm.toFixed(1)} km away
                                    {deliveryQuote.withinLimit
                                      ? deliveryQuote.deliveryFee > 0
                                        ? ` • $${deliveryQuote.deliveryFee.toFixed(2)} fee`
                                        : " • Free delivery"
                                      : ` • Outside ${deliveryQuote.maxDeliveryDistanceKm} km area`}
                                  </div>
                                )}
                              </div>
                            </div>
                            <Button
                              type="button"
                              variant="outline"
                              size="sm"
                              onClick={handleClearAddress}
                              className="text-xs shrink-0 ml-2 h-8 px-2"
                            >
                              Change
                            </Button>
                          </div>
                        )}
                      </div>

                      {selectedLocation && (
                        <>
                          <FormField
                            control={deliveryForm.control}
                            name="street"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel className="flex items-center space-x-2">
                                  <Home className="w-4 h-4" />
                                  <span>Street Address</span>
                                </FormLabel>
                                <FormControl>
                                  <Input placeholder="e.g. 100 Beach Road" className="border-2 focus:border-food-primary" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />

                          <div className="grid grid-cols-3 gap-3">
                            <FormField
                              control={deliveryForm.control}
                              name="city"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel>City</FormLabel>
                                  <FormControl>
                                    <Input placeholder="Singapore" className="border-2 focus:border-food-primary" {...field} />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />

                            <FormField
                              control={deliveryForm.control}
                              name="state"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel>State</FormLabel>
                                  <FormControl>
                                    <Input placeholder="SG" className="border-2 focus:border-food-primary" {...field} />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />

                            <FormField
                              control={deliveryForm.control}
                              name="postalCode"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel>Postal Code</FormLabel>
                                  <FormControl>
                                    <Input placeholder="189702" className="border-2 focus:border-food-primary" {...field} />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                          </div>

                          <FormField
                            control={deliveryForm.control}
                            name="deliveryInstructions"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Delivery Instructions (Optional)</FormLabel>
                                <FormControl>
                                  <Textarea
                                    placeholder="Leave at door, ring bell, etc..."
                                    className="border-2 focus:border-food-primary resize-none"
                                    rows={3}
                                    {...field}
                                  />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </>
                      )}

                      <Button type="submit" className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-3 rounded-xl font-semibold text-lg">
                        Continue to Payment
                      </Button>
                    </form>
                  </Form>
                ) : (
                  <Form {...pickupForm}>
                    <form onSubmit={pickupForm.handleSubmit(onSubmitCustomerInfo)} className="space-y-4">
                      <FormField
                        control={pickupForm.control}
                        name="name"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel className="flex items-center space-x-2">
                              <User className="w-4 h-4" />
                              <span>Full Name</span>
                            </FormLabel>
                            <FormControl>
                              <Input placeholder="e.g. Alex Tan" className="border-2 focus:border-food-primary" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <FormField
                        control={pickupForm.control}
                        name="email"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel className="flex items-center space-x-2">
                              <Mail className="w-4 h-4" />
                              <span>Email</span>
                            </FormLabel>
                            <FormControl>
                              <Input placeholder="alex@example.com" type="email" className="border-2 focus:border-food-primary" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <FormField
                        control={pickupForm.control}
                        name="phone"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel className="flex items-center space-x-2">
                              <Phone className="w-4 h-4" />
                              <span>Phone Number</span>
                            </FormLabel>
                            <FormControl>
                              <Input placeholder="+65 9123 4567" type="tel" className="border-2 focus:border-food-primary" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <div className="bg-gray-50 p-4 rounded-xl">
                        <h4 className="font-semibold text-gray-900 mb-2 flex items-center">
                          <MapPinCheck className="w-4 h-4 mr-1.5 text-food-primary" />
                          Pickup Location
                        </h4>
                        {restaurant?.name && (
                          <p className="font-semibold text-gray-900 text-sm mb-1">{restaurant.name}</p>
                        )}
                        <p className="text-gray-600 text-sm">{restaurant?.address || "Address unavailable"}</p>
                        <p className="text-gray-600 text-sm mt-1 flex items-center">
                          <Clock className="w-3.5 h-3.5 mr-1 text-food-primary" />
                          Estimated pickup time: {getEstimatedTime()}
                        </p>
                      </div>

                      <Button type="submit" className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-3 rounded-xl font-semibold text-lg">
                        Continue to Payment
                      </Button>
                    </form>
                  </Form>
                )}
              </CardContent>
            </Card>
          )}

          {step === 2 && (
            <Card className="border-0 shadow-lg">
              <CardHeader className="pb-4">
                <div className="flex items-center justify-between">
                  <CardTitle className="flex items-center space-x-2">
                    <CreditCard className="w-5 h-5 text-food-primary" />
                    <span>Payment Method</span>
                  </CardTitle>
                  <Button variant="ghost" size="sm" onClick={handleBackStep} className="text-food-primary">
                    <ArrowLeft className="w-4 h-4 mr-1" />
                    Back
                  </Button>
                </div>
              </CardHeader>
              <CardContent>
                {paymentMethod === "creditCard" && currentOrder ? (
                  clientSecret ? (
                    <Elements
                      stripe={stripePromise}
                      options={{
                        clientSecret,
                        appearance: {
                          theme: 'stripe',
                          variables: {
                            colorPrimary: '#f97316',
                          },
                        },
                      }}
                    >
                      <StripePaymentForm
                        orderId={currentOrder.id}
                        total={currentOrder.total}
                        onSuccess={handlePaymentSuccess}
                        onError={handlePaymentError}
                      />
                    </Elements>
                  ) : (
                    <div className="space-y-4">
                      <div className="bg-yellow-50 p-4 rounded-xl border border-yellow-200">
                        <h3 className="text-base font-semibold mb-2 text-yellow-900">Stripe payment isn't ready</h3>
                        <p className="text-yellow-800 text-sm">
                          We couldn't initialize a Stripe PaymentIntent for this order.
                        </p>
                        <p className="text-yellow-800 mt-2 text-sm">
                          In dev, this usually means Stripe keys are not configured for the backend and/or frontend.
                        </p>
                      </div>

                      <Button
                        onClick={retryCreatePaymentIntent}
                        disabled={isSubmitting}
                        className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-3 rounded-xl font-semibold text-lg disabled:opacity-50"
                      >
                        {isSubmitting ? (
                          <span className="flex items-center justify-center">
                            <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                            Initializing Stripe...
                          </span>
                        ) : (
                          "Retry Stripe Payment"
                        )}
                      </Button>
                    </div>
                  )
                ) : (
                  <div className="space-y-6">
                    {/* Payment Method Selection */}
                    <div className="space-y-3">
                      <Label className="text-base font-semibold">Choose Payment Method</Label>
                      <div className="space-y-3">
                        <div
                          className={`p-4 border-2 rounded-xl cursor-pointer transition-all ${
                            paymentMethod === "creditCard" ? 'border-food-primary bg-food-primary/5' : 'border-gray-200 hover:border-gray-300'
                          }`}
                          onClick={() => setPaymentMethod("creditCard")}
                        >
                          <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-3">
                              <CreditCard className="w-5 h-5 text-food-primary" />
                              <span className="font-medium">Credit Card</span>
                            </div>
                            <div className={`w-4 h-4 rounded-full border-2 ${
                              paymentMethod === "creditCard" ? 'border-food-primary bg-food-primary' : 'border-gray-300'
                            }`}>
                              {paymentMethod === "creditCard" && <div className="w-2 h-2 bg-white rounded-full m-0.5" />}
                            </div>
                          </div>
                        </div>

                        <div
                          className={`p-4 border-2 rounded-xl cursor-pointer transition-all ${
                            paymentMethod === "payAtLocation" ? 'border-food-primary bg-food-primary/5' : 'border-gray-200 hover:border-gray-300'
                          }`}
                          onClick={() => setPaymentMethod("payAtLocation")}
                        >
                          <div className="flex items-center justify-between">
                            <div className="flex items-center space-x-3">
                              <MapPin className="w-5 h-5 text-food-primary" />
                              <span className="font-medium">Pay at {orderType === "delivery" ? "Door" : "Pickup"}</span>
                            </div>
                            <div className={`w-4 h-4 rounded-full border-2 ${
                              paymentMethod === "payAtLocation" ? 'border-food-primary bg-food-primary' : 'border-gray-300'
                            }`}>
                              {paymentMethod === "payAtLocation" && <div className="w-2 h-2 bg-white rounded-full m-0.5" />}
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>

                    {paymentMethod === "payAtLocation" && (
                      <div className="bg-gray-50 p-4 rounded-xl text-center">
                        <MapPin className="w-10 h-10 text-food-primary mx-auto mb-3" />
                        <h3 className="text-base font-semibold mb-1">Pay at {orderType === "delivery" ? "Delivery" : "Pickup"}</h3>
                        <p className="text-gray-600 text-sm">
                          You'll pay ${total.toFixed(2)} when you receive your order.
                        </p>
                      </div>
                    )}

                    <Button
                      onClick={onSubmitPayment}
                      disabled={isSubmitting}
                      className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-3 rounded-xl font-semibold text-lg disabled:opacity-50"
                    >
                      {isSubmitting ? (
                        <span className="flex items-center justify-center">
                          <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                          Processing...
                        </span>
                      ) : (
                        `Place Order • $${total.toFixed(2)}`
                      )}
                    </Button>
                  </div>
                )}
              </CardContent>
            </Card>
          )}
        </div>

        {/* Order Summary - Sticky Bottom */}
        <div className="fixed bottom-0 left-1/2 transform -translate-x-1/2 w-full max-w-md bg-white border-t shadow-lg">
          <div className="px-6 py-4">
            <div className="space-y-2">
              <div className="flex justify-between text-sm">
                <span>Subtotal</span>
                <span>${subtotal.toFixed(2)}</span>
              </div>
              {discount > 0 && (
                <div className="flex justify-between text-sm text-green-600 font-medium">
                  <span>Discount ({appliedPromo?.code})</span>
                  <span>-${discount.toFixed(2)}</span>
                </div>
              )}
              {orderType === "delivery" && (
                <div className="flex justify-between text-sm">
                  <span>Delivery Fee</span>
                  {deliveryFee > 0 ? (
                    <span>${deliveryFee.toFixed(2)}</span>
                  ) : (
                    <span className="text-gray-500">Select your address</span>
                  )}
                </div>
              )}
              {tax > 0 && (
                <div className="flex justify-between text-sm">
                  <span>Tax ({(taxRate * 100).toFixed(0)}%)</span>
                  <span>${tax.toFixed(2)}</span>
                </div>
              )}
              <div className="flex justify-between font-bold text-lg pt-2 border-t">
                <span>Total</span>
                <span className="text-food-primary">${total.toFixed(2)}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default CheckoutMobile;
