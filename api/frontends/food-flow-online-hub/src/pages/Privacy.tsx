import React from 'react';
import MarketingPage from '@/components/layout/MarketingPage';
import { Shield } from 'lucide-react';

const Privacy: React.FC = () => {
  return (
    <MarketingPage
      title="Privacy Policy"
      description="Last updated: August 2026"
      icon={Shield}
      maxWidth="max-w-3xl"
    >
      <div className="glass-panel ring-lit rounded-2xl p-6 sm:p-8 border-white/10 space-y-6 text-gray-300 text-sm leading-relaxed">
        <section>
          <h2 className="text-base font-semibold text-white mb-2">1. Overview & Singapore PDPA Compliance</h2>
          <p>
            FoodFlow (&ldquo;we&rdquo;, &ldquo;our&rdquo;, or &ldquo;us&rdquo;) respects your privacy and is committed to protecting your personal data in accordance with the Singapore Personal Data Protection Act 2012 (PDPA). This policy explains how we collect, use, and protect your information when you use our online restaurant ordering platform.
          </p>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">2. Information We Collect</h2>
          <p className="mb-2">To process your food orders and facilitate fulfillment, we collect:</p>
          <ul className="list-disc pl-5 space-y-1 text-gray-300">
            <li><strong className="text-white">Contact Details:</strong> Customer name, contact telephone number, and email address.</li>
            <li><strong className="text-white">Delivery Information:</strong> Delivery street address and geocoordinates (for delivery fulfillment).</li>
            <li><strong className="text-white">Order History:</strong> Ordered items, add-ons, order timestamps, and restaurant preferences.</li>
          </ul>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">3. Payment Security</h2>
          <p>
            All online payments are securely processed through Stripe. FoodFlow does not store your full credit/debit card numbers on our servers. Payment credentials are tokenized directly with PCI-DSS compliant payment gateways.
          </p>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">4. Sharing with Partner Restaurants</h2>
          <p>
            To fulfill your orders, your contact details, delivery address, and order selections are shared with the specific partner restaurant preparing and delivering your food. Partner restaurants are independently responsible for managing order fulfillment.
          </p>
        </section>

        <section>
          <h2 className="text-base font-semibold text-white mb-2">5. Contact Our Team</h2>
          <p>
            If you have any questions regarding your personal data or wish to exercise your rights under the PDPA, please contact us at:
          </p>
          <div className="mt-3 p-3.5 glass-panel rounded-xl text-xs font-mono text-gray-200 border-white/10">
            Email: adil@codercrafters.com | Phone: +65 8371 5877
          </div>
        </section>
      </div>
    </MarketingPage>
  );
};

export default Privacy;
