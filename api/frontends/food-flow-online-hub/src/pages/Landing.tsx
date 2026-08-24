import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import Layout from '@/components/Layout';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Package,
  Clock,
  CreditCard,
  Smartphone,
  Search,
  MessageSquare,
  Mail,
  Store,
  ChevronRight
} from 'lucide-react';

const Landing: React.FC = () => {
  const navigate = useNavigate();
  const [orderSearchId, setOrderSearchId] = useState('');

  const handleTrackSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (orderSearchId.trim()) {
      navigate(`/track-order/${orderSearchId.trim()}`);
    } else {
      navigate('/track-order');
    }
  };

  return (
    <Layout>
      {/* 
        Modern, Light Theme with Thematic Imagery 
        Top matches Header (bg-white)
        Bottom matches Footer (bg-gray-800)
      */}
      <div className="bg-white min-h-screen text-gray-900 font-sans selection:bg-food-primary/20">
        
        {/* ═══════════════════════════════════════════════════════════
            HERO SECTION — Clean White Base
         ═══════════════════════════════════════════════════════════ */}
        <section className="relative pt-16 pb-20 lg:pt-24 lg:pb-32 overflow-hidden">
          <div className="container mx-auto px-4 max-w-7xl">
            <div className="flex flex-col lg:flex-row items-center gap-12 lg:gap-16">
              
              {/* Left Content */}
              <div className="w-full lg:w-1/2 flex flex-col items-center lg:items-start text-center lg:text-left z-10">
                <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-orange-50 border border-orange-100 text-food-primary font-semibold text-sm mb-6 shadow-sm">
                  <span className="relative flex h-2 w-2">
                    <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-food-primary opacity-75" />
                    <span className="relative inline-flex rounded-full h-2 w-2 bg-food-primary" />
                  </span>
                  Now live in Singapore
                </div>
                
                <h1 className="text-5xl sm:text-6xl lg:text-[4.5rem] font-extrabold tracking-tight text-gray-900 leading-[1.05] mb-6">
                  Great food, <br className="hidden lg:block"/>
                  <span className="text-food-primary">zero friction.</span>
                </h1>
                
                <p className="text-lg md:text-xl text-gray-600 max-w-xl mb-10 leading-relaxed">
                  Skip the queue. Scan, order, and pay instantly. Experience the most elegant way to dine in or take out at your favourite restaurants.
                </p>

                {/* Track Order Input */}
                <div className="w-full max-w-md relative group bg-white rounded-2xl shadow-xl shadow-gray-200/60 ring-1 ring-gray-100">
                  <form
                    onSubmit={handleTrackSubmit}
                    className="relative flex items-center p-2"
                  >
                    <Search className="absolute left-6 h-5 w-5 text-gray-400" />
                    <Input
                      type="text"
                      placeholder="Enter Order ID to track..."
                      value={orderSearchId}
                      onChange={(e) => setOrderSearchId(e.target.value)}
                      className="h-14 pl-14 pr-4 bg-transparent border-0 focus-visible:ring-0 text-gray-900 placeholder:text-gray-400 text-base md:text-lg w-full"
                    />
                    <Button
                      type="submit"
                      className="h-12 px-6 md:px-8 bg-food-primary hover:bg-orange-600 text-white font-bold rounded-xl shadow-md transition-all"
                    >
                      Track
                    </Button>
                  </form>
                </div>
              </div>

              {/* Right Image Layout */}
              <div className="w-full lg:w-1/2 relative">
                {/* Decorative background shape */}
                <div className="absolute inset-0 bg-orange-50 rounded-[3rem] transform rotate-3 scale-105 -z-10" />
                
                <div className="relative rounded-[2.5rem] overflow-hidden shadow-2xl shadow-gray-300/50 aspect-[4/3] group">
                  <img 
                    src="https://images.unsplash.com/photo-1517248135467-4c7edcad34c4?auto=format&fit=crop&w=1200&q=80" 
                    alt="Beautiful modern restaurant interior" 
                    className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-1000"
                  />
                  {/* Subtle overlay for contrast */}
                  <div className="absolute inset-0 bg-gradient-to-tr from-gray-900/20 to-transparent" />
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* ═══════════════════════════════════════════════════════════
            THEMATIC FEATURES SECTION — Image Cards
         ═══════════════════════════════════════════════════════════ */}
        <section className="py-20 lg:py-28 bg-gray-50 border-t border-gray-100">
          <div className="container mx-auto px-4 max-w-7xl">
            <div className="text-center mb-16">
              <h2 className="text-3xl md:text-5xl font-extrabold tracking-tight text-gray-900 mb-4">
                Everything you need. <br className="md:hidden" /><span className="text-gray-400">Nothing you don't.</span>
              </h2>
              <p className="text-lg text-gray-600 max-w-2xl mx-auto">
                A seamless dining experience powered by intuitive technology.
              </p>
            </div>
            
            {/* Consistent 3-Card Layout with Background Images */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8">

              {/* Card 1 */}
              <div className="relative rounded-3xl overflow-hidden group aspect-[4/5] md:aspect-auto md:min-h-[480px] shadow-lg">
                <img 
                  src="https://images.unsplash.com/photo-1600147131759-880e94a6185f?auto=format&fit=crop&w=800&q=80" 
                  alt="Fine dining table setup" 
                  className="absolute inset-0 w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
                />
                {/* Heavy dark gradient overlay for text readability */}
                <div className="absolute inset-0 bg-gradient-to-t from-gray-900 via-gray-900/70 to-gray-900/20" />
                
                <div className="absolute inset-0 p-8 md:p-10 flex flex-col justify-end text-white">
                  <div className="w-14 h-14 bg-food-primary rounded-2xl flex items-center justify-center mb-6 shadow-lg shadow-food-primary/30">
                    <Smartphone className="h-7 w-7 text-white" />
                  </div>
                  <h3 className="text-2xl lg:text-3xl font-bold mb-3 tracking-tight">Digital Menus,<br/>Reimagined.</h3>
                  <p className="text-gray-200 leading-relaxed text-sm lg:text-base">
                    Say goodbye to clunky PDFs. Browse interactive, beautifully designed menus with instant dietary filters and effortless customization.
                  </p>
                </div>
              </div>

              {/* Card 2 */}
              <div className="relative rounded-3xl overflow-hidden group aspect-[4/5] md:aspect-auto md:min-h-[480px] shadow-lg">
                <img 
                  src="https://images.unsplash.com/photo-1556742049-0cfed4f6a45d?auto=format&fit=crop&w=800&q=80" 
                  alt="Contactless payment" 
                  className="absolute inset-0 w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-gray-900 via-gray-900/70 to-gray-900/20" />
                
                <div className="absolute inset-0 p-8 md:p-10 flex flex-col justify-end text-white">
                  <div className="w-14 h-14 bg-food-primary rounded-2xl flex items-center justify-center mb-6 shadow-lg shadow-food-primary/30">
                    <CreditCard className="h-7 w-7 text-white" />
                  </div>
                  <h3 className="text-2xl lg:text-3xl font-bold mb-3 tracking-tight">Lightning Fast<br/>Checkout.</h3>
                  <p className="text-gray-200 leading-relaxed text-sm lg:text-base">
                    Integrated instantly with PayNow and all major cards. Secure, bank-level encryption with absolutely no account required.
                  </p>
                </div>
              </div>

              {/* Card 3 */}
              <div className="relative rounded-3xl overflow-hidden group aspect-[4/5] md:aspect-auto md:min-h-[480px] shadow-lg">
                <img 
                  src="https://images.unsplash.com/photo-1728044849280-10a1a75cff83?auto=format&fit=crop&w=800&q=80" 
                  alt="Chef cooking in kitchen" 
                  className="absolute inset-0 w-full h-full object-cover group-hover:scale-110 transition-transform duration-700"
                />
                <div className="absolute inset-0 bg-gradient-to-t from-gray-900 via-gray-900/70 to-gray-900/20" />
                
                <div className="absolute inset-0 p-8 md:p-10 flex flex-col justify-end text-white">
                  <div className="w-14 h-14 bg-food-primary rounded-2xl flex items-center justify-center mb-6 shadow-lg shadow-food-primary/30">
                    <Clock className="h-7 w-7 text-white" />
                  </div>
                  <h3 className="text-2xl lg:text-3xl font-bold mb-3 tracking-tight">Real-Time<br/>Tracking.</h3>
                  <p className="text-gray-200 leading-relaxed text-sm lg:text-base">
                    Watch your order move from the kitchen to your table with live milestone updates. Always know exactly when to pick up.
                  </p>
                </div>
              </div>

            </div>
          </div>
        </section>

        {/* ═══════════════════════════════════════════════════════════
            MERCHANT CTA — Matches Footer (bg-gray-800)
         ═══════════════════════════════════════════════════════════ */}
        <section className="relative bg-gray-800 py-24 lg:py-32 overflow-hidden border-b border-gray-700/50">
          {/* Subtle background image of restaurant owner */}
          <img 
            src="https://images.unsplash.com/photo-1572715376701-98568319fd0b?auto=format&fit=crop&w=2000&q=80" 
            alt="Restaurant environment"
            className="absolute inset-0 w-full h-full object-cover mix-blend-overlay opacity-10"
          />
          {/* Gradients to blend seamlessly into the footer */}
          <div className="absolute inset-0 bg-gradient-to-b from-gray-800/80 via-gray-800/60 to-gray-800" />
          
          <div className="container relative z-10 mx-auto px-4 max-w-4xl text-center">
            <div className="w-16 h-16 bg-gray-700/50 backdrop-blur-md rounded-2xl flex items-center justify-center mx-auto mb-8 ring-1 ring-white/10">
              <Store className="h-8 w-8 text-food-primary" />
            </div>
            
            <h2 className="text-4xl md:text-5xl lg:text-6xl font-extrabold tracking-tight text-white mb-6 leading-tight">
              Run your restaurant <br className="hidden sm:block"/>on FoodFlow.
            </h2>
            
            <p className="text-lg md:text-xl text-gray-300 mb-10 max-w-2xl mx-auto leading-relaxed">
              Take back control of your margins. Zero commissions, powerful insights, and a world-class ordering experience for your guests.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
              <Button
                asChild
                size="lg"
                className="h-14 px-8 bg-food-primary hover:bg-orange-600 text-white font-bold rounded-xl text-base shadow-lg shadow-food-primary/20 w-full sm:w-auto"
              >
                <a href="https://wa.me/6583715877" target="_blank" rel="noopener noreferrer">
                  <MessageSquare className="mr-2 h-5 w-5" />
                  Partner With Us
                </a>
              </Button>
              <Button
                asChild
                size="lg"
                className="h-14 px-8 bg-white hover:bg-gray-100 text-gray-900 font-bold rounded-xl text-base shadow-lg w-full sm:w-auto"
              >
                <a href="mailto:adil@codercrafters.com?subject=FoodFlow%20Restaurant%20Partnership">
                  <Mail className="mr-2 h-5 w-5 text-gray-600" />
                  Contact Sales
                </a>
              </Button>
            </div>
          </div>
        </section>

      </div>
    </Layout>
  );
};

export default Landing;
