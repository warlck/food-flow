import { Toaster } from "@/components/ui/toaster";
import { Toaster as Sonner } from "@/components/ui/sonner";
import { TooltipProvider } from "@/components/ui/tooltip";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter, Routes, Route } from "react-router-dom";
import { CartProvider } from "@/context/CartContext";

import Restaurant from "./pages/Restaurant";
import Login from "./pages/Login";
import NotFound from "./pages/NotFound";
import MobileRestaurant from "./pages/MobileRestaurant";
import CartMobile from "./pages/CartMobile";
import CartDesktop from "./pages/CartDesktop";
import CheckoutMobile from "./pages/CheckoutMobile";
import CheckoutDesktop from "./pages/CheckoutDesktop";

const queryClient = new QueryClient();

const App = () => (
  <QueryClientProvider client={queryClient}>
    <TooltipProvider>
      <CartProvider>
        <Toaster />
        <Sonner />
        <BrowserRouter>
          <Routes>
            <Route path="/restaurant/:restaurantId" element={<Restaurant />} />
            <Route path="/mobile-restaurant/:restaurantId" element={<MobileRestaurant />} />
            <Route path="/cart" element={<CartDesktop />} />
            <Route path="/mobile-cart" element={<CartMobile />} />
            <Route path="/checkout" element={<CheckoutDesktop />} />
            <Route path="/mobile-checkout" element={<CheckoutMobile />} />
            <Route path="/login" element={<Login />} />
            <Route path="*" element={<NotFound />} />
          </Routes>
        </BrowserRouter>
      </CartProvider>
    </TooltipProvider>
  </QueryClientProvider>
);

export default App;

