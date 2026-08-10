import React, { useState } from "react";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { MapPin, Clock, CreditCard, CheckCircle, User, Phone, Mail, Home, MapPinCheck, Package, Loader2, Search } from "lucide-react";
import { loadStripe } from "@stripe/stripe-js";
import { Elements } from "@stripe/react-stripe-js";

import Layout from "@/components/Layout";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Link } from "react-router-dom";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useCart } from "@/context/CartContext";
import { Label } from "@/components/ui/label";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { orderService, Order, DeliveryQuote } from "@/services/orderService";
import { searchAddress, GeocodingResult } from "@/lib/geocoding";
import StripePaymentForm from "@/components/StripePaymentForm";
import { useRestaurantDetails } from "@/hooks/useRestaurantDetails";

// Initialize Stripe (may be null when key isn't configured)
const stripePromise = import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY
  ? loadStripe(import.meta.env.VITE_STRIPE_PUBLISHABLE_KEY)
  : Promise.resolve(null);

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

const paymentFormSchema = z.object({
  paymentMethod: z.enum(["creditCard", "payAtLocation"]),
  cardNumber: z.string().optional(),
  cardExpiry: z.string().optional(),
  cardCvc: z.string().optional(),
});

const DEFAULT_DELIVERY_TIME = { min: 30, max: 45 };
const DEFAULT_PICKUP_TIME = { min: 15, max: 25 };
const DEFAULT_RESTAURANT_ADDRESS = "123 Main Street, City, State 12345";

