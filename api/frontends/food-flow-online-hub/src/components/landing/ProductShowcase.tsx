import React from 'react';
import { ArrowUpRight, Check, Plus, ShoppingBag } from 'lucide-react';

const KPIS = [
  { label: 'Revenue today', value: '$4,182', delta: '+18.2%' },
  { label: 'Orders', value: '164', delta: '+9.4%' },
  { label: 'Avg order value', value: '$25.50', delta: '+4.1%' },
];

const LIVE_ORDERS = [
  { id: '#1841', table: 'Table 04', total: '$38.20', status: 'Preparing', tone: 'amber' },
  { id: '#1840', table: 'Takeaway', total: '$16.80', status: 'Ready', tone: 'emerald' },
  { id: '#1839', table: 'Table 11', total: '$52.00', status: 'Served', tone: 'slate' },
];

const STATUS_TONE: Record<string, string> = {
  amber: 'bg-amber-400/15 text-amber-300 ring-amber-400/25',
  emerald: 'bg-emerald-400/15 text-emerald-300 ring-emerald-400/25',
  slate: 'bg-white/[0.06] text-gray-400 ring-white/10',
};

/**
 * Product surfaces rendered as markup rather than screenshots: an owner-facing
 * dashboard in a browser chrome, with the diner ordering flow overlapping it in
 * a phone frame on large screens.
 */
