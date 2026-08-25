import React from 'react';
import MarketingPage from '@/components/layout/MarketingPage';
import { FileText } from 'lucide-react';

const Terms: React.FC = () => {
  return (
    <MarketingPage
      title="Terms of Service"
      description="Last updated: August 2026"
      icon={FileText}
      maxWidth="max-w-3xl"
    >
      <div className="glass-panel ring-lit rounded-2xl p-6 sm:p-8 border-white/10 space-y-6 text-gray-300 text-sm leading-relaxed">
        <section>
          <h2 className="text-base font-semibold text-white mb-2">1. Agreement to Terms</h2>
          <p>
            By accessing or using the FoodFlow platform, you agree to be bound by these Terms of Service. These terms are governed by the laws of the Republic of Singapore.
          </p>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">2. Role of FoodFlow</h2>
          <p>
            FoodFlow is a technology software platform that connects customers with independent restaurant merchants. <strong className="text-white">FoodFlow is not a food producer, kitchen operator, or delivery fleet operator.</strong> All food preparation, ingredient sourcing, kitchen hygiene, allergen handling, pricing, and physical delivery fulfillment are the sole responsibility of the independent partner restaurant.
          </p>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">3. Orders and Payment</h2>
          <p>
            When placing an order, you agree to provide accurate contact and delivery information. All charges are in Singapore Dollars (SGD) inclusive of applicable taxes, delivery fees, and discounts shown prior to checkout.
          </p>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">4. Cancellations and Modifications</h2>
          <p>
            Because restaurants start cooking food promptly upon receiving your order, cancellations or changes must be requested immediately via WhatsApp (+65 8371 5877) or by calling the restaurant directly. Orders already prepared or out for delivery may not be eligible for refund.
          </p>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">5. Contact Information</h2>
          <p>
            For questions regarding these terms, please contact:
          </p>
          <div className="mt-3 p-3.5 glass-panel rounded-xl text-xs font-mono text-gray-200 border-white/10">
            Email: adil@codercrafters.com | WhatsApp: +65 8371 5877
          </div>
        </section>
      </div>
    </MarketingPage>
  );
};

export default Terms;
