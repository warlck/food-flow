import React from 'react';
import MarketingPage from '@/components/layout/MarketingPage';
import { Card, CardContent, CardHeader, CardTitle, CardDescription } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Phone, MessageSquare, Mail, Clock, Info } from 'lucide-react';

const Contact: React.FC = () => {
  return (
    <MarketingPage
      title="Contact Us"
      description="Have questions about an order, restaurant partnership, or feedback? Reach out directly via WhatsApp, phone, or email."
      icon={Phone}
    >
      {/* Direct Action Grid */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-10">
        {/* WhatsApp Card */}
        <Card className="glass-panel ring-lit hover:border-green-500/30 transition-all flex flex-col justify-between text-white border-white/10">
          <CardHeader className="text-center pb-2">
            <div className="mx-auto w-12 h-12 rounded-full bg-green-500/15 flex items-center justify-center text-green-400 mb-3 ring-1 ring-green-500/20">
              <MessageSquare className="h-6 w-6" />
            </div>
            <CardTitle className="text-lg text-white">WhatsApp Chat</CardTitle>
            <CardDescription className="text-gray-400">Instant messaging & support</CardDescription>
          </CardHeader>
          <CardContent className="text-center pt-2">
            <p className="text-sm font-medium text-white mb-4">+65 8371 5877</p>
            <Button asChild className="w-full bg-green-600 hover:bg-green-500 text-white font-semibold shadow-lg shadow-green-900/20">
              <a href="https://wa.me/6583715877" target="_blank" rel="noopener noreferrer">
                Chat on WhatsApp
              </a>
            </Button>
          </CardContent>
        </Card>

        {/* Phone Call Card */}
        <Card className="glass-panel ring-lit hover:border-blue-500/30 transition-all flex flex-col justify-between text-white border-white/10">
          <CardHeader className="text-center pb-2">
            <div className="mx-auto w-12 h-12 rounded-full bg-blue-500/15 flex items-center justify-center text-blue-400 mb-3 ring-1 ring-blue-500/20">
              <Phone className="h-6 w-6" />
            </div>
            <CardTitle className="text-lg text-white">Phone Call</CardTitle>
            <CardDescription className="text-gray-400">Direct line assistance</CardDescription>
          </CardHeader>
          <CardContent className="text-center pt-2">
            <p className="text-sm font-medium text-white mb-4">+65 8371 5877</p>
            <Button asChild variant="outline" className="w-full border-white/20 bg-white/5 text-white hover:bg-white/15 hover:text-white">
              <a href="tel:+6583715877">Call Now</a>
            </Button>
          </CardContent>
        </Card>

        {/* Email Card */}
        <Card className="glass-panel ring-lit hover:border-amber-500/30 transition-all flex flex-col justify-between text-white border-white/10">
          <CardHeader className="text-center pb-2">
            <div className="mx-auto w-12 h-12 rounded-full bg-amber-500/15 flex items-center justify-center text-amber-400 mb-3 ring-1 ring-amber-500/20">
              <Mail className="h-6 w-6" />
            </div>
            <CardTitle className="text-lg text-white">Email Us</CardTitle>
            <CardDescription className="text-gray-400">For inquiries & partnerships</CardDescription>
          </CardHeader>
          <CardContent className="text-center pt-2">
            <p className="text-xs sm:text-sm font-medium text-white mb-4 truncate">adil@codercrafters.com</p>
            <Button asChild variant="outline" className="w-full border-white/20 bg-white/5 text-white hover:bg-white/15 hover:text-white">
              <a href="mailto:adil@codercrafters.com?subject=FoodFlow%20Inquiry">Send Email</a>
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Service Hours & Platform Role Notice */}
      <div className="space-y-4">
        <div className="flex items-center justify-center p-4 glass-panel rounded-2xl text-sm text-gray-300">
          <Clock className="h-5 w-5 text-food-primary mr-2.5 shrink-0" />
          <span>
            <strong className="text-white">Support Hours:</strong> Monday – Sunday: 10:00 AM – 11:00 PM (Singapore Time / SGT)
          </span>
        </div>

        <div className="p-5 glass-panel rounded-2xl border-amber-500/20 bg-amber-500/5 text-sm text-gray-300 flex items-start">
          <Info className="h-5 w-5 text-amber-400 mr-3 mt-0.5 shrink-0" />
          <div>
            <strong className="font-semibold text-white block mb-1">Platform Notice</strong>
            FoodFlow is a technology platform connecting customers to independent partner restaurants in Singapore. Food preparation, packaging, pricing, and delivery fulfillment are managed directly by the individual partner restaurants. For urgent order modifications, you can also contact the restaurant directly via the contact number listed on your order confirmation.
          </div>
        </div>
      </div>
    </MarketingPage>
  );
};

export default Contact;