const ProductShowcase: React.FC = () => {
  return (
    <section className="relative isolate overflow-hidden border-t border-white/10 bg-transparent py-24 lg:py-32">
      <div aria-hidden="true" className="absolute inset-0 -z-10 overflow-hidden pointer-events-none">
        {/* Subtle warm bloom positioned lower to seamlessly transition into FinalCTA */}
        <div className="absolute left-1/2 bottom-[-4rem] h-[36rem] w-[56rem] -translate-x-1/2 rounded-full bg-food-primary/[0.09] blur-[150px]" />
      </div>

      <div className="container mx-auto max-w-7xl px-4">
        <header className="mx-auto max-w-3xl text-center">
          <p className="text-xs font-semibold uppercase tracking-[0.2em] text-food-primary">
            See it in action
          </p>
          <h2 className="mt-4 text-display-sm font-extrabold text-white">
            Two screens. One source of truth.
          </h2>
          <p className="mt-5 text-lg leading-relaxed text-gray-300">
            Your team works the dashboard while guests order from their own phones. Both
            sides stay in sync, order by order.
          </p>
        </header>

        <div className="relative mt-20 lg:pr-32 xl:pr-40">
          {/* ── Browser frame: owner dashboard ───────────────────────── */}
          <div className="overflow-hidden rounded-2xl glass-panel-strong ring-lit">
            {/* Chrome */}
            <div className="flex items-center gap-2 border-b border-white/10 bg-white/[0.03] px-4 py-3">
              <span className="h-3 w-3 rounded-full bg-red-400/70" />
              <span className="h-3 w-3 rounded-full bg-amber-400/70" />
              <span className="h-3 w-3 rounded-full bg-emerald-400/70" />
              <div className="ml-4 flex-1">
                <div className="mx-auto w-full max-w-xs truncate rounded-md bg-black/30 px-3 py-1 text-center text-[11px] text-gray-300 font-medium">
                  app.foodflow.sg/dashboard
                </div>
              </div>
            </div>

            {/* Dashboard body */}
            <div className="bg-ink-900/80 p-5 sm:p-7">
              {/* KPI row */}
              <div className="grid gap-3 sm:grid-cols-3">
                {KPIS.map((kpi) => (
                  <div
                    key={kpi.label}
                    className="rounded-xl border border-white/10 bg-white/[0.04] p-4"
                  >
                    <p className="text-[11px] font-semibold uppercase tracking-wider text-gray-400">
                      {kpi.label}
                    </p>
                    <p className="mt-2 text-2xl font-bold tabular-nums text-white">
                      {kpi.value}
                    </p>
                    <p className="mt-1 inline-flex items-center gap-1 text-xs font-semibold text-emerald-300">
                      <ArrowUpRight className="h-3.5 w-3.5" />
                      {kpi.delta}
                    </p>
                  </div>
                ))}
              </div>

              {/* Chart + live orders */}
              <div className="mt-5 grid gap-5 lg:grid-cols-[1.4fr_1fr]">
                {/* Sparkline-ish area chart */}
                <div className="rounded-xl border border-white/10 bg-white/[0.03] p-5">
                  <div className="mb-5 flex items-center justify-between">
                    <p className="text-sm font-semibold text-white">Revenue, last 14 days</p>
                    <span className="rounded-md bg-white/[0.06] px-2 py-1 text-[10px] font-semibold text-gray-300">
                      Daily
                    </span>
                  </div>
                  <svg
                    viewBox="0 0 320 96"
                    className="h-28 w-full"
                    preserveAspectRatio="none"
                    aria-hidden="true"
                  >
                    <defs>
                      <linearGradient id="ff-area" x1="0" y1="0" x2="0" y2="1">
                        <stop offset="0%" stopColor="#FF4500" stopOpacity="0.45" />
                        <stop offset="100%" stopColor="#FF4500" stopOpacity="0" />
                      </linearGradient>
                    </defs>
                    <path
                      d="M0 74 L27 66 L53 70 L80 54 L107 60 L133 42 L160 47 L187 32 L213 38 L240 24 L267 30 L293 16 L320 10 L320 96 L0 96 Z"
                      fill="url(#ff-area)"
                    />
                    <path
                      d="M0 74 L27 66 L53 70 L80 54 L107 60 L133 42 L160 47 L187 32 L213 38 L240 24 L267 30 L293 16 L320 10"
                      fill="none"
                      stroke="#FF8C42"
                      strokeWidth="2.5"
                      strokeLinecap="round"
                      strokeLinejoin="round"
                    />
                  </svg>
                </div>

                {/* Live order queue */}
                <div className="rounded-xl border border-white/10 bg-white/[0.03] p-5">
                  <div className="mb-4 flex items-center gap-2">
                    <span className="relative flex h-2 w-2">
                      <span className="absolute inline-flex h-full w-full rounded-full bg-food-primary animate-pulse-ring" />
                      <span className="relative inline-flex h-2 w-2 rounded-full bg-food-primary" />
                    </span>
                    <p className="text-sm font-semibold text-white">Live orders</p>
                  </div>

                  <ul className="space-y-2.5">
                    {LIVE_ORDERS.map((order) => (
                      <li
                        key={order.id}
                        className="flex items-center justify-between gap-3 rounded-lg bg-white/[0.04] px-3 py-2.5"
                      >
                        <div className="min-w-0">
                          <p className="truncate text-xs font-bold text-white">
                            {order.id} &middot; {order.table}
                          </p>
                          <p className="text-[11px] tabular-nums text-gray-400 font-medium">{order.total}</p>
                        </div>
                        <span
                          className={`shrink-0 rounded-full px-2.5 py-1 text-[10px] font-bold ring-1 ring-inset ${STATUS_TONE[order.tone]}`}
                        >
                          {order.status}
                        </span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>
            </div>
          </div>

          {/* ── Phone frame: diner ordering ──────────────────────────── */}
          <div className="mx-auto mt-10 w-[17rem] lg:absolute lg:-right-2 lg:bottom-0 lg:mt-0 xl:right-0">
            <div className="animate-float-y rounded-[2.5rem] border border-white/15 bg-ink-850 p-2.5 shadow-2xl shadow-black/60">
              <div className="overflow-hidden rounded-[2rem] bg-white">
                {/* Notch */}
                <div className="flex justify-center bg-white pt-3">
                  <div className="h-1.5 w-16 rounded-full bg-gray-200" />
                </div>

                {/* Restaurant header */}
                <div className="px-4 pb-3 pt-4">
                  <p className="text-[10px] font-semibold uppercase tracking-wider text-food-primary">
                    Table 12
                  </p>
                  <p className="text-base font-extrabold text-gray-900">Kopitiam Corner</p>
                </div>

                {/* Category chips */}
                <div className="flex gap-1.5 overflow-hidden px-4 pb-3">
                  {['Popular', 'Rice', 'Noodles', 'Drinks'].map((chip, index) => (
                    <span
                      key={chip}
                      className={
                        index === 0
                          ? 'shrink-0 rounded-full bg-food-primary px-2.5 py-1 text-[10px] font-bold text-white'
                          : 'shrink-0 rounded-full bg-gray-100 px-2.5 py-1 text-[10px] font-semibold text-gray-600'
                      }
                    >
                      {chip}
                    </span>
                  ))}
                </div>

                {/* Menu items */}
                <div className="space-y-2 px-4 pb-3">
                  {[
                    { name: 'Chicken Rice', price: '5.50', added: true },
                    { name: 'Laksa Bowl', price: '8.50', added: false },
                    { name: 'Iced Kopi', price: '2.60', added: false },
                  ].map((item) => (
                    <div
                      key={item.name}
                      className="flex items-center gap-3 rounded-xl border border-gray-100 p-2.5"
                    >
                      <div className="h-10 w-10 shrink-0 rounded-lg bg-gradient-to-br from-orange-100 to-amber-50" />
                      <div className="min-w-0 flex-1">
                        <p className="truncate text-xs font-bold text-gray-900">{item.name}</p>
                        <p className="text-[11px] font-semibold tabular-nums text-gray-500">
                          ${item.price}
                        </p>
                      </div>
                      <span
                        className={
                          item.added
                            ? 'flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-food-success text-white'
                            : 'flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-food-primary text-white'
                        }
                      >
                        {item.added ? (
                          <Check className="h-3.5 w-3.5" />
                        ) : (
                          <Plus className="h-3.5 w-3.5" />
                        )}
                      </span>
                    </div>
                  ))}
                </div>

                {/* Sticky checkout bar */}
                <div className="border-t border-gray-100 p-3">
                  <div className="flex items-center justify-between rounded-xl bg-food-primary px-3.5 py-3 text-white">
                    <span className="flex items-center gap-2 text-xs font-bold">
                      <ShoppingBag className="h-4 w-4" />
                      3 items
                    </span>
                    <span className="text-xs font-extrabold tabular-nums">$24.70</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
};

export default ProductShowcase;
