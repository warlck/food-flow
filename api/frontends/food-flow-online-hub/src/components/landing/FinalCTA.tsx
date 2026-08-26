import React from 'react';
import { Button } from '@/components/ui/button';
import { ArrowRight, Mail, MessageSquare } from 'lucide-react';
import { SALES_MAILTO, WHATSAPP_URL } from './constants';

/**
 * Closing conversion band. Uses gray-800 as its base tone so the section bleeds
 * directly into the site footer with no visible seam.
 */
const FinalCTA: React.FC = () => {
  return (
    <section className="relative isolate overflow-hidden border-t border-gray-700/50 bg-gray-800 py-24 lg:py-32">
      <div aria-hidden="true" className="absolute inset-0 -z-10 overflow-hidden pointer-events-none">
        {/* Warm bloom behind the headline */}
        <div className="absolute left-1/2 -top-12 h-[34rem] w-[50rem] -translate-x-1/2 rounded-full bg-food-primary/20 blur-[120px]" />
        <div className="absolute inset-0 bg-grid-fine mask-radial opacity-40" />
        {/* Fade into the footer tone */}
        <div className="absolute inset-x-0 bottom-0 h-32 bg-gradient-to-b from-transparent to-gray-800" />
      </div>

      <div className="container mx-auto max-w-4xl px-4 text-center">
        <h2 className="text-display-sm font-extrabold text-white">
          Stop paying rent on
          <br className="hidden sm:block" /> your own customers.
        </h2>

        <p className="mx-auto mt-6 max-w-2xl text-lg leading-relaxed text-gray-300">
          Bring your menu, keep your margins. We will get your storefront and QR ordering
          live &mdash; usually the same day you sign up.
        </p>

        <div className="mt-10 flex flex-col items-stretch justify-center gap-3 sm:flex-row sm:items-center">
          <Button
            asChild
            size="lg"
            className="group h-14 rounded-xl bg-food-primary px-8 text-base font-bold text-white shadow-lg shadow-food-primary/25 transition-all hover:-translate-y-0.5 hover:bg-orange-600"
          >
            <a href={WHATSAPP_URL} target="_blank" rel="noopener noreferrer">
              <MessageSquare className="mr-2 h-5 w-5" />
              Partner with us
              <ArrowRight className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" />
            </a>
          </Button>

          <Button
            asChild
            size="lg"
            className="h-14 rounded-xl bg-white px-8 text-base font-bold text-gray-900 shadow-lg transition-all hover:-translate-y-0.5 hover:bg-gray-100"
          >
            <a href={SALES_MAILTO}>
              <Mail className="mr-2 h-5 w-5 text-gray-600" />
              Contact sales
            </a>
          </Button>
        </div>

        <p className="mt-8 text-sm text-gray-300 font-medium">
          No setup fee &middot; No lock-in contract &middot; No commission per order
        </p>
      </div>
    </section>
  );
};

export default FinalCTA;
