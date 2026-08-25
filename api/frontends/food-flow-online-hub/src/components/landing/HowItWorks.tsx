import React from 'react';
import { UploadCloud, QrCode, ChefHat, Wallet } from 'lucide-react';
import type { LucideIcon } from 'lucide-react';

interface Step {
  icon: LucideIcon;
  title: string;
  body: string;
  detail: string;
}

const STEPS: Step[] = [
  {
    icon: UploadCloud,
    title: 'Load your menu',
    body: 'Import categories, items, add-ons and pricing once. Photos, dietary tags and modifiers included.',
    detail: 'Step 01',
  },
  {
    icon: QrCode,
    title: 'Publish your storefront',
    body: 'Get a branded ordering page plus printable QR codes for every table, counter and takeaway flyer.',
    detail: 'Step 02',
  },
  {
    icon: ChefHat,
    title: 'Orders hit the kitchen',
    body: 'Paid tickets arrive instantly with modifiers and notes. Update status and diners see it live.',
    detail: 'Step 03',
  },
  {
    icon: Wallet,
    title: 'Keep the whole payment',
    body: 'Card and PayNow settle straight to your account. No commission deducted, no reconciliation games.',
    detail: 'Step 04',
  },
];

/**
 * Four-step onboarding narrative. Steps are connected by a horizontal rule on
 * large screens and stack cleanly on mobile.
 */
const HowItWorks: React.FC = () => {
  return (
    <section className="relative isolate overflow-hidden bg-transparent py-24 lg:py-32">
      <div aria-hidden="true" className="absolute inset-0 -z-10 overflow-hidden pointer-events-none">
        <div className="absolute left-1/4 top-1/2 h-[34rem] w-[46rem] -translate-y-1/2 rounded-full bg-food-primary/[0.08] blur-[140px]" />
      </div>

      <div className="container mx-auto max-w-7xl px-4">
        <header className="mx-auto max-w-3xl text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-food-primary">
            How it works
          </p>
          <h2 className="mt-4 text-display-sm font-extrabold text-white">
            Live by tomorrow&apos;s lunch service.
          </h2>
          <p className="mt-5 text-lg leading-relaxed text-gray-400">
            No hardware to buy, no integrations to schedule, no engineer required. Four
            steps between signing up and taking your first paid order.
          </p>
        </header>

        <div className="relative mt-20">
          {/* Connector line, desktop only */}
          <div
            aria-hidden="true"
            className="absolute left-0 right-0 top-8 hidden h-px bg-gradient-to-r from-transparent via-white/15 to-transparent lg:block"
          />

          <ol className="grid gap-12 sm:grid-cols-2 lg:grid-cols-4 lg:gap-8">
            {STEPS.map((step) => {
              const Icon = step.icon;
              return (
                <li key={step.title} className="relative">
                  {/* Node */}
                  <div className="relative z-10 mb-7 flex h-16 w-16 items-center justify-center rounded-2xl bg-ink-850 glass-panel-strong ring-lit">
                    <Icon className="h-7 w-7 text-food-primary" />
                  </div>

                  <p className="mb-2 text-xs font-bold uppercase tracking-[0.18em] text-gray-600">
                    {step.detail}
                  </p>
                  <h3 className="text-xl font-bold tracking-tight text-white">{step.title}</h3>
                  <p className="mt-3 text-sm leading-relaxed text-gray-400">{step.body}</p>
                </li>
              );
            })}
          </ol>
        </div>
      </div>
    </section>
  );
};

export default HowItWorks;
