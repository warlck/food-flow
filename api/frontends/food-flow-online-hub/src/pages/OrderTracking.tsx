import React, { useEffect, useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { orderService, Order } from '@/services/orderService';
import Layout from '@/components/Layout';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
import { useCart } from '@/context/CartContext';
import {
  Package,
  Search,
  Clock,
  CheckCircle2,
  ChefHat,
  Bike,
  Home,
  Store,
  RefreshCw,
  Copy,
  Share2,
  Phone,
  Mail,
  MapPin,
  AlertTriangle,
  ArrowLeft,
  Loader2,
} from 'lucide-react';

const ORDER_STATUS_LABELS: Record<string, string> = {
  pending: 'Order Placed',
  confirmed: 'Order Confirmed',
  preparing: 'Kitchen Preparing',
  ready: 'Ready for Pickup',
  out_for_delivery: 'Out for Delivery',
  completed: 'Completed',
  cancelled: 'Cancelled',
};

const getOrderStatusLabel = (status: string, isDelivery: boolean): string => {
  if (status === 'ready') {
    return isDelivery ? 'Ready for Delivery' : 'Ready for Pickup';
  }
  return ORDER_STATUS_LABELS[status] || status;
};

const ORDER_STATUS_COLORS: Record<string, string> = {
  pending: 'bg-yellow-100 text-yellow-800 border-yellow-300',
  confirmed: 'bg-blue-100 text-blue-800 border-blue-300',
  preparing: 'bg-orange-100 text-orange-800 border-orange-300',
  ready: 'bg-emerald-100 text-emerald-800 border-emerald-300',
  out_for_delivery: 'bg-purple-100 text-purple-800 border-purple-300',
  completed: 'bg-green-100 text-green-800 border-green-300',
  cancelled: 'bg-red-100 text-red-800 border-red-300',
};

interface StepConfig {
  key: string;
  label: string;
  description: string;
  icon: React.ComponentType<{ className?: string }>;
}

const DELIVERY_STEPS: StepConfig[] = [
  { key: 'placed', label: 'Order Placed', description: 'Received by restaurant', icon: Package },
  { key: 'preparing', label: 'Preparing', description: 'Kitchen is cooking your food', icon: ChefHat },
  { key: 'out_for_delivery', label: 'Out for Delivery', description: 'Courier on the way to you', icon: Bike },
  { key: 'completed', label: 'Delivered', description: 'Enjoy your meal!', icon: Home },
];

const PICKUP_STEPS: StepConfig[] = [
  { key: 'placed', label: 'Order Placed', description: 'Received by restaurant', icon: Package },
  { key: 'preparing', label: 'Preparing', description: 'Kitchen is cooking your food', icon: ChefHat },
  { key: 'ready', label: 'Ready for Pickup', description: 'Ready at restaurant counter', icon: Store },
  { key: 'completed', label: 'Picked Up', description: 'Order completed', icon: CheckCircle2 },
];

function getStepIndex(status: string, isDelivery: boolean): number {
  switch (status) {
    case 'pending':
    case 'confirmed':
      return 0;
    case 'preparing':
      return 1;
    case 'ready':
      return isDelivery ? 1 : 2;
    case 'out_for_delivery':
      return isDelivery ? 2 : 2;
    case 'completed':
      return 3;
    default:
      return 0;
  }
}

const OrderTracking: React.FC = () => {
  const { orderId } = useParams<{ orderId?: string }>();
  const navigate = useNavigate();
  const { toast } = useToast();
  const { restaurantId } = useCart();

  const [inputOrderId, setInputOrderId] = useState(orderId || '');
  const [order, setOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState<boolean>(!!orderId);
  const [isRefreshing, setIsRefreshing] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);
  const [autoRefresh, setAutoRefresh] = useState<boolean>(true);
  const [lastRefreshedAt, setLastRefreshedAt] = useState<Date | null>(null);

  const fetchOrder = useCallback(async (id: string, isBackground = false) => {
    if (!id) return;
    if (isBackground) {
      setIsRefreshing(true);
    } else {
      setLoading(true);
    }
    setError(null);

    try {
      const data = await orderService.getOrder(id);
      setOrder(data);
      setLastRefreshedAt(new Date());
    } catch (err) {
      console.error('Error fetching order:', err);
      if (!isBackground) {
        setError(err instanceof Error ? err.message : 'Could not find order. Please verify your Order ID.');
        setOrder(null);
      }
    } finally {
      setLoading(false);
      setIsRefreshing(false);
    }
  }, []);

  // Sync initial param
  useEffect(() => {
    if (orderId) {
      setInputOrderId(orderId);
      fetchOrder(orderId);
    } else {
      setOrder(null);
      setLoading(false);
    }
  }, [orderId, fetchOrder]);

  // Polling effect
  useEffect(() => {
    if (!orderId || !autoRefresh || !order) return;
    if (order.orderStatus === 'completed' || order.orderStatus === 'cancelled') return;

    const interval = setInterval(() => {
      fetchOrder(orderId, true);
    }, 10000);

    return () => clearInterval(interval);
  }, [orderId, autoRefresh, order, fetchOrder]);

  const handleSearchSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    const trimmed = inputOrderId.trim();
    if (!trimmed) {
      toast({
        title: 'Order ID Required',
        description: 'Please enter a valid order ID to track.',
        variant: 'destructive',
      });
      return;
    }
    navigate(`/track-order/${trimmed}`);
  };

  const handleCopyOrderId = () => {
    if (!order) return;
    navigator.clipboard.writeText(order.id);
    toast({
      title: 'Copied!',
      description: 'Order ID copied to clipboard.',
    });
  };

  const handleShareLink = () => {
    navigator.clipboard.writeText(window.location.href);
    toast({
      title: 'Link Copied!',
      description: 'Tracking link copied to clipboard.',
    });
  };

  const isDelivery = order?.orderType === 'delivery';
  const currentStep = order ? getStepIndex(order.orderStatus, isDelivery) : 0;
  const steps = isDelivery ? DELIVERY_STEPS : PICKUP_STEPS;
  const isCancelled = order?.orderStatus === 'cancelled';
  const isCompleted = order?.orderStatus === 'completed';

  return (
    <Layout>
      <div className="container mx-auto px-4 py-8 max-w-4xl min-h-[75vh]">
        {/* Search Header Card */}
        <Card className="mb-8 border-gray-200 shadow-sm bg-gradient-to-r from-orange-50 to-amber-50 dark:from-gray-800 dark:to-gray-900">
          <CardContent className="pt-6">
            <form onSubmit={handleSearchSubmit} className="flex flex-col sm:flex-row items-center gap-3">
              <div className="relative flex-1 w-full">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 w-5 h-5" />
                <Input
                  type="text"
                  placeholder="Enter Order ID (e.g., 550e8400-e29b-41d4-a716-446655440000)"
                  value={inputOrderId}
                  onChange={(e) => setInputOrderId(e.target.value)}
                  className="pl-10 bg-white dark:bg-gray-800 h-12 text-base shadow-inner border-gray-300 focus-visible:ring-food-primary"
                />
              </div>
              <Button
                type="submit"
                className="w-full sm:w-auto h-12 px-6 bg-food-primary hover:bg-food-accent text-white font-semibold flex items-center justify-center gap-2"
              >
                <Search className="w-4 h-4" />
                Track Order
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Loading State */}
        {loading && (
          <div className="flex flex-col items-center justify-center py-20">
            <Loader2 className="w-12 h-12 animate-spin text-food-primary mb-4" />
            <p className="text-gray-600 font-medium">Fetching order status...</p>
          </div>
        )}

        {/* Error / Not Found State */}
        {!loading && error && (
          <Card className="text-center py-12 border-red-200 bg-red-50/50">
            <CardContent>
              <AlertTriangle className="w-16 h-16 text-red-500 mx-auto mb-4" />
              <h2 className="text-2xl font-bold text-gray-900 mb-2">Order Not Found</h2>
              <p className="text-gray-600 max-w-md mx-auto mb-6">{error}</p>
              <Button
                onClick={() => navigate(restaurantId ? `/restaurant/${restaurantId}` : '/')}
                variant="outline"
                className="border-food-primary text-food-primary hover:bg-food-primary/10"
              >
                <ArrowLeft className="w-4 h-4 mr-2" />
                Back to Home
              </Button>
            </CardContent>
          </Card>
        )}

        {/* No Order Selected Prompt */}
        {!loading && !error && !order && !orderId && (
          <Card className="text-center py-16 border-dashed border-2">
            <CardContent>
              <Package className="w-16 h-16 text-gray-400 mx-auto mb-4" />
              <h2 className="text-xl font-semibold text-gray-800 mb-2">Track Your Food Order Live</h2>
              <p className="text-gray-500 max-w-md mx-auto mb-6">
                Enter your unique Order ID in the search bar above to see live updates, estimated time, and order details.
              </p>
            </CardContent>
          </Card>
        )}

        {/* Order Tracking Dashboard */}
        {!loading && !error && order && (
          <div className="space-y-6 animate-fade-in">
            {/* Live Refresh Control & Header */}
            <div className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm space-y-2">
              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                <div className="flex items-center gap-2.5 flex-wrap min-h-[36px]">
                  <h1 className="text-2xl font-bold text-gray-900 leading-none">
                    Order #{order.id.slice(0, 8).toUpperCase()}
                  </h1>
                  <Badge className={`h-7 px-3 inline-flex items-center text-xs font-semibold rounded-full border ${ORDER_STATUS_COLORS[order.orderStatus] || 'bg-gray-100'}`}>
                    {getOrderStatusLabel(order.orderStatus, isDelivery)}
                  </Badge>
                  <Badge variant="outline" className="h-7 px-3 inline-flex items-center capitalize text-xs font-medium border-gray-300 rounded-full">
                    {order.orderType}
                  </Badge>
                </div>

                <div className="flex items-center gap-2 flex-wrap sm:flex-nowrap min-h-[36px]">
                  {!isCompleted && !isCancelled && (
                    <div className="flex items-center space-x-2 bg-gray-50 px-3 h-9 rounded-lg border border-gray-200">
                      <Switch
                        id="auto-refresh"
                        checked={autoRefresh}
                        onCheckedChange={setAutoRefresh}
                      />
                      <Label htmlFor="auto-refresh" className="text-xs text-gray-600 flex items-center gap-1.5 cursor-pointer select-none">
                        <span className={`w-2 h-2 rounded-full ${autoRefresh ? 'bg-green-500 animate-pulse' : 'bg-gray-400'}`}></span>
                        Auto-Refresh
                      </Label>
                    </div>
                  )}

                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => fetchOrder(order.id, true)}
                    disabled={isRefreshing}
                    className="flex items-center gap-1.5 text-xs h-9 px-3"
                  >
                    <RefreshCw className={`w-3.5 h-3.5 ${isRefreshing ? 'animate-spin' : ''}`} />
                    Refresh
                  </Button>
                  <Button variant="outline" size="icon" onClick={handleShareLink} title="Share Link" className="h-9 w-9">
                    <Share2 className="w-4 h-4 text-gray-600" />
                  </Button>
                  <Button variant="outline" size="icon" onClick={handleCopyOrderId} title="Copy ID" className="h-9 w-9">
                    <Copy className="w-4 h-4 text-gray-600" />
                  </Button>
                </div>
              </div>

              <p className="text-xs text-gray-500 flex items-center gap-2 pt-0.5">
                <span>Placed on {new Date(order.dateCreated).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
                {lastRefreshedAt && (
                  <span>• Updated {lastRefreshedAt.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}</span>
                )}
              </p>
            </div>

            {/* Cancelled Banner */}
            {isCancelled && (
              <div className="bg-red-50 border border-red-200 rounded-xl p-4 flex items-center gap-3 text-red-800">
                <AlertTriangle className="w-6 h-6 text-red-600 flex-shrink-0" />
                <div>
                  <h3 className="font-semibold text-sm">Order Cancelled</h3>
                  <p className="text-xs text-red-700">This order was cancelled. If you have questions, please contact customer support.</p>
                </div>
              </div>
            )}

            {/* Ready State Info Banner */}
            {!isCancelled && order.orderStatus === 'ready' && (
              <div className="bg-emerald-50 border border-emerald-200 rounded-xl p-4 flex items-center gap-3 text-emerald-900">
                {isDelivery ? (
                  <>
                    <Package className="w-6 h-6 text-emerald-600 flex-shrink-0" />
                    <div>
                      <h3 className="font-semibold text-sm">Your order is ready!</h3>
                      <p className="text-xs text-emerald-700">Food preparation is complete. Your order is waiting to be dispatched to a courier for delivery.</p>
                    </div>
                  </>
                ) : (
                  <>
                    <Store className="w-6 h-6 text-emerald-600 flex-shrink-0" />
                    <div>
                      <h3 className="font-semibold text-sm">Ready for Pickup!</h3>
                      <p className="text-xs text-emerald-700">Your order is ready to be collected at the restaurant counter.</p>
                    </div>
                  </>
                )}
              </div>
            )}

            {/* Status Timeline Stepper */}
            {!isCancelled && (
              <Card className="border-gray-200 shadow-sm overflow-hidden">
                <CardHeader className="bg-gray-50 border-b border-gray-100 py-4">
                  <CardTitle className="text-lg font-semibold flex items-center justify-between">
                    <span className="flex items-center gap-2">
                      <Clock className="w-5 h-5 text-food-primary" />
                      Order Status
                    </span>
                    {!isCompleted && (
                      <span className="text-xs font-normal text-gray-500 bg-amber-100 text-amber-900 px-3 py-1 rounded-full">
                        Estimated fulfillment: ~20-30 mins
                      </span>
                    )}
                  </CardTitle>
                </CardHeader>

                <CardContent className="pt-8 pb-10 px-4 sm:px-8">
                  {/* Stepper Progress Line Container */}
                  <div className="relative flex items-start justify-between max-w-2xl mx-auto">
                    {/* Horizontal background track (aligned with 24px center of 48px circles) */}
                    <div
                      className="absolute top-6 -translate-y-1/2 h-1 bg-gray-200 z-0"
                      style={{ left: '24px', right: '24px' }}
                    />
                    {/* Active progress line (spans from center of step 1 to center of current step) */}
                    <div
                      className="absolute top-6 -translate-y-1/2 h-1 bg-food-primary transition-all duration-500 z-0"
                      style={{
                        left: '24px',
                        width: `calc((100% - 48px) * ${currentStep / (steps.length - 1)})`,
                      }}
                    />

                    {steps.map((step, idx) => {
                      const IconComponent = step.icon;
                      const isDone = idx <= currentStep;
                      const isCurrent = idx === currentStep && !isCompleted;

                      const isReadyDeliveryStep = isDelivery && order.orderStatus === 'ready' && step.key === 'preparing';
                      const stepLabel = isReadyDeliveryStep ? 'Food Ready' : step.label;
                      const stepDesc = isReadyDeliveryStep ? 'Awaiting driver pickup' : step.description;

                      return (
                        <div key={step.key} className="relative z-10 flex flex-col items-center group">
                          <div
                            className={`w-12 h-12 rounded-full flex items-center justify-center border-2 transition-all duration-300 ${
                              isDone
                                ? 'bg-food-primary border-food-primary text-white shadow-md'
                                : 'bg-white border-gray-300 text-gray-400'
                            } ${isCurrent ? 'ring-4 ring-orange-200 scale-110' : ''}`}
                          >
                            <IconComponent className="w-5 h-5" />
                          </div>
                          <div className="mt-3 text-center">
                            <p className={`text-xs font-semibold ${isDone ? 'text-gray-900' : 'text-gray-400'}`}>
                              {stepLabel}
                            </p>
                            <p className="text-[10px] text-gray-500 max-w-[90px] hidden sm:block mt-0.5">
                              {stepDesc}
                            </p>
                          </div>
                        </div>
                      );
                    })}
                  </div>
                </CardContent>
              </Card>
            )}

            {/* Main Content Grid: Items & Details */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {/* Left Column: Items Breakdown */}
              <Card className="md:col-span-2 border-gray-200 shadow-sm">
                <CardHeader>
                  <CardTitle className="text-lg flex items-center gap-2">
                    <Package className="w-5 h-5 text-food-primary" />
                    Items Ordered ({order.items.reduce((acc, item) => acc + item.quantity, 0)})
                  </CardTitle>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="divide-y divide-gray-100">
                    {order.items.map((item) => (
                      <div key={item.id} className="py-3 flex justify-between items-start first:pt-0 last:pb-0">
                        <div>
                          <p className="font-medium text-gray-900 text-sm">
                            {item.quantity}x {item.menuItemName}
                          </p>
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
                          {item.specialInstructions && (
                            <p className="text-xs text-gray-500 italic mt-0.5">
                              Note: {item.specialInstructions}
                            </p>
                          )}
                        </div>
                        <span className="font-semibold text-sm text-gray-800">
                          ${(item.menuItemPrice * item.quantity).toFixed(2)}
                        </span>
                      </div>
                    ))}
                  </div>

                  <div className="border-t pt-4 space-y-2 text-sm">
                    <div className="flex justify-between text-gray-600">
                      <span>Subtotal</span>
                      <span>${order.subtotal.toFixed(2)}</span>
                    </div>
                    {order.deliveryFee > 0 && (
                      <div className="flex justify-between text-gray-600">
                        <span>Delivery Fee</span>
                        <span>${order.deliveryFee.toFixed(2)}</span>
                      </div>
                    )}
                    {order.tax > 0 && (
                      <div className="flex justify-between text-gray-600">
                        <span>Tax</span>
                        <span>${order.tax.toFixed(2)}</span>
                      </div>
                    )}
                    <div className="flex justify-between font-bold text-base border-t pt-3 text-gray-900">
                      <span>Total Paid</span>
                      <span className="text-food-primary">${order.total.toFixed(2)}</span>
                    </div>
                  </div>
                </CardContent>
              </Card>

              {/* Right Column: Customer & Delivery Info */}
              <div className="space-y-6">
                <Card className="border-gray-200 shadow-sm">
                  <CardHeader>
                    <CardTitle className="text-lg flex items-center gap-2">
                      {isDelivery ? <MapPin className="w-5 h-5 text-food-primary" /> : <Store className="w-5 h-5 text-food-primary" />}
                      {isDelivery ? 'Delivery Info' : 'Pickup Details'}
                    </CardTitle>
                  </CardHeader>
                  <CardContent className="space-y-3 text-sm">
                    <div>
                      <p className="font-semibold text-gray-900">{order.customerName}</p>
                      <p className="text-gray-600 flex items-center gap-1.5 mt-1">
                        <Mail className="w-3.5 h-3.5 text-gray-400" />
                        {order.customerEmail}
                      </p>
                      <p className="text-gray-600 flex items-center gap-1.5 mt-1">
                        <Phone className="w-3.5 h-3.5 text-gray-400" />
                        {order.customerPhone}
                      </p>
                    </div>

                    {isDelivery && order.deliveryAddress && (
                      <div className="border-t pt-3 mt-3">
                        <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">
                          Delivery Address
                        </p>
                        <p className="text-gray-800">{order.deliveryAddress.street}</p>
                        <p className="text-gray-600">
                          {order.deliveryAddress.city}{order.deliveryAddress.state ? `, ${order.deliveryAddress.state}` : ''} {order.deliveryAddress.postalCode}
                        </p>
                        {order.deliveryAddress.deliveryInstructions && (
                          <p className="text-xs text-gray-500 italic mt-1 bg-gray-50 p-2 rounded border">
                            "{order.deliveryAddress.deliveryInstructions}"
                          </p>
                        )}
                      </div>
                    )}

                    <div className="border-t pt-3 mt-3">
                      <p className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-1">
                        Payment Method
                      </p>
                      <div className="flex items-center justify-between">
                        <span className="capitalize text-gray-800">
                          {order.paymentMethod === 'creditCard' ? 'Credit Card' : 'Pay at Location'}
                        </span>
                        <Badge variant="outline" className={`capitalize text-[11px] ${order.paymentStatus === 'paid' ? 'bg-green-50 text-green-700 border-green-200' : 'bg-yellow-50 text-yellow-700 border-yellow-200'}`}>
                          {order.paymentStatus}
                        </Badge>
                      </div>
                    </div>
                  </CardContent>
                </Card>

                {/* Back / Re-order Action Card */}
                <Card className="border-gray-200 shadow-sm bg-gray-50">
                  <CardContent className="pt-6 space-y-3">
                    <Button
                      onClick={() => navigate(restaurantId ? `/restaurant/${restaurantId}` : '/')}
                      className="w-full bg-food-primary hover:bg-food-accent text-white font-semibold"
                    >
                      <ArrowLeft className="w-4 h-4 mr-2" />
                      Back to Restaurant Menu
                    </Button>
                  </CardContent>
                </Card>
              </div>
            </div>
          </div>
        )}
      </div>
    </Layout>
  );
};

export default OrderTracking;
