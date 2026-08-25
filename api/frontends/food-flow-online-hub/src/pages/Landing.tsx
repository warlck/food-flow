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
    <Layout>
      <div className="relative isolate overflow-hidden bg-ink-950 font-sans text-white selection:bg-food-primary/30">
        {/* Subtle, consistent warm reddish background hue across the whole landing canvas */}
        <div aria-hidden="true" className="pointer-events-none absolute inset-0 -z-10 overflow-hidden">
          <div className="absolute top-[18%] -left-32 h-[48rem] w-[48rem] rounded-full bg-food-primary/[0.07] blur-[150px]" />
          <div className="absolute top-[42%] -right-32 h-[52rem] w-[52rem] rounded-full bg-orange-600/[0.06] blur-[160px]" />
          <div className="absolute top-[68%] -left-24 h-[46rem] w-[46rem] rounded-full bg-food-primary/[0.07] blur-[150px]" />
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
