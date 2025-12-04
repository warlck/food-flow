import React, { useState } from "react";
import * as z from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { MapPin, Clock, CreditCard, ArrowLeft, CheckCircle, User, Phone, Mail, Home, MapPinCheck, Package } from "lucide-react";

import Layout from "@/components/Layout";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useCart } from "@/context/CartContext";
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

// Default values for restaurant settings
const DEFAULT_DELIVERY_FEE = 3.99;
const DEFAULT_DELIVERY_TIME = { min: 30, max: 45 };
const DEFAULT_PICKUP_TIME = { min: 15, max: 25 };
const DEFAULT_RESTAURANT_NAME = "Restaurant";
const DEFAULT_RESTAURANT_ADDRESS = "123 Main Street, City, State 12345";
const DEFAULT_RESTAURANT_PHONE = "(555) 123-4567";

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

const Checkout: React.FC = () => {
  const navigate = useNavigate();
  const { items, getTotalPrice, clearCart, orderType } = useCart();
  const [step, setStep] = useState(1);
  const [paymentMethod, setPaymentMethod] = useState<"creditCard" | "payAtLocation">("creditCard");
  const [orderDetails, setOrderDetails] = useState<any>({});
  
  // Calculate totals
  const subtotal = getTotalPrice();
  const deliveryFee = orderType === "delivery" ? DEFAULT_DELIVERY_FEE : 0;
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

  // Submit handler for customer info step
  const onSubmitCustomerInfo = (data: any) => {
    setOrderDetails({ ...orderDetails, ...data });
    setStep(2); // Go to payment step
  };

  // Submit handler for payment step
  const onSubmitPayment = (data: any) => {
    setOrderDetails({ ...orderDetails, ...data });
    
    // Simulate order submission
    setTimeout(() => {
      clearCart();
      toast.success("Order placed successfully!");
      navigate("/");
    }, 1500);
  };

  // Handle going back to previous step
  const handleBackStep = () => {
    setStep(step - 1);
  };

  // Show estimated time for either delivery or pickup
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
        <div className="mb-6">
          <Button 
            variant="ghost" 
            onClick={() => navigate("/cart")} 
            className="flex items-center mb-4"
          >
            <ArrowLeft size={18} className="mr-2" /> Back to Cart
          </Button>
          <h1 className="text-3xl font-bold">Checkout</h1>
          <p className="text-gray-600 mt-2">
            {orderType === "delivery" ? "Delivery" : "Pickup"} Order
          </p>
        </div>

        {/* Progress indicator */}
        <div className="mb-8">
          <div className="flex items-center justify-center">
            <div className="flex items-center">
              <div className={`rounded-full h-10 w-10 flex items-center justify-center ${step === 1 ? "bg-food-primary text-white" : "bg-green-500 text-white"}`}>
                1
              </div>
              <div className="mx-4 h-1 w-24 bg-gray-200">
                <div className={`h-full ${step === 2 ? "bg-food-primary" : step > 2 ? "bg-green-500" : "bg-gray-200"}`} style={{ width: step > 1 ? "100%" : "0%" }}></div>
              </div>
              <div className={`rounded-full h-10 w-10 flex items-center justify-center ${step === 2 ? "bg-food-primary text-white" : step > 2 ? "bg-green-500 text-white" : "bg-gray-200 text-gray-600"}`}>
                2
              </div>
            </div>
          </div>
          <div className="flex items-center justify-center mt-2">
            <span className={`text-sm font-medium mr-16 ${step === 1 ? "text-food-primary" : "text-green-500"}`}>
              Information
            </span>
            <span className={`text-sm font-medium ${step === 2 ? "text-food-primary" : step > 2 ? "text-green-500" : "text-gray-500"}`}>
              Payment
            </span>
          </div>
        </div>

        <div className="flex flex-col lg:flex-row gap-8">
          {/* Checkout Form */}
          <div className="lg:w-2/3">
            <Card className="p-6">
              {step === 1 && (
                <>
                  <h2 className="text-xl font-semibold mb-4">
                    {orderType === "delivery" ? "Delivery" : "Pickup"} Information
                  </h2>
                  
                  {orderType === "delivery" ? (
                    <Form {...deliveryForm}>
                      <form onSubmit={deliveryForm.handleSubmit(onSubmitCustomerInfo)} className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <FormField
                            control={deliveryForm.control}
                            name="name"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Full Name</FormLabel>
                                <FormControl>
                                  <Input placeholder="John Doe" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                          
                          <FormField
                            control={deliveryForm.control}
                            name="email"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Email</FormLabel>
                                <FormControl>
                                  <Input placeholder="you@example.com" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </div>

                        <FormField
                          control={deliveryForm.control}
                          name="phone"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Phone Number</FormLabel>
                              <FormControl>
                                <Input placeholder="(123) 456-7890" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <div className="border-t pt-4 pb-2">
                          <h3 className="text-lg font-medium flex items-center mb-4">
                            <MapPin className="mr-2" size={18} />
                            Delivery Address
                          </h3>
                        </div>

                        <FormField
                          control={deliveryForm.control}
                          name="street"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Street Address</FormLabel>
                              <FormControl>
                                <Input placeholder="123 Main St" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                          <FormField
                            control={deliveryForm.control}
                            name="city"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>City</FormLabel>
                                <FormControl>
                                  <Input placeholder="San Francisco" {...field} />
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
                                  <Input placeholder="CA" {...field} />
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
                                  <Input placeholder="94103" {...field} />
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
                                  placeholder="Gate code, apartment number, or other special instructions"
                                  className="resize-none" 
                                  {...field} 
                                />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <Button type="submit" className="w-full mt-6 bg-food-primary hover:bg-food-accent">
                          Continue to Payment
                        </Button>
                      </form>
                    </Form>
                  ) : (
                    <Form {...pickupForm}>
                      <form onSubmit={pickupForm.handleSubmit(onSubmitCustomerInfo)} className="space-y-4">
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <FormField
                            control={pickupForm.control}
                            name="name"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Full Name</FormLabel>
                                <FormControl>
                                  <Input placeholder="John Doe" {...field} />
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
                                <FormLabel>Email</FormLabel>
                                <FormControl>
                                  <Input placeholder="you@example.com" {...field} />
                                </FormControl>
                                <FormMessage />
                              </FormItem>
                            )}
                          />
                        </div>

                        <FormField
                          control={pickupForm.control}
                          name="phone"
                          render={({ field }) => (
                            <FormItem>
                              <FormLabel>Phone Number</FormLabel>
                              <FormControl>
                                <Input placeholder="(123) 456-7890" {...field} />
                              </FormControl>
                              <FormMessage />
                            </FormItem>
                          )}
                        />

                        <div className="border-t pt-4 pb-2">
                          <h3 className="text-lg font-medium flex items-center">
                            <MapPinCheck className="mr-2" size={18} />
                            Pickup Location
                          </h3>
                          <div className="bg-gray-50 p-4 mt-3 rounded-lg">
                            <p className="font-medium">{DEFAULT_RESTAURANT_NAME}</p>
                            <p className="text-gray-600">{DEFAULT_RESTAURANT_ADDRESS}</p>
                            <p className="text-gray-600">{DEFAULT_RESTAURANT_PHONE}</p>
                          </div>
                          <div className="mt-4 flex items-center text-gray-600">
                            <Clock size={16} className="mr-2" />
                            <span>Estimated pickup time: {getEstimatedTime()}</span>
                          </div>
                        </div>

                        <Button type="submit" className="w-full mt-6 bg-food-primary hover:bg-food-accent">
                          Continue to Payment
                        </Button>
                      </form>
                    </Form>
                  )}
                </>
              )}

              {step === 2 && (
                <>
                  <div className="flex items-center mb-4">
                    <Button 
                      variant="ghost" 
                      type="button" 
                      onClick={handleBackStep} 
                      className="p-0 mr-3 h-auto"
                    >
                      <ArrowLeft size={18} />
                    </Button>
                    <h2 className="text-xl font-semibold">Payment</h2>
                  </div>

                  <Form {...paymentForm}>
                    <form onSubmit={paymentForm.handleSubmit(onSubmitPayment)} className="space-y-6">
                      <FormField
                        control={paymentForm.control}
                        name="paymentMethod"
                        render={({ field }) => (
                          <FormItem>
                            <FormLabel>Payment Method</FormLabel>
                            <FormControl>
                              <Select 
                                value={field.value} 
                                onValueChange={(value: "creditCard" | "payAtLocation") => {
                                  field.onChange(value);
                                  setPaymentMethod(value);
                                }}
                              >
                                <SelectTrigger>
                                  <SelectValue placeholder="Select payment method" />
                                </SelectTrigger>
                                <SelectContent>
                                  <SelectItem value="creditCard">Credit Card</SelectItem>
                                  <SelectItem value="payAtLocation">
                                    {orderType === "delivery" ? "Pay on Delivery" : "Pay at Pickup"}
                                  </SelectItem>
                                </SelectContent>
                              </Select>
                            </FormControl>
                            <FormMessage />
                          </FormItem>
                        )}
                      />

                      {paymentMethod === "creditCard" && (
                        <>
                          <div className="border-t pt-4 pb-2">
                            <h3 className="text-lg font-medium flex items-center mb-4">
                              <CreditCard className="mr-2" size={18} />
                              Card Details
                            </h3>
                          </div>

                          <FormField
                            control={paymentForm.control}
                            name="cardNumber"
                            render={({ field }) => (
                              <FormItem>
                                <FormLabel>Card Number</FormLabel>
                                <FormControl>
                                  <Input placeholder="1234 5678 9012 3456" {...field} />
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
                                    <Input placeholder="MM/YY" {...field} />
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
                                    <Input placeholder="123" {...field} />
                                  </FormControl>
                                  <FormMessage />
                                </FormItem>
                              )}
                            />
                          </div>
                        </>
                      )}

                      <div className="pt-4">
                        <Button 
                          type="submit" 
                          className="w-full bg-food-primary hover:bg-food-accent"
                        >
                          Place Order - ${total.toFixed(2)}
                        </Button>
                        <p className="text-center text-sm text-gray-500 mt-4">
                          By placing your order, you agree to our Terms of Service and Privacy Policy
                        </p>
                      </div>
                    </form>
                  </Form>
                </>
              )}
            </Card>
          </div>

          {/* Order Summary */}
          <div className="lg:w-1/3">
            <Card className="p-6 sticky top-24">
              <h2 className="text-xl font-semibold mb-4">Order Summary</h2>
              
              {/* Order details */}
              <div className="space-y-4">
                <div className="max-h-64 overflow-y-auto pr-2 space-y-3 mb-4">
                  {items.map((item) => (
                    <div key={item.menuItem.id} className="flex justify-between text-sm">
                      <div className="flex items-start">
                        <span className="bg-gray-100 rounded-full h-5 w-5 flex items-center justify-center text-xs mr-2">
                          {item.quantity}
                        </span>
                        <span>{item.menuItem.name}</span>
                      </div>
                      <span>${(item.menuItem.price * item.quantity).toFixed(2)}</span>
                    </div>
                  ))}
                </div>

                <div className="border-t pt-4 space-y-2">
                  <div className="flex justify-between text-sm">
                    <span>Subtotal</span>
                    <span>${subtotal.toFixed(2)}</span>
                  </div>
                  
                  {orderType === "delivery" && (
                    <div className="flex justify-between text-sm">
                      <span>Delivery Fee</span>
                      <span>${deliveryFee.toFixed(2)}</span>
                    </div>
                  )}
                  
                  <div className="flex justify-between text-sm">
                    <span>Tax</span>
                    <span>${tax.toFixed(2)}</span>
                  </div>
                  
                  <div className="flex justify-between font-medium text-base pt-2 border-t">
                    <span>Total</span>
                    <span>${total.toFixed(2)}</span>
                  </div>
                </div>

                <div className="bg-gray-50 p-3 rounded-md mt-4">
                  <div className="flex items-start">
                    {orderType === "delivery" ? (
                      <MapPin className="h-5 w-5 mr-2 mt-0.5 text-gray-600 flex-shrink-0" />
                    ) : (
                      <Package className="h-5 w-5 mr-2 mt-0.5 text-gray-600 flex-shrink-0" />
                    )}
                    <div>
                      <p className="font-medium">
                        {orderType === "delivery" ? "Delivery" : "Pickup"}
                      </p>
                      <p className="text-sm text-gray-600">
                        Estimated {orderType === "delivery" ? "delivery" : "pickup"} time: {getEstimatedTime()}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </Card>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default Checkout;
