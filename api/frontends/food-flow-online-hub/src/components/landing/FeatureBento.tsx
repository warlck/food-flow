import React from 'react';
import {
  BadgePercent,
  BarChart3,
  CreditCard,
  QrCode,
  Radio,
  UtensilsCrossed,
} from 'lucide-react';

/**
 * Shared shell for every bento tile so hover, glass and radius stay consistent.
 */
const Tile: React.FC<{ className?: string; children: React.ReactNode }> = ({
  className = '',
  children,
}) => (
  <div
    className={`group relative overflow-hidden rounded-3xl p-8 glass-panel ring-lit transition-all duration-500 hover:border-white/20 hover:bg-white/[0.06] ${className}`}
  >
    {/* Warm glow that blooms on hover */}
    <div
      aria-hidden="true"
      className="pointer-events-none absolute -right-16 -top-16 h-48 w-48 rounded-full bg-food-primary/20 opacity-0 blur-3xl transition-opacity duration-500 group-hover:opacity-100"
    />
    <div className="relative">{children}</div>
  </div>
);

const TileIcon: React.FC<{ icon: React.ElementType }> = ({ icon: Icon }) => (
  <div className="mb-6 inline-flex h-12 w-12 items-center justify-center rounded-xl bg-food-primary/15 ring-1 ring-inset ring-food-primary/25">
    <Icon className="h-6 w-6 text-food-primary" />
  </div>
);

/**
 * Asymmetric feature grid. The two wide tiles carry small inline product
 * visualisations; the rest stay text-forward so the section keeps its rhythm.
 */
