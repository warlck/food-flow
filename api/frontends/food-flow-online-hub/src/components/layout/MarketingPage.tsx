import React from 'react';
import { Link } from 'react-router-dom';
import { ArrowLeft, type LucideIcon } from 'lucide-react';
import Layout from '@/components/Layout';

export interface MarketingPageProps {
  eyebrow?: string;
  title: string;
  description?: string;
  icon?: LucideIcon;
  maxWidth?: string;
  children: React.ReactNode;
}

export const MarketingPage: React.FC<MarketingPageProps> = ({
  eyebrow,
  title,
  description,
  icon: Icon,
  maxWidth = 'max-w-4xl',
  children,
}) => {
  return (
    <Layout surface="marketing">
      <div className="relative isolate min-h-[calc(100vh-160px)] bg-ink-950 font-sans text-white selection:bg-food-primary/30 animate-fade-in">
        {/* Background stack */}
        <div aria-hidden="true" className="absolute inset-0 -z-10 overflow-hidden pointer-events-none">
          <div className="absolute -top-40 left-1/2 h-[38rem] w-[38rem] -translate-x-1/2 rounded-full bg-food-primary/20 blur-[120px]" />
          <div className="absolute -right-32 top-24 h-[24rem] w-[24rem] rounded-full bg-food-secondary/15 blur-[110px]" />
          <div className="absolute inset-0 bg-grid mask-fade-y opacity-60" />
          <div className="absolute inset-0 bg-noise opacity-[0.15] mix-blend-soft-light" />
        </div>

        <div className={`container mx-auto px-4 py-12 md:py-16 ${maxWidth}`}>
          {/* Back to Home Link */}
          <div className="mb-8">
            <Link
              to="/"
              className="inline-flex items-center text-sm font-medium text-gray-400 hover:text-white transition-colors gap-2"
            >
              <ArrowLeft className="h-4 w-4" />
              <span>Back to home</span>
            </Link>
          </div>

          {/* Page Masthead */}
          <div className="mb-10 text-left">
            {eyebrow && (
              <div className="mb-3 text-xs font-bold uppercase tracking-wider text-food-primary">
                {eyebrow}
              </div>
            )}
            <div className="flex items-center gap-4 mb-4">
              {Icon && (
                <div className="w-12 h-12 rounded-2xl bg-white/5 ring-1 ring-white/10 flex items-center justify-center text-food-primary shrink-0">
                  <Icon className="h-6 w-6" />
                </div>
              )}
              <h1 className="text-3xl sm:text-4xl md:text-5xl font-extrabold tracking-tight text-white">
                {title}
              </h1>
            </div>
            {description && (
              <p className="text-base sm:text-lg text-gray-400 leading-relaxed max-w-2xl">
                {description}
              </p>
            )}
          </div>

          {/* Page Content */}
          <div className="relative z-10">
            {children}
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default MarketingPage;
