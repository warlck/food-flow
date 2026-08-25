import React, { useEffect, useState } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import { orderService, Order } from '@/services/orderService';
import Layout from '@/components/Layout';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { CheckCircle, XCircle, Clock, MapPin, Phone, Mail, Package, Loader2, ArrowLeft } from 'lucide-react';

const ORDER_STATUS_LABELS: Record<string, string> = {
  pending: 'Pending',
  confirmed: 'Confirmed',
  preparing: 'Preparing',
  ready: 'Ready',
  out_for_delivery: 'Out for delivery',
  completed: 'Completed',
  cancelled: 'Cancelled',
};

const ORDER_STATUS_COLORS: Record<string, string> = {
  pending: 'text-yellow-600',
  confirmed: 'text-green-600',
  preparing: 'text-orange-600',
  ready: 'text-green-600',
  out_for_delivery: 'text-indigo-600',
  completed: 'text-green-600',
  cancelled: 'text-red-600',
};

const OrderConfirmation: React.FC = () => {
  const { orderId } = useParams<{ orderId: string }>();
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Check for redirect status from Stripe
  const redirectStatus = searchParams.get('redirect_status');
  const paymentIntentClientSecret = searchParams.get('payment_intent_client_secret');

  useEffect(() => {
    const fetchOrder = async () => {
      if (!orderId) {
        setError('Order ID not found');
        setLoading(false);
        return;
      }

      try {
        // If we have a payment_intent_client_secret, confirm the payment first
        if (paymentIntentClientSecret && redirectStatus === 'succeeded') {
          await orderService.confirmPayment(orderId);
        }

        const data = await orderService.getOrder(orderId);
        setOrder(data);
      } catch (err) {
        console.error('Failed to fetch order:', err);
        setError('Failed to load order details');
      } finally {
        setLoading(false);
      }
    };

    fetchOrder();
  }, [orderId, paymentIntentClientSecret, redirectStatus]);

  if (loading) {
    return (
      <Layout>
        <div className="container mx-auto px-4 py-16 flex flex-col items-center justify-center min-h-[60vh]">
          <Loader2 className="w-12 h-12 animate-spin text-food-primary mb-4" />
          <p className="text-gray-600">Loading order details...</p>
        </div>
      </Layout>
    );
  }

  if (error || !order) {
    return (
      <Layout>
        <div className="container mx-auto px-4 py-16 flex flex-col items-center justify-center min-h-[60vh]">
          <XCircle className="w-16 h-16 text-red-500 mb-4" />
          <h1 className="text-2xl font-bold text-gray-900 mb-2">Order Not Found</h1>
          <p className="text-gray-600 mb-6">{error || 'Unable to find order details'}</p>
          <Button onClick={() => navigate('/')} className="bg-food-primary hover:bg-food-accent">
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Home
          </Button>
        </div>
      </Layout>
    );
  }

  const isPaid = order.paymentStatus === 'paid';
  const isPending = order.paymentStatus === 'pending' || order.paymentStatus === 'processing';
  const isFailed = order.paymentStatus === 'failed';

  return (
    <Layout>
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        {/* Status Header */}
        <div className="text-center mb-8">
          {isPaid && (
            <>
              <CheckCircle className="w-20 h-20 text-green-500 mx-auto mb-4" />
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Order Confirmed!</h1>
              <p className="text-gray-600">Thank you for your order. We've received your payment.</p>
            </>
          )}
          {isPending && (
            <>
              <Clock className="w-20 h-20 text-yellow-500 mx-auto mb-4" />
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Payment Pending</h1>
              <p className="text-gray-600">Your order is being processed. We'll update you shortly.</p>
            </>
          )}
          {isFailed && (
            <>
              <XCircle className="w-20 h-20 text-red-500 mx-auto mb-4" />
              <h1 className="text-3xl font-bold text-gray-900 mb-2">Payment Failed</h1>
              <p className="text-gray-600">There was an issue with your payment. Please try again.</p>
            </>
          )}
          <p className="text-sm text-gray-500 mt-2">Order #{order.id.slice(0, 8).toUpperCase()}</p>
        </div>

        <div className="grid md:grid-cols-2 gap-6">
          {/* Order Details */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Package className="w-5 h-5 mr-2 text-food-primary" />
                Order Details
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex justify-between">
                <span className="text-gray-600">Status:</span>
                <span className={`font-semibold capitalize ${ORDER_STATUS_COLORS[order.orderStatus] ?? 'text-gray-900'}`}>
                  {ORDER_STATUS_LABELS[order.orderStatus] ?? order.orderStatus}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Payment:</span>
                <span className={`font-semibold capitalize ${
                  isPaid ? 'text-green-600' :
                  isPending ? 'text-yellow-600' :
                  isFailed ? 'text-red-600' : 'text-gray-900'
                }`}>
                  {order.paymentStatus}
                </span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Order Type:</span>
                <span className="font-semibold capitalize">{order.orderType}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-600">Payment Method:</span>
                <span className="font-semibold capitalize">
                  {order.paymentMethod === 'creditCard' ? 'Credit Card' : 'Pay at Location'}
                </span>
              </div>

              <div className="border-t pt-4 mt-4">
                <h3 className="font-semibold mb-3">Items:</h3>
                {order.items.map((item) => (
                  <div key={item.id} className="py-2 border-b border-gray-100 last:border-0">
                    <div className="flex justify-between">
                      <span className="text-gray-700">
                        {item.quantity}x {item.menuItemName}
                      </span>
                      <span className="font-medium">${(item.menuItemPrice * item.quantity).toFixed(2)}</span>
                    </div>
                    {item.addons && item.addons.length > 0 && (
                      <div className="ml-6 mt-1 space-y-0.5">
                        {item.addons.map((addon) => (
                          <div key={addon.id} className="flex justify-between text-xs text-gray-600">
                            <span>+ {addon.addonName} x{addon.quantity}</span>
                            <span className="text-food-primary">
                              +${(addon.addonPrice * addon.quantity * item.quantity).toFixed(2)}
                            </span>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                ))}
              </div>

              <div className="border-t pt-4 space-y-2">
                <div className="flex justify-between">
                  <span className="text-gray-600">Subtotal:</span>
                  <span>${order.subtotal.toFixed(2)}</span>
                </div>
                {order.deliveryFee > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-600">Delivery Fee:</span>
                    <span>${order.deliveryFee.toFixed(2)}</span>
                  </div>
                )}
                {order.tax > 0 && (
                  <div className="flex justify-between">
                    <span className="text-gray-600">Tax:</span>
                    <span>${order.tax.toFixed(2)}</span>
                  </div>
                )}
                <div className="flex justify-between font-bold text-lg border-t pt-2">
                  <span>Total:</span>
                  <span className="text-food-primary">${order.total.toFixed(2)}</span>
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Customer & Delivery Info */}
          <div className="space-y-6">
            {/* Customer Info */}
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Customer Information</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center text-gray-700">
                  <span className="font-medium">{order.customerName}</span>
                </div>
                <div className="flex items-center text-gray-600">
                  <Mail className="w-4 h-4 mr-2 text-food-primary" />
                  {order.customerEmail}
                </div>
                <div className="flex items-center text-gray-600">
                  <Phone className="w-4 h-4 mr-2 text-food-primary" />
                  {order.customerPhone}
                </div>
              </CardContent>
            </Card>

            {/* Delivery Address */}
            {order.deliveryAddress && (
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg flex items-center">
                    <MapPin className="w-5 h-5 mr-2 text-food-primary" />
                    Delivery Address
                  </CardTitle>
                </CardHeader>
                <CardContent>
                  <p className="text-gray-700">{order.deliveryAddress.street}</p>
                  <p className="text-gray-700">
                    {order.deliveryAddress.city}{order.deliveryAddress.state ? `, ${order.deliveryAddress.state}` : ''} {order.deliveryAddress.postalCode}
                  </p>
                  {order.deliveryAddress.deliveryInstructions && (
                    <p className="text-gray-600 mt-2 italic">
                      Note: {order.deliveryAddress.deliveryInstructions}
                    </p>
                  )}
                </CardContent>
              </Card>
            )}

            {/* Confirmation Email Notice */}
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 text-center">
              <p className="text-blue-800">
                A confirmation email has been sent to <strong>{order.customerEmail}</strong>
              </p>
            </div>
          </div>
        </div>

        {/* Action Buttons */}
        <div className="mt-8 flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button
            onClick={() => navigate(`/track-order/${order.id}`)}
            className="w-full sm:w-auto bg-food-primary hover:bg-food-accent text-white font-semibold shadow-md px-6 py-2.5 flex items-center justify-center"
          >
            <Clock className="w-5 h-5 mr-2" />
            Track Your Order Live
          </Button>

          <Button
            onClick={() => navigate(order?.restaurantId ? `/restaurant/${order.restaurantId}` : '/')}
            variant="outline"
            className="w-full sm:w-auto border-gray-300 text-gray-700 hover:bg-gray-100"
          >
            <ArrowLeft className="w-4 h-4 mr-2" />
            Back to Home
          </Button>
        </div>
      </div>
    </Layout>
  );
};

export default OrderConfirmation;
