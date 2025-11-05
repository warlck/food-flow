import React, { useState } from "react";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { MapPin, Clock, CreditCard, ArrowLeft, CheckCircle, User, Phone, Mail, Home, MapPinCheck } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useCart } from "@/context/CartContext";
import { Label } from "@/components/ui/label";
import { mockRestaurant } from "@/data/mockData";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

// Form schemas
const deliveryFormSchema = z.object({
  name: z.string().min(2, { message: "Name must be at least 2 characters" }),
  email: z.string().email({ message: "Invalid email address" }),
  phone: z.string().min(10, { message: "Phone number must be at least 10 characters" }),
  street: z.string().min(5, { message: "Street address must be at least 5 characters" }),
  city: z.string().min(2, { message: "City must be at least 2 characters" }),
  state: z.string().min(2, { message: "State must be at least 2 characters" }),
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

const CheckoutMobile: React.FC = () => {
  const navigate = useNavigate();
  const { items, getTotalPrice, clearCart, orderType } = useCart();
  const [step, setStep] = useState(1);
  const [paymentMethod, setPaymentMethod] = useState<"creditCard" | "payAtLocation">("creditCard");
  const [orderDetails, setOrderDetails] = useState<any>({});
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  // Calculate totals
  const subtotal = getTotalPrice();
  const deliveryFee = orderType === "delivery" ? mockRestaurant.deliveryFee : 0;
  const tax = subtotal * 0.1; // 10% tax
  const total = subtotal + deliveryFee + tax;

  // Setup forms based on order type
  const deliveryForm = useForm<z.infer<typeof deliveryFormSchema>>({
    resolver: zodResolver(deliveryFormSchema),
    defaultValues: {
      name: "",
      email: "",
      phone: "",
      street: "",
      city: "",
      state: "",
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

  // Submit handlers
  const onSubmitCustomerInfo = (data: any) => {
    setOrderDetails({ ...orderDetails, ...data });
    setStep(2);
  };

  const onSubmitPayment = async (data: any) => {
    setIsSubmitting(true);
    setOrderDetails({ ...orderDetails, ...data });
    
    // Simulate order submission
    setTimeout(() => {
      clearCart();
      toast.success("Order placed successfully!");
      navigate("/mobile-restaurant/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d");
      setIsSubmitting(false);
    }, 2000);
  };

  const handleBackStep = () => {
    setStep(step - 1);
  };

  const getEstimatedTime = () => {
    if (orderType === "delivery") {
      return `${mockRestaurant.estimatedDeliveryTime.min}-${mockRestaurant.estimatedDeliveryTime.max} minutes`;
    } else {
      return `${mockRestaurant.estimatedPickupTime?.min || 15}-${mockRestaurant.estimatedPickupTime?.max || 25} minutes`;
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
        <div className="px-6 py-6 pb-32">
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
                    <form onSubmit={deliveryForm.handleSubmit(onSubmitCustomerInfo)} className="space-y-4">
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
                              <Input placeholder="John Doe" className="border-2 focus:border-food-primary" {...field} />
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
                                <Input placeholder="john@example.com" type="email" className="border-2 focus:border-food-primary" {...field} />
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
                                <Input placeholder="+1 (555) 123-4567" type="tel" className="border-2 focus:border-food-primary" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

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
                              <Input placeholder="123 Main Street" className="border-2 focus:border-food-primary" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <div className="grid grid-cols-2 gap-4">
                        <FormField
                          control={deliveryForm.control}
                          name="city"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>City</FormLabel>
                              <FormControl>
                                <Input placeholder="New York" className="border-2 focus:border-food-primary" {...field} />
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
                                <Input placeholder="NY" className="border-2 focus:border-food-primary" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />
                      </div>

                      <FormField
                        control={deliveryForm.control}
                        name="postalCode"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Postal Code</FormLabel>
                            <FormControl>
                              <Input placeholder="10001" className="border-2 focus:border-food-primary" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

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
                              <Input placeholder="John Doe" className="border-2 focus:border-food-primary" {...field} />
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
                              <Input placeholder="john@example.com" type="email" className="border-2 focus:border-food-primary" {...field} />
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
                              <Input placeholder="+1 (555) 123-4567" type="tel" className="border-2 focus:border-food-primary" {...field} />
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      <div className="bg-gray-50 p-4 rounded-xl">
                        <h4 className="font-semibold text-gray-900 mb-2">Pickup Location</h4>
                        <p className="text-gray-600 text-sm">{mockRestaurant.address}</p>
                        <p className="text-gray-600 text-sm mt-1">Estimated pickup time: {getEstimatedTime()}</p>
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
                <Form {...paymentForm}>
                  <form onSubmit={paymentForm.handleSubmit(onSubmitPayment)} className="space-y-6">
                    
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

                    {/* Credit Card Fields */}
                    {paymentMethod === "creditCard" && (
                      <div className="space-y-4">
                        <FormField
                          control={paymentForm.control}
                          name="cardNumber"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Card Number</FormLabel>
                              <FormControl>
                                <Input placeholder="1234 5678 9012 3456" className="border-2 focus:border-food-primary" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <div className="grid grid-cols-2 gap-4">
                          <FormField
                            control={paymentForm.control}
                            name="cardExpiry"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Expiry Date</FormLabel>
                                <FormControl>
                                  <Input placeholder="MM/YY" className="border-2 focus:border-food-primary" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />

                          <FormField
                            control={paymentForm.control}
                            name="cardCvc"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>CVC</FormLabel>
                                <FormControl>
                                  <Input placeholder="123" className="border-2 focus:border-food-primary" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </div>
                      </div>
                    )}

                    <Button 
                      type="submit" 
                      disabled={isSubmitting}
                      className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-3 rounded-xl font-semibold text-lg disabled:opacity-50"
                    >
                      {isSubmitting ? "Processing..." : `Place Order • $${total.toFixed(2)}`}
                    </Button>
                  </form>
                </Form>
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
              {deliveryFee > 0 && (
                <div className="flex justify-between text-sm">
                  <span>Delivery Fee</span>
                  <span>${deliveryFee.toFixed(2)}</span>
                </div>
              )}
              <div className="flex justify-between text-sm">
                <span>Tax</span>
                <span>${tax.toFixed(2)}</span>
              </div>
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
