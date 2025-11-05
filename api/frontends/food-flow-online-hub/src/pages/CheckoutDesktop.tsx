import React, { useState } from "react";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { MapPin, Clock, CreditCard, CheckCircle, User, Phone, Mail, Home, MapPinCheck, Package } from "lucide-react";

import Layout from "@/components/Layout";
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

const CheckoutDesktop: React.FC = () => {
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
      navigate("/restaurant/a1b2c3d4-e5f6-4a5b-8c9d-1e2f3a4b5c6d");
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
    <Layout>
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-7xl mx-auto">
          {/* Header */}
          <div className="mb-8">
            <h1 className="text-3xl font-bold text-gray-900">Checkout</h1>
            <p className="text-gray-600 mt-2">Complete your {orderType === "delivery" ? "delivery" : "pickup"} order</p>
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
                        <form onSubmit={deliveryForm.handleSubmit(onSubmitCustomerInfo)} className="space-y-6">
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

                          <div className="grid grid-cols-3 gap-6">
                            <FormField
                              control={deliveryForm.control}
                              name="city"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel className="text-base">City</FormLabel>
                                  <FormControl>
                                    <Input placeholder="New York" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                  <FormLabel className="text-base">State</FormLabel>
                                  <FormControl>
                                    <Input placeholder="NY" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                    <Input placeholder="10001" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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

                          <Button type="submit" className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg">
                            Continue to Payment
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
                            <p className="text-gray-600">{mockRestaurant.address}</p>
                            <p className="text-gray-600 mt-2 flex items-center">
                              <Clock className="w-4 h-4 mr-2 text-food-primary" />
                              Estimated pickup time: {getEstimatedTime()}
                            </p>
                          </div>

                          <Button type="submit" className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg">
                            Continue to Payment
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
                        <span>Payment Method</span>
                      </CardTitle>
                      <Button variant="outline" size="sm" onClick={handleBackStep} className="text-food-primary border-food-primary hover:bg-food-primary/10">
                        ← Back
                      </Button>
                    </div>
                  </CardHeader>
                  <CardContent className="p-6">
                    <Form {...paymentForm}>
                      <form onSubmit={paymentForm.handleSubmit(onSubmitPayment)} className="space-y-8">
                        
                        {/* Payment Method Selection */}
                        <div className="space-y-4">
                          <Label className="text-lg font-semibold">Choose Payment Method</Label>
                          <div className="grid grid-cols-2 gap-4">
                            <div 
                              className={`p-6 border-2 rounded-xl cursor-pointer transition-all ${
                                paymentMethod === "creditCard" ? 'border-food-primary bg-food-primary/5 shadow-md' : 'border-gray-200 hover:border-gray-300'
                              }`}
                              onClick={() => setPaymentMethod("creditCard")}
                            >
                              <div className="flex flex-col items-center space-y-3">
                                <CreditCard className="w-8 h-8 text-food-primary" />
                                <span className="font-medium text-base">Credit Card</span>
                                <div className={`w-5 h-5 rounded-full border-2 ${
                                  paymentMethod === "creditCard" ? 'border-food-primary bg-food-primary' : 'border-gray-300'
                                }`}>
                                  {paymentMethod === "creditCard" && <div className="w-3 h-3 bg-white rounded-full m-0.5" />}
                                </div>
                              </div>
                            </div>

                            <div 
                              className={`p-6 border-2 rounded-xl cursor-pointer transition-all ${
                                paymentMethod === "payAtLocation" ? 'border-food-primary bg-food-primary/5 shadow-md' : 'border-gray-200 hover:border-gray-300'
                              }`}
                              onClick={() => setPaymentMethod("payAtLocation")}
                            >
                              <div className="flex flex-col items-center space-y-3">
                                <MapPin className="w-8 h-8 text-food-primary" />
                                <span className="font-medium text-base">Pay at {orderType === "delivery" ? "Door" : "Pickup"}</span>
                                <div className={`w-5 h-5 rounded-full border-2 ${
                                  paymentMethod === "payAtLocation" ? 'border-food-primary bg-food-primary' : 'border-gray-300'
                                }`}>
                                  {paymentMethod === "payAtLocation" && <div className="w-3 h-3 bg-white rounded-full m-0.5" />}
                                </div>
                              </div>
                            </div>
                          </div>
                        </div>

                        {/* Credit Card Fields */}
                        {paymentMethod === "creditCard" && (
                          <div className="space-y-6">
                            <FormField
                              control={paymentForm.control}
                              name="cardNumber"
                              render={({ field }) => (
                                <FormItem>
                                  <FormLabel className="text-base">Card Number</FormLabel>
                                  <FormControl>
                                    <Input placeholder="1234 5678 9012 3456" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />

                            <div className="grid grid-cols-2 gap-6">
                              <FormField
                                control={paymentForm.control}
                                name="cardExpiry"
                                render={({ field }) => (
                                  <FormItem>
                                    <FormLabel className="text-base">Expiry Date</FormLabel>
                                    <FormControl>
                                      <Input placeholder="MM/YY" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                                    <FormLabel className="text-base">CVC</FormLabel>
                                    <FormControl>
                                      <Input placeholder="123" className="border-2 focus:border-food-primary h-12 text-base" {...field} />
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
                          className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg disabled:opacity-50"
                        >
                          {isSubmitting ? "Processing..." : `Place Order • $${total.toFixed(2)}`}
                        </Button>
                      </form>
                    </Form>
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
                        {deliveryFee > 0 && (
                          <div className="flex justify-between">
                            <span className="text-gray-600">Delivery Fee</span>
                            <span className="font-medium">${deliveryFee.toFixed(2)}</span>
                          </div>
                        )}
                        <div className="flex justify-between">
                          <span className="text-gray-600">Tax (10%)</span>
                          <span className="font-medium">${tax.toFixed(2)}</span>
                        </div>
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
