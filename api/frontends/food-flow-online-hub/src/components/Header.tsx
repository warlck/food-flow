import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import { useCart } from '@/context/CartContext';
import { ShoppingCart, Menu, X, Home, Book, ChefHat, LogIn } from 'lucide-react';

const Header: React.FC = () => {
  const { getTotalItems } = useCart();
  const [mobileMenuOpen, setMobileMenuOpen] = useState(false);

  const toggleMobileMenu = () => {
    setMobileMenuOpen(!mobileMenuOpen);
  };

  return (
    <header className="sticky top-0 z-40 bg-white shadow-md">
      <div className="container mx-auto px-4 py-4 flex items-center justify-between">
        {/* Logo */}
        <Link to="/" className="flex items-center">
          <div className="text-food-primary font-bold text-2xl flex items-center">
            <ChefHat className="mr-2" size={32} />
            <span>FoodFlow</span>
          </div>
        </Link>

        {/* Desktop Navigation */}
        <nav className="hidden md:flex items-center space-x-6">
          <Link to="/" className="text-gray-700 hover:text-food-primary transition-colors flex items-center">
            <Home className="mr-1" size={18} />
            <span>Home</span>
          </Link>
          <Link to="/menu" className="text-gray-700 hover:text-food-primary transition-colors flex items-center">
            <Book className="mr-1" size={18} />
            <span>Menu</span>
          </Link>
          <Link to="/login" className="text-gray-700 hover:text-food-primary transition-colors flex items-center">
            <LogIn className="mr-1" size={18} />
            <span>Login</span>
          </Link>
        </nav>

        {/* Cart Button */}
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

        {/* Mobile Menu Button */}
        <Button
          variant="ghost"
          size="icon"
          className="md:hidden"
          onClick={toggleMobileMenu}
        >
          {mobileMenuOpen ? <X /> : <Menu />}
        </Button>
      </div>

      {/* Mobile Navigation */}
      {mobileMenuOpen && (
        <div className="md:hidden bg-white border-t py-4 px-4 shadow-lg animate-fade-in">
          <nav className="flex flex-col space-y-4">
            <Link 
              to="/" 
              className="text-gray-700 hover:text-food-primary transition-colors flex items-center p-2"
              onClick={() => setMobileMenuOpen(false)}
            >
              <Home className="mr-2" size={18} />
              <span>Home</span>
            </Link>
            <Link 
              to="/menu" 
              className="text-gray-700 hover:text-food-primary transition-colors flex items-center p-2"
              onClick={() => setMobileMenuOpen(false)}
            >
              <Book className="mr-2" size={18} />
              <span>Menu</span>
            </Link>
            <Link 
              to="/login" 
              className="text-gray-700 hover:text-food-primary transition-colors flex items-center p-2"
              onClick={() => setMobileMenuOpen(false)}
            >
              <LogIn className="mr-2" size={18} />
              <span>Login</span>
            </Link>
          </nav>
        </div>
      )}
    </header>
  );
};

export default Header;
