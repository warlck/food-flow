import React, { useEffect, useState } from 'react';
import {
  PaymentElement,
  useStripe,
  useElements,
} from '@stripe/react-stripe-js';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';
import { Loader2 } from 'lucide-react';

interface StripePaymentFormProps {
  orderId: string;
  total: number;
  onSuccess: () => void;
  onError: (error: Error) => void;
}

const StripePaymentForm: React.FC<StripePaymentFormProps> = ({
  orderId,
  total,
  onSuccess,
  onError,
}) => {
  const stripe = useStripe();
  const elements = useElements();
  const [isProcessing, setIsProcessing] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);
  const [isStripeReady, setIsStripeReady] = useState(false);
  const [showStillLoadingHint, setShowStillLoadingHint] = useState(false);

  useEffect(() => {
    if (stripe && elements) {
      setIsStripeReady(true);
      return;
    }

    setIsStripeReady(false);

    // If Stripe stays uninitialized for a while, show a helpful hint.
    const t = window.setTimeout(() => setShowStillLoadingHint(true), 6000);
    return () => window.clearTimeout(t);
  }, [stripe, elements]);

  useEffect(() => {
    const onWindowError = (event: ErrorEvent) => {
      const msg = event?.error?.message || event?.message || '';
      if (msg.includes("Failed to execute 'observe' on 'MutationObserver'")) {
        setErrorMessage(
          'A browser extension is breaking Stripe Elements on this page (MutationObserver error). Disable extensions for this site (ad blockers/privacy tools/Dark Reader/password managers/Grammarly) or use an incognito window.'
        );
      }
    };

    window.addEventListener('error', onWindowError);
    return () => window.removeEventListener('error', onWindowError);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!stripe || !elements) {
      // Stripe.js hasn't loaded yet
      return;
    }

    setIsProcessing(true);
    setErrorMessage(null);

    try {
      const { error, paymentIntent } = await stripe.confirmPayment({
        elements,
        confirmParams: {
          return_url: `${window.location.origin}/order-confirmation/${orderId}`,
        },
        redirect: 'if_required',
      });

      if (error) {
        // Show error to customer
        setErrorMessage(error.message || 'An error occurred during payment');
        toast.error(error.message || 'Payment failed');
        onError(new Error(error.message));
      } else if (paymentIntent && paymentIntent.status === 'succeeded') {
        // Payment succeeded
        toast.success('Payment successful!');
        onSuccess();
      } else if (paymentIntent && paymentIntent.status === 'processing') {
        // Payment is still processing
        toast.info('Payment is processing...');
        onSuccess();
      } else {
        // Payment requires additional action (like 3D Secure)
        // This is handled by stripe.confirmPayment redirect
        toast.info('Payment requires additional verification');
      }
    } catch (err) {
      const error = err as Error;
      setErrorMessage(error.message || 'An unexpected error occurred');
      toast.error('Payment failed. Please try again.');
      onError(error);
    } finally {
      setIsProcessing(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {!isStripeReady && (
        <div className="p-3 bg-gray-50 border border-gray-200 rounded-lg text-sm text-gray-700">
          <div className="flex items-center">
            <Loader2 className="w-4 h-4 mr-2 animate-spin" />
            Loading Stripe payment form...
          </div>
          {showStillLoadingHint && (
            <div className="mt-2 text-xs text-gray-600">
              Still loading? Check if an ad blocker/privacy extension is blocking Stripe, and confirm
              that <code className="px-1">js.stripe.com</code> is reachable.
            </div>
          )}
        </div>
      )}

      <div className="p-4 border-2 rounded-xl bg-white">
        <PaymentElement
          options={{
            layout: 'tabs',
          }}
        />
      </div>

      {errorMessage && (
        <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-600 text-sm">
          {errorMessage}
        </div>
      )}
      
      <Button
        type="submit"
        disabled={!stripe || isProcessing}
        className="w-full bg-gradient-to-r from-food-primary to-food-accent hover:from-food-accent hover:to-food-primary text-white py-4 rounded-xl font-semibold text-lg disabled:opacity-50"
      >
        {isProcessing ? (
          <span className="flex items-center justify-center">
            <Loader2 className="w-5 h-5 mr-2 animate-spin" />
            Processing...
          </span>
        ) : (
          `Pay $${total.toFixed(2)}`
        )}
      </Button>
      
      <p className="text-xs text-gray-500 text-center">
        Your payment is secured by Stripe. We never store your card details.
      </p>
    </form>
  );
};

export default StripePaymentForm;