const CheckoutDesktop: React.FC = () => {
  const navigate = useNavigate();
  const { items, getTotalPrice, clearCart, orderType, restaurantId } = useCart();
  const { data: restaurant } = useRestaurantDetails(restaurantId || "");
  const [step, setStep] = useState(1);
  const [paymentMethod, setPaymentMethod] = useState<"creditCard" | "payAtLocation">("creditCard");
  const [orderDetails, setOrderDetails] = useState<Record<string, unknown>>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // Stripe-related state
  const [currentOrder, setCurrentOrder] = useState<Order | null>(null);
  const [clientSecret, setClientSecret] = useState<string | null>(null);
  const [isCreatingOrder, setIsCreatingOrder] = useState(false);

  // Delivery destination state
  const [addressQuery, setAddressQuery] = useState("");
  const [searchResults, setSearchResults] = useState<GeocodingResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [selectedLocation, setSelectedLocation] = useState<GeocodingResult | null>(null);
  const [deliveryQuote, setDeliveryQuote] = useState<DeliveryQuote | null>(null);

  const subtotal = getTotalPrice();
  const deliveryFee = orderType === "delivery" && deliveryQuote?.withinLimit ? deliveryQuote.deliveryFee : 0;
  const taxRate = restaurant?.taxRate ?? 0.10;
  const tax = taxRate > 0 ? subtotal * taxRate : 0;
  const total = subtotal + deliveryFee + tax;

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

  const paymentForm = useForm<z.infer<typeof paymentFormSchema>>({
    resolver: zodResolver(paymentFormSchema),
    defaultValues: {
      paymentMethod: "creditCard",
      cardNumber: "",
      cardExpiry: "",
      cardCvc: "",
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
  const onSubmitCustomerInfo = async (data: z.infer<typeof deliveryFormSchema> | z.infer<typeof pickupFormSchema>) => {
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
    setIsCreatingOrder(true);

    try {
      // Create the order
      const orderData = data as z.infer<typeof deliveryFormSchema>;
      const order = await orderService.createOrder({
        restaurantId: restaurantId || "",
        customerName: orderData.name,
        customerEmail: orderData.email,
        customerPhone: orderData.phone,
        orderType: orderType as "delivery" | "pickup",
        paymentMethod: paymentMethod === "creditCard" ? "creditCard" : "cash",
        items: items.map((item) => ({
          menuItemId: item.menuItem.id,
          quantity: item.quantity,
          specialInstructions: item.specialInstructions,
          addons: item.selectedAddons?.map(({ addon, quantity: addonQty }) => ({
            addonId: addon.id,
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
      });

      setCurrentOrder(order);

      if (paymentMethod === "creditCard") {
        // Create payment intent for credit card payment
        const paymentIntent = await orderService.createPaymentIntent(order.id);
        setClientSecret(paymentIntent.clientSecret);
      }

      setStep(2);
    } catch (error) {
      console.error("Failed to create order:", error);
      const message = error instanceof Error ? error.message : "Failed to create order. Please try again.";
      toast.error(message);
    } finally {
      setIsCreatingOrder(false);
    }
  };

  const onSubmitPayAtLocation = async () => {
    setIsSubmitting(true);

    try {
      // For pay at location, confirm the order directly
      if (currentOrder) {
        await orderService.confirmPayment(currentOrder.id);
      }
      clearCart();
      toast.success("Order placed successfully!");
      navigate(`/order-confirmation/${currentOrder?.id}`);
    } catch (error) {
      console.error("Failed to place order:", error);
      toast.error("Failed to place order. Please try again.");
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
      setIsCreatingOrder(true);
      const paymentIntent = await orderService.createPaymentIntent(currentOrder.id);
      setClientSecret(paymentIntent.clientSecret);
    } catch (error) {
      console.error("Failed to create payment intent:", error);
      toast.error("Failed to initialize Stripe payment. Please try again.");
    } finally {
      setIsCreatingOrder(false);
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
    <Layout>
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="mb-8 flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Checkout</h1>
              <p className="text-gray-600 mt-2">Complete your {orderType === "delivery" ? "delivery" : "pickup"} order</p>
            </div>
            <Link to={restaurantId ? `/restaurant/${restaurantId}` : "/"}>
              <Button variant="outline" className="text-food-primary border-food-primary hover:bg-food-primary/10">
                Cancel
              </Button>
            </Link>
          </div>

          {/* Progress Steps */}
          <div className="mb-8">
            <div className="flex items-center justify-center space-x-8 max-w-2xl mx-auto">
              <div className="flex items-center">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center text-lg font-semibold ${
                  step >= 1 ? 'bg-food-primary text-white' : 'bg-gray-200 text-gray-600'
                }`}>
                  {step > 1 ? <CheckCircle className="w-6 h-6" /> : '1'}
                </div>
                <span className={`ml-3 text-base font-medium ${step === 1 ? 'text-food-primary' : 'text-gray-600'}`}>
                  Customer Information
                </span>
              </div>
              <div className="flex-1 h-2 bg-gray-200 rounded-full mx-4">
                <div className={`h-full rounded-full transition-all duration-300 ${
                  step >= 2 ? 'bg-food-primary' : 'bg-gray-200'
                }`} style={{ width: step >= 2 ? '100%' : '0%' }} />
              </div>
              <div className="flex items-center">
                <div className={`w-12 h-12 rounded-full flex items-center justify-center text-lg font-semibold ${
                  step >= 2 ? 'bg-food-primary text-white' : 'bg-gray-200 text-gray-600'
                }`}>
                  {step > 2 ? <CheckCircle className="w-6 h-6" /> : '2'}
                </div>
                <span className={`ml-3 text-base font-medium ${step === 2 ? 'text-food-primary' : 'text-gray-600'}`}>
                  Payment
                </span>
              </div>
            </div>
          </div>

          {/* Two Column Layout */}
          <div className="grid grid-cols-3 gap-8">
            {/* Left Column - Forms */}
            <div className="col-span-2">
              {step === 1 && (
                <Card className="shadow-md border">
                  <CardHeader className="bg-gray-50">
                    <CardTitle className="flex items-center space-x-2 text-xl">
                      {orderType === "delivery" ? <MapPin className="w-6 h-6 text-food-primary" /> : <MapPinCheck className="w-6 h-6 text-food-primary" />}
                      <span>{orderType === "delivery" ? "Delivery" : "Pickup"} Information</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="p-6">
                    {orderType === "delivery" ? (
                      <Form {...deliveryForm}>
                        <form onSubmit={deliveryForm.handleSubmit(onSubmitCustomerInfo, onInvalidCustomerInfo)} className="space-y-6">
                          <FormField
                            control={deliveryForm.control}
                            name="name"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel className="flex items-center space-x-2 text-base">
                                  <User className="w-5 h-5" />
                                  <span>Full Name</span>
                                </FormLabel>
                                <FormControl>
                                  <Input placeholder="John Doe" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />

                          <div className="grid grid-cols-2 gap-6">
                            <FormField
                              control={deliveryForm.control}
                              name="email"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel className="flex items-center space-x-2 text-base">
                                    <Mail className="w-5 h-5" />
                                    <span>Email</span>
                                  </FormLabel>
                                  <FormControl>
                                    <Input placeholder="john@example.com" type="email" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                  <FormLabel className="flex items-center space-x-2 text-base">
                                    <Phone className="w-5 h-5" />
                                    <span>Phone Number</span>
                                  </FormLabel>
                                  <FormControl>
                                    <Input placeholder="+1 (555) 123-4567" type="tel" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                          </div>

                          <div className="space-y-3">
                            <Label className="flex items-center space-x-2 text-base">
                              <Search className="w-5 h-5" />
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
                                className="border-2 focus:border-food-primary h-12 text-base"
                              />
                              <Button
                                type="button"
                                onClick={handleAddressSearch}
                                disabled={isSearching}
                                variant="outline"
                                className="h-12 px-6 text-food-primary border-food-primary hover:bg-food-primary/10"
                              >
                                {isSearching ? <Loader2 className="w-5 h-5 animate-spin" /> : <Search className="w-5 h-5" />}
                              </Button>
                            </div>

                            {searchResults.length > 0 && (
                              <div className="border-2 rounded-xl divide-y max-h-48 overflow-y-auto">
                                {searchResults.map((result, index) => (
                                  <button
                                    key={index}
                                    type="button"
                                    onClick={() => handleSelectAddress(result)}
                                    className="w-full text-left px-4 py-3 hover:bg-food-primary/5 text-sm text-gray-700 transition-colors"
                                  >
                                    {result.displayName}
                                  </button>
                                ))}
                              </div>
                            )}

                            {selectedLocation && (
                              <div className={`flex items-start justify-between rounded-xl border p-4 text-sm ${
                                deliveryQuote && !deliveryQuote.withinLimit
                                  ? "border-red-200 bg-red-50 text-red-800"
                                  : "border-green-200 bg-green-50 text-green-800"
                              }`}>
                                <div className="flex items-start space-x-2">
                                  <MapPin className="w-4 h-4 mt-0.5 shrink-0" />
                                  <div>
                                    <div className="font-medium">Delivery destination</div>
                                    <div>{selectedLocation.displayName}</div>
                                    {deliveryQuote && (
                                      <div className="mt-1">
                                        {deliveryQuote.distanceKm.toFixed(1)} km away
                                        {deliveryQuote.withinLimit
                                          ? deliveryQuote.deliveryFee > 0
                                            ? ` • $${deliveryQuote.deliveryFee.toFixed(2)} delivery fee`
                                            : " • Free delivery"
                                          : ` • Outside the ${deliveryQuote.maxDeliveryDistanceKm} km delivery area`}
                                      </div>
                                    )}
                                  </div>
                                </div>
                                <Button
                                  type="button"
                                  variant="outline"
                                  size="sm"
                                  onClick={handleClearAddress}
                                  className="text-xs shrink-0 ml-3"
                                >
                                  Change Address
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
                                    <FormLabel className="flex items-center space-x-2 text-base">
                                      <Home className="w-5 h-5" />
                                      <span>Street Address</span>
                                    </FormLabel>
                                    <FormControl>
                                      <Input placeholder="123 Main Street" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
                                    </FormControl>
                                    <FormMessage />
                                  </FormItem>
                                )}
                              />

                              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                                <FormField
                                  control={deliveryForm.control}
                                  name="city"
                                  render={({ field }) => (
                                    <FormItem>
                                      <FormLabel className="text-base">City</FormLabel>
                                      <FormControl>
                                        <Input placeholder="Singapore" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                      <FormLabel className="text-base">State (Optional)</FormLabel>
                                      <FormControl>
                                        <Input placeholder="SG" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                      <FormLabel className="text-base">Postal Code</FormLabel>
                                      <FormControl>
                                        <Input placeholder="e.g. 048582" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                    <FormLabel className="text-base">Delivery Instructions (Optional)</FormLabel>
                                    <FormControl>
                                      <Textarea 
                                        placeholder="Leave at door, ring bell, etc..."
                                        className="border-2 focus:border-food-primary resize-none text-base"
                                        rows={4}
                                        {...field} 
                                      />
                                    </FormControl>
                                    <FormMessage />
                                  </FormItem>
                                )}
                              />
                            </>
                          )}

                          <Button 
                            type="submit" 
                            disabled={isCreatingOrder}
                            className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg disabled:opacity-50"
                          >
                            {isCreatingOrder ? (
                              <span className="flex items-center justify-center">
                                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                                Creating Order...
                              </span>
                            ) : (
                              "Continue to Payment"
                            )}
                          </Button>
                        </form>
                      </Form>
                    ) : (
                      <Form {...pickupForm}>
                        <form onSubmit={pickupForm.handleSubmit(onSubmitCustomerInfo)} className="space-y-6">
                          <FormField
                            control={pickupForm.control}
                            name="name"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel className="flex items-center space-x-2 text-base">
                                  <User className="w-5 h-5" />
                                  <span>Full Name</span>
                                </FormLabel>
                                <FormControl>
                                  <Input placeholder="John Doe" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />

                          <div className="grid grid-cols-2 gap-6">
                            <FormField
                              control={pickupForm.control}
                              name="email"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel className="flex items-center space-x-2 text-base">
                                    <Mail className="w-5 h-5" />
                                    <span>Email</span>
                                  </FormLabel>
                                  <FormControl>
                                    <Input placeholder="john@example.com" type="email" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                  <FormLabel className="flex items-center space-x-2 text-base">
                                    <Phone className="w-5 h-5" />
                                    <span>Phone Number</span>
                                  </FormLabel>
                                  <FormControl>
                                    <Input placeholder="+1 (555) 123-4567" type="tel" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                          </div>

                          <div className="bg-gray-50 p-6 rounded-xl border">
                            <h4 className="font-semibold text-gray-900 mb-3 text-lg flex items-center">
                              <MapPinCheck className="w-5 h-5 mr-2 text-food-primary" />
                              Pickup Location
                            </h4>
                            <p className="text-gray-600">{DEFAULT_RESTAURANT_ADDRESS}</p>
                            <p className="text-gray-600 mt-2 flex items-center">
                              <Clock className="w-4 h-4 mr-2 text-food-primary" />
                              Estimated pickup time: {getEstimatedTime()}
                            </p>
                          </div>

                          <Button 
                            type="submit" 
                            disabled={isCreatingOrder}
                            className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg disabled:opacity-50"
                          >
                            {isCreatingOrder ? (
                              <span className="flex items-center justify-center">
                                <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                                Creating Order...
                              </span>
                            ) : (
                              "Continue to Payment"
                            )}
                          </Button>
                        </form>
                      </Form>
                    )}
                  </CardContent>
                </Card>
              )}

              {step === 2 && (
                <Card className="shadow-md border">
                  <CardHeader className="bg-gray-50">
                    <div className="flex items-center justify-between">
                      <CardTitle className="flex items-center space-x-2 text-xl">
                        <CreditCard className="w-6 h-6 text-food-primary" />
                        <span>Payment</span>
                      </CardTitle>
                      <Button variant="outline" size="sm" onClick={handleBackStep} className="text-food-primary border-food-primary hover:bg-food-primary/10">
                        ← Back
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent className="p-6">
                    {paymentMethod === "creditCard" ? (
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
                            orderId={currentOrder?.id || ""}
                            total={currentOrder?.total || total}
                            onSuccess={handlePaymentSuccess}
                            onError={handlePaymentError}
                          />
                        </Elements>
                      ) : (
                        <div className="space-y-4">
                          <div className="bg-yellow-50 p-6 rounded-xl border border-yellow-200">
                            <h3 className="text-lg font-semibold mb-2 text-yellow-900">Stripe payment isn't ready</h3>
                            <p className="text-yellow-800">
                              We couldn't initialize a Stripe PaymentIntent for this order.
                            </p>
                            <p className="text-yellow-800 mt-2 text-sm">
                              In dev, this usually means Stripe keys are not configured for the backend and/or frontend.
                            </p>
                          </div>

                          <Button
                            onClick={retryCreatePaymentIntent}
                            disabled={isCreatingOrder}
                            className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg disabled:opacity-50"
                          >
                            {isCreatingOrder ? (
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
                        <div className="bg-gray-50 p-6 rounded-xl border text-center">
                          <MapPin className="w-12 h-12 text-food-primary mx-auto mb-4" />
                          <h3 className="text-lg font-semibold mb-2">Pay at {orderType === "delivery" ? "Delivery" : "Pickup"}</h3>
                          <p className="text-gray-600 mb-4">
                            You'll pay ${(currentOrder?.total || total).toFixed(2)} when you receive your order.
                          </p>
                        </div>
                        <Button
                          onClick={onSubmitPayAtLocation}
                          disabled={isSubmitting}
                          className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg disabled:opacity-50"
                        >
                          {isSubmitting ? (
                            <span className="flex items-center justify-center">
                              <Loader2 className="w-5 h-5 mr-2 animate-spin" />
                              Placing Order...
                            </span>
                          ) : (
                            `Place Order • $${(currentOrder?.total || total).toFixed(2)}`
                          )}
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              )}
            </div>

            {/* Right Column - Order Summary (Sticky) */}
            <div className="col-span-1">
              <div className="sticky top-24">
                <Card className="shadow-md border">
                  <CardHeader className="bg-gray-50">
                    <CardTitle className="flex items-center space-x-2">
                      <Package className="w-5 h-5 text-food-primary" />
                      <span>Order Summary</span>
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="p-6">
                    <div className="space-y-4">
                      <div>
                        <div className="text-sm text-gray-600 mb-2">Order Type</div>
                        <div className="font-semibold text-base capitalize">{orderType}</div>
                      </div>
                      
                      <div className="border-t pt-4 space-y-3">
                        <div className="flex justify-between">
                          <span className="text-gray-600">Subtotal</span>
                          <span className="font-medium">${subtotal.toFixed(2)}</span>
                        </div>
                        {orderType === "delivery" && (
                          <div className="flex justify-between">
                            <span className="text-gray-600">Delivery Fee</span>
                            <span className="font-medium">
                              {deliveryQuote && deliveryQuote.withinLimit
                                ? deliveryQuote.deliveryFee > 0
                                  ? `$${deliveryQuote.deliveryFee.toFixed(2)}`
                                  : "Free"
                                : "Select address"}
                            </span>
                          </div>
                        )}
                        {tax > 0 && (
                          <div className="flex justify-between">
                            <span className="text-gray-600">Tax ({(taxRate * 100).toFixed(0)}%)</span>
                            <span className="font-medium">${tax.toFixed(2)}</span>
                          </div>
                        )}
                      </div>

                      <div className="border-t pt-4">
                        <div className="flex justify-between items-center">
                          <span className="font-bold text-lg">Total</span>
                          <span className="font-bold text-2xl text-food-primary">${total.toFixed(2)}</span>
                        </div>
                      </div>

                      <div className="border-t pt-4">
                        <div className="flex items-center text-sm text-gray-600">
                          <Clock className="w-4 h-4 mr-2 text-food-primary" />
                          <span>Estimated time: {getEstimatedTime()}</span>
                        </div>
                      </div>
                    </div>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default CheckoutDesktop;