const FeatureBento: React.FC = () => {
  return (
    <section className="relative isolate overflow-hidden border-t border-white/10 bg-ink-900 py-24 lg:py-32">
      <div aria-hidden="true" className="absolute inset-0 -z-10 bg-noise opacity-[0.12] mix-blend-soft-light" />

      <div className="container mx-auto max-w-7xl px-4">
        <header className="mx-auto max-w-3xl text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-food-primary">
            The platform
          </p>
          <h2 className="mt-4 text-display-sm font-extrabold text-white">
            Everything a modern restaurant
            <br className="hidden sm:block" /> actually needs.
          </h2>
          <p className="mt-5 text-lg leading-relaxed text-gray-400">
            One system for menus, ordering, payments and reporting &mdash; instead of four
            vendors and a spreadsheet.
          </p>
        </header>

        <div className="mt-16 grid gap-5 lg:grid-cols-3">
          {/* ── Wide: commission economics ───────────────────────────── */}
          <Tile className="lg:col-span-2">
            <TileIcon icon={BadgePercent} />
            <h3 className="text-2xl font-bold tracking-tight text-white">
              Zero commission, permanently
            </h3>
            <p className="mt-3 max-w-lg text-sm leading-relaxed text-gray-400">
              Aggregators take 25&ndash;30% of every order. FoodFlow charges a flat
              subscription, so a busier month makes you more money instead of a bigger bill.
            </p>

            {/* Comparison bars */}
            <div className="mt-8 space-y-4">
              <div>
                <div className="mb-2 flex items-baseline justify-between text-xs font-semibold">
                  <span className="text-gray-500">Typical aggregator &mdash; you keep</span>
                  <span className="tabular-nums text-gray-400">$720 of $1,000</span>
                </div>
                <div className="h-2.5 overflow-hidden rounded-full bg-white/[0.06]">
                  <div className="h-full w-[72%] rounded-full bg-gray-600" />
                </div>
              </div>
              <div>
                <div className="mb-2 flex items-baseline justify-between text-xs font-semibold">
                  <span className="text-white">FoodFlow &mdash; you keep</span>
                  <span className="tabular-nums text-food-secondary">$1,000 of $1,000</span>
                </div>
                <div className="relative h-2.5 overflow-hidden rounded-full bg-white/[0.06]">
                  <div className="h-full w-full rounded-full bg-gradient-to-r from-food-secondary via-food-accent to-food-primary" />
                  <div
                    aria-hidden="true"
                    className="absolute inset-0 animate-shimmer bg-gradient-to-r from-transparent via-white/30 to-transparent"
                  />
                </div>
              </div>
            </div>
          </Tile>

          {/* ── Live order tracking ──────────────────────────────────── */}
          <Tile>
            <TileIcon icon={Radio} />
            <h3 className="text-xl font-bold tracking-tight text-white">Live order tracking</h3>
            <p className="mt-3 text-sm leading-relaxed text-gray-400">
              Diners follow their order from confirmed to ready without calling the counter.
            </p>

            <ul className="mt-7 space-y-3">
              {[
                { label: 'Confirmed', done: true },
                { label: 'In the kitchen', done: true },
                { label: 'Ready for pickup', done: false },
              ].map((stage) => (
                <li key={stage.label} className="flex items-center gap-3 text-sm">
                  <span
                    className={
                      stage.done
                        ? 'flex h-5 w-5 items-center justify-center rounded-full bg-food-primary text-[10px] font-bold text-white'
                        : 'relative flex h-5 w-5 items-center justify-center rounded-full border border-dashed border-white/25'
                    }
                  >
                    {stage.done ? '✓' : null}
                    {!stage.done && (
                      <span className="absolute h-2 w-2 rounded-full bg-food-secondary animate-pulse-ring" />
                    )}
                  </span>
                  <span className={stage.done ? 'font-medium text-white' : 'text-gray-500'}>
                    {stage.label}
                  </span>
                </li>
              ))}
            </ul>
          </Tile>

          {/* ── Payments ─────────────────────────────────────────────── */}
          <Tile>
            <TileIcon icon={CreditCard} />
            <h3 className="text-xl font-bold tracking-tight text-white">Payments built in</h3>
            <p className="mt-3 text-sm leading-relaxed text-gray-400">
              Cards and PayNow settle to your own account, secured by Stripe. No terminal
              rental, no chasing cash.
            </p>
            <div className="mt-7 flex flex-wrap gap-2">
              {['Visa', 'Mastercard', 'Amex', 'PayNow', 'Apple Pay'].map((method) => (
                <span
                  key={method}
                  className="rounded-lg border border-white/10 bg-white/[0.04] px-3 py-1.5 text-xs font-semibold text-gray-300"
                >
                  {method}
                </span>
              ))}
            </div>
          </Tile>

          {/* ── QR ordering ──────────────────────────────────────────── */}
          <Tile>
            <TileIcon icon={QrCode} />
            <h3 className="text-xl font-bold tracking-tight text-white">QR table ordering</h3>
            <p className="mt-3 text-sm leading-relaxed text-gray-400">
              One code per table turns every seat into a self-serve till. Guests order and
              pay without waiting for staff.
            </p>
            <p className="mt-7 inline-flex items-center gap-2 rounded-lg bg-food-success/10 px-3 py-2 text-xs font-semibold text-emerald-300 ring-1 ring-inset ring-emerald-400/20">
              Frees up floor staff at peak
            </p>
          </Tile>

          {/* ── Digital menus ────────────────────────────────────────── */}
          <Tile>
            <TileIcon icon={UtensilsCrossed} />
            <h3 className="text-xl font-bold tracking-tight text-white">Menus you control</h3>
            <p className="mt-3 text-sm leading-relaxed text-gray-400">
              Change prices, hide sold-out dishes and launch specials instantly. No
              reprints, no waiting on a vendor.
            </p>
            <div className="mt-7 space-y-2">
              {[
                { name: 'Signature Laksa', state: 'Available', good: true },
                { name: 'Chilli Crab', state: 'Sold out', good: false },
              ].map((item) => (
                <div
                  key={item.name}
                  className="flex items-center justify-between rounded-xl border border-white/10 bg-white/[0.03] px-3.5 py-2.5 text-xs"
                >
                  <span className="font-semibold text-gray-200">{item.name}</span>
                  <span
                    className={
                      item.good
                        ? 'font-bold text-emerald-300'
                        : 'font-bold text-gray-500'
                    }
                  >
                    {item.state}
                  </span>
                </div>
              ))}
            </div>
          </Tile>

          {/* ── Wide: insights ───────────────────────────────────────── */}
          <Tile className="lg:col-span-2">
            <TileIcon icon={BarChart3} />
            <h3 className="text-2xl font-bold tracking-tight text-white">
              Know what actually sells
            </h3>
            <p className="mt-3 max-w-lg text-sm leading-relaxed text-gray-400">
              Revenue, average order value, peak hours and top items &mdash; computed from
              your own orders, not an aggregator&apos;s black box.
            </p>

            {/* Bar chart sketch */}
            <div className="mt-8 flex items-end gap-2 sm:gap-3" aria-hidden="true">
              {[38, 52, 44, 68, 58, 82, 100].map((height, index) => (
                <div key={index} className="flex flex-1 flex-col items-center gap-2">
                  <div
                    className="w-full rounded-t-md bg-gradient-to-t from-food-primary/25 to-food-primary transition-all duration-500 group-hover:from-food-primary/40"
                    style={{ height: `${Math.round(height * 0.9)}px` }}
                  />
                  <span className="text-[10px] font-medium text-gray-600">
                    {['M', 'T', 'W', 'T', 'F', 'S', 'S'][index]}
                  </span>
                </div>
              ))}
            </div>
          </Tile>
        </div>
      </div>
    </section>
  );
};

export default FeatureBento;
