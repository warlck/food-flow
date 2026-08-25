
import React from 'react';
import Header from './Header';
import Footer from './Footer';
import { type Surface } from '@/lib/surfaces';

interface LayoutProps {
  children: React.ReactNode;
  surface?: Surface;
}

const Layout: React.FC<LayoutProps> = ({ children, surface }) => {
  return (
    <div className="flex flex-col min-h-screen">
      <Header surface={surface} />
      <main className="flex-grow">
        {children}
      </main>
      <Footer />
    </div>
  );
};

export default Layout;
