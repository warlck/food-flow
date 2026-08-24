import React from 'react';
import Layout from '@/components/Layout';
import { ArrowLeft, Shield } from 'lucide-react';
import { Link } from 'react-router-dom';

const Privacy: React.FC = () => {
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
            <Shield className="h-5 w-5" />
          </div>
          <div>
            <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Privacy Policy</h1>
            <p className="text-xs text-gray-500">Last updated: August 2026</p>
          </div>
        </div>

        <div className="bg-white rounded-xl border p-6 sm:p-8 shadow-sm space-y-6 text-gray-700 text-sm leading-relaxed">
          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">1. Overview & Singapore PDPA Compliance</h2>
            <p>
              FoodFlow (&ldquo;we&rdquo;, &ldquo;our&rdquo;, or &ldquo;us&rdquo;) respects your privacy and is committed to protecting your personal data in accordance with the Singapore Personal Data Protection Act 2012 (PDPA). This policy explains how we collect, use, and protect your information when you use our online restaurant ordering platform.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">2. Information We Collect</h2>
            <p className="mb-2">To process your food orders and facilitate fulfillment, we collect:</p>
            <ul className="list-disc pl-5 space-y-1">
              <li><strong>Contact Details:</strong> Customer name, contact telephone number, and email address.</li>
              <li><strong>Delivery Information:</strong> Delivery street address and geocoordinates (for delivery fulfillment).</li>
              <li><strong>Order History:</strong> Ordered items, add-ons, order timestamps, and restaurant preferences.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">3. Payment Security</h2>
            <p>
              All online payments are securely processed through Stripe. FoodFlow does not store your full credit/debit card numbers on our servers. Payment credentials are tokenized directly with PCI-DSS compliant payment gateways.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">4. Sharing with Partner Restaurants</h2>
            <p>
              To fulfill your orders, your contact details, delivery address, and order selections are shared with the specific partner restaurant preparing and delivering your food. Partner restaurants are independently responsible for managing order fulfillment.
            </p>
          </section>

          <section>
            <h2 className="text-base font-semibold text-gray-900 mb-2">5. Contact Our Team</h2>
            <p>
              If you have any questions regarding your personal data or wish to exercise your rights under the PDPA, please contact us at:
            </p>
            <div className="mt-2 p-3 bg-gray-50 rounded-lg text-xs font-mono text-gray-800">
              Email: adil@codercrafters.com | Phone: +65 8371 5877
            </div>
          </section>
        </div>
      </div>
    </Layout>
  );
};

export default Privacy;
