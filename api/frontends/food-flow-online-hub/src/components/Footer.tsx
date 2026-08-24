import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Facebook, Instagram, Twitter, Heart } from 'lucide-react';
import { useCart } from '@/context/CartContext';

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();
  const { restaurantId } = useCart();
  const location = useLocation();
  const isLandingPage = location.pathname === '/';

  return (
    <footer className="bg-gray-800 text-white mt-auto border-t border-gray-700/50">
      <div className="container mx-auto px-4 py-12 max-w-7xl">
        <div className="flex flex-col md:flex-row md:justify-between gap-10 lg:gap-16">
          
          {/* Logo & Description */}
          <div className="md:w-1/3 flex flex-col items-start">
            <h3 className="text-food-primary font-extrabold text-2xl mb-4 tracking-tight">FoodFlow</h3>
            <p className="text-gray-400 leading-relaxed mb-6 max-w-sm">
              Connecting restaurants with hungry customers through a seamless food delivery experience.
            </p>
            <div className="flex space-x-5">
              <a href="#" className="text-gray-400 hover:text-food-primary transition-colors">
                <Facebook className="h-5 w-5" />
              </a>
              <a href="#" className="text-gray-400 hover:text-food-primary transition-colors">
                <Instagram className="h-5 w-5" />
              </a>
              <a href="#" className="text-gray-400 hover:text-food-primary transition-colors">
                <Twitter className="h-5 w-5" />
              </a>
            </div>
          </div>

          {/* For Customers */}
          <div className="md:w-1/3 flex flex-col md:items-center">
            <div className="w-full md:w-auto">
              <h3 className="font-bold text-lg mb-5 text-white">For Customers</h3>
              <ul className="space-y-3">
                {!isLandingPage && restaurantId && (
                  <li>
                    <Link to={`/restaurant/${restaurantId}`} className="text-gray-400 hover:text-white transition-colors">
                      Browse Menu
                    </Link>
                  </li>
                )}
                <li>
                  <Link to="/track-order" className="text-gray-400 hover:text-white transition-colors">
                    Track Order
                  </Link>
                </li>
              </ul>
            </div>
          </div>

          {/* Contact & Support */}
          <div className="md:w-1/3 flex flex-col md:items-end">
            <div className="w-full md:w-auto">
              <h3 className="font-bold text-lg mb-5 text-white">Contact & Support</h3>
              <ul className="space-y-3">
                <li>
                  <Link to="/contact" className="text-gray-400 hover:text-white transition-colors">
                    Contact Us
                  </Link>
                </li>
                <li>
                  <Link to="/faq" className="text-gray-400 hover:text-white transition-colors">
                    FAQs
                  </Link>
                </li>
                <li>
                  <Link to="/privacy" className="text-gray-400 hover:text-white transition-colors">
                    Privacy Policy
                  </Link>
                </li>
                <li>
                  <Link to="/terms" className="text-gray-400 hover:text-white transition-colors">
                    Terms of Service
                  </Link>
                </li>
              </ul>
            </div>
          </div>
          
        </div>

        {/* Divider & Copyright */}
        <div className="border-t border-gray-700/60 mt-16 pt-8 text-center space-y-2 flex flex-col md:flex-row justify-between items-center text-sm">
          <p className="text-gray-500">
            &copy; {currentYear} FoodFlow. All rights reserved.
          </p>
          <p className="text-gray-400 flex items-center gap-1.5 mt-4 md:mt-0">
            <span>Created with</span>
            <Heart className="h-4 w-4 text-red-500 fill-red-500 inline-block mx-0.5" />
            <span>by <strong className="text-white font-semibold">CodeCrafters</strong></span>
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
