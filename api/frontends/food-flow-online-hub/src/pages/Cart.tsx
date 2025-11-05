import React from 'react';
import { useIsMobile } from '@/hooks/use-mobile';
import CartDesktop from './CartDesktop';
import CartMobile from './CartMobile';

const Cart: React.FC = () => {
  const isMobile = useIsMobile();
  
  return isMobile ? <CartMobile /> : <CartDesktop />;
};

export default Cart;
