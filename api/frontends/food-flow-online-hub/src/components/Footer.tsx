
import React from 'react';
import { Link } from 'react-router-dom';
import { Facebook, Instagram, Twitter } from 'lucide-react';

const Footer: React.FC = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-gray-800 text-white mt-auto">
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          {/* Logo & Description */}
          <div className="col-span-1 md:col-span-1">
            <h3 className="text-food-primary font-bold text-xl mb-4">FoodFlow</h3>
            <p className="text-gray-300">
              Connecting restaurants with hungry customers through a seamless food delivery experience.
            </p>
            <div className="flex space-x-4 mt-4">
              <a href="#" className="text-gray-300 hover:text-food-primary">
                <Facebook className="h-5 w-5" />
              </a>
              <a href="#" className="text-gray-300 hover:text-food-primary">
                <Instagram className="h-5 w-5" />
              </a>
              <a href="#" className="text-gray-300 hover:text-food-primary">
                <Twitter className="h-5 w-5" />
              </a>
            </div>
          </div>

          {/* For Restaurants */}
          <div className="col-span-1">
            <h3 className="font-semibold text-lg mb-4">For Restaurants</h3>
            <ul className="space-y-2">
              <li>
                <Link to="/restaurant/join" className="text-gray-300 hover:text-food-primary">
                  Join FoodFlow
                </Link>
              </li>
              <li>
                <Link to="/restaurant/login" className="text-gray-300 hover:text-food-primary">
                  Restaurant Login
                </Link>
              </li>
              <li>
                <Link to="/restaurant/dashboard" className="text-gray-300 hover:text-food-primary">
                  Restaurant Dashboard
                </Link>
              </li>
              <li>
                <Link to="/restaurant/resources" className="text-gray-300 hover:text-food-primary">
                  Resources
                </Link>
              </li>
            </ul>
          </div>

          {/* For Customers */}
          <div className="col-span-1">
            <h3 className="font-semibold text-lg mb-4">For Customers</h3>
            <ul className="space-y-2">
              <li>
                <Link to="/menu" className="text-gray-300 hover:text-food-primary">
                  Browse Menu
                </Link>
              </li>
              <li>
                <Link to="/login" className="text-gray-300 hover:text-food-primary">
                  Login
                </Link>
              </li>
              <li>
                <Link to="/signup" className="text-gray-300 hover:text-food-primary">
                  Sign Up
                </Link>
              </li>
              <li>
                <Link to="/track-order" className="text-gray-300 hover:text-food-primary">
                  Track Order
                </Link>
              </li>
            </ul>
          </div>

          {/* Contact & Support */}
          <div className="col-span-1">
            <h3 className="font-semibold text-lg mb-4">Contact & Support</h3>
            <ul className="space-y-2">
              <li>
                <Link to="/contact" className="text-gray-300 hover:text-food-primary">
                  Contact Us
                </Link>
              </li>
              <li>
                <Link to="/faq" className="text-gray-300 hover:text-food-primary">
                  FAQs
                </Link>
              </li>
              <li>
                <Link to="/privacy" className="text-gray-300 hover:text-food-primary">
                  Privacy Policy
                </Link>
              </li>
              <li>
                <Link to="/terms" className="text-gray-300 hover:text-food-primary">
                  Terms of Service
                </Link>
              </li>
            </ul>
          </div>
        </div>

        <div className="border-t border-gray-700 mt-8 pt-6 text-center">
          <p className="text-gray-400 text-sm">
            &copy; {currentYear} FoodFlow. All rights reserved.
          </p>
        </div>
      </div>
    </footer>
  );
};

export default Footer;
