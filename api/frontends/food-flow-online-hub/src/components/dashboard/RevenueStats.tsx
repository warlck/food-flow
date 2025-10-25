
import React from 'react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { DollarSign, CreditCard, TrendingUp } from 'lucide-react';

export const RevenueStats: React.FC = () => {
  return (
    <Card>
      <CardHeader className="pb-2">
        <CardTitle className="text-lg font-medium">
          <div className="flex items-center">
            <DollarSign className="mr-2 h-5 w-5 text-muted-foreground" />
            Revenue Overview
          </div>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div>
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-muted-foreground">Monthly Revenue</p>
              <p className="text-2xl font-bold">$12,548</p>
            </div>
            <div className="flex items-center text-green-500">
              <TrendingUp className="h-4 w-4 mr-1" />
              <span className="text-sm font-medium">+18.2%</span>
            </div>
          </div>
          
          <div className="mt-4 grid grid-cols-2 gap-4 border-t border-muted/60 pt-4">
            <div>
              <p className="text-sm text-muted-foreground">Average Order</p>
              <p className="text-xl font-bold">$24.35</p>
            </div>
            <div>
              <p className="text-sm text-muted-foreground">Orders This Month</p>
              <p className="text-xl font-bold">486</p>
            </div>
          </div>
          
          <div className="mt-4 space-y-2">
            <div className="flex items-center space-x-2">
              <CreditCard className="h-4 w-4 text-muted-foreground" />
              <div className="flex-1 space-y-1">
                <p className="text-sm font-medium leading-none">
                  Payment Methods
                </p>
                <div className="flex items-center text-xs text-muted-foreground">
                  <div className="flex items-center">
                    <span className="inline-block w-2 h-2 rounded-full bg-green-500 mr-1"></span>
                    <span>Card (68%)</span>
                  </div>
                  <span className="mx-2">•</span>
                  <div className="flex items-center">
                    <span className="inline-block w-2 h-2 rounded-full bg-blue-500 mr-1"></span>
                    <span>Cash (32%)</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  );
};
