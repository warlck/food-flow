import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useCart } from '@/context/CartContext';
import { ShoppingCart, Menu, X, Book, ChefHat, Package } from 'lucide-react';
import { WHATSAPP_URL } from '@/components/landing/constants';

const Header: React.FC = () => {
  const { getTotalItems, restaurantId } = useCart();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const location = useLocation();
  const isLandingPage = location.pathname === '/';

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  // The landing page hero is a dark, full-bleed surface, so the header switches
  // to a translucent dark treatment there and stays light everywhere else.
  const headerClass = isLandingPage
    ? 'sticky top-0 z-40 border-b border-white/10 bg-ink-950/70 backdrop-blur-xl'
    : 'sticky top-0 z-40 bg-white shadow-md';

  const navLinkClass = isLandingPage
    ? 'text-gray-300 hover:text-white transition-colors flex items-center'
    : 'text-gray-700 hover:text-food-primary transition-colors flex items-center';

  return (
    <header className={headerClass}>
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        {/* Logo */}
        <Link to={restaurantId ? `/restaurant/${restaurantId}` : "/"} className="flex items-center">
          <div className="text-food-primary font-bold text-2xl flex items-center">
            <ChefHat className="mr-2" size={32} />
            <span className={isLandingPage ? 'text-white' : undefined}>FoodFlow</span>
          </div>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex items-center space-x-6">
          {!isLandingPage && restaurantId && (
            <Link to={`/restaurant/${restaurantId}`} className={navLinkClass}>
              <Book className="mr-1" size={18} />
              <span>Menu</span>
            </Link>
          )}
          <Link to="/track-order" className={navLinkClass}>
            <Package className="mr-1" size={18} />
            <span>Track Order</span>
          </Link>
          {isLandingPage && (
            <a
              href={WHATSAPP_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-lg bg-food-primary px-4 py-2 text-sm font-bold text-white transition-colors hover:bg-orange-600"
            >
              Partner with us
            </a>
          )}
        </nav>

        {/* Cart Button */}
        {!isLandingPage && (
          <Link to="/cart">
            <Button variant="outline" className="ml-4 relative">
              <ShoppingCart className="h-5 w-5" />
              {getTotalItems() > 0 && (
                <span className="absolute -top-2 -right-2 bg-food-primary text-white text-xs rounded-full h-5 w-5 flex items-center justify-center">
                  {getTotalItems()}
                </span>
              )}
            </Button>
          </Link>
        )}

        {/* Mobile Menu Button */}
        <Button
          variant="ghost"
          size="icon"
          className={isLandingPage ? 'md:hidden text-white hover:bg-white/10 hover:text-white' : 'md:hidden'}
          onClick={toggleMobileMenu}
        >
          {mobileMenuOpen ? <X /> : <Menu />}
        </Button>
      </div>

      {/* Mobile Navigation */}
      {mobileMenuOpen && (
        <div
          className={
            isLandingPage
              ? 'md:hidden border-t border-white/10 bg-ink-950/95 backdrop-blur-xl py-4 px-4 animate-fade-in'
              : 'md:hidden bg-white border-t py-4 px-4 shadow-lg animate-fade-in'
          }
        >
          <nav className="flex flex-col space-y-4">
            {!isLandingPage && restaurantId && (
              <Link 
                to={`/mobile-restaurant/${restaurantId}`} 
                className={`${navLinkClass} p-2`}
                onClick={() => setMobileMenuOpen(false)}
              >
                <Book className="mr-2" size={18} />
                <span>Menu</span>
              </Link>
            )}
            <Link 
              to="/track-order" 
              className={`${navLinkClass} p-2`}
              onClick={() => setMobileMenuOpen(false)}
            >
              <Package className="mr-2" size={18} />
              <span>Track Order</span>
            </Link>
            {isLandingPage && (
              <a
                href={WHATSAPP_URL}
                target="_blank"
                rel="noopener noreferrer"
                className="rounded-lg bg-food-primary px-4 py-3 text-center text-sm font-bold text-white transition-colors hover:bg-orange-600"
                onClick={() => setMobileMenuOpen(false)}
              >
                Partner with us
              </a>
            )}
          </nav>
        </div>
      )}
    </header>
  );
};

export default Header;
