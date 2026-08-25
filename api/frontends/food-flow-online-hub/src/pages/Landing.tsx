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
      <div className="bg-ink-950 font-sans text-white selection:bg-food-primary/30">
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
