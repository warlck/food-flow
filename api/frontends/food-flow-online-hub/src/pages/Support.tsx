import React, { useState } from 'react';
import Layout from '@/components/Layout';
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
  ArrowLeft,
  ChevronRight,
} from 'lucide-react';
import { Link } from 'react-router-dom';
import { useCart } from '@/context/CartContext';
import { Input } from '@/components/ui/input';

interface FAQItem {
  id: string;
  question: string;
  answer: string;
  category: string;
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
  const { restaurantId } = useCart();
  const [searchQuery, setSearchQuery] = useState('');

  const filteredFaqs = faqs.filter(
    (faq) =>
      faq.question.toLowerCase().includes(searchQuery.toLowerCase()) ||
      faq.answer.toLowerCase().includes(searchQuery.toLowerCase()) ||
      faq.category.toLowerCase().includes(searchQuery.toLowerCase())
  );

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
          <h1 className="text-3xl sm:text-4xl font-bold text-gray-900 tracking-tight">Help & Support</h1>
          <p className="mt-3 text-lg text-gray-600 max-w-2xl mx-auto">
            Find answers to frequently asked questions, track your orders, or reach out to our team in Singapore.
          </p>
        </div>

        {/* Quick Action Banner */}
        <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-10">
          <Card className="bg-food-primary/5 border-food-primary/20 hover:border-food-primary transition-colors">
            <CardHeader className="p-5 pb-3">
              <div className="flex items-center space-x-3">
                <Package className="h-6 w-6 text-food-primary" />
                <CardTitle className="text-base">Track Order</CardTitle>
              </div>
            </CardHeader>
            <CardContent className="p-5 pt-0">
              <p className="text-xs text-gray-600 mb-3">Check live delivery progress of your food.</p>
              <Button asChild size="sm" className="w-full bg-food-primary hover:bg-food-primary/90 text-white">
                <Link to="/track-order">
                  Track Now <ChevronRight className="h-3 w-3 ml-1" />
                </Link>
              </Button>
            </CardContent>
          </Card>

          <Card className="bg-green-50/50 border-green-200 hover:border-green-400 transition-colors">
            <CardHeader className="p-5 pb-3">
              <div className="flex items-center space-x-3">
                <MessageSquare className="h-6 w-6 text-green-600" />
                <CardTitle className="text-base">WhatsApp</CardTitle>
              </div>
            </CardHeader>
            <CardContent className="p-5 pt-0">
              <p className="text-xs text-gray-600 mb-3">Quick chat with our support team.</p>
              <Button asChild size="sm" className="w-full bg-green-600 hover:bg-green-700 text-white">
                <a href="https://wa.me/6583715877" target="_blank" rel="noopener noreferrer">
                  Chat Now <ChevronRight className="h-3 w-3 ml-1" />
                </a>
              </Button>
            </CardContent>
          </Card>

          <Card className="bg-blue-50/50 border-blue-200 hover:border-blue-400 transition-colors">
            <CardHeader className="p-5 pb-3">
              <div className="flex items-center space-x-3">
                <Mail className="h-6 w-6 text-blue-600" />
                <CardTitle className="text-base">Contact Us</CardTitle>
              </div>
            </CardHeader>
            <CardContent className="p-5 pt-0">
              <p className="text-xs text-gray-600 mb-3">View all support & direct channels.</p>
              <Button asChild size="sm" variant="outline" className="w-full border-blue-600 text-blue-700 hover:bg-blue-50">
                <Link to="/contact">
                  Contact Page <ChevronRight className="h-3 w-3 ml-1" />
                </Link>
              </Button>
            </CardContent>
          </Card>
        </div>

        {/* Search Input */}
        <div className="relative mb-8">
          <Search className="absolute left-3.5 top-3 h-4 w-4 text-gray-400" />
          <Input
            type="text"
            placeholder="Search FAQs (e.g. delivery, payments, cancel, tracking)..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10 h-11 bg-white"
          />
        </div>

        {/* FAQ Accordion List */}
        <div className="bg-white rounded-xl border p-6 shadow-sm mb-10">
          <div className="flex items-center space-x-2 mb-4 pb-2 border-b">
            <HelpCircle className="h-5 w-5 text-food-primary" />
            <h2 className="text-lg font-semibold text-gray-900">Frequently Asked Questions</h2>
          </div>

          {filteredFaqs.length > 0 ? (
            <Accordion type="single" collapsible className="w-full">
              {filteredFaqs.map((faq) => (
                <AccordionItem key={faq.id} value={faq.id}>
                  <AccordionTrigger className="text-left font-medium text-gray-800 hover:text-food-primary py-4">
                    {faq.question}
                  </AccordionTrigger>
                  <AccordionContent className="text-gray-600 leading-relaxed text-sm">
                    {faq.answer}
                  </AccordionContent>
                </AccordionItem>
              ))}
            </Accordion>
          ) : (
            <div className="text-center py-8 text-gray-500 text-sm">
              No matching FAQs found for &ldquo;{searchQuery}&rdquo;. Feel free to{' '}
              <Link to="/contact" className="text-food-primary font-medium underline">
                contact us
              </Link>{' '}
              directly!
            </div>
          )}
        </div>
      </div>
    </Layout>
  );
};

export default Support;
