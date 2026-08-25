import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { CartProvider } from "@/context/CartContext";
import { RestaurantContextProvider } from "@/context/RestaurantContextProvider";

import { ScrollToTop } from "@/components/ScrollToTop";

import Landing from "./pages/Landing";
import Restaurant from "./pages/Restaurant";
import NotFound from "./pages/NotFound";
import MobileRestaurant from "./pages/MobileRestaurant";
import CartMobile from "./pages/CartMobile";
import CartDesktop from "./pages/CartDesktop";
import CheckoutMobile from "./pages/CheckoutMobile";
import CheckoutDesktop from "./pages/CheckoutDesktop";
import OrderConfirmation from "./pages/OrderConfirmation";
import OrderTracking from "./pages/OrderTracking";
import Contact from "./pages/Contact";
import Support from "./pages/Support";
import Privacy from "./pages/Privacy";
import Terms from "./pages/Terms";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <CartProvider>
        <RestaurantContextProvider>
          <Toaster />
          <Sonner />
          <BrowserRouter>
            <ScrollToTop />
            <Routes>
              <Route path="/" element={<Landing />} />
              <Route path="/restaurant/:restaurantId" element={<Restaurant />} />
              <Route path="/mobile-restaurant/:restaurantId" element={<MobileRestaurant />} />
              <Route path="/cart" element={<CartDesktop />} />
              <Route path="/mobile-cart" element={<CartMobile />} />
              <Route path="/checkout" element={<CheckoutDesktop />} />
              <Route path="/mobile-checkout" element={<CheckoutMobile />} />
              <Route path="/order-confirmation/:orderId" element={<OrderConfirmation />} />
              <Route path="/track-order" element={<OrderTracking />} />
              <Route path="/track-order/:orderId" element={<OrderTracking />} />
              <Route path="/order-tracking" element={<OrderTracking />} />
              <Route path="/order-tracking/:orderId" element={<OrderTracking />} />
              <Route path="/contact" element={<Contact />} />
              <Route path="/support" element={<Support />} />
              <Route path="/faq" element={<Support />} />
              <Route path="/privacy" element={<Privacy />} />
              <Route path="/terms" element={<Terms />} />
              <Route path="*" element={<NotFound />} />
            </Routes>
          </BrowserRouter>
        </RestaurantContextProvider>
      </CartProvider>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;

