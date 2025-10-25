
import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import Layout from '@/components/Layout';
import { ChevronRight, ShoppingBag, Store, Clock, Star } from 'lucide-react';
import { mockMenuItems, mockRestaurant } from '@/data/mockData';

const Home: React.FC = () => {
  const featuredItems = mockMenuItems.slice(0, 3);

  return (
    <Layout>
      {/* Hero Section */}
      <section 
        className="relative bg-cover bg-center h-[500px]" 
        style={{ 
          backgroundImage: `url(${mockRestaurant.coverImage})`,
          backgroundPosition: 'center'
        }}
      >
        <div className="absolute inset-0 bg-black/60"></div>
        <div className="container mx-auto px-4 h-full flex flex-col justify-center relative z-10">
          <div className="max-w-xl">
            <h1 className="text-4xl md:text-5xl font-bold text-white mb-4 animate-fade-in">
              Fresh Food Delivered to Your Doorstep
            </h1>
            <p className="text-xl text-white/90 mb-6 animate-slide-in">
              Order from {mockRestaurant.name} and enjoy delicious meals without leaving your home.
            </p>
            <div className="flex flex-col sm:flex-row gap-4 animate-slide-in" style={{animationDelay: '0.2s'}}>
              <Link to="/menu">
                <Button size="lg" className="bg-food-primary hover:bg-food-accent w-full">
                  Order Now <ShoppingBag className="ml-2 h-5 w-5" />
                </Button>
              </Link>
              <Link to="/restaurant/dashboard">
                <Button variant="outline" size="lg" className="bg-white text-food-primary border-food-primary hover:bg-food-primary/10 w-full">
                  Restaurant Login <Store className="ml-2 h-5 w-5" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </section>

      {/* Featured Menu Items */}
      <section className="py-16 bg-food-light">
        <div className="container mx-auto px-4">
          <div className="flex justify-between items-center mb-8">
            <h2 className="text-3xl font-bold">Featured Menu Items</h2>
            <Link to="/menu" className="text-food-primary hover:text-food-accent flex items-center">
              View Full Menu <ChevronRight size={16} />
            </Link>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
            {featuredItems.map((item) => (
              <div key={item.id} className="bg-white rounded-lg shadow-md overflow-hidden transition-transform hover:scale-105 hover:shadow-lg">
                <img
                  src={item.image}
                  alt={item.name}
                  className="w-full h-48 object-cover"
                />
                <div className="p-4">
                  <div className="flex justify-between items-start mb-2">
                    <h3 className="font-semibold text-xl">{item.name}</h3>
                    <span className="font-bold text-food-primary">${item.price.toFixed(2)}</span>
                  </div>
                  <p className="text-gray-600 mb-4">{item.description}</p>
                  <div className="flex justify-between items-center">
                    <div className="flex items-center text-sm text-gray-500">
                      <Clock size={14} className="mr-1" /> {item.preparationTime} min
                    </div>
                    <Link to="/menu">
                      <Button className="bg-food-primary hover:bg-food-accent">Order Now</Button>
                    </Link>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Restaurant Info */}
      <section className="py-16">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row items-center gap-8">
            <div className="md:w-1/2">
              <h2 className="text-3xl font-bold mb-4">{mockRestaurant.name}</h2>
              <div className="flex items-center mb-4">
                <Star className="h-5 w-5 text-food-secondary fill-food-secondary" />
                <span className="ml-1 font-medium">{mockRestaurant.rating} Rating</span>
              </div>
              <p className="text-gray-700 mb-6">{mockRestaurant.description}</p>
              <div className="mb-6">
                <h3 className="font-semibold mb-2">Location</h3>
                <p className="text-gray-600">{mockRestaurant.address}</p>
              </div>
              <div className="mb-6">
                <h3 className="font-semibold mb-2">Hours</h3>
                <ul className="text-gray-600">
                  <li>Monday-Thursday: {mockRestaurant.openingHours.monday.open} - {mockRestaurant.openingHours.monday.close}</li>
                  <li>Friday: {mockRestaurant.openingHours.friday.open} - {mockRestaurant.openingHours.friday.close}</li>
                  <li>Saturday: {mockRestaurant.openingHours.saturday.open} - {mockRestaurant.openingHours.saturday.close}</li>
                  <li>Sunday: {mockRestaurant.openingHours.sunday.open} - {mockRestaurant.openingHours.sunday.close}</li>
                </ul>
              </div>
              <Link to="/menu">
                <Button size="lg" className="bg-food-primary hover:bg-food-accent">
                  View Menu & Order
                </Button>
              </Link>
            </div>
            <div className="md:w-1/2">
              <img
                src={mockRestaurant.coverImage}
                alt={mockRestaurant.name}
                className="rounded-lg shadow-lg w-full"
              />
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-16 bg-food-primary text-white">
        <div className="container mx-auto px-4 text-center">
          <h2 className="text-3xl font-bold mb-6">Are You a Restaurant Owner?</h2>
          <p className="text-xl mb-8 max-w-2xl mx-auto">
            Join FoodFlow and start offering your delicious menu items for online delivery today!
          </p>
          <Link to="/restaurant/join">
            <Button size="lg" variant="outline" className="bg-white text-food-primary hover:bg-food-primary/10 border-white hover:border-white">
              Register Your Restaurant
            </Button>
          </Link>
        </div>
      </section>
    </Layout>
  );
};

export default Home;
