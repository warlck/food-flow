import React from 'react';
import Layout from '@/components/Layout';
import LandingHero from '@/components/landing/LandingHero';
import StatsBar from '@/components/landing/StatsBar';
import HowItWorks from '@/components/landing/HowItWorks';
import FeatureBento from '@/components/landing/FeatureBento';
import ProductShowcase from '@/components/landing/ProductShowcase';
import FinalCTA from '@/components/landing/FinalCTA';

/**
 * Marketing landing page for restaurant owners.
 *
 * The page is a dark, high-contrast surface that flows from the hero through to
 * the existing gray-800 footer. Each section owns its own background so the
 * composition here stays declarative.
 */
const Landing: React.FC = () => {
  return (
    <Layout surface="marketing" className="bg-gray-800 text-white">
      <div className="relative isolate min-h-screen bg-gray-800 font-sans text-white selection:bg-food-primary/30 -mt-[65px] animate-fade-in">
        {/* Background stack matching MarketingPage & Track Order */}
        <div aria-hidden="true" className="pointer-events-none absolute inset-0 -z-10 overflow-hidden">
          {/* Warm bloom behind hero headline and header */}
          <div className="absolute left-1/2 -top-24 h-[40rem] w-[56rem] -translate-x-1/2 rounded-full bg-food-primary/20 blur-[120px]" />
          {/* Continuous warm ambient glows down the landing canvas */}
          <div className="absolute top-[28%] -left-32 h-[48rem] w-[48rem] rounded-full bg-food-primary/[0.12] blur-[150px]" />
          <div className="absolute top-[52%] -right-32 h-[52rem] w-[52rem] rounded-full bg-orange-600/[0.10] blur-[160px]" />
          <div className="absolute top-[76%] -left-24 h-[46rem] w-[46rem] rounded-full bg-food-primary/[0.12] blur-[150px]" />
          {/* Fine grid overlay with radial mask */}
          <div className="absolute inset-0 bg-grid-fine opacity-40" />
        </div>

        <LandingHero />
        <StatsBar />
        <HowItWorks />
        <FeatureBento />
        <ProductShowcase />
        <FinalCTA />
      </div>
    </Layout>
  );
};

export default Landing;
