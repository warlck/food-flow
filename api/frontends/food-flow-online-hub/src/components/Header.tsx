import React, { useState, useEffect } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useCart } from '@/context/CartContext';
import { useSurface } from '@/hooks/useSurface';
import { useRestaurantContext } from '@/hooks/useRestaurantContext';
import { type Surface } from '@/lib/surfaces';
import { ShoppingCart, Menu, X, Book, ChefHat, Package } from 'lucide-react';
import { WHATSAPP_URL } from '@/components/landing/constants';

interface HeaderProps {
  surface?: Surface;
}

const Header: React.FC<HeaderProps> = ({ surface: propSurface }) => {
  const { getTotalItems } = useCart();
  const { restaurantId } = useRestaurantContext();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);
  const [scrolled, setScrolled] = useState(false);
  const location = useLocation();
  const hookSurface = useSurface();
  const surface = propSurface ?? hookSurface;
  const isMarketing = surface === 'marketing';
  const isLandingPage = location.pathname === '/';

  const isCurrentRestaurantPage = Boolean(
    restaurantId &&
      (location.pathname === `/restaurant/${restaurantId}` ||
        location.pathname === `/mobile-restaurant/${restaurantId}`)
  );
  const showMenuLink = Boolean(restaurantId && !isCurrentRestaurantPage);

  useEffect(() => {
    if (!isLandingPage) {
      setScrolled(false);
      return;
    }

    let ticking = false;
    const onScroll = () => {
      if (!ticking) {
        window.requestAnimationFrame(() => {
          setScrolled(window.scrollY > 24);
          ticking = false;
        });
        ticking = true;
      }
    };

    window.addEventListener('scroll', onScroll, { passive: true });
    setScrolled(window.scrollY > 24);

    return () => {
      window.removeEventListener('scroll', onScroll);
    };
  }, [isLandingPage]);

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  let headerClass = 'sticky top-0 z-40 transition-colors duration-300 ';
  if (isMarketing) {
    if (isLandingPage && !scrolled) {
      headerClass += 'bg-transparent border-b border-transparent';
    } else {
      headerClass += 'border-b border-white/10 bg-ink-950/70 backdrop-blur-xl';
    }
  } else {
    headerClass += 'bg-white shadow-md';
  }

  const navLinkClass = isMarketing
    ? 'text-gray-300 hover:text-white transition-colors flex items-center'
    : 'text-gray-700 hover:text-food-primary transition-colors flex items-center';

  return (
    <header className={headerClass}>
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        {/* Logo */}
        <Link to={!isMarketing && restaurantId ? `/restaurant/${restaurantId}` : "/"} className="flex items-center">
          <div className="text-food-primary font-bold text-2xl flex items-center">
            <ChefHat className="mr-2" size={32} />
            <span className={isMarketing ? 'text-white' : undefined}>FoodFlow</span>
          </div>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex items-center space-x-6">
          {showMenuLink && (
            <Link to={`/restaurant/${restaurantId}`} className={navLinkClass}>
              <Book className="mr-1" size={18} />
              <span>Menu</span>
            </Link>
          )}
          <Link to="/track-order" className={navLinkClass}>
            <Package className="mr-1" size={18} />
            <span>Track Order</span>
          </Link>
          {isMarketing && (
            <a
              href={WHATSAPP_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="rounded-lg bg-food-primary px-4 py-2 text-sm font-bold text-white transition-colors hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-food-primary focus:ring-offset-2 focus:ring-offset-ink-950"
            >
              Partner with us
            </a>
          )}
        </nav>

        {/* Cart Button */}
        {!isMarketing && (
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
          className={isMarketing ? 'md:hidden text-white hover:bg-white/10 hover:text-white' : 'md:hidden'}
          onClick={toggleMobileMenu}
        >
          {mobileMenuOpen ? <X /> : <Menu />}
        </Button>
      </div>

      {/* Mobile Navigation */}
      {mobileMenuOpen && (
        <div
          className={
            isMarketing
              ? 'md:hidden border-t border-white/10 bg-ink-950/95 backdrop-blur-xl py-4 px-4 animate-fade-in'
              : 'md:hidden bg-white border-t py-4 px-4 shadow-lg animate-fade-in'
          }
        >
          <nav className="flex flex-col space-y-4">
            {showMenuLink && (
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
            {isMarketing && (
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
