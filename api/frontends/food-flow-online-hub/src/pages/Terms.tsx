import React from 'react';
import Layout from '@/components/Layout';
import { ArrowLeft, FileText } from 'lucide-react';
import { Link } from 'react-router-dom';

const Terms: React.FC = () => {
  return (
    <Layout>
      <div className="container mx-auto px-4 py-12 max-w-3xl">
        <div className="mb-6">
          <Link
            to="/"
            className="inline-flex items-center text-sm font-medium text-gray-600 hover:text-food-primary transition-colors"
          >
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Home
          </Link>
        </div>

        <div className="flex items-center space-x-3 mb-6">
          <div className="w-10 h-10 rounded-full bg-food-primary/10 flex items-center justify-center text-food-primary">
            <FileText className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Terms of Service</h1>
            <p className="text-xs text-gray-500">Last updated: August 2026</p>
          </div>
        </div>

        <div className="bg-white rounded-xl border p-6 sm:p-8 shadow-sm space-y-6 text-gray-700 text-sm leading-relaxed">
          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">1. Agreement to Terms</h2>
            <p>
              By accessing or using the FoodFlow platform, you agree to be bound by these Terms of Service. These terms are governed by the laws of the Republic of Singapore.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">2. Role of FoodFlow</h2>
            <p>
              FoodFlow is a technology software platform that connects customers with independent restaurant merchants. <strong>FoodFlow is not a food producer, kitchen operator, or delivery fleet operator.</strong> All food preparation, ingredient sourcing, kitchen hygiene, allergen handling, pricing, and physical delivery fulfillment are the sole responsibility of the independent partner restaurant.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">3. Orders and Payment</h2>
            <p>
              When placing an order, you agree to provide accurate contact and delivery information. All charges are in Singapore Dollars (SGD) inclusive of applicable taxes, delivery fees, and discounts shown prior to checkout.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">4. Cancellations and Modifications</h2>
            <p>
              Because restaurants start cooking food promptly upon receiving your order, cancellations or changes must be requested immediately via WhatsApp (+65 8371 5877) or by calling the restaurant directly. Orders already prepared or out for delivery may not be eligible for refund.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">5. Contact Information</h2>
            <p>
              For questions regarding these terms, please contact:
            </p>
            <div className="mt-2 p-3 bg-gray-50 rounded-lg text-xs font-mono text-gray-800">
              Email: adil@codercrafters.com | WhatsApp: +65 8371 5877
            </div>
          </section>
        </div>
      </div>
    </Layout>
  );
};

export default Terms;
