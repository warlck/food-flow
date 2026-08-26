import React, { useState } from 'react';
import MarketingPage from '@/components/layout/MarketingPage';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from '@/components/ui/accordion';
import {
  HelpCircle,
  Package,
  MessageSquare,
  Mail,
  Search,
  ChevronRight,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { Input } from '@/components/ui/input';

interface FAQItem {
  id: string;
  category: string;
  question: string;
  answer: string;
}

const faqs: FAQItem[] = [
  {
    id: 'how-it-works',
    category: 'General',
    question: 'How does FoodFlow work in Singapore?',
    answer:
      'FoodFlow is an online food ordering platform that connects hungry diners with independent partner restaurants across Singapore. You can easily browse digital menus, customize items with add-ons, and place orders for convenient self-pickup or doorstep delivery.',
  },
  {
    id: 'food-delivery-responsibility',
    category: 'Delivery & Quality',
    question: 'Who is responsible for food preparation and delivery?',
    answer:
      'All food preparation, cooking, kitchen hygiene, packaging, and delivery fulfillment are managed directly by each independent partner restaurant. FoodFlow provides the ordering technology and tracking interface connecting you to the restaurant.',
  },
  {
    id: 'track-order',
    category: 'Orders',
    question: 'How do I track my active order?',
    answer:
      'After placing your order, you are provided with a unique Order ID. You can visit the Track Order page anytime, enter your Order ID or phone number, and see live status updates from order confirmation to kitchen prep and delivery.',
  },
  {
    id: 'cancellation-policy',
    category: 'Orders',
    question: 'How do I modify or cancel an order?',
    answer:
      'Because restaurants often begin preparing fresh food immediately, please contact us via WhatsApp at +65 8371 5877 or call the restaurant directly as soon as possible with your Order ID.',
  },
  {
    id: 'payment-methods',
    category: 'Payments',
    question: 'What payment methods are accepted?',
    answer:
      'We support major credit/debit cards (Visa, Mastercard, American Express) and PayNow through secure, encrypted Stripe checkout. All prices are listed in Singapore Dollars (SGD).',
  },
  {
    id: 'restaurant-partnership',
    category: 'Partnerships',
    question: 'How can my restaurant partner with FoodFlow?',
    answer:
      'We welcome new restaurant partners in Singapore! Contact us directly at adil@codercrafters.com or message +65 8371 5877 on WhatsApp to discuss onboarding and platform setup.',
  },
];

const Support: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');

  const filteredFaqs = faqs.filter(
    (faq) =>
      faq.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
      faq.answer.toLowerCase().includes(searchQuery.toLowerCase()) ||
      faq.category.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <MarketingPage
      title="Help & Support"
      description="Find answers to frequently asked questions, track your orders, or reach out to our team in Singapore."
      icon={HelpCircle}
    >
      {/* Quick Action Banner */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-10">
        <Card className="glass-panel ring-lit border-white/10 hover:border-food-primary/40 transition-colors text-white">
          <CardHeader className="p-5 pb-3">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 rounded-lg bg-food-primary/15 flex items-center justify-center text-food-primary">
                <Package className="h-5 w-5" />
              </div>
              <CardTitle className="text-base text-white">Track Order</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="p-5 pt-0">
            <p className="text-xs text-gray-400 mb-3">Check live delivery progress of your food.</p>
            <Button asChild size="sm" className="w-full bg-food-primary hover:bg-food-primary/90 text-white font-semibold">
              <Link to="/track-order">
                Track Now <ChevronRight className="h-3 w-3 ml-1" />
              </Link>
            </Button>
          </CardContent>
        </Card>

        <Card className="glass-panel ring-lit border-white/10 hover:border-green-500/40 transition-colors text-white">
          <CardHeader className="p-5 pb-3">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 rounded-lg bg-green-500/15 flex items-center justify-center text-green-400">
                <MessageSquare className="h-5 w-5" />
              </div>
              <CardTitle className="text-base text-white">WhatsApp</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="p-5 pt-0">
            <p className="text-xs text-gray-400 mb-3">Quick chat with our support team.</p>
            <Button asChild size="sm" className="w-full bg-green-600 hover:bg-green-500 text-white font-semibold">
              <a href="https://wa.me/6583715877" target="_blank" rel="noopener noreferrer">
                Chat Now <ChevronRight className="h-3 w-3 ml-1" />
              </a>
            </Button>
          </CardContent>
        </Card>

        <Card className="glass-panel ring-lit border-white/10 hover:border-blue-500/40 transition-colors text-white">
          <CardHeader className="p-5 pb-3">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 rounded-lg bg-blue-500/15 flex items-center justify-center text-blue-400">
                <Mail className="h-5 w-5" />
              </div>
              <CardTitle className="text-base text-white">Contact Us</CardTitle>
            </div>
          </CardHeader>
          <CardContent className="p-5 pt-0">
            <p className="text-xs text-gray-400 mb-3">View all support & direct channels.</p>
            <Button asChild size="sm" variant="outline" className="w-full border-white/20 bg-white/5 text-white hover:bg-white/15 hover:text-white">
              <Link to="/contact">
                Contact Page <ChevronRight className="h-3 w-3 ml-1" />
              </Link>
            </Button>
          </CardContent>
        </Card>
      </div>

      {/* Search Input */}
      <div className="relative mb-8">
        <Search className="absolute left-3.5 top-3.5 h-4 w-4 text-gray-400" />
        <Input
          type="text"
          placeholder="Search FAQs (e.g. delivery, payments, cancel, tracking)..."
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="pl-10 h-12 glass-panel border-white/10 text-white placeholder:text-gray-500 focus-visible:ring-food-primary"
        />
      </div>

      {/* FAQ Accordion List */}
      <div className="glass-panel ring-lit rounded-2xl p-6 sm:p-8 mb-10 border-white/10">
        <div className="flex items-center space-x-2.5 mb-6 pb-4 border-b border-white/10">
          <HelpCircle className="h-5 w-5 text-food-primary" />
          <h2 className="text-lg font-bold text-white">Frequently Asked Questions</h2>
        </div>

        {filteredFaqs.length > 0 ? (
          <Accordion type="single" collapsible className="w-full">
            {filteredFaqs.map((faq) => (
              <AccordionItem key={faq.id} value={faq.id} className="border-b border-white/10">
                <AccordionTrigger className="text-left font-medium text-white hover:text-food-primary py-4 hover:no-underline">
                  {faq.question}
                </AccordionTrigger>
                <AccordionContent className="text-gray-300 leading-relaxed text-sm pb-4">
                  {faq.answer}
                </AccordionContent>
              </AccordionItem>
            ))}
          </Accordion>
        ) : (
          <div className="text-center py-8 text-gray-400 text-sm">
            No matching FAQs found for &ldquo;{searchQuery}&rdquo;. Feel free to{' '}
            <Link to="/contact" className="text-food-primary font-medium underline hover:text-orange-400">
              contact us
            </Link>{' '}
            directly!
          </div>
        )}
      </div>
    </MarketingPage>
  );
};

export default Support;
