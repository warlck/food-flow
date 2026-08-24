import React from 'react';
import Layout from '@/components/Layout';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Phone, MessageSquare, Mail, Clock, Info, ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';
import { useCart } from '@/context/CartContext';

const Contact: React.FC = () => {
  const { restaurantId } = useCart();

  return (
    <Layout>
      <div className="container mx-auto px-4 py-12 max-w-4xl">
        {/* Back Link */}
        <div className="mb-6">
          <Link
            to={restaurantId ? `/restaurant/${restaurantId}` : '/'}
            className="inline-flex items-center text-sm font-medium text-gray-600 hover:text-food-primary transition-colors"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Menu
          </Link>
        </div>

        {/* Hero Header */}
        <div className="text-center mb-10">
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 tracking-tight">Contact Us</h1>
          <p className="mt-3 text-lg text-gray-600 max-w-2xl mx-auto">
            Have questions about an order, restaurant partnership, or feedback? Reach out directly via WhatsApp, phone, or email.
          </p>
        </div>

        {/* Direct Action Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
          {/* WhatsApp Card */}
          <Card className="hover:shadow-md transition-shadow border-green-200 bg-green-50/30 flex flex-col justify-between">
            <CardHeader className="text-center pb-2">
              <div className="mx-auto w-12 h-12 rounded-full bg-green-100 flex items-center justify-center text-green-600 mb-3">
                <MessageSquare className="h-6 w-6" />
              </div>
              <CardTitle className="text-lg">WhatsApp Chat</CardTitle>
              <CardDescription>Instant messaging & support</CardDescription>
            </CardHeader>
            <CardContent className="text-center pt-2">
              <p className="text-sm font-medium text-gray-900 mb-4">+65 8371 5877</p>
              <Button asChild className="w-full bg-green-600 hover:bg-green-700 text-white">
                <a href="https://wa.me/6583715877" target="_blank" rel="noopener noreferrer">
                  Chat on WhatsApp
                </a>
              </Button>
            </CardContent>
          </Card>

          {/* Phone Call Card */}
          <Card className="hover:shadow-md transition-shadow border-blue-200 bg-blue-50/30 flex flex-col justify-between">
            <CardHeader className="text-center pb-2">
              <div className="mx-auto w-12 h-12 rounded-full bg-blue-100 flex items-center justify-center text-blue-600 mb-3">
                <Phone className="h-6 w-6" />
              </div>
              <CardTitle className="text-lg">Phone Call</CardTitle>
              <CardDescription>Direct line assistance</CardDescription>
            </CardHeader>
            <CardContent className="text-center pt-2">
              <p className="text-sm font-medium text-gray-900 mb-4">+65 8371 5877</p>
              <Button asChild variant="outline" className="w-full border-blue-600 text-blue-700 hover:bg-blue-50">
                <a href="tel:+6583715877">Call Now</a>
              </Button>
            </CardContent>
          </Card>

          {/* Email Card */}
          <Card className="hover:shadow-md transition-shadow border-amber-200 bg-amber-50/30 flex flex-col justify-between">
            <CardHeader className="text-center pb-2">
              <div className="mx-auto w-12 h-12 rounded-full bg-amber-100 flex items-center justify-center text-amber-600 mb-3">
                <Mail className="h-6 w-6" />
              </div>
              <CardTitle className="text-lg">Email Us</CardTitle>
              <CardDescription>For inquiries & partnerships</CardDescription>
            </CardHeader>
            <CardContent className="text-center pt-2">
              <p className="text-xs sm:text-sm font-medium text-gray-900 mb-4 truncate">adil@codercrafters.com</p>
              <Button asChild variant="outline" className="w-full border-amber-600 text-amber-800 hover:bg-amber-50">
                <a href="mailto:adil@codercrafters.com?subject=FoodFlow%20Inquiry">Send Email</a>
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* Service Hours & Platform Role Notice */}
        <div className="space-y-4">
          <div className="flex items-center justify-center p-4 bg-gray-50 rounded-xl border text-sm text-gray-700">
            <Clock className="h-5 w-5 text-gray-500 mr-2 shrink-0" />
            <span>
              <strong>Support Hours:</strong> Monday – Sunday: 10:00 AM – 11:00 PM (Singapore Time / SGT)
            </span>
          </div>

          <div className="p-5 bg-amber-50/70 rounded-xl border border-amber-200 text-sm text-amber-900 flex items-start">
            <Info className="h-5 w-5 text-amber-600 mr-3 mt-0.5 shrink-0" />
            <div>
              <strong className="font-semibold block mb-1">Platform Notice</strong>
              FoodFlow is a technology platform connecting customers to independent partner restaurants in Singapore. Food preparation, packaging, pricing, and delivery fulfillment are managed directly by the individual partner restaurants. For urgent order modifications, you can also contact the restaurant directly via the contact number listed on your order confirmation.
            </div>
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default Contact;
