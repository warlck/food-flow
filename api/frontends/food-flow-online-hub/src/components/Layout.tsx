
import React from 'react';
import Header from './Header';
import Footer from './Footer';
import { useSurface } from '@/hooks/useSurface';
import { type Surface } from '@/lib/surfaces';

interface LayoutProps {
  children: React.ReactNode;
  surface?: Surface;
}

const Layout: React.FC<LayoutProps> = ({ children, surface: propSurface }) => {
  const hookSurface = useSurface();
  const surface = propSurface ?? hookSurface;
  const isMarketing = surface === 'marketing';

  return (
    <div className={`flex flex-col min-h-screen ${isMarketing ? 'bg-ink-950 text-white' : 'bg-background text-foreground'}`}>
      <Header surface={surface} />
      <main className="flex-grow">
        {children}
      </main>
      <Footer />
    </div>
  );
};

export default Layout;
