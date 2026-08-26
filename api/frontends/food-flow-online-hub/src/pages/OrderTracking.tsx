import React, { useEffect, useState, useCallback } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { orderService, Order } from '@/services/orderService';
import Layout from '@/components/Layout';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/use-toast';
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
  pending: 'bg-amber-500/15 text-amber-300 border-amber-500/30',
  confirmed: 'bg-blue-500/15 text-blue-300 border-blue-500/30',
  preparing: 'bg-orange-500/15 text-orange-300 border-orange-500/30',
  ready: 'bg-emerald-500/15 text-emerald-300 border-emerald-500/30',
  out_for_delivery: 'bg-purple-500/15 text-purple-300 border-purple-500/30',
  completed: 'bg-green-500/15 text-green-300 border-green-500/30',
  cancelled: 'bg-red-500/15 text-red-300 border-red-500/30',
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
        setError('Order not found');
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
      setError(null);
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
    <Layout surface="marketing" className="bg-gray-800 text-white">
      <div className="relative isolate min-h-[calc(100vh-160px)] bg-gray-800 font-sans text-white selection:bg-food-primary/30 animate-fade-in -mt-[65px] pt-[65px] pb-20">
        {/* Background stack matching MarketingPage & FinalCTA */}
        <div aria-hidden="true" className="absolute inset-0 -z-10 overflow-hidden pointer-events-none">
          {/* Warm bloom behind headline and header */}
          <div className="absolute left-1/2 -top-24 h-[36rem] w-[52rem] -translate-x-1/2 rounded-full bg-food-primary/20 blur-[120px]" />
          {/* Fine grid overlay with radial mask */}
          <div className="absolute inset-0 bg-grid-fine mask-radial opacity-40" />
          {/* Fade into the footer tone */}
          <div className="absolute inset-x-0 bottom-0 h-32 bg-gradient-to-b from-transparent to-gray-800" />
        </div>

        <div className="container mx-auto px-4 max-w-4xl pt-8">
          {/* Back to Home Link */}
          <div className="mb-6">
            <Link
              to="/"
              className="inline-flex items-center text-sm font-medium text-gray-400 hover:text-white transition-colors gap-2"
            >
              <ArrowLeft className="h-4 w-4" />
              <span>Back to home</span>
            </Link>
          </div>

          <div className="text-center mb-8">
            <div className="mb-2 text-xs font-bold uppercase tracking-wider text-food-primary">
              Live Status
            </div>
            <h1 className="text-3xl sm:text-4xl font-extrabold tracking-tight text-white mb-2">
              Track your order
            </h1>
            <p className="text-gray-300 text-sm sm:text-base max-w-md mx-auto">
              Follow real-time milestones from kitchen preparation to your doorstep.
            </p>
          </div>

          {/* Search Bar */}
          <form
            onSubmit={handleSearchSubmit}
            className="glass-panel-strong ring-lit rounded-2xl p-2 sm:p-2.5 flex flex-col sm:flex-row items-center gap-2 max-w-2xl mx-auto mb-10"
          >
            <div className="relative flex-1 w-full">
              <Search className="absolute left-3.5 top-1/2 -translate-y-1/2 text-gray-400 w-5 h-5" />
              <Input
                type="text"
                placeholder="Enter Order ID (e.g., 550e8400-e29b-41d4-a716-446655440000)"
                value={inputOrderId}
                onChange={(e) => setInputOrderId(e.target.value)}
                className="pl-11 bg-transparent border-0 text-white placeholder:text-gray-400 h-12 text-base focus-visible:ring-0 focus-visible:ring-offset-0 shadow-none"
              />
            </div>
            <Button
              type="submit"
              className="w-full sm:w-auto h-12 px-6 bg-food-primary hover:bg-orange-600 text-white font-semibold flex items-center justify-center gap-2 rounded-xl shadow-lg shadow-food-primary/25 transition-colors"
            >
              <Search className="w-4 h-4" />
              Track Order
            </Button>
          </form>

          {/* Loading State */}
          {loading && (
            <div className="flex flex-col items-center justify-center py-20">
              <Loader2 className="w-12 h-12 animate-spin text-food-primary mb-4" />
              <p className="text-gray-300 font-medium">Fetching order status...</p>
            </div>
          )}

          {/* Error / Not Found State */}
          {!loading && error && (
            <div className="glass-panel rounded-2xl p-10 text-center border border-red-500/30 bg-red-950/20">
              <AlertTriangle className="w-14 h-14 text-red-400 mx-auto mb-4" />
              <h2 className="text-2xl font-bold text-white mb-2">Order Not Found</h2>
              <p className="text-gray-300 max-w-md mx-auto mb-6">
                We could not find the order you were looking for.
              </p>
              <Button
                onClick={() => navigate('/')}
                className="border border-white/20 bg-white/10 hover:bg-white/20 text-white rounded-xl h-11 px-6 font-semibold transition-colors"
              >
                <ArrowLeft className="w-4 h-4 mr-2" />
                Back to Home
              </Button>
            </div>
          )}

          {/* No Order Selected Prompt */}
          {!loading && !error && !order && !orderId && (
            <div className="glass-panel ring-lit rounded-2xl p-12 text-center border border-white/10 bg-white/[0.03]">
              <Package className="w-16 h-16 text-gray-400 mx-auto mb-4" />
              <h2 className="text-xl font-bold text-white mb-2">Track Your Food Order Live</h2>
              <p className="text-gray-300 max-w-md mx-auto">
                Enter your unique Order ID in the search bar above to see live updates, estimated time, and order details.
              </p>
            </div>
          )}

          {/* Order Tracking Dashboard */}
          {!loading && !error && order && (
            <div className="space-y-6 animate-fade-in">
              {/* Live Refresh Control & Header */}
              <div className="glass-panel ring-lit rounded-2xl p-6 border border-white/10 bg-white/[0.04] space-y-3">
                <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
                  <div className="flex items-center gap-2.5 flex-wrap min-h-[36px]">
                    <h1 className="text-2xl font-bold text-white leading-none">
                      Order #{order.id.slice(0, 8).toUpperCase()}
                    </h1>
                    <Badge className={`h-7 px-3 inline-flex items-center text-xs font-semibold rounded-full border ${ORDER_STATUS_COLORS[order.orderStatus] || 'bg-white/10 text-white'}`}>
                      {getOrderStatusLabel(order.orderStatus, isDelivery)}
                    </Badge>
                    <Badge variant="outline" className="h-7 px-3 inline-flex items-center capitalize text-xs font-medium border-white/20 text-gray-300 bg-white/5 rounded-full">
                      {order.orderType}
                    </Badge>
                  </div>

                  <div className="flex items-center gap-2 flex-wrap sm:flex-nowrap min-h-[36px]">
                    {!isCompleted && !isCancelled && (
                      <div className="flex items-center space-x-2 bg-white/5 px-3 h-9 rounded-lg border border-white/10">
                        <Switch
                          id="auto-refresh"
                          checked={autoRefresh}
                          onCheckedChange={setAutoRefresh}
                        />
                        <Label htmlFor="auto-refresh" className="text-xs text-gray-300 flex items-center gap-1.5 cursor-pointer select-none">
                          <span className={`w-2 h-2 rounded-full ${autoRefresh ? 'bg-green-400 animate-pulse' : 'bg-gray-500'}`}></span>
                          Auto-Refresh
                        </Label>
                      </div>
                    )}

                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => fetchOrder(order.id, true)}
                      disabled={isRefreshing}
                      className="flex items-center gap-1.5 text-xs h-9 px-3 border-white/15 bg-white/5 text-white hover:bg-white/10"
                    >
                      <RefreshCw className={`w-3.5 h-3.5 ${isRefreshing ? 'animate-spin' : ''}`} />
                      Refresh
                    </Button>
                    <Button variant="outline" size="icon" onClick={handleShareLink} title="Share Link" className="h-9 w-9 border-white/15 bg-white/5 text-white hover:bg-white/10">
                      <Share2 className="w-4 h-4 text-gray-300" />
                    </Button>
                    <Button variant="outline" size="icon" onClick={handleCopyOrderId} title="Copy ID" className="h-9 w-9 border-white/15 bg-white/5 text-white hover:bg-white/10">
                      <Copy className="w-4 h-4 text-gray-300" />
                    </Button>
                  </div>
                </div>

                <p className="text-xs text-gray-400 flex items-center gap-2 pt-0.5">
                  <span>Placed on {new Date(order.dateCreated).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}</span>
                  {lastRefreshedAt && (
                    <span>• Updated {lastRefreshedAt.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' })}</span>
                  )}
                </p>
              </div>

              {/* Cancelled Banner */}
              {isCancelled && (
                <div className="bg-red-950/40 border border-red-500/30 rounded-2xl p-4 flex items-center gap-3 text-red-200">
                  <AlertTriangle className="w-6 h-6 text-red-400 flex-shrink-0" />
                  <div>
                    <h3 className="font-semibold text-sm text-white">Order Cancelled</h3>
                    <p className="text-xs text-red-300">This order was cancelled. If you have questions, please contact customer support.</p>
                  </div>
                </div>
              )}

              {/* Ready State Info Banner */}
              {!isCancelled && order.orderStatus === 'ready' && (
                <div className="bg-emerald-950/40 border border-emerald-500/30 rounded-2xl p-4 flex items-center gap-3 text-emerald-200">
                  {isDelivery ? (
                    <>
                      <Package className="w-6 h-6 text-emerald-400 flex-shrink-0" />
                      <div>
                        <h3 className="font-semibold text-sm text-white">Your order is ready!</h3>
                        <p className="text-xs text-emerald-300">Food preparation is complete. Your order is waiting to be dispatched to a courier for delivery.</p>
                      </div>
                    </>
                  ) : (
                    <>
                      <Store className="w-6 h-6 text-emerald-400 flex-shrink-0" />
                      <div>
                        <h3 className="font-semibold text-sm text-white">Ready for Pickup!</h3>
                        <p className="text-xs text-emerald-300">Your order is ready to be collected at the restaurant counter.</p>
                      </div>
                    </>
                  )}
                </div>
              )}

              {/* Status Timeline Stepper */}
              {!isCancelled && (
                <div className="glass-panel ring-lit rounded-2xl border border-white/10 bg-white/[0.04] overflow-hidden">
                  <div className="border-b border-white/10 px-6 py-4 bg-white/[0.02] flex items-center justify-between">
                    <span className="text-base sm:text-lg font-semibold text-white flex items-center gap-2">
                      <Clock className="w-5 h-5 text-food-primary" />
                      Order Status
                    </span>
                    {!isCompleted && (
                      <span className="text-xs font-medium text-amber-300 bg-amber-400/15 border border-amber-400/30 px-3 py-1 rounded-full">
                        Estimated fulfillment: ~20-30 mins
                      </span>
                    )}
                  </div>

                  <div className="pt-8 pb-10 px-4 sm:px-8">
                    {/* Stepper Progress Line Container */}
                    <div className="relative flex items-start justify-between max-w-2xl mx-auto">
                      {/* Horizontal background track */}
                      <div
                        className="absolute top-6 -translate-y-1/2 h-1 bg-white/10 z-0"
                        style={{ left: '24px', right: '24px' }}
                      />
                      {/* Active progress line */}
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
                                  ? 'bg-food-primary border-food-primary text-white shadow-lg shadow-food-primary/25'
                                  : 'bg-gray-800 border-white/20 text-gray-400'
                              } ${isCurrent ? 'ring-4 ring-food-primary/30 scale-110' : ''}`}
                            >
                              <IconComponent className="w-5 h-5" />
                            </div>
                            <div className="mt-3 text-center">
                              <p className={`text-xs font-semibold ${isDone ? 'text-white' : 'text-gray-400'}`}>
                                {stepLabel}
                              </p>
                              <p className="text-[10px] text-gray-400 max-w-[90px] hidden sm:block mt-0.5">
                                {stepDesc}
                              </p>
                            </div>
                          </div>
                        );
                      })}
                    </div>
                  </div>
                </div>
              )}

              {/* Main Content Grid: Items & Details */}
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                {/* Left Column: Items Breakdown */}
                <div className="glass-panel ring-lit rounded-2xl border border-white/10 bg-white/[0.04] p-6 space-y-4 md:col-span-2">
                  <div className="flex items-center gap-2 text-lg font-semibold text-white pb-2 border-b border-white/10">
                    <Package className="w-5 h-5 text-food-primary" />
                    Items Ordered ({order.items.reduce((acc, item) => acc + item.quantity, 0)})
                  </div>
                  <div className="divide-y divide-white/10">
                    {order.items.map((item) => (
                      <div key={item.id} className="py-3 flex justify-between items-start first:pt-0 last:pb-0">
                        <div>
                          <p className="font-medium text-white text-sm">
                            {item.quantity}x {item.menuItemName}
                          </p>
                          {item.addons && item.addons.length > 0 && (
                            <div className="ml-6 mt-1 space-y-0.5">
                              {item.addons.map((addon) => (
                                <div key={addon.id} className="flex justify-between text-xs text-gray-400">
                                  <span>+ {addon.addonName} x{addon.quantity}</span>
                                  <span className="text-food-primary">
                                    +${(addon.addonPrice * addon.quantity * item.quantity).toFixed(2)}
                                  </span>
                                </div>
                              ))}
                            </div>
                          )}
                          {item.specialInstructions && (
                            <p className="text-xs text-gray-400 italic mt-1 bg-white/5 p-2 rounded border border-white/10">
                              Note: {item.specialInstructions}
                            </p>
                          )}
                        </div>
                        <span className="font-semibold text-sm text-gray-200">
                          ${(item.menuItemPrice * item.quantity).toFixed(2)}
                        </span>
                      </div>
                    ))}
                  </div>

                  <div className="border-t border-white/10 pt-4 space-y-2 text-sm">
                    <div className="flex justify-between text-gray-300">
                      <span>Subtotal</span>
                      <span>${order.subtotal.toFixed(2)}</span>
                    </div>
                    {order.deliveryFee > 0 && (
                      <div className="flex justify-between text-gray-300">
                        <span>Delivery Fee</span>
                        <span>${order.deliveryFee.toFixed(2)}</span>
                      </div>
                    )}
                    {order.tax > 0 && (
                      <div className="flex justify-between text-gray-300">
                        <span>Tax</span>
                        <span>${order.tax.toFixed(2)}</span>
                      </div>
                    )}
                    <div className="flex justify-between font-bold text-base border-t border-white/10 pt-3 text-white">
                      <span>Total Paid</span>
                      <span className="text-food-primary text-lg">${order.total.toFixed(2)}</span>
                    </div>
                  </div>
                </div>

                {/* Right Column: Customer & Delivery Info */}
                <div className="space-y-6">
                  <div className="glass-panel ring-lit rounded-2xl border border-white/10 bg-white/[0.04] p-6 space-y-3 text-sm">
                    <div className="flex items-center gap-2 text-lg font-semibold text-white pb-2 border-b border-white/10">
                      {isDelivery ? <MapPin className="w-5 h-5 text-food-primary" /> : <Store className="w-5 h-5 text-food-primary" />}
                      {isDelivery ? 'Delivery Info' : 'Pickup Details'}
                    </div>
                    <div>
                      <p className="font-semibold text-white">{order.customerName}</p>
                      <p className="text-gray-300 flex items-center gap-1.5 mt-1 text-xs">
                        <Mail className="w-3.5 h-3.5 text-gray-400" />
                        {order.customerEmail}
                      </p>
                      <p className="text-gray-300 flex items-center gap-1.5 mt-1 text-xs">
                        <Phone className="w-3.5 h-3.5 text-gray-400" />
                        {order.customerPhone}
                      </p>
                    </div>

                    {isDelivery && order.deliveryAddress && (
                      <div className="border-t border-white/10 pt-3 mt-3">
                        <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">
                          Delivery Address
                        </p>
                        <p className="text-white">{order.deliveryAddress.street}</p>
                        <p className="text-gray-300">
                          {order.deliveryAddress.city}{order.deliveryAddress.state ? `, ${order.deliveryAddress.state}` : ''} {order.deliveryAddress.postalCode}
                        </p>
                        {order.deliveryAddress.deliveryInstructions && (
                          <p className="text-xs text-gray-300 italic mt-2 bg-white/5 p-2 rounded border border-white/10">
                            "{order.deliveryAddress.deliveryInstructions}"
                          </p>
                        )}
                      </div>
                    )}

                    <div className="border-t border-white/10 pt-3 mt-3">
                      <p className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-1">
                        Payment Method
                      </p>
                      <div className="flex items-center justify-between">
                        <span className="capitalize text-gray-200">
                          {order.paymentMethod === 'creditCard' ? 'Credit Card' : 'Pay at Location'}
                        </span>
                        <Badge variant="outline" className={`capitalize text-[11px] ${order.paymentStatus === 'paid' ? 'bg-green-500/15 text-green-300 border-green-500/30' : 'bg-yellow-500/15 text-yellow-300 border-yellow-500/30'}`}>
                          {order.paymentStatus}
                        </Badge>
                      </div>
                    </div>
                  </div>

                  {/* Back Action Card */}
                  <div className="glass-panel ring-lit rounded-2xl border border-white/10 bg-white/[0.04] p-6">
                    {order.restaurantId ? (
                      <Button
                        onClick={() => navigate(`/restaurant/${order.restaurantId}`)}
                        className="w-full h-12 bg-food-primary hover:bg-orange-600 text-white font-bold rounded-xl shadow-lg shadow-food-primary/25 transition-colors"
                      >
                        <ArrowLeft className="w-4 h-4 mr-2" />
                        Back to Restaurant Menu
                      </Button>
                    ) : (
                      <Button
                        onClick={() => navigate('/')}
                        variant="outline"
                        className="w-full h-12 border-white/15 bg-white/5 hover:bg-white/10 text-white font-bold rounded-xl transition-colors"
                      >
                        <ArrowLeft className="w-4 h-4 mr-2" />
                        Back to Home
                      </Button>
                    )}
                  </div>
                </div>
              </div>
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default OrderTracking;
