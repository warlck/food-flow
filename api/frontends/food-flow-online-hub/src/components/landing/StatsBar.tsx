import React from 'react';

const STATS = [
  { value: '0%', label: 'Commission on every order', sub: 'Flat subscription, no revenue share' },
  { value: '<1 day', label: 'From signup to first order', sub: 'Menu import and QR codes included' },
  { value: '3 taps', label: 'Average diner checkout', sub: 'No app, no account, no friction' },
  { value: '24/7', label: 'Orders without extra staff', sub: 'Self-serve QR and online ordering' },
];

const FORMATS = [
  'Hawker stalls',
  'Cafes',
  'Fine dining',
  'Bubble tea',
  'Bakeries',
  'Ramen bars',
  'Cloud kitchens',
  'Food courts',
  'Bistros',
  'Dessert shops',
];

/**
 * Proof band directly under the hero: hard numbers first, then a marquee of the
 * restaurant formats the platform serves. Rendered on the dark surface so the
 * hero flows into it without a visible seam.
 */
const StatsBar: React.FC = () => {
  return (
    <section className="relative isolate overflow-hidden border-y border-white/10 bg-ink-900">
      <div aria-hidden="true" className="absolute inset-0 -z-10 bg-grid-fine opacity-40" />

      <div className="container mx-auto max-w-7xl px-4 py-16 lg:py-20">
        {/* Numbers */}
        <dl className="grid grid-cols-2 gap-x-6 gap-y-10 lg:grid-cols-4">
          {STATS.map((stat) => (
            <div key={stat.label} className="text-center lg:text-left">
              <dd className="text-4xl font-extrabold tracking-tight text-gradient-brand sm:text-5xl">
                {stat.value}
              </dd>
              <dt className="mt-3 text-sm font-semibold text-white sm:text-base">
                {stat.label}
              </dt>
              <p className="mt-1 text-xs leading-relaxed text-gray-500 sm:text-sm">{stat.sub}</p>
            </div>
          ))}
        </dl>

        {/* Format marquee */}
        <div className="mt-16 border-t border-white/10 pt-10">
          <p className="mb-6 text-center text-xs font-semibold uppercase tracking-[0.2em] text-gray-500">
            Built for every kind of kitchen
          </p>

          <div className="mask-fade-x pause-on-hover overflow-hidden">
            <div className="marquee-track animate-marquee-slow gap-4">
              {/* Duplicated so the -50% translation loops seamlessly */}
              {[...FORMATS, ...FORMATS].map((format, index) => (
                <span
                  key={`${format}-${index}`}
                  className="shrink-0 rounded-full px-5 py-2.5 text-sm font-medium text-gray-300 glass-panel"
                >
                  {format}
                </span>
              ))}
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default StatsBar;
