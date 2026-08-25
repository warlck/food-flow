import React from 'react';
import { Link } from 'react-router-dom';
import { Button } from '@/components/ui/button';
import {
  ArrowRight,
  BadgePercent,
  MessageSquare,
  Package,
  QrCode,
  ShieldCheck,
  Sparkles,
} from 'lucide-react';
import { SALES_MAILTO, WHATSAPP_URL } from './constants';

/**
 * Full-bleed dark hero aimed at restaurant owners.
 *
 * The background is composed of three stacked layers (aurora blobs, blueprint
 * grid, film grain) so the section keeps depth without shipping any imagery.
 */
const LandingHero: React.FC = () => {
  return (
    <section className="relative isolate overflow-hidden bg-ink-950 pt-[65px]">
      {/* ── Background layers ─────────────────────────────────────────── */}
      <div aria-hidden="true" className="absolute inset-0 -z-10">
        {/* Aurora blobs */}
        <div className="absolute -top-40 left-1/2 h-[38rem] w-[38rem] -translate-x-1/2 rounded-full bg-food-primary/25 blur-[120px] animate-aurora" />
        <div
          className="absolute -right-32 top-24 h-[26rem] w-[26rem] rounded-full bg-food-secondary/20 blur-[110px] animate-aurora"
          style={{ animationDelay: '-4s' }}
        />
        <div
          className="absolute -left-24 bottom-0 h-[24rem] w-[24rem] rounded-full bg-orange-600/20 blur-[110px] animate-aurora"
          style={{ animationDelay: '-8s' }}
        />
        {/* Blueprint grid, dissolved towards the bottom */}
        <div className="absolute inset-0 bg-grid mask-fade-y opacity-70" />
        {/* Grain */}
        <div className="absolute inset-0 bg-noise opacity-[0.15] mix-blend-soft-light" />
        {/* Hand off cleanly to the next (light) section */}
        <div className="absolute inset-x-0 bottom-0 h-40 bg-gradient-to-b from-transparent to-ink-950" />
      </div>

      <div className="container mx-auto max-w-7xl px-4 pb-24 pt-16 sm:pb-28 lg:pb-36 lg:pt-24">
        <div className="grid items-center gap-16 lg:grid-cols-[minmax(0,1fr)_minmax(0,26rem)] lg:gap-12">
          {/* ── Copy column ─────────────────────────────────────────── */}
          <div className="animate-rise-in text-center lg:text-left">
            {/* Eyebrow */}
            <div className="mb-8 inline-flex items-center gap-2 rounded-full px-4 py-1.5 text-sm font-semibold text-orange-200 glass-panel">
              <span className="relative flex h-2 w-2">
                <span className="absolute inline-flex h-full w-full rounded-full bg-food-secondary animate-pulse-ring" />
                <span className="relative inline-flex h-2 w-2 rounded-full bg-food-secondary" />
              </span>
              Live with restaurants across Singapore
            </div>

            <h1 className="text-display font-extrabold text-white">
              Your restaurant.
              <br />
              Your customers.
              <br />
              <span className="text-gradient-brand">Your margins.</span>
            </h1>

            <p className="mx-auto mt-8 max-w-xl text-lg leading-relaxed text-gray-400 sm:text-xl lg:mx-0">
              FoodFlow is the ordering platform restaurants own end to end. Launch a
              branded storefront, take QR and online orders, and keep{' '}
              <span className="font-semibold text-white">every dollar</span> of every
              order &mdash; no commissions, ever.
            </p>

            {/* Primary actions */}
            <div className="mt-10 flex flex-col items-stretch gap-3 sm:flex-row sm:items-center sm:justify-center lg:justify-start">
              <Button
                asChild
                size="lg"
                className="group h-14 rounded-xl bg-food-primary px-8 text-base font-bold text-white shadow-lg shadow-food-primary/25 transition-all hover:-translate-y-0.5 hover:bg-orange-600 hover:shadow-xl hover:shadow-food-primary/30"
              >
                <a href={WHATSAPP_URL} target="_blank" rel="noopener noreferrer">
                  <MessageSquare className="mr-2 h-5 w-5" />
                  Start taking orders
                  <ArrowRight className="ml-2 h-5 w-5 transition-transform group-hover:translate-x-1" />
                </a>
              </Button>

              <Button
                asChild
                size="lg"
                variant="outline"
                className="h-14 rounded-xl border-white/15 bg-white/5 px-8 text-base font-bold text-white backdrop-blur-md transition-all hover:-translate-y-0.5 hover:border-white/25 hover:bg-white/10 hover:text-white"
              >
                <a href={SALES_MAILTO}>Book a walkthrough</a>
              </Button>
            </div>

            {/* Trust row */}
            <ul className="mt-10 flex flex-wrap items-center justify-center gap-x-6 gap-y-3 text-sm text-gray-500 lg:justify-start">
              <li className="flex items-center gap-2">
                <BadgePercent className="h-4 w-4 text-food-secondary" />
                0% commission
              </li>
              <li className="flex items-center gap-2">
                <ShieldCheck className="h-4 w-4 text-food-secondary" />
                PCI-compliant payments
              </li>
              <li className="flex items-center gap-2">
                <Sparkles className="h-4 w-4 text-food-secondary" />
                Live in under a day
              </li>
            </ul>

            {/* Secondary path for diners */}
            <p className="mt-8 text-sm text-gray-500">
              Placed an order?{' '}
              <Link
                to="/track-order"
                className="inline-flex items-center gap-1 font-semibold text-gray-300 underline decoration-white/20 underline-offset-4 transition-colors hover:text-white hover:decoration-food-primary"
              >
                <Package className="h-4 w-4" />
                Track it here
              </Link>
            </p>
          </div>

          {/* ── Visual column: floating order ticket ─────────────────── */}
          <div className="relative mx-auto w-full max-w-sm lg:max-w-none">
            <div className="animate-float-y">
              <div className="relative rounded-3xl p-6 glass-panel-strong ring-lit">
                {/* Ticket header */}
                <div className="flex items-center justify-between border-b border-white/10 pb-5">
                  <div>
                    <p className="text-xs font-medium uppercase tracking-widest text-gray-500">
                      Incoming order
                    </p>
                    <p className="mt-1 text-lg font-bold text-white">Table 12</p>
                  </div>
                  <span className="rounded-full bg-food-success/15 px-3 py-1 text-xs font-bold text-emerald-300 ring-1 ring-inset ring-emerald-400/25">
                    Paid
                  </span>
                </div>

                {/* Line items */}
                <ul className="space-y-4 py-5 text-sm">
                  {[
                    { qty: 2, name: 'Chicken Rice', note: 'Extra chilli', price: '11.00' },
                    { qty: 1, name: 'Laksa Bowl', note: 'No cockles', price: '8.50' },
                    { qty: 2, name: 'Iced Kopi', note: 'Less sweet', price: '5.20' },
                  ].map((item) => (
                    <li key={item.name} className="flex items-start justify-between gap-4">
                      <div className="flex items-start gap-3">
                        <span className="mt-0.5 flex h-6 w-6 shrink-0 items-center justify-center rounded-md bg-white/[0.08] text-xs font-bold text-gray-300">
                          {item.qty}
                        </span>
                        <div>
                          <p className="font-semibold text-white">{item.name}</p>
                          <p className="text-xs text-gray-500">{item.note}</p>
                        </div>
                      </div>
                      <span className="font-semibold tabular-nums text-gray-300">
                        ${item.price}
                      </span>
                    </li>
                  ))}
                </ul>

                {/* Totals: the pitch, quantified */}
                <div className="space-y-2.5 border-t border-white/10 pt-5 text-sm">
                  <div className="flex items-center justify-between text-gray-400">
                    <span>Order total</span>
                    <span className="font-semibold tabular-nums text-white">$24.70</span>
                  </div>
                  <div className="flex items-center justify-between text-gray-400">
                    <span>Platform commission</span>
                    <span className="font-semibold tabular-nums text-emerald-300">
                      &minus;$0.00
                    </span>
                  </div>
                  <div className="flex items-center justify-between border-t border-white/10 pt-3 text-base">
                    <span className="font-bold text-white">You keep</span>
                    <span className="text-xl font-extrabold tabular-nums text-gradient-brand">
                      $24.70
                    </span>
                  </div>
                </div>
              </div>
            </div>

            {/* QR badge pinned to the ticket */}
            <div
              className="absolute -bottom-6 -left-4 hidden animate-float-y items-center gap-3 rounded-2xl px-4 py-3 glass-panel-strong ring-lit sm:flex"
              style={{ animationDelay: '-3s' }}
            >
              <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-food-primary">
                <QrCode className="h-5 w-5 text-white" />
              </div>
              <div className="pr-1">
                <p className="text-sm font-bold leading-tight text-white">Scan to order</p>
                <p className="text-xs text-gray-400">No app download</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default LandingHero;
