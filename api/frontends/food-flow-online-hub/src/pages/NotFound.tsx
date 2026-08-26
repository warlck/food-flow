import React, { useEffect } from 'react';
import { useLocation, Link } from 'react-router-dom';
import MarketingPage from '@/components/layout/MarketingPage';
import { Button } from '@/components/ui/button';
import { FileQuestion, ArrowLeft, Package } from 'lucide-react';

const NotFound: React.FC = () => {
  const location = useLocation();

  useEffect(() => {
    console.error(
      '404 Error: User attempted to access non-existent route:',
      location.pathname
    );
  }, [location.pathname]);

  return (
    <MarketingPage
      eyebrow="404 Error"
      title="Page not found"
      description="The page you are looking for doesn't exist, has been removed, or is temporarily unavailable."
      icon={FileQuestion}
      maxWidth="max-w-2xl"
    >
      <div className="glass-panel ring-lit rounded-2xl p-8 border-white/10 text-center space-y-6">
        <p className="text-gray-300 text-base">
          Looking for a food delivery or an active order? Use the quick actions below:
        </p>

        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Button
            asChild
            className="w-full sm:w-auto bg-food-primary hover:bg-food-accent text-white font-semibold shadow-lg shadow-food-primary/25 px-6"
          >
            <Link to="/">
              <ArrowLeft className="mr-2 h-4 w-4" />
              Back to home
            </Link>
          </Button>

          <Button
            asChild
            variant="outline"
            className="w-full sm:w-auto border-white/20 bg-white/5 text-white hover:bg-white/15 px-6"
          >
            <Link to="/track-order">
              <Package className="mr-2 h-4 w-4" />
              Track an order
            </Link>
          </Button>
        </div>
      </div>
    </MarketingPage>
  );
};

export default NotFound;
