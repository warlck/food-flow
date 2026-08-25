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
    <Layout surface="marketing" className="bg-gray-800 text-white">
      <div className="relative isolate min-h-[calc(100vh-160px)] bg-gray-800 font-sans text-white selection:bg-food-primary/30 animate-fade-in">
        {/* Background stack matching FinalCTA ("Stop paying rent on your own customers") */}
        <div aria-hidden="true" className="absolute inset-0 -z-10 overflow-hidden pointer-events-none">
          {/* Warm bloom behind headline and header */}
          <div className="absolute left-1/2 -top-24 h-[36rem] w-[52rem] -translate-x-1/2 rounded-full bg-food-primary/20 blur-[120px]" />
          {/* Fine grid overlay with radial mask */}
          <div className="absolute inset-0 bg-grid-fine mask-radial opacity-40" />
          {/* Fade into the footer tone */}
          <div className="absolute inset-x-0 bottom-0 h-32 bg-gradient-to-b from-transparent to-gray-800" />
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
